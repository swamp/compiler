/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ExtraTypeParameterError struct {
	extraParameter *TypeParameter
	searchedType   Type
}

func NewExtraTypeParameterError(extraParameter *TypeParameter, searchedType Type) *ExtraTypeParameterError {
	return &ExtraTypeParameterError{extraParameter: extraParameter, searchedType: searchedType}
}

func (e *ExtraTypeParameterError) Error() string {
	return fmt.Sprintf("you defined %v but wasn't used in type %v", e.extraParameter, e.searchedType)
}

func (e *ExtraTypeParameterError) FetchPositionLength() token.PositionLength {
	return e.extraParameter.ident.Symbol().FetchPositionLength()
}

type UndefinedTypeParameterError struct {
	extraParameter *TypeParameter
	context        *TypeParameterIdentifierContext
}

func NewUndefinedTypeParameterError(extraParameter *TypeParameter,
	context *TypeParameterIdentifierContext) *UndefinedTypeParameterError {
	return &UndefinedTypeParameterError{extraParameter: extraParameter, context: context}
}

func (e *UndefinedTypeParameterError) FetchPositionLength() token.PositionLength {
	return e.extraParameter.ident.Symbol().FetchPositionLength()
}

func (e *UndefinedTypeParameterError) Error() string {
	return fmt.Sprintf("you referenced %v but it wasn't declared in context %v", e.extraParameter, e.context)
}
