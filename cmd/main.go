package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "code-gen-cli/agents/fs"
    "code-gen-cli/agents/code-gen"
    "code-gen-cli/agents/ask"
)

func main() {
    fmt.Println("Welcome to LLM Codegen Agent!")
    fmt.Println("Type a prompt and I'll generate code for you.")
    fmt.Println("Type 'quit' or 'exit' to leave.")
    fmt.Println("Type ask to ask a question with no code generation")
    fmt.Println()

    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("You: ")
        if !scanner.Scan() {
            break
        }

        input := strings.TrimSpace(scanner.Text())
        if input == "" {
            continue
        }

        if input == "quit" || input == "exit" {
            fmt.Println("Goodbye!")
            break
        }

        if input == "ask"{
            fmt.Print("Ask: ")

            if !scanner.Scan() {
                break
            }
            input = strings.TrimSpace(scanner.Text())
            ask.GenerateAnswer(input)
            fmt.Println()
            fmt.Println(strings.Repeat("-", 50))
            fmt.Println()
            continue
        }

        fmt.Println("Generating code...")

        codeMap, err := code.GenerateCode(input)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        outputDir := "./example"
        err = fs.WriteFiles(outputDir, codeMap)
        if err != nil {
            fmt.Printf("Failed to write files: %v\n", err)
            continue
        }

        fmt.Printf("Code written to '%s'\n", outputDir)
        fmt.Println(strings.Repeat("-", 50))
        fmt.Println()
        
    }
}
