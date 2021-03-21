/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type TypeReferenceMaker interface {
	// CreateTypeReference(typeIdentifier *ast.TypeIdentifier) (*dectype.TypeReference, *dectype.NamedDefinitionTypeReference, decshared.DecoratedError)
	// CreateTypeScopedReference(typeIdentifier *ast.TypeIdentifierScoped) (*dectype.TypeReferenceScoped, *dectype.NamedDefinitionTypeReference, decshared.DecoratedError)
	CreateSomeTypeReference(someTypeIdentifier ast.TypeIdentifierNormalOrScoped) (dectype.TypeReferenceScopedOrNormal, decshared.DecoratedError)
}

type TypeAddAndReferenceMaker interface {
	TypeReferenceMaker
	AddTypeAlias(alias *dectype.Alias) TypeError
	AddCustomType(customType *dectype.CustomTypeAtom) TypeError
	FindBuiltInType(s string) dtype.Type
	SourceModule() *Module
}

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

type ExternalFunctionDeclaration struct {
	AstExternalFunction *ast.ExternalFunction
}

func (d *ExternalFunctionDeclaration) FetchPositionLength() token.SourceFileReference {
	return d.AstExternalFunction.FetchPositionLength()
}

func (d *ExternalFunctionDeclaration) StatementString() string {
	return "external func"
}

func (d *ExternalFunctionDeclaration) String() string {
	return "external func"
}

type Module struct {
	localTypes        *ModuleTypes
	localDefinitions  *ModuleDefinitions
	localDeclarations *ModuleDeclarations

	importedTypes       *ExposedTypes
	importedModules     *ModuleImports
	importedDefinitions *ModuleReferenceDefinitions

	exposedTypes       *ExposedTypes
	exposedDefinitions *ModuleReferenceDefinitions

	program *ast.SourceFile

	externalFunctions []*ExternalFunctionDeclaration

	isInternal               bool
	sourceFileUri            *token.SourceFileDocument
	fullyQualifiedModuleName dectype.ArtifactFullyQualifiedModuleName
	rootNodes                []Node
	nodes                    []TypeOrToken
	references               []*ModuleReference
	errors                   []decshared.DecoratedError
	warnings                 []decshared.DecoratedWarning
}

func NewModule(fullyQualifiedModuleName dectype.ArtifactFullyQualifiedModuleName, sourceFileUri *token.SourceFileDocument) *Module {
	m := &Module{
		fullyQualifiedModuleName: fullyQualifiedModuleName,
		sourceFileUri:            sourceFileUri,
		importedModules:          NewModuleImports(),
	}

	m.exposedTypes = NewExposedTypes(m)
	m.importedTypes = NewExposedTypes(m)
	m.localTypes = NewModuleTypes(m)
	m.localDefinitions = NewModuleDefinitions(m)
	m.localDeclarations = NewModuleDeclarations(m)

	m.importedDefinitions = NewModuleReferenceDefinitions(m)
	m.exposedDefinitions = NewModuleReferenceDefinitions(m)

	return m
}

func (m *Module) FetchPositionLength() token.SourceFileReference {
	return token.MakeSourceFileReference(m.sourceFileUri, token.NewPositionLength(token.NewPositionTopLeft(), 0))
}

func (m *Module) AddReference(ref *ModuleReference) {
	m.references = append(m.references, ref)
}

func (m *Module) AddWarning(warning decshared.DecoratedWarning) {
	log.Printf("%s Warning: '%v' ", warning.FetchPositionLength().ToCompleteReferenceString(), warning.Warning())
	m.warnings = append(m.warnings, warning)
}

func (m *Module) Warnings() []decshared.DecoratedWarning {
	return m.warnings
}

func (m *Module) References() []*ModuleReference {
	return m.references
}

func (m *Module) SetProgram(program *ast.SourceFile) {
	m.program = program
}

func (m *Module) SetErrors(errors []decshared.DecoratedError) {
	m.errors = errors
}

func (m *Module) Errors() []decshared.DecoratedError {
	return m.errors
}

func (m *Module) Document() *token.SourceFileDocument {
	return m.sourceFileUri
}

func (m *Module) SetRootNodes(nodes []Node) {
	m.rootNodes = nodes
	if false {
		log.Printf("all root nodes in: %v\n", m.FullyQualifiedModuleName())
		for _, x := range m.rootNodes {
			log.Printf("root node: %v %v (%T)\n", x.FetchPositionLength(), x, x)
		}
	}
	expandedNodes := ExpandAllChildNodes(nodes)
	m.nodes = nil

	moduleLocalPath, _ := m.sourceFileUri.Uri.ToLocalFilePath()
	// log.Printf("all nodes in: %v\n", m.FullyQualifiedModuleName())
	for _, node := range expandedNodes {
		sourceRef := node.FetchPositionLength()
		if sourceRef.Document == nil {
			// panic(fmt.Sprintf("source ref can not be nil for %v", node))
			continue
		}
		localPath, _ := sourceRef.Document.Uri.ToLocalFilePath()
		if localPath != moduleLocalPath {
			// panic(fmt.Sprintf("how can this be %v vs %v", localPath, m.sourceFileUri))
			// log.Printf("not matching %v and %v\n", localPath, moduleLocalPath)
			continue
		}
		// log.Printf("node: %v %v (%T)\n", x.FetchPositionLength(), x, x)
		m.nodes = append(m.nodes, node)
	}
}

func (m *Module) RootNodes() []Node {
	return m.rootNodes
}

func (m *Module) Nodes() []TypeOrToken {
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

func (m *Module) AddExternalFunction(function *ast.ExternalFunction) *ExternalFunctionDeclaration {
	externalFunc := &ExternalFunctionDeclaration{AstExternalFunction: function}
	m.externalFunctions = append(m.externalFunctions, externalFunc)

	return externalFunc
}

func (m *Module) ExternalFunctions() []*ExternalFunctionDeclaration {
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

func (m *Module) TypeRepo() *ModuleTypes {
	return m.localTypes
}

func (m *Module) ExposedTypes() *ExposedTypes {
	return m.exposedTypes
}

func (m *Module) ImportedTypes() *ExposedTypes {
	return m.importedTypes
}

func (m *Module) ImportedModules() *ModuleImports {
	return m.importedModules
}

func (m *Module) Definitions() *ModuleDefinitions {
	return m.localDefinitions
}

func (m *Module) Declarations() *ModuleDeclarations {
	return m.localDeclarations
}

func (m *Module) ImportedDefinitions() *ModuleReferenceDefinitions {
	return m.importedDefinitions
}

func (m *Module) ExposedDefinitions() *ModuleReferenceDefinitions {
	return m.exposedDefinitions
}

func (m *Module) LocalAndImportedDefinitions() *ModuleDefinitionsCombine {
	importAndLocal := NewModuleDefinitionsCombine(m.Definitions(), m.ImportedDefinitions(), m.importedModules)

	return importAndLocal
}

func (m *Module) DebugOutput(debug string) {
	fmt.Printf("%v: \n", debug)
	fmt.Println(m.DebugString())
}

func (m *Module) ShortString() string {
	return m.localTypes.String() + "\n" + m.localDefinitions.ShortString()
}

func (m *Module) String() string {
	s := "------------ " + m.fullyQualifiedModuleName.String() + " ----------- \nexposed:\n"
	s += m.exposedTypes.DebugString()
	s += "\n"
	s += m.exposedDefinitions.ShortString()
	s += "\nmodule:\n"
	s += m.ShortString()
	s += "\nimported definitions:\n"
	s += m.importedDefinitions.ShortString()
	s += "\nimported modules:\n"
	s += m.importedModules.String()
	s += "\n--------------------------------\n"

	return s
}

func (m *Module) DebugString() string {
	s := "---DEBUG--------- "
	s += m.fullyQualifiedModuleName.String()
	s += " ----------- \n"
	s += m.localTypes.DebugString()
	s += "\n"
	s += m.localDefinitions.DebugString()
	s += "\n-----------------------\n"

	return s
}
