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

func (p *Parser) parseExpressionStatement(precedingComments *ast.MultilineComment) (ast.Expression, parerr.ParseError) {
	keywordSymbol, keywordSymbolErr := p.stream.tokenizer.ParseStartingKeyword()
	if keywordSymbolErr != nil {
		return nil, keywordSymbolErr
	}

	_, isExternalFunction := keywordSymbol.(token.ExternalFunctionToken)
	if isExternalFunction {
		p.stream.eatOneSpace("after external")
		keywordSymbol, keywordSymbolErr = p.stream.tokenizer.ParseVariableSymbol()
	}

	_, isExternalVarFunction := keywordSymbol.(token.ExternalVarFunction)
	if isExternalVarFunction {
		p.stream.eatOneSpace("after external")
		keywordSymbol, keywordSymbolErr = p.stream.tokenizer.ParseVariableSymbol()
		isExternalFunction = true
	}

	variableSymbol, wasVariableSymbol := keywordSymbol.(token.VariableSymbolToken)
	if !wasVariableSymbol {
		return nil, parerr.NewUnknownStatement(variableSymbol)
	}
	p.stream.nodes = append(p.stream.nodes, variableSymbol)

	switch variableSymbol.Name() {
	case "type":
		keywordType := token.NewKeyword(variableSymbol.Raw(), token.Import, variableSymbol.SourceFileReference)
		return parseCustomType(p.stream, keywordType, precedingComments, variableSymbol.Indentation)
	case "import":
		keywordImport := token.NewKeyword(variableSymbol.Raw(), token.Import, variableSymbol.SourceFileReference)
		return parseImport(p.stream, keywordImport, 0, precedingComments)
	default:
		return checkAndParseAnnotationOrDefinition(p.stream, variableSymbol, isExternalFunction, isExternalVarFunction, precedingComments)
	}
}
