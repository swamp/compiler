package generate_ir

import (
	"fmt"
	"github.com/llir/llvm/ir/value"
	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateFunctionReference(funcRef *decorated.FunctionReference, genContext *generateContext) (value.Value, error) {
	fullyQualifiedName := funcRef.NameReference().FullyQualified()
	irFunction := genContext.irFunctions.GetFunc(fullyQualifiedName)
	if irFunction == nil {
		panic(fmt.Errorf("can not find function:%v", fullyQualifiedName))
	}

	return irFunction, nil
}
