/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type Alias struct {
	aliasName       *TypeIdentifier
	xreferencedType Type
	keywordType     token.Keyword
	keywordAlias    token.Keyword
	comment         *MultilineComment
}

func (i *Alias) String() string {
	return fmt.Sprintf("[AliasType %v %v]", i.aliasName, i.xreferencedType)
}

func (i *Alias) DebugString() string {
	return fmt.Sprintf("[AliasType %v %v]", i.aliasName, i.xreferencedType)
}

func (i *Alias) DecoratedName() string {
	return i.aliasName.Name()
}

func (i *Alias) Name() string {
	return i.aliasName.Name()
}

func (i *Alias) Identifier() *TypeIdentifier {
	return i.aliasName
}

func (i *Alias) ReferencedType() Type {
	return i.xreferencedType
}

func (i *Alias) FetchPositionLength() token.SourceFileReference {
	return i.aliasName.FetchPositionLength()
}

func (i *Alias) KeywordType() token.Keyword {
	return i.keywordType
}

func (i *Alias) KeywordAlias() token.Keyword {
	return i.keywordAlias
}

func (i *Alias) Comment() *MultilineComment {
	return i.comment
}

func NewAlias(keywordType token.Keyword, keywordAlias token.Keyword, aliasName *TypeIdentifier,
	referenced Type, comment *MultilineComment) *Alias {
	return &Alias{
		keywordType:     keywordType,
		keywordAlias:    keywordAlias,
		aliasName:       aliasName,
		comment:         comment,
		xreferencedType: referenced,
	}
}
