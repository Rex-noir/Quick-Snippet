package cmd

import (
	"QuickSnip/db"
	"QuickSnip/mapper"
	"QuickSnip/ui"
	"bytes"
	"errors"
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
		appDir := viper.GetString("data_dir")
		dbConn, err := db.Open(appDir)
		if err != nil {
			return err
		}

		err = db.RunMigrations(db.GetDBPath(appDir))
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
		defer db.Close(dbConn)

		snippets, err := db.FetchSnippets(dbConn)
		if err != nil {
			return err
		}
		return ui.RunBrowse(dbConn, mapper.ToUISnippets(snippets))
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/snip/config.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")

	err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	if err != nil {
		return
	}
}

func initConfig() {
	var configPath string
	if cfgFile != "" {
		configPath = cfgFile
		viper.SetConfigFile(configPath)
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("Failed to get user config dir:", err)
			os.Exit(1)
		}
		appDir := filepath.Join(configDir, "snip")
		_ = os.MkdirAll(appDir, 0755)
		configPath = filepath.Join(appDir, "config.yaml")
		homeDir, _ := os.UserHomeDir()
		viper.SetConfigName("config")
		viper.AddConfigPath(appDir)
		viper.AddConfigPath(configDir)

		viper.AddConfigPath(homeDir)
		viper.SetConfigType("yaml")
	}
	err := viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			fmt.Println("No config file found, creating one...")
			homeDir, _ := os.UserHomeDir()

			defaultConfig := []byte(fmt.Sprintf(`data_dir: %s/.config/snip
debug: false`, homeDir))

			err := viper.ReadConfig(bytes.NewBuffer(defaultConfig))

			if err != nil {
				fmt.Println("Failed to create default config:", err)
				os.Exit(1)
			}
			err = viper.WriteConfigAs(configPath)
			if err != nil {
				fmt.Println("Failed to write default config:", err)
				os.Exit(1)
			}

		} else {
			fmt.Println("Error reading config file:", err)
			os.Exit(1)
		}

		fmt.Println("Using config file:", viper.ConfigFileUsed())

	}
}

func Root() *cobra.Command {
	return rootCmd
}
