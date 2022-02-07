package generate_c

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"testing"

	deccy "github.com/swamp/compiler/src/decorated"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
)

func testGenerateInternal(code string, writer io.Writer) error {
	const useCores = true
	const errorsAsWarnings = false
	module, compileErr := deccy.CompileToModuleOnceForTest(code, useCores, errorsAsWarnings)
	if compileErr != nil {
		return compileErr
	}

	gen := NewGenerator()

	const verboseFlag = verbosity.None
	_, _, _, typeInfoErr := typeinfo.GenerateModule(module)
	if typeInfoErr != nil {
		return typeInfoErr
	}
	genErr := gen.GenerateAllLocalDefinedFunctions(module, writer, verboseFlag)
	if genErr != nil {
		return genErr
	}
	return genErr
}

func checkGeneratedC(generatedC string, expectedC string) error {
	generatedC = strings.TrimSpace(generatedC)
	expectedC = strings.TrimSpace(expectedC)
	asmLines := strings.Split(expectedC, "\n")
	expectedCleanedUp := ""
	for _, asmLine := range asmLines {
		foundIndex := strings.Index(asmLine, "; ")
		if foundIndex != -1 {
			asmLine = strings.TrimSpace(asmLine[:foundIndex])
		}
		expectedCleanedUp = expectedCleanedUp + asmLine + "\n"
	}
	expectedCleanedUp = strings.TrimSpace(expectedCleanedUp)
	if expectedCleanedUp != generatedC {
		return fmt.Errorf("not matching, generated:\n%v\nExpected:\n%v", generatedC, expectedCleanedUp)
	}
	return nil
}

func testGenerateInternalWithAssemblerCheck(code string, expectedC string) error {
	buf := &strings.Builder{}
	generateErr := testGenerateInternal(code, buf)
	if generateErr != nil {
		return generateErr
	}
	checkErr := checkGeneratedC(buf.String(), expectedC)
	return checkErr
}

func testGenerate(t *testing.T, code string, expectedCCode string) {
	code = strings.TrimSpace(code)
	decorateErr := testGenerateInternalWithAssemblerCheck(code, expectedCCode)
	if decorateErr != nil {
		log.Printf("problem %v\n", decorateErr)
		t.Error(decorateErr)
	}
}

func testGenerateFail(t *testing.T, code string, expectedError interface{}) {
	code = strings.TrimSpace(code)
	buf := &strings.Builder{}
	testErr := testGenerateInternal(code, buf)
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
