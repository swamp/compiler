/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import "fmt"

type ModuleNamePart struct {
	part *TypeIdentifier
}

type ModuleReference struct {
	parts []*ModuleNamePart
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

func (m *ModuleReference) Parts() []*ModuleNamePart {
	return m.parts
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

func NewModuleReference(parts []*ModuleNamePart) *ModuleReference {
	return &ModuleReference{parts: parts}
}
