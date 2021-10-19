//go:build mage
// +build mage

package main

import (
	"github.com/scottames/cmder"
)

var (
	g0 = cmder.New("go").RunFn()
)

// Check run checks
func Check() {
	Lint()
	Test()
}

// Lint run golangci-lint
func Lint() error {
	return cmder.New("golangci-lint", "run").Run()
}

// Test run go tests
func Test() error {
	return g0("test", "./...")
}

// Tidy run go mod tidy
func Tidy() error {
	return g0("mod", "tidy")
}
