/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"fmt"
	"reflect"
	"strings"

	parerr "github.com/swamp/compiler/src/parser/errors"

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
	if rootModuleErr != nil {
		return nil, rootModuleErr
	}

	importRepository := &NoImportModuleRepository{}

	const verbose = verbosity.None

	const enforceStyle = true

	return InternalCompileToModule(decorated.ModuleTypeNormal, importRepository, rootModule, dectype.MakeArtifactFullyQualifiedModuleName(nil), "for test", code,
		enforceStyle, verbose, errorsAsWarnings)
}

func InternalCompileToProgram(absoluteFilename string, code string, enforceStyle bool, verbose verbosity.Verbosity) (*tokenize.Tokenizer, *ast.SourceFile, decshared.DecoratedError) {
	ioReader := strings.NewReader(code)

	var errors []decshared.DecoratedError

	runeReader, runeReaderErr := runestream.NewRuneReader(ioReader, absoluteFilename)
	if runeReaderErr != nil {
		return nil, nil, decorated.NewInternalError(runeReaderErr)
	}

	tokenizer, tokenizerErr := tokenize.NewTokenizer(runeReader, enforceStyle)
	if parser.IsCompileError(tokenizerErr) {
		return tokenizer, nil, tokenizerErr
	}
	if tokenizerErr != nil {
		errors = append(errors, tokenizerErr)
	}

	p := parser.NewParser(tokenizer, enforceStyle)
	program, programErr := p.Parse()
	if programErr != nil {
		errors = append(errors, programErr)
	}
	if parser.IsCompileError(programErr) {
		var returnErr parerr.ParseError
		if len(errors) > 0 {
			returnErr = decorated.NewMultiErrors(errors)
		}
		return tokenizer, program, returnErr
	}

	if program == nil {
		panic("program must be valid since it was not a compile error")
	}

	parserErrors := p.Errors()
	for _, parserErr := range parserErrors {
		errors = append(errors, parserErr)
		if parser.IsCompileError(parserErr) {
			return tokenizer, nil, parserErr
		}
	}

	program.SetNodes(p.Nodes())

	var returnErr parerr.ParseError
	if len(errors) > 0 {
		returnErr = decorated.NewMultiErrors(errors)
	}

	return tokenizer, program, returnErr
}

func isInternalType(typeToCheck dtype.Type) bool {
	_, isPrimitive := typeToCheck.(*dectype.PrimitiveAtom)

	return isPrimitive
}

func checkUnusedImports(module *decorated.Module) []decshared.DecoratedError {
	var errors []decshared.DecoratedError
	for _, definition := range module.ImportedModules().AllInOrderModules() {
		if !definition.WasReferenced() && !definition.ReferencedModule().IsInternal() {
			warning := decorated.NewUnusedImportWarning(definition, "importedModules.All()")
			errors = append(errors, warning)
		}
	}

	for _, definition := range module.ImportedDefinitions().ReferencedDefinitions() {
		if !definition.WasReferenced() && !definition.IsInternal() {
			warning := decorated.NewUnusedImportWarning(definition.CreatedBy(), "referenced definitions")
			errors = append(errors, warning)
		}
	}
	for _, importedType := range module.ImportedTypes().AllInOrderTypes() {
		if !importedType.WasReferenced() && importedType.CreatedByModuleImport() != nil && !isInternalType(importedType.ReferencedType()) {
			//warning := decorated.NewUnusedImportWarning(importedType.CreatedByModuleImport(), fmt.Sprintf("imported types '%v'", importedTypeName))
			//module.AddWarning(warning)
		}
	}

	return errors
}

func InternalCompileToModule(moduleType decorated.ModuleType, moduleRepository ModuleRepository, rootModule *decorated.Module,
	moduleName dectype.ArtifactFullyQualifiedModuleName, absoluteFilename string, code string,
	enforceStyle bool, verbose verbosity.Verbosity, errorAsWarning bool) (*decorated.Module, decshared.DecoratedError) {
	tokenizer, program, programErr := InternalCompileToProgram(absoluteFilename, code, enforceStyle, verbose)
	if parser.IsCompileError(programErr) {
		return nil, programErr
	}

	module := decorated.NewModule(moduleType, moduleName, tokenizer.Document())

	// relativeModuleName := dectype.MakePackageRelativeModuleName(importModule.FullyQualifiedModuleName().Path())
	fakeModuleReference := decorated.NewModuleReference(rootModule.FullyQualifiedModuleName().Path(), rootModule)
	const exposeAllImports = true
	keyword := token.NewKeyword("", 0, token.SourceFileReference{})
	i := ast.NewImport(keyword, nil, nil, fakeModuleReference.AstModuleReference(), nil, nil, nil, true, nil)
	fakeImportStatement := decorated.NewImport(i, fakeModuleReference, fakeModuleReference, exposeAllImports)
	importErr := ImportModuleToModule(module, fakeImportStatement)
	if importErr != nil {
		return nil, decorated.NewInternalError(importErr)
	}

	for _, importedSubModule := range rootModule.ImportedModules().AllInOrderModules() {
		fakeModuleReference := decorated.NewModuleReference(rootModule.FullyQualifiedModuleName().Path(), importedSubModule.ReferencedModule())
		i := ast.NewImport(keyword, nil, nil, fakeModuleReference.AstModuleReference(), nil, nil, nil, true, nil)
		fakeImportStatement := decorated.NewImport(i, fakeModuleReference, fakeModuleReference, exposeAllImports)
		module.ImportedModules().ImportModule(importedSubModule.ModuleName(), importedSubModule.ReferencedModule(), fakeImportStatement)
	}
	module.ImportedModules().ImportModule(rootModule.FullyQualifiedModuleName().Path(), rootModule, fakeImportStatement)

	typeLookup := decorated.NewTypeLookup(module.ImportedModules(), module.LocalTypes(), module.ImportedTypes())
	createAndLookup := decorated.NewTypeCreateAndLookup(typeLookup, module.LocalTypes())

	converter := NewDecorator(moduleRepository, module, createAndLookup)

	rootStatementHandler := decorator.NewRootStatementHandler(converter, createAndLookup, moduleType, "compiletomodule")
	var allErrors []decshared.DecoratedError
	rootNodes, generateErr := rootStatementHandler.HandleStatements(program)
	if generateErr != nil {
		allErrors = append(allErrors, generateErr)
	}
	allErrors = append(allErrors, converter.Errors()...)

	//importErrors := checkUnusedImports(module)
	//allErrors = append(allErrors, importErrors...)

	var returnErr decshared.DecoratedError

	if len(allErrors) > 0 {
		returnErr = decorated.NewMultiErrors(allErrors)
	}

	// log.Printf("before EXPOSING LOCAL TYPES:%v\n", module.ExposedTypes().DebugString())
	module.ExposedTypes().AddTypesFromModule(module.LocalTypes().AllInOrderTypes(), module)
	// log.Printf("AFTER EXPOSING LOCAL TYPES:%v\n", module.ExposedTypes().DebugString())
	module.ExposedDefinitions().AddDefinitions(module.LocalDefinitions().Definitions())
	module.SetProgram(program)

	var rootNodesConverted []decorated.Node
	for _, rootNode := range rootNodes {
		converted, couldConver := rootNode.(decorated.Node)
		if !couldConver {
			panic(fmt.Sprintf("can not convert %T", rootNode))
		}
		if converted == nil || reflect.ValueOf(converted).IsNil() {
			panic("can not be nil")
		}
		rootNodesConverted = append(rootNodesConverted, converted)
	}
	module.SetRootNodes(rootNodesConverted)

	return module, returnErr
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
