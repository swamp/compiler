package generate_sp

import (
	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateTuple(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	tupleLiteral *decorated.TupleLiteral, genContext *generateContext) error {
	tuplePointer := genContext.context.stackMemory.Allocate(tupleLiteral.TupleType().MemorySize(),
		tupleLiteral.TupleType().MemoryAlignment(), "tuple")
	fields := tupleLiteral.TupleType().Fields()
	for index, expr := range tupleLiteral.Expressions() {
		tupleField := fields[index]
		fieldTarget := assembler_sp.TargetStackPosRange{
			Pos:  assembler_sp.TargetStackPos(uint(tuplePointer.Pos) + tupleField.MemoryOffset()),
			Size: assembler_sp.StackRange(tupleField.MemorySize()),
		}

		genErr := generateExpression(code, fieldTarget, expr, genContext)
		if genErr != nil {
			return genErr
		}
	}

	return nil
}
