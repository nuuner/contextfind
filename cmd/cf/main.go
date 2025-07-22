package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/nuuner/contextfind/utils"
)

func main() {
	if err := utils.CheckDependencies(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

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
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found.")
		os.Exit(0)
	}

	input := strings.Join(files, "\n")
	cmd := exec.Command("fzf", "--multi")
	cmd.Stdin = strings.NewReader(input)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running fzf: %v\n", err)
		os.Exit(1)
	}

	selectedStr := strings.TrimSpace(string(output))
	if selectedStr == "" {
		fmt.Println("No files selected.")
		os.Exit(0)
	}

	selected := strings.Split(selectedStr, "\n")

	for _, file := range selected {
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
}
