package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "github.com/joho/godotenv"

    "code-gen-cli/agents/fs"
    "code-gen-cli/agents/code-gen"
    "code-gen-cli/agents/ask"
)

func readMultiline(scanner *bufio.Scanner) string {
    var lines []string
    for {
        if !scanner.Scan() {
            break
        }
        line := scanner.Text()
        if strings.TrimSpace(line) == "---" {
            break
        } else if  line == "ask"|| line == "quit" ||  line == "exit"{
            return line
        }  
        lines = append(lines, line)
    }
    return strings.Join(lines, "\n")
}


func main() {
    if err := godotenv.Load(); err != nil {
        fmt.Println("No .env file found")
    }
    fmt.Println("Welcome to LLM Codegen Agent!")
    fmt.Println("Type a prompt and I'll generate code for you.")
    fmt.Println("Type 'quit' or 'exit' to leave.")
    fmt.Println("Type ask to ask a question with no code generation")
    fmt.Println()
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Buffer(make([]byte, 1024), 1024*1024) 

    for {
        fmt.Print("You (end with ---):")
        input := readMultiline(scanner)

        if input == "" {
            continue
        }

        if input == "quit" || input == "exit" {
            fmt.Println("Goodbye!")
            break
        }

        if input == "ask"{
            fmt.Print("You ask (end with ---):")
            input := readMultiline(scanner)

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