//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

var sourceDirs = []string{
	"archetypes/",
	"content/",
	"resources/",
	"static/",
	"themes/",
}

// Lint the markdown source
func Lint() error {
	return sh.Run("markdownlint", "content", "archetypes")
}

// Build the project
func Build() error {
	isChanged, err := target.Dir("public/", sourceDirs...)
	if err != nil {
		return err
	}
	if !isChanged {
		return nil
	}
	return sh.Run("hugo", "-v")
}

// Start the debug server for blog development
func Serve() error {
	// Serve drafts
	return sh.Run("hugo", "server", "-D")
}

func Proof() error {
	// We require the build output for linting html
	mg.Deps(Build)
	// Proofread the generated HTML
	return sh.Run("htmlproofer", "./public", "--allow-hash-href", "--check-html")
}

func SpellCheck() error {
	// We require the build output for linting html
	mg.Deps(Build)
	// Proofread the generated HTML
	matches, err := filepath.Glob("./public/**/*.html")
	if err != nil {
		return err
	}
	for _, match := range matches {
		// Read the file contents
		contents, err := ioutil.ReadFile(match)
		if err != nil {
			return err
		}
		// Pipe the file into aspell
		cmd := exec.Command("aspell", "--mode=html", "list")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		// Pipe the file into stdin
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, string(contents))
		}()
		// Run the command
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		fmt.Printf(string(out))
	}
	return nil
}
