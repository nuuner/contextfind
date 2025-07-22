# ContextFind

A tool to quickly find and get contents of files in a directory. Useful for giving context to LLMs.

![cf demo](./demonstration.gif)

## Prerequisites

- Go 1.22+ installed.
- [fzf](https://github.com/junegunn/fzf) (install via `brew install fzf` on macOS).
- [markitdown](https://github.com/microsoft/markitdown) (see [installation instructions](https://github.com/microsoft/markitdown#installation)).

## Installation

Install directly via Go (binary named `cf`):

```bash
go install github.com/nuuner/contextfind/cmd/cf@latest
```

Or clone and install locally:

```bash
git clone https://github.com/nuuner/contextfind.git
cd contextfind/cmd/cf
go install
```

Ensure `$GOPATH/bin` is in your `$PATH` (e.g., add `export PATH=$PATH:$(go env GOPATH)/bin` to `~/.zshrc` or `~/.bash_profile`).

## Usage

Run `cf [directory]` to select and display files.

`cf | pbcopy` will copy the selected files to your clipboard.
