package cmd

import (
	"QuickSnip/db"
	"QuickSnip/mapper"
	"QuickSnip/ui"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var firstCmd = &cobra.Command{
	Use:   "first [number]",
	Short: "Print the first snippet(s)",
	Long:  `Print the oldest snippet(s) from the database. You can specify a number to fetch multiple.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		appDir := viper.GetString("app_dir")
		dbConn, err := db.Open(appDir)
		if err != nil {
			return err
		}
		number := 1
		if len(args) > 0 {
			n, convErr := strconv.Atoi(args[0])
			if convErr != nil {
				return fmt.Errorf("invalid number: %v", convErr)
			}
			number = n
		}

		snippets, err := db.FetchFirst(dbConn, number)
		if err != nil {
			return err
		}
		err = ui.RunListModel(dbConn, mapper.ToUISnippets(snippets))
		if err != nil {
			return err
		}
		defer db.Close(dbConn)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(firstCmd)
}
