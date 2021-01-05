/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import "fmt"

type LocalType struct {
	typeParameterReference *TypeParameter
}

func (i *LocalType) String() string {
	return fmt.Sprintf("[local-type: %v]", i.typeParameterReference)
}

func (i *LocalType) Name() string {
	return i.typeParameterReference.Identifier().Name()
}

func (i *LocalType) TypeParameter() *TypeParameter {
	return i.typeParameterReference
}
