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
	"github.com/piot/lsp-server/lspserv"

	swampcompiler "github.com/swamp/compiler/src/compiler"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/lspservice"
)

var Version string

func buildCommandLine(fileOrDirectory string, outputDirectory string, enforceStyle bool, verboseFlag bool) error {
	filenameToCompile := fileOrDirectory

	return swampcompiler.BuildMain(filenameToCompile, outputDirectory, enforceStyle, verboseFlag)
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

type LspCmd struct{}

func (c *LspCmd) Run() error {
	fileSystem := loader.NewFileSystemDocumentProvider()
	lspService := lspservice.NewLspImpl(fileSystem)
	service := lspservice.NewService(lspService, lspService, lspService, lspService)
	fmt.Fprintf(os.Stderr, "LSP Server initiated. Will receive commands from stdin and send reply on stdout")
	lspServ := lspserv.NewService(service)
	const logOutput = false
	lspServ.RunUntilClose(lspserv.StdInOutReadWriteCloser{}, logOutput)

	return nil
}

type BuildCmd struct {
	Path         string `help:"path to file or directory" arg:"" default:"." type:"path"`
	DisableStyle bool   `help:"disable enforcing of style" default:"false"`
	Output       string `help:"output directory" type:"existingdir" short:"o" default:"."`
	IsVerbose    bool   `help:"verbose output"`
	Modules      string
}

func (c *BuildCmd) Run() error {
	if c.Path == "" {
		return fmt.Errorf("must specify build directory")
	}

	err := buildCommandLine(c.Path, c.Output, !c.DisableStyle, c.IsVerbose)
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
		decErr, wasDecorated := err.(decshared.DecoratedError)
		if wasDecorated {
			moduleErr, wasModuleErr := decErr.(*decorated.ModuleError)
			if wasModuleErr {
				decErr = moduleErr.WrappedError()
			}
		}
		multiErr, _ := decErr.(*decorated.MultiErrors)
		if multiErr != nil {
			for _, subErr := range multiErr.Errors() {
				fmt.Printf("%v ERROR:%v\n", subErr.FetchPositionLength(), subErr)
			}
		} else if decErr != nil {
			fmt.Printf("%v ERROR:%v\n", decErr.FetchPositionLength(), err)
		} else {
			fmt.Printf("Unknown ERROR: '%v'\n", err)
		}
		os.Exit(-1)
	}

	os.Exit(0)
}
