package db

import (
	"sort"
)

func (d *Domain) sortPathMap() {
	var fr By
	if d.PathMapping == nil {
		return
	}
	fr = func(p1, p2 *PathMappingEntry) bool {
		return len(p1.From) > len(p2.From)
	}
	fr.Sort(*d.PathMapping)
}

// planetSorter joins a By function and a slice of Paths to be sorted.
type pathMapSorter struct {
	paths PathList
	by    By
}

// By is the type of a "less" function that defines the ordering of its arguments.
type By func(p1, p2 *PathMappingEntry) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(paths PathList) {
	ps := &pathMapSorter{
		paths: paths,
		by:    by,
	}
	sort.Sort(ps)
}

// Len is part of sort.Interface
func (s *pathMapSorter) Len() int {
	return len(s.paths)
}

// Swap is part of sort.Interface
func (s *pathMapSorter) Swap(i, j int) {
	s.paths[i], s.paths[j] = s.paths[j], s.paths[i]
}

// Less is part of sort.Interface
func (s *pathMapSorter) Less(i, j int) bool {
	return s.by(&s.paths[i], &s.paths[j])
}
