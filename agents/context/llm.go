package llmCtx

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Builder struct {
	MaxBytes int
	MaxFiles int
	RootDir  string
}

func New(maxBytes int, maxFiles int, root string) *Builder {
	return &Builder{
		MaxBytes: maxBytes,
		MaxFiles: maxFiles,
		RootDir:  root,
	}
}

func (b *Builder) Build(prompt string) string {
	root := b.RootDir
	if _, err := os.Stat(root); os.IsNotExist(err) {
		root, _ = os.Getwd()
	}

	files, _ := b.scanFiles(root)
	ranked := b.rankFiles(files, prompt)
	selected := b.limitFiles(ranked)

	content := b.loadFiles(selected)
	tree := b.buildTree(root, selected)

	return fmt.Sprintf(`
You are working inside a real codebase.

Project Root:
%s

Project Structure:
%s

Relevant Code:
%s

User Prompt:
%s
`, root, tree, content, prompt)
}

func (b *Builder) scanFiles(root string) ([]string, error) {
	var files []string

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			switch info.Name() {
			case ".git", "node_modules", "vendor", "dist", "build":
				return filepath.SkipDir
			}
			return nil
		}

		if isRelevantFile(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, nil
}

func isRelevantFile(path string) bool {
	ext := filepath.Ext(path)

	allowed := map[string]bool{
		".go":   true,
		".js":   true,
		".ts":   true,
		".json": true,
		".md":   true,
		".txt":  true,
		".yaml": true,
		".yml":  true,
	}

	return allowed[ext]
}

func (b *Builder) rankFiles(files []string, prompt string) []string {
	type scored struct {
		path  string
		score int
	}

	var results []scored
	p := strings.ToLower(prompt)
	words := strings.Fields(p)

	for _, f := range files {
		score := 0

		name := strings.ToLower(filepath.Base(f))

		for _, w := range words {
			if strings.Contains(name, w) {
				score += 5
			}
		}

		data, err := os.ReadFile(f)
		if err == nil && len(data) > 0 {
			sample := strings.ToLower(string(data[:min(2048, len(data))]))
			for _, w := range words {
				if strings.Contains(sample, w) {
					score++
				}
			}
		}

		if score > 0 {
			results = append(results, scored{f, score})
		}
	}

	if len(results) == 0 {
		limit := min(len(files), b.MaxFiles)
		return files[:limit]
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	var ranked []string
	for _, r := range results {
		ranked = append(ranked, r.path)
	}

	return ranked
}

func (b *Builder) limitFiles(files []string) []string {
	if len(files) <= b.MaxFiles {
		return files
	}
	return files[:b.MaxFiles]
}

func (b *Builder) loadFiles(files []string) string {
	var builder strings.Builder
	total := 0

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}

		chunks := chunkBytes(data, 1500)

		for i, chunk := range chunks {
			if total+len(chunk) > b.MaxBytes {
				return builder.String()
			}

			builder.WriteString(fmt.Sprintf(
				"\n--- FILE: %s (chunk %d) ---\n",
				f, i+1,
			))
			builder.Write(chunk)

			total += len(chunk)
		}
	}

	return builder.String()
}

func (b *Builder) buildTree(root string, files []string) string {
	var builder strings.Builder
	seen := make(map[string]bool)

	for _, f := range files {
		rel, err := filepath.Rel(root, f)
		if err != nil {
			continue
		}

		dir := filepath.Dir(rel)

		if !seen[dir] {
			builder.WriteString(fmt.Sprintf("📁 %s\n", dir))
			seen[dir] = true
		}

		builder.WriteString(fmt.Sprintf("  └── %s\n", filepath.Base(f)))
	}

	return builder.String()
}

func chunkBytes(data []byte, size int) [][]byte {
	var chunks [][]byte

	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}

	return chunks
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}