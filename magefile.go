// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"time"
)

const (
	entryPoint = "cmd/lbs.go"
	ldFlags    = "-X $PACKAGE/version/version.commitHash=$COMMIT_HASH " +
		"-X $PACKAGE/version/version.buildDate=$BUILD_DATE"
)

// allow user to override go executable by running as GOEXE=xxx mage ... on
// UNIX-like systems.
var goexe = "go"

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

// Build binary
func Build() error {
	return runCmd(flagEnv(), goexe, "build", "-ldflags", ldFlags, entryPoint)
}

// Build binary with race detector enabled
func BuildRace() error {
	return runCmd(flagEnv(), goexe, "build", "-race", "-ldflags", ldFlags,
		entryPoint)
}

// Run tests
func Test() error {
	return runCmd(flagEnv(), goexe, "test", "./...")
}

// Run tests with race detector
func TestRace() error {
	return runCmd(flagEnv(), goexe, "test", "-race", "./...")
}

// Run basic golangci-lint check.
func Lint() error {
	linterArgs := []string{
		"run",
		"--disable-all",
		"--enable=govet",
		"--enable=gofmt",
		"--enable=gosec",
	}

	if err := runCmd(flagEnv(), "golangci-lint", linterArgs...); err != nil {
		return err
	}

	return nil
}

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"PACKAGE":     entryPoint,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

func runCmd(env map[string]string, cmd string, args ...string) error {
	if mg.Verbose() {
		return sh.RunWith(env, cmd, args...)
	}

	if output, err := sh.OutputWith(env, cmd, args...); err != nil {
		fmt.Fprint(os.Stderr, output)
	}

	return nil
}