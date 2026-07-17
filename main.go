package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"agente-bancario/agent"
	"agente-bancario/llm"
	"agente-bancario/mcp"
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
	answer, answerErr := llm.GenerateFinalAnswer(prompt, call, response, err)
	if answerErr != nil {
		fmt.Println("Erro ao gerar resposta final:", answerErr)
		return
	}

	fmt.Println(answer)
}

func confirmPendingAction(orchestrator *agent.Orchestrator, user policy.AuthenticatedUser) {
	response, err := orchestrator.ConfirmPendingAction(user)
	call := mcp.ToolCall{Name: "create_pix"}

	answer, answerErr := llm.GenerateFinalAnswer("confirmo", call, response, err)
	if answerErr != nil {
		fmt.Println("Erro ao gerar resposta final:", answerErr)
		return
	}

	fmt.Println(answer)
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
