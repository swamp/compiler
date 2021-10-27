package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateTuple(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	tupleLiteral *decorated.TupleLiteral, genContext *generateContext) error {
	fields := tupleLiteral.TupleType().Fields()
	for index, expr := range tupleLiteral.Expressions() {
		tupleField := fields[index]
		fieldTarget := assembler_sp.TargetStackPosRange{
			Pos:  assembler_sp.TargetStackPos(uint(target.Pos) + uint(tupleField.MemoryOffset())),
			Size: assembler_sp.StackRange(tupleField.MemorySize()),
		}

		genErr := generateExpression(code, fieldTarget, expr, false, genContext)
		if genErr != nil {
			return genErr
		}
	}

	return nil
}

func handleTuple(code *assembler_sp.Code,
	tupleLiteral *decorated.TupleLiteral, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	tuplePointer := genContext.context.stackMemory.Allocate(uint(tupleLiteral.TupleType().MemorySize()),
		uint32(tupleLiteral.TupleType().MemoryAlignment()), "tuple")

	if err := generateTuple(code, tuplePointer, tupleLiteral, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(tuplePointer), nil
}
