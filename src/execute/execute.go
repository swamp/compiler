/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package execute

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/swamp/compiler/src/generate_sp"
	"github.com/swamp/compiler/src/parser"

	"github.com/fatih/color"
	swampcompiler "github.com/swamp/compiler/src/compiler"
	"github.com/swamp/compiler/src/environment"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/verbosity"
)

func Execute(verboseFlag verbosity.Verbosity, executableName string, arguments ...string) (string, error) {
	command := exec.Command(executableName, arguments...)
	color.Magenta("executing swamp %s %s", executableName, strings.Join(arguments, " "))
	output, err := command.Output()
	if err != nil {
		color.Red("execute failed: '%v'", err)
		return "", err
	}
	if verboseFlag > verbosity.None {
		log.Printf("execute returned %s", output)
	}
	if !command.ProcessState.Success() {
		color.Red("couldn't execute")
		return "", fmt.Errorf("couldn't execute")
	}

	return string(output), err
}

func FindProjectDirectory() (string, error) {
	wd, _ := os.Getwd()
	for i := 0; !strings.HasSuffix(wd, "swamp-compiler"); i++ {
		wd = path.Dir(wd)

		if i > 4 {
			return "", fmt.Errorf("could not find project directory")
		}
	}

	return wd, nil
}

func ExecuteSwamp(swampCode string) (string, error) {
	const tempOutputFileTemplate = "temp.swamp-pack"

	tempDir, err := os.MkdirTemp("", "swamptest")
	if err != nil {
		return "", err
	}

	tempFileName := path.Join(tempDir, "Main.swamp")

	tmpFile, tmpFileErr := os.Create(tempFileName)
	if tmpFileErr != nil {
		return "", tmpFileErr
	}

	tempOutputFile := path.Join(tempDir, tempOutputFileTemplate)

	tempSwampFilename := tmpFile.Name()
	const verbose = verbosity.None
	if verbose > verbosity.None {
		log.Printf("=========== TEMPFILE:%v =============\n", tempSwampFilename)
	}

	if _, err := tmpFile.WriteString(swampCode); err != nil {
		return "", err
	}

	tmpFile.Close()
	const enforceStyle = true
	const showAssembly = false
	resourceNameLookup := resourceid.NewResourceNameLookupImpl()

	gen := generate_sp.NewGenerator()
	gen.PrepareForNewPackage()
	_, compileErr := swampcompiler.CompileAndLink(gen, resourceNameLookup, environment.Environment{}, "temp", tempSwampFilename, tempOutputFile, enforceStyle, verbose, showAssembly)
	if parser.IsCompileError(compileErr) {
		return "", compileErr
	}

	projectPath, projectPathErr := FindProjectDirectory()
	if projectPathErr != nil {
		return "", projectPathErr
	}
	completePath := path.Join(projectPath, "bin/swamp_run_linux_amd64")
	return Execute(verbose, completePath, "-v", tempOutputFile)
}
