package embedding

import (
	"bytes"
	"encoding/json"
	"memoria/internal/config"
	"net/http"
)

type OllamaEmbedder struct {
	BaseURL string
	Model   string
}

func NewOllamaEmbedder() *OllamaEmbedder {
	return &OllamaEmbedder{
		BaseURL: config.Get("OLLAMA_URL", "http://localhost:11434"),
		Model:   config.Get("OLLAMA_MODEL", "nomic-embed-text-v2-moe:latest"),
	}
}

type ollamaReq struct {
	Model string `json:"model"`
	Input string `json:"prompt"`
}

type ollamaResp struct {
	Embedding []float32 `json:"embedding"`
}

func (o *OllamaEmbedder) Embed(text string) ([]float32, error) {

	reqBody := ollamaReq{
		Model: o.Model,
		Input: text,
	}

	b, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		o.BaseURL+"/api/embeddings",
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out ollamaResp
	json.NewDecoder(resp.Body).Decode(&out)

	return out.Embedding, nil
}
