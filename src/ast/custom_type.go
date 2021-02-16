/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import "fmt"

type CustomType struct {
	name           *TypeIdentifier
	typeParameters []*TypeParameter
	variants       []*CustomTypeVariant
}

func (i *CustomType) String() string {
	return fmt.Sprintf("[custom-type %v %v]", i.name, i.variants)
}

func (i *CustomType) Identifier() *TypeIdentifier {
	return i.name
}

func (i *CustomType) Name() string {
	return i.name.Name()
}

func (i *CustomType) Variants() []*CustomTypeVariant {
	return i.variants
}

func (i *CustomType) FindAllLocalTypes() []*TypeParameter {
	return i.typeParameters
}
