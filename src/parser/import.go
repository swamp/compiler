/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
)

func parseModuleName(p ParseStream) ([]*ast.TypeIdentifier, parerr.ParseError) {
	var lookups []*ast.TypeIdentifier
	importName, importNameErr := p.readTypeIdentifier()
	if importNameErr != nil {
		return nil, parerr.NewImportMustHaveUppercaseIdentifierError(importNameErr)
	}
	lookups = append(lookups, importName)
	for {
		if !p.maybeAccessor() {
			break
		}
		lookup, subNameErr := p.readTypeIdentifier()
		if subNameErr != nil {
			return nil, parerr.NewImportMustHaveUppercasePathError(subNameErr)
		}
		lookups = append(lookups, lookup)
	}

	return lookups, nil
}

func parseImport(p ParseStream, keyword token.VariableSymbolToken,
	precedingComments token.CommentBlock) (ast.Expression, parerr.ParseError) {
	_, spaceAfterKeywordErr := p.eatOneSpace("space after IMPORT")
	if spaceAfterKeywordErr != nil {
		return nil, spaceAfterKeywordErr
	}

	moduleReference, moduleReferenceErr := parseModuleName(p)
	if moduleReferenceErr != nil {
		return nil, moduleReferenceErr
	}

	var alias *ast.TypeIdentifier

	wasNewLine := p.detectNewLine()
	if !wasNewLine {
		_, spaceAfterModuleReferenceErr := p.eatOneSpace("space after IMPORT")
		if spaceAfterModuleReferenceErr != nil {
			return nil, spaceAfterModuleReferenceErr
		}

		if p.maybeKeywordAs() {
			_, spaceAfterAliasErr := p.eatOneSpace("space after alias")
			if spaceAfterAliasErr != nil {
				return nil, spaceAfterAliasErr
			}

			foundAlias, aliasErr := p.readTypeIdentifier()
			if aliasErr != nil {
				return nil, aliasErr
			}
			alias = foundAlias

			if alias.Name() != moduleReference[len(moduleReference)-1].Name() {
				p.addWarning("it is advised to use the last part of the import as alias. `import Some.Long.Name as Name`", p.positionLength())
			}

			wasNewLine := p.detectNewLine()
			if !wasNewLine {
				return nil, parerr.NewExpectedContinuationLineOrOneSpace(p.positionLength())
			}
		} else {
			return nil, parerr.NewExpectedContinuationLineOrOneSpace(p.positionLength())
		}
	}

	return ast.NewImport(keyword, moduleReference, alias, precedingComments), nil
}
