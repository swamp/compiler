/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"fmt"

	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/typeinfo"
	opcode_sp_type "github.com/swamp/opcodes/type"
)

func preparePackageConstants(compiledPackage *loader.Package, packageConstants *assembler_sp.PackageConstants,
	typeInformationChunk typeinfo.TypeLookup) decshared.DecoratedError {
	for _, module := range compiledPackage.AllModules() {
		for _, named := range module.LocalDefinitions().Definitions() {
			unknownExpression := named.Expression()
			maybeFunction, _ := unknownExpression.(*decorated.FunctionValue)
			if maybeFunction != nil {
				fullyQualifiedName := module.FullyQualifiedName(named.Identifier())
				isExternal := maybeFunction.IsSomeKindOfExternal()
				if isExternal {
					var paramPosRanges []assembler_sp.SourceStackPosRange
					hasLocalTypes := dectype.TypeIsTemplateHasLocalTypes(maybeFunction.Type())
					// parameterCount := len(maybeFunction.Parameters())
					pos := dectype.MemoryOffset(0)
					if hasLocalTypes {
						returnPosRange := assembler_sp.SourceStackPosRange{
							Pos:  assembler_sp.SourceStackPos(0),
							Size: assembler_sp.SourceStackRange(0),
						}
						paramPosRanges = make([]assembler_sp.SourceStackPosRange, len(maybeFunction.Parameters()))
						if _, err := packageConstants.AllocatePrepareExternalFunctionConstant(fullyQualifiedName.String(),
							returnPosRange, paramPosRanges); err != nil {
							return decorated.NewInternalError(err)
						}
						continue
					}
					returnSize, _ := dectype.GetMemorySizeAndAlignment(dectype.ResolveToFunctionAtom(maybeFunction.Type()).ReturnType())
					returnPosRange := assembler_sp.SourceStackPosRange{
						Pos:  assembler_sp.SourceStackPos(pos),
						Size: assembler_sp.SourceStackRange(returnSize),
					}

					pos += dectype.MemoryOffset(returnSize)

					parameterTypes, _ := dectype.ResolveToFunctionAtom(maybeFunction.Type()).ParameterAndReturn()

					for _, param := range parameterTypes {
						unaliased := dectype.Unalias(param)
						if dectype.ArgumentNeedsTypeIdInsertedBefore(unaliased) || dectype.IsTypeIdRef(unaliased) {
							pos = align(pos, dectype.MemoryAlign(opcode_sp_type.AlignOfSwampInt))
							typeIndexPosRange := assembler_sp.SourceStackPosRange{
								Pos:  assembler_sp.SourceStackPos(pos),
								Size: assembler_sp.SourceStackRange(opcode_sp_type.SizeofSwampInt),
							}
							paramPosRanges = append(paramPosRanges, typeIndexPosRange)
							pos += dectype.MemoryOffset(typeIndexPosRange.Size)
							if dectype.IsTypeIdRef(unaliased) {
								continue
							}
						}
						size, alignment := dectype.GetMemorySizeAndAlignment(param)
						pos = align(pos, alignment)
						posRange := assembler_sp.SourceStackPosRange{
							Pos:  assembler_sp.SourceStackPos(pos),
							Size: assembler_sp.SourceStackRange(size),
						}
						paramPosRanges = append(paramPosRanges, posRange)
						pos += dectype.MemoryOffset(size)
					}

					if _, err := packageConstants.AllocatePrepareExternalFunctionConstant(fullyQualifiedName.String(),
						returnPosRange, paramPosRanges); err != nil {
						return decorated.NewInternalError(err)
					}
				} else {
					// parameterTypes, _ := maybeFunction.UnaliasedDeclaredFunctionType().ParameterAndReturn()
					returnSize, returnAlign := dectype.GetMemorySizeAndAlignment(dectype.ResolveToFunctionAtom(maybeFunction.Type()).ReturnType())
					parameterCount := uint(len(maybeFunction.Parameters())) // parameterTypes

					functionTypeIndex, lookupErr := typeInformationChunk.Lookup(maybeFunction.Type())
					if lookupErr != nil {
						return decorated.NewInternalError(lookupErr)
					}

					pos := dectype.MemoryOffset(0)
					for _, param := range maybeFunction.Parameters() {
						paramSize, paramAlign := dectype.GetMemorySizeAndAlignment(param.Type())
						pos = align(pos, paramAlign)
						pos += dectype.MemoryOffset(paramSize)
					}
					parameterOctetSize := dectype.MemorySize(pos)
					if _, err := packageConstants.AllocatePrepareFunctionConstant(fullyQualifiedName.String(),
						opcode_sp_type.MemorySize(returnSize), opcode_sp_type.MemoryAlign(returnAlign), parameterCount,
						opcode_sp_type.MemorySize(parameterOctetSize), uint(functionTypeIndex)); err != nil {
						return decorated.NewInternalError(err)
					}

				}
			} else {
				if _, isConstant := unknownExpression.(*decorated.Constant); !isConstant {
					panic(fmt.Errorf("unknown thing here: %T", unknownExpression))
				}
			}
		}
	}

	return nil
}
