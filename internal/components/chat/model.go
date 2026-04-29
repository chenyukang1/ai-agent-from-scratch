package chat

import (
	"context"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
)

type BaseChatModel interface {
	Generate(ctx context.Context, input []*schema.Message, opts ...Option) (*schema.Message, error)
}
