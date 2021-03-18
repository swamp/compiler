/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type CustomTypeVariantConstructor struct {
	arguments                  []Expression
	customTypeVariantReference *CustomTypeVariantReference
}

func NewCustomTypeVariantConstructor(customTypeVariantReference *CustomTypeVariantReference,
	arguments []Expression) *CustomTypeVariantConstructor {
	if customTypeVariantReference == nil {
		panic("custom customTypeVariant is nil")
	}

	customTypeVariant := customTypeVariantReference.CustomTypeVariant()
	if customTypeVariant.InCustomType() == nil {
		panic("custom type is nil")
	}

	if customTypeVariant.ParameterCount() != len(arguments) {
		panic(fmt.Sprintf("%v custom type variant constructor. wrong number of arguments %v %v %v", customTypeVariantReference.FetchPositionLength(), customTypeVariantReference, customTypeVariant.ParameterCount(), arguments))
	}

	return &CustomTypeVariantConstructor{
		customTypeVariantReference: customTypeVariantReference,
		arguments:                  arguments,
	}
}

func (c *CustomTypeVariantConstructor) Reference() *CustomTypeVariantReference {
	return c.customTypeVariantReference
}

func (c *CustomTypeVariantConstructor) CustomTypeVariantIndex() int {
	return c.customTypeVariantReference.customTypeVariant.Index()
}

func (c *CustomTypeVariantConstructor) CustomTypeVariant() *dectype.CustomTypeVariant {
	return c.customTypeVariantReference.customTypeVariant
}

func (c *CustomTypeVariantConstructor) Arguments() []Expression {
	return c.arguments
}

func (c *CustomTypeVariantConstructor) Type() dtype.Type {
	var resolvedTypes []dtype.Type
	for _, resolved := range c.arguments {
		resolvedTypes = append(resolvedTypes, resolved.Type())
	}

	resolvedType, callErr := dectype.CallType(c.customTypeVariantReference.customTypeVariant, resolvedTypes)
	if callErr != nil {
		panic(callErr)
	}
	return resolvedType
}

func (c *CustomTypeVariantConstructor) String() string {
	return fmt.Sprintf("[variant-constructor %v %v]", c.customTypeVariantReference.customTypeVariant, c.arguments)
}

func (c *CustomTypeVariantConstructor) HumanReadable() string {
	return fmt.Sprintf("Custom Type Variant Constructor %v", c.customTypeVariantReference.customTypeVariant)
}

func (c *CustomTypeVariantConstructor) FetchPositionLength() token.SourceFileReference {
	return c.customTypeVariantReference.FetchPositionLength()
}
