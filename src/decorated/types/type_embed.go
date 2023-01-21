/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type InvokerType struct {
	typeToInvoke dtype.Type
	params       []dtype.Type
	inclusive    token.SourceFileReference
}

func (u *InvokerType) HumanReadable() string {
	return fmt.Sprintf("%v %v", u.typeToInvoke.HumanReadable(), TypesToHumanReadable(u.params))
}

func (u *InvokerType) TypeGenerator() dtype.Type {
	return u.typeToInvoke
}

func (u *InvokerType) Params() []dtype.Type {
	return u.params
}

func (u *InvokerType) FetchPositionLength() token.SourceFileReference {
	return u.inclusive
}

func (u *InvokerType) String() string {
	return fmt.Sprintf("%v<%v>", u.typeToInvoke.String(), TypesToString(u.params))
}

func (u *InvokerType) Resolve() (dtype.Atom, error) {
	anotherType, callErr := CallType(u.typeToInvoke, u.params)
	if callErr != nil {
		return nil, callErr
	}

	return anotherType.Resolve()
}

func (u *InvokerType) ParameterCount() int {
	return 0
}

func (u *InvokerType) Next() dtype.Type {
	return nil
}

func NewInvokerType(typeToInvoke dtype.Type, params []dtype.Type) (*InvokerType, decshared.DecoratedError) {
	if len(params) != typeToInvoke.ParameterCount() {
		return nil, &InternalError{fmt.Errorf("wrong parameter count")}
	}
	for _, param := range params {
		if param == nil {
			panic("sorry we have nil here in InvokerType")
		}
	}
	log.Printf("invoker %T %v (%v) ", typeToInvoke, typeToInvoke.FetchPositionLength().ToStartAndEndReferenceString(), params)
	for _, x := range params {
		log.Printf(".. params %T %v (%v) ", x, x.FetchPositionLength().Range, x)

	}
	inclusive := token.MakeInclusiveSourceFileReference(params[0].FetchPositionLength(), params[len(params)-1].FetchPositionLength())
	return &InvokerType{params: params, typeToInvoke: typeToInvoke, inclusive: inclusive}, nil
}

func (u *InvokerType) WasReferenced() bool {
	return false // Invoker types are not reused
}

type InternalError struct {
	Err error
}

func (e *InternalError) FetchPositionLength() token.SourceFileReference {
	return token.SourceFileReference{}
}

func (e *InternalError) Error() string {
	return e.Err.Error()
}
