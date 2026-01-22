package providers

import (
    "bytes"
    "encoding/json"
    "net/http"
	"io"
    "errors"
	"code-gen-cli/internal/llm"
)

type OpenAI struct {
    baseURL string
    apiKey  string
    model   string
}

func NewOpenAI(cfg llm.Config) llm.LLMClient {
    return &OpenAI{
        baseURL: cfg.OpenAIBaseURL,
        apiKey:  cfg.OpenAIAPIKey,
        model:   cfg.Model,
    }
}

func (o *OpenAI) Generate(prompt string) (*http.Response, error) {
    body := map[string]interface{}{
        "model": o.model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
    }

    b, _ := json.Marshal(body)

    req, _ := http.NewRequest(
        "POST",
        o.baseURL+"/chat/completions",
        bytes.NewBuffer(b),
    )

    req.Header.Set("Authorization", "Bearer "+o.apiKey)
    req.Header.Set("Content-Type", "application/json")

    return http.DefaultClient.Do(req)
}

func DecodeOpenAI(body io.Reader) (string, error) {
    var res struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }

    if err := json.NewDecoder(body).Decode(&res); err != nil {
        return "", err
    }

    if len(res.Choices) == 0 {
        return "", errors.New("openai: empty response")
    }

    return res.Choices[0].Message.Content, nil
}
