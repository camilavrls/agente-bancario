package mcp

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []ToolParameter
}

type ToolParameter struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

type ToolCall struct {
	Name      string
	Arguments map[string]string
}
