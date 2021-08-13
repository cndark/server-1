package evtmgr

// ============================================================================

var (
	evt_map = make(map[string][]func(...interface{}))
)

// ============================================================================

// setting up events at initialization stage is preferred.
// be care of race conditions if you don't do like this.
func On(evt string, f func(...interface{})) {
	evt_map[evt] = append(evt_map[evt], f)
}

func Fire(evt string, args ...interface{}) {
	arr := evt_map[evt]
	for _, f := range arr {
		f(args...)
	}
}
