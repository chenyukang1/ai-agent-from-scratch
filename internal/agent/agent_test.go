package agent_test

import (
	"context"
	"strings"
	"testing"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/agent"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/chat"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/tool"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/tool/builtin"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
)

type fakeChatModel struct {
	responses []*schema.Message
	calls     int
}

func (f *fakeChatModel) Generate(_ context.Context, _ []*schema.Message, _ ...chat.Option) (*schema.Message, error) {
	if f.calls >= len(f.responses) {
		return &schema.Message{Role: schema.Assistant, Content: "done"}, nil
	}
	msg := f.responses[f.calls]
	f.calls++
	return msg, nil
}

func TestAgentRunWithToolLoop(t *testing.T) {
	model := &fakeChatModel{
		responses: []*schema.Message{
			{
				Role: schema.Assistant,
				ToolCalls: []schema.ToolCall{
					{ID: "call_1", Name: builtin.CalculatorToolName, Arguments: `{"a":2,"b":3,"op":"mul"}`},
				},
			},
			{Role: schema.Assistant, Content: "The answer is 6."},
		},
	}

	var tools []tool.BaseTool
	tools = append(tools, builtin.Calculator{})

	ag, err := agent.New(agent.Config{
		Model: model,
		Tools: tools,
	})
	if err != nil {
		t.Fatalf("new agent: %v", err)
	}

	result, err := ag.Run(context.Background(), "What is 2 times 3?")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Message.Content != "The answer is 6." {
		t.Fatalf("final content: got %q", result.Message.Content)
	}
	if model.calls != 2 {
		t.Fatalf("expected 2 model calls, got %d", model.calls)
	}
	if len(result.Messages) != 4 {
		t.Fatalf("expected 4 messages (user, assistant+toolcall, tool, assistant), got %d", len(result.Messages))
	}
	if result.Messages[2].Role != schema.Tool || result.Messages[2].Content != "6" {
		t.Fatalf("tool message: %+v", result.Messages[2])
	}
}

func TestAgentUnknownTool(t *testing.T) {
	model := &fakeChatModel{
		responses: []*schema.Message{
			{
				Role: schema.Assistant,
				ToolCalls: []schema.ToolCall{
					{ID: "call_x", Name: "missing", Arguments: `{}`},
				},
			},
			{Role: schema.Assistant, Content: "ok"},
		},
	}
	ag, err := agent.New(agent.Config{Model: model, Tools: nil})
	if err != nil {
		t.Fatalf("new agent: %v", err)
	}
	result, err := ag.Run(context.Background(), "hi")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if result.Message.Content != "ok" {
		t.Fatalf("unexpected final: %+v", result.Message)
	}
	if result.Messages[2].Role != schema.Tool || !strings.Contains(result.Messages[2].Content, "unknown tool") {
		t.Fatalf("expected unknown tool error, got %+v", result.Messages[2])
	}
}
