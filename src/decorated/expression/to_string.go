/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func AtomizersFromTypes(types []dtype.Type) string {
	var atoms []dtype.Atom

	for _, singleType := range types {
		atom, _ := singleType.Resolve()
		atoms = append(atoms, atom)
	}

	return Atomizers(atoms)
}

func Atomizers(atoms []dtype.Atom) string {
	s := ""
	for index, atom := range atoms {
		if index > 0 {
			s += ", "
		}
		s += Atomizer(atom)
	}

	return s
}

func Atomizer(atom dtype.Atom) string {
	switch t := atom.(type) {
	case *dectype.FunctionAtom:
		return "function: " + AtomizersFromTypes(t.FunctionParameterTypes()) + "\n"
	case *dectype.RecordAtom:
		{
			s := fmt.Sprintf("record: %v\n", t.ShortName())
			for _, field := range t.SortedFields() {
				fieldAtom, _ := field.Type().Resolve()
				s += fmt.Sprintf("  %v : %v\n", field.Name(), Atomizer(fieldAtom))
			}
			return s
		}
	case *dectype.CustomTypeAtom:
		{
			s := fmt.Sprintf("customype: %v\n", t.Name())
			for _, field := range t.Variants() {
				fieldAtom, _ := field.Resolve()
				s += fmt.Sprintf("  %v : %v\n", field.Name().Name(), Atomizer(fieldAtom))
			}
			return s
		}
	}

	return atom.AtomName() + "\n"
}

func AtomizerFromType(someType dtype.Type) string {
	atom, _ := someType.Resolve()
	return Atomizer(atom)
}
