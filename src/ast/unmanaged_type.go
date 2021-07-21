/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type UnmanagedType struct {
	startParen             token.ParenToken
	endParen               token.ParenToken
	comment                *MultilineComment
	nativeLanguageTypeName *TypeIdentifier
	inclusive              token.SourceFileReference
}

func NewUnmanagedType(startParen token.ParenToken, endParen token.ParenToken, nativeLanguageTypeName *TypeIdentifier, comment *MultilineComment) *UnmanagedType {
	inclusive := token.MakeInclusiveSourceFileReference(startParen.SourceFileReference, endParen.SourceFileReference)
	return &UnmanagedType{nativeLanguageTypeName: nativeLanguageTypeName, inclusive: inclusive, comment: comment}
}

func (i *UnmanagedType) NativeLanguageTypeName() *TypeIdentifier {
	return i.nativeLanguageTypeName
}

func (i *UnmanagedType) Name() string {
	return "UnmanagedType"
}

func (i *UnmanagedType) String() string {
	return fmt.Sprintf("[unmanaged-type %v]", i.nativeLanguageTypeName)
}

func (i *UnmanagedType) Comment() *MultilineComment {
	return i.comment
}

func (i *UnmanagedType) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}
