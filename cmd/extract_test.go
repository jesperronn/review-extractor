package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jesper/review-extractor/pkg/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoadConfig(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	config := &models.Config{
		APIToken:   "test-token",
		OutputFile: "test-output.json",
		Repositories: []models.RepositoryConfig{
			{
				Provider: models.ProviderGitHub,
				URL:      "https://github.com/test/repo",
			},
		},
	}

	data, err := yaml.Marshal(config)
	assert.NoError(t, err)

	err = os.WriteFile(configPath, data, 0644)
	assert.NoError(t, err)

	// Test loading config
	loadedConfig, err := loadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, config.APIToken, loadedConfig.APIToken)
	assert.Equal(t, config.OutputFile, loadedConfig.OutputFile)
	assert.Len(t, loadedConfig.Repositories, 1)
	assert.Equal(t, config.Repositories[0], loadedConfig.Repositories[0])
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := loadConfig("nonexistent.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create temporary config file with invalid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid-config.yaml")

	err := os.WriteFile(configPath, []byte("invalid: yaml: content: [}"), 0644)
	assert.NoError(t, err)

	_, err = loadConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config file")
}

func TestWriteOutput(t *testing.T) {
	// Create temporary output directory
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-output.json")

	result := &models.ExtractionResult{
		TotalComments:         2,
		RepositoriesProcessed: 1,
		Reviews: []models.Review{
			{
				PRID:    1,
				PRTitle: "Test PR",
			},
		},
	}

	// Test writing output
	err := writeOutput(result, outputPath)
	assert.NoError(t, err)

	// Verify file exists and contains expected content
	data, err := os.ReadFile(outputPath)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "Test PR")
	assert.Contains(t, string(data), "total_comments")
}

func TestWriteOutput_CreateDirectory(t *testing.T) {
	// Create temporary base directory
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "subdir", "test-output.json")

	result := &models.ExtractionResult{
		TotalComments:         1,
		RepositoriesProcessed: 1,
	}

	// Test writing output with directory creation
	err := writeOutput(result, outputPath)
	assert.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)
}

func TestWriteOutput_InvalidPath(t *testing.T) {
	result := &models.ExtractionResult{
		TotalComments:         1,
		RepositoriesProcessed: 1,
	}

	// Test writing to invalid path
	err := writeOutput(result, "/nonexistent/path/output.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create output directory")
}
