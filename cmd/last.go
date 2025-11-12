package cmd

import (
	"QuickSnip/db"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var lastCmd = &cobra.Command{
	Use:   "last [number]",
	Short: "Print the last snippet",
	Long:  `Print the last snippet saved in the database.`,
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
				return fmt.Errorf("invalid number: %w", convErr)
			}
			number = n
		}

		snippets, err := db.FetchLatest(dbConn, &number)
		if err != nil {
			return err
		}

		for _, s := range snippets {
			fmt.Printf("[%d] %s\n%s\n\n", s.ID, s.Title, s.Body)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lastCmd)
}
