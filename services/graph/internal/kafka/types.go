package kafka

import (
	"encoding/json"
	"fmt"
)

type Request []byte

func ExtractMessage[T any](req Request) (*T, error) {
	var extracted T

	if err := json.Unmarshal(req, &extracted); err != nil {
		return nil, fmt.Errorf("extract message: %w", err)
	}

	return &extracted, nil
}

type PathsProcessingMessage struct {
	Source  int    `json:"source"`
	GraphID string `json:"graphID"`
}
