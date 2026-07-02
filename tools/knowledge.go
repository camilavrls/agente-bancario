package tools

import (
	"agente-bancario/knowledge"
	"agente-bancario/mcp"
)

func SearchKnowledgeBaseToolDefinition() mcp.ToolDefinition {
	return mcp.ToolDefinition{
		Name:        "search_knowledge_base",
		Description: "Busca informacoes na base de conhecimento do banco para responder duvidas sobre politicas, tarifas, produtos e FAQ.",
		Parameters: []mcp.ToolParameter{
			{
				Name:        "query",
				Type:        "string",
				Description: "Pergunta ou termo de busca do usuario.",
				Required:    true,
			},
		},
	}
}

func SearchKnowledgeBaseTool(query string) ([]knowledge.KnowledgeResult, error) {
	return knowledge.SearchKnowledgeBase(query)
}
