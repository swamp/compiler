/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"log"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

// decorateProbableConstructorCall needed for no parameter custom type
// constructors
func decorateProbableConstructorCall(d DecorateStream, call ast.TypeIdentifierNormalOrScoped) (decorated.Expression, decshared.DecoratedError) {
	variantConstructor, err := d.TypeReferenceMaker().CreateSomeTypeReference(call)
	if err != nil {
		return nil, err
	}

	unaliasedConstructor := dectype.Unalias(variantConstructor)

	switch unaliasedConstructor.(type) {
	case *dectype.CustomTypeVariantAtom:
		variantRef, wasVariantRef := variantConstructor.(*dectype.CustomTypeVariantReference)
		if !wasVariantRef {
			panic("illegal variant constructor")
		}
		return decorated.NewCustomTypeVariantConstructor(variantRef, nil), nil
	default:
		log.Printf("expected a constructor here %T", unaliasedConstructor)

		return nil, decorated.NewInternalError(fmt.Errorf("expected a constructor here %v", unaliasedConstructor))
	}
}
