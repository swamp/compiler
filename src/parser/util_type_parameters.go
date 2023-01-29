/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func readOptionalTypeParameters(p ParseStream, keywordIndentation int, typeParameterContext ast.LocalTypeNameDefinitionContextDynamic) ([]ast.Type, parerr.ParseError) {
	var parameterTypes []ast.Type

	for {
		wasTermination := p.detectOneSpaceAndTermination()
		if wasTermination {
			return parameterTypes, nil
		}
		_, eatErr := p.eatOneSpace("before optional type parameters")
		if eatErr != nil {
			return nil, eatErr
		}

		parameterType, parseTypeErr := internalParseTypeTermReference(p, keywordIndentation, typeParameterContext, false, nil)
		if parseTypeErr != nil {
			_, wasOk := parseTypeErr.(parerr.ExpectedTypeReferenceError)
			if !wasOk {
				return nil, parseTypeErr
			}
			break
		}

		parameterTypes = append(parameterTypes, parameterType)
	}

	return parameterTypes, nil
}
