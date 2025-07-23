# ContextFind

A CLI tool for selecting files with fzf and outputting their contents. Includes context management for saving and reusing file selections.

![cf demo](./demonstration.gif)

## Prerequisites

- Go 1.24+ installed
- [fzf](https://github.com/junegunn/fzf) (install via `brew install fzf` on macOS)
- [markitdown](https://github.com/microsoft/markitdown) for binary file conversion

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
cf [directory]         # Select files from current or specified directory
cf | pbcopy            # Copy output to clipboard
```

### Context Management

```bash
cf save [name]         # Save current file selection as named context
cf from [name]         # Load files from saved context
cf update [name]       # Update existing context with new files
cf delete [name]       # Delete saved contexts
```

If no name is provided, fzf will open for interactive selection. Context configurations are stored in `.contextfind.toml` files.
