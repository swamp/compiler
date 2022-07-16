package generate_ir

import (
	"fmt"
	"github.com/llir/llvm/ir"
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

type IrTypeRepo struct {
	string    *types.PointerType
	blob      *types.PointerType
	array     *types.PointerType
	list      *types.PointerType
	unmanaged *types.PointerType

	typeDefs map[string]types.Type
}

func NewIrTypeRepo() *IrTypeRepo {
	stringStruct := types.NewStruct()
	stringStruct.SetName("String")
	stringPointer := types.NewPointer(stringStruct)

	blobStruct := types.NewStruct()
	blobStruct.SetName("Blob")
	blobPointer := types.NewPointer(blobStruct)

	listStruct := types.NewStruct()
	listStruct.SetName("List")
	listPointer := types.NewPointer(listStruct)

	arrayStruct := types.NewStruct()
	arrayStruct.SetName("Array")
	arrayPointer := types.NewPointer(arrayStruct)

	unmanagedStruct := types.NewStruct()
	unmanagedStruct.SetName("Unmanaged")
	unmanagedPointer := types.NewPointer(unmanagedStruct)

	return &IrTypeRepo{
		string:    stringPointer,
		blob:      blobPointer,
		list:      listPointer,
		array:     arrayPointer,
		unmanaged: unmanagedPointer,
		typeDefs:  make(map[string]types.Type),
	}
}

func (r *IrTypeRepo) AddTypeDef(name string, newType types.Type) {
	_, hasType := r.typeDefs[name]
	if hasType {
		log.Printf("skipping %v", name)
		return
	}

	r.typeDefs[name] = newType
}

func (r *IrTypeRepo) GetTypeRef(name string) (types.Type, error) {
	log.Printf("can not find typeref '%v'", name)
	foundType, hasType := r.typeDefs[name]
	if !hasType {
		return nil, fmt.Errorf("can not find %v %v", name, r.typeDefs)
	}

	return foundType, nil
}

func makeIrForType(irModule *ir.Module, repo *IrTypeRepo, p dtype.Type) types.Type {
	unaliased := dectype.UnaliasWithResolveInvoker(p)
	switch t := unaliased.(type) {
	case *dectype.RecordAtom:
		return generateRecordType(irModule, repo, t)
	case *dectype.TupleTypeAtom:
		return generateTupleType(irModule, repo, t)
	case *dectype.FunctionAtom:
		return generateFunctionType(irModule, repo, t)
	case *dectype.CustomTypeAtom:
		foundType, foundErr := repo.GetTypeRef(t.ArtifactTypeName().String())
		if foundErr != nil {
			generateCustomType(irModule, repo, t)
			foundType2, foundErr2 := repo.GetTypeRef(t.ArtifactTypeName().String())
			if foundErr2 != nil {
				panic(foundErr2)
			}
			return foundType2
		}
		return foundType
	case *dectype.PrimitiveAtom:
		switch t.AtomName() {
		case "Int":
			return types.I32
		case "Fixed":
			return types.I32
		case "Bool":
			return types.I8
		case "String":
			return repo.string
		case "ResourceName":
			return repo.string
		case "Blob":
			return repo.blob
		case "Array":
			return repo.array
		case "List":
			return repo.list
		default:
			panic(fmt.Errorf("unknown atom %v", t))
		}
	case *dectype.UnmanagedType:
		foundType, err := repo.GetTypeRef(t.Identifier().NativeLanguageTypeName().Name())
		if err != nil {

			panic(fmt.Errorf("what is this %v", err))
		}
		return foundType
	default:
		panic(fmt.Errorf("what is this %T", t))
		return types.I1
	}
}

func makeIrType(irModule *ir.Module, repo *IrTypeRepo, p dtype.Type) types.Type {
	return makeIrForType(irModule, repo, p)
}

func generateFunctionParameter(irModule *ir.Module, repo *IrTypeRepo, functionParam *decorated.FunctionParameterDefinition) *ir.Param {
	irType := makeIrType(irModule, repo, functionParam.Type())
	newParam := ir.NewParam(functionParam.Identifier().Name(), irType)

	return newParam
}

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName,
	f *decorated.FunctionValue, lookup typeinfo.TypeLookup, resourceNameLookup resourceid.ResourceNameLookup, fileCache *assembler_sp.FileUrlCache, irModule *ir.Module, repo *IrTypeRepo, verboseFlag verbosity.Verbosity) (*ir.Func, error) {
	functionType := f.Type().(*dectype.FunctionTypeReference).FunctionAtom()
	irReturnType := makeIrType(irModule, repo, functionType.ReturnType())
	//unaliasedReturnType := dectype.UnaliasWithResolveInvoker()

	var irParams []*ir.Param
	for _, parameter := range f.Parameters() {
		irParam := generateFunctionParameter(irModule, repo, parameter)
		irParams = append(irParams, irParam)
		//		log.Println(irParam)
	}

	paramContext := newParameterContext(irParams)

	log.Printf("paramContext %v", paramContext)

	genContext := &generateContext{
		irModule:           irModule,
		block:              ir.NewBlock("function"),
		parameterContext:   paramContext,
		lookup:             lookup,
		resourceNameLookup: resourceNameLookup,
		fileCache:          fileCache,
	}

	newIrFunc := irModule.NewFunc(f.Annotation().Annotation().Identifier().Name(), irReturnType, irParams...)

	result, genErr := generateExpression(f.Expression(), genContext)
	if genErr != nil {
		return nil, genErr
	}
	if result == nil {
		return nil, nil
	}

	genContext.block.Term = genContext.block.NewRet(result)

	log.Printf("result of function was %v in block %v", result.Ident(), genContext.block.LLString())

	return newIrFunc, nil
}
