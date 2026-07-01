package llm

import (
	"agente-bancario/mcp"
	"agente-bancario/policy"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	llmProviderGemini    = "gemini"
	llmProviderHeuristic = "heuristic"
)

func PlanToolCall(message string, user policy.AuthenticatedUser, availableTools []mcp.ToolDefinition) (mcp.ToolCall, error) {
	provider := strings.ToLower(os.Getenv("LLM_PROVIDER"))
	if provider == "" {
		provider = llmProviderHeuristic
	}

	switch provider {
	case llmProviderGemini:
		return PlanToolCallGemini(message, user, availableTools)
	case llmProviderHeuristic:
		return PlanToolCallHeuristic(message, user, availableTools)
	default:
		return mcp.ToolCall{}, fmt.Errorf("LLM_PROVIDER invalido: %s", provider)
	}
}

func PlanToolCallHeuristic(message string, user policy.AuthenticatedUser, availableTools []mcp.ToolDefinition) (mcp.ToolCall, error) {
	normalized := strings.ToLower(message)
	customerID := user.CustomerID

	if strings.Contains(normalized, "joao") || strings.Contains(normalized, "joão") {
		customerID = "cust-456"
	}

	if isCardLimitUpdateIntent(normalized) {
		newLimit, err := extractFirstNumber(message)
		if err != nil {
			return mcp.ToolCall{}, err
		}

		return mcp.ToolCall{
			Name: "update_card_limit",
			Arguments: map[string]string{
				"customer_id": customerID,
				"new_limit":   newLimit,
			},
		}, nil
	}

	if strings.Contains(normalized, "limite") || strings.Contains(normalized, "cartao") || strings.Contains(normalized, "cartão") {
		return mcp.ToolCall{
			Name: "get_card_limit",
			Arguments: map[string]string{
				"customer_id": customerID,
			},
		}, nil
	}

	if strings.Contains(normalized, "perfil") || strings.Contains(normalized, "meu") {
		return mcp.ToolCall{
			Name: "get_customer_profile",
			Arguments: map[string]string{
				"customer_id": customerID,
			},
		}, nil
	}

	return mcp.ToolCall{}, fmt.Errorf("nao sei qual tool chamar")
}

func isCardLimitUpdateIntent(message string) bool {
	hasLimitContext := strings.Contains(message, "limite") || strings.Contains(message, "cartao") || strings.Contains(message, "cartão")
	hasUpdateVerb := strings.Contains(message, "aumentar") ||
		strings.Contains(message, "alterar") ||
		strings.Contains(message, "ajustar") ||
		strings.Contains(message, "mudar")

	return hasLimitContext && hasUpdateVerb
}

func extractFirstNumber(message string) (string, error) {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(message)
	if match == "" {
		return "", fmt.Errorf("nao encontrei o novo limite no prompt")
	}

	return match, nil
}
