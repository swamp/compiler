/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"fmt"
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
	"github.com/swamp/compiler/src/tokenize"
	"github.com/swamp/compiler/src/verbosity"
)

type NoImportModuleRepository struct{}

func (*NoImportModuleRepository) FetchModuleInPackage(moduleType decorated.ModuleType, moduleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
	return nil, decorated.NewInternalError(fmt.Errorf("this is a no import module. Imports are not allowed"))
}

func CompileToModuleOnceForTest(code string, useCores bool, errorsAsWarnings bool) (*decorated.Module, decshared.DecoratedError) {
	rootModules, importModules, rootModuleErr := CreateDefaultRootModule(useCores)
	if rootModuleErr != nil {
		return nil, rootModuleErr
	}

	importRepository := &NoImportModuleRepository{}

	const verbose = verbosity.None

	const enforceStyle = true

	return InternalCompileToModule(decorated.ModuleTypeNormal, importRepository, rootModules, importModules, dectype.MakeArtifactFullyQualifiedModuleName(nil), "for test", code,
		enforceStyle, verbose, errorsAsWarnings)
}

func InternalCompileToProgram(absoluteFilename string, code string, enforceStyle bool, verbose verbosity.Verbosity) (*tokenize.Tokenizer, *ast.SourceFile, decshared.DecoratedError) {
	ioReader := strings.NewReader(code)
	runeReader, runeReaderErr := runestream.NewRuneReader(ioReader, absoluteFilename)
	if runeReaderErr != nil {
		return nil, nil, decorated.NewInternalError(runeReaderErr)
	}

	tokenizer, tokenizerErr := tokenize.NewTokenizer(runeReader, enforceStyle)
	if tokenizerErr != nil {
		const errorsAsWarnings = false
		parser.ShowError(tokenizer, absoluteFilename, tokenizerErr, verbose, errorsAsWarnings)
		return tokenizer, nil, tokenizerErr
	}

	p := parser.NewParser(tokenizer, enforceStyle)
	program, programErr := p.Parse()
	if programErr != nil {
		return tokenizer, nil, programErr
	}
	program.SetNodes(p.Nodes())

	return tokenizer, program, nil
}

func isInternalType(typeToCheck dtype.Type) bool {
	_, isPrimitive := typeToCheck.(*dectype.PrimitiveAtom)

	return isPrimitive
}

func checkUnusedImports(module *decorated.Module, importModules []*decorated.Module) {
	for _, definition := range module.ImportedModules().AllModules() {
		if !definition.WasReferenced() && !definition.ReferencedModule().IsInternal() {
			warning := decorated.NewUnusedImportWarning(definition, "importedModules.All()")
			module.AddWarning(warning)
		}
	}

	for _, definition := range module.ImportedDefinitions().ReferencedDefinitions() {
		if !definition.WasReferenced() && !definition.IsInternal() {
			warning := decorated.NewUnusedImportWarning(definition.CreatedBy(), "referenced definitions")
			module.AddWarning(warning)
		}
	}
	/*
		for importedTypeName, importedType := range module.ImportedTypes().AllTypes() {
			if !importedType.WasReferenced() && importedType.CreatedByModuleImport() != nil && !isInternalType(importedType.ReferencedType()) {
				warning := decorated.NewUnusedImportWarning(importedType.CreatedByModuleImport(), fmt.Sprintf("imported types '%v'", importedTypeName))
				module.AddWarning(warning)
			}
		}

	*/
}

func InternalCompileToModule(moduleType decorated.ModuleType, moduleRepository ModuleRepository, aliasModules []*decorated.Module,
	importModules []*decorated.Module, moduleName dectype.ArtifactFullyQualifiedModuleName, absoluteFilename string, code string,
	enforceStyle bool, verbose verbosity.Verbosity, errorAsWarning bool) (*decorated.Module, decshared.DecoratedError) {
	tokenizer, program, programErr := InternalCompileToProgram(absoluteFilename, code, enforceStyle, verbose)
	if programErr != nil {
		parser.ShowError(tokenizer, absoluteFilename, programErr, verbose, errorAsWarning)
		return nil, programErr
	}

	module := decorated.NewModule(moduleType, moduleName, tokenizer.Document())

	for _, aliasModule := range aliasModules {
		if err := CopyModuleToModule(module, aliasModule); err != nil {
			return nil, err
		}
	}

	for _, importModule := range importModules {
		if importModule == nil {
			panic("importModule is nil")
		}
		// relativeModuleName := dectype.MakePackageRelativeModuleName(importModule.FullyQualifiedModuleName().Path())
		fakeModuleReference := decorated.NewModuleReference(importModule.FullyQualifiedModuleName().Path(), importModule)
		fakeImportStatement := decorated.NewImport(nil, fakeModuleReference, fakeModuleReference, false)

		importErr := ImportModuleToModule(module, fakeImportStatement)
		if importErr != nil {
			return nil, decorated.NewInternalError(importErr)
		}
	}

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

	checkUnusedImports(module, importModules)

	var returnErr decshared.DecoratedError

	if len(allErrors) > 0 {
		returnErr = decorated.NewMultiErrors(allErrors)
	}

	module.ExposedTypes().AddTypesFromModule(module.LocalTypes().AllTypes(), module)
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

	exposedTypes := source.ExposedTypes().AllTypes()
	exposedDefinitions := source.ExposedDefinitions().ReferencedDefinitions()

	importedModule := target.ImportedModules().ImportModule(statement.ImportAsName().AstModuleReference(), source, statement.Module())

	if statement.ExposeAll() {
		target.ImportedTypes().AddTypes(exposedTypes, importedModule)
		for _, exposedDefinition := range exposedDefinitions {
			importedModule := decorated.NewImportedModule(target,
				statement.Module())
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
	return target.LocalTypes().CopyTypes(source.LocalTypes().AllTypes())
}

func ExposeEverythingInModule(target *decorated.Module) {
	target.ExposedTypes().AddTypesFromModule(target.LocalTypes().AllTypes(), target)
}
