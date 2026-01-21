package providers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "io"
    "errors"
	"fmt"
	"code-gen-cli/internal/llm"
)

type GoogleClient struct {
    APIKey string
    Model  string
}

func NewGoogle(cfg llm.Config) *GoogleClient {
    return &GoogleClient{
        APIKey: cfg.GoogleAPIKey,
        Model:  cfg.GoogleModel,
    }
}

func (c *GoogleClient) Generate(prompt string) (*http.Response, error) {
    body := map[string]interface{}{
        "contents": []map[string]interface{}{
            {
                "parts": []map[string]string{
                    {"text": prompt},
                },
            },
        },
    }

    b, _ := json.Marshal(body)

    url := fmt.Sprintf(
        "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
        c.Model,
        c.APIKey,
    )

    return http.Post(url, "application/json", bytes.NewBuffer(b))
}

func decodeGemini(body io.Reader) (string, error) {
    var res struct {
        Candidates []struct {
            Content struct {
                Parts []struct {
                    Text string `json:"text"`
                } `json:"parts"`
            } `json:"content"`
        } `json:"candidates"`
    }

    if err := json.NewDecoder(body).Decode(&res); err != nil {
        return "", err
    }

    if len(res.Candidates) == 0 ||
       len(res.Candidates[0].Content.Parts) == 0 {
        return "", errors.New("gemini: empty response")
    }

    return res.Candidates[0].Content.Parts[0].Text, nil
}
