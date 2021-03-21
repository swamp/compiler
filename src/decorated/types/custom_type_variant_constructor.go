/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type CustomTypeVariantConstructorType struct {
	variant *CustomTypeVariant
}

func NewCustomTypeVariantConstructorType(variant *CustomTypeVariant) *CustomTypeVariantConstructorType {
	return &CustomTypeVariantConstructorType{variant: variant}
}

func (s *CustomTypeVariantConstructorType) Variant() *CustomTypeVariant {
	return s.variant
}

func (s *CustomTypeVariantConstructorType) String() string {
	return fmt.Sprintf("[variantconstr %v]", s.variant)
}

func (s *CustomTypeVariantConstructorType) HumanReadable() string {
	return fmt.Sprintf("%v", s.variant.HumanReadable())
}

func (s *CustomTypeVariantConstructorType) ParameterCount() int {
	return s.variant.ParameterCount()
}

func (s *CustomTypeVariantConstructorType) Resolve() (dtype.Atom, error) {
	return nil, fmt.Errorf("could not generate")
}

func (s *CustomTypeVariantConstructorType) Next() dtype.Type {
	return s.variant
}

func (s *CustomTypeVariantConstructorType) FetchPositionLength() token.SourceFileReference {
	return s.variant.FetchPositionLength()
}
