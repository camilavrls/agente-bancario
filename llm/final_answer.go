package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"agente-bancario/banking"
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
		context.SafeMessage = "Para continuar com o PIX, peca ao usuario para digitar exatamente: confirmo"
		return context
	}

	context.SafeMessage = buildSuccessSafeMessage(response)
	return context
}

func GenerateFinalAnswerLocal(context finalAnswerContext) string {
	if context.SafeMessage != "" {
		return context.SafeMessage
	}

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
	case banking.PixResult:
		return fmt.Sprintf(
			"Resposta: PIX %s enviado para %s no valor de %s. Saldo restante: %s.",
			value.Transaction.ID,
			value.Transaction.ToPixKey,
			formatCents(value.Transaction.AmountCents),
			formatCents(value.RemainingBalance.BalanceCents),
		)
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

	prompt := fmt.Sprintf(`Mensagem segura obrigatoria produzida pelo backend:
%s

Contexto tecnico apenas para referencia:
%s`, context.SafeMessage, payload)

	return callGeminiWithSystemInstruction(finalAnswerSystemInstruction(), prompt, "text/plain")
}

func finalAnswerSystemInstruction() string {
	return `Voce e um assistente bancario.
Reescreva a mensagem segura obrigatoria com linguagem clara, objetiva e educada.
Preserve exatamente o sentido, status, valores, destino, saldo, confirmacao e fontes da mensagem segura.
Nao altere decisoes de autorizacao, negacao, erro, sucesso ou pendencia.
Nao transforme pendencia em falha.
Nao transforme sucesso em erro.
Nao invente links, canais de atendimento, politicas, prazos, nomes de produtos ou fontes.
Nao exponha detalhes internos como nomes de structs, JSON, ToolCall ou codigos tecnicos.
Responda em no maximo 3 frases.`
}

func buildSuccessSafeMessage(response any) string {
	switch value := response.(type) {
	case banking.PixResult:
		return fmt.Sprintf(
			"PIX executado com sucesso. Transacao %s enviada para %s no valor de %s. Saldo restante: %s.",
			value.Transaction.ID,
			value.Transaction.ToPixKey,
			formatCents(value.Transaction.AmountCents),
			formatCents(value.RemainingBalance.BalanceCents),
		)
	case []knowledge.KnowledgeResult:
		return formatKnowledgeAnswer(value)
	default:
		return fmt.Sprintf("Operacao executada com sucesso. Resultado: %+v", response)
	}
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
	case "novo limite excede o limite maximo permitido":
		return "Nao foi possivel aumentar o limite porque o valor solicitado excede o limite maximo permitido para este cliente."
	case "novo limite nao pode ser menor que o valor ja utilizado":
		return "Nao foi possivel alterar o limite porque o novo valor e menor que o valor ja utilizado no cartao."
	case "novo limite deve ser maior que zero":
		return "Nao foi possivel alterar o limite porque o novo valor precisa ser maior que zero."
	case "saldo insuficiente":
		return "Nao foi possivel concluir o PIX porque o saldo disponivel e insuficiente."
	case "valor do PIX deve ser maior que zero":
		return "Nao foi possivel concluir o PIX porque o valor precisa ser maior que zero."
	case "chave PIX de destino obrigatoria":
		return "Nao foi possivel concluir o PIX porque a chave de destino nao foi informada."
	default:
		if strings.HasPrefix(reason, "missing required argument") || strings.HasPrefix(reason, "invalid ") {
			return "Nao consegui executar a operacao porque faltam informacoes validas."
		}
		if strings.HasPrefix(reason, "limite de cartão não encontrado") {
			return "Nao encontrei informacoes de limite para este cliente."
		}
		if strings.HasPrefix(reason, "cliente não encontrado") {
			return "Nao encontrei o cliente informado."
		}
		if strings.HasPrefix(reason, "saldo não encontrado") {
			return "Nao encontrei informacoes de saldo para este cliente."
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
