/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
)

type DecorateStream interface {
	TypeRepo() *dectype.TypeRepo
	AddDefinition(identifier *ast.VariableIdentifier, expr decorated.DecoratedExpression) error
	AddDeclaration(identifier *ast.VariableIdentifier, declaredType dtype.Type) error
	NewVariableContext() *VariableContext
	AddImport(moduleName dectype.PackageRelativeModuleName, alias dectype.SingleModuleName, exposeAll bool, verboseFlag bool) decshared.DecoratedError
	AddExternalFunction(functionName string, parameterCount uint) decshared.DecoratedError
}
