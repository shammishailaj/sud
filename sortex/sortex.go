package sortex

import (
	"sort"
)

type sString struct {
	array []string
	less  func(i, j int) bool
}

func (a sString) Len() int           { return len(a.array) }
func (a sString) Swap(i, j int)      { a.array[i], a.array[j] = a.array[j], a.array[i] }
func (a sString) Less(i, j int) bool { return a.less(i, j) }
func SortStrings(arr []string, less func(i, j int) bool) {
	sort.Sort(sString{array: arr, less: less})
}
