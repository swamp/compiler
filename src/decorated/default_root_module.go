/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"log"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/parser"
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
__externalvarexfn maybe : b -> (a -> b) -> Maybe a -> b
`

const tupleCode = `
__externalfn first : (a, b) -> a
__externalfn second : (a, b) -> b
__externalfn third : (a, b, c) -> c
__externalfn forth : (a, b, c, d) -> d
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

const stdCode = `
type Maybe a =
    Nothing
    | Just a


type Result value error =
    Ok value
    | Err error
`

func compileToModule(globalModule *decorated.Module, name string, code string) (*decorated.Module, decshared.DecoratedError) {
	const verbose = verbosity.Low
	const enforceStyle = true
	const errorAsWarning = false

	nameTypeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken(name, globalModule.FetchPositionLength(), 0))

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
	var allErr decshared.DecoratedError

	newModule, err := compileToModule(targetModule, name, code)
	if parser.IsCompileError(err) {
		return err
	}
	allErr = decorated.AppendError(allErr, err)

	reference := newModule.FullyQualifiedModuleName().ModuleName.Path()
	exposeAllImports := true

	fakeSourceFileReference := token.SourceFileReference{
		Range:    token.Range{},
		Document: targetModule.Document(),
	}
	keyword := token.NewKeyword("", 0, fakeSourceFileReference)
	i := ast.NewImport(keyword, nil, nil, reference, nil, nil, nil, true, nil)

	fakeImportStatement := decorated.NewImport(i, nil, nil, exposeAllImports)

	targetModule.ImportedModules().ImportModule(reference, newModule, fakeImportStatement)

	return allErr
}

func addCores(globalPrimitiveModule *decorated.Module) decshared.DecoratedError {
	var err decshared.DecoratedError

	maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Maybe", maybeCode)
	if parser.IsCompileError(maybeModuleErr) {
		return maybeModuleErr
	}
	err = decorated.AppendError(err, maybeModuleErr)

	mathModuleErr := compileAndAddToModule(globalPrimitiveModule, "Math", mathCode)
	if parser.IsCompileError(mathModuleErr) {
		return mathModuleErr
	}
	err = decorated.AppendError(err, mathModuleErr)

	listModuleErr := compileAndAddToModule(globalPrimitiveModule, "List", listCode)
	if parser.IsCompileError(listModuleErr) {
		return listModuleErr
	}
	err = decorated.AppendError(err, listModuleErr)

	intModuleErr := compileAndAddToModule(globalPrimitiveModule, "Int", intCode)
	if parser.IsCompileError(intModuleErr) {
		return intModuleErr
	}
	err = decorated.AppendError(err, intModuleErr)

	debugModuleErr := compileAndAddToModule(globalPrimitiveModule, "Debug", debugCode)
	if parser.IsCompileError(debugModuleErr) {
		return debugModuleErr
	}
	err = decorated.AppendError(err, debugModuleErr)

	arrayModuleErr := compileAndAddToModule(globalPrimitiveModule, "Array", arrayCode)
	if parser.IsCompileError(arrayModuleErr) {
		return arrayModuleErr
	}
	err = decorated.AppendError(err, arrayModuleErr)

	blobModuleErr := compileAndAddToModule(globalPrimitiveModule, "Blob", blobCode)
	if parser.IsCompileError(blobModuleErr) {
		return blobModuleErr
	}
	err = decorated.AppendError(err, blobModuleErr)

	charModuleErr := compileAndAddToModule(globalPrimitiveModule, "Char", charCode)
	if parser.IsCompileError(charModuleErr) {
		return charModuleErr
	}
	err = decorated.AppendError(err, charModuleErr)

	stringModuleErr := compileAndAddToModule(globalPrimitiveModule, "String", stringCode)
	if parser.IsCompileError(stringModuleErr) {
		return stringModuleErr
	}
	err = decorated.AppendError(err, stringModuleErr)

	/*
		maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "Tuple", tupleCode); maybeModuleErr != nil {
			return maybeModuleErr
		}

		if maybeModuleErr := compileAndAddToModule(globalPrimitiveModule, "TypeRef", typeIdCode); maybeModuleErr != nil {
			return maybeModuleErr
		}
	*/

	return err
}

func createTypeIdentifier(name string) *ast.TypeIdentifier {
	symbol := token.NewTypeSymbolToken(name, token.NewInternalSourceFileReference(), 0)

	return ast.NewTypeIdentifier(symbol)
}

func addPrimitive(types *decorated.ModuleTypes, atom *dectype.PrimitiveAtom) {
	types.InternalAddPrimitive(atom.PrimitiveName(), atom)
}

func kickstartPrimitives() *decorated.Module {
	newSourceFileUri := token.MakeDocumentURI("file://internal/")
	doc := &token.SourceFileDocument{Uri: newSourceFileUri}
	sourceFileReference := token.SourceFileReference{Document: doc, Range: token.Range{}}

	nameTypeIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("", sourceFileReference, 0))
	sourceFileDocument := &token.SourceFileDocument{
		Uri: newSourceFileUri}
	rootPrimitiveModule := decorated.NewModule(decorated.ModuleTypeNormal, dectype.MakeArtifactFullyQualifiedModuleName(ast.NewModuleReference([]*ast.ModuleNamePart{ast.NewModuleNamePart(nameTypeIdentifier)})), sourceFileDocument)
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

	listIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("List", token.NewInternalSourceFileReference(), 0))

	localTypeVariable := ast.NewVariableIdentifier(token.NewVariableSymbolToken("a", token.NewInternalSourceFileReference(), 0))
	typeParameter := ast.NewTypeParameter(localTypeVariable)
	localType := dectype.NewLocalType(typeParameter)
	listType := dectype.NewPrimitiveType(listIdentifier, []dtype.Type{localType})

	addPrimitive(primitiveModuleLocalTypes, listType)

	arrayIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("Array", token.NewInternalSourceFileReference(), 0))
	arrayType := dectype.NewPrimitiveType(arrayIdentifier, []dtype.Type{localType})

	addPrimitive(primitiveModuleLocalTypes, arrayType)

	typeRefIdentifier := ast.NewTypeIdentifier(token.NewTypeSymbolToken("TypeRef", token.NewInternalSourceFileReference(), 0))
	typeRefType := dectype.NewPrimitiveType(typeRefIdentifier, []dtype.Type{localType})

	addPrimitive(primitiveModuleLocalTypes, typeRefType)

	ExposeEverythingInModule(rootPrimitiveModule)

	return rootPrimitiveModule
}

func CreateDefaultRootModule(includeCores bool) (*decorated.Module, decshared.DecoratedError) {
	primitiveModule := kickstartPrimitives()

	var err decshared.DecoratedError
	stdModule, stdModuleErr := compileToModule(primitiveModule, "", stdCode)
	if parser.IsCompileError(stdModuleErr) {
		return nil, stdModuleErr
	}
	err = decorated.AppendError(err, stdModuleErr)

	if err := primitiveModule.LocalTypes().CopyTypes(stdModule.LocalTypes().AllInOrderTypes()); err != nil {
		return nil, err
	}

	primitiveModule.LocalDefinitions().CopyFrom(stdModule.LocalDefinitions())

	ExposeEverythingInModule(primitiveModule)

	if includeCores {
		importModulesErr := addCores(primitiveModule)
		if parser.IsCompileError(importModulesErr) {
			log.Printf("ERROR:%v\n", importModulesErr)
			return nil, importModulesErr
		}
		err = decorated.AppendError(err, importModulesErr)
	}

	// log.Printf("rootPrimitiveModule is finally %v\n", primitiveModule)

	return primitiveModule, err
}
