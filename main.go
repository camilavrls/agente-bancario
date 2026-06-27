package main

import (
	"fmt"

	"agente-bancario/agent"
	"agente-bancario/llm"
	"agente-bancario/policy"
	"agente-bancario/tools"
)

func main() {
	maria := policy.AuthenticatedUser{
		ID:         "123",
		Name:       "Maria",
		Role:       policy.RoleCustomer,
		CustomerID: "cust-123",
	}

	runPrompt("Maria acessando o proprio perfil", maria, "quero consultar meu perfil")
	runPrompt("Maria tentando acessar Joao", maria, "quero consultar o perfil do joao")
}

func runPrompt(label string, user policy.AuthenticatedUser, prompt string) {
	availableTools := tools.ListTools()

	fmt.Println("Cenario:", label)
	fmt.Println("Prompt recebido:", prompt)

	call, err := llm.PlanToolCall(prompt, user, availableTools)
	if err != nil {
		fmt.Println("Erro ao planejar tool call:", err)
		return
	}

	fmt.Printf("Tool call planejada: %+v\n", call)

	response, err := agent.HandleToolCall(user, call)
	if err != nil {
		fmt.Println("Execucao negada ou falhou:", err)
		fmt.Println()
		return
	}

	fmt.Printf("Execucao permitida: %+v\n\n", response)
}
