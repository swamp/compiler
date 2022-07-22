package generate_ir

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
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

func (g *Generator) GenerateAllLocalDefinedFunctions(module *decorated.Module, irModule *ir.Module, repo *IrTypeRepo,
	lookup typeinfo.TypeLookup, resourceNameLookup resourceid.ResourceNameLookup, fileUrlCache *assembler_sp.FileUrlCache, verboseFlag verbosity.Verbosity) (*ir.Module, error) {
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
				lookup, resourceNameLookup, fileUrlCache, irModule, repo, verboseFlag)
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

func generateRecordType(irModule *ir.Module, repo *IrTypeRepo, recordType *dectype.RecordAtom) *types.StructType {
	var recordFieldIrTypes []types.Type

	for _, recordField := range recordType.SortedFields() {
		recordFieldIrType := makeIrType(irModule, repo, recordField.Type())
		recordFieldIrTypes = append(recordFieldIrTypes, recordFieldIrType)
	}
	recordStruct := types.NewStruct(recordFieldIrTypes...)
	// Note: Not allowed to set a name for the struct. We need a literal structure that is compared by the contents and not the typename

	return recordStruct
}

func generateTupleType(irModule *ir.Module, repo *IrTypeRepo, tupleType *dectype.TupleTypeAtom) *types.StructType {
	var tupleFieldIrTypes []types.Type

	for _, tupleField := range tupleType.Fields() {
		tupleFieldIrType := makeIrType(irModule, repo, tupleField.Type())
		tupleFieldIrTypes = append(tupleFieldIrTypes, tupleFieldIrType)
	}
	tupleStruct := types.NewStruct(tupleFieldIrTypes...)
	// Note: Not allowed to set a name for the struct. We need a literal structure that is compared by the contents and not the typename

	return tupleStruct
}

func generateFunctionType(irModule *ir.Module, repo *IrTypeRepo, functionType *dectype.FunctionAtom) *types.StructType {
	var functionFieldIrTypes []types.Type

	for _, functionField := range functionType.FunctionParameterTypes() {
		functionFieldIrType := makeIrType(irModule, repo, functionField)
		functionFieldIrTypes = append(functionFieldIrTypes, functionFieldIrType)
	}
	functionStruct := types.NewStruct(functionFieldIrTypes...)
	// Note: Not allowed to set a name for the struct. We need a literal structure that is compared by the contents and not the typename

	log.Printf("   function %v", functionStruct)

	return functionStruct
}

// generateCustomType generates Ir types for a swamp custom type.
// A custom type in swamp is in principle the same as a tagged union.
// https://mapping-high-level-constructs-to-llvm-ir.readthedocs.io/en/latest/basic-constructs/unions.html#tagged-unions
// Tagged union has a single octet and then a array of octets with the max size of the struct (including padding)
// Each variant needs to bitcast from the completeUnion to the specific variant:
// Example %1 = bitcast %CustomTypeName* %x to %CustomTypeName_VariantName*
func generateCustomType(irModule *ir.Module, repo *IrTypeRepo, customType *dectype.CustomTypeAtom) error {

	memSize, _ := dectype.GetMemorySizeAndAlignment(customType)
	maximumPaddedSize := memSize - 1
	unionPayloadArray := types.NewArray(uint64(maximumPaddedSize), types.I8)
	completeUnionStruct := types.NewStruct(types.I8, unionPayloadArray)
	irCompleteUnionName := customType.ArtifactTypeName().String()
	completeUnionStruct.SetName(irCompleteUnionName)
	completeUnionTypeDef := irModule.NewTypeDef(irCompleteUnionName, completeUnionStruct)

	repo.AddTypeDef(customType, completeUnionTypeDef)

	for _, variant := range customType.Variants() {
		var variantParamIrTypes []types.Type
		variantParamIrTypes = append(variantParamIrTypes, types.I8)
		for _, variantParam := range variant.ParameterTypes() {
			variantParamIrType := makeIrType(irModule, repo, variantParam)
			variantParamIrTypes = append(variantParamIrTypes, variantParamIrType)
		}
		ilVariantName := customType.DecoratedName() + "_" + variant.DecoratedName()
		variantStruct := types.NewStruct(variantParamIrTypes...)
		variantStruct.SetName(ilVariantName)
		irModule.NewTypeDef(ilVariantName, variantStruct)
	}

	return nil
}

func generateAlias(irModule *ir.Module, repo *IrTypeRepo, alias *dectype.Alias) error {
	log.Printf("alias: %T", alias.Next())
	switch t := alias.Next().(type) {
	case *dectype.RecordAtom:
		irType := generateRecordType(irModule, repo, t)
		log.Printf("irType:%v", irType)
	}
	return nil
}

func generateType(irModule *ir.Module, repo *IrTypeRepo, definedType dtype.Type) error {
	unAliased := dectype.UnaliasWithResolveInvoker(definedType)
	switch t := unAliased.(type) {
	case *dectype.CustomTypeAtom:
		return generateCustomType(irModule, repo, t)
	case *dectype.CustomTypeVariant:
		return nil // All variants are generated along side customType
	case *dectype.RecordAtom:
		return nil
	case *dectype.PrimitiveAtom:
		return nil
	case *dectype.UnmanagedType:
		log.Printf("unmanagedType %v", t.Identifier())
	default:
		log.Printf("what is this %T", unAliased)
	}

	return nil
}

func (g *Generator) GenerateModule(module *decorated.Module,
	lookup typeinfo.TypeLookup, resourceNameLookup resourceid.ResourceNameLookup, fileUrlCache *assembler_sp.FileUrlCache, verboseFlag verbosity.Verbosity) (*ir.Module, error) {
	irModule := ir.NewModule()

	repo := NewIrTypeRepo()

	for _, definedType := range module.LocalTypes().AllTypes() {
		if err := generateType(irModule, repo, definedType); err != nil {
			return nil, err
		}
	}

	g.GenerateAllLocalDefinedFunctions(module, irModule, repo, lookup, resourceNameLookup, fileUrlCache, verboseFlag)

	return irModule, nil
}
