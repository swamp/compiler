package decorated

import (
	"log"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type DecoratedTypeOrToken interface {
	String() string
	FetchPositionLength() token.SourceFileReference
}

type DecoratedToken interface {
	DecoratedTypeOrToken
	Type() dtype.Type
}

func expandChildNodesFunctionValue(fn *FunctionValue) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Expression())...)
	for _, parameter := range fn.Parameters() {
		log.Printf("checking param:%v '%v'\n", parameter.FetchPositionLength(), parameter.String())
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesAnnotation(fn *Annotation) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	tokens = append(tokens, expandChildNodes(fn.Type())...)
	return tokens
}

func expandChildNodesFunctionType(fn *dectype.FunctionAtom) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	for _, parameter := range fn.FunctionParameterTypes() {
		log.Printf("TYPES:%v (%v)\n", parameter, parameter.FetchPositionLength())
		tokens = append(tokens, expandChildNodes(parameter)...)
	}
	return tokens
}

func expandChildNodesLetAssignment(assignment *LetAssignment) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	tokens = append(tokens, expandChildNodes(assignment.LetVariable())...)
	tokens = append(tokens, expandChildNodes(assignment.Expression())...)

	return tokens
}

func expandChildNodesLet(let *Let) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	for _, assignment := range let.Assignments() {
		log.Printf("checking assignment:%v '%v'\n", assignment.Expression().FetchPositionLength(), assignment.String())
		tokens = append(tokens, expandChildNodes(assignment)...)
	}

	tokens = append(tokens, expandChildNodes(let.Consequence())...)

	return tokens
}

func expandChildNodes(node Node) []DecoratedTypeOrToken {
	tokens := []DecoratedTypeOrToken{node}
	switch t := node.(type) {
	case *Annotation:
		return append(tokens, expandChildNodesAnnotation(t)...)
	case *FunctionValue:
		return append(tokens, expandChildNodesFunctionValue(t)...)
	case *dectype.FunctionAtom:
		return append(tokens, expandChildNodesFunctionType(t)...)
	case *Let:
		return append(tokens, expandChildNodesLet(t)...)
	case *LetAssignment:
		return append(tokens, expandChildNodesLetAssignment(t)...)
	default:
		log.Printf("do not know how to fix this: %T %v\n", t, t)
		return tokens
	}
}

func ExpandAllChildNodes(nodes []Node) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	for _, node := range nodes {
		tokens = append(tokens, expandChildNodes(node)...)
	}

	return tokens
}
