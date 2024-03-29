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
}

func NewGuardExpression(items []GuardItem, defaultItem *GuardDefault) *GuardExpression {
	return &GuardExpression{items: items, defaultItem: defaultItem}
}

func (i *GuardExpression) Items() []GuardItem {
	return i.items
}

func (i *GuardExpression) Default() *GuardDefault {
	return i.defaultItem
}

func (i *GuardExpression) FetchPositionLength() token.SourceFileReference {
	return i.items[0].Condition.FetchPositionLength()
}

func (i *GuardExpression) String() string {
	return fmt.Sprintf("[Guard %v %v]", i.items, i.defaultItem)
}

func (i *GuardExpression) DebugString() string {
	return "[guard]"
}
