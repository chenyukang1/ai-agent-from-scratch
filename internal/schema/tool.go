package schema

// ToolInfo describes a tool for the model (function calling).
type ToolInfo struct {
	// Name is the unique name of the tool that clearly communicates its purpose.
	Name string
	// Desc tells the model how/when/why to use the tool.
	Desc string
	// Parameters is a JSON Schema object describing the tool arguments.
	// Example: {"type":"object","properties":{...},"required":[...]}
	// If nil, the tool accepts no parameters.
	Parameters map[string]any
	// Extra is optional provider-specific metadata.
	Extra map[string]any
}

// ToolCall is a model-requested tool invocation on an assistant message.
type ToolCall struct {
	ID        string
	Name      string
	Arguments string // JSON object string
}
