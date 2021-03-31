/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type AnyMatchingType struct {
	operatorToken token.OperatorToken
}

func (i *AnyMatchingType) String() string {
	return fmt.Sprintf("[anymatching-type: %v]", i.operatorToken)
}

func (i *AnyMatchingType) Name() string {
	return "AnyMatchingType"
}

func (i *AnyMatchingType) AsteriskToken() token.OperatorToken {
	return i.operatorToken
}

func (i *AnyMatchingType) FetchPositionLength() token.SourceFileReference {
	return i.operatorToken.FetchPositionLength()
}

func NewAnyMatchingType(operatorToken token.OperatorToken) *AnyMatchingType {
	return &AnyMatchingType{operatorToken: operatorToken}
}
