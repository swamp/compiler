package generate_c

import (
	"fmt"
	"io"

	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateExpression(expr decorated.Expression, writer io.Writer, returnPrefix string, indentation int) error {
	switch e := expr.(type) {
	case *decorated.Let:
		return generateLet(e, writer, returnPrefix, indentation)

	/*
		case *decorated.ArithmeticOperator:
			return generateArithmeticMultiple(code, target, e, genContext)

		case *decorated.BitwiseOperator:
			return generateBitwise(code, target, e, genContext)

		case *decorated.BitwiseUnaryOperator:
			return generateUnaryBitwise(code, target, e, genContext)

		case *decorated.LogicalUnaryOperator:
			return generateUnaryLogical(code, target, e, genContext)

		case *decorated.ArithmeticUnaryOperator:
			return generateUnaryArithmetic(code, target, e, genContext)


	*/
	case *decorated.LogicalOperator:
		return generateLogical(e, writer, indentation)

	case *decorated.BooleanOperator:
		return generateBinaryOperatorBooleanResult(e, writer, indentation)
	case *decorated.If:
		return generateIf(e, writer, returnPrefix, indentation)
		/*
			case *decorated.PipeLeftOperator:
				return generatePipeLeft(code, target, e, genContext)

			case *decorated.PipeRightOperator:
				return generatePipeRight(code, target, e, genContext)

			case *decorated.RecordLookups:
				return generateLookups(code, target, e, genContext)

			case *decorated.CaseCustomType:
				return generateCaseCustomType(code, target, e, genContext)

			case *decorated.CaseForPatternMatching:
				return generateCasePatternMatchingMultiple(code, target, e, genContext)

			case *decorated.RecordLiteral:
				return generateRecordLiteral(code, target, e, genContext)

			case *decorated.Guard:
				return generateGuard(code, target, e, genContext)

			case *decorated.TypeIdLiteral:
				return generateTypeIdLiteral(code, target, e, genContext)
		*/
	case *decorated.CharacterLiteral:
		return generateCharacterLiteral(e, writer, indentation)
	case *decorated.StringLiteral:
		return generateStringLiteral(e, writer, indentation)
	case *decorated.IntegerLiteral:
		return generateIntLiteral(e, writer, indentation)
		/*
			case *decorated.FixedLiteral:
				return generateFixedLiteral(code, target, e)

			case *decorated.ResourceNameLiteral:
				return generateResourceNameLiteral(code, target, e, genContext.context.Constants())
		*/
	case *decorated.BooleanLiteral:
		return generateBoolLiteral(e, writer, indentation)
		/*
			case *decorated.ListLiteral:
				return generateList(code, target, e, genContext)

			case *decorated.TupleLiteral:
				return generateTuple(code, target, e, genContext)

			case *decorated.ArrayLiteral:
				return generateArray(code, target, e, genContext)

			case *decorated.FunctionCall:
				return generateFunctionCall(code, target, e, genContext)

			case *decorated.RecurCall:
				return generateRecurCall(code, e, genContext)

			case *decorated.CurryFunction:
				return generateCurry(code, target, e, genContext)

			case *decorated.StringInterpolation:
				return generateExpression(code, target, e.Expression(), genContext)

			case *decorated.CustomTypeVariantConstructor:
				return generateCustomTypeVariantConstructor(code, target, e, genContext)

			case *decorated.Constant:
				return generateConstant(code, target, e, genContext)

			case *decorated.ConstantReference:
				return generateExpression(code, target, e.Constant(), genContext)
		*/
	case *decorated.FunctionParameterReference:
		return generateLocalFunctionParameterReference(e, writer, indentation)
	case *decorated.LetVariableReference:
		return generateLetVariableReference(e, writer, indentation)
		/*
			case *decorated.FunctionReference:
				return generateFunctionReference(code, target, e, genContext.context.constants)

			case *decorated.CaseConsequenceParameterReference:
				return generateLocalConsequenceParameterReference(code, target, e, genContext.context)

			case *decorated.ConsOperator:
				return generateListCons(code, target, e, genContext)

			case *decorated.RecordConstructorFromRecord:
				return generateExpression(code, target, e.Expression(), genContext)

			case *decorated.RecordConstructorFromParameters:
				return generateRecordConstructorSortedAssignments(code, target, e, genContext)

		*/
	}

	panic(fmt.Errorf("generate_sp: unknown node %T %v", expr, expr))
}
