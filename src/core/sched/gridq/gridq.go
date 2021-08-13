package gridq

// ============================================================================

var grids = make(map[int32]*queue_t)

// ============================================================================

type queue_t []func()

// ============================================================================

func Push(id int32, f func(done func())) {
	// get q
	q := grids[id]
	if q == nil {
		q = new(queue_t)
		grids[id] = q
	}

	// append
	*q = append(*q, func() {
		done := func() {
			// remove from q
			(*q)[0] = nil
			*q = (*q)[1:]

			// check q
			if len(*q) == 0 {
				// remove the entire q
				delete(grids, id)
			} else {
				// continue executing
				(*q)[0]()
			}
		}

		f(done)
	})

	// execute if we're the only one
	if len(*q) == 1 {
		(*q)[0]()
	}
}
