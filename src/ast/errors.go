/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ExtraTypeParametersError struct {
	extraParameter []*LocalTypeNameDefinition
	searchedType   Type
	inclusive      token.SourceFileReference
}

func NewExtraTypeNameParametersError(extraParameter []*LocalTypeNameDefinition, searchedType Type) *ExtraTypeParametersError {
	inclusive := token.MakeInclusiveSourceFileReference(extraParameter[0].FetchPositionLength(), extraParameter[len(extraParameter)-1].FetchPositionLength())
	return &ExtraTypeParametersError{extraParameter: extraParameter, searchedType: searchedType, inclusive: inclusive}
}

func (e *ExtraTypeParametersError) Error() string {
	return fmt.Sprintf("you defined %v but wasn't used in type %v", e.extraParameter, e.searchedType)
}

func (e *ExtraTypeParametersError) FetchPositionLength() token.SourceFileReference {
	return e.inclusive
}

type UndefinedTypeParameterError struct {
	extraParameter *LocalTypeNameDefinition
	context        *LocalTypeNameDefinitionContext
}

func NewUndefinedTypeParameterError(extraParameter *LocalTypeNameDefinition,
	context *LocalTypeNameDefinitionContext) *UndefinedTypeParameterError {
	return &UndefinedTypeParameterError{extraParameter: extraParameter, context: context}
}

func (e *UndefinedTypeParameterError) FetchPositionLength() token.SourceFileReference {
	return e.extraParameter.ident.FetchPositionLength()
}

func (e *UndefinedTypeParameterError) Error() string {
	return fmt.Sprintf("you referenced %v but it wasn't declared in context %v", e.extraParameter, e.context)
}

type UnknownTypeParameterError struct {
	extraParameter *LocalTypeName
	context        *LocalTypeNameDefinitionContext
}

func NewUnknownTypeParameterError(extraParameter *LocalTypeName,
	context *LocalTypeNameDefinitionContext) *UnknownTypeParameterError {
	return &UnknownTypeParameterError{extraParameter: extraParameter, context: context}
}

func (e *UnknownTypeParameterError) FetchPositionLength() token.SourceFileReference {
	return e.extraParameter.ident.FetchPositionLength()
}

func (e *UnknownTypeParameterError) Error() string {
	return fmt.Sprintf("you referenced %v but it wasn't declared in context %v", e.extraParameter, e.context)
}
