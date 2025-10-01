package fs

import (
    "fmt"
    "os"
    "path/filepath"
)

func WriteFiles(rootDir string, files map[string]string) error {
    for relPath, content := range files {
        absPath := filepath.Join(rootDir, relPath)

        if err := os.MkdirAll(filepath.Dir(absPath), os.ModePerm); err != nil {
            return fmt.Errorf("failed to create dir %s: %w", absPath, err)
        }

        if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
            return fmt.Errorf("failed to write file %s: %w", absPath, err)
        }

        fmt.Printf("Written: %s\n", absPath)
    }

    return nil
}
