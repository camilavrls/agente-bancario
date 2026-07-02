package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"agente-bancario/agent"
	"agente-bancario/knowledge"
	"agente-bancario/llm"
	"agente-bancario/policy"
	"agente-bancario/tools"
)

func main() {
	orchestrator := agent.NewOrchestrator()

	user := policy.AuthenticatedUser{
		ID:         "123",
		Name:       "Maria",
		Role:       policy.RoleCustomer,
		CustomerID: "cust-123",
	}

	printStartup(user)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		message := strings.TrimSpace(scanner.Text())
		if message == "" {
			continue
		}

		if isExitCommand(message) {
			fmt.Println("Encerrando.")
			break
		}

		if isConfirmation(message) {
			confirmPendingAction(orchestrator, user)
			continue
		}

		runPrompt(orchestrator, user, message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler entrada:", err)
	}
}

func printStartup(user policy.AuthenticatedUser) {
	provider := os.Getenv("LLM_PROVIDER")
	if provider == "" {
		provider = "local"
	}

	fmt.Println("Agente bancario inteligente")
	fmt.Println("Usuario autenticado:", user.Name, "("+user.CustomerID+")")
	fmt.Println("LLM_PROVIDER:", provider)
	fmt.Println("Digite uma mensagem ou 'sair' para encerrar.")
	fmt.Println("Exemplos:")
	fmt.Println("- quero consultar meu perfil")
	fmt.Println("- qual e o limite do meu cartao?")
	fmt.Println("- quero aumentar meu limite para 12000")
	fmt.Println("- fazer pix de 2000 para joao")
	fmt.Println()
}

func runPrompt(orchestrator *agent.Orchestrator, user policy.AuthenticatedUser, prompt string) {
	availableTools := tools.ListTools()

	call, err := llm.PlanToolCall(prompt, user, availableTools)
	if err != nil {
		fmt.Println("Erro ao planejar tool call:", err)
		return
	}

	fmt.Printf("Tool call planejada: %+v\n", call)

	response, err := orchestrator.HandleToolCall(user, call)
	if err != nil {
		fmt.Println("Execucao negada ou falhou:", err)
		return
	}

	if response == "pix_requires_confirmation" {
		fmt.Println("PIX pendente de confirmacao. Digite \"confirmo\" para executar.")
		return
	}

	printResponse(response)
}

func confirmPendingAction(orchestrator *agent.Orchestrator, user policy.AuthenticatedUser) {
	response, err := orchestrator.ConfirmPendingAction(user)
	if err != nil {
		fmt.Println("Confirmacao falhou:", err)
		return
	}

	fmt.Printf("Acao confirmada e executada: %+v\n", response)
}

func printResponse(response any) {
	switch value := response.(type) {
	case []knowledge.KnowledgeResult:
		printKnowledgeResponse(value)
	default:
		fmt.Printf("Resposta: %+v\n", response)
	}
}

func printKnowledgeResponse(results []knowledge.KnowledgeResult) {
	if len(results) == 0 {
		fmt.Println("Resposta: nenhuma informacao encontrada na base de conhecimento.")
		return
	}

	result := results[0]
	fmt.Println("Resposta:")
	fmt.Println(strings.TrimSpace(result.Answer))

	if result.Source != "" {
		fmt.Println()
		fmt.Println("Fonte:", result.Source)
	}
}

func isConfirmation(message string) bool {
	normalized := strings.ToLower(message)
	return normalized == "sim" ||
		normalized == "confirmo" ||
		normalized == "confirmar" ||
		normalized == "confirmo o pix"
}

func isExitCommand(message string) bool {
	normalized := strings.ToLower(message)
	return normalized == "sair" || normalized == "exit" || normalized == "quit"
}
