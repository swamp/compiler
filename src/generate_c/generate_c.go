package generate_c

import (
	"fmt"
	"io"

	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/verbosity"
)

type Generator struct {
	code *assembler_sp.Code
}

func NewGenerator() *Generator {
	return &Generator{code: assembler_sp.NewCode()}
}

func (g *Generator) GenerateAllLocalDefinedFunctions(module *decorated.Module,
	writer io.Writer, verboseFlag verbosity.Verbosity) error {
	if false {
		fmt.Fprintf(writer, `
typedef int32_t Int;
typedef uint8_t Bool;


`)
	}
	for _, named := range module.LocalDefinitions().Definitions() {
		unknownType := named.Expression()
		_, isConstant := unknownType.(*decorated.Constant)
		if isConstant {
			continue
		}
		fullyQualifiedName := module.FullyQualifiedName(named.Identifier())
		maybeFunction, _ := unknownType.(*decorated.FunctionValue)
		if maybeFunction != nil {
			if maybeFunction.Annotation().Annotation().IsSomeKindOfExternal() {
				continue
			}
			if verboseFlag >= verbosity.Mid {
				fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
			}

			if maybeFunction.Annotation().Annotation().IsSomeKindOfExternal() {
				continue
			}
			genFuncErr := generateFunction(fullyQualifiedName, maybeFunction, writer, "return ", 0, verboseFlag)
			if genFuncErr != nil {
				return genFuncErr
			}
		} else {
			maybeConstant, _ := unknownType.(*decorated.Constant)
			if maybeConstant != nil {
				if verboseFlag >= verbosity.Mid {
					fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
				}
			} else {
				return fmt.Errorf("generate: unknown type %T", unknownType)
			}
		}
	}

	return nil
}
