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
)

func stringComparisonToCCode(operatorType decorated.BooleanOperatorType) string {
	switch operatorType {
	case decorated.BooleanEqual:
		return "=="
	case decorated.BooleanNotEqual:
		return "!="
	case decorated.BooleanLess:
		return "<"
	case decorated.BooleanLessOrEqual:
		return "<="
	case decorated.BooleanGreater:
		return ">"
	case decorated.BooleanGreaterOrEqual:
		return ">="
	default:
		panic("not supported")
	}
}

func booleanToCCode(operatorType decorated.BooleanOperatorType) string {
	switch operatorType {
	case decorated.BooleanEqual:
		return "=="
	case decorated.BooleanNotEqual:
		return "!="
	case decorated.BooleanLess:
		return "<"
	case decorated.BooleanLessOrEqual:
		return "<="
	case decorated.BooleanGreater:
		return ">"
	case decorated.BooleanGreaterOrEqual:
		return ">="
	default:
		panic("not supported")
	}
}

func generateBinaryOperatorBooleanResult(operator *decorated.BooleanOperator, writer io.Writer, indentation int) error {
	fmt.Fprintf(writer, "(")
	unaliasedTypeLeft := dectype.UnaliasWithResolveInvoker(operator.Left().Type())
	foundPrimitive, _ := unaliasedTypeLeft.(*dectype.PrimitiveAtom)
	isStringComparison := foundPrimitive != nil && foundPrimitive.AtomName() == "String"

	if isStringComparison {
		fmt.Fprintf(writer, "strcmp(")
		leftErr := generateExpression(operator.Left(), writer, "", indentation)
		if leftErr != nil {
			return leftErr
		}
		fmt.Fprintf(writer, ", ")
		rightErr := generateExpression(operator.Right(), writer, "", indentation)
		if rightErr != nil {
			return rightErr
		}
		fmt.Fprintf(writer, ")")
		comparisonString := booleanToCCode(operator.OperatorType())
		fmt.Fprintf(writer, "%s 0)", comparisonString)
		return nil
	}

	leftErr := generateExpression(operator.Left(), writer, "", indentation)
	if leftErr != nil {
		return leftErr
	}

	if foundPrimitive == nil {
		foundCustomType, _ := unaliasedTypeLeft.(*dectype.CustomTypeAtom)
		if foundCustomType == nil {
			panic(fmt.Errorf("not implemented binary operator boolean %v", unaliasedTypeLeft.HumanReadable()))
		} else {
			// unaliasedTypeRight := dectype.UnaliasWithResolveInvoker(operator.Right().Type())
			panic("todo")
			//			panic(fmt.Errorf("not implemented yet %v", unaliasedTypeRight))
		}
	} else if foundPrimitive.AtomName() == "Int" || foundPrimitive.AtomName() == "Char" {
		comparisonString := booleanToCCode(operator.OperatorType())
		fmt.Fprintf(writer, " "+comparisonString+" ")
	} else {
		panic(fmt.Errorf("what operator is this for %v", foundPrimitive.AtomName()))
	}

	rightErr := generateExpression(operator.Right(), writer, "", indentation)
	if rightErr != nil {
		return rightErr
	}

	fmt.Fprintf(writer, ")")

	return nil
}
