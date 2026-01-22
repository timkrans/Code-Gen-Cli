package providers

import (
    "bytes"
    "encoding/json"
    "net/http"
	"io"
    "errors"
	"code-gen-cli/internal/llm"
)

type HuggingFaceClient struct {
    BaseURL string
    APIKey  string
    Model   string
}

func NewHuggingFace(cfg llm.Config) *HuggingFaceClient {
    return &HuggingFaceClient{
        BaseURL: cfg.HFBaseURL,
        APIKey:  cfg.HFAPIKey,
        Model:   cfg.Model,
    }
}


func (c *HuggingFaceClient) Generate(prompt string) (*http.Response, error) {
    body := map[string]string{"inputs": prompt}
    b, _ := json.Marshal(body)

    req, _ := http.NewRequest(
        "POST",
        c.BaseURL+"/models/"+c.Model,
        bytes.NewBuffer(b),
    )

    req.Header.Set("Authorization", "Bearer "+c.APIKey)
    req.Header.Set("Content-Type", "application/json")

    return http.DefaultClient.Do(req)
}

func DecodeHuggingFace(body io.Reader) (string, error) {
    var res []struct {
        GeneratedText string `json:"generated_text"`
    }

    if err := json.NewDecoder(body).Decode(&res); err != nil {
        return "", err
    }

    if len(res) == 0 {
        return "", errors.New("huggingface: empty response")
    }

    return res[0].GeneratedText, nil
}

