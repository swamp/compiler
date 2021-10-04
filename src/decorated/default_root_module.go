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



map x y =
    __asm callexternal 00 coreListMap 01 02


map2 : (a -> b -> c) -> List a -> List b -> List c
map2 x y z =
    __asm callexternal 00 coreListMap2 01 02 03





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


isEmpty lst =
    __asm callexternal 00 coreListIsEmpty 01


range : Int -> Int -> List Int
range start end =
    __asm callexternal 00 coreListRange 01 02


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

const listCode = `
__externalvarfn head : List a -> Maybe a
__externalvarfn map : (a -> b) -> List a -> List b
__externalvarfn isEmpty : List a -> Bool
__externalvarfn length : List a -> Int
__externalvarexfn foldl : (a -> b -> b) -> b -> List a -> b
__externalvarexfn foldlstop : (a -> b -> Maybe b) -> b -> List a -> b
__externalvarexfn filterMap : (a -> Maybe b) -> List a -> List b
__externalvarexfn indexedMap : (Int -> a -> b) -> List a -> List b
`

const mathCode = `
__externalfn remainderBy : Int -> Int -> Int
__externalfn sin : Fixed -> Fixed
__externalfn cos : Fixed -> Fixed
__externalfn rnd : Int -> Int -> Int
__externalfn atan2 : Int -> Int -> Fixed
__externalfn mid : Int -> Int -> Int
__externalfn abs : Int -> Int
__externalfn sign : Int -> Int
__externalfn clamp : Int -> Int -> Int -> Int
__externalfn lerp : Fixed -> Int -> Int -> Int
__externalfn metronome : Int -> Int -> Int -> Int -> Bool
__externalfn drunk : Int -> Int -> Int -> Int
__externalfn mod : Int -> Int -> Int
`

const blobCode = `
__externalfn map : (Int -> Int) -> Blob -> Blob
__externalfn indexedMap : (Int -> Int -> Int) -> Blob -> Blob
__externalvarfn filterIndexedMap : (Int -> Int -> Maybe a) -> Blob -> List a
__externalfn toString2d : { width : Int, height : Int } -> Blob -> String
__externalfn get2d : { x : Int, y : Int } -> { width : Int, height : Int } -> Blob -> Maybe Int
__externalfn fromArray : Array Int -> Blob
-- __externalfn isEmpty : Blob -> Bool
-- __externalfn fromList : List Int -> Blob
-- __externalfn toList : Blob -> List Int
-- __externalfn length : Blob -> Int
-- __externalfn get : Int -> Blob -> Maybe Int
-- __externalfn grab : Int -> Blob -> Int
-- __externalfn grab2d : { x : Int, y : Int } -> { width : Int, height : Int } -> Blob -> Int
-- __externalfn set : Int -> Int -> Blob -> Blob
-- __externalfn set2d : { x : Int, y : Int } -> { width : Int, height : Int } -> Int -> Blob -> Blob
`

const arrayCode = `
__externalvarfn fromList : List a -> Array a
__externalvarfn toList : Array a -> List a
__externalvarfn grab : Int -> Array a -> a
__externalvarfn length : Array a -> Int
-- __externalvarfn slice : Int -> Int -> Array a -> Array a
-- __externalvarfn repeat : Int -> a -> Array a
-- __externalvarfn get : Int -> Array a -> Maybe a
-- __externalvarfn set : Int -> a -> Array a -> Array a
`

const maybeCode = `
__externalvarexfn withDefault : a -> Maybe a -> a
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
__externalfn log : String -> String
__externalvarfn logAny : Any -> String
__externalvarfn toString : Any -> String
`

const intCode = `
__externalfn toFixed : Int -> Fixed
__externalfn round : Fixed -> Int
`

const charCode = `
__externalfn ord : Char -> Int
__externalfn toCode : Char -> Int
__externalfn fromCode : Int -> Char
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

	newModule, err := InternalCompileToModule(decorated.ModuleTypeNormal, nil, []*decorated.Module{globalModule, targetModule},
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
	mathModule, mathModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Math", mathCode)
	if mathModuleErr != nil {
		return nil, mathModuleErr
	}
	importModules = append(importModules, mathModule)
	listModule, listModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "List", listCode)
	if listModuleErr != nil {
		return nil, listModuleErr
	}
	importModules = append(importModules, listModule)

	maybeModule, maybeModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Maybe", maybeCode)
	if maybeModuleErr != nil {
		return nil, maybeModuleErr
	}
	importModules = append(importModules, maybeModule)

	intModule, intModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Int", intCode)
	if intModuleErr != nil {
		return nil, intModuleErr
	}
	importModules = append(importModules, intModule)

	debugModule, debugModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Debug", debugCode)
	if debugModuleErr != nil {
		return nil, debugModuleErr
	}
	importModules = append(importModules, debugModule)

	arrayModule, arrayModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Array", arrayCode)
	if arrayModuleErr != nil {
		return nil, arrayModuleErr
	}
	importModules = append(importModules, arrayModule)

	blobModule, blobModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Blob", blobCode)
	if blobModuleErr != nil {
		return nil, blobModuleErr
	}
	importModules = append(importModules, blobModule)

	charModule, charModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Char", charCode)
	if charModuleErr != nil {
		return nil, charModuleErr
	}
	importModules = append(importModules, charModule)

	/*
		tupleModule, tupleModuleErr := compileToGlobal(targetModule, globalPrimitiveModule, "Tuple", tupleCode)
		if tupleModuleErr != nil {
			return nil, listModuleErr
		}
		importModules = append(importModules, tupleModule)

		typeId, typeIdErr := compileToGlobal(targetModule, globalPrimitiveModule, "TypeRef", typeIdCode)
		if typeIdErr != nil {
			return nil, typeIdErr
		}
		importModules = append(importModules, typeId)
	*/

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
	rootPrimitiveModule := decorated.NewModule(decorated.ModuleTypeNormal, dectype.MakeArtifactFullyQualifiedModuleName(ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(nameTypeIdentifier)})), nil)
	rootPrimitiveModule.MarkAsInternal()
	globalModuleTypes := rootPrimitiveModule.LocalTypes()

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

	globalModule, globalModuleErr := InternalCompileToModule(decorated.ModuleTypeNormal, nil, nil, nil,
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
