// Author: Yoshiyuki Koyanagi <moutend@gmail.com>
// License: mIT

package main

import (
	"path/filepath"
	"testing"
)

func TestParseArticle(t *testing.T) {
	filename := filepath.Join("testdata", "a.md")
	a, n, err := parseArticle(filename)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("publication number should be 0 but %d", n)
	}
	expectedTitle := "Test article A"
	actualTitle := a.Title
	if expectedTitle != actualTitle {
		t.Errorf("\nexpected: %s\nactual: %s\n", expectedTitle, actualTitle)
	}
	return
}
