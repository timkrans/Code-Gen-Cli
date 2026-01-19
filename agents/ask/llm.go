package ask

import(
	"encoding/json"
	"strings"
	"code-gen-cli/agents/models"
	"fmt"
	"errors"
	"net/http"
	"bytes"
	"io"
)

func GenerateAnswer(prompt string)(error){
	reqBody := models.OllamaRequest{
        Model:  "llama3.2",
        Prompt: prompt,
        Stream: true,
    }

    reqBytes, err := json.Marshal(reqBody)
    if err != nil {
        return err
    }

    resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(reqBytes))
    if err != nil {
        return fmt.Errorf("failed to connect to Ollama: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Ollama returned non-200: %d", resp.StatusCode)
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
            return fmt.Errorf("failed to decode Ollama stream: %w", err)
        }

        fullResponse.WriteString(chunk.Response)
    }
	fmt.Printf(fullResponse.String())
	return nil
}