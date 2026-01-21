package ctx

import (
    "fmt"
    "net/http"
    "io"
    "os"
    "code-gen-cli/internal/llm"
    "code-gen-cli/internal/llm/factory"
)

func GenerateContext(prompt string) error {
    cfg := llm.LoadConfig()
    client := factory.NewClient(cfg)

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
        output, err = client.decodeOllamaStream(resp.Body)
        if err != nil {
            return err
        }
    case "openai":
        output, err = client.decodeOpenAI(resp.Body)

    case "anthropic":
        output, err = client.decodeAnthropic(resp.Body)

    case "google":
        output, err = client.decodeGemini(resp.Body)

    case "huggingface":
        output, err = client.decodeHuggingFace(resp.Body)

    default:
        return fmt.Errorf("unsupported provider: %s", provider)
    }

    if err != nil {
        return err
    }

    fmt.Print(output)
    return nil
}
