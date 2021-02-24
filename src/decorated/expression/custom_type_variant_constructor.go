/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type CustomTypeVariantConstructor struct {
	arguments         []DecoratedExpression
	customTypeVariant *dectype.CustomTypeVariant
	typeIdentifier    *ast.TypeIdentifier
}

func NewCustomTypeVariantConstructor(typeIdentifier *ast.TypeIdentifier, customTypeVariant *dectype.CustomTypeVariant,
	arguments []DecoratedExpression) *CustomTypeVariantConstructor {
	if customTypeVariant == nil {
		panic("custom customTypeVariant is nil")
	}

	if customTypeVariant.InCustomType() == nil {
		panic("custom type is nil")
	}

	if customTypeVariant.ParameterCount() != len(arguments) {
		panic(fmt.Sprintf("custom type variant constructor. wrong number of arguments %v %v", customTypeVariant.ParameterCount(), arguments))
	}

	return &CustomTypeVariantConstructor{
		typeIdentifier: typeIdentifier, customTypeVariant: customTypeVariant,
		arguments: arguments,
	}
}

func (c *CustomTypeVariantConstructor) CustomTypeVariantIndex() int {
	return c.customTypeVariant.Index()
}

func (c *CustomTypeVariantConstructor) CustomTypeVariant() *dectype.CustomTypeVariant {
	return c.customTypeVariant
}

func (c *CustomTypeVariantConstructor) Arguments() []DecoratedExpression {
	return c.arguments
}

func (c *CustomTypeVariantConstructor) Type() dtype.Type {
	var resolvedTypes []dtype.Type
	for _, resolved := range c.arguments {
		resolvedTypes = append(resolvedTypes, resolved.Type())
	}

	resolvedType, callErr := dectype.CallType(c.customTypeVariant, resolvedTypes)
	if callErr != nil {
		panic(callErr)
	}
	return resolvedType
}

func (c *CustomTypeVariantConstructor) String() string {
	return fmt.Sprintf("[variant-constructor %v %v]", c.customTypeVariant, c.arguments)
}

func (c *CustomTypeVariantConstructor) FetchPositionLength() token.SourceFileReference {
	return c.typeIdentifier.Symbol().SourceFileReference
}
