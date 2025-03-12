package slices

func Map[V, U any](in []V, f func(int, V) U) []U {
	out := make([]U, 0, len(in))

	for i, el := range in {
		out = append(out, f(i, el))
	}

	return out
}
