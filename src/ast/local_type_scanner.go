/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

func containsLocalType(params []*TypeParameter, newType *TypeParameter) bool {
	for _, param := range params {
		if param.Identifier().Name() == newType.Identifier().Name() {
			return true
		}
	}
	return false
}

func scanLocalTypes(types []Type, params *[]*TypeParameter) {
	for _, foundType := range types {
		switch e := foundType.(type) {
		case *FunctionType:
			scanLocalTypes(e.FunctionParameters(), params)
		case *TypeReference:
			scanLocalTypes(e.Arguments(), params)
		case *LocalType:
			if !containsLocalType(*params, e.TypeParameter()) {
				*params = append(*params, e.TypeParameter())
			}
		}
	}
}

func ScanLocalTypeInTypes(types []Type) []*TypeParameter {
	params := make([]*TypeParameter, 0)

	scanLocalTypes(types, &params)

	return params
}
