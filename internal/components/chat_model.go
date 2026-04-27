package components

import "context"

type BaseChatModel interface {
	Generate(ctx context.Context, input []string, opts ...Option) (string, error)
}
