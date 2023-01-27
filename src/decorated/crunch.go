/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"fmt"
	"github.com/swamp/compiler/src/semantic"
	"log"
	"reflect"
	"strings"

	"github.com/swamp/compiler/src/ast"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/parser"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
	"github.com/swamp/compiler/src/verbosity"
)

type NoImportModuleRepository struct{}

func (*NoImportModuleRepository) FetchModuleInPackage(moduleType decorated.ModuleType, moduleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
	return nil, decorated.NewInternalError(fmt.Errorf("this is a no import module. Imports are not allowed"))
}

func CompileToModuleOnceForTest(code string, useCores bool, errorsAsWarnings bool) (*decorated.Module, decshared.DecoratedError) {
	rootModule, rootModuleErr := CreateDefaultRootModule(useCores)
	if parser.IsCompileError(rootModuleErr) {
		return nil, rootModuleErr
	}

	if rootModule == nil {
		panic("not allowed to be nil since it wasnt a compileError")
	}
	var totalErr decshared.DecoratedError

	totalErr = decorated.AppendError(totalErr, rootModuleErr)

	importRepository := &NoImportModuleRepository{}

	const verbose = verbosity.None

	const enforceStyle = true

	absoluteFilename := "fortest.swamp"
	module, compileErr := InternalCompileToModule(decorated.ModuleTypeNormal, importRepository, rootModule, dectype.MakeArtifactFullyQualifiedModuleName(nil), absoluteFilename, code,
		enforceStyle, verbose, errorsAsWarnings)

	totalErr = decorated.AppendError(totalErr, compileErr)

	return module, totalErr

}

func InternalCompileToProgram(absoluteFilename string, code string, enforceStyle bool, verbose verbosity.Verbosity) (*tokenize.Tokenizer, *ast.SourceFile, decshared.DecoratedError) {
	ioReader := strings.NewReader(code)

	var errors decshared.DecoratedError

	runeReader, runeReaderErr := runestream.NewRuneReader(ioReader, absoluteFilename)
	if runeReaderErr != nil {
		return nil, nil, decorated.NewInternalError(runeReaderErr)
	}

	tokenizer, tokenizerErr := tokenize.NewTokenizer(runeReader, enforceStyle)
	if parser.IsCompileError(tokenizerErr) {
		return tokenizer, nil, tokenizerErr
	}
	errors = decorated.AppendError(errors, tokenizerErr)

	p := parser.NewParser(tokenizer, enforceStyle)
	program, programErr := p.Parse()
	errors = decorated.AppendError(errors, programErr)
	if parser.IsCompileError(programErr) {
		return tokenizer, program, errors
	}

	if program == nil {
		panic("program must be valid since it was not a compile error")
	}

	parserErrors := p.Errors()
	errors = decorated.AppendError(errors, parserErrors)
	if parser.IsCompileError(parserErrors) {
		return tokenizer, nil, errors
	}

	program.SetNodes(p.Nodes())

	return tokenizer, program, errors
}

func isInternalType(typeToCheck dtype.Type) bool {
	_, isPrimitive := typeToCheck.(*dectype.PrimitiveAtom)

	return isPrimitive
}

func checkUnusedImports(module *decorated.Module) decshared.DecoratedError {
	var errors decshared.DecoratedError
	for _, definition := range module.ImportedModules().AllInOrderModules() {
		if !definition.WasReferenced() && !definition.ReferencedModule().IsInternal() {
			warning := decorated.NewUnusedImportWarning(definition, "importedModules.All()")
			errors = decorated.AppendError(errors, warning)
		}
	}

	for _, definition := range module.ImportedDefinitions().ReferencedDefinitions() {
		if !definition.WasReferenced() && !definition.IsInternal() {
			warning := decorated.NewUnusedImportWarning(definition.CreatedBy(), "referenced definitions")
			errors = decorated.AppendError(errors, warning)
		}
	}
	for _, importedType := range module.ImportedTypes().AllInOrderTypes() {
		if !importedType.WasReferenced() && importedType.CreatedByModuleImport() != nil && !isInternalType(importedType.ReferencedType()) {
			warning := decorated.NewUnusedImportWarning(importedType.CreatedByModuleImport(), fmt.Sprintf("imported types '%v'", importedType.String()))
			errors = decorated.AppendError(errors, warning)
		}
	}

	return errors
}

func InternalCompileToModule(moduleType decorated.ModuleType, moduleRepository ModuleRepository, rootModule *decorated.Module,
	moduleName dectype.ArtifactFullyQualifiedModuleName, absoluteFilename string, code string,
	enforceStyle bool, verbose verbosity.Verbosity, errorAsWarning bool) (*decorated.Module, decshared.DecoratedError) {
	var errors decshared.DecoratedError

	tokenizer, program, programErr := InternalCompileToProgram(absoluteFilename, code, enforceStyle, verbose)
	errors = decorated.AppendError(errors, programErr)
	if parser.IsCompileError(programErr) {
		return nil, programErr
	}

	module := decorated.NewModule(moduleType, moduleName, tokenizer.Document())

	// relativeModuleName := dectype.MakePackageRelativeModuleName(importModule.FullyQualifiedModuleName().Path())
	fakeModuleReference := decorated.NewModuleReference(rootModule.FullyQualifiedModuleName().Path(), rootModule)
	const exposeAllImports = true
	sourceFileReference := token.SourceFileReference{
		Range:    token.Range{},
		Document: rootModule.Document(),
	}

	keyword := token.NewKeyword("", 0, sourceFileReference)
	i := ast.NewImport(keyword, nil, nil, fakeModuleReference.AstModuleReference(), nil, nil, nil, true, nil)
	fakeImportStatement := decorated.NewImport(i, fakeModuleReference, fakeModuleReference, exposeAllImports)

	importErr := ImportModuleToModule(module, fakeImportStatement)
	if parser.IsCompileErr(importErr) {
		return nil, decorated.NewInternalError(importErr)
	}
	if importErr != nil {
		errors = decorated.AppendError(errors, decorated.NewInternalError(importErr))
	}

	for _, importedSubModule := range rootModule.ImportedModules().AllInOrderModules() {
		fakeModuleReference := decorated.NewModuleReference(rootModule.FullyQualifiedModuleName().Path(), importedSubModule.ReferencedModule())
		i := ast.NewImport(keyword, nil, nil, fakeModuleReference.AstModuleReference(), nil, nil, nil, true, nil)
		fakeImportStatement := decorated.NewImport(i, fakeModuleReference, fakeModuleReference, exposeAllImports)
		module.ImportedModules().ImportModule(importedSubModule.ModuleName(), importedSubModule.ReferencedModule(), fakeImportStatement)
	}
	importedModule := module.ImportedModules().ImportModule(rootModule.FullyQualifiedModuleName().Path(), rootModule, fakeImportStatement)

	typeLookup := decorated.NewTypeLookup(module.ImportedModules(), module.LocalTypes(), module.ImportedTypes())
	createAndLookup := decorated.NewTypeCreateAndLookup(typeLookup, module.LocalTypes(), dectype.NewTypeParameterContext())

	converter := NewDecorator(moduleRepository, module, importedModule, createAndLookup)

	rootStatementHandler := decorator.NewRootStatementHandler(converter, createAndLookup, moduleType, "compiletomodule")

	rootNodes, generateErr := rootStatementHandler.HandleStatements(program)
	errors = decorated.AppendError(errors, generateErr)
	if parser.IsCompileErr(generateErr) {
		return nil, generateErr
	}
	errors = decorated.AppendError(errors, converter.Errors())

	//importErrors := checkUnusedImports(module)
	//errors = decorated.AppendError(errors, importErrors)

	// log.Printf("before EXPOSING LOCAL TYPES:%v\n", module.ExposedTypes().DebugString())
	module.ExposedTypes().AddTypesFromModule(module.LocalTypes().AllInOrderTypes(), module)
	// log.Printf("AFTER EXPOSING LOCAL TYPES:%v\n", module.ExposedTypes().DebugString())
	module.ExposedDefinitions().AddDefinitions(module.LocalDefinitions().Definitions())
	module.SetProgram(program)

	var rootNodesConverted []decorated.Node
	for _, rootNode := range rootNodes {
		converted, couldConvert := rootNode.(decorated.Node)
		if !couldConvert {
			panic(fmt.Sprintf("can not convert %T", rootNode))
		}
		if converted == nil || reflect.ValueOf(converted).IsNil() {
			panic("can not be nil")
		}
		rootNodesConverted = append(rootNodesConverted, converted)
	}
	module.SetRootNodes(rootNodesConverted)

	log.Printf("expandedNodes : %d", len(module.ExpandedNodes()))

	_, semanticErr := semantic.GenerateTokensEncodedValues(module.ExpandedNodes(), module.Document())
	if semanticErr != nil {
		panic(semanticErr)
		return nil, decorated.NewInternalError(semanticErr)
	}

	return module, errors
}

func ImportModuleToModule(target *decorated.Module, statement *decorated.ImportStatement) error {
	if target == nil {
		panic("no target")
	}

	source := statement.Module()
	if source == nil {
		panic("no source")
	}

	// sourceMountedModuleName := source.FullyQualifiedModuleName()

	exposedTypes := source.ExposedTypes().AllInOrderTypes()
	exposedDefinitions := source.ExposedDefinitions().ReferencedDefinitions()

	importedModule := target.ImportedModules().ImportModule(statement.ImportAsName().AstModuleReference(), source, statement)

	if statement.ExposeAll() {
		target.ImportedTypes().AddTypes(exposedTypes, importedModule)
		for _, exposedDefinition := range exposedDefinitions {
			importedModule := decorated.NewImportedModule(target,
				statement)
			importedDefinition := decorated.NewImportedDefinition(importedModule, exposedDefinition.Identifier(), exposedDefinition)
			target.ImportedDefinitions().AddDefinition(exposedDefinition.Identifier(), importedDefinition)
		}

		// HACK
		target.ExposedTypes().AddTypes(exposedTypes, importedModule)
		target.ExposedDefinitions().AddDefinitions(exposedDefinitions)
	}

	return nil
}

func CopyModuleToModule(target *decorated.Module, source *decorated.Module) decshared.DecoratedError {
	if err := target.LocalTypes().CopyTypes(source.LocalTypes().AllInOrderTypes()); err != nil {
		return err
	}

	if err := target.LocalDefinitions().CopyFrom(source.LocalDefinitions()); err != nil {
		return decorated.NewInternalError(err)
	}

	return nil
}

func ExposeEverythingInModule(target *decorated.Module) decshared.DecoratedError {
	target.ExposedTypes().AddTypesFromModule(target.LocalTypes().AllInOrderTypes(), target)
	if err := target.ExposedDefinitions().AddDefinitions(target.LocalDefinitions().Definitions()); err != nil {
		return decorated.NewInternalError(err)
	}
	//	target.ExposedDefinitions().DebugOutput("expose everything in module")
	return nil
}
