/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/parser"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/verbosity"
)

type RootStatementHandler struct {
	typeRepo          decorated.TypeAddAndReferenceMaker
	localComments     []decorated.Comment
	localCommentBlock *ast.MultilineComment
	verboseFlag       verbosity.Verbosity
	decorateStream    DecorateStream
	parentModuletype  decorated.ModuleType
}

func NewRootStatementHandler(dectorateStream DecorateStream, typeRepo decorated.TypeAddAndReferenceMaker, parentModuletype decorated.ModuleType, debugName string) *RootStatementHandler {
	g := &RootStatementHandler{verboseFlag: verbosity.None, decorateStream: dectorateStream, typeRepo: typeRepo, parentModuletype: parentModuletype}
	return g
}

func (g *RootStatementHandler) handleAliasStatement(astAlias *ast.Alias) (*dectype.Alias, decshared.DecoratedError) {
	referencedType, referencedTypeErr := g.findTypeFromAstType(astAlias)
	if referencedTypeErr != nil {
		return nil, decorated.NewUnknownTypeAliasType(astAlias, referencedTypeErr)
	}
	alias, wasAlias := referencedType.(*dectype.Alias)
	if !wasAlias {
		return nil, decorated.NewInternalError(fmt.Errorf("was supposed to be an alias here"))
	}
	// g.typeRepo.AddTypeAlias(alias)
	return alias, nil
}

func (g *RootStatementHandler) findTypeFromAstType(constantType ast.Type) (dtype.Type, decorated.TypeError) {
	t, tErr := ConvertFromAstToDecorated(constantType, g.typeRepo)
	if tErr != nil {
		return nil, tErr
	}
	return t, nil
}

func ConvertWrappedOrNormalCustomTypeStatement(hopefullyCustomType ast.Type, typeRepo decorated.TypeAddAndReferenceMaker, localComments []ast.LocalComment) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	customType2, _ := hopefullyCustomType.(*ast.CustomType)
	resultType, tErr := DecorateCustomType(customType2, typeRepo)
	if tErr != nil {
		return nil, tErr
	}
	return resultType, nil
}

func (g *RootStatementHandler) handleCustomTypeStatement(customTypeStatement *ast.CustomType) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	customType, convertErr := ConvertWrappedOrNormalCustomTypeStatement(customTypeStatement, g.typeRepo, nil)
	g.localComments = nil
	if convertErr != nil {
		return nil, convertErr
	}
	typeErr := g.typeRepo.AddCustomType(customType)
	if typeErr != nil {
		panic(typeErr)
	}
	return customType, nil
}

func (g *RootStatementHandler) handleImport(d DecorateStream, importAst *ast.Import) (*decorated.ImportStatement, decshared.DecoratedError) {
	alias := dectype.MakeSingleModuleName(nil)
	if importAst.Alias() != nil {
		alias = dectype.MakeSingleModuleName(importAst.Alias())
	}
	packageRelative := dectype.MakePackageRelativeModuleName(importAst.ModuleName())
	return d.ImportModule(g.parentModuletype, importAst, packageRelative, alias, importAst.ExposeAll(), g.verboseFlag)
}

func (g *RootStatementHandler) handleSinglelineComment(d DecorateStream, singleLineComment *ast.SingleLineComment) (*decorated.SingleLineComment, decshared.DecoratedError) {
	decoratedComment := decorated.NewSingleLineComment(singleLineComment)
	g.localComments = append(g.localComments, decoratedComment)
	return decoratedComment, nil
}

func (g *RootStatementHandler) handleMultilineComment(d DecorateStream, multilineComment *ast.MultilineComment) (*decorated.MultilineComment, decshared.DecoratedError) {
	decoratedComment := decorated.NewMultilineComment(multilineComment)
	g.localComments = append(g.localComments, decoratedComment)
	g.localCommentBlock = multilineComment
	return decoratedComment, nil
}

func (g *RootStatementHandler) handleConstantDefinition(d DecorateStream, constant *ast.ConstantDefinition) (*decorated.Constant, decshared.DecoratedError) {
	name := constant.Identifier()
	variableContext := d.NewVariableContext()
	namedConstant, decoratedExpressionErr := decorateConstant(d, name, constant, variableContext, constant.Comment())
	d.AddDefinition(name, namedConstant)
	if decoratedExpressionErr != nil {
		return nil, decoratedExpressionErr
	}
	g.localComments = nil

	return namedConstant, nil
}

func createParameterDefinitions(referenceMaker decorated.TypeAddAndReferenceMaker, forcedFunctionType *dectype.FunctionAtom, functionValue *ast.FunctionValue) ([]*decorated.FunctionParameterDefinition, decshared.DecoratedError) {
	var parameters []*decorated.FunctionParameterDefinition
	functionParameterTypes, _ := forcedFunctionType.ParameterAndReturn()
	identifiers := functionValue.Parameters()
	if len(identifiers) < len(functionParameterTypes) {
		return nil, decorated.NewTooFewIdentifiersForFunctionType(identifiers, forcedFunctionType, functionValue)
	} else if len(identifiers) > len(functionParameterTypes) {
		return nil, decorated.NewTooManyIdentifiersForFunctionType(identifiers, forcedFunctionType, functionValue)
	}
	for index, parameterType := range functionParameterTypes {
		identifier := identifiers[index]
		argDef := decorated.NewFunctionParameterDefinition(identifier, parameterType)
		parameters = append(parameters, argDef)
	}
	return parameters, nil
}

func (g *RootStatementHandler) declareNamedFunctionValue(d DecorateStream, assignment *ast.FunctionValueNamedDefinition) (*decorated.NamedFunctionValue, decshared.DecoratedError) {

	/*
		forcedFunctionTypeLike := targetFunctionValue.ForcedFunctionType()

		potentialFunc := targetFunctionValue.AstFunctionValue()

		if err := checkParameterCount(forcedFunctionType, potentialFunc); err != nil {
			return err
		}

		_, expectedReturnType := forcedFunctionType.ParameterAndReturn()

		if len(functionParameterTypes) == 0 {
			log.Printf("no input parameters, is this a constant? %v", targetFunctionValue.AstFunctionValue())
		}

	*/

	annotatedType, convertedTypeErr := g.findTypeFromAstType(assignment.FunctionValue().Type())
	if convertedTypeErr != nil {
		return nil, decorated.NewUnknownAnnotationTypeReference(assignment.Identifier(), convertedTypeErr)
	}

	name := assignment.Identifier()
	if annotatedType == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("can not have nil in local annotation"))
	}

	foundFunctionType := DerefFunctionType(annotatedType)
	if foundFunctionType == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("expected function type"))
	}
	forcedFunctionType := DerefFunctionTypeLike(foundFunctionType)
	if forcedFunctionType == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("no forced function type %v", assignment))
	}

	parameters, parametersErr := createParameterDefinitions(d.TypeReferenceMaker(), forcedFunctionType, assignment.FunctionValue())
	if parametersErr != nil {
		return nil, parametersErr
	}

	preparedFunctionValue := decorated.NewPrepareFunctionValue(assignment.FunctionValue(), foundFunctionType, parameters, g.localCommentBlock)

	d.AddDefinition(name, preparedFunctionValue)

	namedFunctionValue := decorated.NewNamedFunctionValue(name, preparedFunctionValue)

	return namedFunctionValue, nil
}

func (g *RootStatementHandler) defineNamedFunctionValue(d DecorateStream, target *decorated.FunctionValue) decshared.DecoratedError {
	variableContext := d.NewVariableContext()

	decoratedExpressionErr := DefineExpressionInPreparedFunctionValue(d, target, variableContext)
	if decoratedExpressionErr != nil {
		return decoratedExpressionErr
	}

	return nil
}

func (g *RootStatementHandler) convertStatement(statement ast.Expression) (decorated.Statement, decshared.DecoratedError) {
	switch v := statement.(type) {
	case *ast.Alias:
		return g.handleAliasStatement(v)
	case *ast.CustomType:
		return g.handleCustomTypeStatement(v)
	case *ast.FunctionValueNamedDefinition:
		return g.declareNamedFunctionValue(g.decorateStream, v)
	case *ast.MultilineComment:
		return g.handleMultilineComment(g.decorateStream, v)
	case *ast.SingleLineComment:
		return g.handleSinglelineComment(g.decorateStream, v)
	case *ast.Import:
		return g.handleImport(g.decorateStream, v)
	case *ast.ConstantDefinition:
		return g.handleConstantDefinition(g.decorateStream, v)
	default:
		return nil, decorated.NewUnknownStatement(token.SourceFileReference{}, statement)
	}
}

func (g *RootStatementHandler) compileFunctionExpression(preparedFunctionNamedValue *decorated.NamedFunctionValue) decshared.DecoratedError {
	preparedFunctionValue := preparedFunctionNamedValue.Value()
	return g.defineNamedFunctionValue(g.decorateStream, preparedFunctionValue)
}

func (g *RootStatementHandler) convertStatements(program *ast.SourceFile) ([]decorated.TypeOrToken, decshared.DecoratedError) {
	var rootNodes []decorated.TypeOrToken

	var errors decshared.DecoratedError
	for _, statement := range program.Statements() {
		convertedStatement, err := g.convertStatement(statement)
		if err != nil {
			if parser.IsCompileErr(err) {
				return nil, err
			}
			errors = decorated.AppendError(errors, err)
		}

		if convertedStatement != nil && !reflect.ValueOf(convertedStatement).IsNil() {
			rootNodes = append(rootNodes, convertedStatement)
		} else {
			panic("not allowed")
		}
	}

	for _, statement := range rootNodes {
		if v, ok := statement.(*decorated.NamedFunctionValue); ok {
			if err := g.compileFunctionExpression(v); err != nil {
				if parser.IsCompileErr(err) {
					return nil, err
				}
				errors = decorated.AppendError(errors, err)
			}
		}
	}

	return rootNodes, errors
}

func (g *RootStatementHandler) HandleStatements(program *ast.SourceFile) ([]decorated.TypeOrToken, decshared.DecoratedError) {
	return g.convertStatements(program)
}
