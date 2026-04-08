package brain

import (
	"fmt"
	"strings"

	"code-gen-cli/agents/ask"
	code "code-gen-cli/agents/code-gen"
	"code-gen-cli/agents/fs"
	"code-gen-cli/agents/context"
)

type TaskType int

const (
	TaskAsk TaskType = iota
	TaskCode
)

type Brain struct {
	ctx       *llmCtx.Builder
	outputDir string
}

func New() *Brain {
	return &Brain{
		ctx:       llmCtx.New(20000, 8,"./example"),
		outputDir: "./example",
	}
}


func (b *Brain) Run(prompt string) error {
	task := b.classify(prompt)

	switch task {
	case TaskAsk:
		return b.runAsk(prompt)

	case TaskCode:
		return b.runCode(prompt)
	}

	return nil
}

func (b *Brain) runAsk(prompt string) error {
	enriched := b.ctx.Build(prompt)

	err := ask.GenerateAnswer(enriched)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println()

	return nil
}

func (b *Brain) runCode(prompt string) error {
	fmt.Println("Building context...")
	enriched := b.ctx.Build(prompt)

	fmt.Println("Generating code...")
	codeMap, err := code.GenerateCode(enriched)
	if err != nil {
		return err
	}

	fmt.Println("Writing files...")
	err = fs.WriteFiles(b.outputDir, codeMap)
	if err != nil {
		return err
	}

	fmt.Printf("Code written to '%s'\n", b.outputDir)
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println()

	return nil
}

func (b *Brain) classify(prompt string) TaskType {
	p := strings.ToLower(prompt)

	codeSignals := []string{
		"write", "build", "create", "generate",
		"api", "server", "function", "golang",
		"cli", "app", "service",
	}

	for _, s := range codeSignals {
		if strings.Contains(p, s) {
			return TaskCode
		}
	}

	return TaskAsk
}