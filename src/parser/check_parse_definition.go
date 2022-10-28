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

func checkAndParseAnnotationOrDefinition(stream ParseStream, keywordIndentation int, variableSymbol token.VariableSymbolToken,
	annotationFunctionType token.AnnotationFunctionType, precedingComments *ast.MultilineComment) (ast.Expression, parerr.ParseError) {
	variableIdentifier := ast.NewVariableIdentifier(variableSymbol)
	_, spaceBeforeAnnotationOrDefinitionErr := stream.eatOneSpace("space after annotation or definition")
	if spaceBeforeAnnotationOrDefinitionErr != nil {
		return nil, parerr.NewExpectedSpacingAfterAnnotationOrDefinition(spaceBeforeAnnotationOrDefinitionErr)
	}

	stream.debugInfo("check maybeColon")

	return parseDefinition(stream, variableIdentifier, annotationFunctionType, precedingComments)
}
