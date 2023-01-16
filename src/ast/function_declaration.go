/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type FunctionDeclarationExpression struct {
	annotationFunctionType token.AnnotationFunctionType
	debugAssignedValue     token.VariableSymbolToken
}

func NewFunctionDeclarationExpression(debugAssignedValue token.VariableSymbolToken, annotationFunctionType token.AnnotationFunctionType) *FunctionDeclarationExpression {
	return &FunctionDeclarationExpression{
		annotationFunctionType: annotationFunctionType,
		debugAssignedValue:     debugAssignedValue,
	}
}

func (i *FunctionDeclarationExpression) AnnotationFunctionType() token.AnnotationFunctionType {
	return i.annotationFunctionType
}

func (i *FunctionDeclarationExpression) IsSomeKindOfExternal() bool {
	return i.annotationFunctionType == token.AnnotationFunctionTypeExternal || i.annotationFunctionType == token.AnnotationFunctionTypeExternalVarEx || i.annotationFunctionType == token.AnnotationFunctionTypeExternalVar
}

func (f *FunctionDeclarationExpression) IsExternalVarFunction() bool {
	return f.annotationFunctionType == token.AnnotationFunctionTypeExternalVar
}

func (f *FunctionDeclarationExpression) IsExternalVarExFunction() bool {
	return f.annotationFunctionType == token.AnnotationFunctionTypeExternalVarEx
}

func (i *FunctionDeclarationExpression) FetchPositionLength() token.SourceFileReference {
	return i.debugAssignedValue.FetchPositionLength()
}

func (i *FunctionDeclarationExpression) String() string {
	return fmt.Sprintf("[FnDeclExpr %v]", i.annotationFunctionType)
}

func (i *FunctionDeclarationExpression) DebugString() string {
	return "[function-declaration]"
}
