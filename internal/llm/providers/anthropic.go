package providers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "io"
    "errors"
	"code-gen-cli/internal/llm"
)

type Anthropic struct {
    apiKey  string
    baseURL string
    model   string
}

func NewAnthropic(cfg llm.Config) llm.LLMClient {
    return &Anthropic{
        apiKey:  cfg.AnthropicAPIKey,
        baseURL: cfg.AnthropicBaseURL,
        model:   cfg.Model,
    }
}

func (a *Anthropic) Generate(prompt string) (*http.Response, error) {
    body := map[string]interface{}{
        "model": a.model,
        "max_tokens": 1024,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
    }

    b, _ := json.Marshal(body)

    req, _ := http.NewRequest(
        "POST",
        a.baseURL+"/v1/messages",
        bytes.NewBuffer(b),
    )

    req.Header.Set("x-api-key", a.apiKey)
    req.Header.Set("anthropic-version", "2023-06-01")
    req.Header.Set("Content-Type", "application/json")

    return http.DefaultClient.Do(req)
}

func DecodeAnthropic(body io.Reader) (string, error) {
    var res struct {
        Content []struct {
            Text string `json:"text"`
        } `json:"content"`
    }

    if err := json.NewDecoder(body).Decode(&res); err != nil {
        return "", err
    }

    if len(res.Content) == 0 {
        return "", errors.New("anthropic: empty response")
    }

    return res.Content[0].Text, nil
}
