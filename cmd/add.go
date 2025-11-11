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
	Args: func(cmd *cobra.Command, args []string) error {
		// If interactive mode is on, skip strict arg validation
		if interactive {
			return nil
		}
		if len(args) != 2 {
			return fmt.Errorf("requires exactly 2 arguments unless --interactive is used")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var title, body string
		appDir := viper.GetString("app_dir")
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

		if title == "" || body == "" {
			return fmt.Errorf("title and body cannot be empty")
		}

		_, err = db.CreateSnippet(conn, title, body)
		if err != nil {
			return err
		}

		fmt.Printf("✅ Added snippet: %q -> %q\n", title, body)
		return nil
	},
}

func init() {
	addCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run in interactive mode to enter snippet details")
	rootCmd.AddCommand(addCmd)
}
