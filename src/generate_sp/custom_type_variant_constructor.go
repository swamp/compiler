package generate_sp

import (
	"fmt"
	"log"

	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func generateCustomTypeVariantConstructor(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	constructor *decorated.CustomTypeVariantConstructor, genContext *generateContext) error {
	smashedCustomType := constructor.Type().(*dectype.CustomTypeAtom)

	smashedVariant := smashedCustomType.FindVariant(constructor.CustomTypeVariant().Name().Name())
	unionMemorySize, _ := dectype.GetMemorySizeAndAlignment(smashedCustomType)
	if smashedVariant.Name().Name() != "Nothing" && uint(target.Size) != uint(unionMemorySize) {
		log.Printf("smashedVariant:%v\n\nsmashedCustomType:%v\n\n", smashedVariant, smashedCustomType)
		return fmt.Errorf("internal error, target size is not exactly right at %v, target is:%v and unionMemorySize is:%v", constructor.FetchPositionLength().ToCompleteReferenceString(), target.Size, unionMemorySize)
	}

	filePosition := genContext.toFilePosition(constructor.FetchPositionLength())
	code.SetEnum(target.Pos, uint8(constructor.CustomTypeVariant().Index()), target.Size, filePosition)

	for index, arg := range constructor.Arguments() {
		variantField := smashedVariant.Fields()[index]
		fieldTarget := assembler_sp.TargetStackPosRange{
			Pos:  assembler_sp.TargetStackPos(uint(target.Pos) + uint(variantField.MemoryOffset())),
			Size: assembler_sp.StackRange(variantField.MemorySize()),
		}
		argRegErr := generateExpression(code, fieldTarget, arg, false, genContext)
		if argRegErr != nil {
			log.Printf("encountered error %v", argRegErr)
			return argRegErr
		}
	}

	return nil
}

func handleCustomTypeVariantConstructor(code *assembler_sp.Code,
	customTypeVariantConstructor *decorated.CustomTypeVariantConstructor, genContext *generateContext) (assembler_sp.SourceStackPosRange, error) {
	posRange := allocMemoryForType(genContext.context.stackMemory, customTypeVariantConstructor.CustomTypeVariant().InCustomType(), "variant constructor target")
	if err := generateCustomTypeVariantConstructor(code, posRange, customTypeVariantConstructor, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
