# GoWebp

A cross-platform batch image-to-WebP converter with a modern GUI, written in Go.

## Features

- Drag-and-drop or browse to select images/folders
- Batch convert to WebP with preserved folder structure
- Adjustable quality and concurrency
- Progress bar and log pane
- Cross-platform: Windows & Linux

## Installation

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Fyne](https://fyne.io/) (Go module will fetch automatically)

### Build

```bash
go mod tidy
go build -o GoWebp-windows.exe ./cmd/GoWebp
GOOS=linux GOARCH=amd64 go build -o GoWebp-linux ./cmd/GoWebp