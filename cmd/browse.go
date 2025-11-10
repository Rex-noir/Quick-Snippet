package cmd

import (
	"QuickSnip/ui"

	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse your saved snippets",
	RunE: func(cmd *cobra.Command, args []string) error {
		snippets := []ui.Snippet{
			{ID: 1, Title: "Golang Tips", Body: "Use defer for cleanup"},
			{ID: 2, Title: "Docker", Body: "docker ps -a"},
			{ID: 3, Title: "Java Snip", Body: "java snip"},
		}
		return ui.RunBrowse(snippets)
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
