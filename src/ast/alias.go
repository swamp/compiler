/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import "fmt"

type Alias struct {
	aliasName      *TypeIdentifier
	referencedType Type
}

func (i *Alias) String() string {
	return fmt.Sprintf("[alias-type %v]", i.referencedType)
}

func (i *Alias) Name() string {
	return "Alias"
}

func (i *Alias) Identifier() *TypeIdentifier {
	return i.aliasName
}

func (i *Alias) ReferencedType() Type {
	return i.referencedType
}
