package agent

import (
	"context"
	"fmt"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/chat"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/tool"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
)

const defaultMaxSteps = 10

// Agent runs a tool-calling loop over a chat model.
type Agent struct {
	model    chat.BaseChatModel
	tools    map[string]tool.BaseTool
	toolList []*schema.ToolInfo
	maxSteps int
}

// Config configures an Agent.
type Config struct {
	Model    chat.BaseChatModel
	Tools    []tool.BaseTool
	MaxSteps int
}

// New builds an Agent from config.
func New(cfg Config) (*Agent, error) {
	if cfg.Model == nil {
		return nil, fmt.Errorf("agent: model is required")
	}
	byName, err := tool.ByName(cfg.Tools)
	if err != nil {
		return nil, err
	}
	toolInfos, err := tool.CollectInfos(context.Background(), cfg.Tools)
	if err != nil {
		return nil, err
	}
	maxSteps := cfg.MaxSteps
	if maxSteps <= 0 {
		maxSteps = defaultMaxSteps
	}
	return &Agent{
		model:    cfg.Model,
		tools:    byName,
		toolList: toolInfos,
		maxSteps: maxSteps,
	}, nil
}

// RunResult holds the final assistant message and full conversation history.
type RunResult struct {
	Message  *schema.Message
	Messages []*schema.Message
}

// Run executes the agent for a single user message.
func (a *Agent) Run(ctx context.Context, userInput string) (*RunResult, error) {
	messages := []*schema.Message{
		{Role: schema.User, Content: userInput},
	}

	var lastAssistant *schema.Message
	for step := 0; step < a.maxSteps; step++ {
		var opts []chat.Option
		if len(a.toolList) > 0 {
			opts = append(opts, chat.WithTools(a.toolList))
		}

		resp, err := a.model.Generate(ctx, messages, opts...)
		if err != nil {
			return nil, err
		}
		messages = append(messages, resp)
		lastAssistant = resp

		if len(resp.ToolCalls) == 0 {
			return &RunResult{Message: resp, Messages: messages}, nil
		}

		for _, tc := range resp.ToolCalls {
			toolMsg, err := a.invokeTool(ctx, tc)
			if err != nil {
				return nil, err
			}
			messages = append(messages, toolMsg)
		}
	}

	if lastAssistant != nil {
		return &RunResult{Message: lastAssistant, Messages: messages},
			fmt.Errorf("agent: exceeded max steps (%d)", a.maxSteps)
	}
	return nil, fmt.Errorf("agent: exceeded max steps (%d)", a.maxSteps)
}

func (a *Agent) invokeTool(ctx context.Context, tc schema.ToolCall) (*schema.Message, error) {
	t, ok := a.tools[tc.Name]
	if !ok {
		return &schema.Message{
			Role:       schema.Tool,
			ToolCallID: tc.ID,
			Content:    fmt.Sprintf("error: unknown tool %q", tc.Name),
		}, nil
	}
	result, err := t.Invoke(ctx, tc.Arguments)
	if err != nil {
		return &schema.Message{
			Role:       schema.Tool,
			ToolCallID: tc.ID,
			Content:    fmt.Sprintf("error: %v", err),
		}, nil
	}
	return &schema.Message{
		Role:       schema.Tool,
		ToolCallID: tc.ID,
		Content:    result,
	}, nil
}
