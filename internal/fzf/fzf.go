package fzf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

func isBatAvailable() bool {
	_, err := exec.LookPath("bat")
	return err == nil
}

func isFdAvailable() bool {
	_, err := exec.LookPath("fd")
	return err == nil
}

type fileInfo struct {
	path string
	mod  time.Time
}

func SelectFiles(dir string) ([]string, error) {
	var paths []string
	var err error

	if isFdAvailable() {
		cmd := exec.Command("fd", "--type", "f", "--strip-cwd-prefix")
		cmd.Dir = dir
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("error running fd: %w", err)
		}
		selectedStr := strings.TrimSpace(string(output))
		if selectedStr == "" {
			return nil, fmt.Errorf("no files found")
		}
		paths = strings.Split(selectedStr, "\n")
	} else {
		var files []fileInfo
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
				files = append(files, fileInfo{path: path, mod: info.ModTime()})
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error walking directory: %w", err)
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("no files found")
		}
		// Extract paths after collecting (sorting happens later)
		paths = make([]string, len(files))
		for i, f := range files {
			paths[i] = f.path
		}
	}

	// Now sort all paths by mod time, recent first
	var fileInfos []fileInfo
	for _, p := range paths {
		fullPath := p
		if isFdAvailable() {
			fullPath = filepath.Join(dir, p)
		}
		info, statErr := os.Stat(fullPath)
		if statErr == nil {
			fileInfos = append(fileInfos, fileInfo{path: p, mod: info.ModTime()})
		}
	}
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].mod.After(fileInfos[j].mod)
	})
	sortedPaths := make([]string, len(fileInfos))
	for i, fi := range fileInfos {
		sortedPaths[i] = fi.path
	}

	if len(sortedPaths) == 0 {
		return nil, fmt.Errorf("no files found")
	}

	return runFzf(sortedPaths, true)
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

	var previewCmd string
	if isBatAvailable() {
		previewCmd = "sh -c 'bat --color=always --style=numbers --line-range=:500 \"{}\" 2>/dev/null || cat \"{}\" 2>/dev/null || echo \"Preview not available\"'"
	} else {
		previewCmd = "sh -c 'cat \"{}\" 2>/dev/null || echo \"Preview not available\"'"
	}
	args = append(args, "--preview", previewCmd)
	args = append(args, "--preview-window", "right:50%:wrap")

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
