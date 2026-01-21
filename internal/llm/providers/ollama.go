package providers

import (
    "bytes"
    "encoding/json"
    "net/http"

    "code-gen-cli/internal/llm"
)

type Ollama struct {
    baseURL string
    model   string
}

func NewOllama(cfg llm.Config) llm.LLMClient {
    return &Ollama{
        baseURL: cfg.OllamaBaseURL,
        model:   cfg.Model,
    }
}

func (o *Ollama) Generate(prompt string) (*http.Response, error) {
    body := map[string]interface{}{
        "model":  o.model,
        "prompt": prompt,
        "stream": true,
    }

    b, _ := json.Marshal(body)

    return http.Post(
        o.baseURL+"/api/generate",
        "application/json",
        bytes.NewBuffer(b),
    )
}

func decodeOllamaStream(body io.Reader) (string, error) {
    var result strings.Builder
    decoder := json.NewDecoder(body)

    for {
        var chunk models.OllamaResponse
        err := decoder.Decode(&chunk)
        if err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            return "", fmt.Errorf("ollama decode error: %w", err)
        }

        result.WriteString(chunk.Response)
    }

    return result.String(), nil
}
