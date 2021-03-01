package decorated

import (
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type TypeOrToken interface {
	String() string
	FetchPositionLength() token.SourceFileReference
}

type Token interface {
	TypeOrToken
	Type() dtype.Type
}

func expandChildNodesFunctionValue(fn *FunctionValue) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Expression())...)
	for _, parameter := range fn.Parameters() {
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesFunctionCall(fn *FunctionCall) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.FunctionValue())...)
	for _, argument := range fn.Arguments() {
		tokens = append(tokens, expandChildNodes(argument)...)
	}
	return tokens
}

func expandChildNodesAnnotation(fn *Annotation) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Type())...)
	return tokens
}

func expandChildNodesFunctionType(fn *dectype.FunctionAtom) []TypeOrToken {
	var tokens []TypeOrToken
	for _, parameter := range fn.FunctionParameterTypes() {
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesFunctionTypeReference(fn *dectype.FunctionTypeReference) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.FunctionAtom())...)
	return tokens
}

func expandChildNodesTypeReference(fn *dectype.TypeReference) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Next())...)
	return tokens
}

func expandChildNodesPrimitive(fn *dectype.PrimitiveAtom) []TypeOrToken {
	var tokens []TypeOrToken
	for _, parameter := range fn.GenericTypes() {
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesInvokerType(fn *dectype.InvokerType) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(fn.TypeGenerator())...)
	for _, param := range fn.Params() {
		tokens = append(tokens, expandChildNodes(param)...)
	}
	return tokens
}

func expandChildNodesLetAssignment(assignment *LetAssignment) []TypeOrToken {
	var tokens []TypeOrToken
	tokens = append(tokens, expandChildNodes(assignment.LetVariable())...)
	tokens = append(tokens, expandChildNodes(assignment.Expression())...)

	return tokens
}

func expandChildNodesListLiteral(listLiteral *ListLiteral) []TypeOrToken {
	var tokens []TypeOrToken
	for _, expression := range listLiteral.Expressions() {
		tokens = append(tokens, expandChildNodes(expression)...)
	}

	return tokens
}

func expandChildNodesLet(let *Let) []TypeOrToken {
	var tokens []TypeOrToken
	for _, assignment := range let.Assignments() {
		tokens = append(tokens, expandChildNodes(assignment)...)
	}

	tokens = append(tokens, expandChildNodes(let.Consequence())...)

	return tokens
}

func expandChildNodes(node Node) []TypeOrToken {
	tokens := []TypeOrToken{node}
	switch t := node.(type) {
	case *Annotation:
		return append(tokens, expandChildNodesAnnotation(t)...)
	case *FunctionValue:
		return append(tokens, expandChildNodesFunctionValue(t)...)
	case *FunctionCall:
		return append(tokens, expandChildNodesFunctionCall(t)...)
	case *Let:
		return append(tokens, expandChildNodesLet(t)...)
	case *LetAssignment:
		return append(tokens, expandChildNodesLetAssignment(t)...)
	case *ListLiteral:
		return append(tokens, expandChildNodesListLiteral(t)...)
	case *dectype.Alias:
		return append(tokens, expandChildNodes(t.Next())...)
	case *dectype.PrimitiveAtom:
		return append(tokens, expandChildNodesPrimitive(t)...)
	case *dectype.InvokerType:
		return append(tokens, expandChildNodesInvokerType(t)...)
	case *dectype.FunctionAtom:
		return append(tokens, expandChildNodesFunctionType(t)...)
	case *dectype.FunctionTypeReference:
		return append(tokens, expandChildNodesFunctionTypeReference(t)...)
	case *dectype.TypeReference:
		return append(tokens, expandChildNodesTypeReference(t)...)
	default:
		//		log.Printf("not handled nodes for expansion: %T %v\n", t, t)
		return tokens
	}
}

func ExpandAllChildNodes(nodes []Node) []TypeOrToken {
	var tokens []TypeOrToken
	for _, node := range nodes {
		tokens = append(tokens, expandChildNodes(node)...)
	}

	return tokens
}
