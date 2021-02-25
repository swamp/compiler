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

type Annotation struct {
	name      *ast.VariableIdentifier
	t         dtype.Type
	inclusive token.SourceFileReference
}

func NewLocalAnnotation(identifier *ast.VariableIdentifier, t dtype.Type) *Annotation {
	if reflect.ValueOf(t).IsNil() {
		panic("not great")
	}
	inclusive := token.MakeInclusiveSourceFileReference(identifier.FetchPositionLength(), t.FetchPositionLength())
	return &Annotation{name: identifier, t: t, inclusive: inclusive}
}

func (d *Annotation) Identifier() *ast.VariableIdentifier {
	return d.name
}

func (d *Annotation) String() string {
	return fmt.Sprintf("[def %v=%v]", d.name, d.t)
}

func (d *Annotation) Type() dtype.Type {
	return d.t
}

func (d *Annotation) FetchPositionLength() token.SourceFileReference {
	return d.inclusive
}
