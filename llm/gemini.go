package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GeminiClassifier struct{}

type geminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"contents"`
	} `json:"candidates"`
}

func (g GeminiClassifier) Classify(input string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	prompt := `
		Classifique a intenção em:
		RAG, TOOL, CRITICAL, BLOCKED.

		Regras:
		- informação → RAG
		- ação conta → TOOL
		- pix/transferência → CRITICAL
		- acesso indevido → BLOCKED

		Retorne apenas a categoria.

		Input: ` + input

	body := geminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest(
		"POST",
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key="+apiKey,
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyResp, _ := io.ReadAll(resp.Body)
	fmt.Println("RAW RESPONSE:", string(bodyResp))

	var result geminiResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("no response")
	}

	text := result.Candidates[0].Content.Parts[0].Text

	return text, nil
}
