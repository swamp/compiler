package generate_sp

/*
func Pack(functions []*Function, externalFunctions []*ExternalFunction, typeInfoPayload []byte,
	lookup typeinfo.TypeLookup) ([]byte, error) {
	constantRepo := swamppack.NewConstantRepo()

	for _, externalFunction := range externalFunctions {
		constantRepo.AddExternalFunction(externalFunction.name.ResolveToString(), externalFunction.parameterCount)
	}

	for _, declareFunction := range functions {
		constantRepo.AddFunctionDeclaration(declareFunction.name.ResolveToString(),
			declareFunction.signature, declareFunction.parameterCount)
	}

	for _, function := range functions {
		var packConstants []*swamppack.Constant

		for _, subConstant := range function.constants {
			var packConstant *swamppack.Constant

			switch subConstant.ConstantType() {
			case assembler_sp.ConstantTypeInteger:
				packConstant = constantRepo.AddInteger(subConstant.IntegerValue())
			case assembler_sp.ConstantTypeResourceName:
				packConstant = constantRepo.AddResourceName(subConstant.StringValue())
			case assembler_sp.ConstantTypeString:
				packConstant = constantRepo.AddString(subConstant.StringValue())
			case assembler_sp.ConstantTypeBoolean:
				packConstant = constantRepo.AddBoolean(subConstant.BooleanValue())
			case assembler_sp.ConstantTypeFunction:
				refConstant, functionRefErr := constantRepo.AddFunctionReference(
					subConstant.FunctionReferenceFullyQualifiedName())
				if functionRefErr != nil {
					return nil, functionRefErr
				}
				packConstant = refConstant
			case assembler_sp.ConstantTypeFunctionExternal:
				refConstant, functionRefErr := constantRepo.AddExternalFunctionReference(
					subConstant.FunctionReferenceFullyQualifiedName())
				if functionRefErr != nil {
					return nil, functionRefErr
				}

				packConstant = refConstant
			default:
				return nil, fmt.Errorf("not handled constanttype %v", subConstant)
			}

			if packConstant == nil {
				return nil, fmt.Errorf("internal error: not handled constanttype %v", subConstant)
			}

			packConstants = append(packConstants, packConstant)
		}

		constantRepo.AddFunction(function.name.ResolveToString(), function.signature, function.parameterCount,
			function.variableCount, packConstants, function.opcodes)
	}

	octets, packErr := swamppack.Pack(constantRepo, typeInfoPayload)
	if packErr != nil {
		return nil, packErr
	}

	return octets, nil
}


*/
