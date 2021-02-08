/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"fmt"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

const listContent = `
__externalfn coreListMap 2
__externalfn coreListMap2 3
__externalfn coreListAny 2
__externalfn coreListFilter 2
__externalfn coreListFilterMap 2
__externalfn coreListRemove 2
__externalfn coreListMapcat 2
__externalfn coreListHead 1
__externalfn coreListConcat 1
__externalfn coreListIsEmpty 1
__externalfn coreListRange 2
__externalfn coreListLength 1
__externalfn coreListNth 2
__externalfn coreListTail 1
__externalfn coreListReduce 2
__externalfn coreListFoldl 3
__externalfn coreListFoldlStop 3


map : (a -> b) -> List a -> List b
map x y =
    __asm callexternal 00 coreListMap 01 02


map2 : (a -> b -> c) -> List a -> List b -> List c
map2 x y z =
    __asm callexternal 00 coreListMap2 01 02 03


__externalfn coreListIndexedMap 2
indexedMap : (Int -> a -> b) -> List a -> List b
indexedMap x y =
    __asm callexternal 00 coreListIndexedMap 01 02


any : (a -> Bool) -> List a -> Bool
any pred coll =
    __asm callexternal 00 coreListAny 01 02


__externalfn coreListFind 2
find : (a -> Bool) -> List a -> Maybe a
find pred coll =
    __asm callexternal 00 coreListFind 01 02


__externalfn coreListMember 2
member : a -> List a -> Bool
member pred coll =
    __asm callexternal 00 coreListMember 01 02


filterMap : (a -> Maybe b) -> List a -> List b
filterMap pred coll =
    __asm callexternal 00 coreListFilterMap 01 02


filter : (a -> Bool) -> List a -> List a
filter pred coll =
    __asm callexternal 00 coreListFilter 01 02


remove : (a -> Bool) -> List a -> List a
remove pred coll =
    __asm callexternal 00 coreListRemove 01 02


concatMap : (a -> List b) -> List a -> List b
concatMap conv lst =
    __asm callexternal 00 coreListMapcat 01 02


concat : List (List a) -> List a
concat lst =
    __asm callexternal 00 coreListConcat 01


isEmpty : List a -> Bool
isEmpty lst =
    __asm callexternal 00 coreListIsEmpty 01


range : Int -> Int -> List Int
range start end =
    __asm callexternal 00 coreListRange 01 02


length : List a -> Int
length lst =
    __asm callexternal 00 coreListLength 01


foldl : (a -> b -> b) -> b -> List a -> b
foldl fn initial lst =
    __asm callexternal 00 coreListFoldl 01 02 03


foldlstop : (a -> b -> Maybe b) -> b -> List a -> b
foldlstop pred initial lst =
    __asm callexternal 00 coreListFoldlStop 01 02 03


tail : List a -> Maybe (List a)
tail lst =
    __asm callexternal 00 coreListTail 01


head : List a -> Maybe a
head lst =
    __asm callexternal 00 coreListHead 01
`

const mathCode = `
__externalfn coreMathRemainderBy 2
remainderBy : Int -> Int -> Int
remainderBy rem v =
    __asm callexternal 00 coreMathRemainderBy 01 02


__externalfn coreMathSin 1
sin : Fixed -> Fixed
sin angle =
    __asm callexternal 00 coreMathSin 01


__externalfn coreMathCos 1
cos : Fixed -> Fixed
cos angle =
    __asm callexternal 00 coreMathCos 01


__externalfn coreMathRnd 2
rnd : Int -> Int -> Int
rnd t mod =
    __asm callexternal 00 coreMathRnd 01 02


__externalfn coreMathAtan2 2
atan2 : Int -> Int -> Fixed
atan2 x y =
    __asm callexternal 00 coreMathAtan2 01 02


__externalfn coreMathMid 2
mid : Int -> Int -> Int
mid x y =
    __asm callexternal 00 coreMathMid 01 02


__externalfn coreMathAbs 1
abs : Int -> Int
abs x =
    __asm callexternal 00 coreMathAbs 01


__externalfn coreMathSign 1
sign : Int -> Int
sign x =
    __asm callexternal 00 coreMathSign 01


__externalfn coreMathClamp 3
clamp : Int -> Int -> Int -> Int
clamp v min max =
    __asm callexternal 00 coreMathClamp 01 02 03


__externalfn coreMathLerp 3
lerp : Fixed -> Int -> Int -> Int
lerp t from to =
    __asm callexternal 00 coreMathLerp 01 02 03


__externalfn coreMathMetronome 4
metronome : Int -> Int -> Int -> Int -> Bool
metronome a t from to =
    __asm callexternal 00 coreMathMetronome 01 02 03 04


__externalfn coreMathDrunk 3
drunk : Int -> Int -> Int -> Int
drunk t from to =
    __asm callexternal 00 coreMathDrunk 01 02 03


__externalfn coreMathMod 2
mod : Int -> Int -> Int
mod rem v =
    __asm callexternal 00 coreMathMod 01 02
`

const blobCode = `
__externalfn coreBlobIsEmpty 1
isEmpty : Blob -> Bool
isEmpty x =
    __asm callexternal 00 coreBlobIsEmpty 01


__externalfn coreBlobFromList 1
fromList : List Int -> Blob
fromList x =
    __asm callexternal 00 coreBlobFromList 01
`

const arrayCode = `
__externalfn coreArrayFromList 1
fromList : List a -> Array a
fromList list =
    __asm callexternal 00 coreArrayFromList 01


__externalfn coreArrayToList 1
toList : Array a -> List a
toList list =
    __asm callexternal 00 coreArrayToList 01


__externalfn coreArraySlice 3
slice : Int -> Int -> Array a -> Array a
slice startIndex endIndex array =
    __asm callexternal 00 coreArraySlice 01 02 03


__externalfn coreArrayRepeat 2
repeat : Int -> a -> Array a
repeat count item =
    __asm callexternal 00 coreArrayRepeat 01 02


__externalfn coreArrayLength 1
length : Array a -> Int
length lst =
    __asm callexternal 00 coreArrayLength 01


__externalfn coreArrayGet 2
get : Int -> Array a -> Maybe a
get index array =
    __asm callexternal 00 coreArrayGet 01 02


__externalfn coreArrayGrab 2
grab : Int -> Array a -> a
grab index array =
    __asm callexternal 00 coreArrayGrab 01 02


__externalfn coreArraySet 3
set : Int -> a -> Array a -> Array a
set index item array =
    __asm callexternal 00 coreArraySet 01 02 03

`

const maybeCode = `
__externalfn coreMaybeWithDefault 2
withDefault : a -> Maybe a -> a
withDefault default maybe =
    __asm callexternal 00 coreMaybeWithDefault 01 02
`

const debugCode = `
__externalfn coreDebugLog 1
log : String -> String
log output =
    __asm callexternal 00 coreDebugLog 01


__externalfn coreDebugToString 1
toString : any -> String
toString output =
    __asm callexternal 00 coreDebugToString 01

`

const intCode = `
__externalfn coreIntToFixed 1
toFixed : Int -> Fixed
toFixed a =
    __asm callexternal 00 coreIntToFixed 01


__externalfn coreFixedToInt 1
round : Fixed -> Int
round a =
    __asm callexternal 00 coreFixedToInt 01

`


const typeIdCode = `
`

const globalCode = `


type Maybe a =
    Nothing
    | Just a
`

func compileToGlobal(rootModule *decorated.Module, globalModule *decorated.Module, name string, code string) (*decorated.Module, decshared.DecoratedError) {
	const verbose = true
	const enforceStyle = true
	const errorAsWarning = false
	nameTypeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken(name, token.PositionLength{}, 0))
	newModule, err := InternalCompileToModule(nil, []*decorated.Module{globalModule, rootModule},
		nil, dectype.MakeArtifactFullyQualifiedModuleName([]*ast.TypeIdentifier{nameTypeIdentifier}),
		name, strings.TrimSpace(code), enforceStyle, verbose, errorAsWarning)
	if err != nil {
		return nil, err
	}
	ExposeEverythingInModule(newModule)
	return newModule, nil
}

func addCores(m *decorated.Module, globalModule *decorated.Module) ([]*decorated.Module, decshared.DecoratedError) {
	var importModules []*decorated.Module
	listModule, listModuleErr := compileToGlobal(m, globalModule, "List", listContent)
	if listModuleErr != nil {
		return nil, listModuleErr
	}
	importModules = append(importModules, listModule)
	mathModule, mathModuleErr := compileToGlobal(m, globalModule, "Math", mathCode)
	if mathModuleErr != nil {
		return nil, mathModuleErr
	}
	importModules = append(importModules, mathModule)

	blobModule, blobModuleErr := compileToGlobal(m, globalModule, "Blob", blobCode)
	if blobModuleErr != nil {
		return nil, blobModuleErr
	}
	importModules = append(importModules, blobModule)

	intModule, intModuleErr := compileToGlobal(m, globalModule, "Int", intCode)
	if intModuleErr != nil {
		return nil, intModuleErr
	}
	importModules = append(importModules, intModule)

	typeId, typeIdErr := compileToGlobal(m, globalModule, "TypeRef", typeIdCode)
	if typeIdErr != nil {
		return nil, typeIdErr
	}
	importModules = append(importModules, typeId)

	arrayModule, arrayModuleErr := compileToGlobal(m, globalModule, "Array", arrayCode)
	if arrayModuleErr != nil {
		return nil, arrayModuleErr
	}
	importModules = append(importModules, arrayModule)

	maybeModule, maybeModuleErr := compileToGlobal(m, globalModule, "Maybe", maybeCode)
	if maybeModuleErr != nil {
		return nil, maybeModuleErr
	}
	importModules = append(importModules, maybeModule)

	debugModule, arrayModuleErr := compileToGlobal(m, globalModule, "Debug", debugCode)
	if arrayModuleErr != nil {
		return nil, arrayModuleErr
	}
	importModules = append(importModules, debugModule)

	return importModules, nil
}

func createTypeIdentifier(name string) *ast.TypeIdentifier {
	symbol := token.NewTypeSymbolToken(name, token.PositionLength{}, 0)

	return ast.NewTypeIdentifier(symbol)
}

func CreateDefaultRootModule(includeCores bool) ([]*decorated.Module, []*decorated.Module, decshared.DecoratedError) {
	var importModules []*decorated.Module
	var copyModules []*decorated.Module
	nameTypeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("root-module", token.PositionLength{}, 0))
	m := decorated.NewModule(dectype.MakeArtifactFullyQualifiedModuleName([]*ast.TypeIdentifier{nameTypeIdentifier}))
	r := m.TypeRepo()
	integerType := dectype.NewPrimitiveType(createTypeIdentifier("Int"), nil)
	fixedType := dectype.NewPrimitiveType(createTypeIdentifier("Fixed"), nil)
	resourceNameType := dectype.NewPrimitiveType(createTypeIdentifier("ResourceName"), nil)
	stringType := dectype.NewPrimitiveType(createTypeIdentifier("String"), nil)
	boolType := dectype.NewPrimitiveType(createTypeIdentifier("Bool"), nil)
	blobType := dectype.NewPrimitiveType(createTypeIdentifier("Blob"), nil)

	r.DeclareType(integerType)
	r.DeclareType(fixedType)
	r.DeclareType(resourceNameType)
	r.DeclareType(stringType)
	r.DeclareType(boolType)
	r.DeclareType(blobType)

	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.PositionLength{}, 0))
	arrayIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("Array", token.PositionLength{}, 0))

	localTypeVariable := ast.NewVariableIdentifier(token.NewVariableSymbolToken("a", token.PositionLength{}, 0))
	typeParameter := ast.NewTypeParameter(localTypeVariable)
	localType := dectype.NewLocalType(typeParameter)
	listType := dectype.NewPrimitiveType(listIdentifier, []dtype.Type{localType})

	r.DeclareType(listType)
	r.DeclareAlias(listIdentifier, listType, nil)

	arrayType := dectype.NewPrimitiveType(arrayIdentifier, []dtype.Type{localType})
	r.DeclareType(arrayType)
	r.DeclareAlias(arrayIdentifier, arrayType, nil)

	typeRefIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("TypeRef", token.PositionLength{}, 0))
	typeRefType := dectype.NewPrimitiveType(arrayIdentifier, []dtype.Type{localType})
	r.DeclareType(typeRefType)
	r.DeclareAlias(typeRefIdentifier, typeRefType, nil)

	const verbose = true

	const enforceStyle = true

	defaultImportName := ast.NewTypeIdentifier(token.NewTypeSymbolToken("DefaultImport", token.PositionLength{}, 0))

	globalModule, globalModuleErr := InternalCompileToModule(nil, nil, nil,
		dectype.MakeArtifactFullyQualifiedModuleName([]*ast.TypeIdentifier{defaultImportName}), "",
		strings.TrimSpace(globalCode), enforceStyle, verbose, false)
	if globalModuleErr != nil {
		return nil, nil, globalModuleErr
	}

	if includeCores {
		var importModulesErr decshared.DecoratedError
		importModules, importModulesErr = addCores(m, globalModule)
		if importModulesErr != nil {
			fmt.Printf("ERROR:%v\n", importModulesErr)
			return nil, nil, importModulesErr
		}
	}

	CopyModuleToModule(m, globalModule)
	ExposeEverythingInModule(m)
	copyModules = append(copyModules, m)

	return copyModules, importModules, nil
}
