/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type AliasStatement struct {
	name           *TypeIdentifier
	referencedType Type
	keywordAlias   token.Keyword
}

func NewAliasStatement(keywordAlias token.Keyword, name *TypeIdentifier, referencedType Type) *AliasStatement {
	if referencedType == nil {
		panic("alias statement can not be nil")
	}
	return &AliasStatement{keywordAlias: keywordAlias, referencedType: referencedType, name: name}
}

func (i *AliasStatement) String() string {
	return fmt.Sprintf("[alias %v %v]", i.name, i.referencedType)
}

func (i *AliasStatement) Name() string {
	return i.name.Name()
}

func (i *AliasStatement) TypeIdentifier() *TypeIdentifier {
	return i.name
}

func (i *AliasStatement) FetchPositionLength() token.SourceFileReference {
	return i.name.symbolToken.SourceFileReference
}

func (i *AliasStatement) Type() Type {
	return i.referencedType
}

func (i *AliasStatement) DebugString() string {
	return fmt.Sprintf("[alias %v]", i.name.Name())
}
