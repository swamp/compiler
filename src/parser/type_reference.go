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

func parseTypeReference(p ParseStream, keywordIndentation int,
	typeParameterContext *ast.TypeParameterIdentifierContext,
	precedingComments token.CommentBlock) (ast.Type, parerr.ParseError) {
	term, tErr := parseTypeTermReference(p, keywordIndentation, typeParameterContext, precedingComments)
	if tErr != nil {
		return nil, tErr
	}

	someTerminationFound := p.detectOneSpaceAndTermination()

	if someTerminationFound {
		if p.maybeOneSpaceAndRightArrow() {

			if _, eatErr := p.eatOneSpace("after right arrow"); eatErr != nil {
				return nil, eatErr
			}
			var functionTypes []ast.Type
			functionTypes = append(functionTypes, term)
			for {
				t, tErr := parseTypeTermReference(p, keywordIndentation, typeParameterContext, precedingComments)
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
				p.eatOneSpace("after continuing right arrow")
			}
			newFunctionType := ast.NewFunctionType(functionTypes)

			return newFunctionType, nil
		}
	}

	return term, nil
}

func parseTypeReferenceFunc(p ParseStream, keywordIndentation int,
	typeParameterContext *ast.TypeParameterIdentifierContext,
	precedingComments token.CommentBlock) (ast.Type, parerr.ParseError) {
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
