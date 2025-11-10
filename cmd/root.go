package cmd

import (
	"QuickSnip/ui"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
