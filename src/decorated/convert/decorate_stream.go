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
	"github.com/swamp/compiler/src/verbosity"
)

type DecorateStream interface {
	TypeRepo() decorated.TypeAddAndReferenceMaker
	AddDefinition(identifier *ast.VariableIdentifier, expr decorated.Expression) error
	AddDeclaration(identifier *ast.VariableIdentifier, declaredType dtype.Type) error
	NewVariableContext() *VariableContext
	ImportModule(moduleType decorated.ModuleType, importAst *ast.Import, moduleName dectype.PackageRelativeModuleName, alias dectype.SingleModuleName, exposeAll bool, verboseFlag verbosity.Verbosity) (*decorated.ImportStatement, decshared.DecoratedError)
	AddExternalFunction(function *ast.ExternalFunction) (*decorated.ExternalFunctionDeclaration, decshared.DecoratedError)
	FindNamedFunctionValue(identifier *ast.VariableIdentifier) *decorated.FunctionValue
	AddDecoratedError(decoratedError decshared.DecoratedError)
}
