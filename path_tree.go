package main

import (
	"slices"
	"strings"
)

// pathTree holds values in a tree-like hierarchy
// defined by /-separate paths (e.g. import paths).
//
// It's expensive to add items to the tree,
// but lookups are fast.
//
// Its zero value is a valid empty tree.
type pathTree[T any] struct {
	// We track two representations of the tree:
	//
	//  1. A list of key-value pairs, sorted by key.
	//  2. A map of keys to values.
	//
	// (1) is used for fast scanning of all descendants of a path.
	//
	// (2) is used for fast lookups of specific paths.

	// paths is a list of all paths in the tree
	// that have a value set for them.
	paths []string // sorted

	// values[i] is the value for paths[i]
	values []T

	// byPath is a map of all paths in the tree
	// to their corresponding values.
	byPath map[string]T // path => node
}

// Set sets the value for path to value.
// If path already has a value, it is overwritten.
func (t *pathTree[T]) Set(path string, value T) {
	if t.byPath == nil {
		t.byPath = make(map[string]T)
	}

	t.byPath[path] = value
	idx, ok := slices.BinarySearch(t.paths, path)
	if ok {
		// t.paths[idx] already contains path.
		t.values[idx] = value
	} else {
		t.paths = slices.Insert(t.paths, idx, path)
		t.values = slices.Insert(t.values, idx, value)
	}
}

// Lookup retrieves the value for the given path.
//
// If the path doesn't have an explicit value set,
// the value for the closest ancestor with a value is returned.
//
// Suffix is the remaining, unmatched part of path.
// It has a leading '/' if the path wasn't an exact match.
//
// If no value is set for the path or its ancestors,
// ok is false.
func (t *pathTree[T]) Lookup(path string) (value T, suffix string, ok bool) {
	idx := len(path)
	for idx > 0 {
		if value, ok := t.byPath[path[:idx]]; ok {
			return value, path[idx:], true
		}

		// No match. Trim the last path component.
		//	"foo/bar" => "foo"
		//	"foo" => ""
		idx = strings.LastIndexByte(path[:idx], '/')
	}

	return value, "", false
}

// ListByPath returns all descendants of path in the tree,
// including the path itself if it exists.
//
// The values are returned in an unspecified order.
func (t *pathTree[T]) ListByPath(path string) []T {
	start, end := t.rangeOf(path)
	if start == end {
		return nil
	}

	descendants := make([]T, end-start)
	for i := start; i < end; i++ {
		descendants[i-start] = t.values[i]
	}
	return descendants
}

func (t *pathTree[T]) rangeOf(path string) (start, end int) {
	if len(path) == 0 {
		return 0, len(t.paths)
	}

	start, _ = slices.BinarySearch(t.paths, path)
	for idx := start; idx < len(t.paths); idx++ {
		if descends(path, t.paths[idx]) {
			continue // path is an ancestor of p
		}

		// End of matching sequences.
		// The next path is not a descendant of path.
		return start, idx
	}

	// All paths following start are descendants of path.
	return start, len(t.paths)
}

func descends(from, to string) bool {
	return to == from || (strings.HasPrefix(to, from) && to[len(from)] == '/')
}
