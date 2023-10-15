package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestPathTree_empty(t *testing.T) {
	var tree pathTree[int]

	_, _, ok := tree.Lookup("")
	assert.False(t, ok)

	_, _, ok = tree.Lookup("foo")
	assert.False(t, ok)

	assert.Empty(t, tree.ListByPath(""))
}

func TestPathTree(t *testing.T) {
	var tree pathTree[int]
	mustHave := func(path string, want int, wantSuffix string) {
		t.Helper()

		v, suffix, ok := tree.Lookup(path)
		require.True(t, ok, "path %q", path)
		assert.Equal(t, v, want, "path %q", path)
		assert.Equal(t, wantSuffix, suffix, "path %q", path)
	}

	mustNotHave := func(path string) {
		t.Helper()

		_, _, ok := tree.Lookup(path)
		require.False(t, ok, "path %q", path)
	}

	mustList := func(path string, want ...int) {
		t.Helper()
		slices.Sort(want)

		got := tree.ListByPath(path)
		slices.Sort(got)

		assert.Equal(t, want, got, "path %q", path)
	}

	tree.Set("foo", 10)
	t.Run("single", func(t *testing.T) {
		mustHave("foo", 10, "")
		mustHave("foo/bar", 10, "/bar")
		mustHave("foo/bar/baz", 10, "/bar/baz")
		mustNotHave("")
		mustNotHave("bar")
		mustNotHave("bar/baz")

		t.Run("list", func(t *testing.T) {
			mustList("", 10)
			mustList("foo", 10)
			mustList("foo/bar")
		})
	})

	// Override a descendant value.
	t.Run("descendant", func(t *testing.T) {
		tree.Set("foo/bar", 20)
		mustHave("foo", 10, "")
		mustHave("foo/bar", 20, "")
		mustHave("foo/bar/baz", 20, "/baz")

		t.Run("list", func(t *testing.T) {
			mustList("", 10, 20)
			mustList("foo", 10, 20)
			mustList("foo/bar", 20)
			mustList("foo/bar/baz")
		})
	})

	// Add a sibling.
	t.Run("sibling", func(t *testing.T) {
		tree.Set("bar", 30)
		mustHave("bar", 30, "")
		mustHave("bar/baz", 30, "/baz")

		t.Run("list", func(t *testing.T) {
			mustList("", 10, 20, 30)
			mustList("foo", 10, 20)
			mustList("bar", 30)
			mustList("bar/baz")
		})
	})

	// Replace an existing value.
	t.Run("replace", func(t *testing.T) {
		tree.Set("bar", 40)
		mustHave("bar", 40, "")
		mustHave("bar/baz", 40, "/baz")

		t.Run("list", func(t *testing.T) {
			mustList("", 10, 20, 40)
			mustList("foo", 10, 20)
			mustList("bar", 40)
			mustList("bar/baz")
		})
	})
}

func TestPathTreeRapid(t *testing.T) {
	pathGen := rapid.StringMatching(`[a-z]+(/[a-z]+)*`)
	valueGen := rapid.Int()

	rapid.Check(t, func(t *rapid.T) {
		var tree pathTree[int]

		// Exact lookup table.
		exact := make(map[string]int)
		var paths []string // known paths

		knownPathGen := rapid.Deferred(func() *rapid.Generator[string] {
			return rapid.SampledFrom(paths)
		})

		drawKnownPath := func(t *rapid.T) string {
			if len(paths) == 0 {
				t.Skip()
			}
			return knownPathGen.Draw(t, "knownPath")
		}

		t.Repeat(map[string]func(*rapid.T){
			"Set": func(t *rapid.T) {
				path := pathGen.Draw(t, "path")
				if _, ok := exact[path]; ok {
					// Already set.
					// Overwrite will handle this.
					t.Skip()
				}

				value := valueGen.Draw(t, "value")
				tree.Set(path, value)
				exact[path] = value
				paths = append(paths, path)
			},
			"Overwrite": func(t *rapid.T) {
				path := drawKnownPath(t)
				value := valueGen.Draw(t, "value")

				tree.Set(path, value)
				exact[path] = value
				// paths already contains path.
			},
			"ExactLookup": func(t *rapid.T) {
				if len(paths) == 0 {
					t.Skip()
				}

				path := drawKnownPath(t)
				want := exact[path]

				got, suffix, ok := tree.Lookup(path)
				assert.True(t, ok, "path %q", path)
				assert.Equal(t, want, got, "path %q", path)
				assert.Empty(t, suffix, "path %q", path)
			},
			"DescendantLookup": func(t *rapid.T) {
				parentPath := drawKnownPath(t)
				want := exact[parentPath]

				var childPath string
				for {
					childPath = "/" + pathGen.Draw(t, "childPath")
					if _, ok := exact[parentPath+childPath]; !ok {
						break // found a unique child path
					}
				}

				path := parentPath + childPath
				got, suffix, ok := tree.Lookup(path)
				assert.True(t, ok, "path %q", path)
				assert.Equal(t, want, got, "path %q", path)
				assert.Equal(t, childPath, suffix, "path %q", path)
			},
			"ListAll": func(t *rapid.T) {
				var want []int
				for _, v := range exact {
					want = append(want, v)
				}
				slices.Sort(want)

				got := tree.ListByPath("")
				slices.Sort(got)

				assert.Equal(t, want, got)
			},
			"ListSubset": func(t *rapid.T) {
				path := drawKnownPath(t)

				var want []int
				for p, v := range exact {
					if descends(path, p) {
						want = append(want, v)
					}
				}
				slices.Sort(want)

				got := tree.ListByPath(path)
				slices.Sort(got)

				if len(want) == 0 {
					// Guard against nil != empty.
					assert.Empty(t, got, "path %q", path)
				} else {
					assert.Equal(t, want, got, "path %q", path)
				}
			},
		})
	})
}

func BenchmarkPathTreeDeep(b *testing.B) {
	depths := []int{10, 100}
	widths := []int{10, 100}
	for _, depth := range depths {
		b.Run(fmt.Sprintf("depth=%d", depth), func(b *testing.B) {
			for _, width := range widths {
				b.Run(fmt.Sprintf("width=%d", width), func(b *testing.B) {
					benchmarkPathTree(b, depth, width)
				})
			}
		})
	}
}

func benchmarkPathTree(b *testing.B, Depth, Width int) {
	var (
		tree    pathTree[int]
		depthpb strings.Builder
	)
	paths := make([]string, 0, Depth*Width)
	for i := 0; i < Depth; i++ {
		if depthpb.Len() > 0 {
			depthpb.WriteByte('/')
		}
		depthpb.WriteString("a")

		depthPath := depthpb.String()
		for j := 0; j < Width; j++ {
			path := depthPath + "/" + strconv.Itoa(j)
			paths = append(paths, path)
			tree.Set(path, i+i)
		}
	}

	b.ResetTimer()

	b.Run("LookupExact", func(b *testing.B) {
		path := paths[rand.Intn(len(paths))]
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _, ok := tree.Lookup(path)
				require.True(b, ok)
			}
		})
	})

	b.Run("LookupDescendant", func(b *testing.B) {
		path := paths[rand.Intn(len(paths))] + strings.Repeat("/xyz", 10)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _, ok := tree.Lookup(path)
				require.True(b, ok)
			}
		})
	})

	b.Run("ListSubtree", func(b *testing.B) {
		path := paths[rand.Intn(len(paths))]
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if len(tree.ListByPath(path)) == 0 {
					b.Fatal("unexpected empty list")
				}
			}
		})
	})

	b.Run("ListAll", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if len(tree.ListByPath("")) == 0 {
					b.Fatal("unexpected empty list")
				}
			}
		})
	})
}
