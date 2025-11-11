package cmd

import (
	"QuickSnip/db"
	"QuickSnip/mapper"
	"QuickSnip/ui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse your saved snippets",
	RunE: func(cmd *cobra.Command, args []string) error {
		appDir := viper.GetString("app_dir")
		dbConn, err := db.Open(appDir)
		if err != nil {
			return err
		}
		snippets, err := db.FetchSnippets(dbConn)
		if err != nil {
			return err
		}
		return ui.RunBrowse(dbConn, mapper.ToUISnippets(snippets))
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
