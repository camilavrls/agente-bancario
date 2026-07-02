package knowledge

import (
	"fmt"
	"os"
	"strings"
)

const (
	knowledgeProviderMemory = "memory"
	knowledgeProviderChroma = "chroma"
)

func SearchKnowledgeBase(query string) ([]KnowledgeResult, error) {
	provider := strings.ToLower(os.Getenv("KNOWLEDGE_PROVIDER"))
	if provider == "" {
		provider = knowledgeProviderMemory
	}

	switch provider {
	case knowledgeProviderMemory:
		return SearchKnowledgeBaseMemory(query)
	// case knowledgeProviderChroma:
	// 	return SearchKnowledgeBaseChroma(query)
	default:
		return nil, fmt.Errorf("KNOWLEDGE_PROVIDER invalido: %s", provider)
	}
}
