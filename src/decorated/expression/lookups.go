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

type LookupField struct {
	structField *dectype.RecordField
}

func (l LookupField) String() string {
	return fmt.Sprintf("[lookup %v]", l.structField)
}

func (l LookupField) Index() int {
	return l.structField.Index()
}

func NewLookupField(structField *dectype.RecordField) LookupField {
	return LookupField{structField: structField}
}

type LookupVariable struct {
	name *ast.VariableIdentifier
	expr dtype.Type
}

func (l LookupVariable) Identifier() *ast.VariableIdentifier {
	return l.name
}

func (l LookupVariable) String() string {
	return fmt.Sprintf("[lookupvar %v (%v)]", l.name, l.expr)
}

func (l LookupVariable) DecoratedExpression() dtype.Type {
	return l.expr
}

func NewLookupVariable(name *ast.VariableIdentifier, expr dtype.Type) LookupVariable {
	return LookupVariable{name: name, expr: expr}
}

type Lookups struct {
	DecoratedExpressionNode
	variableLookup LookupVariable
	lookupFields   []LookupField
}

func NewLookups(variableLookup LookupVariable, lookupFields []LookupField) *Lookups {
	l := &Lookups{variableLookup: variableLookup, lookupFields: lookupFields}
	count := len(lookupFields)
	l.decoratedType = lookupFields[count-1].structField.Type()
	return l
}

func (l *Lookups) Variable() LookupVariable {
	return l.variableLookup
}

func (l *Lookups) LookupFields() []LookupField {
	return l.lookupFields
}

func (l *Lookups) String() string {
	return fmt.Sprintf("[lookups %v %v]", l.variableLookup, l.lookupFields)
}

func (l *Lookups) FetchPositionAndLength() token.PositionLength {
	return token.PositionLength{}
}
