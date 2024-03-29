/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
)

func parseCustomTypeVariantTypesUntilNewline(p ParseStream, keywordIndentation int, typeParameterContext *ast.TypeParameterIdentifierContext) ([]ast.Type, parerr.ParseError) {
	var customTypeVariantTypes []ast.Type
	for {
		foundSomething, wasNewLine := p.detectNewLineOrSpace()
		if foundSomething {
			if wasNewLine {
				return customTypeVariantTypes, nil
			}
			p.eatOneSpace("")
		}

		if p.detectAssign() {
			return customTypeVariantTypes, nil
		}
		p.eatOneSpace("")

		userType, userTypeErr := parseTypeVariantParameter(p, keywordIndentation, typeParameterContext)
		if userTypeErr != nil {
			return nil, userTypeErr
		}

		customTypeVariantTypes = append(customTypeVariantTypes, userType)
	}
}
