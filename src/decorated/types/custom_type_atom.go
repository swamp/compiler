/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type CustomTypeAtom struct {
	nameToField      map[string]*CustomTypeVariantAtom
	variants         []*CustomTypeVariantAtom `debug:"true"`
	astCustomType    *ast.CustomType
	arguments        []dtype.Type
	artifactTypeName ArtifactFullyQualifiedTypeName `debug:"true"`
	references       []*CustomTypeReference
	memorySize       MemorySize
	memoryAlign      MemoryAlign
}

func (s *CustomTypeAtom) AstCustomType() *ast.CustomType {
	return s.astCustomType
}

func (s *CustomTypeAtom) Arguments() []dtype.Type {
	return s.arguments
}

func (s *CustomTypeAtom) MemorySize() MemorySize {
	return s.memorySize
}

func (s *CustomTypeAtom) MemoryAlignment() MemoryAlign {
	return s.memoryAlign
}

func (s *CustomTypeAtom) String() string {
	return fmt.Sprintf("[CustomType %v%v %v]", s.artifactTypeName, TypesToStringSuffix(s.arguments), s.variants)
}

func (s *CustomTypeAtom) HumanReadable() string {
	return fmt.Sprintf("%v%v", s.astCustomType.Identifier().Name(), TypesToHumanReadableWithinBrackets(s.arguments))
}

func (s *CustomTypeAtom) FetchPositionLength() token.SourceFileReference {
	return s.astCustomType.FetchPositionLength()
}

func (s *CustomTypeAtom) TypeIdentifier() *ast.TypeIdentifier {
	return s.astCustomType.Identifier()
}

func (s *CustomTypeAtom) DecoratedName() string {
	return s.astCustomType.Identifier().Name()
}

func (s *CustomTypeAtom) AtomName() string {
	return s.DecoratedName()
}

func (s *CustomTypeAtom) StatementString() string {
	return s.DecoratedName()
}

func (s *CustomTypeAtom) Name() string {
	return s.DecoratedName()
}

func (s *CustomTypeAtom) ArtifactTypeName() ArtifactFullyQualifiedTypeName {
	return s.artifactTypeName
}

func calculateTotalSizeAndAlignment(variants []*CustomTypeVariantAtom) (MemorySize, MemoryAlign) {
	maxVariantSize := MemorySize(1)
	maxVariantAlign := MemoryAlign(1)
	for _, variant := range variants {
		offset := MemoryOffset(1) // The union custom type starts with uint8
		maxAlign := MemoryAlign(1)
		for index, field := range variant.parameterFields {
			fieldType := variant.parameterFields[index].Type()
			_, wasLocalType := fieldType.(*ResolvedLocalType)
			if wasLocalType {
				return 0, 0
			}
			memorySize, memoryAlign := GetMemorySizeAndAlignment(fieldType)
			if memorySize == 0 || memoryAlign == 0 {
				//panic(fmt.Errorf("illegal size or align values. index %d, field %T %v in variant %v", index, fieldType, fieldType, variant))
			}

			rest := MemoryAlign(uint32(offset) % uint32(memoryAlign))
			if rest != 0 {
				offset += MemoryOffset(memoryAlign - rest)
			}
			if memoryAlign > maxAlign {
				maxAlign = memoryAlign
			}

			field.memoryOffset = offset
			field.memorySize = memorySize

			offset += MemoryOffset(memorySize)
		}

		rest := MemoryAlign(uint32(offset) % uint32(maxAlign))
		if rest != 0 {
			offset += MemoryOffset(maxAlign - rest)
		}

		variant.debugMemorySize = MemorySize(offset)
		variant.debugMemoryAlign = maxAlign

		if offset > MemoryOffset(maxVariantSize) {
			maxVariantSize = MemorySize(offset)
		}
		if maxAlign > maxVariantAlign {
			maxVariantAlign = maxAlign
		}
	}

	return maxVariantSize, maxVariantAlign
}

func NewCustomTypePrepare(astCustomType *ast.CustomType,
	artifactTypeName ArtifactFullyQualifiedTypeName) *CustomTypeAtom {
	s := &CustomTypeAtom{
		astCustomType: astCustomType, artifactTypeName: artifactTypeName,
	}

	return s
}

func (s *CustomTypeAtom) FinalizeVariants(variants []*CustomTypeVariantAtom) {
	nameToField := make(map[string]*CustomTypeVariantAtom)
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

	memorySize, memoryAlign := calculateTotalSizeAndAlignment(variants)
	if memorySize == 0 {
		memorySize = 1
	}
	if memoryAlign == 0 {
		memoryAlign = 1
	}

	s.variants = variants
	s.nameToField = nameToField
	s.memorySize = memorySize
	s.memoryAlign = memoryAlign
}

func (s *CustomTypeAtom) HasVariant(variantToLookFor *CustomTypeVariantAtom) bool {
	for _, variant := range s.variants {
		if variant == variantToLookFor {
			return true
		}
	}
	return false
}

func (s *CustomTypeAtom) Resolve() (dtype.Atom, error) {
	return s, nil
}

func (s *CustomTypeAtom) Next() dtype.Type {
	return nil
}

func (s *CustomTypeAtom) IsVariantEqual(otherVariant *CustomTypeVariantAtom) error {
	for _, variant := range s.variants {
		if variant.index == otherVariant.index && variant.astCustomTypeVariant.Name() == otherVariant.astCustomTypeVariant.Name() &&
			len(variant.parameterFields) == len(otherVariant.parameterFields) {
			for index, variantParam := range variant.parameterFields {
				otherParam := otherVariant.parameterFields[index]
				compatibleErr := CompatibleTypes(variantParam.parameterType, otherParam.parameterType)
				if compatibleErr != nil {
					return compatibleErr
				}
			}
			return nil
		}
	}

	return fmt.Errorf("couldn't find it")
}

func compareCustomType(u *CustomTypeAtom, other *CustomTypeAtom) error {
	otherVariants := other.variants
	if len(u.variants) != len(otherVariants) {
		return fmt.Errorf("different number of variants %v %v", u.variants, otherVariants)
	}

	for index, variant := range u.variants {
		otherParam := otherVariants[index]
		if variant.Name().Name() != otherParam.Name().Name() {
			return fmt.Errorf("not same variants %v %v", variant, otherParam)
		}
		types := variant.parameterFields
		otherTypes := otherParam.parameterFields
		if len(types) != len(otherTypes) {
			return fmt.Errorf("variants had different number of type params %v %v", types, otherTypes)
		}

		for index, resolveType := range types {
			if err := CompatibleTypes(resolveType.parameterType, otherTypes[index].parameterType); err != nil {
				return fmt.Errorf("wrong in custom type '%s' variant: '%s' parameter:\n%v\nvs\n%v\n%w", u.Name(),
					variant.Name().Name(), resolveType, otherTypes[index], err)
			}
		}
	}

	return nil
}

func (u *CustomTypeAtom) IsEqual(other_ dtype.Atom) error {
	otherCustomType, wasCustomType := other_.(*CustomTypeAtom)
	if wasCustomType {
		return compareCustomType(u, otherCustomType)
	}

	otherVariant, wasCustomTypeVariant := other_.(*CustomTypeVariantAtom)
	if wasCustomTypeVariant {
		if TypeIsTemplateHasLocalTypes(u) {
			if otherVariant.inCustomType.Name() == "Maybe" && otherVariant.Name().Name() == "Nothing" {
				return nil
			}
			log.Printf("I can not compare, since I ama a template %v vs %v", u, other_)
			panic("i can not compare, since I am a template")
		}

		if TypeIsTemplateHasLocalTypes(otherVariant) {
			panic(fmt.Errorf("can not compare variants that are not invoked %v %v", u.Name(), otherVariant.Name()))
		}
		return u.IsVariantEqual(otherVariant)
	}

	return fmt.Errorf("was not even a custom type or variant %T %v", other_, other_)
}

func (s *CustomTypeAtom) Variants() []*CustomTypeVariantAtom {
	return s.variants
}

func (s *CustomTypeAtom) VariantCount() int {
	return len(s.variants)
}

func (s *CustomTypeAtom) FindVariant(name string) *CustomTypeVariantAtom {
	return s.nameToField[name]
}

func (s *CustomTypeAtom) AddReferee(reference *CustomTypeReference) {
	s.references = append(s.references, reference)
}

func (s *CustomTypeAtom) References() []*CustomTypeReference {
	return s.references
}

func (s *CustomTypeAtom) WasReferenced() bool {
	return len(s.references) > 0
}
