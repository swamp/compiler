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

func parseMultiLineComment(p ParseStream, commentToken token.MultiLineCommentToken) (*ast.MultilineComment, parerr.ParseError) {
	c := ast.NewMultilineComment(commentToken)
	return c, nil
}
