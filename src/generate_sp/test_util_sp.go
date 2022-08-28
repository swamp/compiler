/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	decorated "github.com/swamp/compiler/src/decorated/expression"

	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/loader"

	"github.com/swamp/compiler/src/parser"

	"github.com/swamp/assembler/lib/assembler_sp"
	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	swampdisasmsp "github.com/swamp/disassembler/lib"
)

func testGenerateInternal(code string) (*assembler_sp.PackageConstants, []*assembler_sp.Constant, error) {
	const useCores = true
	const errorsAsWarnings = false
	module, compileErr := deccy.CompileToModuleOnceForTest(code, useCores, errorsAsWarnings)
	if parser.IsCompileError(compileErr) {
		return nil, nil, compileErr
	}

	fileSystemRoot := loader.LocalFileSystemRoot("")
	pack := loader.NewPackage(fileSystemRoot, "someName")
	fullyQualifiedName := dectype.MakeArtifactFullyQualifiedModuleName(nil)
	pack.AddModule(fullyQualifiedName, module)
	gen := NewGenerator()
	gen.PrepareForNewPackage()
	const verboseFlag = verbosity.None
	_, _, resourceLookup, typeInfoErr := typeinfo.GenerateModule(module)
	if typeInfoErr != nil {
		return nil, nil, typeInfoErr
	}

	genErr := gen.GenerateFromPackage(pack, resourceLookup, verboseFlag)
	if parser.IsCompileErr(genErr) {
		return nil, nil, genErr
	}
	if genErr != nil {
		compileErr = decorated.AppendError(compileErr, decorated.NewInternalError(genErr))
	}
	return gen.PackageConstants(), gen.LastFunctionConstants(), compileErr
}

func checkGeneratedAssembler(constants *assembler_sp.PackageConstants, functions []*assembler_sp.Constant, expectedAsm string) error {
	assemblerOutput := constants.DebugString([]assembler_sp.ConstantType{assembler_sp.ConstantTypeString, assembler_sp.ConstantTypeResourceName})
	if len(assemblerOutput) > 0 {
		assemblerOutput += "\n"
	}
	for _, f := range functions {
		opcodes := constants.FetchOpcodes(f)
		lines := swampdisasmsp.Disassemble(opcodes, false)
		assemblerOutput = assemblerOutput + fmt.Sprintf("func %v\n%s\n\n", f, strings.Join(lines[:], "\n"))
	}

	assemblerOutput = strings.TrimSpace(assemblerOutput)
	expectedAsm = strings.TrimSpace(expectedAsm)
	asmLines := strings.Split(expectedAsm, "\n")
	asmCleanedUp := ""
	for _, asmLine := range asmLines {
		foundIndex := strings.Index(asmLine, "; ")
		if foundIndex != -1 {
			asmLine = strings.TrimSpace(asmLine[:foundIndex])
		}
		asmCleanedUp = asmCleanedUp + asmLine + "\n"
	}
	asmCleanedUp = strings.TrimSpace(asmCleanedUp)
	if asmCleanedUp != assemblerOutput {
		return fmt.Errorf("not matching, generated:\n%v\nExpected:\n%v", assemblerOutput, asmCleanedUp)
	}
	return nil
}

func testGenerateInternalWithAssemblerCheck(code string, expectedAsm string) error {
	constants, functions, generateErr := testGenerateInternal(code)
	if generateErr != nil {
		if parser.TypeOfWarningRecursive(generateErr) >= parser.ReportAsSeverityError {
			return generateErr
		}
	}
	checkErr := checkGeneratedAssembler(constants, functions, expectedAsm)
	return checkErr
}

func testGenerate(t *testing.T, code string, expectedAsm string) {
	code = strings.TrimSpace(code)
	decorateErr := testGenerateInternalWithAssemblerCheck(code, expectedAsm)
	if decorateErr != nil {
		log.Printf("problem %v\n", decorateErr)
		t.Error(decorateErr)
	}
}

func testGenerateFail(t *testing.T, code string, expectedError interface{}) {
	code = strings.TrimSpace(code)
	_, _, testErr := testGenerateInternal(code)
	if testErr == nil {
		log.Printf("problem, should fail")
		t.Errorf("was supposed to fail")
		return
	}
	isSameErr := reflect.TypeOf(expectedError) == reflect.TypeOf(testErr)
	if !isSameErr {
		t.Errorf("generate: unexpected fail: expected %T but received %T", expectedError, testErr)
	}
}
