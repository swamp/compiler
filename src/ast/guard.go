/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type GuardItemBasic struct {
	Comment     token.Comment
	Consequence Expression
	GuardToken  token.GuardToken
}

func NewGuardItemBasic(comment token.Comment, guardToken token.GuardToken, consequence Expression) GuardItemBasic {
	return GuardItemBasic{
		Comment:     comment,
		Consequence: consequence,
		GuardToken:  guardToken,
	}
}

func (i GuardItemBasic) String() string {
	return fmt.Sprintf("%v", i.Consequence)
}

type GuardItem struct {
	GuardItemBasic
	Condition Expression
}

func (i GuardItem) String() string {
	return fmt.Sprintf("[%v => %v]", i.Condition, i.GuardItemBasic)
}

type GuardDefault struct {
	GuardItemBasic
}

func (i *GuardDefault) String() string {
	return fmt.Sprintf("[_ => %v]", i.GuardItemBasic)
}

type GuardExpression struct {
	items       []GuardItem
	defaultItem *GuardDefault
	inclusive   token.SourceFileReference
}

func NewGuardExpression(items []GuardItem, defaultItem *GuardDefault) *GuardExpression {
	lastRange := items[len(items)-1].Consequence.FetchPositionLength()
	if defaultItem != nil {
		lastRange = defaultItem.Consequence.FetchPositionLength()
	}
	start := items[0].GuardToken.FetchPositionLength()
	inclusive := token.MakeInclusiveSourceFileReference(start, lastRange)
	return &GuardExpression{items: items, defaultItem: defaultItem, inclusive: inclusive}
}

func (i *GuardExpression) Items() []GuardItem {
	return i.items
}

func (i *GuardExpression) Default() *GuardDefault {
	return i.defaultItem
}

func (i *GuardExpression) FetchPositionLength() token.SourceFileReference {
	return i.inclusive
}

func (i *GuardExpression) String() string {
	return fmt.Sprintf("[Guard %v %v]", i.items, i.defaultItem)
}

func (i *GuardExpression) DebugString() string {
	return "[guard]"
}
