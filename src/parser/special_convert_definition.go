/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
)

func isConstant(expression ast.Expression) bool {
	switch expression.(type) {
	case *ast.IntegerLiteral:
		return true
	case *ast.StringLiteral:
		return true
	case *ast.CharacterLiteral:
		return true
	case *ast.TypeId:
		return true
	case *ast.ResourceNameLiteral:
		return true
	case *ast.RecordLiteral:
		return true
	case *ast.ListLiteral:
		return true
	case *ast.ArrayLiteral:
		return true
	case *ast.FixedLiteral:
		return true
	}

	return false
}

func parseParameters(p ParseStream, keywordIndentation int) ([]*ast.FunctionParameter, parerr.ParseError) {
	var parameters []*ast.FunctionParameter

	for {
		if p.maybeRightParen() {
			break
		}

		identifier, wasVariable := p.maybeVariableIdentifierWithColon()
		if wasVariable {
			_, skipAfterIdentifierErr := p.eatOneSpace("space after skip identifier in definition")
			if skipAfterIdentifierErr != nil {
				return nil, skipAfterIdentifierErr
			}
		}

		if !wasVariable {
			fakeSymbol := token.NewVariableSymbolToken("_", token.SourceFileReference{
				Range:    token.Range{},
				Document: nil,
			}, 0)
			identifier = ast.NewVariableIdentifier(fakeSymbol)
		}

		var astType ast.Type
		typeParameterContext := ast.NewTypeParameterIdentifierContext(nil)

		var tErr parerr.ParseError
		astType, tErr = parseTypeReference(p, keywordIndentation, typeParameterContext, nil)
		if tErr != nil {
			return nil, tErr
		}

		if _, wasComma := p.maybeComma(); wasComma {
			p.eatOneSpace("after comma")
		}

		parameters = append(parameters, ast.NewFunctionParameter(identifier, astType))
	}

	return parameters, nil
}

func parseDefinition(p ParseStream, ident *ast.VariableIdentifier,
	annotationFunctionType token.AnnotationFunctionType, precedingComments *ast.MultilineComment) (ast.Expression, parerr.ParseError) {
	keywordIndentation := ident.Symbol().FetchIndentation()

	var returnType ast.Type
	var parameters []*ast.FunctionParameter

	if p.maybeColon() {
		p.eatOneSpace("after colon")
	}

	expressionFollows := false
	if !p.maybeAssign() {
		_, foundLeftParen := p.maybeLeftParen()

		if foundLeftParen {
			var paramErr parerr.ParseError
			parameters, paramErr = parseParameters(p, keywordIndentation)
			if paramErr != nil {
				return nil, paramErr
			}

			p.eatOneSpace("after parameters")

			if err := p.eatRightArrow(); err != nil {
				return nil, err
			}
		}

		p.eatOneSpace("Return type")

		typeParameterContext := ast.NewTypeParameterIdentifierContext(nil)

		var tErr parerr.ParseError
		returnType, tErr = parseTypeReference(p, keywordIndentation, typeParameterContext, nil)
		if tErr != nil {
			return nil, tErr
		}

		p.eatOneSpace("after arrow")

		expressionFollows = p.maybeAssign()
	} else {
		expressionFollows = true
	}

	var expression ast.Expression
	if expressionFollows {
		newIndentation, _, indentationErr := p.eatContinuationReturnIndentationAllowComment(keywordIndentation)
		if indentationErr != nil {
			return nil, indentationErr
		}

		expressionIndentation := newIndentation
		var exprErr parerr.ParseError

		expression, exprErr = p.parseExpressionNormal(expressionIndentation)
		if exprErr != nil {
			return nil, exprErr
		}

		if len(parameters) == 0 && isConstant(expression) {
			return ast.NewConstantDefinition(ident, expression, precedingComments), nil
		}
	} else {
		expression = ast.NewFunctionDeclarationExpression(ident.Symbol(), annotationFunctionType)
	}

	newFunction := ast.NewFunctionValue(ident.Symbol(), parameters, returnType, expression, precedingComments)

	return ast.NewFunctionValueNamedDefinition(ident, newFunction), nil
}
