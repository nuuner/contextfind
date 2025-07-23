package fzf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func SelectFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found")
	}

	return runFzf(files, true)
}

func SelectFromList(items []string, multi bool) ([]string, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("no items to select from")
	}
	return runFzf(items, multi)
}

func runFzf(items []string, multi bool) ([]string, error) {
	input := strings.Join(items, "\n")
	
	args := []string{}
	if multi {
		args = append(args, "--multi")
	}
	
	cmd := exec.Command("fzf", args...)
	cmd.Stdin = strings.NewReader(input)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running fzf: %w", err)
	}

	selectedStr := strings.TrimSpace(string(output))
	if selectedStr == "" {
		return []string{}, nil
	}

	return strings.Split(selectedStr, "\n"), nil
}

func OutputFiles(files []string) error {
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file, err)
			continue
		}

		ext := filepath.Ext(file)
		lang := strings.TrimPrefix(ext, ".")
		if lang == "" {
			lang = "text"
		}

		isText := utf8.Valid(content)
		var finalContent string
		var finalLang string

		if isText {
			finalContent = string(content)
			finalLang = lang
		} else {
			cmd := exec.Command("markitdown", file)
			mdBytes, err := cmd.Output()
			if err != nil {
				finalContent = fmt.Sprintf("Error converting %s with markitdown: %v", file, err)
				finalLang = "text"
			} else {
				finalContent = string(mdBytes)
				finalLang = "markdown"
			}
		}

		fmt.Printf("## %s\n\n```%s\n%s\n```\n\n", file, finalLang, finalContent)
	}
	return nil
}