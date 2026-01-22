package code

import (
    "bufio"
    "fmt"
    "net/http"
    "path/filepath"
    "strings"
    "code-gen-cli/internal/llm"
    "code-gen-cli/internal/llm/providers"
    "code-gen-cli/internal/llm/factory"
    "os"
)


func GenerateCode(prompt string) (map[string]string, error) {
    fmt.Println("Connecting to Ollama...")

    fullPrompt := `You are a code generation assistant.
Generate Go code based on the following prompt.
Split the output into files using this format and NOTHING other then this format starting with no spaces:

/// FILE: <relative_path>
<code>

Prompt:
` + prompt

    cfg := llm.LoadConfig()
    client := factory.NewClient(cfg)

    resp, err := client.Generate(fullPrompt)
    if err != nil {
        return fmt.Errorf("failed to connect to LLM: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Ollama returned non-200: %d", resp.StatusCode)
    }

    provider := os.Getenv("LLM_PROVIDER")
    var output string

    switch provider {

    case "ollama":
        output, err = providers.DecodeOllamaStream(resp.Body)
        if err != nil {
            return nil, err
        }

    case "openai":
        output, err = providers.DecodeOpenAI(resp.Body)

    case "anthropic":
        output, err = providers.DecodeAnthropic(resp.Body)

    case "google":
        output, err = providers.DecodeGemini(resp.Body)

    case "huggingface":
        output, err = providers.DecodeHuggingFace(resp.Body)

    default:
        return nil, fmt.Errorf("unsupported provider: %s", provider)
    }

    if err != nil {
        return nil, err
    }

    fmt.Print(output)
    return parseMultiFileResponse(output), nil
}


func parseMultiFileResponse(output string) map[string]string {
    files := make(map[string]string)
    var currentFile string
    var contentBuilder strings.Builder

    scanner := bufio.NewScanner(strings.NewReader(output))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "/// FILE:") {
            if currentFile != "" {
                files[filepath.Clean(currentFile)] = contentBuilder.String()
                contentBuilder.Reset()
            }

            currentFile = strings.TrimSpace(strings.TrimPrefix(line, "/// FILE:"))
        } else if currentFile != "" {
            contentBuilder.WriteString(line + "\n")
        }
    }

    if currentFile != "" {
        files[filepath.Clean(currentFile)] = contentBuilder.String()
    }

    return files
}
