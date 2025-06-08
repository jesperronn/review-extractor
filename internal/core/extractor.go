package core

import (
	"context"
	"fmt"
	"sort"
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

// ExtractReviews extracts reviews from all configured repositories
func (e *ReviewExtractor) ExtractReviews(ctx context.Context) (*models.ExtractionResult, error) {
	var allReviews []models.Review

	// Process each repository
	for _, repo := range e.config.Repositories {
		extractor, ok := e.extractors[repo.Provider]
		if !ok {
			return nil, fmt.Errorf("no extractor available for provider: %s", repo.Provider)
		}

		reviews, err := extractor.ExtractReviews(ctx, repo.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to extract reviews from %s: %w", repo.URL, err)
		}

		allReviews = append(allReviews, reviews...)
	}

	// Generate statistics
	stats := generateStatistics(allReviews)

	// Create result
	result := &models.ExtractionResult{
                Reviews:               allReviews,
                Statistics:            stats,
                ExtractedAt:           time.Now(),
                TotalComments:         len(allReviews),
                RepositoriesProcessed: len(e.config.Repositories),
	}

	return result, nil
}

// generateStatistics analyzes the reviews and returns aggregated statistics
func generateStatistics(reviews []models.Review) models.Statistics {
	reviewerCounts := make(map[string]int)
	repoCounts := make(map[string]int)
	prCounts := make(map[int]bool)
	prSizes := make(map[int]int)

	for _, review := range reviews {
		reviewerCounts[review.CommentAuthor]++
		repoCounts[review.Repository]++
		prCounts[review.PRID] = true
		prSizes[review.PRID]++
	}

	// Calculate average PR size
	var totalPRSize int
	for _, size := range prSizes {
		totalPRSize += size
	}
	averagePRSize := 0.0
	if len(prCounts) > 0 {
		averagePRSize = float64(totalPRSize) / float64(len(prCounts))
	}

	// Calculate review frequency (reviews per PR)
	reviewFrequency := 0.0
	if len(prCounts) > 0 {
		reviewFrequency = float64(len(reviews)) / float64(len(prCounts))
	}

	return models.Statistics{
		TotalReviews:    len(reviews),
		TotalPRs:        len(prCounts),
		TopReviewers:    getTopN(reviewerCounts, 5),
		TopRepositories: getTopN(repoCounts, 5),
		AveragePRSize:   averagePRSize,
		ReviewFrequency: reviewFrequency,
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
