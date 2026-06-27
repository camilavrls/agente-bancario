package tools

import (
	"agente-bancario/mcp"
	"agente-bancario/policy"
	"fmt"
)

func ListTools() []mcp.ToolDefinition {
	return []mcp.ToolDefinition{
		GetCustomerProfileToolDefinition(),
	}
}

func ExecuteTool(user policy.AuthenticatedUser, toolcall mcp.ToolCall) (any, error) {
	switch toolcall.Name {
	case "get_customer_profile":
		customerID, ok := toolcall.Arguments["customer_id"]
		if !ok {
			return nil, fmt.Errorf("missing required argument: customer_id")
		}

		return GetCustomerProfileTool(user, customerID)
	default:
		return nil, fmt.Errorf("unknown tool")
	}
}
