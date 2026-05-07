package schema

// ObjectParams builds a JSON Schema object for tool parameters.
func ObjectParams(properties map[string]any, required ...string) map[string]any {
	p := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		p["required"] = required
	}
	return p
}
