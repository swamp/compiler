/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_ir

import (
	"fmt"
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func booleanIntToIPred(operatorType decorated.BooleanOperatorType) enum.IPred {
	switch operatorType {
	case decorated.BooleanEqual:
		return enum.IPredEQ
	case decorated.BooleanNotEqual:
		return enum.IPredNE
	case decorated.BooleanLess:
		return enum.IPredSLT
	case decorated.BooleanLessOrEqual:
		return enum.IPredSLE
	case decorated.BooleanGreater:
		return enum.IPredSGT
	case decorated.BooleanGreaterOrEqual:
		return enum.IPredSGE
	default:
		panic(fmt.Errorf("not allowed int operator type"))
	}

	return 0
}

func generateBinaryOperatorBooleanResult(operator *decorated.BooleanOperator, genContext *generateContext) (
	value.Value, error,
) {
	leftVar, leftErr := generateExpression(operator.Left(), false, genContext)
	if leftErr != nil {
		return nil, leftErr
	}

	rightVar, rightErr := generateExpression(operator.Right(), false, genContext)
	if rightErr != nil {
		return nil, rightErr
	}

	unaliasedTypeLeft := dectype.ResolveToAtom(operator.Left().Type())
	foundPrimitive, _ := unaliasedTypeLeft.(*dectype.PrimitiveAtom)
	if foundPrimitive == nil {
		foundCustomType, _ := unaliasedTypeLeft.(*dectype.CustomTypeAtom)
		if foundCustomType == nil {
			panic(fmt.Errorf("not implemented binary operator boolean %v", operator.Left().Type().HumanReadable()))
		} else {
			zeroIndex := constant.NewInt(types.I32, 0)
			var indices []value.Value
			indices = append(indices, zeroIndex)
			indices = append(indices, zeroIndex)
			leftPtrToOctet := ir.NewGetElementPtr(types.I8Ptr, leftVar, indices...)
			rightPtrToOctet := ir.NewGetElementPtr(types.I8Ptr, rightVar, indices...)
			predicate := booleanIntToIPred(operator.OperatorType())
			leftOctet := ir.NewLoad(types.I8, leftPtrToOctet)
			rightOctet := ir.NewLoad(types.I8, rightPtrToOctet)

			ir.NewICmp(predicate, leftOctet, rightOctet)
		}
	} else if foundPrimitive.AtomName() == "String" {
		stringCompare := &ir.Func{} // TODO: Fetch string compare functions
		return ir.NewCall(stringCompare, leftVar, rightVar), nil
	} else if foundPrimitive.AtomName() == "Int" || foundPrimitive.AtomName() == "Char" || foundPrimitive.AtomName() == "Fixed" {
		predicate := booleanIntToIPred(operator.OperatorType())
		log.Printf("%T left:%T, right:%T", operator.Left(), leftVar, rightVar)
		return ir.NewICmp(predicate, leftVar, rightVar), nil
	} else if foundPrimitive.AtomName() == "Bool" {
		isNotEqual := genContext.block.NewXor(leftVar, rightVar)
		if operator.OperatorType() == decorated.BooleanEqual {
			return ir.NewXor(isNotEqual, constant.NewBool(true)), nil
		} else if operator.OperatorType() == decorated.BooleanNotEqual {
			return isNotEqual, nil
		} else {
			panic(fmt.Errorf("illegal boolean operator for bool %v", operator))
		}
	} else {
		panic(fmt.Errorf("generate sp: what operator is this for %v", foundPrimitive.AtomName()))
	}

	panic(fmt.Errorf("unknwon"))
}
