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
		CreatePixToolDefinition(),
		SearchKnowledgeBaseToolDefinition(),
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

	case "create_pix":
		customerID, ok := toolcall.Arguments["customer_id"]
		if !ok {
			return nil, fmt.Errorf("missing required argument: customer_id")
		}

		pixKey, ok := toolcall.Arguments["pix_key"]
		if !ok {
			return nil, fmt.Errorf("missing required argument: pix_key")
		}

		amountCentsRaw, ok := toolcall.Arguments["amount_cents"]
		if !ok {
			return nil, fmt.Errorf("missing required argument: amount_cents")
		}

		amountCents, err := strconv.Atoi(amountCentsRaw)
		if err != nil {
			return nil, fmt.Errorf("invalid amount_cents: %w", err)
		}

		confirmed := toolcall.Arguments["confirmed"] == "true"

		return CreatePixTool(user, customerID, pixKey, amountCents, confirmed)

	case "search_knowledge_base":
		query, ok := toolcall.Arguments["query"]
		if !ok {
			return nil, fmt.Errorf("missing required argument: query")
		}

		return SearchKnowledgeBaseTool(query)

	default:
		return nil, fmt.Errorf("unknown tool")
	}
}
