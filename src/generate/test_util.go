/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	deccy "github.com/swamp/compiler/src/decorated"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/typeinfo"

	swampdisasm "github.com/swamp/disassembler/lib"
)

func testGenerateInternal(code string) ([]*Function, error) {
	const useCores = true
	const errorsAsWarnings = false
	module, compileErr := deccy.CompileToModuleOnce(code, useCores, errorsAsWarnings)
	if compileErr != nil {
		return nil, compileErr
	}

	gen := NewGenerator()
	rootContext := decorator.NewVariableContext(module.LocalAndImportedDefinitions())
	const verboseFlag = false
	_, lookup, typeInfoErr := typeinfo.GenerateModule(module)
	if typeInfoErr != nil {
		return nil, typeInfoErr
	}
	functions, genErr := gen.GenerateAllLocalDefinedFunctions(module, rootContext, lookup, verboseFlag)
	if genErr != nil {
		return nil, genErr
	}
	return functions, genErr
}

func checkGeneratedAssembler(functions []*Function, expectedAsm string) error {
	var assemblerOutput string
	for _, f := range functions {
		lines := swampdisasm.Disassemble(f.opcodes)
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
	functions, generateErr := testGenerateInternal(code)
	if generateErr != nil {
		return generateErr
	}
	checkErr := checkGeneratedAssembler(functions, expectedAsm)
	return checkErr
}

func testGenerate(t *testing.T, code string, expectedAsm string) {
	code = strings.TrimSpace(code)
	decorateErr := testGenerateInternalWithAssemblerCheck(code, expectedAsm)
	if decorateErr != nil {
		fmt.Printf("problem %v\n", decorateErr)
		t.Error(decorateErr)
	}
}

func testGenerateFail(t *testing.T, code string, expectedError interface{}) {
	code = strings.TrimSpace(code)
	_, testErr := testGenerateInternal(code)
	if testErr == nil {
		fmt.Printf("problem, should fail")
		t.Errorf("was supposed to fail")
		return
	}
	isSameErr := reflect.TypeOf(expectedError) == reflect.TypeOf(testErr)
	if !isSameErr {
		t.Errorf("generate: unexpected fail: expected %T but received %T", expectedError, testErr)
	}
}
