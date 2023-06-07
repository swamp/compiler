/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type Statement interface {
	TypeOrToken
	StatementString() string
}

type TypeOrToken interface {
	String() string
	FetchPositionLength() token.SourceFileReference
}

type HumanReadEnabler interface {
	HumanReadable() string
}

type Token interface {
	TypeOrToken
	HumanReadEnabler
	Type() dtype.Type
}

func expandFunctionValue(fn *FunctionValue, list *ExpandedNode) {
	if fn.Expression() == nil {
		panic(fmt.Errorf("a function must have an expression '%v'", fn.FetchPositionLength().ToReferenceString()))
	}

	//for _, parameter := range fn.Parameters() {
	//	expand(parameter, list)
	//}

	//expand(fn.UnaliasedDeclaredFunctionType().ReturnType(), list)
	expand(fn.Expression(), list)

}

func expandConstant(constant *Constant, list *ExpandedNode) {
	expand(constant.Expression(), list)
}

func expandFunctionReference(fn *FunctionReference, list *ExpandedNode) {
	optionalModuleRef := fn.NameReference().ModuleReference()
	if optionalModuleRef != nil {
		//list.AddNode(optionalModuleRef)
	}
}

func expandFunctionCall(fn *FunctionCall, list *ExpandedNode) {
	expand(fn.FunctionExpression(), list)
	for _, argument := range fn.Arguments() {
		expand(argument, list)
	}
}

func expandCurryFunction(fn *CurryFunction, list *ExpandedNode) {
	expand(fn.FunctionValue(), list)
	for _, argument := range fn.ArgumentsToSave() {
		expand(argument, list)
	}
}

func expandImportStatement(importStatement *ImportStatement, list *ExpandedNode) {
	expand(importStatement.ModuleReference(), list)
	if importStatement.Alias() != nil {
		expand(importStatement.Alias(), list)
	}
}

func expandFunctionType(fn *dectype.FunctionAtom, list *ExpandedNode) {
	for _, parameter := range fn.FunctionParameterTypes() {
		expand(parameter, list)
	}
}

func expandTupleType(fn *dectype.TupleTypeAtom, list *ExpandedNode) {
	for _, parameter := range fn.ParameterTypes() {
		expand(parameter, list)
	}
}

func expandNamedTypeReferenceModule(name *dectype.NamedDefinitionTypeReference, list *ExpandedNode) {
	optionalModuleRef := name.ModuleReference()
	if optionalModuleRef != nil {
		//list.AddNode(optionalModuleRef)
	}
}

func expandTypeReference(reference dectype.TypeReferenceScopedOrNormal, list *ExpandedNode) {
	named := reference.NameReference()
	expandNamedTypeReferenceModule(named, list)

	//expand(reference.Next(), list)
}

func expandCustomTypeVariant(variant *dectype.CustomTypeVariantAtom, list *ExpandedNode) {
	//	expand(variant.Name(), list)
	//	for _, param := range variant.ParameterTypes() {
	//		expand(param, list)
	//	}
}

func expandCustomType(fn *dectype.CustomTypeAtom, list *ExpandedNode) {
	//expand(fn.TypeIdentifier(), list)
	for _, variant := range fn.Variants() {
		expandCustomTypeVariant(variant, list)
	}
}

func expandRecordType(fn *dectype.RecordAtom, list *ExpandedNode) {
	for _, field := range fn.ParseOrderedFields() {
		expand(field.FieldName(), list)
		expand(field.Type(), list)
	}
}

func expandUnmanagedType(fn *dectype.UnmanagedType, list *ExpandedNode) {
}

func expandFunctionTypeReference(fn *dectype.FunctionTypeReference, list *ExpandedNode) {
	expand(fn.FunctionAtom(), list)
}

func expandPrimitive(fn *dectype.PrimitiveAtom, list *ExpandedNode) {
	for _, parameter := range fn.ParameterTypes() {
		expand(parameter, list)
	}
}

func expandLetAssignment(assignment *LetAssignment, list *ExpandedNode) {
	for _, param := range assignment.LetVariables() {
		expand(param, list)
	}
	expand(assignment.Expression(), list)
}

func expandListLiteral(listLiteral *ListLiteral, list *ExpandedNode) {
	for _, expression := range listLiteral.Expressions() {
		expand(expression, list)
	}
}

func expandStringInterpolation(stringInterpolation *StringInterpolation, list *ExpandedNode) {
	for _, expression := range stringInterpolation.IncludedExpressions() {
		expand(expression, list)
	}

}

func expandTupleLiteral(tupleLiteral *TupleLiteral, list *ExpandedNode) {
	for _, expression := range tupleLiteral.Expressions() {
		expand(expression, list)
	}
}

func expandArrayLiteral(arrayLiteral *ArrayLiteral, list *ExpandedNode) {
	for _, expression := range arrayLiteral.Expressions() {
		expand(expression, list)
	}
}

func expandRecordLiteral(recordLiteral *RecordLiteral, list *ExpandedNode) {
	if recordLiteral.RecordTemplate() != nil {
		expand(recordLiteral.RecordTemplate(), list)
	}

	for _, assignment := range recordLiteral.ParseOrderedAssignments() {
		expand(assignment.FieldName(), list)
		expand(assignment.Expression(), list)
	}
}

func expandNamedFunctionValue(namedFunctionValue *NamedFunctionValue, list *ExpandedNode) {
	//expand(namedFunctionValue.FunctionName(), list)
	expand(namedFunctionValue.Value(), list)
}

func expandNamedCustomType(namedCustomType *NamedCustomType, list *ExpandedNode) {
	//expand(namedCustomType.FunctionName(), list)
	expand(namedCustomType.customTypeAtom, list)
}

func expandFunctionParameterDefinition(parameter *FunctionParameterDefinition, list *ExpandedNode) {
	expand(parameter.Type(), list)
}

func expandCustomTypeVariantConstructor(constructor *CustomTypeVariantConstructor, list *ExpandedNode) {
	expand(constructor.Reference(), list)

	for _, arg := range constructor.arguments {
		expand(arg, list)
	}
}

func expandRecordConstructor(constructor *RecordConstructorFromParameters, list *ExpandedNode) {
	optionalModuleRef := constructor.NamedTypeReference().ModuleReference()
	if optionalModuleRef != nil {
		//		list.AddNode(optionalModuleRef)
	}

	for _, arg := range constructor.arguments {
		expand(arg.Expression(), list)
	}
}

func expandRecordConstructorRecord(constructor *RecordConstructorFromRecord, list *ExpandedNode) {
	optionalModuleRef := constructor.NamedTypeReference().ModuleReference()
	if optionalModuleRef != nil {
		//list.AddNode(optionalModuleRef)
	}

	expand(constructor.Expression(), list)
}

func expandGuard(guard *Guard, list *ExpandedNode) {
	for _, item := range guard.Items() {
		expand(item.Condition(), list)
		expand(item.Expression(), list)
	}

	if guard.DefaultGuard() != nil {
		expand(guard.DefaultGuard().Expression(), list)
	}
}

func expandCustomTypeVariantReference(constructor *dectype.CustomTypeVariantReference) {
	// expand(constructor.typeIdentifier, list) // TODO: Need meaning
}

func expandCaseForTypeAlias(typeAlias *dectype.Alias, list *ExpandedNode) {
	//tokens = append(tokens, expand(typeAlias.TypeIdentifier())...)
	expand(typeAlias.Next(), list)

}

func expandCaseForCustomType(caseForCustomType *CaseCustomType, list *ExpandedNode) {

	expand(caseForCustomType.Test(), list)

	for _, consequence := range caseForCustomType.Consequences() {
		expand(consequence.VariantReference(), list)
		for _, param := range consequence.Parameters() {
			expand(param, list)
		}
		expand(consequence.Expression(), list)
	}

	if caseForCustomType.DefaultCase() != nil {
		expand(caseForCustomType.DefaultCase(), list)
	}
}

func expandCaseForPatternMatching(caseForCustomType *CaseForPatternMatching, list *ExpandedNode) {
	expand(caseForCustomType.Test(), list)

	for _, consequence := range caseForCustomType.Consequences() {
		expand(consequence.Literal(), list)
		expand(consequence.Expression(), list)
	}

	expand(caseForCustomType.DefaultCase(), list)
}

func expandBinaryOperator(namedFunctionValue *BinaryOperator, list *ExpandedNode) {
	expand(namedFunctionValue.Left(), list)
	expand(namedFunctionValue.Right(), list)
}

func expandForCastOperator(castOperator *CastOperator, list *ExpandedNode) {
	expand(castOperator.Expression(), list)
	expand(castOperator.AliasReference(), list)
}

func expandRecordLookups(lookup *RecordLookups, list *ExpandedNode) {
	expand(lookup.Expression(), list)
	for _, lookupField := range lookup.LookupFields() {
		expand(lookupField.reference, list)
	}
}

func expandLet(let *Let, list *ExpandedNode) {
	for _, assignment := range let.Assignments() {
		expand(assignment, list)
	}

	inConsequnce := let.Consequence()
	expand(inConsequnce, list)
}

func expandIf(ifExpression *If, list *ExpandedNode) {
	expand(ifExpression.Condition(), list)
	expand(ifExpression.Consequence(), list)
	expand(ifExpression.Alternative(), list)
}

type ExpandedNode struct {
	node          TypeOrToken
	children      []*ExpandedNode
	debugDocument *token.SourceFileDocument
}

func (n *ExpandedNode) String() string {
	return fmt.Sprintf("[expand %v %v]", n.node, n.children)
}

func (n *ExpandedNode) AddChildNode(node TypeOrToken) {
	newNode := NewExpandedNode(node, n.debugDocument)
	n.addChild(newNode)
}

func (n *ExpandedNode) debugLog(indent int) {
	log.Printf("%s %v %T ('%v')", strings.Repeat("..", indent), n.node.FetchPositionLength().Range, n.node,
		n.node.FetchPositionLength().ToStartAndEndReferenceString())
	for _, childNode := range n.children {
		childNode.debugLog(indent + 1)
	}
}

func (n *ExpandedNode) DebugLog() {
	log.Printf("showing nodes for document %v", n.debugDocument)
	n.debugLog(0)
}

func (n *ExpandedNode) Children() []*ExpandedNode {
	return n.children
}

func (n *ExpandedNode) HasChildren() bool {
	return len(n.children) > 0
}

func (n *ExpandedNode) TypeOrToken() TypeOrToken {
	return n.node
}

func (n *ExpandedNode) addChild(expandedNode *ExpandedNode) {
	nodeRange := expandedNode.node.FetchPositionLength().Range

	//log.Printf("adding %T %v", expandedNode.node, nodeRange)

	if !n.node.FetchPositionLength().Range.ContainsRange(nodeRange) {
		err := fmt.Errorf("can not add a child that is not within the range of the parent %v %v (%T and %T)",
			n.node.FetchPositionLength().Range, nodeRange, n.node, expandedNode.node)
		log.Print(err)
		panic(err)
	}

	if len(n.children) > 0 && !nodeRange.IsAfter(n.children[len(n.children)-1].node.FetchPositionLength().Range) {
		log.Printf("serious error. tokens are not expanded in range")

		//log.Printf("error in sematic code generation for document %v", n.debugDocument)
		for _, existingNode := range n.children {
			existingRange := existingNode.node.FetchPositionLength().Range
			log.Printf("  expandedNode: %v (%T)", existingRange, existingNode.node)
		}
		log.Printf("--> added incorrect expandedNode: %v (%T) %v", expandedNode.node.FetchPositionLength().Range,
			expandedNode.node, expandedNode.node.FetchPositionLength().ToCompleteReferenceString())
		panic(fmt.Errorf("not in order"))
	}

	n.children = append(n.children, expandedNode)
}

func (n *ExpandedNode) NewGroupNode(node TypeOrToken) *ExpandedNode {
	newRoot := NewExpandedNode(node, n.debugDocument)
	n.addChild(newRoot)

	return newRoot
}

func (n *ExpandedNode) AddChildNodes(addNodes []TypeOrToken) {
	for _, node := range addNodes {
		n.AddChildNode(node)
	}
}

func NewExpandedNode(node TypeOrToken, debugDocument *token.SourceFileDocument) *ExpandedNode {
	return &ExpandedNode{
		node:          node,
		children:      nil,
		debugDocument: debugDocument,
	}
}

type NodeList struct {
	//	expandedRootNodes []TypeOrToken
	rootNode *ExpandedNode

	debugDocument *token.SourceFileDocument
}

type FakeTypeOrToken struct {
	document *token.SourceFileDocument
}

func (n FakeTypeOrToken) String() string {
	return ""
}

func (n FakeTypeOrToken) FetchPositionLength() token.SourceFileReference {
	return token.SourceFileReference{
		Range:    token.MakeRange(token.NewPositionTopLeft(), token.MakePosition(9999, 9999, -1)),
		Document: n.document,
	}
}

func NewNodeList(debugDocument *token.SourceFileDocument) *NodeList {
	return &NodeList{
		rootNode:      NewExpandedNode(FakeTypeOrToken{document: debugDocument}, debugDocument),
		debugDocument: debugDocument,
	}
}

func debugThisDocument(debugDocument *token.SourceFileDocument) bool {
	filepath, _ := debugDocument.Uri.ToLocalFilePath()
	return strings.HasSuffix(filepath, "FindGraphicTiles.swamp")
}

func expand(node Node, parentNode *ExpandedNode) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		panic("can not be nil")
	}

	//	log.Printf("%v expand %T (%v) parent: %T (%v)", node.FetchPositionLength().ToCompleteReferenceString(), node, node, parentNode, parentNode)
	newParentNode := parentNode.NewGroupNode(node)
	switch t := node.(type) {
	case *ModuleReference:
	case *ImportStatement:
		expandImportStatement(t, newParentNode)
	case *FunctionValue:
		expandFunctionValue(t, newParentNode)
	case *FunctionReference:
		expandFunctionReference(t, newParentNode)
	case *FunctionCall:
		expandFunctionCall(t, newParentNode)
	case *CurryFunction:
		expandCurryFunction(t, newParentNode)
	case *Let:
		expandLet(t, newParentNode)
	case *If:
		expandIf(t, newParentNode)
	case *LetAssignment:
		expandLetAssignment(t, newParentNode)
	case *ListLiteral:
		expandListLiteral(t, newParentNode)
	case *TupleLiteral:
		expandTupleLiteral(t, newParentNode)
	case *ArrayLiteral:
		expandArrayLiteral(t, newParentNode)
	case *RecordLiteral:
		expandRecordLiteral(t, newParentNode)
	case *FunctionParameterDefinition:
		expandFunctionParameterDefinition(t, newParentNode)
	case *NamedFunctionValue:
		expandNamedFunctionValue(t, newParentNode)
	case *NamedCustomType:
		expandNamedCustomType(t, newParentNode)
	case *CustomTypeVariantConstructor:
		expandCustomTypeVariantConstructor(t, newParentNode)
	case *RecordConstructorFromParameters:
		expandRecordConstructor(t, newParentNode)
	case *RecordConstructorFromRecord:
		expandRecordConstructorRecord(t, newParentNode)
	case *Guard:
		expandGuard(t, newParentNode)
	case *CaseCustomType:
		expandCaseForCustomType(t, newParentNode)
	case *CaseForPatternMatching:
		expandCaseForPatternMatching(t, newParentNode)
	case *PipeRightOperator:
		expand(&t.BinaryOperator, newParentNode)
	case *PipeLeftOperator:
		expand(&t.BinaryOperator, newParentNode)
	case *ArithmeticOperator:
		expand(&t.BinaryOperator, newParentNode)
	case *LogicalOperator:
		expand(&t.BinaryOperator, newParentNode)
	case *ConsOperator:
		expand(&t.BinaryOperator, newParentNode)
	case *BooleanOperator:
		expand(&t.BinaryOperator, newParentNode)
	case *BitwiseOperator:
		expand(&t.BinaryOperator, newParentNode)
	case *CastOperator:
		expandForCastOperator(t, newParentNode)
	case *CaseConsequenceParameterForCustomType:
	case *ArithmeticUnaryOperator:
		expand(&t.UnaryOperator, newParentNode)
	case *FunctionName: // Should not be expanded

	case *LetVariableReference: // Should not be expanded

	case *LetVariable: // Should not be expanded

	case *RecordTypeFieldReference: // Should not be expanded

	case *FunctionParameterReference: // Should not be expanded

	case *CaseConsequenceParameterReference: // Should not be expanded

	case *IntegerLiteral:
	case *FixedLiteral: // Should not be expanded

	case *CharacterLiteral: // Should not be expanded

	case *TypeIdLiteral: // Should not be expanded

	case *ResourceNameLiteral: // Should not be expanded

	case *StringInterpolation:
		expandStringInterpolation(t, newParentNode)
	case *BooleanLiteral: // Should not be expanded

	case *StringLiteral: // Should not be expanded

	case *MultilineComment: // Should not be expanded

	case *RecordLiteralField: // Should not be expanded

	case *BitwiseUnaryOperator:
		expand(&t.UnaryOperator, newParentNode)
	case *LogicalUnaryOperator:
		expand(&t.UnaryOperator, newParentNode)
	case *UnaryOperator:
		expand(t.Left(), newParentNode)
	case *Constant:
		expandConstant(t, newParentNode)
	case *ConstantReference:

	case *AliasReference:
		expand(t.Type(), newParentNode)
	case *BinaryOperator:
		expandBinaryOperator(t, newParentNode)
	case *RecordLookups:
		expandRecordLookups(t, newParentNode)
	case *IncompleteFunctionCall:
		expand(t.functionValueExpression, newParentNode)
		for _, argument := range t.Arguments() {
			expand(argument, newParentNode)
		}

	case *ExternalFunctionDeclarationExpression:

	case *dectype.ResolvedLocalType:

	case *dectype.AnyMatchingTypes:

	case *dectype.Alias:
		expandCaseForTypeAlias(t, newParentNode)
	case *dectype.PrimitiveAtom:
		expandPrimitive(t, newParentNode)
	case *dectype.FunctionAtom:
		expandFunctionType(t, newParentNode)
	case *dectype.CustomTypeAtom:
		expandCustomType(t, newParentNode)
	case *dectype.RecordAtom:
		expandRecordType(t, newParentNode)
	case *dectype.TupleTypeAtom:
		expandTupleType(t, newParentNode)
	case *dectype.RecordFieldName:
	case *dectype.AliasReference:
		expandTypeReference(t, newParentNode)
	case *dectype.LocalTypeNameOnlyContext:
		expand(t.Next(), newParentNode)
	case *dectype.ResolvedLocalTypeContext:
		expand(t.Next(), newParentNode)
	case *dectype.LocalTypeNameReference:
		break
	case *dectype.LocalTypeNameOnlyContextReference:
		break
	case *dectype.CustomTypeReference:
		expandTypeReference(t, newParentNode)
	case *dectype.PrimitiveTypeReference:
		expandTypeReference(t, newParentNode)
	case *dectype.CustomTypeVariantReference:
		expandTypeReference(t, newParentNode)
	case *dectype.FunctionTypeReference:
		expandFunctionTypeReference(t, newParentNode)
	case *dectype.UnmanagedType:
		expandUnmanagedType(t, newParentNode)
	default:
		panic(fmt.Errorf("expand_nodes: could not expand: %T", t))

	}
}

func ExpandAllChildNodes(nodes []Node) []*ExpandedNode {
	document := nodes[0].FetchPositionLength().Document
	list := NewNodeList(document)
	for _, node := range nodes {
		expand(node, list.rootNode)
	}

	//list.rootNode.DebugLog()

	return list.rootNode.children
}
