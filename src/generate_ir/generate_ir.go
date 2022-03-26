package generate_ir

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
	"github.com/swamp/assembler/lib/assembler_sp"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	"log"
)

// This example produces LLVM IR code equivalent to the following C code, which
// implements a pseudo-random number generator.
//
//    int abs(int x);
//
//    // ref: https://en.wikipedia.org/wiki/Linear_congruential_generator
//    //    a = 0x15A4E35
//    //    c = 1
//    int rand(void) {
//       seed = seed*0x15A4E35 + 1;
//       return abs(seed);
//    }

import (
	"fmt"
)

func Example() {
	// Create convenience types and constants.
	i32 := types.I32
	a := constant.NewInt(i32, 0x15A4E35) // multiplier of the PRNG.
	c := constant.NewInt(i32, 1)         // increment of the PRNG.

	// Create a new LLVM IR module.
	m := ir.NewModule()

	// Create an external function declaration and append it to the module.
	//
	//    int abs(int x);
	abs := m.NewFunc("abs", i32, ir.NewParam("x", i32))
	// Create a function definition and append it to the module.
	//
	//    int rand(void) { ... }

	diCompileUnit := &metadata.DICompileUnit{
		MetadataID:   -1,
		Distinct:     true,
		Language:     enum.DwarfLangC11,
		Producer:     "clang version 12.0.0 (tags/RELEASE_800/final)",
		EmissionKind: enum.EmissionKindFullDebug,
	}

	x := ir.NewParam("x", i32)


	rand := m.NewFunc("rand", i32, x)


	// Create an unnamed entry basic block and append it to the `rand` function.
	block := rand.NewBlock("")

	// Create instructions and append them to the basic block.
	tmp1 := block.NewLoad(i32, x)

	diFile := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "foo.c",
		Directory:  "/home/u/Desktop/foo",
	}

	diCompileUnit.File = diFile



	tmp2 := block.NewMul(tmp1, a)
	tmp3 := block.NewAdd(tmp2, c)

	diBasicTypeI32 := &metadata.DIBasicType{
		MetadataID: -1,
		Name:       "int",
		Size:       32,
		Encoding:   enum.DwarfAttEncodingSigned,
	}

	diLocalVarA := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "a",
		Arg:        1,
		Scope:      tmp3,
		File:       diFile,
		Line:       1,
		Type:       diBasicTypeI32,
	}
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarA)

	tmp4 := block.NewCall(abs, tmp3)
	block.NewRet(tmp4)

	rand2 := m.NewFunc("rand2", i32, x)
	block2 := rand2.NewBlock("")
	tmp11 := block.NewLoad(i32, x)
	block2.NewRet(tmp11)
	// Print the LLVM IR assembly of the module.
	fmt.Println(m)
}

type Generator struct {

}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateAllLocalDefinedFunctions(module *decorated.Module,
	lookup typeinfo.TypeLookup, resourceNameLookup resourceid.ResourceNameLookup, fileUrlCache *assembler_sp.FileUrlCache, verboseFlag verbosity.Verbosity) (*ir.Module, error) {
	irModule := ir.NewModule()

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
			_, genFuncErr := generateFunction(fullyQualifiedName, maybeFunction,
				lookup, resourceNameLookup, fileUrlCache, irModule, verboseFlag)
			if genFuncErr != nil {
				return nil, genFuncErr
			}

			if verboseFlag >= verbosity.High {
				log.Printf("---------- generated code for '%v'", fullyQualifiedName.String())
			}

		} else {
			maybeConstant, _ := unknownType.(*decorated.Constant)
			if maybeConstant != nil {
				if verboseFlag >= verbosity.Mid {
					fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
				}
			} else {
				return nil, fmt.Errorf("generate: unknown type %T", unknownType)
			}
		}
	}

	return irModule, nil
}
