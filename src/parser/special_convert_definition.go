/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"

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

func parseParameters(p ParseStream, keywordIndentation int, nameOnlyContext ast.LocalTypeNameDefinitionContextDynamic) (
	[]*ast.FunctionParameter, parerr.ParseError,
) {
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
			identifier = nil
			/*
				fakeSymbol := token.NewVariableSymbolToken("_", token.SourceFileReference{
					Range:    p.positionLength(),
					Document: nil,
				}, 0)
				identifier = ast.NewVariableIdentifier(fakeSymbol)

			*/
		}

		var astType ast.Type

		var tErr parerr.ParseError
		astType, tErr = parseTypeReference(p, keywordIndentation, nameOnlyContext, nil)
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
	annotationFunctionType token.AnnotationFunctionType, precedingComments *ast.MultilineComment) (
	ast.Expression, parerr.ParseError,
) {
	keywordIndentation := ident.Symbol().FetchIndentation()

	var returnType ast.Type
	var parameters []*ast.FunctionParameter

	if p.maybeColon() {
		p.eatOneSpace("after colon")
	}

	expressionFollows := false
	nameOnlyDynamicContext := ast.NewLocalTypeNameContext(nil, nil)

	//var functionType *ast.FunctionType
	var typeToUse ast.Type
	if !p.maybeAssign() {
		_, foundLeftParen := p.maybeLeftParen()

		var allTypes []ast.Type

		if foundLeftParen {
			var paramErr parerr.ParseError
			parameters, paramErr = parseParameters(p, keywordIndentation, nameOnlyDynamicContext)
			if paramErr != nil {
				return nil, paramErr
			}

			p.eatOneSpace("after parameters")

			if err := p.eatRightArrow(); err != nil {
				return nil, err
			}

			for _, parameter := range parameters {
				allTypes = append(allTypes, parameter.Type())
			}
		}

		p.eatOneSpace("Return type")

		var tErr parerr.ParseError
		returnType, tErr = parseTypeReference(p, keywordIndentation, nameOnlyDynamicContext, nil)
		if tErr != nil {
			return nil, tErr
		}
		allTypes = append(allTypes, returnType)

		p.eatOneSpace("after arrow")

		expressionFollows = p.maybeAssign()

		functionType := ast.NewFunctionType(allTypes)

		typeToUse = functionType
		if !nameOnlyDynamicContext.IsEmpty() {
			nameOnlyDynamicContext.SetNextType(typeToUse)
			typeToUse = nameOnlyDynamicContext
		}
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

		if returnType == nil {
			return nil, parerr.NewInternalError(p.positionLength(), fmt.Errorf("unsupported assignment / constant"))
		}
	} else {
		expression = ast.NewFunctionDeclarationExpression(ident.Symbol(), annotationFunctionType)
	}

	newFunction := ast.NewFunctionValue(ident.Symbol(), parameters, typeToUse, expression, precedingComments)

	return ast.NewFunctionValueNamedDefinition(ident, newFunction), nil
}
