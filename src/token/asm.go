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
	functionName   string
	parameterCount uint
	raw            string
	commentToken   IndentationReport
}

func NewExternalFunctionToken(raw string, functionName string, parameterCount uint, commentToken IndentationReport, startPosition SourceFileReference) ExternalFunctionToken {
	return ExternalFunctionToken{raw: raw, functionName: functionName, parameterCount: parameterCount, commentToken: commentToken, SourceFileReference: startPosition}
}

func (s ExternalFunctionToken) Raw() string {
	return s.raw
}

func (s ExternalFunctionToken) Type() Type {
	return ExternalFunction
}

func (s ExternalFunctionToken) ExternalFunction() string {
	return s.functionName
}

func (s ExternalFunctionToken) ParameterCount() uint {
	return s.parameterCount
}

func (s ExternalFunctionToken) String() string {
	return fmt.Sprintf("[externalfn %v %d]", s.functionName, s.parameterCount)
}

func (s ExternalFunctionToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
