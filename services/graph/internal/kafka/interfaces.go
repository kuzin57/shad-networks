package kafka

import "context"

type Processor interface {
	Process(ctx context.Context, request Request) error
}

type ProcessorsHolder interface {
	GetProcessor(name string) Processor
}
