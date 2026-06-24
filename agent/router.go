package agent

import "fmt"

type SimpleRouter struct{}

type Router interface {
    Route(intent string, input string) (string, error)
}

func (r SimpleRouter) Route(intent string, input string) (string, error) {

    switch intent {

    case "RAG":
        return "Resposta via RAG (mock)", nil

    case "TOOL":
        return "Execução de tool (mock MCP)", nil

    case "CRITICAL":
        return "Fluxo crítico: pedir confirmação", nil

    case "BLOCKED":
        return "Ação negada por política", nil

    default:
        return "", fmt.Errorf("intent desconhecido")
    }
}