// Resolver for transforming a URI path into a route.
//
// This is a more sophisticated path resolver, aware of
// heirarchyand wildcards.
//
// Examples:
// - URI path `/foo` matches the entry `/foo`
// - URI path `/foo/bar` could match entries like `/foo/*`, `/foo/**`, and `/foo/bar`
// - URI path `/foo/bar/baz` could match `/foo/*/baz` and `/foo/**`
//
// The most exact match "wins". E.g. for registry items `/foo/bar` and `/foo/**`, if the 
// URI path is `/foo/bar`, the `/foo/bar` entry will match first.

package web

import (
	"strings"
	"fmt"
)

func ParsePaths(paths []string) map[string]*pathEntry {
	tree := make(map[string]*pathEntry)
	for _, item := range paths {
		pieces := strings.Split(item, "/")
		fmt.Printf("Pieces: %v\n", pieces)
		entry  := parsePaths(pieces, tree)
		fmt.Printf("Adding %s to map.", entry.path)
		tree[entry.path] = entry
	}
	return tree;
	//root := newPathEntry("")
	//root.SubPaths = tree
	//fmt.Printf("Tree: %+v\n", root.SubPaths)
	//return root
}

func parsePaths(paths []string, tree map[string]*pathEntry) *pathEntry {
	head := paths[0]

	entry := newPathEntry(head)

	if len(paths) > 1 {
		tail := paths[1:]
		subentry := parsePaths(tail, tree)
		entry.SubPaths[subentry.path] = subentry
	}

	return entry
}

type pathEntry struct {
	path string
	SubPaths map[string]*pathEntry
}

func newPathEntry(path string) *pathEntry {
	entry := new(pathEntry)
	entry.path = path
	entry.SubPaths = make(map[string]*pathEntry)
	return entry
}

func (p *pathEntry) addEntry(entry *pathEntry) {
	p.SubPaths[entry.path] = entry
}

