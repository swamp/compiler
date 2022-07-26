/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type AnnotationStatement struct {
	astAnnotation *ast.Annotation
	t             dtype.Type
	hasLocalTypes bool
}

func NewAnnotation(astAnnotation *ast.Annotation, t dtype.Type) *AnnotationStatement {
	funcType := t.(*dectype.FunctionTypeReference)
	hasLocalTypes := false
	for _, param := range funcType.FunctionAtom().FunctionParameterTypes() {
		_, wasLocalType := param.(*dectype.LocalType)
		if wasLocalType {
			hasLocalTypes = true
			break
		}
	}
	return &AnnotationStatement{astAnnotation: astAnnotation, t: t, hasLocalTypes: hasLocalTypes}
}

func (d *AnnotationStatement) Identifier() *ast.VariableIdentifier {
	return d.astAnnotation.Identifier()
}

func (d *AnnotationStatement) HasLocalTypes() bool {
	return d.hasLocalTypes
}

func (d *AnnotationStatement) Annotation() *ast.Annotation {
	return d.astAnnotation
}

func (d *AnnotationStatement) String() string {
	return fmt.Sprintf("[annotation %v=%v]", d.astAnnotation.Identifier(), d.t)
}

func (d *AnnotationStatement) StatementString() string {
	return fmt.Sprintf("[annotation %v=%v]", d.astAnnotation.Identifier(), d.t)
}

func (d *AnnotationStatement) Type() dtype.Type {
	return d.t
}

func (d *AnnotationStatement) FetchPositionLength() token.SourceFileReference {
	return d.astAnnotation.FetchPositionLength()
}
