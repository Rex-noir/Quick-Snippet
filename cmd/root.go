package cmd

import (
	"QuickSnip/db"
	"QuickSnip/ui"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "snip",
	Short: "Quick Snip is a tool to save your thought snippets",
	Long:  `A fast and flexible cli tool to save your thought snippets and read them again`,
	RunE: func(cmd *cobra.Command, args []string) error {
		snippets := []ui.Snippet{
			{ID: 1, Title: "Golang Tips", Body: "Use defer for cleanup"},
			{ID: 2, Title: "Docker", Body: "docker ps -a"},
			{ID: 3, Title: "Java Snip", Body: "java snip"},
		}
		return ui.RunBrowse(snippets)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.snip.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")

	err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		return
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return
		}
		appDir := filepath.Join(configDir, "quicksnip")
		if err := os.MkdirAll(appDir, 0755); err != nil {
			return
		}
		dbPath := db.GetDBPath(appDir)
		err = db.RunMigrations(dbPath)
		if err != nil {
			return
		}

		path := configDir + "/snip"
		viper.SetConfigName("config")
		viper.AddConfigPath(path)

		configPath := filepath.Join(appDir, "config.yaml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			err := viper.SafeWriteConfigAs(configPath)
			if err != nil {
				return
			}
		}
	}
}
