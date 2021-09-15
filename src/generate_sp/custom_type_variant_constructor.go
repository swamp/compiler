package generate_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func customTypeVariantConstructor() {
	/*

		unionMemorySize := variant.InCustomType().MemorySize()
		unionMemoryAlignment := variant.InCustomType().MemoryAlignment()
		unionMemory := genContext.context.stackMemory.Allocate(unionMemorySize, unionMemoryAlignment, "variant constructor")
	*/
}

func generateCustomTypeVariantConstructor(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	constructor *decorated.CustomTypeVariantConstructor, genContext *generateContext) error {
	variant := constructor.CustomTypeVariant()
	unionMemorySize := variant.InCustomType().MemorySize()
	if uint(target.Size) != unionMemorySize {
		return fmt.Errorf("internal error, target size is not exactly right")
	}

	for index, arg := range constructor.Arguments() {
		variantField := variant.Fields()[index]
		fieldTarget := assembler_sp.TargetStackPosRange{
			Pos:  assembler_sp.TargetStackPos(uint(target.Pos) + variantField.MemoryOffset()),
			Size: assembler_sp.StackRange(variantField.MemorySize()),
		}
		argRegErr := generateExpression(code, fieldTarget, arg, genContext)
		if argRegErr != nil {
			return argRegErr
		}
	}

	return nil
}
