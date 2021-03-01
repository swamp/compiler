/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type GuardItem struct {
	Condition   Expression
	Consequence Expression
}

type GuardExpression struct {
	items       []GuardItem
	defaultItem Expression
}

func NewGuardExpression(items []GuardItem, defaultItem Expression) *GuardExpression {
	return &GuardExpression{items: items, defaultItem: defaultItem}
}

func (i *GuardExpression) Items() []GuardItem {
	return i.items
}

func (i *GuardExpression) DefaultExpression() Expression {
	return i.defaultItem
}

func (i *GuardExpression) FetchPositionLength() token.SourceFileReference {
	return i.items[0].Condition.FetchPositionLength()
}

func (i *GuardExpression) String() string {
	return fmt.Sprintf("[guard: %v %v]", i.items, i.defaultItem)
}

func (i *GuardExpression) DebugString() string {
	return "[guard]"
}
