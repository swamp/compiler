package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func isIntLike(typeToCheck dtype.Type) bool {
	unaliasType := dectype.UnaliasWithResolveInvoker(typeToCheck)

	primitiveAtom, _ := unaliasType.(*dectype.PrimitiveAtom)
	if primitiveAtom == nil {
		return false
	}

	name := primitiveAtom.AtomName()

	return name == "Int" || name == "Fixed" || name == "Char"
}

func isListLike(typeToCheck dtype.Type) bool {
	unaliasType := dectype.UnaliasWithResolveInvoker(typeToCheck)

	primitiveAtom, _ := unaliasType.(*dectype.PrimitiveAtom)
	if primitiveAtom == nil {
		return false
	}

	name := primitiveAtom.PrimitiveName().Name()

	return name == "List"
}

func handleArithmeticMultiple(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, e *decorated.ArithmeticOperator,
	genContext *generateContext) error {
	leftPrimitive, _ := dectype.UnReference(e.Left().Type()).(*dectype.PrimitiveAtom)
	switch {
	case isListLike(e.Left().Type()) && e.OperatorType() == decorated.ArithmeticAppend:
		return generateListAppend(code, target, e, genContext)
	case leftPrimitive != nil && leftPrimitive.AtomName() == "String" && e.OperatorType() == decorated.ArithmeticAppend:
		return generateStringAppend(code, target, e, genContext)
	case isIntLike(e.Left().Type()):
		return generateArithmetic(code, target, e, genContext)
	default:
		return fmt.Errorf("cant generate arithmetic for type: %v <-> %v (%v)",
			e.Left().Type(), e.Right().Type(), e.OperatorType())
	}
}