/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package execute

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	swampcompiler "github.com/swamp/compiler/src/compiler"
)

func Execute(verboseFlag bool, executableName string, arguments ...string) (string, error) {
	command := exec.Command(executableName, arguments...)
	color.Magenta("executing swamp %s %s", executableName, strings.Join(arguments, " "))
	output, err := command.Output()
	if err != nil {
		color.Red("execute failed: '%v'", err)
		return "", err
	}
	if verboseFlag {
		fmt.Printf("execute returned %s", output)
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
		wd = filepath.Dir(wd)
		if i > 4 {
			return "", fmt.Errorf("could not find project directory")
		}
	}
	return wd, nil
}

func ExecuteSwamp(swampCode string) (string, error) {
	const tempOutputFileTemplate = "temp.swamp-pack"
	const tempSwampFileSuffix = "temp.swamp"

	tempDir, err := ioutil.TempDir("", "test.swamp")
	if err != nil {
		return "", err
	}
	tempFileName := filepath.Join(tempDir, "Main.swamp")
	tmpFile, tmpFileErr := os.Create(tempFileName)
	if tmpFileErr != nil {
		return "", tmpFileErr
	}

	tempOutputFile := filepath.Join(tempDir, tempOutputFileTemplate)

	tempSwampFilename := tmpFile.Name()
	const verbose = false
	if verbose {
		fmt.Printf("=========== TEMPFILE:%v =============\n", tempSwampFilename)
	}
	tmpFile.WriteString(swampCode)
	tmpFile.Close()
	const enforceStyle = true
	compileErr := swampcompiler.CompileAndLink(tempSwampFilename, tempOutputFile, enforceStyle, verbose)
	if compileErr != nil {
		return "", compileErr
	}

	projectPath, projectPathErr := FindProjectDirectory()
	if projectPathErr != nil {
		return "", projectPathErr
	}
	completePath := path.Join(projectPath, "bin/swamp_run_linux_amd64")
	return Execute(verbose, completePath, "-v", tempOutputFile)
}
