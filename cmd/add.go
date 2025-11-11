package cmd

import (
	"QuickSnip/db"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add <title> <body>",
	Short: "Add a new snippet",
	Long:  `Directly add a new snippet to the database.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		body := args[1]

		// Example: You can call your own save function here
		if title == "" || body == "" {
			return fmt.Errorf("title and body cannot be empty")
		}

		appDir := viper.GetString("app_dir")
		conn, err := db.Open(appDir)
		if err != nil {
			return err
		}

		_, err = db.CreateSnippet(conn, title, body)
		if err != nil {
			return err
		}

		fmt.Printf("Added snippet: %q -> %q\n", title, body)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
