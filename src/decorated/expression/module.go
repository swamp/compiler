/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type FullyQualifiedVariableName struct {
	module     *Module
	identifier *ast.VariableIdentifier
}

func NewFullyQualifiedVariableName(module *Module, identifier *ast.VariableIdentifier) *FullyQualifiedVariableName {
	if module == nil {
		panic(fmt.Sprintf("must have module in %v", identifier))
	}

	return &FullyQualifiedVariableName{module: module, identifier: identifier}
}

func (q *FullyQualifiedVariableName) ResolveToString() string {
	return q.module.fullyQualifiedModuleName.JoinLocalName(q.identifier)
}

func (q *FullyQualifiedVariableName) String() string {
	return fmt.Sprintf("%v", q.ResolveToString())
}

type ExternalFunction struct {
	FunctionName   string
	ParameterCount uint
}

type Module struct {
	typeRepo *dectype.TypeRepo

	definitions  *ModuleDefinitions
	declarations *ModuleDeclarations

	importedTypes       *dectype.ExposedTypes
	importedDefinitions *ModuleReferenceDefinitions

	exposedTypes       *dectype.ExposedTypes
	exposedDefinitions *ModuleReferenceDefinitions

	program *ast.SourceFile

	externalFunctions []*ExternalFunction

	isInternal               bool
	sourceFile               *token.SourceFileURI
	fullyQualifiedModuleName dectype.ArtifactFullyQualifiedModuleName
	nodes                    []Node
}

func NewModule(fullyQualifiedModuleName dectype.ArtifactFullyQualifiedModuleName, sourceFile *token.SourceFileURI) *Module {
	importedTypes := dectype.NewExposedTypes()
	m := &Module{
		fullyQualifiedModuleName: fullyQualifiedModuleName, exposedTypes: dectype.NewExposedTypes(),
		importedTypes: importedTypes,
		sourceFile:    sourceFile,
	}
	m.typeRepo = dectype.NewTypeRepo(fullyQualifiedModuleName, importedTypes)
	m.definitions = NewModuleDefinitions(m)
	m.importedDefinitions = NewModuleReferenceDefinitions(m)
	m.exposedDefinitions = NewModuleReferenceDefinitions(m)
	m.declarations = NewModuleDeclarations(m)

	return m
}

func (m *Module) SetProgram(program *ast.SourceFile) {
	m.program = program
}

func (m *Module) SetNodes(nodes []Node) {
	m.nodes = nodes

	log.Printf("all nodes in: %v\n", m.FullyQualifiedModuleName())
	for _, x := range m.nodes {
		log.Printf("node: %v %v (%T)\n", x.FetchPositionLength(), x, x)
	}
}

func (m *Module) Nodes() []Node {
	return m.nodes
}

func (m *Module) Program() *ast.SourceFile {
	return m.program
}

func (m *Module) MarkAsInternal() {
	m.isInternal = true
}

func (m *Module) IsInternal() bool {
	return m.isInternal
}

func (m *Module) AddExternalFunction(name string, parameterCount uint) {
	m.externalFunctions = append(m.externalFunctions,
		&ExternalFunction{FunctionName: name, ParameterCount: parameterCount})
}

func (m *Module) ExternalFunctions() []*ExternalFunction {
	return m.externalFunctions
}

func (m *Module) FullyQualifiedName(identifier *ast.VariableIdentifier) *FullyQualifiedVariableName {
	if m == nil {
		panic(fmt.Sprintf("how is this possible %v\n", identifier))
	}
	return NewFullyQualifiedVariableName(m, identifier)
}

func (m *Module) FullyQualifiedModuleName() dectype.ArtifactFullyQualifiedModuleName {
	return m.fullyQualifiedModuleName
}

func (m *Module) TypeRepo() *dectype.TypeRepo {
	return m.typeRepo
}

func (m *Module) ExposedTypes() *dectype.ExposedTypes {
	return m.exposedTypes
}

func (m *Module) ImportedTypes() *dectype.ExposedTypes {
	return m.importedTypes
}

func (m *Module) Definitions() *ModuleDefinitions {
	return m.definitions
}

func (m *Module) Declarations() *ModuleDeclarations {
	return m.declarations
}

func (m *Module) ImportedDefinitions() *ModuleReferenceDefinitions {
	return m.importedDefinitions
}

func (m *Module) ExposedDefinitions() *ModuleReferenceDefinitions {
	return m.exposedDefinitions
}

func (m *Module) LocalAndImportedDefinitions() *ModuleDefinitionsCombine {
	importAndLocal := NewModuleDefinitionsCombine(m.Definitions(), m.ImportedDefinitions())

	return importAndLocal
}

func (m *Module) DebugOutput(debug string) {
	fmt.Printf("%v: \n", debug)
	fmt.Println(m.DebugString())
}

func (m *Module) ShortString() string {
	return m.typeRepo.ShortString() + "\n" + m.definitions.ShortString()
}

func (m *Module) String() string {
	s := "------------ " + m.fullyQualifiedModuleName.String() + " ----------- \nexposed:\n"
	s += m.exposedTypes.ShortString()
	s += "\n"
	s += m.exposedDefinitions.ShortString()
	s += "\nmodule:\n"
	s += m.ShortString()
	s += "\nimported:\n"
	s += m.importedDefinitions.ShortString()
	s += "\n--------------------------------\n"

	return s
}

func (m *Module) DebugString() string {
	s := "---DEBUG--------- "
	s += m.fullyQualifiedModuleName.String()
	s += " ----------- \n"
	s += m.typeRepo.DebugString()
	s += "\n"
	s += m.definitions.DebugString()
	s += "\n-----------------------\n"

	return s
}
