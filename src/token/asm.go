/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import (
	"fmt"
)

// ExternalFunctionToken :
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
	return "[externalfn]"
}

func (s ExternalFunctionToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}

// ExternalVarFunction :
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

// ExternalVarExFunction :
type ExternalVarExFunction struct {
	SourceFileReference
}

func NewExternalVarExFunction(startPosition SourceFileReference) ExternalVarExFunction {
	return ExternalVarExFunction{SourceFileReference: startPosition}
}

func (s ExternalVarExFunction) Type() Type {
	return ExternalFunction
}

func (s ExternalVarExFunction) String() string {
	return fmt.Sprintf("[externalvarfn]")
}

func (s ExternalVarExFunction) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
