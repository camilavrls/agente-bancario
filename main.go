package main

import (
	"agente-bancario/agent"
	"agente-bancario/llm"
	"fmt"
)

func main() {

	ag := agent.Agent{
		Classifier: llm.GeminiClassifier{},
		Router:     agent.SimpleRouter{},
	}

	resp, err := ag.Handle("quero fazer um pix de 20000")

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
