/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"fmt"
	"strings"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/ast/codewriter"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/parser"
	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/tokenize"
	"github.com/swamp/compiler/src/verbosity"
)

func parseToProgram(moduleName string, x string, enforceStyle bool, verbose verbosity.Verbosity) (*tokenize.Tokenizer, *ast.SourceFile, decshared.DecoratedError) {
	ioReader := strings.NewReader(x)
	runeReader, runeReaderErr := runestream.NewRuneReader(ioReader, "")
	if runeReaderErr != nil {
		return nil, nil, decorated.NewInternalError(runeReaderErr)
	}
	tokenizer, tokenizerErr := tokenize.NewTokenizerInternal(runeReader, enforceStyle)
	if tokenizerErr != nil {
		parser.ShowAsError(tokenizer, tokenizerErr)
		return tokenizer, nil, tokenizerErr
	}
	p := parser.NewParser(tokenizer, enforceStyle)

	program, programErr := p.Parse()
	if programErr != nil {
		return tokenizer, nil, programErr
	}

	return tokenizer, program, nil
}

func beautify(moduleName string, code string) (string, decshared.DecoratedError) {
	const doNotForceStyle = false
	const verbose = verbosity.None

	_, program, programErr := parseToProgram(moduleName, code, doNotForceStyle, verbose)
	if programErr != nil {
		return "", programErr
	}

	codeWithColor, withColorErr := codewriter.WriteCode(program, true)
	if withColorErr != nil {
		return "", decorated.NewInternalError(withColorErr)
	}

	fmt.Println(codeWithColor)

	codeWithoutColor, withoutColorErr := codewriter.WriteCode(program, false)
	if withoutColorErr != nil {
		return "", decorated.NewInternalError(withoutColorErr)
	}

	return codeWithoutColor, nil
}
