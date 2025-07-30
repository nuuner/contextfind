<h1>
  <img src="logo_100x100.png" alt="ContextFind Logo" width="100" height="100" align="left">
  ContextFind
</h1>

A CLI tool for selecting files with fzf and outputting their contents. Includes context management for saving and reusing file selections.

![cf demo](./demonstration.gif)

## Prerequisites

- Go 1.24+ installed
- [fzf](https://github.com/junegunn/fzf) (install via `brew install fzf` on macOS)
- [markitdown](https://github.com/microsoft/markitdown) for binary file conversion
- [bat](https://github.com/sharkdp/bat) (optional, for syntax highlighting in preview - install via `brew install bat` on macOS)

## Installation

Install directly via Go:

```bash
go install github.com/nuuner/contextfind/cmd/cf@latest
```

Or clone and build locally:

```bash
git clone https://github.com/nuuner/contextfind.git
cd contextfind
go build -o cf ./cmd/cf
```

## Usage

### Basic File Selection

```bash
cf [directory]         # Select files from current or specified directory, automatically saves as last selection
cf | pbcopy            # Copy output to clipboard
```

### Last Selection

```bash
cf last                # Output files from the last selection
cf save-last [name]    # Save the last selection as a named context (prompts if no name provided)
```

### Context Management

```bash
cf save [name]         # Save current file selection as named context (prompts if no name provided)
cf from [name]         # Load files from saved context
cf update [name]       # Update existing context with new files
cf delete [name]       # Delete saved contexts
```

If no name is provided for commands like `from`, `update`, or `delete`, fzf will open for interactive selection. Context configurations are stored in `.contextfind.toml` files.
