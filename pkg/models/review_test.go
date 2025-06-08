package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReviewCreation(t *testing.T) {
	now := time.Now()
	review := Review{
		PRID:           1,
		PRTitle:        "Test PR",
		PRAuthor:       "testuser",
		Repository:     "testrepo",
		Provider:       ProviderGitHub,
		CommentID:      "123",
		CommentAuthor:  "reviewer",
		CommentText:    "Test comment",
		CommentCreated: now,
		FilePath:       "test.go",
		LineNumber:     10,
		DiffContext:    "test diff",
	}

	assert.Equal(t, 1, review.PRID)
	assert.Equal(t, "Test PR", review.PRTitle)
	assert.Equal(t, "testuser", review.PRAuthor)
	assert.Equal(t, "testrepo", review.Repository)
	assert.Equal(t, ProviderGitHub, review.Provider)
	assert.Equal(t, "123", review.CommentID)
	assert.Equal(t, "reviewer", review.CommentAuthor)
	assert.Equal(t, "Test comment", review.CommentText)
	assert.Equal(t, now, review.CommentCreated)
	assert.Equal(t, "test.go", review.FilePath)
	assert.Equal(t, 10, review.LineNumber)
	assert.Equal(t, "test diff", review.DiffContext)
}

func TestProviderValidation(t *testing.T) {
	tests := []struct {
		name     string
		provider Provider
		valid    bool
	}{
		{
			name:     "GitHub provider",
			provider: ProviderGitHub,
			valid:    true,
		},
		{
			name:     "GitLab provider",
			provider: ProviderGitLab,
			valid:    true,
		},
		{
			name:     "Invalid provider",
			provider: "invalid",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.provider == ProviderGitHub || tt.provider == ProviderGitLab
			assert.Equal(t, tt.valid, valid)
		})
	}
}
