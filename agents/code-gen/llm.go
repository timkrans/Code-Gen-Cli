package code

import (
    "bufio"
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "path/filepath"
    "strings"
    "code-gen-cli/agents/models"
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

    reqBody := models.OllamaRequest{
        Model:  "llama3.2",
        Prompt: fullPrompt,
        Stream: true,
    }

    reqBytes, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(reqBytes))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Ollama returned non-200: %d", resp.StatusCode)
    }

    var fullResponse strings.Builder
    decoder := json.NewDecoder(resp.Body)

    for decoder.More() {
        var chunk models.OllamaResponse
        err := decoder.Decode(&chunk)
        if err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            return nil, fmt.Errorf("failed to decode Ollama stream: %w", err)
        }

        fullResponse.WriteString(chunk.Response)
    }
    fmt.Printf(fullResponse.String())
    return parseMultiFileResponse(fullResponse.String()), nil
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
