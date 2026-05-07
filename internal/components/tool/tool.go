package tool

import (
	"context"
	"fmt"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
)

// BaseTool exposes metadata and execution for agent tool calling.
type BaseTool interface {
	Info(ctx context.Context) (*schema.ToolInfo, error)
	Invoke(ctx context.Context, argumentsInJSON string) (string, error)
}

// CollectInfos returns tool metadata for all tools in order.
func CollectInfos(ctx context.Context, tools []BaseTool) ([]*schema.ToolInfo, error) {
	infos := make([]*schema.ToolInfo, 0, len(tools))
	for _, t := range tools {
		if t == nil {
			continue
		}
		info, err := t.Info(ctx)
		if err != nil {
			return nil, fmt.Errorf("tool info: %w", err)
		}
		if info == nil {
			return nil, fmt.Errorf("tool %T returned nil info", t)
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// ByName indexes tools by ToolInfo.Name.
func ByName(tools []BaseTool) (map[string]BaseTool, error) {
	out := make(map[string]BaseTool, len(tools))
	for _, t := range tools {
		if t == nil {
			continue
		}
		info, err := t.Info(context.Background())
		if err != nil {
			return nil, err
		}
		if info.Name == "" {
			return nil, fmt.Errorf("tool %T has empty name", t)
		}
		if _, exists := out[info.Name]; exists {
			return nil, fmt.Errorf("duplicate tool name %q", info.Name)
		}
		out[info.Name] = t
	}
	return out, nil
}
