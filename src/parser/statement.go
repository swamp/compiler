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

func (p *Parser) parseExpressionStatement(precedingComments token.CommentBlock) (ast.Expression, parerr.ParseError) {
	keywordSymbol, keywordSymbolErr := p.stream.tokenizer.ParseStartingKeyword()
	if keywordSymbolErr != nil {
		return nil, keywordSymbolErr

	}

	externalFunction, isExternalFunction := keywordSymbol.(token.ExternalFunctionToken)
	if isExternalFunction {
		return parseExternalFunction(p.stream, externalFunction)
	}

	asm, isAsm := keywordSymbol.(token.AsmToken)
	if isAsm {
		return parseAsm(p.stream, asm)
	}

	variableSymbol, wasVariableSymbol := keywordSymbol.(token.VariableSymbolToken)
	if !wasVariableSymbol {
		return nil, parerr.NewUnknownStatement(variableSymbol)
	}

	switch variableSymbol.Name() {
	case "type":
		return parseCustomType(p.stream, variableSymbol, precedingComments)
	case "import":
		return parseImport(p.stream, variableSymbol, precedingComments)
	case "module":
		p.stream.tokenizer.ReadStringUntilEndOfLine()
		p.stream.tokenizer.MaybeOneNewLine()
		p.stream.tokenizer.MaybeOneNewLine()
		p.stream.tokenizer.MaybeOneNewLine()
		return p.parseExpressionStatement(precedingComments)
	default:
		return checkAndParseAnnotationOrDefinition(p.stream, variableSymbol, precedingComments)
	}
}
