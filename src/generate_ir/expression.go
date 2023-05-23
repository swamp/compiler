/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_ir

import (
	"fmt"
	"log"

	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateExpressionHelper(expr decorated.Expression, isLeafNode bool, genContext *generateContext) (
	value.Value, error,
) {
	switch e := expr.(type) {
	case *decorated.Let:
		return generateLet(e, genContext)
	case *decorated.ArithmeticOperator:
		return generateArithmeticMultiple(e, genContext)
	case *decorated.RecordConstructorFromParameters:
		return generateRecordConstructorSortedAssignments(e, genContext)
	case *decorated.RecordLookups:
		return generateLookups(e, genContext)
	case *decorated.FunctionParameterReference:
		return generateLocalFunctionParameterReference(e, genContext)
	case *decorated.CaseCustomType:
		return generateCaseCustomType(e, genContext)
	case *decorated.BooleanOperator:
		return generateBinaryOperatorBooleanResult(e, genContext)
	case *decorated.If:
		return generateIf(e, isLeafNode, genContext)
	case *decorated.IntegerLiteral:
		return generateIntLiteral(e, genContext)
	case *decorated.LogicalOperator:
		return generateLogical(e, genContext)
	case *decorated.FunctionCall:
		return generateFunctionCall(e, isLeafNode, genContext)
	case *decorated.FunctionReference:
		return generateFunctionReference(e, genContext)
	}

	/*
		case *decorated.BitwiseOperator:
			return generateBitwise(code, target, e, genContext)

		case *decorated.BitwiseUnaryOperator:
			return generateUnaryBitwise(code, target, e, genContext)

		case *decorated.LogicalUnaryOperator:
			return generateUnaryLogical(code, target, e, genContext)

		case *decorated.ArithmeticUnaryOperator:
			return generateUnaryArithmetic(code, target, e, genContext)



		case *decorated.PipeLeftOperator:
			return generatePipeLeft(code, target, e, genContext)

		case *decorated.PipeRightOperator:
			return generatePipeRight(code, target, e, genContext)


		case *decorated.CaseForPatternMatching:
			return generateCasePatternMatchingMultiple(code, target, e, genContext)

		case *decorated.RecordLiteral:
			return generateRecordLiteral(code, target, e, genContext)


		case *decorated.Guard:
			return generateGuard(code, target, e, genContext)

		case *decorated.StringLiteral:
			return generateStringLiteral(code, target, e, genContext)

		case *decorated.CharacterLiteral:
			return generateCharacterLiteral(code, target, e, genContext)

		case *decorated.TypeIdLiteral:
			return generateTypeIdLiteral(code, target, e, genContext)


		case *decorated.FixedLiteral:
			return generateFixedLiteral(code, target, e, genContext)

		case *decorated.ResourceNameLiteral:
			return generateResourceNameLiteral(code, target, e, genContext)

		case *decorated.BooleanLiteral:
			return generateBoolLiteral(code, target, e, genContext)

		case *decorated.ListLiteral:
			return generateList(code, target, e, genContext)

		case *decorated.TupleLiteral:
			return generateTuple(code, target, e, genContext)

		case *decorated.ArrayLiteral:
			return generateArray(code, target, e, genContext)

		case *decorated.RecurCall:
			return generateRecurCall(code, e, genContext)

		case *decorated.CurryFunction:
			return generateCurry(code, target, e, genContext)

		case *decorated.StringInterpolation:
			return generateExpression(code, target, e.Expression(), leafNode, genContext)

		case *decorated.CustomTypeVariantConstructor:
			return generateCustomTypeVariantConstructor(code, target, e, genContext)

		case *decorated.Constant:
			return generateConstant(code, target, e, genContext)

		case *decorated.ConstantReference:
			return generateExpression(code, target, e.Constant(), leafNode, genContext)

		case *decorated.LetVariableReference:
			return generateLetVariableReference(code, target, e, genContext)

		case *decorated.CaseConsequenceParameterReference:
			return generateLocalConsequenceParameterReference(code, target, e, genContext)

		case *decorated.ConsOperator:
			return generateListCons(code, target, e, genContext)

		case *decorated.RecordConstructorFromRecord:
			return generateExpression(code, target, e.Expression(), leafNode, genContext)


		case *decorated.CastOperator:
			return generateExpression(code, target, e.Expression(), leafNode, genContext)
		}

		panic(fmt.Errorf("generate_sp: unknown node %T %v %v", expr, expr, genContext))

	*/
	panic(fmt.Errorf("not implemented yet %T", expr))
	return nil, fmt.Errorf("not implemented yet %T", expr)

}

func generateExpression(expr decorated.Expression, isLeafNode bool, genContext *generateContext) (value.Value, error) {
	log.Printf("!generating %T %v", expr, expr)

	result, resultErr := generateExpressionHelper(expr, isLeafNode, genContext)

	log.Printf("result from %T:  %v, err:%v", expr, result, resultErr)

	return result, resultErr
}
