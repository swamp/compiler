/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

type AnnotationFunctionType uint8

const (
	AnnotationFunctionTypeNormal AnnotationFunctionType = iota
	AnnotationFunctionTypeExternal
	AnnotationFunctionTypeExternalVar
	AnnotationFunctionTypeExternalVarEx
)

// VariableSymbolToken :
type VariableSymbolToken struct {
	SourceFileReference
	raw         string `debug:"true"`
	Indentation int
}

func NewVariableSymbolToken(raw string, startPosition SourceFileReference, indentation int) VariableSymbolToken {
	return VariableSymbolToken{raw: raw, SourceFileReference: startPosition, Indentation: indentation}
}

func (s VariableSymbolToken) Type() Type {
	return VariableSymbol
}

func (s VariableSymbolToken) Name() string {
	return s.raw
}

func (s VariableSymbolToken) IsIgnore() bool {
	return s.raw == "_"
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

func (s VariableSymbolToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
