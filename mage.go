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
	"assets/",
}

const (
	dictionaryFile = "dictionary.txt"
	aspellLang     = "en_GB"
)

// Make a namespace for the spell checking methods
type Spell mg.Namespace

// Make a type for spell checking logic
type spellCheckingLogic func(string) error

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

func (Spell) Interactive() error {
	return runChecker(func(filename string) error {
		// Run the interactive aspell checker
		return sh.Run(
			"aspell",
			"--home-dir=.",
			fmt.Sprintf("--personal=%s", dictionaryFile),
			"--mode=markdown",
			fmt.Sprintf("--lang=%s", aspellLang),
			"check",
			filename)
	})
}

func (Spell) Check() error {
	return runChecker(func(filename string) error {
		// Read the file contents
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		// Pipe the file into aspell
		cmd := exec.Command(
			"aspell",
			"--home-dir=.",
			fmt.Sprintf("--personal=%s", dictionaryFile),
			"--mode=markdown",
			fmt.Sprintf("--lang=%s", aspellLang),
			"list")
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
		return nil
	})
}

func runChecker(handler spellCheckingLogic) error {
	// We require the build output for linting html
	mg.Deps(Build)
	// Find all of the html files
	matches, err := filepath.Glob("./content/**/*.md")
	if err != nil {
		return err
	}
	// Run the spellchecker on each
	for _, match := range matches {
		// Run the logic on the file
		err := handler(match)
		if err != nil {
			return err
		}
	}
	return nil
}
