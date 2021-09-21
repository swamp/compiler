/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
)

// Keyword :
type AsmToken struct {
	SourceFileReference
	asm string
}

func NewAsmToken(asm string, startPosition SourceFileReference) AsmToken {
	return AsmToken{asm: asm, SourceFileReference: startPosition}
}

func (s AsmToken) Type() Type {
	return Asm
}

func (s AsmToken) Raw() string {
	return s.asm
}

func (s AsmToken) Asm() string {
	return s.asm
}

func (s AsmToken) String() string {
	return fmt.Sprintf("[asm %v]", s.asm)
}

func (s AsmToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}

// Keyword :
type ExternalFunctionToken struct {
	SourceFileReference
}

func NewExternalFunctionToken(startPosition SourceFileReference) ExternalFunctionToken {
	return ExternalFunctionToken{SourceFileReference: startPosition}
}

func (s ExternalFunctionToken) Type() Type {
	return ExternalFunction
}

func (s ExternalFunctionToken) String() string {
	return fmt.Sprintf("[externalfn]")
}

func (s ExternalFunctionToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}

// Keyword :
type ExternalVarFunction struct {
	SourceFileReference
}

func NewExternalVarFunction(startPosition SourceFileReference) ExternalVarFunction {
	return ExternalVarFunction{SourceFileReference: startPosition}
}

func (s ExternalVarFunction) Type() Type {
	return ExternalFunction
}

func (s ExternalVarFunction) String() string {
	return fmt.Sprintf("[externalvarfn]")
}

func (s ExternalVarFunction) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
