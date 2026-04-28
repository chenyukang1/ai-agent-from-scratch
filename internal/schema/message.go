package schema

// RoleType is the type of the role of a message.
type RoleType string

const (
	// Assistant is the role of an assistant, means the message is returned by ChatModel.
	Assistant RoleType = "assistant"
	// User is the role of a user, means the message is a user message.
	User RoleType = "user"
	// System is the role of a system, means the message is a system message.
	System RoleType = "system"
	// Tool is the role of a tool, means the message is a tool call output.
	Tool RoleType = "tool"
)

type Message struct {
	// Role is the role of a message.
	Role RoleType `json:"role"`
	// Content is for user text input and model text output.
	Content string `json:"content"`
}
