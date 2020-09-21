package fs

import (
	"sort"
	"strings"
)

type lessFunc func(f1, f2 File) bool

// multiSorter implements the Sort interface, sorting the changes within.
type multiSorter struct {
	files []File
	less  []lessFunc
}

// sort sorts the argument slice according to the less functions passed to orderedBy.
func (ms *multiSorter) sort(files []File) {
	ms.files = files
	sort.Sort(ms)
}

// orderedBy returns a Sorter that sorts using the less functions, in order.
// Call its sort method to sort the data.
func orderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.files)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.files[i], ms.files[j] = ms.files[j], ms.files[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := ms.files[i], ms.files[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

func name(f1, f2 File) bool {
	return strings.ToUpper(f1.Name()) < strings.ToUpper(f2.Name())
}

func dirFirst(f1, f2 File) bool {
	return f1.IsDir() && !f2.IsDir()
}
