package llm

import (
	"agente-bancario/mcp"
	"agente-bancario/policy"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultGeminiModel = "gemini-2.5-flash"

type geminiGenerateContentRequest struct {
	SystemInstruction geminiContent          `json:"system_instruction"`
	Contents          []geminiContent        `json:"contents"`
	GenerationConfig  geminiGenerationConfig `json:"generationConfig"`
}

type geminiGenerationConfig struct {
	Temperature      float64 `json:"temperature"`
	ResponseMimeType string  `json:"response_mime_type"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerateContentResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

type rawToolCall struct {
	Name      string         `json:"Name"`
	Arguments map[string]any `json:"Arguments"`
}

func toolPlanningSystemInstruction() string {
	return `Voce e um planejador de ferramentas para um agente bancario.
Sua tarefa e escolher uma das tools disponiveis.
Nao execute autorizacao.
Nao invente tools.
Use somente os parametros declarados na tool.
Responda somente JSON valido, sem markdown.
O JSON deve seguir exatamente este formato:
{
  "Name": "nome_da_tool",
  "Arguments": {
    "parametro": "valor"
  }
}`
}

func buildToolPlanningPrompt(message string, user policy.AuthenticatedUser, availableTools []mcp.ToolDefinition) (string, error) {
	toolsJSON, err := json.MarshalIndent(availableTools, "", "  ")
	if err != nil {
		return "", err
	}

	userJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`Usuario autenticado:
%s

Tools disponiveis:
%s

Contexto de demonstracao:
- Quando o usuario falar de si mesmo, use o customer_id do usuario autenticado.
- Quando o usuario mencionar Joao ou João, o customer_id conhecido para este demo e cust-456.

Mensagem do usuario:
%s`, userJSON, toolsJSON, message), nil
}

func PlanToolCallGemini(message string, user policy.AuthenticatedUser, availableTools []mcp.ToolDefinition) (mcp.ToolCall, error) {
	prompt, err := buildToolPlanningPrompt(message, user, availableTools)
	if err != nil {
		return mcp.ToolCall{}, err
	}

	output, err := callGemini(prompt)
	if err != nil {
		return mcp.ToolCall{}, err
	}

	var raw rawToolCall
	if err := json.Unmarshal([]byte(output), &raw); err != nil {
		return mcp.ToolCall{}, fmt.Errorf("resposta invalida do Gemini: %w; resposta: %s", err, output)
	}

	if raw.Name == "" {
		return mcp.ToolCall{}, fmt.Errorf("Gemini nao retornou o nome da tool")
	}

	arguments := normalizeToolArguments(raw.Arguments)

	return mcp.ToolCall{
		Name:      raw.Name,
		Arguments: arguments,
	}, nil
}

func normalizeToolArguments(arguments map[string]any) map[string]string {
	normalized := map[string]string{}
	for key, value := range arguments {
		normalized[key] = stringifyToolArgument(value)
	}

	return normalized
}

func stringifyToolArgument(value any) string {
	switch typedValue := value.(type) {
	case string:
		return typedValue
	case float64:
		if typedValue == float64(int64(typedValue)) {
			return fmt.Sprintf("%d", int64(typedValue))
		}
		return fmt.Sprintf("%f", typedValue)
	case bool:
		return fmt.Sprintf("%t", typedValue)
	default:
		return fmt.Sprint(typedValue)
	}
}

func callGemini(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("configure GEMINI_API_KEY para chamar o Gemini")
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = defaultGeminiModel
	}

	requestBody := geminiGenerateContentRequest{
		SystemInstruction: geminiContent{
			Parts: []geminiPart{{Text: toolPlanningSystemInstruction()}},
		},
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: prompt}}},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature:      0,
			ResponseMimeType: "application/json",
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", model)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var geminiResponse geminiGenerateContentResponse
	if err := json.Unmarshal(responseBody, &geminiResponse); err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if geminiResponse.Error != nil {
			return "", fmt.Errorf("Gemini retornou erro %d/%s: %s", geminiResponse.Error.Code, geminiResponse.Error.Status, geminiResponse.Error.Message)
		}
		return "", fmt.Errorf("Gemini retornou status %d: %s", resp.StatusCode, string(responseBody))
	}

	if len(geminiResponse.Candidates) == 0 || len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini nao retornou conteudo")
	}

	return strings.TrimSpace(geminiResponse.Candidates[0].Content.Parts[0].Text), nil
}
