/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
)

type CustomTypeAtom struct {
	nameToField           map[string]*CustomTypeVariant
	parameters            []dtype.Type
	variants              []*CustomTypeVariant
	name                  *ast.TypeIdentifier
	artifactTypeName      ArtifactFullyQualifiedTypeName
	genericLocalTypeNames []*dtype.TypeArgumentName
}

func (s *CustomTypeAtom) String() string {
	return fmt.Sprintf("[custom-type %v]", s.variants)
}

func (s *CustomTypeAtom) HumanReadable() string {
	str := fmt.Sprintf("type %v = ", s.name)
	for _, variant := range s.variants {
		str += "\n"
		str += variant.HumanReadable()
	}
	return str
}

func (s *CustomTypeAtom) ShortString() string {
	str := "[custom-type "
	for _, variant := range s.variants {
		str += " " + variant.ShortString()
	}

	str += "]"
	return str
}

func (s *CustomTypeAtom) TypeIdentifier() *ast.TypeIdentifier {
	return s.name
}

func (s *CustomTypeAtom) DecoratedName() string {
	return s.name.Name()
}

func (s *CustomTypeAtom) AtomName() string {
	return s.DecoratedName()
}

func (s *CustomTypeAtom) ShortName() string {
	return s.DecoratedName()
}

func (s *CustomTypeAtom) Name() string {
	return s.DecoratedName()
}

func (s *CustomTypeAtom) ArtifactTypeName() ArtifactFullyQualifiedTypeName {
	return s.artifactTypeName
}

func (s *CustomTypeAtom) ConcretizedName() string {
	return s.DecoratedName()
}

func NewCustomType(name *ast.TypeIdentifier, artifactTypeName ArtifactFullyQualifiedTypeName, genericLocalTypeNames []*dtype.TypeArgumentName, variants []*CustomTypeVariant) *CustomTypeAtom {
	nameToField := make(map[string]*CustomTypeVariant)
	for index, variant := range variants {
		key := variant.Name().Name()
		if index != variant.Index() {
			panic("internal error. index in variant enum")
		}
		if nameToField[key] != nil {
			panic("can not have several enum with same name")
		}
		nameToField[key] = variant
	}

	return &CustomTypeAtom{name: name, artifactTypeName: artifactTypeName, genericLocalTypeNames: genericLocalTypeNames, variants: variants, nameToField: nameToField}
}

func (s *CustomTypeAtom) HasVariant(variantToLookFor *CustomTypeVariant) bool {
	for _, variant := range s.variants {
		if variant == variantToLookFor {
			return true
		}
	}
	return false
}

func (s *CustomTypeAtom) ParameterCount() int {
	return len(s.genericLocalTypeNames)
}

func (s *CustomTypeAtom) Generate(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("sorry")
}

func (u *CustomTypeAtom) Apply(params []dtype.Type) (dtype.Type, error) {
	return nil, fmt.Errorf("custom type does not have apply")
}

func (s *CustomTypeAtom) Resolve() (dtype.Atom, error) {
	return s, nil
}

func (s *CustomTypeAtom) Next() dtype.Type {
	return nil
}

func (s *CustomTypeAtom) IsVariantEqual(otherVariant *CustomTypeVariant) error {
	for _, variant := range s.variants {
		if variant.index == otherVariant.index && variant.name.Name() == otherVariant.name.Name() &&
			len(variant.parameterTypes) == len(otherVariant.parameterTypes) {
			for index, variantParam := range variant.parameterTypes {
				otherParam := otherVariant.parameterTypes[index]
				compatibleErr := CompatibleTypes(variantParam, otherParam)
				if compatibleErr == nil {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("couldn't find it")
}

func (u *CustomTypeAtom) IsEqual(other_ dtype.Atom) error {
	other, wasFunctionAtom := other_.(*CustomTypeAtom)
	if !wasFunctionAtom {
		return fmt.Errorf("was not even a custom type %v", other)
	}
	otherParams := other.variants
	if len(u.variants) != len(otherParams) {
		return fmt.Errorf("different number of variants %v %v", u.variants, otherParams)
	}
	for index, variant := range u.variants {
		otherParam := otherParams[index]
		if variant.Name() != otherParam.Name() {
			return fmt.Errorf("not same variants %v %v", variant, otherParam)
		}
		types := variant.ParameterTypes()
		otherTypes := otherParam.ParameterTypes()
		if len(types) != len(otherTypes) {
			return fmt.Errorf("variants had different number of type params %v %v", types, otherTypes)
		}

		for index, resolveType := range types {
			otherType, otherErr := otherTypes[index].Resolve()
			if otherErr != nil {
				return fmt.Errorf("variant had different type params %v %v", resolveType, otherType)
			}

			resolveType, resolveErr := resolveType.Resolve()
			if resolveErr != nil {
				return fmt.Errorf("variant params resolved to different types %w", resolveErr)
			}
			equalErr := resolveType.IsEqual(otherType)
			if equalErr != nil {
				return equalErr
			}
		}
	}

	return nil
}

func (s *CustomTypeAtom) Variants() []*CustomTypeVariant {
	return s.variants
}

func (s *CustomTypeAtom) VariantCount() int {
	return len(s.variants)
}

func (s *CustomTypeAtom) FindVariant(name string) *CustomTypeVariant {
	return s.nameToField[name]
}
