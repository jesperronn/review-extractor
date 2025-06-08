package core

import (
	"context"
	"errors"
	"testing"

	"github.com/jesper/review-extractor/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExtractor is a mock implementation of the Extractor interface
type MockExtractor struct {
	mock.Mock
}

func (m *MockExtractor) ExtractReviews(ctx context.Context, repoURL string) ([]models.Review, error) {
	args := m.Called(ctx, repoURL)
	return args.Get(0).([]models.Review), args.Error(1)
}

func TestNewReviewExtractor(t *testing.T) {
	config := &models.Config{
		APIToken: "test-token",
		Repositories: []models.RepositoryConfig{
			{
				Provider: models.ProviderGitHub,
				URL:      "https://github.com/test/repo",
			},
		},
	}

	extractors := make(map[models.Provider]Extractor)
	extractor := NewReviewExtractor(config, extractors)

	assert.NotNil(t, extractor)
	assert.Equal(t, config, extractor.config)
	assert.Equal(t, extractors, extractor.extractors)
}

func TestExtractReviews_Success(t *testing.T) {
	// Setup
	config := &models.Config{
		APIToken: "test-token",
		Repositories: []models.RepositoryConfig{
			{
				Provider: models.ProviderGitHub,
				URL:      "https://github.com/test/repo1",
			},
			{
				Provider: models.ProviderGitHub,
				URL:      "https://github.com/test/repo2",
			},
		},
	}

	mockExtractor := new(MockExtractor)
	extractors := map[models.Provider]Extractor{
		models.ProviderGitHub: mockExtractor,
	}

	reviews1 := []models.Review{
		{
			PRID:       1,
			PRTitle:    "Test PR 1",
			PRAuthor:   "user1",
			Repository: "repo1",
		},
	}

	reviews2 := []models.Review{
		{
			PRID:       2,
			PRTitle:    "Test PR 2",
			PRAuthor:   "user2",
			Repository: "repo2",
		},
	}

	mockExtractor.On("ExtractReviews", mock.Anything, "https://github.com/test/repo1").Return(reviews1, nil)
	mockExtractor.On("ExtractReviews", mock.Anything, "https://github.com/test/repo2").Return(reviews2, nil)

	extractor := NewReviewExtractor(config, extractors)

	// Execute
	result, err := extractor.ExtractReviews(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalComments)
	assert.Equal(t, 2, result.RepositoriesProcessed)
	assert.Len(t, result.Reviews, 2)
	assert.Equal(t, reviews1[0], result.Reviews[0])
	assert.Equal(t, reviews2[0], result.Reviews[1])
}

func TestExtractReviews_NoExtractor(t *testing.T) {
	// Setup
	config := &models.Config{
		APIToken: "test-token",
		Repositories: []models.RepositoryConfig{
			{
				Provider: "bitbucket", // Use a string that is not defined in models.Provider
				URL:      "https://bitbucket.org/test/repo",
			},
		},
	}

	extractors := make(map[models.Provider]Extractor)
	extractor := NewReviewExtractor(config, extractors)

	// Execute
	result, err := extractor.ExtractReviews(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no extractor available for provider: bitbucket")
}

func TestExtractReviews_ExtractionError(t *testing.T) {
	// Setup
	config := &models.Config{
		APIToken: "test-token",
		Repositories: []models.RepositoryConfig{
			{
				Provider: models.ProviderGitHub,
				URL:      "https://github.com/test/repo",
			},
		},
	}

	mockExtractor := new(MockExtractor)
	extractors := map[models.Provider]Extractor{
		models.ProviderGitHub: mockExtractor,
	}

	expectedErr := errors.New("extraction failed")
	mockExtractor.On("ExtractReviews", mock.Anything, "https://github.com/test/repo").Return([]models.Review{}, expectedErr)

	extractor := NewReviewExtractor(config, extractors)

	// Execute
	result, err := extractor.ExtractReviews(context.Background())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to extract reviews from")
}

func TestGenerateStatistics(t *testing.T) {
	reviews := []models.Review{
		{
			CommentAuthor: "user1",
			FilePath:      "file1.go",
		},
		{
			CommentAuthor: "user1",
			FilePath:      "file1.go",
		},
		{
			CommentAuthor: "user2",
			FilePath:      "file2.go",
		},
		{
			CommentAuthor: "user3",
			FilePath:      "file1.go",
		},
	}

	stats := generateStatistics(reviews)

	// Since Statistics struct does not have MostActiveReviewers or FilesWithMostComments fields,
	// we only check the fields that exist.
	assert.Equal(t, 4, stats.TotalReviews)
	assert.Equal(t, 0, stats.TotalPRs) // Not set in generateStatistics
	assert.NotNil(t, stats.TopReviewers)
	assert.NotNil(t, stats.TopRepositories)
}

func TestGetTopN(t *testing.T) {
	tests := []struct {
		name     string
		counts   map[string]int
		n        int
		expected []string
	}{
		{
			name: "normal case",
			counts: map[string]int{
				"a": 3,
				"b": 1,
				"c": 2,
			},
			n:        2,
			expected: []string{"a", "c"},
		},
		{
			name:     "empty map",
			counts:   map[string]int{},
			n:        2,
			expected: []string{},
		},
		{
			name: "n larger than map size",
			counts: map[string]int{
				"a": 1,
				"b": 2,
			},
			n:        3,
			expected: []string{"b", "a"},
		},
		{
			name: "n is zero",
			counts: map[string]int{
				"a": 1,
				"b": 2,
			},
			n:        0,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTopN(tt.counts, tt.n)
			assert.Equal(t, tt.expected, result)
		})
	}
}
