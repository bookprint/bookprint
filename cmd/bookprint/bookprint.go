/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"

	"stefanco.de/bookprint/internal/bookprint"
	"stefanco.de/bookprint/internal/util/fs"
)

const usage = `
Usage:
	bookprint [options...] <file>

Options:
	-t, --template-dir <dir>    Path to the directory containing custom templates used for generating the book.
	-o, --output-dir <dir>      Path to the directory where the generated book pages will be stored.
	-s, --static-dir <dir>      Path to the directory with additional files for the book. Copied to output directory.
	-v, --version               Print the version number.
	-h, --help                  Print the help message.

Examples:
	Reading from HTML file:
	$ bookprint --template-dir templates --output-dir out examples/index.html
	> Created book in 'out' directory

	Reading from STDIN:
	$ echo "<html>...</html>" | bookprint --template-dir templates --output-dir out --
	> Created book in 'out' directory
`

// Version can be set at link time to override debug.BuildInfo.Main.Version,
// which is "(devel)" when building from within the module.
//
// See: https://golang.org/issue/29814
// See: https://golang.org/issue/29228
// See: Dockerfile
var Version string

func main() {
	flag.Usage = getUsage

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	var (
		templateDirectoryFlag string
		outputDirectoryFlag   string
		staticDirectoryFlag   string
		versionFlag           bool
		helpFlag              bool
	)

	flag.StringVar(&templateDirectoryFlag, "t", "templates", "Path to the directory containing custom templates used for generating the book.")
	flag.StringVar(&templateDirectoryFlag, "template-dir", "templates", "Path to the directory containing custom templates used for generating the book.")
	flag.StringVar(&staticDirectoryFlag, "s", "", "Path to the directory with additional files for the book. Copied to output directory.")
	flag.StringVar(&staticDirectoryFlag, "static-dir", "", "Path to the directory with additional files for the book. Copied to output directory.")
	flag.StringVar(&outputDirectoryFlag, "o", "out", "Path to the directory where the generated book pages will be stored.")
	flag.StringVar(&outputDirectoryFlag, "output-dir", "out", "Path to the directory where the generated book pages will be stored.")
	flag.BoolVar(&versionFlag, "v", false, "Print the version number.")
	flag.BoolVar(&versionFlag, "version", false, "Print the version number.")
	flag.BoolVar(&helpFlag, "h", false, "Print the help message.")
	flag.BoolVar(&helpFlag, "help", false, "Print the help message.")

	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Println(getVersion())
		os.Exit(0)
	}

	file, err := getFile(flag.Arg(0))
	if err != nil {
		fail(err)
	}

	// Missing template directory
	if !fs.ExistDir(templateDirectoryFlag) {
		fail(fmt.Errorf("directory '%s' does not exist", templateDirectoryFlag))
	}

	// Create output directory
	err = fs.RemoveDir(outputDirectoryFlag)
	if err != nil {
		fail(err)
	}
	err = fs.MakeDir(outputDirectoryFlag)
	if err != nil {
		fail(err)
	}

	// Copy everything from static to output directory
	if staticDirectoryFlag != "" {
		err := fs.CopyDir(staticDirectoryFlag, outputDirectoryFlag)
		if err != nil {
			fail(err)
		}
	}

	err = bookprint.New(&bookprint.Config{
		File:        file,
		OutputDir:   outputDirectoryFlag,
		TemplateDir: templateDirectoryFlag,
	})
	if err != nil {
		fail(err)
	}

	fmt.Printf("Created book in '%s' directory", outputDirectoryFlag)
}

func getUsage() {
	_, err := fmt.Fprintf(os.Stderr, "%s\n\n", strings.TrimSpace(usage))
	if err != nil {
		fail(err)
	}
}

func getVersion() string {
	if Version != "" {
		return fmt.Sprintf("BookPrint (%s) %s/%s", Version, runtime.GOOS, runtime.GOARCH)
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		return fmt.Sprintf("BookPrint %s %s/%s", buildInfo.Main.Version, runtime.GOOS, runtime.GOARCH)
	}

	return fmt.Sprintf("BookPrint (unknown) %s/%s", runtime.GOOS, runtime.GOARCH)
}

func getFile(name string) ([]byte, error) {
	var file []byte

	if name != "" {
		if !fs.ExistFile(name) {
			return file, fmt.Errorf("file '%s' does not exist", name)
		}

		file, err := os.ReadFile(name)
		if err != nil {
			return file, err
		}

		return file, nil
	}

	file, err := fs.StdinAll()
	if err != nil {
		return file, err
	}

	return file, nil
}

func fail(err error) {
	fmt.Printf("Error: %s", err)
	os.Exit(1)
}
