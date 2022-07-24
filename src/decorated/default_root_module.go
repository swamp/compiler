/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"github.com/swamp/compiler/src/parser"
	"log"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/verbosity"
)

const listCode = `
__externalvarfn head : List a -> Maybe a
__externalvarfn map : (a -> b) -> List a -> List b
__externalvarfn map2 : (a -> b -> c) -> List a -> List b -> List c
__externalvarfn concatMap : (a -> List b) -> List a -> List b
__externalvarfn isEmpty : List a -> Bool
__externalvarfn length : List a -> Int
__externalvarexfn foldl : (a -> b -> b) -> b -> List a -> b
__externalvarexfn foldlstop : (a -> b -> Maybe b) -> b -> List a -> b
__externalvarexfn reduce : (a -> a -> a) -> List a -> a
__externalvarexfn filterMap : (a -> Maybe b) -> List a -> List b
__externalvarfn indexedMap : (Int -> a -> b) -> List a -> List b
__externalvarfn find : (a -> Bool) -> List a -> Maybe a
__externalfn range : Int -> Int -> List Int
__externalfn range0 : Int -> List Int
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
__externalfn mapToBlob : (Int -> Int) -> Blob -> Blob
__externalfn indexedMapToBlob : (Int -> Int -> Int) -> Blob -> Blob
__externalfn indexedMapToBlob! : (Int -> Int -> Int) -> Blob -> Blob
__externalvarfn filterIndexedMap2d : ({ x : Int, y : Int } -> Int -> Maybe a) ->
    { width : Int, height : Int } -> Blob -> List a
__externalvarfn filterIndexedMap : (Int -> Int -> Maybe a) -> Blob -> List a
__externalfn toString2d : { width : Int, height : Int } -> Blob -> String
__externalfn get2d : { x : Int, y : Int } -> { width : Int, height : Int } -> Blob -> Maybe Int
__externalfn slice2d : { x : Int, y : Int } -> { width : Int, height : Int } ->
    { width : Int, height : Int } -> Blob -> Blob
__externalfn fill2d! : { x : Int, y : Int } -> { width : Int, height : Int } -> Int ->
    { width : Int, height : Int } -> Blob -> Blob
__externalfn copy2d! : { x : Int, y : Int } -> { width : Int, height : Int } ->
    { width : Int, height : Int } -> Blob -> Blob -> Blob
__externalfn drawWindow2d! : { x : Int, y : Int } -> { width : Int, height : Int } ->
    { width : Int, height : Int } -> Blob -> Blob
__externalfn member : Int -> Blob -> Bool
__externalfn any : (Int -> Bool) -> Blob -> Bool
__externalfn fromArray : Array Int -> Blob
__externalfn make : Int -> Blob
__externalvarfn map2d : ({ x : Int, y : Int } -> Int -> a) -> { width : Int, height : Int } -> Blob -> List a
__externalfn fromList : List Int -> Blob
-- __externalfn isEmpty : Blob -> Bool
-- __externalvarfn map : (Int -> a) -> Blob -> List a
-- __externalvarfn indexedMap : (Int -> Int -> a) -> Blob -> List a
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
__externalvarfn get : Int -> Array a -> Maybe a
-- __externalvarfn slice : Int -> Int -> Array a -> Array a
-- __externalvarfn repeat : Int -> a -> Array a
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
__externalfn log : Any -> String
__externalvarfn toString : Any -> String
__externalfn panic : Any -> Any

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

const stringCode = `
__externalfn fromInt : Int -> String
`

const typeIdCode = `
`

/*
TODO: Add this
type Result a b =
    Ok a
    | Err b
*/
const stdCode = `
type Maybe a =
    Nothing
    | Just a
`

func compileToModule(globalModule *decorated.Module, name string, code string) (*decorated.Module, decshared.DecoratedError) {
	const verbose = verbosity.Low
	const enforceStyle = true
	const errorAsWarning = false

	nameTypeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken(name, token.SourceFileReference{}, 0))

	var fullyQualifiedName dectype.ArtifactFullyQualifiedModuleName
	if name == "" {
		fullyQualifiedName = dectype.MakeArtifactFullyQualifiedModuleName(nil)
	} else {
		fullyQualifiedName = dectype.MakeArtifactFullyQualifiedModuleName(ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(nameTypeIdentifier)}))
	}

	newModule, err := InternalCompileToModule(decorated.ModuleTypeNormal, nil, globalModule,
		fullyQualifiedName,
		name+"_internal", strings.TrimSpace(code), enforceStyle, verbose, errorAsWarning)
	if parser.IsCompileError(err) {
		return nil, err
	}

	if newModule == nil {
		panic("newModule can not be nil, if it isnt a compile err")
	}

	newModule.MarkAsInternal()

	// moduleReference := ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(nameTypeIdentifier)})

	// newModule.DebugOutput("before EXPOSEEVERYTHINGINMODULE")
	/*
	   if exposeErr := ExposeEverythingInModule(newModule); exposeErr != nil {
	   		return nil, exposeErr
	   	}
	*/

	return newModule, err
}

func compileAndAddToModule(targetModule *decorated.Module, name string, code string) decshared.DecoratedError {
	newModule, err := compileToModule(targetModule, name, code)
	if err != nil {
		return err
	}

	reference := newModule.FullyQualifiedModuleName().ModuleName.Path()
	exposeAllImports := true

	keyword := token.NewKeyword("", 0, token.SourceFileReference{})
	i := ast.NewImport(keyword, nil, nil, reference, nil, nil, nil, true, nil)

	fakeImportStatement := decorated.NewImport(i, nil, nil, exposeAllImports)

	targetModule.ImportedModules().ImportModule(reference, newModule, fakeImportStatement)

	return nil
}

func addCores(globalPrimitiveModule *decorated.Module) decshared.DecoratedError {
	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Maybe", maybeCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Math", mathCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "List", listCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Int", intCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Debug", debugCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Array", arrayCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Blob", blobCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Char", charCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "String", stringCode); maybeModuleErr != nil {
		return maybeModuleErr
	}

	/*
		tupleModule, tupleModuleErr := compileToModule(stdModule, globalPrimitiveModule, "Tuple", tupleCode)
		if tupleModuleErr != nil {
			return nil, listModuleErr
		}
		importModules = append(importModules, tupleModule)

		typeId, typeIdErr := compileToModule(stdModule, globalPrimitiveModule, "TypeRef", typeIdCode)
		if typeIdErr != nil {
			return nil, typeIdErr
		}
		importModules = append(importModules, typeId)
	*/

	return nil
}

func createTypeIdentifier(name string) *ast.TypeIdentifier {
	symbol := token.NewTypeSymbolToken(name, token.SourceFileReference{}, 0)

	return ast.NewTypeIdentifier(symbol)
}

func addPrimitive(types *decorated.ModuleTypes, atom *dectype.PrimitiveAtom) {
	types.InternalAddPrimitive(atom.PrimitiveName(), atom)
}

func kickstartPrimitives() *decorated.Module {
	nameTypeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("", token.SourceFileReference{}, 0))
	rootPrimitiveModule := decorated.NewModule(decorated.ModuleTypeNormal, dectype.MakeArtifactFullyQualifiedModuleName(ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(nameTypeIdentifier)})), nil)
	rootPrimitiveModule.MarkAsInternal()
	primitiveModuleLocalTypes := rootPrimitiveModule.LocalTypes()

	anyType := dectype.NewAnyType()
	integerType := dectype.NewPrimitiveType(createTypeIdentifier("Int"), nil)
	fixedType := dectype.NewPrimitiveType(createTypeIdentifier("Fixed"), nil)
	resourceNameType := dectype.NewPrimitiveType(createTypeIdentifier("ResourceName"), nil)
	stringType := dectype.NewPrimitiveType(createTypeIdentifier("String"), nil)
	charType := dectype.NewPrimitiveType(createTypeIdentifier("Char"), nil)
	boolType := dectype.NewPrimitiveType(createTypeIdentifier("Bool"), nil)
	blobType := dectype.NewPrimitiveType(createTypeIdentifier("Blob"), nil)
	unmanagedType := dectype.NewPrimitiveType(createTypeIdentifier("Unmanaged"), nil)

	addPrimitive(primitiveModuleLocalTypes, anyType)
	addPrimitive(primitiveModuleLocalTypes, integerType)
	addPrimitive(primitiveModuleLocalTypes, fixedType)
	addPrimitive(primitiveModuleLocalTypes, resourceNameType)
	addPrimitive(primitiveModuleLocalTypes, stringType)
	addPrimitive(primitiveModuleLocalTypes, charType)
	addPrimitive(primitiveModuleLocalTypes, boolType)
	addPrimitive(primitiveModuleLocalTypes, blobType)
	addPrimitive(primitiveModuleLocalTypes, unmanagedType)

	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.SourceFileReference{}, 0))

	localTypeVariable := ast.NewVariableIdentifier(token.NewVariableSymbolToken("a", token.SourceFileReference{}, 0))
	typeParameter := ast.NewTypeParameter(localTypeVariable)
	localType := dectype.NewLocalType(typeParameter)
	listType := dectype.NewPrimitiveType(listIdentifier, []dtype.Type{localType})

	addPrimitive(primitiveModuleLocalTypes, listType)

	arrayIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("Array", token.SourceFileReference{}, 0))
	arrayType := dectype.NewPrimitiveType(arrayIdentifier, []dtype.Type{localType})

	addPrimitive(primitiveModuleLocalTypes, arrayType)

	typeRefIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("TypeRef", token.SourceFileReference{}, 0))
	typeRefType := dectype.NewPrimitiveType(typeRefIdentifier, []dtype.Type{localType})

	addPrimitive(primitiveModuleLocalTypes, typeRefType)

	ExposeEverythingInModule(rootPrimitiveModule)

	return rootPrimitiveModule
}

func CreateDefaultRootModule(includeCores bool) (*decorated.Module, decshared.DecoratedError) {
	primitiveModule := kickstartPrimitives()

	stdModule, stdModuleErr := compileToModule(primitiveModule, "", stdCode)
	if stdModuleErr != nil {
		return nil, stdModuleErr
	}
	if err := primitiveModule.LocalTypes().CopyTypes(stdModule.LocalTypes().AllInOrderTypes()); err != nil {
		return nil, err
	}
	primitiveModule.LocalDefinitions().CopyFrom(stdModule.LocalDefinitions())
	/*
		primitiveModule.LocalDefinitions().CopyFrom(stdModule.LocalDefinitions())
		importModules = append(importModules, stdModule)
	*/

	ExposeEverythingInModule(primitiveModule)

	if includeCores {
		var importModulesErr decshared.DecoratedError
		importModulesErr = addCores(primitiveModule)
		if importModulesErr != nil {
			log.Printf("ERROR:%v\n", importModulesErr)
			return nil, importModulesErr
		}
	}

	// log.Printf("rootPrimitiveModule is finally %v\n", primitiveModule)

	return primitiveModule, nil
}
