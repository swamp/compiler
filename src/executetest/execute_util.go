/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package executetest

import (
	"strings"
	"testing"

	"github.com/swamp/compiler/src/execute"
)

func internalExecuteTest(code string) (string, error) {
	output, err := execute.ExecuteSwamp(
		code)
	if err != nil {
		return "", err
	}
	return output, nil
}

func executeTest(t *testing.T, code string, expectedResult string) {
	code = strings.TrimSpace(code)
	output, internalErr := internalExecuteTest(code)
	if internalErr != nil {
		t.Fatal(internalErr)
	}

	expectedResult = strings.TrimSpace(expectedResult)
	output = strings.TrimSpace(output)
	if output == "" {
		t.Errorf("sorry, test must produce output to be valid")
	}

	if output != expectedResult {
		t.Errorf("no match. Received\n%v\n but expected \n%v\n", output, expectedResult)
	}
}
