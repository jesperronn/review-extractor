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
		Short: "Extract code review comments from Git platforms",
	}

	rootCmd.AddCommand(cmd.ExtractCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
} 