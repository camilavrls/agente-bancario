package knowledge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	defaultChromaURL        = "http://localhost:8000"
	defaultChromaCollection = "banking_kb"
)

var chromaEmbeddingKeywords = []string{
	"emprestimo",
	"consignado",
	"aposentado",
	"taxa",
	"tarifa",
	"conta",
	"ted",
	"pix",
	"limite",
	"cartao",
	"aumento",
	"risco",
	"confirmacao",
	"seguranca",
	"cliente",
	"autorizacao",
}

type chromaCollection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type chromaQueryRequest struct {
	QueryEmbeddings [][]float64 `json:"query_embeddings"`
	NResults        int         `json:"n_results"`
	Include         []string    `json:"include"`
}

type chromaQueryResponse struct {
	IDs       [][]string                 `json:"ids"`
	Documents [][]string                 `json:"documents"`
	Metadatas [][]map[string]interface{} `json:"metadatas"`
	Distances [][]float64                `json:"distances"`
}

func SearchKnowledgeBaseChroma(query string) ([]KnowledgeResult, error) {
	collection, err := getChromaCollection(defaultChromaCollection)
	if err != nil {
		return nil, err
	}

	response, err := queryChromaCollection(collection.ID, query)
	if err != nil {
		return nil, err
	}

	return chromaResponseToKnowledgeResults(response), nil
}

func chromaResponseToKnowledgeResults(response chromaQueryResponse) []KnowledgeResult {
	if len(response.Documents) == 0 || len(response.Documents[0]) == 0 {
		return []KnowledgeResult{
			{
				Answer:     "Nao encontrei informacoes na base de conhecimento para essa pergunta.",
				Source:     "chroma",
				DocumentID: "not_found",
			},
		}
	}

	results := make([]KnowledgeResult, 0, len(response.Documents[0]))
	for index, document := range response.Documents[0] {
		result := KnowledgeResult{Answer: document}

		if len(response.IDs) > 0 && len(response.IDs[0]) > index {
			result.DocumentID = response.IDs[0][index]
		}

		if len(response.Metadatas) > 0 && len(response.Metadatas[0]) > index {
			result.Source = metadataValue(response.Metadatas[0][index], "source")
		}

		results = append(results, result)
	}

	return results
}

func getChromaCollection(collectionName string) (chromaCollection, error) {
	baseURL := strings.TrimRight(os.Getenv("CHROMA_URL"), "/")
	if baseURL == "" {
		baseURL = defaultChromaURL
	}

	client := &http.Client{Timeout: 10 * time.Second}
	endpoint := fmt.Sprintf(
		"%s/api/v2/tenants/default_tenant/databases/default_database/collections",
		baseURL,
	)

	var collections []chromaCollection
	if err := getJSON(client, endpoint, &collections); err != nil {
		return chromaCollection{}, fmt.Errorf("erro ao listar collections no Chroma: %w", err)
	}

	for _, collection := range collections {
		if collection.Name == collectionName {
			return collection, nil
		}
	}

	return chromaCollection{}, fmt.Errorf("collection %q nao encontrada no Chroma", collectionName)
}

func queryChromaCollection(collectionID string, query string) (chromaQueryResponse, error) {
	baseURL := strings.TrimRight(os.Getenv("CHROMA_URL"), "/")
	if baseURL == "" {
		baseURL = defaultChromaURL
	}

	client := &http.Client{Timeout: 10 * time.Second}
	endpoint := fmt.Sprintf(
		"%s/api/v2/tenants/default_tenant/databases/default_database/collections/%s/query",
		baseURL,
		url.PathEscape(collectionID),
	)

	request := chromaQueryRequest{
		QueryEmbeddings: [][]float64{embedForChroma(query)},
		NResults:        1,
		Include:         []string{"documents", "metadatas", "distances"},
	}

	var response chromaQueryResponse
	if err := postJSON(client, endpoint, request, &response); err != nil {
		return chromaQueryResponse{}, fmt.Errorf("erro ao consultar Chroma: %w", err)
	}

	return response, nil
}

func embedForChroma(text string) []float64 {
	normalized := normalizeEmbeddingText(text)
	vector := make([]float64, 0, len(chromaEmbeddingKeywords))

	var squaredSum float64
	for _, keyword := range chromaEmbeddingKeywords {
		value := float64(strings.Count(normalized, keyword))
		vector = append(vector, value)
		squaredSum += value * value
	}

	norm := math.Sqrt(squaredSum)
	if norm == 0 {
		return vector
	}

	for index := range vector {
		vector[index] = vector[index] / norm
	}

	return vector
}

func normalizeEmbeddingText(text string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "ã", "a", "â", "a",
		"é", "e", "ê", "e",
		"í", "i",
		"ó", "o", "õ", "o", "ô", "o",
		"ú", "u",
		"ç", "c",
	)

	return replacer.Replace(strings.ToLower(text))
}

func getJSON(client *http.Client, endpoint string, target interface{}) error {
	response, err := client.Get(endpoint)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("status %d: %s", response.StatusCode, string(body))
	}

	return json.NewDecoder(response.Body).Decode(target)
}

func metadataValue(metadata map[string]interface{}, key string) string {
	value, ok := metadata[key]
	if !ok {
		return ""
	}

	return fmt.Sprint(value)
}

func postJSON(client *http.Client, endpoint string, payload interface{}, target interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	response, err := client.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("status %d: %s", response.StatusCode, string(body))
	}

	return json.NewDecoder(response.Body).Decode(target)
}
