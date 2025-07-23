package context

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const ContextFileName = ".contextfind.toml"

type ContextFile struct {
	Contexts []Context `toml:"contexts"`
}

type Context struct {
	Name  string   `toml:"name"`
	Files []string `toml:"files"`
}

func GetContextFilePath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, ContextFileName)
}

func LoadContextFile() (*ContextFile, error) {
	path := GetContextFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &ContextFile{Contexts: []Context{}}, nil
	}

	var cf ContextFile
	if _, err := toml.DecodeFile(path, &cf); err != nil {
		return nil, fmt.Errorf("failed to decode context file: %w", err)
	}
	return &cf, nil
}

func (cf *ContextFile) Save() error {
	path := GetContextFilePath()
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create context file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(cf); err != nil {
		return fmt.Errorf("failed to encode context file: %w", err)
	}
	return nil
}

func (cf *ContextFile) AddOrUpdateContext(name string, files []string) {
	for i, ctx := range cf.Contexts {
		if ctx.Name == name {
			cf.Contexts[i].Files = files
			return
		}
	}
	cf.Contexts = append(cf.Contexts, Context{Name: name, Files: files})
}

func (cf *ContextFile) GetContext(name string) (*Context, bool) {
	for _, ctx := range cf.Contexts {
		if ctx.Name == name {
			return &ctx, true
		}
	}
	return nil, false
}

func (cf *ContextFile) DeleteContext(name string) bool {
	for i, ctx := range cf.Contexts {
		if ctx.Name == name {
			cf.Contexts = append(cf.Contexts[:i], cf.Contexts[i+1:]...)
			return true
		}
	}
	return false
}

func (cf *ContextFile) GetContextNames() []string {
	names := make([]string, len(cf.Contexts))
	for i, ctx := range cf.Contexts {
		names[i] = ctx.Name
	}
	return names
}

func (cf *ContextFile) IsEmpty() bool {
	return len(cf.Contexts) == 0
}

func DeleteContextFileIfEmpty() error {
	cf, err := LoadContextFile()
	if err != nil {
		return err
	}
	if cf.IsEmpty() {
		path := GetContextFilePath()
		if _, err := os.Stat(path); err == nil {
			return os.Remove(path)
		}
	}
	return nil
}