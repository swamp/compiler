/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/parser"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/piot/lsp-server/lspserv"
	"github.com/swamp/compiler/src/doc"
	"github.com/swamp/compiler/src/environment"
	"github.com/swamp/compiler/src/verbosity"

	swampcompiler "github.com/swamp/compiler/src/compiler"
	"github.com/swamp/compiler/src/file"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/lspservice"
)

var Version string

func buildCommandLine(fileOrDirectory string, outputDirectory string, enforceStyle bool, assembler bool, target swampcompiler.Target, verbosity verbosity.Verbosity) ([]*loader.Package, error) {
	filenameToCompile := fileOrDirectory

	return swampcompiler.BuildMain(filenameToCompile, outputDirectory, enforceStyle, assembler, target, verbosity)
}

func buildCommandLineNoOutput(fileOrDirectory string, enforceStyle bool, verbosity verbosity.Verbosity) ([]*loader.Package, error) {
	filenameToCompile := fileOrDirectory

	return swampcompiler.BuildMainOnlyCompile(filenameToCompile, enforceStyle, verbosity)
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
		octets, err := os.ReadFile(match)
		if err != nil {
			return err
		}

		beautifiedCode, beautyErr := beautify("", string(octets))
		if beautyErr != nil {
			return beautyErr
		}

		os.WriteFile(match, []byte(beautifiedCode), 0o600)

		fmt.Println(beautifiedCode)
	}
	return nil
}

type DocCmd struct {
	Path         string `help:"path to file or directory" arg:"" default:"." type:"path"`
	Verbosity    int    `help:"verbose output" type:"counter" short:"v"`
	DisableStyle bool   `help:"disable enforcing of style" default:"false"`
}

func (c *DocCmd) Run() error {
	compiledPackages, err := buildCommandLineNoOutput(c.Path, !c.DisableStyle, verbosity.Verbosity(c.Verbosity))
	if err != nil {
		return err
	}

	if err := doc.PackagesToHtmlPage(os.Stdout, compiledPackages); err != nil {
		return err
	}

	if c.Verbosity > 0 {
		color.Green(fmt.Sprintf("done. %v", len(compiledPackages)))
	}

	return nil
}

type LspCmd struct{}

func (c *LspCmd) Run() error {
	fileSystem := loader.NewFileSystemDocumentProvider()
	config, _, configErr := environment.LoadFromConfig()
	if configErr != nil {
		return configErr
	}
	lspService := lspservice.NewLspImpl(fileSystem, config)
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
	Target       string `help:"target platform" enum:"swamp-pack,llvm-ir" short:"t" default:"swamp-pack"`
	Verbosity    int    `help:"verbose output" type:"counter" short:"v"`
	Assembler    bool   `help:"output assembler" short:"s" default:"false"`
	Modules      string
}

func (c *BuildCmd) Run() error {
	if c.Path == "" {
		return fmt.Errorf("must specify build directory")
	}

	c.Path = filepath.ToSlash(c.Path)

	target := swampcompiler.SwampOpcode
	if c.Target == "llvm-ir" {
		target = swampcompiler.LlvmIr
	}

	compiledPackages, err := buildCommandLine(c.Path, c.Output, !c.DisableStyle, c.Assembler, target, verbosity.Verbosity(c.Verbosity))
	if err != nil {
		return err
	}

	if c.Verbosity > 0 {
		color.Green(fmt.Sprintf("done. %v", len(compiledPackages)))
	}

	return nil
}

type EnvironmentSetCmd struct {
	Name string `help:"fmt" arg:""`
	Path string `help:"fmt" arg:""`
}

func (c *EnvironmentSetCmd) Run() error {
	fmt.Printf("setting '%v'='%v'\n", c.Name, c.Path)

	configuration, _, err := environment.LoadFromConfig()
	if err != nil {
		return err
	}

	absolutePath, absErr := filepath.Abs(c.Path)
	if absErr != nil {
		return absErr
	}

	configuration.AddOrSet(c.Name, absolutePath)

	return configuration.SaveToConfig()
}

type EnvironmentListCmd struct{}

func (c *EnvironmentListCmd) Run() error {
	configuration, _, err := environment.LoadFromConfig()
	if err != nil {
		return err
	}

	fmt.Printf("Environment:\n")
	for _, module := range configuration.Package {
		fmt.Printf("'%v'='%v'\n", module.Name, module.Path)
	}

	return nil
}

type EnvironmentCmd struct {
	Set  EnvironmentSetCmd  `cmd:"" help:"set swamp environment package"`
	List EnvironmentListCmd `cmd:"" default:"1" help:"list swamp environment packages"`
}

func (c *EnvironmentCmd) Run() error {
	return nil
}

type VersionCmd struct{}

func (c *VersionCmd) Run() error {
	fmt.Printf("swamp v%v\n", Version)
	return nil
}

type Options struct {
	Lsp     LspCmd         `help:"lsp" cmd:""`
	Fmt     FmtCmd         `help:"fmt" cmd:""`
	Doc     DocCmd         `help:"fmt" cmd:""`
	Build   BuildCmd       `cmd:"" help:"builds a swamp application"`
	Env     EnvironmentCmd `cmd:"" help:"manage swamp environment"`
	Version VersionCmd     `cmd:"" help:"shows the version information"`
}

/*
highestError := parser.ReportAsSeverityNote
multiErr, wasMultiErr := decErr.(*decorated.MultiErrors)

	if wasMultiErr {
		for _, subErr := range multiErr.Errors() {

			if detectedError > highestError {
				highestError = detectedError
			}
		}
	} else {

		multiErrPar, wasMultiErr := decErr.(parerr.MultiError)
		if wasMultiErr {
			for _, subErr := range multiErrPar.Errors() {
				detectedError := parser.ShowWarningOrError(nil, subErr)
				if detectedError > highestError {
					highestError = detectedError
				}
			}
		} else if decErr != nil {
			detectedError := parser.ShowWarningOrError(nil, decErr)
			if detectedError > highestError {
				highestError = detectedError
			}
		} else {
			fmt.Fprintf(os.Stderr, "Unknown ERROR!!: '%v'\n", err)
		}
	}
*/
func main() {
	ctx := kong.Parse(&Options{})

	err := ctx.Run()
	if err != nil {
		log.Print(err)
		decErr, wasDecorated := err.(decshared.DecoratedError)
		highestError := parser.ReportAsSeverityError
		if wasDecorated {
			highestError = parser.ShowWarningOrError(nil, decErr)
		}
		if highestError >= parser.ReportAsSeverityError {
			os.Exit(-1)
		}
	}

	os.Exit(0)
}
