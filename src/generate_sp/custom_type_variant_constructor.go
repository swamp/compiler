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
	unaliasedTypeVariant := dectype.UnaliasWithResolveInvoker(constructor.Type())
	smashedVariant := unaliasedTypeVariant.(*dectype.CustomTypeVariantAtom)

	variantMemorySize, _ := dectype.GetMemorySizeAndAlignment(smashedVariant.InCustomType())
	if uint(variantMemorySize) > uint(target.Size) {
		log.Printf("smashedVariant:%v\n\nsmashedCustomType:%v\n\n", smashedVariant, smashedVariant.InCustomType())
		return fmt.Errorf("internal error, target size is not exactly right at %v, target is:%v and unionMemorySize is:%v", constructor.FetchPositionLength().ToCompleteReferenceString(), target.Size, variantMemorySize)
	}

	filePosition := genContext.toFilePosition(constructor.FetchPositionLength())
	log.Printf("setting enum %d for size:%d", constructor.CustomTypeVariantIndex(), target.Size)
	code.SetEnum(target.Pos, uint8(constructor.CustomTypeVariant().Index()), target.Size, filePosition)

	for index, arg := range constructor.Arguments() {
		variantField := smashedVariant.Fields()[index]
		log.Printf(" setting enum argument %d at %d with size %d", index, variantField.MemoryOffset(), variantField.MemorySize())
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
	unaliasedTypeVariant := dectype.UnaliasWithResolveInvoker(customTypeVariantConstructor.Type())
	smashedVariant := unaliasedTypeVariant.(*dectype.CustomTypeVariantAtom)
	posRange := allocMemoryForType(genContext.context.stackMemory, smashedVariant.InCustomType(), "variant constructor target")
	if err := generateCustomTypeVariantConstructor(code, posRange, customTypeVariantConstructor, genContext); err != nil {
		return assembler_sp.SourceStackPosRange{}, err
	}

	return targetToSourceStackPosRange(posRange), nil
}
