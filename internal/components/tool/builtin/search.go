package builtin

import (
	"context"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
)

type SearchTool struct{}

func (s SearchTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "search",
		Desc: "Search the web for information.",
		Parameters: schema.ObjectParams(map[string]any{
			"query": map[string]any{"type": "string", "description": "The search query."},
		}, "query"),
	}, nil
}

func (s SearchTool) Invoke(ctx context.Context, argumentsInJSON string) (string, error) {
	// For demonstration purposes, we will just return a static response.
	// In a real implementation, you would parse the arguments and perform a web search.
	return "Search results for query: " + argumentsInJSON, nil
}
