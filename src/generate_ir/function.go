/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

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

type IrFunctions struct {
	functions map[string]*ir.Func
}

func NewIrFunctions() *IrFunctions {
	return &IrFunctions{functions: make(map[string]*ir.Func)}
}

func (i *IrFunctions) AddFunc(name *decorated.FullyQualifiedPackageVariableName, p *ir.Func) {
	i.functions[name.ResolveToString()] = p
}

func (i *IrFunctions) GetFunc(name *decorated.FullyQualifiedPackageVariableName) *ir.Func {
	return i.functions[name.ResolveToString()]
}

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

func (r *IrTypeRepo) AddTypeDef(decoratedType dtype.Type, newType types.Type) {
	unreferenced := dectype.UnReference(decoratedType)
	_, hasType := r.typeDefs[unreferenced.String()]
	if hasType {
		log.Printf("skipping %v", decoratedType)
		return
	}
	typeName := decoratedType.String()
	switch t := unreferenced.(type) {
	case *dectype.CustomTypeAtom:
		typeName = t.ArtifactTypeName().String()
	}

	log.Printf("**** [%v] = %T", typeName, newType)
	r.typeDefs[typeName] = newType
}

func (r *IrTypeRepo) GetTypeRef(decoratedType dtype.Type) (types.Type, error) {
	unreferenced := dectype.UnReference(decoratedType)
	unaliased := dectype.UnaliasWithResolveInvoker(decoratedType)
	switch t := unaliased.(type) {
	case *dectype.RecordAtom:
		{
			return generateRecordType(nil, r, t), nil
		}
	case *dectype.PrimitiveAtom:
		switch t.AtomName() {
		case "Int":
			return types.I32, nil
		case "Fixed":
			return types.I32, nil
		case "Bool":
			return types.I8, nil
		case "String":
			return r.string, nil
		case "ResourceName":
			return r.string, nil
		case "Blob":
			return r.blob, nil
		case "Array":
			return r.array, nil
		case "List":
			return r.list, nil
		default:
			panic(fmt.Errorf("unknown primitive atom %v", t))
		}
	case *dectype.CustomTypeAtom:
		typeName := t.ArtifactTypeName().String()
		foundType, hasType := r.typeDefs[typeName]
		if !hasType {
			panic(fmt.Errorf("GetTypeRef: can not CustomTypeAtom '%v' '%v' '%v'", typeName, unreferenced, r.typeDefs))
			return nil, fmt.Errorf("can not find %v %v", unreferenced, r.typeDefs)
		}
		return foundType, nil
	default:

		foundType, hasType := r.typeDefs[unreferenced.String()]
		if !hasType {
			panic(fmt.Errorf("GetTypeRef: can not find '%v' '%v'", unreferenced, r.typeDefs))
			return nil, fmt.Errorf("can not find %v %v", unreferenced, r.typeDefs)
		}

		return foundType, nil
	}
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
		foundType, foundErr := repo.GetTypeRef(t)
		if foundErr != nil {
			generateCustomType(irModule, repo, t)
			foundType2, foundErr2 := repo.GetTypeRef(t)
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
			panic(fmt.Errorf("unknown primitive atom %v", t))
		}
	case *dectype.UnmanagedType:
		voidPointer := types.NewPointer(types.Void)
		return voidPointer
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
	if types.IsStruct(irType) {
		irType = types.NewPointer(irType)
	}
	newParam := ir.NewParam(functionParam.Parameter().Name(), irType)

	return newParam
}

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName,
	f *decorated.FunctionValue, lookup typeinfo.TypeLookup, resourceNameLookup resourceid.ResourceNameLookup, fileCache *assembler_sp.FileUrlCache, irModule *ir.Module, repo *IrTypeRepo, irFunctions *IrFunctions, verboseFlag verbosity.Verbosity) (*ir.Func, error) {
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
		irTypeRepo:         repo,
		lookup:             lookup,
		resourceNameLookup: resourceNameLookup,
		fileCache:          fileCache,
		irFunctions:        irFunctions,
	}

	newIrFunc := irModule.NewFunc(fullyQualifiedVariableName.Identifier().Name(), irReturnType, irParams...)
	irFunctions.AddFunc(fullyQualifiedVariableName, newIrFunc)

	result, genErr := generateExpression(f.Expression(), true, genContext)
	if genErr != nil {
		return nil, genErr
	}
	if result == nil {
		return nil, nil
	}

	genContext.block.Term = genContext.block.NewRet(result)

	log.Printf("result of function was %v in block %v", result.Ident(), genContext.block.LLString())

	newIrFunc.Blocks = append(newIrFunc.Blocks, genContext.block)
	return newIrFunc, nil
}
