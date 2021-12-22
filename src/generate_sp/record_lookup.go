package generate_sp

import (
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func handleRecordLookup(code *assembler_sp.Code, lookups *decorated.RecordLookups,
	genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	startOfStruct, err := generateExpressionWithSourceVar(code, lookups.Expression(), genContext, "lookups")
	if err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	indexOffset := uint(0)

	var lastLookup decorated.LookupField
	for _, indexLookup := range lookups.LookupFields() {
		indexOffset += uint(indexLookup.MemoryOffset())
		lastLookup = indexLookup
	}

	sourcePosRange := assembler_sp.SourceStackPosRange{
		Pos:  assembler_sp.SourceStackPos(uint(startOfStruct.Pos) + indexOffset),
		Size: assembler_sp.SourceStackRange(lastLookup.MemorySize()),
	}

	return sourcePosRange, nil
}

func generateLookups(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, lookups *decorated.RecordLookups,
	genContext *generateContext) error {
	sourcePosRange, err := handleRecordLookup(code, lookups, genContext)
	if err != nil {
		return err
	}

	filePosition := genContext.toFilePosition(lookups.FetchPositionLength())
	code.CopyMemory(target.Pos, sourcePosRange, filePosition)

	return nil
}
