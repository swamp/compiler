/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	dectype "github.com/swamp/compiler/src/decorated/types"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/dtype"
	"github.com/swamp/compiler/src/token"
)

type ExternalFunctionDeclarationExpression struct {
	astFunctionDeclaration *ast.FunctionDeclarationExpression
	functionType           dtype.Type `debug:"true"`
}

func NewExternalFunctionDeclarationExpression(astFunctionDeclaration *ast.FunctionDeclarationExpression) *ExternalFunctionDeclarationExpression {
	return &ExternalFunctionDeclarationExpression{
		astFunctionDeclaration: astFunctionDeclaration, functionType: dectype.NewAnyType(),
	}
}

func (i *ExternalFunctionDeclarationExpression) Type() dtype.Type {
	return i.functionType
}

func (i *ExternalFunctionDeclarationExpression) String() string {
	return fmt.Sprintf("[ExternalFunctionDeclarationExpression %v %v]", i.astFunctionDeclaration, i.functionType)
}

func (i *ExternalFunctionDeclarationExpression) FetchPositionLength() token.SourceFileReference {
	return i.astFunctionDeclaration.FetchPositionLength()
}
