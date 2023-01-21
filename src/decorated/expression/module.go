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
	CreateSomeTypeReference(someTypeIdentifier ast.TypeIdentifierNormalOrScoped) (dectype.TypeReferenceScopedOrNormal, decshared.DecoratedError)
}

type TypeAddAndReferenceMaker interface {
	TypeReferenceMaker
	AddTypeAlias(alias *dectype.Alias) TypeError
	AddCustomType(customType *dectype.CustomTypeAtom) TypeError
	FindBuiltInType(s string) dtype.Type
	SourceModule() *Module
}

type FullyQualifiedPackageVariableName struct {
	module     *Module
	identifier *ast.VariableIdentifier
}

func NewFullyQualifiedVariableName(module *Module, identifier *ast.VariableIdentifier) *FullyQualifiedPackageVariableName {
	if module == nil {
		panic(fmt.Sprintf("must have module in %v", identifier))
	}

	return &FullyQualifiedPackageVariableName{module: module, identifier: identifier}
}

func (q *FullyQualifiedPackageVariableName) ResolveToString() string {
	return q.module.fullyQualifiedModuleName.JoinLocalName(q.identifier)
}

func (q *FullyQualifiedPackageVariableName) String() string {
	return fmt.Sprintf("%v", q.ResolveToString())
}

func (q *FullyQualifiedPackageVariableName) Identifier() *ast.VariableIdentifier {
	return q.identifier
}

func ModuleTypeToString(moduleType ModuleType) string {
	switch moduleType {
	case ModuleTypeNormal:
		return "normal"
	case ModuleTypeFromPath:
		return "shared library"
	case ModuleTypeFromEnvironment:
		return "environment"
	}
	panic("unknown module type")
}

type ModuleType = int

const (
	ModuleTypeNormal ModuleType = iota
	ModuleTypeFromPath
	ModuleTypeFromEnvironment
)

type Module struct {
	localTypes        *ModuleTypes
	localDefinitions  *ModuleDefinitions
	localDeclarations *ModuleDeclarations

	importedTypes       *ExposedTypes
	importedModules     *ModuleImports
	importedDefinitions *ModuleImportedDefinitions

	exposedTypes       *ExposedTypes
	exposedDefinitions *ModuleReferenceDefinitions

	program *ast.SourceFile

	// externalFunctions []*ExternalFunctionDeclaration

	isInternal               bool
	sourceFileUri            *token.SourceFileDocument
	fullyQualifiedModuleName dectype.ArtifactFullyQualifiedModuleName
	rootNodes                []Node
	expandedRootNodes        []*ExpandedNode
	references               []*ModuleReference
	moduleType               ModuleType
}

func NewModule(moduleType ModuleType, fullyQualifiedModuleName dectype.ArtifactFullyQualifiedModuleName, sourceFileUri *token.SourceFileDocument) *Module {
	m := &Module{
		fullyQualifiedModuleName: fullyQualifiedModuleName,
		sourceFileUri:            sourceFileUri,
		moduleType:               moduleType,
		importedModules:          NewModuleImports(),
	}

	m.exposedTypes = NewExposedTypes(m)
	m.importedTypes = NewExposedTypes(m)
	m.localTypes = NewModuleTypes(m)
	m.localDefinitions = NewModuleDefinitions(m)
	m.localDeclarations = NewModuleDeclarations(m)

	m.importedDefinitions = NewModuleImportedDefinitions(m)
	m.exposedDefinitions = NewModuleReferenceDefinitions(m)

	return m
}

func (m *Module) FetchPositionLength() token.SourceFileReference {
	return token.MakeSourceFileReference(m.sourceFileUri, token.NewPositionLength(token.NewPositionTopLeft(), 1))
}

func (m *Module) ModuleType() ModuleType {
	return m.moduleType
}

func (m *Module) AddReference(ref *ModuleReference) {
	m.references = append(m.references, ref)
}

func (m *Module) References() []*ModuleReference {
	return m.references
}

func (m *Module) SetProgram(program *ast.SourceFile) {
	m.program = program
}

func (m *Module) Document() *token.SourceFileDocument {
	return m.sourceFileUri
}

func (m *Module) SetRootNodes(nodes []Node) {
	if len(nodes) == 0 {
		panic("must have root expandedRootNodes expandedRootNodes")
	}
	m.rootNodes = nodes
	if false {
		log.Printf("all root expandedRootNodes in: %v\n", m.FullyQualifiedModuleName())
		for _, x := range m.rootNodes {
			log.Printf("root node: %v %v (%T)\n", x.FetchPositionLength(), x, x)
		}
	}
	expandedNodes := ExpandAllChildNodes(nodes)
	if len(expandedNodes) == 0 {
		log.Printf("must have expanded expandedRootNodes")
	}

	m.expandedRootNodes = nil

	moduleLocalPath, _ := m.sourceFileUri.Uri.ToLocalFilePath()
	log.Printf("all root nodes in: %v root node count: %d", m.FullyQualifiedModuleName(), len(nodes))
	for _, node := range expandedNodes {
		sourceRef := node.node.FetchPositionLength()
		if sourceRef.Document == nil {
			panic(fmt.Sprintf("source ref can not be nil for %v", node))
			continue
		}
		localPath, _ := sourceRef.Document.Uri.ToLocalFilePath()
		if localPath != moduleLocalPath {
			panic(fmt.Sprintf("how can this be %v vs %v", localPath, m.sourceFileUri))
			log.Printf("not matching %v and %v\n", localPath, moduleLocalPath)
			continue
		}
		//log.Printf("   node: %v %v (%T)", node.FetchPositionLength(), node, node)
	}

	m.expandedRootNodes = expandedNodes
}

func (m *Module) RootNodes() []Node {
	return m.rootNodes
}

func (m *Module) ExpandedNodes() []*ExpandedNode {
	return m.expandedRootNodes
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

func (m *Module) FullyQualifiedName(identifier *ast.VariableIdentifier) *FullyQualifiedPackageVariableName {
	if m == nil {
		panic(fmt.Sprintf("how is this possible %v\n", identifier))
	}
	return NewFullyQualifiedVariableName(m, identifier)
}

func (m *Module) FullyQualifiedModuleName() dectype.ArtifactFullyQualifiedModuleName {
	return m.fullyQualifiedModuleName
}

func (m *Module) LocalTypes() *ModuleTypes {
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

func (m *Module) LocalDefinitions() *ModuleDefinitions {
	return m.localDefinitions
}

func (m *Module) Declarations() *ModuleDeclarations {
	return m.localDeclarations
}

func (m *Module) ImportedDefinitions() *ModuleImportedDefinitions {
	return m.importedDefinitions
}

func (m *Module) ExposedDefinitions() *ModuleReferenceDefinitions {
	return m.exposedDefinitions
}

func (m *Module) LocalAndImportedDefinitions() *ModuleDefinitionsCombine {
	importAndLocal := NewModuleDefinitionsCombine(m.LocalDefinitions(), m.ImportedDefinitions(), m.importedModules)

	return importAndLocal
}

func (m *Module) DebugOutput(debug string) {
	log.Printf("%v: \n", debug)
	log.Println(m.String())
}

func (m *Module) ShortString() string {
	s := m.localTypes.DebugString()
	s += "\n" + m.localDefinitions.ShortString()
	return s
}

func (m *Module) String() string {
	s := "------------ " + m.fullyQualifiedModuleName.String() + " ----------- \nexposed:\n"
	s += m.exposedTypes.DebugString()
	s += "\nexposed definitions:\n"
	s += m.exposedDefinitions.ShortString()
	s += "\nimported modules:\n"
	s += m.importedModules.String()
	s += "\n--------------------------------\n"

	s += "\nimported definitions:\n"
	s += m.importedDefinitions.ShortString()

	s += "\nlocal types:\n"
	s += m.localTypes.DebugString()
	s += "\nlocal definitions:\n"
	s += m.localDefinitions.String()
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
