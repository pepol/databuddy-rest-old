package cmd

import (
	"os"

	"github.com/pepol/databuddy/internal/db"
	"github.com/pepol/databuddy/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the command to initialize the database.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Long:  `Initialize the database by creating all the required structures.`,
	Run:   initializeDatabase,
}

func init() {
	rootCmd.AddCommand(initCmd)

	viper.SetDefault("bucket", db.DefaultBucketName)

	// Parse environment variables.
	viper.SetEnvPrefix(configEnvPrefix)
	viper.AutomaticEnv()

	// Parse commandline arguments.

	// Default bucket name.
	initCmd.Flags().String("bucket", db.DefaultBucketName, "name of the default bucket")
	if err := viper.BindPFlag("bucket", initCmd.Flags().Lookup("bucket")); err != nil {
		log.Fatal(err)
	}
}

func initializeDatabase(cmd *cobra.Command, args []string) {
	datadir := viper.GetString("datadir")
	bucket := viper.GetString("bucket")

	if err := db.InitDatabase(datadir, bucket); err != nil {
		log.Error("initializing database", err)
		os.Exit(1)
	}
}
