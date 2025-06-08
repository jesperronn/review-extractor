package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jesper/review-extractor/internal/adapters/github"
	"github.com/jesper/review-extractor/internal/core"
	"github.com/jesper/review-extractor/pkg/models"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configFile string
	outputFile string
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract code review comments from Git platforms",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		config, err := loadConfig(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override output file if specified
		if outputFile != "" {
			config.OutputFile = outputFile
		}

		// Create extractors
		extractors := make(map[models.Provider]core.Extractor)
		
		// Add GitHub extractor
		extractors[models.ProviderGitHub] = github.NewExtractor(config.APIToken)

		// Create main extractor
		extractor := core.NewReviewExtractor(config, extractors)

		// Extract reviews
		ctx := context.Background()
		result, err := extractor.ExtractReviews(ctx)
		if err != nil {
			return fmt.Errorf("failed to extract reviews: %w", err)
		}

		// Write output
		if err := writeOutput(result, config.OutputFile); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}

		fmt.Printf("Successfully extracted %d reviews from %d repositories\n",
			result.TotalComments, result.RepositoriesProcessed)
		return nil
	},
}

func init() {
	extractCmd.Flags().StringVar(&configFile, "config", "", "Path to configuration file")
	extractCmd.Flags().StringVar(&outputFile, "output", "", "Path to output file (overrides config)")
	extractCmd.MarkFlagRequired("config")
}

// loadConfig loads the configuration from a YAML file
func loadConfig(path string) (*models.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// writeOutput writes the extraction result to a JSON file
func writeOutput(result *models.ExtractionResult, path string) error {
	// Create output directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
} 