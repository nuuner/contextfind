package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	contextpkg "github.com/nuuner/contextfind/internal/context"
	"github.com/nuuner/contextfind/internal/fzf"
	"github.com/urfave/cli/v3"
)

func DefaultAction(_ context.Context, cmd *cli.Command) error {
	dir := "."
	if cmd.Args().Len() > 0 {
		dir = cmd.Args().First()
	}

	files, err := fzf.SelectFiles(dir)
	if err != nil {
		return fmt.Errorf("file selection failed: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No files selected.")
		return nil
	}

	return fzf.OutputFiles(files)
}

func SaveAction(_ context.Context, cmd *cli.Command) error {
	dir := "."
	if cmd.Args().Len() > 1 {
		dir = cmd.Args().Get(1)
	}

	files, err := fzf.SelectFiles(dir)
	if err != nil {
		return fmt.Errorf("file selection failed: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No files selected.")
		return nil
	}

	var name string
	if cmd.Args().Len() > 0 {
		name = cmd.Args().First()
	} else {
		fmt.Print("Enter context name: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			name = strings.TrimSpace(scanner.Text())
		}
		if name == "" {
			return fmt.Errorf("context name cannot be empty")
		}
	}

	cf, err := contextpkg.LoadContextFile()
	if err != nil {
		return fmt.Errorf("failed to load context file: %w", err)
	}

	cf.AddOrUpdateContext(name, files)
	if err := cf.Save(); err != nil {
		return fmt.Errorf("failed to save context: %w", err)
	}

	fmt.Printf("Context '%s' saved with %d files.\n", name, len(files))
	return nil
}

func FromAction(_ context.Context, cmd *cli.Command) error {
	cf, err := contextpkg.LoadContextFile()
	if err != nil {
		return fmt.Errorf("failed to load context file: %w", err)
	}

	if cf.IsEmpty() {
		return fmt.Errorf("no saved contexts found")
	}

	var selectedContext *contextpkg.Context
	if cmd.Args().Len() > 0 {
		name := cmd.Args().First()
		var found bool
		selectedContext, found = cf.GetContext(name)
		if !found {
			return fmt.Errorf("context '%s' not found", name)
		}
	} else {
		names := cf.GetContextNames()
		selected, err := fzf.SelectFromList(names, false)
		if err != nil {
			return fmt.Errorf("context selection failed: %w", err)
		}
		if len(selected) == 0 {
			fmt.Println("No context selected.")
			return nil
		}
		selectedContext, _ = cf.GetContext(selected[0])
	}

	return fzf.OutputFiles(selectedContext.Files)
}

func DeleteAction(_ context.Context, cmd *cli.Command) error {
	cf, err := contextpkg.LoadContextFile()
	if err != nil {
		return fmt.Errorf("failed to load context file: %w", err)
	}

	if cf.IsEmpty() {
		return fmt.Errorf("no saved contexts found")
	}

	if cmd.Args().Len() > 0 {
		name := cmd.Args().First()
		if cf.DeleteContext(name) {
			if err := cf.Save(); err != nil {
				return fmt.Errorf("failed to save context file: %w", err)
			}
			fmt.Printf("Context '%s' deleted.\n", name)
		} else {
			return fmt.Errorf("context '%s' not found", name)
		}
	} else {
		names := cf.GetContextNames()
		selected, err := fzf.SelectFromList(names, true)
		if err != nil {
			return fmt.Errorf("context selection failed: %w", err)
		}
		if len(selected) == 0 {
			fmt.Println("No contexts selected for deletion.")
			return nil
		}

		for _, name := range selected {
			cf.DeleteContext(name)
		}
		if err := cf.Save(); err != nil {
			return fmt.Errorf("failed to save context file: %w", err)
		}
		fmt.Printf("Deleted %d context(s).\n", len(selected))
	}

	return contextpkg.DeleteContextFileIfEmpty()
}

func UpdateAction(_ context.Context, cmd *cli.Command) error {
	cf, err := contextpkg.LoadContextFile()
	if err != nil {
		return fmt.Errorf("failed to load context file: %w", err)
	}

	if cf.IsEmpty() {
		return fmt.Errorf("no saved contexts found")
	}

	var selectedContext *contextpkg.Context
	if cmd.Args().Len() > 0 {
		name := cmd.Args().First()
		var found bool
		selectedContext, found = cf.GetContext(name)
		if !found {
			return fmt.Errorf("context '%s' not found", name)
		}
	} else {
		names := cf.GetContextNames()
		selected, err := fzf.SelectFromList(names, false)
		if err != nil {
			return fmt.Errorf("context selection failed: %w", err)
		}
		if len(selected) == 0 {
			fmt.Println("No context selected.")
			return nil
		}
		selectedContext, _ = cf.GetContext(selected[0])
	}

	dir := "."
	files, err := fzf.SelectFiles(dir)
	if err != nil {
		return fmt.Errorf("file selection failed: %w", err)
	}

	cf.AddOrUpdateContext(selectedContext.Name, files)
	if err := cf.Save(); err != nil {
		return fmt.Errorf("failed to save context: %w", err)
	}

	fmt.Printf("Context '%s' updated with %d files.\n", selectedContext.Name, len(files))
	return nil
}