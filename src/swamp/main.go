/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"

	swampcompiler "github.com/swamp/compiler/src/compiler"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/lspservice"
	"github.com/swamp/compiler/src/token"
)

var Version string

func compileAndLink(fileOrDirectory string, outputFilename string, enforceStyle bool, verboseFlag bool) error {
	filenameToCompile := fileOrDirectory

	statInfo, statErr := os.Stat(fileOrDirectory)
	if statErr != nil {
		return statErr
	}

	if statInfo.IsDir() {
		swampDirectory := fileOrDirectory
		filenameToCompile = swampDirectory
	}

	world, _, compileErr := swampcompiler.CompileMain(filenameToCompile, enforceStyle, verboseFlag)
	if compileErr != nil {
		return compileErr
	}
	return swampcompiler.GenerateAndLink(world, outputFilename, verboseFlag)
}

type FmtCmd struct {
	Path string `help:"fmt" arg:""`
}

func (c *FmtCmd) Run() error {
	matches, globErr := filepath.Glob(c.Path)
	if globErr != nil {
		return globErr
	}

	for _, match := range matches {
		if file.IsDir(match) {
			continue
		}
		octets, err := ioutil.ReadFile(match)
		if err != nil {
			return err
		}

		beautifiedCode, beautyErr := beautify("", string(octets))
		if beautyErr != nil {
			return beautyErr
		}

		ioutil.WriteFile(match, []byte(beautifiedCode), 0o600)

		fmt.Println(beautifiedCode)
	}
	return nil
}

type LspImpl struct {
	world  *loader.World
	module *decorated.Module
}

func (l *LspImpl) Compile(filename string) error {
	const enforceStyle = true
	const verboseFlag = false
	world, _, err := swampcompiler.CompileMain(filename, enforceStyle, verboseFlag)
	if err != nil {
		return err
	}

	l.world = world
	return nil
}

func (l *LspImpl) FindToken(position token.Position) lspservice.DecoratedTypeOrToken {
	return nil
}

type LspCmd struct{}

func (c *LspCmd) Run() error {
	fmt.Fprintf(os.Stderr, "LSP Server initiated. Will receive commands from stdin and send reply on stdout")
	return nil
}

type BuildCmd struct {
	Path         string `help:"path to file or directory" arg:"" default:"." type:"path"`
	DisableStyle bool   `help:"disable enforcing of style" default:"false"`
	Output       string `help:"output file name" short:"o" default:"out.swamp-pack"`
	IsVerbose    bool   `help:"verbose output"`
	Modules      string
}

func (c *BuildCmd) Run() error {
	if c.Path == "" {
		return fmt.Errorf("must specify build directory")
	}

	err := compileAndLink(c.Path, c.Output, !c.DisableStyle, c.IsVerbose)
	if err != nil {
		return err
	}

	if c.IsVerbose {
		color.Green("done.")
	}

	return nil
}

type VersionCmd struct{}

func (c *VersionCmd) Run() error {
	fmt.Printf("swamp v%v\n", Version)
	return nil
}

type Options struct {
	Lsp     LspCmd     `help:"lsp" cmd:""`
	Fmt     FmtCmd     `help:"fmt" cmd:""`
	Build   BuildCmd   `cmd:"" default:"1" help:"builds a swamp application"`
	Version VersionCmd `cmd:"" help:"shows the version information"`
}

func main() {
	ctx := kong.Parse(&Options{})

	err := ctx.Run()
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(-1)
	}

	os.Exit(0)
}
