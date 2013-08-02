package web

import (
	"testing"
	//"fmt"
)

func TestParsePaths(t *testing.T) {
	paths := []string{"/foo/bar", "/foo/*", "/foo/**", "/foo/bar/baz", "/foo/bar/baz/qux", "/"}
	pathMap := ParsePaths(paths)

	if len(pathMap) != 2 {
		t.Error("! Expected 2 items, got ", len(pathMap))
	}

}
