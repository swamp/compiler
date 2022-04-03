package generate_ir

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/swamp/assembler/lib/assembler_sp"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
)

func makeIrType(p dtype.Type) types.Type {
	switch t := p.(type) {
	case *dectype.PrimitiveAtom:
		if t.AtomName() == "Int" {
			return types.I32
		}
		return types.I1
	default:
		return types.I1
	}
}

func generateFunctionParameter(functionParam *decorated.FunctionParameterDefinition) *ir.Param {
	irType := makeIrType(functionParam.Type())
	newParam := ir.NewParam("x", irType)

	return newParam
}

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName,
	f *decorated.FunctionValue, lookup typeinfo.TypeLookup, resourceNameLookup resourceid.ResourceNameLookup, fileCache *assembler_sp.FileUrlCache, module *ir.Module, verboseFlag verbosity.Verbosity) (*ir.Func, error) {
	functionType := f.Type().(*dectype.FunctionTypeReference).FunctionAtom()
	irReturnType := makeIrType(functionType.ReturnType())
	//unaliasedReturnType := dectype.UnaliasWithResolveInvoker()

	var irParams []*ir.Param
	for _, parameter := range f.Parameters() {
		irParam := generateFunctionParameter(parameter)
		irParams = append(irParams, irParam)
		//		log.Println(irParam)
	}

	newIrFunc := module.NewFunc(f.Annotation().Annotation().Identifier().Name(), irReturnType, irParams...)

	/*
		genErr := generateExpression(code, returnValueTargetPointer, f.Expression(), true, genContext)
		if genErr != nil {
			return nil, genErr
		}
	*/

	return newIrFunc, nil
}