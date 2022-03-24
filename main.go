/*
Copyright © 2022 Peter Polacik

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/ansrivas/fiberprometheus/v2"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	expvarmw "github.com/gofiber/fiber/v2/middleware/expvar"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	api "github.com/pepol/databuddy/api/v1alpha3"
	_ "github.com/pepol/databuddy/docs"
	"github.com/pepol/databuddy/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var rootCmd = &cobra.Command{
	Use:   "databuddy",
	Short: "DataBuddy Global Datastore",
	Long:  `Service that handles API requests for databuddy storage model`,
	Run:   serve,
}

var cfgFile string

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".orders-api" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".orders-api")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.databuddy.yaml)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//go:generate swag init
// @title        DataBuddy
// @version      1.0.0-alpha3
// @description  API to use DataBuddy data storage system

// @contact.name  Peter Polacik

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1alpha3
// Serve HTTP requests.
func serve(cmd *cobra.Command, args []string) {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error("error shutting down tracer provider", err)
		}
	}()

	app := fiber.New()

	prometheus := fiberprometheus.New("databuddy")
	prometheus.RegisterAt(app, "/_internal/metrics")

	app.Use("/", prometheus.Middleware)
	app.Use(otelfiber.Middleware("databuddy"))
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(expvarmw.New())
	app.Use(etag.New())

	app.Get("/_internal/dashboard", monitor.New())
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API definition.
	v1alpha3 := app.Group("/v1alpha3")

	v1alpha3.Use("/", prometheus.Middleware)

	// API controller.
	apiCtl := api.NewController(&api.Config{})

	apiCtl.Route(v1alpha3)

	if err := app.Listen(":8080"); err != nil {
		log.Error("error listening", err)
	}
}

func initTracer() *sdktrace.TracerProvider {
	exporter, err := stdout.New()
	if err != nil {
		log.Fatal(err)
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatal("cannot retrieve build info")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("databuddy"),
				semconv.ServiceVersionKey.String(info.Main.Version),
			),
		),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
