package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateList(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	list *decorated.ListLiteral, genContext *generateContext) error {
	variables := make([]assembler_sp.SourceStackPos, len(list.Expressions()))
	for index, expr := range list.Expressions() {
		debugName := fmt.Sprintf("listliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar.Pos
	}
	primitive, _ := list.Type().(*dectype.PrimitiveAtom)
	itemSize, _ := getMemorySizeAndAlignment(primitive.Next())
	code.ListLiteral(target.Pos, variables, assembler_sp.StackRange(itemSize))
	return nil
}
