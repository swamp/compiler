/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type LocalAnnotation struct {
	name      *ast.VariableIdentifier
	t         dtype.Type
	inclusive token.SourceFileReference
}

func NewLocalAnnotation(identifier *ast.VariableIdentifier, t dtype.Type) *LocalAnnotation {
	if reflect.ValueOf(t).IsNil() {
		panic("not great")
	}
	inclusive := token.MakeInclusiveSourceFileReference(identifier.FetchPositionLength(), t.FetchPositionLength())
	return &LocalAnnotation{name: identifier, t: t, inclusive: inclusive}
}

func (d *LocalAnnotation) Identifier() *ast.VariableIdentifier {
	return d.name
}

func (d *LocalAnnotation) String() string {
	return fmt.Sprintf("[def %v=%v]", d.name, d.t)
}

func (d *LocalAnnotation) Type() dtype.Type {
	return d.t
}

func (d *LocalAnnotation) FetchPositionLength() token.SourceFileReference {
	return d.inclusive
}
