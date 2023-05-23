/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"bytes"
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type GuardItemDefault struct {
	consequence   Expression `debug:"true"`
	internalIndex int
	guardDefault  *ast.GuardDefault
}

func NewGuardItemDefault(guardDefault *ast.GuardDefault, internalIndex int, consequence Expression) *GuardItemDefault {
	return &GuardItemDefault{guardDefault: guardDefault, internalIndex: internalIndex, consequence: consequence}
}

func (c *GuardItemDefault) String() string {
	return c.consequence.String()
}

func (c *GuardItemDefault) Expression() Expression {
	return c.consequence
}

func (c *GuardItemDefault) InternalIndex() int {
	return c.internalIndex
}

func (c *GuardItemDefault) AstGuardDefault() *ast.GuardDefault {
	return c.guardDefault
}

type GuardItem struct {
	condition     Expression `debug:"true"`
	consequence   Expression `debug:"true"`
	internalIndex int
	astGuardItem  ast.GuardItem
}

func NewGuardItem(astGuardItem ast.GuardItem, internalIndex int, condition Expression,
	consequence Expression) *GuardItem {
	return &GuardItem{astGuardItem: astGuardItem, internalIndex: internalIndex, condition: condition,
		consequence: consequence}
}

func (c *GuardItem) Expression() Expression {
	return c.consequence
}

func (c *GuardItem) Condition() Expression {
	return c.condition
}

func (c *GuardItem) InternalIndex() int {
	return c.internalIndex
}

func (c *GuardItem) AstGuardItem() ast.GuardItem {
	return c.astGuardItem
}

func (c *GuardItem) String() string {
	return fmt.Sprintf("[DGuardItem %v %v]", c.condition, c.consequence)
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
	items        []*GuardItem      `debug:"true"`
	defaultGuard *GuardItemDefault `debug:"true"`
	astGuard     *ast.GuardExpression
}

func NewGuard(astGuard *ast.GuardExpression, items []*GuardItem, defaultGuard *GuardItemDefault) (*Guard,
	decshared.DecoratedError) {
	return &Guard{astGuard: astGuard, items: items, defaultGuard: defaultGuard}, nil
}

func (i *Guard) Type() dtype.Type {
	if len(i.items) == 0 {
		return i.defaultGuard.consequence.Type()
	}
	firstGuard := i.items[0]
	return firstGuard.Expression().Type()
}

func (i *Guard) String() string {
	if i.defaultGuard != nil {
		return fmt.Sprintf("[DGuard: %v default: %v]", guardConsequenceArrayToStringEx(i.items, ";"), i.defaultGuard)
	}
	return fmt.Sprintf("[DGuard %v]", guardConsequenceArrayToStringEx(i.items, ";"))
}

func (i *Guard) Items() []*GuardItem {
	return i.items
}

func (i *Guard) DefaultGuard() *GuardItemDefault {
	return i.defaultGuard
}

func (i *Guard) AstGuard() *ast.GuardExpression {
	return i.astGuard
}

func (i *Guard) DebugString() string {
	return "[dguard]"
}

func (i *Guard) FetchPositionLength() token.SourceFileReference {
	return i.astGuard.FetchPositionLength()
}
