package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"agente-bancario/knowledge"
	"agente-bancario/mcp"
)

type finalAnswerContext struct {
	UserMessage string       `json:"user_message"`
	ToolCall    mcp.ToolCall `json:"tool_call"`
	Status      string       `json:"status"`
	Data        any          `json:"data,omitempty"`
	SafeMessage string       `json:"safe_message,omitempty"`
}

func GenerateFinalAnswer(userMessage string, toolCall mcp.ToolCall, response any, executionErr error) (string, error) {
	provider := strings.ToLower(os.Getenv("LLM_PROVIDER"))
	if provider == "" {
		provider = llmProviderLocal
	}

	context := buildFinalAnswerContext(userMessage, toolCall, response, executionErr)

	switch provider {
	case llmProviderGemini:
		return GenerateFinalAnswerGemini(context)
	case llmProviderLocal, "heuristic":
		return GenerateFinalAnswerLocal(context), nil
	default:
		return "", fmt.Errorf("LLM_PROVIDER invalido: %s", provider)
	}
}

func buildFinalAnswerContext(userMessage string, toolCall mcp.ToolCall, response any, executionErr error) finalAnswerContext {
	context := finalAnswerContext{
		UserMessage: userMessage,
		ToolCall:    toolCall,
		Status:      "success",
		Data:        response,
	}

	if executionErr != nil {
		context.Status = "error"
		context.Data = nil
		context.SafeMessage = safeMessageForError(executionErr.Error())
		return context
	}

	if response == "pix_requires_confirmation" {
		context.Status = "pending"
		context.SafeMessage = "Para continuar com o PIX, solicite confirmacao explicita do usuario."
	}

	return context
}

func GenerateFinalAnswerLocal(context finalAnswerContext) string {
	switch context.Status {
	case "pending":
		return "PIX pendente de confirmacao. Digite \"confirmo\" para executar."
	case "error":
		if context.SafeMessage != "" {
			return fmt.Sprintf("Execucao negada ou falhou: %s", context.SafeMessage)
		}
		return "Execucao negada ou falhou."
	}

	switch value := context.Data.(type) {
	case []knowledge.KnowledgeResult:
		return formatKnowledgeAnswer(value)
	default:
		return fmt.Sprintf("Resposta: %+v", context.Data)
	}
}

func GenerateFinalAnswerGemini(context finalAnswerContext) (string, error) {
	payload, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`Mensagem original e resultado seguro produzido pelo backend:
%s`, payload)

	return callGeminiWithSystemInstruction(finalAnswerSystemInstruction(), prompt, "text/plain")
}

func finalAnswerSystemInstruction() string {
	return `Voce e um assistente bancario.
Redija a resposta final para o cliente com linguagem clara, objetiva e educada.
Use somente o resultado seguro produzido pelo backend.
Nao altere decisoes de autorizacao, negacao, erro ou pendencia.
Se o status for denied ou error, explique a negativa sem sugerir formas de contornar autorizacao.
Se houver fonte da base de conhecimento, cite a fonte.
Nao exponha detalhes internos como nomes de structs, JSON, ToolCall ou codigos tecnicos.`
}

func safeMessageForError(reason string) string {
	switch reason {
	case "denied_customer_other_resource", "denied_customer_update_other_card_limit", "denied_customer_create_other_pix":
		return "Nao posso acessar ou alterar dados de outro cliente com o seu perfil de acesso."
	case "denied_manager_create_pix", "denied_admin_create_pix":
		return "Nao posso executar PIX usando um perfil que nao representa o proprio cliente da conta."
	case "denied_unknown_role":
		return "Nao foi possivel validar sua permissao para esta operacao."
	case "no_pending_action":
		return "Nao encontrei nenhuma operacao pendente para confirmar."
	default:
		if strings.HasPrefix(reason, "missing required argument") || strings.HasPrefix(reason, "invalid ") {
			return "Nao consegui executar a operacao porque faltam informacoes validas."
		}
		return "Nao foi possivel concluir a operacao solicitada."
	}
}

func formatKnowledgeAnswer(results []knowledge.KnowledgeResult) string {
	if len(results) == 0 {
		return "Nao encontrei informacoes suficientes na base de conhecimento."
	}

	result := results[0]
	answer := strings.TrimSpace(result.Answer)
	if result.Source == "" {
		return answer
	}

	return fmt.Sprintf("%s\n\nFonte: %s", answer, result.Source)
}

func formatCents(amountCents int) string {
	return fmt.Sprintf("R$ %.2f", float64(amountCents)/100)
}
