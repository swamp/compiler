/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

func decorateString(d DecorateStream, str *ast.StringLiteral) (decorated.Expression, decshared.DecoratedError) {
	stringType := d.TypeRepo().FindBuiltInType("String")
	if stringType == nil {
		panic("internal error. String is an unknown type")
	}
	decoratedString := decorated.NewStringLiteral(str, stringType.(*dectype.PrimitiveAtom))
	return decoratedString, nil
}

func decorateStringInterpolation(d DecorateStream, str *ast.StringInterpolation, context *VariableContext) (decorated.Expression, decshared.DecoratedError) {
	decoratedExpression, err := DecorateExpression(d, str.Expression(), context)
	if err != nil {
		return nil, err
	}

	return decoratedExpression, nil
}
