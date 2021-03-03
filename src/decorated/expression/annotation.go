/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type AnnotationStatement struct {
	name      *ast.VariableIdentifier
	t         dtype.Type
	inclusive token.SourceFileReference
}

func NewAnnotation(identifier *ast.VariableIdentifier, t dtype.Type) *AnnotationStatement {
	if reflect.ValueOf(t).IsNil() {
		panic("not great")
	}
	inclusive := token.MakeInclusiveSourceFileReference(identifier.FetchPositionLength(), t.FetchPositionLength())
	return &AnnotationStatement{name: identifier, t: t, inclusive: inclusive}
}

func (d *AnnotationStatement) Identifier() *ast.VariableIdentifier {
	return d.name
}

func (d *AnnotationStatement) String() string {
	return fmt.Sprintf("[annotation %v=%v]", d.name, d.t)
}

func (d *AnnotationStatement) StatementString() string {
	return fmt.Sprintf("[annotation %v=%v]", d.name, d.t)
}

func (d *AnnotationStatement) Type() dtype.Type {
	return d.t
}

func (d *AnnotationStatement) FetchPositionLength() token.SourceFileReference {
	return d.inclusive
}
