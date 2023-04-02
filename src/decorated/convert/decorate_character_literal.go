/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func decorateCharacter(d DecorateStream, ch *ast.CharacterLiteral) (decorated.Expression, decshared.DecoratedError) {
	characterType := d.TypeReferenceMaker().FindBuiltInType("Char", ch.FetchPositionLength())
	if characterType == nil {
		panic("internal error. String is an unknown type")
	}
	decoratedCharacter := decorated.NewCharacterLiteral(ch, characterType)
	return decoratedCharacter, nil
}
