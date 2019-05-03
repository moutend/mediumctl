package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArticle_Markdown(t *testing.T) {
	a, n, err := parseArticle(filepath.Join("testdata", "article.md"))

	assert.Nil(t, err)
	assert.Equal(t, 0, n, "publication number should be 0")
	assert.Equal(t, "markdown", a.ContentFormat, "ContentFormat should be \"markdown\"")
	assert.Equal(t, "Test article Markdown", a.Title)
}

func TestParseArticle_HTML(t *testing.T) {
	a, n, err := parseArticle(filepath.Join("testdata", "article.html"))

	assert.Nil(t, err)
	assert.Equal(t, 0, n, "publication number should be 0")
	assert.Equal(t, "html", a.ContentFormat, "ContentFormat should be \"html\"")
	assert.Equal(t, "Test article HTML", a.Title)
}
