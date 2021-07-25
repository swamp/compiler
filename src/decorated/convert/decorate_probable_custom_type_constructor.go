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
	variantConstructor, err := d.TypeRepo().CreateSomeTypeReference(call)
	if err != nil {
		return nil, err
	}

	unaliasedConstructor := dectype.Unalias(variantConstructor)

	switch e := unaliasedConstructor.(type) {
	case *dectype.CustomTypeVariantReference:
		return decorated.NewCustomTypeVariantConstructor(e, nil), nil
	default:
		log.Printf("expected a constructor here %T", unaliasedConstructor)

		return nil, decorated.NewInternalError(fmt.Errorf("expected a constructor here %v", unaliasedConstructor))
	}
}
