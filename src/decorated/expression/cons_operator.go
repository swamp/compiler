/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type ConsOperator struct {
	BinaryOperator
}

func NewConsOperator(leftItem Expression, rightList Expression, refMaker TypeAddAndReferenceMaker) (*ConsOperator,
	decshared.DecoratedError) {
	a := &ConsOperator{}
	a.BinaryOperator.left = leftItem
	a.BinaryOperator.right = rightList
	resultType := rightList.Type()
	if dectype.IsListAny(rightList.Type()) {
		if dectype.IsAny(leftItem.Type()) {
			return nil, NewInternalError(fmt.Errorf("cons, both sides are any"))
		}
		listType := refMaker.FindBuiltInType("List", rightList.FetchPositionLength())
		if listType == nil {
			panic("container literal must have a container type defined to use literals")
		}
		unaliasListType := dectype.Unalias(listType)
		collectionType, wasCollectionType := unaliasListType.(*dectype.PrimitiveAtom)
		if !wasCollectionType {
			panic("must have a List type defined to use [] list literals")
		}
		resultType = dectype.NewPrimitiveType(collectionType.PrimitiveName(), []dtype.Type{leftItem.Type()})
	}
	a.BinaryOperator.ExpressionNode.decoratedType = resultType
	return a, nil
}

func (a *ConsOperator) Left() Expression {
	return a.left
}

func (a *ConsOperator) Right() Expression {
	return a.right
}

func (a *ConsOperator) String() string {
	return fmt.Sprintf("[cons left:%v right:%v]", a.left, a.right)
}

func (a *ConsOperator) FetchPositionLength() token.SourceFileReference {
	inclusive := token.MakeInclusiveSourceFileReference(a.left.FetchPositionLength(), a.right.FetchPositionLength())
	return inclusive
}
