/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	decshared "github.com/swamp/compiler/src/decorated/decshared"
)

func testDecorateInternal(code string, useCores bool, errorsAsWarnings bool) (string, decshared.DecoratedError) {
	code = strings.TrimSpace(code)
	module, compileErr := CompileToModuleOnceForTest(code, useCores, errorsAsWarnings)
	if compileErr != nil {
		return "", compileErr
	}
	return module.ShortString(), nil
}

func testDecorate(t *testing.T, code string, ast string) {
	decorationString, decorateErr := testDecorateInternal(code, true, false)
	if decorateErr != nil {
		fmt.Printf("problem %v\n", decorateErr)

		t.Fatal(decorateErr)
	}
	ast = strings.TrimSpace(ast)
	decorationString = strings.TrimSpace(decorationString)
	if ast != decorationString {
		astLines := strings.Split(ast, "\n")
		decorationStringLines := strings.Split(decorationString, "\n")
		for index, astLine := range astLines {
			decoratedLine := ""
			if index < len(decorationStringLines) {
				decoratedLine = decorationStringLines[index]
			}

			trimmedDecoratedLine := strings.TrimSpace(decoratedLine)
			trimmedAstLine := strings.TrimSpace(astLine)
			if trimmedAstLine != trimmedDecoratedLine {
				fmt.Printf("detected line diff: \ncorrect: \n%v\n", trimmedAstLine)
				fmt.Printf("wrong: \n%v\n", trimmedDecoratedLine)
				break
			}
		}
		t.Errorf("Mismatch strings received \n%v\n but expected\n%v", decorationString, ast)
	}
}

func testDecorateWithoutDefault(t *testing.T, code string, ast string) {
	decorationString, decorateErr := testDecorateInternal(code, false, false)
	if decorateErr != nil {
		fmt.Printf("problem %v\n", decorateErr)

		t.Fatal(decorateErr)
	}
	ast = strings.TrimSpace(ast)
	decorationString = strings.TrimSpace(decorationString)
	if ast != decorationString {
		astLines := strings.Split(ast, "\n")
		decorationStringLines := strings.Split(decorationString, "\n")
		for index, astLine := range astLines {
			decoratedLine := ""
			if index < len(decorationStringLines) {
				decoratedLine = decorationStringLines[index]
			}

			trimmedDecoratedLine := strings.TrimSpace(decoratedLine)
			trimmedAstLine := strings.TrimSpace(astLine)
			if trimmedAstLine != trimmedDecoratedLine {
				fmt.Printf("correct: \n%v\n", trimmedAstLine)
				fmt.Printf("wrong: \n%v\n", trimmedDecoratedLine)
				break
			}
		}
		t.Errorf("Mismatch strings received \n%v\n but expected\n%v", decorationString, ast)
	}
}

func testDecorateFail(t *testing.T, code string, expectedError interface{}) {
	const errorsAsWarnings = true
	_, testErr := testDecorateInternal(code, true, errorsAsWarnings)
	if testErr == nil {
		t.Errorf("it was supposed to fail, but didn't")
		return
	}
	isSameErr := reflect.TypeOf(expectedError) == reflect.TypeOf(testErr)
	if !isSameErr {
		t.Errorf("unexpected fail: %v %T but expected %T", testErr, testErr, expectedError)
	}
}
