/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_ir

import (
	"fmt"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateArithmeticMultiple(operator *decorated.ArithmeticOperator, genContext *generateContext) (value.Value, error) {
	leftPrimitive, _ := dectype.UnReference(operator.Left().Type()).(*dectype.PrimitiveAtom)
	switch {
	case dectype.IsListLike(operator.Left().Type()) && operator.OperatorType() == decorated.ArithmeticAppend:
		//return generateListAppend(code, target, operator, genContext)
	case leftPrimitive != nil && leftPrimitive.AtomName() == "String" && operator.OperatorType() == decorated.ArithmeticAppend:
		//return generateStringAppend(code, target, operator, genContext)
	case dectype.IsIntLike(operator.Left().Type()):
		return generateArithmeticInt(operator, genContext)
	default:
		return nil, fmt.Errorf("cant generate arithmetic for type: %v <-> %v (%v)",
			operator.Left().Type(), operator.Right().Type(), operator.OperatorType())
	}

	return nil, fmt.Errorf("internal error")
}
