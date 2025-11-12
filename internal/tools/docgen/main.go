package main

import (
	"QuickSnip/cmd"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"
)

func main() {
	out := flag.String("out", "./Readme.md", "output file (for single readme)")
	dir := flag.String("dir", "./docs", "output directory for per-command docs")
	front := flag.Bool("frontmatter", false, "prepend YAML front matter to markdown")
	flag.Parse()

	root := cmd.Root()
	root.DisableAutoGenTag = true

	// Ensure directory exists for per-command docs
	if err := os.MkdirAll(*dir, 0o755); err != nil {
		log.Fatal(err)
	}

	// 1. Generate the per-command Markdown files first
	err := doc.GenMarkdownTreeCustom(root, *dir,
		func(filename string) string {
			if *front {
				base := filepath.Base(filename)
				name := strings.TrimSuffix(base, filepath.Ext(base))
				title := strings.ReplaceAll(name, "_", " ")
				return fmt.Sprintf("---\ntitle: %q\nslug: %q\ndescription: \"CLI reference for %s\"\n---\n\n",
					title, name, title)
			}
			return ""
		},
		func(name string) string {
			// Customize the link format: from "snip_add.md" to "docs/snip_add.md"
			return fmt.Sprintf("docs/%s", strings.ToLower(name))
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Generate the main README.md
	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var buf strings.Builder
	if err := doc.GenMarkdown(root, &buf); err != nil {
		log.Fatal(err)
	}

	content := buf.String()

	if *front {
		title := strings.TrimSuffix(filepath.Base(*out), filepath.Ext(*out))
		title = strings.ReplaceAll(title, "_", " ")
		f.WriteString(fmt.Sprintf("---\ntitle: %q\ndescription: \"CLI reference for %s\"\n---\n\n",
			title, title))
	}

	if _, err := f.WriteString(content); err != nil {
		log.Fatal(err)
	}

	log.Printf("README generated at %s\nCommand docs generated in %s\n", *out, *dir)
}
