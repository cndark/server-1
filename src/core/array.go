package core

import (
	"reflect"
	"sort"
)

// sorted insert.
// less: a[i] < e. in ascending order
func SInsert(a interface{}, e interface{}, less func(i int) bool) {
	av := reflect.ValueOf(a)
	ev := reflect.ValueOf(e)
	n := av.Elem().Len()

	// search
	i := sort.Search(n, func(i int) bool { return !less(i) })

	// a[:i], e, a[i:]
	r := reflect.Append(av.Elem(), ev)                      // insure cap
	r = reflect.AppendSlice(r.Slice(0, i+1), r.Slice(i, n)) // leave blank
	r.Index(i).Set(ev)                                      // insert into blank

	// set
	av.Elem().Set(r)
}

func ArrayFind(n int, f func(int) bool) int {
	for i := 0; i < n; i++ {
		if f(i) {
			return i
		}
	}
	return -1
}
