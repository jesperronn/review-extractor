package main

import (
	"fmt"
	"os"

	"github.com/jesper/review-extractor/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "review-extractor",
		Short: "Extract code reviews from various platforms",
		Long:  `A tool to extract code reviews from various platforms like GitHub, GitLab, etc.`,
	}

	// Add commands
	rootCmd.AddCommand(cmd.NewExtractCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
