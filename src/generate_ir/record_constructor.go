package generate_ir

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"log"
)

func generateRecordConstructorSortedAssignmentsHelper(
	recordType *dectype.RecordAtom, sortedAssignments []*decorated.RecordLiteralAssignment, genContext *generateContext) (value.Value, error) {
	irType, lookupErr := genContext.irTypeRepo.GetTypeRef(recordType)
	if lookupErr != nil {
		return nil, lookupErr
	}

	instAlloc := ir.NewAlloca(irType)
	log.Printf("record Constructor : %v, %v", irType, instAlloc)
	for index, assignment := range sortedAssignments {
		//recordField := recordType.SortedFields()[index]
		sourceValue, genErr := generateExpression(assignment.Expression(), false, genContext)
		if genErr != nil {
			return nil, genErr
		}
		zeroIndex := constant.NewInt(types.I32, int64(0))
		constantIndex := constant.NewInt(types.I32, int64(index))
		destination := ir.NewGetElementPtr(irType, instAlloc, zeroIndex, constantIndex)
		log.Printf("record constructor store destination:%v", destination)
		stored := ir.NewStore(sourceValue, destination)
		log.Printf("record constructor store:%v", stored)
	}

	return instAlloc, nil
}

func generateRecordConstructorSortedAssignments(recordConstructor *decorated.RecordConstructorFromParameters, genContext *generateContext) (value.Value, error) {
	recordType := recordConstructor.RecordType()
	return generateRecordConstructorSortedAssignmentsHelper(recordType, recordConstructor.SortedAssignments(), genContext)
}
