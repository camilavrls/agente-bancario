package tools

import (
	"strconv"

	"agente-bancario/mcp"
	"agente-bancario/policy"
	"fmt"
)

func ListTools() []mcp.ToolDefinition {
	return []mcp.ToolDefinition{
		GetCustomerProfileToolDefinition(),
		GetCardLimitToolDefinition(),
		UpdateCardLimitToolDefinition(),
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

	case "get_card_limit":
		customerID, ok := toolcall.Arguments["customer_id"]
		if !ok {
			return nil, fmt.Errorf("missing required argument: customer_id")
		}

		return GetCardLimitTool(user, customerID)

	case "update_card_limit":
		customerID, ok := toolcall.Arguments["customer_id"]
		if !ok {
			return nil, fmt.Errorf("missing required argument: customer_id")
		}

		newLimitRaw, ok := toolcall.Arguments["new_limit"]
		if !ok {
			return nil, fmt.Errorf("missing required argument: new_limit")
		}

		newLimit, err := strconv.Atoi(newLimitRaw)
		if err != nil {
			return nil, fmt.Errorf("invalid new_limit: %w", err)
		}

		return UpdateCardLimitTool(user, customerID, newLimit)

	default:
		return nil, fmt.Errorf("unknown tool")
	}
}
