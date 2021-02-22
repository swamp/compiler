/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// VariableSymbolToken :
type VariableSymbolToken struct {
	Range
	raw         string
	Indentation int
	sourceFile  *SourceFile
}

func NewVariableSymbolToken(raw string, sourceFile *SourceFile, startPosition Range, indentation int) VariableSymbolToken {
	return VariableSymbolToken{raw: raw, Range: startPosition, sourceFile: sourceFile, Indentation: indentation}
}

func (s VariableSymbolToken) Type() Type {
	return VariableSymbol
}

func (s VariableSymbolToken) Name() string {
	return s.raw
}

func (s VariableSymbolToken) Raw() string {
	return s.raw
}

func (s VariableSymbolToken) FetchIndentation() int {
	return s.Indentation
}

func (s VariableSymbolToken) String() string {
	return fmt.Sprintf("$%s", s.raw)
}

func (s VariableSymbolToken) SourceFile() *SourceFile {
	return s.sourceFile
}
