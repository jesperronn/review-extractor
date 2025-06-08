package core

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/jesper/review-extractor/pkg/models"
)

// Extractor defines the interface for extracting reviews from a Git platform
type Extractor interface {
	ExtractReviews(ctx context.Context, repoURL string) ([]models.Review, error)
}

// ReviewExtractor orchestrates the extraction process across multiple repositories
type ReviewExtractor struct {
	extractors map[models.Provider]Extractor
	config     *models.Config
}

// NewReviewExtractor creates a new ReviewExtractor instance
func NewReviewExtractor(config *models.Config, extractors map[models.Provider]Extractor) *ReviewExtractor {
	return &ReviewExtractor{
		extractors: extractors,
		config:     config,
	}
}

// ExtractReviews processes all configured repositories and returns the combined results
func (re *ReviewExtractor) ExtractReviews(ctx context.Context) (*models.ExtractionResult, error) {
	var (
		wg          sync.WaitGroup
		mu          sync.Mutex
		allReviews  []models.Review
		errorChan   = make(chan error, len(re.config.Repositories))
		reviewsChan = make(chan []models.Review, len(re.config.Repositories))
	)

	// Process each repository concurrently
	for _, repo := range re.config.Repositories {
		wg.Add(1)
		go func(repo models.RepositoryConfig) {
			defer wg.Done()

			extractor, exists := re.extractors[repo.Provider]
			if !exists {
				errorChan <- fmt.Errorf("no extractor found for provider: %s", repo.Provider)
				return
			}

			reviews, err := extractor.ExtractReviews(ctx, repo.URL)
			if err != nil {
				errorChan <- fmt.Errorf("failed to extract reviews from %s: %w", repo.URL, err)
				return
			}

			reviewsChan <- reviews
		}(repo)
	}

	// Wait for all extractions to complete
	go func() {
		wg.Wait()
		close(errorChan)
		close(reviewsChan)
	}()

	// Collect results and errors
	for reviews := range reviewsChan {
		mu.Lock()
		allReviews = append(allReviews, reviews...)
		mu.Unlock()
	}

	// Check for errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("encountered %d errors during extraction: %v", len(errors), errors)
	}

	// Generate statistics
	stats := generateStatistics(allReviews)

	return &models.ExtractionResult{
		ExtractionDate:        time.Now(),
		TotalComments:         len(allReviews),
		RepositoriesProcessed: len(re.config.Repositories),
		Reviews:               allReviews,
		Statistics:            stats,
	}, nil
}

// generateStatistics analyzes the reviews and returns aggregated statistics
func generateStatistics(reviews []models.Review) models.Statistics {
	reviewerCount := make(map[string]int)
	fileCount := make(map[string]int)

	for _, review := range reviews {
		reviewerCount[review.CommentAuthor]++
		fileCount[review.FilePath]++
	}

	// Get most active reviewers
	mostActiveReviewers := getTopN(reviewerCount, 5)

	// Get files with most comments
	filesWithMostComments := getTopN(fileCount, 5)

	return models.Statistics{
		MostActiveReviewers:   mostActiveReviewers,
		FilesWithMostComments: filesWithMostComments,
		// Note: CommonCommentTypes would require NLP analysis
		// This is a placeholder for future enhancement
		CommonCommentTypes: []string{},
	}
}

// getTopN returns the top N keys from a map based on their values
func getTopN(counts map[string]int, n int) []string {
	if n <= 0 {
		return []string{}
	}

	type kv struct {
		key   string
		value int
	}

	var ss []kv
	for k, v := range counts {
		ss = append(ss, kv{k, v})
	}

	// Sort by value in descending order
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].value > ss[j].value
	})

	// Take top N
	result := make([]string, 0, n)
	for i := 0; i < n && i < len(ss); i++ {
		result = append(result, ss[i].key)
	}

	return result
}
