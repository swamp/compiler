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
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/parser"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/tokenize"
	"github.com/swamp/compiler/src/verbosity"
)

type NoImportModuleRepository struct{}

func (*NoImportModuleRepository) FetchModuleInPackage(moduleName dectype.PackageRelativeModuleName, verboseFlag verbosity.Verbosity) (*decorated.Module, decshared.DecoratedError) {
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

	return InternalCompileToModule(importRepository, rootModules, importModules, dectype.MakeArtifactFullyQualifiedModuleName(nil), "for test", code,
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

func InternalCompileToModule(moduleRepository ModuleRepository, aliasModules []*decorated.Module,
	importModules []*decorated.Module, moduleName dectype.ArtifactFullyQualifiedModuleName, absoluteFilename string, code string,
	enforceStyle bool, verbose verbosity.Verbosity, errorAsWarning bool) (*decorated.Module, decshared.DecoratedError) {
	tokenizer, program, programErr := InternalCompileToProgram(absoluteFilename, code, enforceStyle, verbose)
	if programErr != nil {
		parser.ShowError(tokenizer, absoluteFilename, programErr, verbose, errorAsWarning)
		return nil, programErr
	}

	module := decorated.NewModule(moduleName, tokenizer.Document())

	for _, aliasModule := range aliasModules {
		if err := CopyModuleToModule(module, aliasModule); err != nil {
			return nil, err
		}
	}

	for _, importModule := range importModules {
		if importModule == nil {
			panic("importModule is nil")
		}
		relativeModuleName := dectype.MakePackageRelativeModuleName(importModule.FullyQualifiedModuleName().Path())
		importErr := ImportModuleToModule(module, importModule, relativeModuleName, false)
		if importErr != nil {
			return nil, decorated.NewInternalError(importErr)
		}

	}

	typeLookup := decorated.NewTypeLookup(module.ImportedModules(), module.TypeRepo(), module.ImportedTypes())
	createAndLookup := decorated.NewTypeCreateAndLookup(typeLookup, module.TypeRepo())

	converter := NewDecorator(moduleRepository, module, createAndLookup)

	definerScan := decorator.NewDefiner(converter, createAndLookup, "compiletomodule")
	var allErrors []decshared.DecoratedError
	rootNodes, generateErr := definerScan.Define(program)
	if generateErr != nil {
		allErrors = append(allErrors, generateErr)
	}
	allErrors = append(allErrors, converter.Errors()...)

	var returnErr decshared.DecoratedError

	if len(allErrors) > 0 {
		returnErr = decorated.NewMultiErrors(allErrors)
	}

	module.ExposedTypes().AddTypes(module.TypeRepo().AllTypes())
	module.ExposedDefinitions().AddDefinitions(module.Definitions().Definitions())
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

func ImportModuleToModule(target *decorated.Module, source *decorated.Module, sourceMountedModuleName dectype.PackageRelativeModuleName, exposeAll bool) error {
	if target == nil {
		panic("no target")
	}

	if source == nil {
		panic("no source")
	}

	exposedTypes := source.ExposedTypes().AllTypes()
	exposedDefinitions := source.ExposedDefinitions().ReferencedDefinitions()

	target.ImportedModules().ImportModule(sourceMountedModuleName.Path(), source)

	if exposeAll {
		target.ImportedTypes().AddTypes(exposedTypes)
		target.ImportedDefinitions().AddDefinitions(exposedDefinitions)

		// HACK
		target.ExposedTypes().AddTypes(exposedTypes)
		target.ExposedDefinitions().AddDefinitions(exposedDefinitions)
	}

	return nil
}

func CopyModuleToModule(target *decorated.Module, source *decorated.Module) decshared.DecoratedError {
	return target.TypeRepo().CopyTypes(source.TypeRepo().AllTypes())
}

func ExposeEverythingInModule(target *decorated.Module) {
	target.ExposedTypes().AddTypes(target.TypeRepo().AllTypes())
}
