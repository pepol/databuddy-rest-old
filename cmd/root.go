// Package cmd contains the cobra and viper related CLI implementations.
package cmd

import (
	"os"

	"github.com/pepol/databuddy/internal/log"
	"github.com/pepol/databuddy/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configEnvPrefix = "APP"
	defaultDevel    = false
	defaultDataDir  = "/var/lib/databuddy"
	defaultPort     = 6543
	defaultHost     = "127.0.0.1"
	defaultLogLevel = "info"
	defaultSerfPort = 6544
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
	viper.SetDefault("devel", defaultDevel)
	viper.SetDefault("datadir", defaultDataDir)
	viper.SetDefault("port", defaultPort)
	viper.SetDefault("host", defaultHost)
	viper.SetDefault("loglevel", defaultLogLevel)
	viper.SetDefault("join", []string{})
	viper.SetDefault("serfport", defaultSerfPort)

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
	rootCmd.Flags().BoolP("devel", "D", defaultDevel, "enable development mode")
	if err := viper.BindPFlag("devel", rootCmd.Flags().Lookup("devel")); err != nil {
		log.Fatal(err)
	}

	// Serf settings.
	rootCmd.Flags().Int("serfport", defaultSerfPort, "port which Serf/Gossip traffic will use")
	if err := viper.BindPFlag("serfport", rootCmd.Flags().Lookup("serfport")); err != nil {
		log.Fatal(err)
	}

	rootCmd.Flags().StringSliceP("join", "j", []string{}, "comma-separated list of connection strings (address and port) to connect to for cluster bootstrap")
	if err := viper.BindPFlag("join", rootCmd.Flags().Lookup("join")); err != nil {
		log.Fatal(err)
	}
}

// Serve HTTP requests.
func serve(cmd *cobra.Command, args []string) {
	server.Serve(version)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
