/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type Annotation struct {
	symbol                 *VariableIdentifier
	annotatedType          Type
	precedingComments      *MultilineComment
	annotationFunctionType token.AnnotationFunctionType
}

func NewAnnotation(variableIdentifier *VariableIdentifier, annotatedType Type, annotationFunctionType token.AnnotationFunctionType,
	precedingComments *MultilineComment) *Annotation {
	if annotatedType == nil {
		panic("must set annotated type")
	}

	return &Annotation{
		symbol: variableIdentifier, annotatedType: annotatedType,
		annotationFunctionType: annotationFunctionType, precedingComments: precedingComments,
	}
}

func (d *Annotation) CommentBlock() *MultilineComment {
	return d.precedingComments
}

func (d *Annotation) AnnotatedType() Type {
	return d.annotatedType
}

func (d *Annotation) Identifier() *VariableIdentifier {
	return d.symbol
}

func (d *Annotation) IsSomeKindOfExternal() bool {
	return d.annotationFunctionType == token.AnnotationFunctionTypeExternal ||
		d.annotationFunctionType == token.AnnotationFunctionTypeExternalVar ||
		d.annotationFunctionType == token.AnnotationFunctionTypeExternalVarEx
}

func (d *Annotation) IsNormalExternal() bool {
	return d.annotationFunctionType == token.AnnotationFunctionTypeExternal
}

func (d *Annotation) IsExternalVarFunction() bool {
	return d.annotationFunctionType == token.AnnotationFunctionTypeExternalVar
}

func (d *Annotation) IsExternalVarExFunction() bool {
	return d.annotationFunctionType == token.AnnotationFunctionTypeExternalVarEx
}

func (d *Annotation) FetchPositionLength() token.SourceFileReference {
	return d.symbol.FetchPositionLength()
}

func (d *Annotation) String() string {
	return fmt.Sprintf("[Annotation %v %v]", d.symbol, d.annotatedType)
}

func (d *Annotation) DebugString() string {
	return d.String()
}
