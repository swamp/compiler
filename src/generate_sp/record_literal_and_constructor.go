package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateRecordConstructorSortedAssignments(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	recordConstructor *decorated.RecordConstructorFromParameters, genContext *generateContext) error {
	recordType := recordConstructor.RecordType()

	if uint(target.Size) != recordType.MemorySize() {
		return fmt.Errorf("wrong target size for record constructor")
	}

	for index, assignment := range recordConstructor.SortedAssignments() {
		recordField := recordType.SortedFields()[index]
		fieldTarget := createTargetWithMemoryOffsetAndSize(target, recordField.MemoryOffset(), recordField.MemorySize())
		genErr := generateExpression(code, fieldTarget, assignment.Expression(), genContext)
		if genErr != nil {
			return genErr
		}
	}

	return nil
}

func generateRecordLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	record *decorated.RecordLiteral, genContext *generateContext) error {
	if record.RecordTemplate() != nil {
		recordType := record.RecordType()
		if uint(target.Size) != recordType.MemorySize() {
			return fmt.Errorf("target for record literal is wrong")
		}

		structToCopyVar, genErr := generateExpressionWithSourceVar(code, record.RecordTemplate(), genContext, "gopher")
		if genErr != nil {
			return genErr
		}

		code.CopyMemory(target.Pos, structToCopyVar)

		recordFields := recordType.SortedFields()
		for index, assignment := range record.SortedAssignments() {
			recordField := recordFields[index]
			fieldTarget := createTargetWithMemoryOffsetAndSize(target, recordField.MemoryOffset(), recordField.MemorySize())
			genErr := generateExpression(code, fieldTarget, assignment.Expression(), genContext)
			if genErr != nil {
				return genErr
			}
		}
	} else {
		return generateRecordConstructorSortedAssignments(code, target, nil, genContext)
	}
	return nil
}
