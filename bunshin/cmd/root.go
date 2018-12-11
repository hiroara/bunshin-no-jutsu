package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var srcDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bunshin",
	Short: "Sync files to another location",
	Long: `Bunshin is a simple tooljto sync files to anoother location.
This can be used for buckup files to a disk or a directory which is watched by Dropbox or other sync tools.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		destDir := viper.GetString("dest")
		dryrun, err := flags.GetBool("dry-run")
		if err != nil {
			log.Fatal(err)
		}
		del, err := flags.GetBool("delete")
		if err != nil {
			log.Fatal(err)
		}
		if err := runSync(srcDir, destDir, dryrun, del); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "run without actual changes.")
	rootCmd.PersistentFlags().BoolP("delete", "d", false, "delete non-existent files in destination.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	dir, err := findConfig()
	if err != nil {
		log.Fatal(err)
	}

	srcDir = dir

	// Use config file from the flag.
	viper.SetConfigFile(filepath.Join(dir, configFileName))

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
