package agent

import (
	"agente-bancario/mcp"
	"agente-bancario/policy"
	"agente-bancario/tools"
)

func HandleToolCall(user policy.AuthenticatedUser, toolcall mcp.ToolCall) (any, error) {

	responseTool, err := tools.ExecuteTool(user, toolcall)
	if err != nil {
		return nil, err
	}

	return responseTool, nil
}
