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
	"github.com/swamp/compiler/src/verbosity"
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


__externalfn coreListFilterMap2 3
filterMap2 : (a -> b -> Maybe c) -> List a -> List b -> List c
filterMap2 pred a b =
    __asm callexternal 00 coreListFilterMap2 01 02 03


filter : (a -> Bool) -> List a -> List a
filter pred coll =
    __asm callexternal 00 coreListFilter 01 02


__externalfn coreListFilter2 3
filter2 : (a -> b -> Bool) -> List a -> List b -> List b
filter2 pred a b =
    __asm callexternal 00 coreListFilter2 01 02 03


remove : (a -> Bool) -> List a -> List a
remove pred coll =
    __asm callexternal 00 coreListRemove 01 02


__externalfn coreListRemove2 3
remove2 : (a -> b -> Bool) -> List a -> List b -> List b
remove2 pred a b =
    __asm callexternal 00 coreListRemove2 01 02 03


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


__externalfn coreListUnzip 1
unzip : List (a, b) -> (List a, List b)
unzip lst =
    __asm callexternal 00 coreListUnzip 01


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


__externalfn coreBlobMap 2
map : (Int -> Int) -> Blob -> Blob
map fn blob =
    __asm callexternal 00 coreBlobMap 01 02


__externalfn coreBlobIndexedMap 2
indexedMap : (Int -> Int -> Int) -> Blob -> Blob
indexedMap fn blob =
    __asm callexternal 00 coreBlobIndexedMap 01 02


__externalfn coreBlobLength 1
length : Blob -> Int
length blob =
    __asm callexternal 00 coreBlobLength 01


__externalfn coreBlobGet 2
get : Int -> Blob -> Maybe Int
get index blob =
    __asm callexternal 00 coreBlobGet 01 02


__externalfn coreBlobGrab 2
grab : Int -> Blob -> Int
grab index blob =
    __asm callexternal 00 coreBlobGrab 01 02


__externalfn coreBlobGet2d 3
get2d : { x : Int, y : Int } -> { width : Int, height : Int } -> Blob -> Maybe Int
get2d position size blob =
    __asm callexternal 00 coreBlobGet2d 01 02 03


__externalfn coreBlobGrab2d 3
grab2d : { x : Int, y : Int } -> { width : Int, height : Int } -> Blob -> Int
grab2d position size blob =
    __asm callexternal 00 coreBlobGrab2d 01 02 03


__externalfn coreBlobSet 3
set : Int -> Int -> Blob -> Blob
set index item blob =
    __asm callexternal 00 coreBlobSet 01 02 03


__externalfn coreBlobSet2d 4
set2d : { x : Int, y : Int } -> { width : Int, height : Int } -> Int -> Blob -> Blob
set2d position size item blob =
    __asm callexternal 00 coreBlobSet2d 01 02 03 04


__externalfn coreBlobToString2d 2
toString2d : { width : Int, height : Int } -> Blob -> String
toString2d size blob =
    __asm callexternal 00 coreBlobToString2d 01 02

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

const tupleCode = `
__externalfn coreTupleFirst 1
first : (a, b) -> a
first tuple =
    __asm callexternal 00 coreTupleFirst 01


__externalfn coreTupleSecond 1
second : (a, b) -> b
second tuple =
    __asm callexternal 00 coreTupleSecond 01


__externalfn coreTupleThird 1
third : (a, b, c) -> c
third tuple =
    __asm callexternal 00 coreTupleThird 01


__externalfn coreTupleFourth 1
third : (a, b, c, d) -> d
third tuple =
    __asm callexternal 00 coreTupleFourth 01
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

const charCode = `
__externalfn coreCharOrd 1
ord : Char -> Int
ord a =
    __asm callexternal 00 coreCharOrd 01


__externalfn coreCharToCode 1
toCode : Char -> Int
toCode a =
    __asm callexternal 00 coreCharToCode 01


__externalfn coreCharFromCode 1
fromCode : Int -> Char
fromCode a =
    __asm callexternal 00 coreCharFromCode 01

`

const typeIdCode = `
`

const globalCode = `

type Maybe a =
    Nothing
    | Just a
`

func compileToGlobal(targetModule *decorated.Module, globalModule *decorated.Module, name string, code string) (*decorated.Module, decshared.DecoratedError) {
	const verbose = verbosity.Low
	const enforceStyle = true
	const errorAsWarning = false

	nameTypeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken(name, token.SourceFileReference{}, 0))

	newModule, err := InternalCompileToModule(nil, []*decorated.Module{globalModule, targetModule},
		nil, dectype.MakeArtifactFullyQualifiedModuleName(ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(nameTypeIdentifier)})),
		name, strings.TrimSpace(code), enforceStyle, verbose, errorAsWarning)
	if err != nil {
		return nil, err
	}

	newModule.MarkAsInternal()

	// moduleReference := ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(nameTypeIdentifier)})

	ExposeEverythingInModule(newModule)

	return newModule, nil
}

func addCores(targetModule *decorated.Module, globalPrimitiveModule *decorated.Module) ([]*decorated.Module, decshared.DecoratedError) {
	var importModules []*decorated.Module
	listModule, listModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "List", listContent)
	if listModuleErr != nil {
		return nil, listModuleErr
	}
	importModules = append(importModules, listModule)

	tupleModule, tupleModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Tuple", tupleCode)
	if tupleModuleErr != nil {
		return nil, listModuleErr
	}
	importModules = append(importModules, tupleModule)

	mathModule, mathModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Math", mathCode)
	if mathModuleErr != nil {
		return nil, mathModuleErr
	}
	importModules = append(importModules, mathModule)

	blobModule, blobModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Blob", blobCode)
	if blobModuleErr != nil {
		return nil, blobModuleErr
	}
	importModules = append(importModules, blobModule)

	intModule, intModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Int", intCode)
	if intModuleErr != nil {
		return nil, intModuleErr
	}
	importModules = append(importModules, intModule)

	charModule, charModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Char", charCode)
	if charModuleErr != nil {
		return nil, charModuleErr
	}
	importModules = append(importModules, charModule)

	typeId, typeIdErr := compileToGlobal(targetModule, globalPrimitiveModule, "TypeRef", typeIdCode)
	if typeIdErr != nil {
		return nil, typeIdErr
	}
	importModules = append(importModules, typeId)

	arrayModule, arrayModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Array", arrayCode)
	if arrayModuleErr != nil {
		return nil, arrayModuleErr
	}
	importModules = append(importModules, arrayModule)

	maybeModule, maybeModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Maybe", maybeCode)
	if maybeModuleErr != nil {
		return nil, maybeModuleErr
	}
	importModules = append(importModules, maybeModule)

	debugModule, arrayModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Debug", debugCode)
	if arrayModuleErr != nil {
		return nil, arrayModuleErr
	}
	importModules = append(importModules, debugModule)

	return importModules, nil
}

func createTypeIdentifier(name string) *ast.TypeIdentifier {
	symbol := token.NewTypeSymbolToken(name, token.SourceFileReference{}, 0)

	return ast.NewTypeIdentifier(symbol)
}

func addPrimitive(types *decorated.ModuleTypes, atom *dectype.PrimitiveAtom) {
	types.InternalAddPrimitive(atom.PrimitiveName(), atom)
}

func CreateDefaultRootModule(includeCores bool) ([]*decorated.Module, []*decorated.Module, decshared.DecoratedError) {
	var importModules []*decorated.Module
	var copyModules []*decorated.Module
	nameTypeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("root-module", token.SourceFileReference{}, 0))
	rootPrimitiveModule := decorated.NewModule(dectype.MakeArtifactFullyQualifiedModuleName(ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(nameTypeIdentifier)})), nil)
	rootPrimitiveModule.MarkAsInternal()
	globalModuleTypes := rootPrimitiveModule.TypeRepo()

	anyType := dectype.NewAnyType()
	integerType := dectype.NewPrimitiveType(createTypeIdentifier("Int"), nil)
	fixedType := dectype.NewPrimitiveType(createTypeIdentifier("Fixed"), nil)
	resourceNameType := dectype.NewPrimitiveType(createTypeIdentifier("ResourceName"), nil)
	stringType := dectype.NewPrimitiveType(createTypeIdentifier("String"), nil)
	charType := dectype.NewPrimitiveType(createTypeIdentifier("Char"), nil)
	boolType := dectype.NewPrimitiveType(createTypeIdentifier("Bool"), nil)
	blobType := dectype.NewPrimitiveType(createTypeIdentifier("Blob"), nil)
	unmanagedType := dectype.NewPrimitiveType(createTypeIdentifier("Unmanaged"), nil)

	addPrimitive(globalModuleTypes, anyType)
	addPrimitive(globalModuleTypes, integerType)
	addPrimitive(globalModuleTypes, fixedType)
	addPrimitive(globalModuleTypes, resourceNameType)
	addPrimitive(globalModuleTypes, stringType)
	addPrimitive(globalModuleTypes, charType)
	addPrimitive(globalModuleTypes, boolType)
	addPrimitive(globalModuleTypes, blobType)
	addPrimitive(globalModuleTypes, unmanagedType)

	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.SourceFileReference{}, 0))

	localTypeVariable := ast.NewVariableIdentifier(token.NewVariableSymbolToken("a", token.SourceFileReference{}, 0))
	typeParameter := ast.NewTypeParameter(localTypeVariable)
	localType := dectype.NewLocalType(typeParameter)
	listType := dectype.NewPrimitiveType(listIdentifier, []dtype.Type{localType})

	addPrimitive(globalModuleTypes, listType)

	arrayIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("Array", token.SourceFileReference{}, 0))
	arrayType := dectype.NewPrimitiveType(arrayIdentifier, []dtype.Type{localType})

	addPrimitive(globalModuleTypes, arrayType)

	typeRefIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("TypeRef", token.SourceFileReference{}, 0))
	typeRefType := dectype.NewPrimitiveType(typeRefIdentifier, []dtype.Type{localType})

	addPrimitive(globalModuleTypes, typeRefType)

	const verbose = verbosity.Mid

	const enforceStyle = true

	defaultImportName := ast.NewTypeIdentifier(token.NewTypeSymbolToken("DefaultImport", token.SourceFileReference{}, 0))

	globalModule, globalModuleErr := InternalCompileToModule(nil, nil, nil,
		dectype.MakeArtifactFullyQualifiedModuleName(ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(defaultImportName)})), "(internal root)",
		strings.TrimSpace(globalCode), enforceStyle, verbose, false)
	if globalModuleErr != nil {
		return nil, nil, globalModuleErr
	}
	globalModule.MarkAsInternal()

	if includeCores {
		var importModulesErr decshared.DecoratedError
		importModules, importModulesErr = addCores(rootPrimitiveModule, globalModule)
		if importModulesErr != nil {
			fmt.Printf("ERROR:%v\n", importModulesErr)
			return nil, nil, importModulesErr
		}
	}

	CopyModuleToModule(rootPrimitiveModule, globalModule)
	ExposeEverythingInModule(rootPrimitiveModule)
	copyModules = append(copyModules, rootPrimitiveModule)

	return copyModules, importModules, nil
}
