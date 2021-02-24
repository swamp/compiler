/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"bytes"
	"fmt"

	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type GuardItem struct {
	condition     DecoratedExpression
	consequence   DecoratedExpression
	internalIndex int
}

func NewGuardItem(internalIndex int, condition DecoratedExpression, consequence DecoratedExpression) *GuardItem {
	return &GuardItem{internalIndex: internalIndex, condition: condition, consequence: consequence}
}

func (c *GuardItem) Expression() DecoratedExpression {
	return c.consequence
}

func (c *GuardItem) Condition() DecoratedExpression {
	return c.condition
}

func (c *GuardItem) InternalIndex() int {
	return c.internalIndex
}

func (c *GuardItem) String() string {
	return fmt.Sprintf("[dguarditem %v %v]", c.condition, c.consequence)
}

func guardConsequenceArrayToStringEx(expressions []*GuardItem, ch string) string {
	var out bytes.Buffer

	for index, expression := range expressions {
		if index > 0 {
			out.WriteString(ch)
		}
		out.WriteString(expression.String())
	}
	return out.String()
}

type Guard struct {
	items        []*GuardItem
	defaultGuard DecoratedExpression
}

func NewGuard(items []*GuardItem, defaultGuard DecoratedExpression) (*Guard, decshared.DecoratedError) {
	return &Guard{items: items, defaultGuard: defaultGuard}, nil
}

func (i *Guard) Type() dtype.Type {
	if len(i.items) == 0 {
		return i.defaultGuard.Type()
	}
	firstGuard := i.items[0]
	return firstGuard.Expression().Type()
}

func (i *Guard) String() string {
	if i.defaultGuard != nil {
		return fmt.Sprintf("[dguard: %v default: %v]", guardConsequenceArrayToStringEx(i.items, ";"), i.defaultGuard)
	}
	return fmt.Sprintf("[dguard: %v]", guardConsequenceArrayToStringEx(i.items, ";"))
}

func (i *Guard) Test() DecoratedExpression {
	return i.defaultGuard
}

func (i *Guard) Items() []*GuardItem {
	return i.items
}

func (i *Guard) DefaultGuard() DecoratedExpression {
	return i.defaultGuard
}

func (i *Guard) DebugString() string {
	return fmt.Sprintf("[dguard]")
}

func (i *Guard) FetchPositionLength() token.SourceFileReference {
	return i.defaultGuard.FetchPositionLength()
}
