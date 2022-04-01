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
	"os"

	"github.com/pepol/databuddy/internal/log"
	"github.com/pepol/databuddy/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configEnvPrefix = "APP"
	defaultDataDir  = "/var/lib/databuddy"
	defaultPort     = 6543
	defaultHost     = "127.0.0.1"
	defaultLogLevel = "debug"
)

var rootCmd = &cobra.Command{
	Use:     "databuddy",
	Short:   "DataBuddy Global Datastore",
	Long:    `Service that handles API requests for databuddy storage model`,
	Run:     serve,
	Version: version,
}

// Build info variables set by goreleaser.
var version = "latest"

func init() {
	viper.SetDefault("datadir", defaultDataDir)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("host", defaultHost)
	viper.SetDefault("loglevel", defaultLogLevel)

	// Parse environment variables.
	viper.SetEnvPrefix(configEnvPrefix)
	viper.AutomaticEnv()

	// Parse commandline arguments.

	// Data storage settings - global.
	rootCmd.PersistentFlags().StringP("datadir", "d", defaultDataDir, "directory where all data is stored")
	if err := viper.BindPFlag("datadir", rootCmd.PersistentFlags().Lookup("datadir")); err != nil {
		log.Fatal(err)
	}

	// RESP server settings.
	rootCmd.Flags().IntP("port", "p", defaultPort, "port to listen on")
	if err := viper.BindPFlag("port", rootCmd.Flags().Lookup("port")); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().StringP("host", "H", defaultHost, "host to listen on")
	if err := viper.BindPFlag("host", rootCmd.Flags().Lookup("host")); err != nil {
		log.Fatal(err)
	}

	// Observability settings.
	rootCmd.Flags().String("loglevel", defaultLogLevel, "level of logs to display")
	if err := viper.BindPFlag("loglevel", rootCmd.Flags().Lookup("loglevel")); err != nil {
		log.Fatal(err)
	}
}

// Serve HTTP requests.
func serve(cmd *cobra.Command, args []string) {
	server.Serve(version)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
