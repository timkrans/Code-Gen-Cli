package ctx

import (
    "encoding/json"
    "strings"
    "fmt"
    "errors"
    "net/http"
    "io"
    "os"

    "code-gen-cli/internal/llm"
)

func GenerateContext(prompt string) error {
    cfg := llm.LoadConfig()
    client := llm.NewClient(cfg)

    resp, err := client.Generate(prompt)
    if err != nil {
        return fmt.Errorf("failed to connect to LLM: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("LLM returned %d: %s", resp.StatusCode, string(body))
    }

    provider := os.Getenv("LLM_PROVIDER")
    var output string

    switch provider {

    case "ollama":
        output, err = decodeOllamaStream(resp.Body)
        if err != nil {
            return err
        }
    case "openai":
        output, err = decodeOpenAI(resp.Body)

    case "anthropic":
        output, err = decodeAnthropic(resp.Body)

    case "google":
        output, err = decodeGemini(resp.Body)

    case "huggingface":
        output, err = decodeHuggingFace(resp.Body)

    default:
        return fmt.Errorf("unsupported provider: %s", provider)
    }

    if err != nil {
        return err
    }

    fmt.Print(output)
    return nil
}
