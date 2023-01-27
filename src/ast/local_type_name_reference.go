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
	typeParameterReference *LocalTypeNameDefinition
}

func (i *LocalTypeNameReference) String() string {
	return fmt.Sprintf("[LocalTypeNameRef %v]", i.typeParameterReference)
}

func (i *LocalTypeNameReference) Name() string {
	return i.typeParameterReference.Identifier().Name()
}

func (i *LocalTypeNameReference) LocalTypeNameDefinition() *LocalTypeNameDefinition {
	return i.typeParameterReference
}

func (i *LocalTypeNameReference) FetchPositionLength() token.SourceFileReference {
	return i.typeParameterReference.FetchPositionLength()
}

func NewLocalTypeNameReference(localTypeDefinition *LocalTypeNameDefinition) *LocalTypeNameReference {
	x := &LocalTypeNameReference{typeParameterReference: localTypeDefinition}
	localTypeDefinition.AddReference(x)
	return x
}
