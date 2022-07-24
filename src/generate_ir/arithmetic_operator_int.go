package generate_ir

import (
	"fmt"
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateArithmeticInt(operator *decorated.ArithmeticOperator, genContext *generateContext) (value.Value, error) {
	irValueLeft, leftErr := generateExpression(operator.Left(), false, genContext)
	if leftErr != nil {
		return nil, leftErr
	}
	irValueRight, rightErr := generateExpression(operator.Left(), false, genContext)
	if rightErr != nil {
		return nil, rightErr
	}

	log.Printf("left is %T and right is %T", irValueLeft, irValueRight)

	switch operator.OperatorType() {
	case decorated.ArithmeticPlus:
		return ir.NewAdd(irValueLeft, irValueRight), nil
	case decorated.ArithmeticMinus:
		return ir.NewSub(irValueLeft, irValueRight), nil
	case decorated.ArithmeticMultiply:
		return ir.NewMul(irValueLeft, irValueRight), nil
	case decorated.ArithmeticDivide:
		return ir.NewSDiv(irValueLeft, irValueRight), nil
	case decorated.ArithmeticRemainder:
		return ir.NewSRem(irValueLeft, irValueRight), nil
	//  case decorated.ArithmeticAppend: // n/a for ints
	//	case decorated.ArithmeticCons: // n/a for ints
	case decorated.ArithmeticFixedMultiply:
		return ir.NewSDiv(ir.NewMul(irValueLeft, irValueRight), constant.NewInt(types.I32, 100)), nil
	case decorated.ArithmeticFixedDivide:
		return ir.NewSDiv(ir.NewSRem(irValueLeft, irValueRight), constant.NewInt(types.I32, 100)), nil
	default:
		return nil, fmt.Errorf("unknown int operator")
	}
}
