package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nuuner/contextfind/internal/commands"
	"github.com/nuuner/contextfind/utils"
	"github.com/urfave/cli/v3"
)

func main() {
	if err := utils.CheckDependencies(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := &cli.Command{
		Name:   "cf",
		Usage:  "Context finder - select files with fzf and output to context",
		Action: commands.DefaultAction,
		Commands: []*cli.Command{
			{
				Name:      "save",
				Usage:     "Save selected files as a named context",
				ArgsUsage: "[name]",
				Action:    commands.SaveAction,
			},
			{
				Name:      "from",
				Usage:     "Load and output files from a saved context",
				ArgsUsage: "[name]",
				Action:    commands.FromAction,
			},
			{
				Name:      "delete",
				Usage:     "Delete saved contexts",
				ArgsUsage: "[name]",
				Action:    commands.DeleteAction,
			},
			{
				Name:      "update",
				Usage:     "Update an existing context with new file selection",
				ArgsUsage: "[name]",
				Action:    commands.UpdateAction,
			},
			{
				Name:   "last",
				Usage:  "Output files from the last selection",
				Action: commands.LastAction,
			},
			{
				Name:      "save-last",
				Usage:     "Save the last selection as a named context",
				ArgsUsage: "[name]",
				Action:    commands.SaveLastAction,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
