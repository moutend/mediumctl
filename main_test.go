// Author: Yoshiyuki Koyanagi <moutend@gmail.com>
// License: mIT

package main

import (
	"path/filepath"
	"testing"
)

func TestParseArticle_Markdown(t *testing.T) {
	filename := filepath.Join("testdata", "a.md")
	a, n, err := parseArticle(filename)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("publication number should be 0 but %d", n)
	}
	if a.ContentFormat != "markdown" {
		t.Fatalf("ContentFormat should be \"markdown\", but %s\n", a.ContentFormat)
	}
	expectedTitle := "Test article Markdown"
	actualTitle := a.Title
	if expectedTitle != actualTitle {
		t.Errorf("\nexpected: %s\nactual: %s\n", expectedTitle, actualTitle)
	}
	return
}

func TestParseArticle_HTML(t *testing.T) {
	filename := filepath.Join("testdata", "a.html")
	a, n, err := parseArticle(filename)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("publication number should be 0 but %d", n)
	}
	if a.ContentFormat != "html" {
		t.Fatalf("ContentFormat should be \"html\" but %s\n", a.ContentFormat)
	}
	expectedTitle := "Test article HTML"
	actualTitle := a.Title
	if expectedTitle != actualTitle {
		t.Errorf("\nexpected: %s\nactual: %s\n", expectedTitle, actualTitle)
	}
	return
}
