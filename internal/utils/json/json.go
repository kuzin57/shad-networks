package jsonutils

import "encoding/json"

func Serialize[T any](v T) map[string]any {
	var (
		b, _ = json.Marshal(&v)
		m    = make(map[string]any)
	)

	_ = json.Unmarshal(b, &m)

	return m
}
