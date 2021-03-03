/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type ModuleNamePart struct {
	part *TypeIdentifier
}

func NewModuleNamePart(part *TypeIdentifier) *ModuleNamePart {
	return &ModuleNamePart{part: part}
}

func (m *ModuleNamePart) Name() string {
	return m.part.Name()
}

func (m *ModuleNamePart) TypeIdentifier() *TypeIdentifier {
	return m.part
}

func (m *ModuleNamePart) String() string {
	return m.Name()
}

func (m *ModuleNamePart) FetchPositionLength() token.SourceFileReference {
	return m.part.FetchPositionLength()
}

type ModuleReference struct {
	parts     []*ModuleNamePart
	inclusive token.SourceFileReference
}

func (m *ModuleReference) String() string {
	return fmt.Sprintf("[moduleref %v]", m.parts)
}

func (m *ModuleReference) ModuleName() string {
	s := ""
	for index, part := range m.parts {
		if index > 0 {
			s += "."
		}
		s += part.Name()
	}
	return s
}

func (m *ModuleReference) FetchPositionLength() token.SourceFileReference {
	return m.inclusive
}

func (m *ModuleReference) Parts() []*ModuleNamePart {
	return m.parts
}

func NewModuleReference(parts []*ModuleNamePart) *ModuleReference {
	if len(parts) == 0 {
		panic("must have parts in module reference")
	}
	inclusive := token.MakeInclusiveSourceFileReference(parts[0].FetchPositionLength(), parts[len(parts)-1].FetchPositionLength())
	return &ModuleReference{parts: parts, inclusive: inclusive}
}
