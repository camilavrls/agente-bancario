package knowledge

import "strings"

type memoryDocument struct {
	ID      string
	Source  string
	Content string
	Terms   []string
}

var memoryDocuments = []memoryDocument{
	{
		ID:      "kb-consignado-001",
		Source:  "politica-emprestimo-consignado.md",
		Content: "A taxa do emprestimo consignado para aposentados com bom historico varia a partir de 1,59% ao mes, sujeita a analise de credito.",
		Terms:   []string{"emprestimo", "consignado", "aposentado", "taxa"},
	},
	{
		ID:      "kb-tarifas-001",
		Source:  "tabela-tarifas-conta.md",
		Content: "A conta possui pacote essencial gratuito. Transferencias TED podem ter tarifa de R$ 10,00 fora do pacote contratado.",
		Terms:   []string{"tarifa", "tarifas", "conta", "ted", "pacote"},
	},
	{
		ID:      "kb-limite-001",
		Source:  "politica-limite-cartao.md",
		Content: "Aumento de limite do cartao depende do limite maximo pre-aprovado, historico de pagamento e analise de risco.",
		Terms:   []string{"limite", "cartao", "aumento", "risco"},
	},
	{
		ID:      "kb-pix-001",
		Source:  "seguranca-pix.md",
		Content: "Operacoes PIX sao consideradas criticas e exigem confirmacao explicita do usuario antes da execucao.",
		Terms:   []string{"pix", "seguranca", "confirmacao", "critica"},
	},
}

func SearchKnowledgeBaseMemory(query string) ([]KnowledgeResult, error) {
	normalized := strings.ToLower(query)
	results := []KnowledgeResult{}

	for _, document := range memoryDocuments {
		if documentMatches(normalized, document) {
			results = append(results, KnowledgeResult{
				Answer:     document.Content,
				Source:     document.Source,
				DocumentID: document.ID,
			})
		}
	}

	if len(results) == 0 {
		return []KnowledgeResult{
			{
				Answer:     "Nao encontrei informacoes na base de conhecimento para essa pergunta.",
				Source:     "memory",
				DocumentID: "not_found",
			},
		}, nil
	}

	return results, nil
}

func documentMatches(query string, document memoryDocument) bool {
	for _, term := range document.Terms {
		if strings.Contains(query, term) {
			return true
		}
	}

	return false
}
