# Command - Better Cmd for Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-zoox/command)](https://pkg.go.dev/github.com/go-zoox/command)
[![Build Status](https://github.com/go-zoox/command/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/go-zoox/command/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-zoox/command)](https://goreportcard.com/report/github.com/go-zoox/command)
[![Coverage Status](https://coveralls.io/repos/github/go-zoox/command/badge.svg?branch=master)](https://coveralls.io/github/go-zoox/command?branch=master)
[![GitHub issues](https://img.shields.io/github/issues/go-zoox/command.svg)](https://github.com/go-zoox/command/issues)
[![Release](https://img.shields.io/github/tag/go-zoox/command.svg?label=Release)](https://github.com/go-zoox/command/tags)

## Installation
To install the package, run:
```bash
go get -u github.com/go-zoox/command
```

## Getting Started

```go
// hello world
cmd := &command.Command{
  Script: `echo "hello world"`,
}
if err := cmd.Run(); err != nil {
  log.Fatal("Failed to run command: %s", err)
}
```

```go
// multiline scripts
	cmd := &Command{
		Script: `
echo 1
echo 2
echo 3
`,
	}
if err := cmd.Run(); err != nil {
  log.Fatal("Failed to run command: %s", err)
}
```

## License
GoZoox is released under the [MIT License](./LICENSE).
