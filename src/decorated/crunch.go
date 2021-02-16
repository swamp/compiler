/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"fmt"
	"strings"

	dectype "github.com/swamp/compiler/src/decorated/types"

	"github.com/swamp/compiler/src/ast"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/parser"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/tokenize"
)

type NoImportModuleRepository struct{}

func (*NoImportModuleRepository) FetchModuleInPackage(moduleName dectype.PackageRelativeModuleName, verboseFlag bool) (*decorated.Module, decshared.DecoratedError) {
	return nil, decorated.NewInternalError(fmt.Errorf("this is a no import module. Imports are not allowed"))
}

func CompileToModuleOnce(x string, useCores bool, errorsAsWarnings bool) (*decorated.Module, decshared.DecoratedError) {
	rootModules, importModules, rootModuleErr := CreateDefaultRootModule(useCores)
	if rootModuleErr != nil {
		return nil, rootModuleErr
	}
	importRepository := &NoImportModuleRepository{}
	const verbose = true
	const enforceStyle = true
	return InternalCompileToModule(importRepository, rootModules, importModules, dectype.MakeArtifactFullyQualifiedModuleName(nil), "for test", x,
		enforceStyle, verbose, errorsAsWarnings)
}

func InternalCompileToProgram(moduleName dectype.ArtifactFullyQualifiedModuleName, relativeFilename string, x string, enforceStyle bool, verbose bool) (*tokenize.Tokenizer, *ast.Program, decshared.DecoratedError) {
	ioReader := strings.NewReader(x)
	runeReader, runeReaderErr := runestream.NewRuneReader(ioReader, relativeFilename)
	if runeReaderErr != nil {
		return nil, nil, decorated.NewInternalError(runeReaderErr)
	}
	tokenizer, tokenizerErr := tokenize.NewTokenizer(runeReader, enforceStyle)
	if tokenizerErr != nil {
		const errorsAsWarnings = false
		parser.ShowError(tokenizer, relativeFilename, tokenizerErr, verbose, errorsAsWarnings)
		return tokenizer, nil, tokenizerErr
	}
	p := parser.NewParser(tokenizer, enforceStyle)

	program, programErr := p.Parse()
	if programErr != nil {
		return tokenizer, nil, programErr
	}

	return tokenizer, program, nil
}

func InternalCompileToModule(moduleRepository ModuleRepository, aliasModules []*decorated.Module,
	importModules []*decorated.Module, moduleName dectype.ArtifactFullyQualifiedModuleName, relativeFilename string, x string,
	enforceStyle bool, verbose bool, errorAsWarning bool) (*decorated.Module, decshared.DecoratedError) {
	tokenizer, program, programErr := InternalCompileToProgram(moduleName, relativeFilename, x, enforceStyle, verbose)
	if programErr != nil {
		parser.ShowError(tokenizer, relativeFilename, programErr, verbose, errorAsWarning)
		return nil, programErr
	}

	module := decorated.NewModule(moduleName)

	converter := NewDecorator(moduleRepository, module)

	for _, aliasModule := range aliasModules {
		CopyModuleToModule(module, aliasModule)
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

	definerScan := decorator.NewDefiner(converter, module.TypeRepo(), "compiletomodule")
	generateErr := definerScan.Define(program)
	if generateErr != nil {
		parser.ShowError(tokenizer, relativeFilename, generateErr, verbose, errorAsWarning)
		return nil, generateErr
	}

	module.ExposedTypes().AddTypes(module.TypeRepo().AllLocalTypes())
	module.ExposedDefinitions().AddDefinitions(module.Definitions().Definitions())

	return module, nil
}

func ImportModuleToModule(target *decorated.Module, source *decorated.Module, sourceMountedModuleName dectype.PackageRelativeModuleName, exposeAll bool) error {
	if target == nil {
		panic("no target")
	}

	if source == nil {
		panic("no source")
	}

	exposedTypes := source.ExposedTypes().AllExposedTypes()
	exposedDefinitions := source.ExposedDefinitions().ReferencedDefinitions()

	target.ImportedTypes().AddTypesWithModulePrefix(exposedTypes, sourceMountedModuleName)

	importErr := target.ImportedDefinitions().AddDefinitionsWithModulePrefix(source.ExposedDefinitions().ReferencedDefinitions(), sourceMountedModuleName)
	if importErr != nil {
		return importErr
	}

	if exposeAll {
		target.ImportedTypes().AddTypes(exposedTypes)
		target.ImportedDefinitions().AddDefinitions(exposedDefinitions)
	}

	return nil
}

func CopyModuleToModule(target *decorated.Module, source *decorated.Module) {
	target.TypeRepo().CopyTypes(source.TypeRepo().AllLocalTypes())
}

func ExposeEverythingInModule(target *decorated.Module) {
	target.ExposedTypes().AddTypes(target.TypeRepo().AllLocalTypes())
}
