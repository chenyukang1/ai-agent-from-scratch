package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/tool"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
)

const CalculatorToolName = "calculator"

// Calculator evaluates simple binary arithmetic (add, sub, mul, div).
type Calculator struct{}

func (Calculator) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: CalculatorToolName,
		Desc: "Evaluate binary arithmetic on two numbers. Use op one of: add, sub, mul, div.",
		Parameters: schema.ObjectParams(map[string]any{
			"a":  map[string]any{"type": "number", "description": "First operand"},
			"b":  map[string]any{"type": "number", "description": "Second operand"},
			"op": map[string]any{"type": "string", "enum": []any{"add", "sub", "mul", "div"}, "description": "Operation"},
		}, "a", "b", "op"),
	}, nil
}

type calculatorArgs struct {
	A  float64 `json:"a"`
	B  float64 `json:"b"`
	Op string  `json:"op"`
}

func (Calculator) Invoke(_ context.Context, argumentsInJSON string) (string, error) {
	var args calculatorArgs
	if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
		return "", fmt.Errorf("calculator: invalid arguments: %w", err)
	}
	var result float64
	switch args.Op {
	case "add":
		result = args.A + args.B
	case "sub":
		result = args.A - args.B
	case "mul":
		result = args.A * args.B
	case "div":
		if args.B == 0 {
			return "", fmt.Errorf("calculator: division by zero")
		}
		result = args.A / args.B
	default:
		return "", fmt.Errorf("calculator: unknown op %q", args.Op)
	}
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return "", fmt.Errorf("calculator: non-finite result")
	}
	return fmt.Sprintf("%g", result), nil
}

var _ tool.BaseTool = Calculator{}
