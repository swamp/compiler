/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type LocalTypeNameReference struct {
	variable               *LocalTypeName           `debug:"true"`
	typeParameterReference *LocalTypeNameDefinition `debug:"true"`
}

func (i *LocalTypeNameReference) String() string {
	return fmt.Sprintf("[LocalTypeNameRef %v]", i.typeParameterReference.ident.Name())
}

func (i *LocalTypeNameReference) LocalTypeName() *LocalTypeName {
	return i.variable
}

func (i *LocalTypeNameReference) Name() string {
	return i.typeParameterReference.Identifier().Name()
}

func (i *LocalTypeNameReference) LocalTypeNameDefinition() *LocalTypeNameDefinition {
	return i.typeParameterReference
}

func (i *LocalTypeNameReference) FetchPositionLength() token.SourceFileReference {
	return i.variable.FetchPositionLength()
}

func NewLocalTypeNameReference(variable *LocalTypeName,
	localTypeDefinition *LocalTypeNameDefinition) *LocalTypeNameReference {
	x := &LocalTypeNameReference{typeParameterReference: localTypeDefinition, variable: variable}
	localTypeDefinition.AddReference(x)
	return x
}
