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

func parseTupleTypeReference(p ParseStream, keywordIndentation int,
	startParen token.ParenToken, typeParameterContext *ast.TypeParameterIdentifierContext,
	precedingComments *ast.MultilineComment, term ast.Type) (ast.Type, parerr.ParseError) {
	var types []ast.Type
	var endParen token.ParenToken
	types = append(types, term)
	for {
		term, err := parseTypeTermReference(p, keywordIndentation, typeParameterContext, precedingComments)
		if err != nil {
			return nil, err
		}
		types = append(types, term)
		if _, wasComma := p.maybeComma(); wasComma {
			p.eatOneSpace("after comma")
		} else {
			foundParen, endParenErr := p.readRightParen()
			if endParenErr != nil {
				return nil, err
			}
			endParen = foundParen
			break
		}
	}

	return ast.NewTupleType(startParen, endParen, types), nil
}

func parseTypeReference(p ParseStream, keywordIndentation int,
	typeParameterContext *ast.TypeParameterIdentifierContext,
	precedingComments *ast.MultilineComment) (ast.Type, parerr.ParseError) {
	term, tErr := parseTypeTermReference(p, keywordIndentation, typeParameterContext, precedingComments)
	if tErr != nil {
		return nil, tErr
	}

	someTerminationFound := p.detectOneSpaceAndTermination()
	currentIndentation := keywordIndentation
	if someTerminationFound {
		if p.maybeOneSpaceAndRightArrow() {
			newIndentation, _, eatErr := p.eatContinuationReturnIndentation(currentIndentation)
			if eatErr != nil {
				return nil, eatErr
			}
			currentIndentation = newIndentation
			var functionTypes []ast.Type
			functionTypes = append(functionTypes, term)
			for {
				t, tErr := parseTypeTermReference(p, currentIndentation, typeParameterContext, precedingComments)
				if tErr != nil {
					return nil, tErr
				}
				functionTypes = append(functionTypes, t)
				_, _, beforeSpaceErr := p.maybeOneSpace()
				if beforeSpaceErr != nil {
					return nil, beforeSpaceErr
				}
				continues := p.maybeRightArrow()
				if !continues {
					break
				}
				_, _, spaceAfterArrowErr := p.eatContinuationReturnIndentation(keywordIndentation)
				if spaceAfterArrowErr != nil {
					return nil, spaceAfterArrowErr
				}
			}
			newFunctionType := ast.NewFunctionType(functionTypes)

			return newFunctionType, nil
		}
	}

	return term, nil
}

func parseTypeReferenceFunc(p ParseStream, keywordIndentation int,
	typeParameterContext *ast.TypeParameterIdentifierContext,
	precedingComments *ast.MultilineComment) (ast.Type, parerr.ParseError) {
	t, tErr := parseTypeReference(p, keywordIndentation, typeParameterContext, precedingComments)
	if tErr != nil {
		return nil, tErr
	}

	// hack
	_, isFunc := t.(*ast.FunctionType)
	if !isFunc {
		t = ast.NewFunctionType([]ast.Type{t})
	}

	return t, nil
}
