package llm

import (
	"agente-bancario/mcp"
	"agente-bancario/policy"
	"fmt"
	"strings"
)

func PlanToolCall(message string, user policy.AuthenticatedUser, availableTools []mcp.ToolDefinition) (mcp.ToolCall, error) {
	normalized := strings.ToLower(message)

	if strings.Contains(normalized, "joao") || strings.Contains(normalized, "joão") {
		return mcp.ToolCall{
			Name: "get_customer_profile",
			Arguments: map[string]string{
				"customer_id": "cust-456",
			},
		}, nil
	}

	if strings.Contains(normalized, "perfil") || strings.Contains(normalized, "meu") {
		return mcp.ToolCall{
			Name: "get_customer_profile",
			Arguments: map[string]string{
				"customer_id": user.CustomerID,
			},
		}, nil
	}

	return mcp.ToolCall{}, fmt.Errorf("nao sei qual tool chamar")
}
