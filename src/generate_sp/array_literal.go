package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateArray(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	array *decorated.ArrayLiteral, genContext *generateContext) error {
	variables := make([]assembler_sp.SourceStackPos, len(array.Expressions()))
	for index, expr := range array.Expressions() {
		debugName := fmt.Sprintf("arrayliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar.Pos
	}
	primitive, _ := array.Type().(*dectype.PrimitiveAtom)
	firstPrimitiveType := primitive.GenericTypes()[0]
	itemSize, _ := dectype.GetMemorySizeAndAlignment(firstPrimitiveType)
	code.ArrayLiteral(target.Pos, variables, assembler_sp.StackRange(itemSize))
	return nil
}
