package cmd

import (
	"QuickSnip/db"
	"QuickSnip/ui"
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var interactive bool

var addCmd = &cobra.Command{
	Use:   "add [title] [body]",
	Short: "Add a new snippet",
	Long:  `Add a new snippet directly or interactively if -i flag is provided.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var title, body string
		appDir := viper.GetString("data_dir")
		conn, err := db.Open(appDir)
		if err != nil {
			return err
		}

		defer func(conn *sql.DB) {
			err := conn.Close()
			if err != nil {
				fmt.Println("Error closing db connection", err)
			}
		}(conn)

		if len(args) > 0 {
			title = args[0]
		}
		if len(args) > 1 {
			body = args[1]
		}

		if interactive {
			return ui.RunAddInteractive(conn, &title, &body)
		}

		if title == "" {
			return fmt.Errorf("title cannot be empty")
		}

		_, err = db.CreateSnippet(conn, title, body)
		if err != nil {
			return err
		}

		if body == "" {
			fmt.Printf("Added snippet: %q\n", title)
			return nil
		}

		fmt.Printf("✅ Added snippet: %q -> %q\n", title, body)
		return nil
	},
}

func init() {
	addCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run in interactive mode to enter snippet details")
	rootCmd.AddCommand(addCmd)
}
