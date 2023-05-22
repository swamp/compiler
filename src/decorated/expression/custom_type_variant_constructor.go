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
	arguments                  []Expression                        `debug:"true"`
	customTypeVariantReference *dectype.CustomTypeVariantReference `debug:"true"`
	returnType                 dtype.Type
	inclusive                  token.SourceFileReference
}

func NewCustomTypeVariantConstructor(customTypeVariantReference *dectype.CustomTypeVariantReference,
	arguments []Expression) *CustomTypeVariantConstructor {
	inclusive := customTypeVariantReference.FetchPositionLength()
	if len(arguments) > 0 {
		inclusive = token.MakeInclusiveSourceFileReference(customTypeVariantReference.FetchPositionLength(), arguments[len(arguments)-1].FetchPositionLength())
	}

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

	var returnType dtype.Type

	returnType = customTypeVariantReference.CustomTypeVariant()

	if len(arguments) > 0 {
		if len(arguments) != customTypeVariant.ParameterCount() {
			panic("wrong number of parameters to variant constructor")
		}
		var types []dtype.Type
		//foundLocal := false
		for index, x := range arguments {
			originalType := customTypeVariant.ParameterTypes()[index]
			_, wasLocal := originalType.(*dectype.ResolvedLocalType)
			if wasLocal {
				//		foundLocal = true
			}
			types = append(types, x.Type())
		}

	}
	/*

		genericContext := customTypeVariantReference.CustomTypeVariant().ResolvedLocalTypeContext()
		if genericContext.HasDefinitions() {
			concretizedCustomTypeVariant := concretize.ConcretizeCustomTypeVariant(customTypeVariantReference, arguments)
			invokerType, typeErr := dectype.NewInvokerType(customTypeVariantReference, types)
			if typeErr != nil {
				panic(typeErr)
			}
			returnType = invokerType
		}
	*/
	return &CustomTypeVariantConstructor{
		customTypeVariantReference: customTypeVariantReference,
		arguments:                  arguments,
		inclusive:                  inclusive,
		returnType:                 returnType,
	}
}

func (c *CustomTypeVariantConstructor) Reference() *dectype.CustomTypeVariantReference {
	return c.customTypeVariantReference
}

func (c *CustomTypeVariantConstructor) CustomTypeVariantIndex() int {
	return c.customTypeVariantReference.CustomTypeVariant().Index()
}

func (c *CustomTypeVariantConstructor) CustomTypeVariant() *dectype.CustomTypeVariantAtom {
	return c.customTypeVariantReference.CustomTypeVariant()
}

func (c *CustomTypeVariantConstructor) Arguments() []Expression {
	return c.arguments
}

func (c *CustomTypeVariantConstructor) Type() dtype.Type {
	return c.returnType
}

func (c *CustomTypeVariantConstructor) String() string {
	return fmt.Sprintf("[VariantConstructor %v %v]", c.customTypeVariantReference.CustomTypeVariant(), c.arguments)
}

func (c *CustomTypeVariantConstructor) HumanReadable() string {
	return "Custom Type Variant Constructor"
}

func (c *CustomTypeVariantConstructor) FetchPositionLength() token.SourceFileReference {
	return c.inclusive
}
