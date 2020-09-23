# structgen [![GoDoc](https://godoc.org/github.com/benoitkugler/structgen?status.svg)](https://godoc.org/github.com/benoitkugler/structgen)

An extremely simple and powerful Go to Typescript definitions.

Inspired by [OneOfOne/structgen](https://github.com/OneOfOne/structgen), but build with [go/types](https://golang.org/pkg/go/types).

## Install

    go get -u -v github.com/benoitkugler/structgen/...

## Command Line Usage

```
➤ structgen -h

Usage of ./structgen:
  -output string
        ts file to write to
  -source string
        go source file to convert

```

All types in `source` are converted (no only structs)
