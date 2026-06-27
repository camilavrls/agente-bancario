package llm

import (
	"agente-bancario/mcp"
	"agente-bancario/policy"
)

type geminiInteractionRequest struct {
	Model             string `json:"model"`
	SystemInstruction string `json:"system_instruction"`
	Input             string `json:"input"`
}

type geminiInteractionResponse struct {
	OutputText string `json:"output_text"`
}

func toolPlanningSystemInstruction() string {
	return `Voce e um planejador de ferramentas para um agente bancario.
     Sua tarefa e escolher uma das tools disponiveis.
     Nao execute autorizacao.
     Nao invente tools.
     Responda somente JSON validos.
     O JSON deve ter:
       Name
       Arguments `
}

func buildToolPlanningPrompt(mensagem string, user policy.AuthenticatedUser, availableTools []mcp.ToolDefinition) (string, error) {

	toolsJSON, err := 

}
