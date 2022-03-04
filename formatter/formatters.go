// Package is a utility wrapper around command line tools
// to format Go, Dart and TypeScript code.
package formatter

import (
	"log"
	"os/exec"
)

// Formatters provides format commands for Go, Dart and TypeScript.
type Formatters struct {
	hasGoFmt, hasDartFmt, hasTsFmt *bool
}

type Format uint8

const (
	NoFormat Format = iota
	Go
	Dart
	Ts
)

// check if the goimports command is working
// and caches the result
func (fmts *Formatters) hasGo() bool {
	if fmts.hasGoFmt == nil {
		err := exec.Command("which", "goimports").Run()
		if err != nil {
			log.Printf("No formatter for Go (%s)", err)
		} else {
			log.Println("Formatter for Go detected")
		}
		fmts.hasGoFmt = new(bool)
		*fmts.hasGoFmt = err == nil
	}
	return *fmts.hasGoFmt
}

// check if the dart command is working
// and caches the result
func (fmts *Formatters) hasDart() bool {
	if fmts.hasDartFmt == nil {
		err := exec.Command("dart", "format", "--help").Run()
		if err != nil {
			log.Printf("No formatter for Dart (%s)", err)
		} else {
			log.Println("Formatter for Dart detected")
		}
		fmts.hasDartFmt = new(bool)
		*fmts.hasDartFmt = err == nil
	}
	return *fmts.hasDartFmt
}

// check if the prettier command is working
// and caches the result
func (fmts *Formatters) hasTypescript() bool {
	if fmts.hasTsFmt == nil {
		err := exec.Command("npx", "prettier", "-v").Run()
		if err != nil {
			log.Printf("No formatter for Typescript (%s)", err)
		} else {
			log.Println("Formatter for Typescript detected")
		}
		fmts.hasTsFmt = new(bool)
		*fmts.hasTsFmt = err == nil
	}
	return *fmts.hasTsFmt
}

// FormatFile format `filename`, if a formatter for `format` is found.
// It returns an error if the command failed, not if no formatter is found.
func (fr *Formatters) FormatFile(format Format, filename string) error {
	switch format {
	case Go:
		if fr.hasGo() {
			return exec.Command("goimports", "-w", filename).Run()
		}
	case Dart:
		if fr.hasDart() {
			return exec.Command("dart", "format", filename).Run()
		}
	case Ts:
		if fr.hasTypescript() {
			return exec.Command("npx", "prettier", "--write", filename).Run()
		}
	}
	return nil
}
