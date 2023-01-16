/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_c

import (
	"fmt"
	"io"

	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/verbosity"
)

func expressionHasOwnReturn(e decorated.Expression) bool {
	switch e.(type) {
	case *decorated.Let:
		return true
	case *decorated.If:
		return true
	}
	return false
}

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName,
	f *decorated.FunctionValue, writer io.Writer, returnString string, indentation int,
	verboseFlag verbosity.Verbosity) error {
	functionType := f.Type().(*dectype.FunctionTypeReference).FunctionAtom()

	returnTypeName := functionType.ReturnType().HumanReadable()

	parameters := ""
	for index, parameter := range f.Parameters() {
		if index > 0 {
			parameters += ", "
		}
		parameterTypeName := parameter.Type().HumanReadable()
		parameterName := parameter.Parameter().Name()
		parameters += fmt.Sprintf("%v %v", parameterTypeName, parameterName)
	}

	fmt.Fprintf(writer, fmt.Sprintf("%s%v %v(%v)\n%s{\n%s", indentationString(indentation), returnTypeName, fullyQualifiedVariableName.ResolveToString(), parameters, indentationString(indentation), indentationString(indentation+1)))

	shouldInsertReturn := !expressionHasOwnReturn(f.Expression())
	if shouldInsertReturn {
		fmt.Fprintf(writer, returnString)
	}

	if genErr := generateExpression(f.Expression(), writer, "return ", indentation+1); genErr != nil {
		return genErr
	}
	if shouldInsertReturn {
		fmt.Fprintf(writer, ";")
	}

	fmt.Fprintf(writer, "\n%s}\n%s", indentationString(indentation), indentationString(indentation))

	return nil
}
