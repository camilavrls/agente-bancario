package agent

import (
	"strings"
)

type Classifier interface {
	Classify(input string) (string, error)
}

type MockClassifier struct{}

func (c MockClassifier) Classify(input string) (string, error) {

	switch {

	case contains(input, "pix"):
		return "CRITICAL", nil

	case contains(input, "limite"):
		return "TOOL", nil

	case contains(input, "joao"):
		return "BLOCKED", nil

	default:
		return "RAG", nil
	}
}

func contains(text, substr string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}
