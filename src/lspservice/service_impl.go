package lspservice

import (
	"fmt"
	"log"

	swampcompiler "github.com/swamp/compiler/src/compiler"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/token"
)

type LspImpl struct {
	world  *loader.World
	module *decorated.Module
}

func (l *LspImpl) Compile(filename string) error {
	const enforceStyle = true
	const verboseFlag = false
	world, module, err := swampcompiler.CompileFile(filename, enforceStyle, verboseFlag)
	if err != nil {
		return err
	}
	if module == nil {
		return fmt.Errorf("module can not be nil!")
	}
	l.world = world
	l.module = module
	return nil
}

func splitUpFunctionValue(fn *decorated.FunctionValue) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	tokens = append(tokens, splitUp(fn.Expression())...)
	for _, parameter := range fn.Parameters() {
		log.Printf("checking param:%v '%v'\n", parameter.FetchPositionLength(), parameter.String())
		tokens = append(tokens, splitUp(parameter)...)
	}
	return tokens
}

func splitUpAnnotation(fn *decorator.LocalAnnotation) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	tokens = append(tokens, splitUp(fn.Type())...)
	return tokens
}

func splitUpFunctionType(fn *dectype.FunctionAtom) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	for _, parameter := range fn.FunctionParameterTypes() {
		log.Printf("TYPES:%v (%v)\n", parameter, parameter.FetchPositionLength())
		tokens = append(tokens, splitUp(parameter)...)
	}
	return tokens
}

func splitUpLetAssignment(assignment *decorated.LetAssignment) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	tokens = append(tokens, splitUp(assignment.LetVariable())...)
	tokens = append(tokens, splitUp(assignment.Expression())...)

	return tokens
}

func splitUpLet(let *decorated.Let) []DecoratedTypeOrToken {
	var tokens []DecoratedTypeOrToken
	for _, assignment := range let.Assignments() {
		log.Printf("checking assignment:%v '%v'\n", assignment.Expression().FetchPositionLength(), assignment.String())
		tokens = append(tokens, splitUp(assignment)...)
	}

	tokens = append(tokens, splitUp(let.Consequence())...)

	return tokens
}

func splitUp(node decorated.Node) []DecoratedTypeOrToken {
	tokens := []DecoratedTypeOrToken{node}
	switch t := node.(type) {
	case *decorator.LocalAnnotation:
		return append(tokens, splitUpAnnotation(t)...)
	case *decorated.FunctionValue:
		return append(tokens, splitUpFunctionValue(t)...)
	case *dectype.FunctionAtom:
		return append(tokens, splitUpFunctionType(t)...)
	case *decorated.Let:
		return append(tokens, splitUpLet(t)...)
	case *decorated.LetAssignment:
		return append(tokens, splitUpLetAssignment(t)...)
	default:
		log.Printf("do not know how to fix this: %T %v\n", t, t)
		return tokens
	}
}

func (l *LspImpl) FindToken(position token.Position) DecoratedTypeOrToken {
	if l.module == nil {
		return nil
	}
	allNodes := l.module.Nodes()

	var tokens []DecoratedTypeOrToken
	for _, node := range allNodes {
		tokens = append(tokens, splitUp(node)...)
	}

	smallestRange := token.MakeRange(
		token.MakePosition(0, 0),
		token.MakePosition(9999999, 0))

	var bestToken DecoratedTypeOrToken
	for _, decoratedToken := range tokens {
		log.Printf("checking node:%v '%v'\n", decoratedToken.FetchPositionLength(), decoratedToken.String())
		foundRange := decoratedToken.FetchPositionLength().Range
		if foundRange.Contains(position) {
			if foundRange.SmallerThan(smallestRange) {
				smallestRange = foundRange
				bestToken = decoratedToken
			}
		}

	}
	return bestToken
}
