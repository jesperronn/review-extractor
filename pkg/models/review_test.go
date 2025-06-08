package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReviewCreation(t *testing.T) {
	now := time.Now()
	review := Review{
		PRID:           123,
		PRTitle:        "Test PR",
		PRAuthor:       "testuser",
		Repository:     "testrepo",
		Provider:       ProviderGitHub,
		CommentID:      "456",
		CommentAuthor:  "reviewer",
		CommentText:    "Test comment",
		CommentCreated: now,
		FilePath:       "test.go",
		LineNumber:     42,
		DiffContext:    "- old line\n+ new line",
	}

	assert.Equal(t, 123, review.PRID)
	assert.Equal(t, "Test PR", review.PRTitle)
	assert.Equal(t, "testuser", review.PRAuthor)
	assert.Equal(t, "testrepo", review.Repository)
	assert.Equal(t, ProviderGitHub, review.Provider)
	assert.Equal(t, "456", review.CommentID)
	assert.Equal(t, "reviewer", review.CommentAuthor)
	assert.Equal(t, "Test comment", review.CommentText)
	assert.Equal(t, now, review.CommentCreated)
	assert.Equal(t, "test.go", review.FilePath)
	assert.Equal(t, 42, review.LineNumber)
	assert.Equal(t, "- old line\n+ new line", review.DiffContext)
}

func TestProviderValidation(t *testing.T) {
	tests := []struct {
		name     string
		provider Provider
		valid    bool
	}{
		{"GitHub", ProviderGitHub, true},
		{"Bitbucket", ProviderBitbucket, true},
		{"GitLab", ProviderGitLab, true},
		{"Invalid", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.provider {
			case ProviderGitHub, ProviderBitbucket, ProviderGitLab:
				assert.True(t, tt.valid)
			default:
				assert.False(t, tt.valid)
			}
		})
	}
}
