/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type Definer struct {
	typeRepo        *dectype.TypeRepo
	localAnnotation *LocalAnnotation
	localComments   []ast.LocalComment
	verboseFlag     bool
	dectorateStream DecorateStream
}

func NewDefiner(dectorateStream DecorateStream, typeRepo *dectype.TypeRepo, debugName string) *Definer {
	g := &Definer{verboseFlag: false, localAnnotation: nil, dectorateStream: dectorateStream, typeRepo: typeRepo}
	return g
}

func (g *Definer) createAliasTypeFromType(aliasName *ast.TypeIdentifier, subType dtype.Type) (dtype.Type, decshared.DecoratedError) {
	existingType := g.typeRepo.FindTypeFromAlias(aliasName.Name())
	if existingType != nil {
		panic(fmt.Sprintf("type alias already defined %v", aliasName))
	}
	t, typeErr := g.typeRepo.DeclareAlias(aliasName, subType, g.localComments)
	if typeErr != nil {
		panic(typeErr)
	}
	g.localComments = nil
	return t, nil
}

func (g *Definer) AnnotateConstant(identifier *ast.VariableIdentifier, realType dtype.Type) decshared.DecoratedError {
	g.localAnnotation = NewLocalAnnotation(identifier, realType)
	return nil
}

func (g *Definer) AnnotateFunc(identifier *ast.VariableIdentifier, funcType dtype.Type) error {
	g.localAnnotation = NewLocalAnnotation(identifier, funcType)
	return nil
}

func (g *Definer) convertAnnotation(identifier *ast.VariableIdentifier, constantType ast.Type) (dtype.Type, decshared.DecoratedError) {
	if constantType == nil {
		return nil, decorated.NewUnknownAnnotationTypeReference(identifier, nil)
	}
	convertedType, convertedTypeErr := g.findTypeFromAstType(constantType)
	if convertedTypeErr != nil {
		return nil, decorated.NewUnknownAnnotationTypeReference(identifier, convertedTypeErr)
	}

	return convertedType, nil
}

func (g *Definer) functionAnnotation(identifier *ast.VariableIdentifier, constantType ast.Type) (dtype.Type, decshared.DecoratedError) {
	convertedType, decErr := g.convertAnnotation(identifier, constantType)
	if decErr != nil {
		return nil, decErr
	}
	checkedType, _ := convertedType.(dtype.Type)
	if checkedType != nil {
		g.dectorateStream.AddDeclaration(identifier, checkedType)
	}
	g.AnnotateFunc(identifier, convertedType)

	return convertedType, nil
	//	return nil, decorated.NewInternalError(fmt.Errorf("do not know what to do:%v", constantType))
}

func (g *Definer) handleAliasStatement(alias *ast.AliasStatement) decshared.DecoratedError {
	t := alias.Type()
	referencedType, referencedTypeErr := g.findTypeFromAstType(t)
	if referencedTypeErr != nil {
		return decorated.NewUnknownTypeAliasType(alias, referencedTypeErr)
	}
	_, tErr := g.createAliasTypeFromType(alias.TypeIdentifier(), referencedType)
	return tErr
}

func (g *Definer) findTypeFromAstType(constantType ast.Type) (dtype.Type, dectype.DecoratedTypeError) {
	t, tErr := ConvertFromAstToDecorated(constantType, g.typeRepo)
	if tErr != nil {
		return nil, tErr
	}
	return t, nil
}

func (g *Definer) createCustomType(ast *ast.CustomType) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	customType, createErr := DecorateCustomType(ast, g.typeRepo)
	if createErr != nil {
		return nil, createErr
	}

	return customType, nil
}

func ConvertWrappedOrNormalCustomTypeStatement(hopefullyCustomType ast.Type, typeRepo *dectype.TypeRepo, localComments []ast.LocalComment) (dtype.Type, decshared.DecoratedError) {
	customType2, _ := hopefullyCustomType.(*ast.CustomType)
	resultType, tErr := DecorateCustomType(customType2, typeRepo)
	if tErr != nil {
		return nil, tErr
	}
	return resultType, nil
}

func (g *Definer) handleCustomTypeStatement(customTypeStatement *ast.CustomTypeStatement) decshared.DecoratedError {
	hopefullyCustomType := customTypeStatement.Type()
	customType, convertErr := ConvertWrappedOrNormalCustomTypeStatement(hopefullyCustomType, g.typeRepo, g.localComments)
	g.localComments = nil
	if convertErr != nil {
		return convertErr
	}
	typeErr := g.typeRepo.DeclareType(customType)
	if typeErr != nil {
		panic(typeErr)
	}
	return nil
}

func (g *Definer) handleImport(d DecorateStream, importAst *ast.Import) decshared.DecoratedError {
	alias := dectype.MakeSingleModuleName(nil)
	if importAst.Alias() != nil {
		alias = dectype.MakeSingleModuleName(importAst.Alias())
	}
	packageRelative := dectype.MakePackageRelativeModuleName(importAst.ModuleName())
	return d.AddImport(packageRelative, alias, importAst.ExposeAll(), g.verboseFlag)
}

func (g *Definer) handleExternalFunction(d DecorateStream, externalFunction *ast.ExternalFunction) decshared.DecoratedError {
	g.localComments = nil
	return d.AddExternalFunction(externalFunction.ExternalFunction(), externalFunction.ParameterCount())
}

func (g *Definer) handleSinglelineComment(d DecorateStream, singleLineComment *ast.SingleLineComment) decshared.DecoratedError {
	g.localComments = append(g.localComments, ast.LocalComment{Singleline: singleLineComment})
	return nil
}

func (g *Definer) handleMultilineComment(d DecorateStream, multilineComment *ast.MultilineComment) decshared.DecoratedError {
	g.localComments = append(g.localComments, ast.LocalComment{Multiline: multilineComment})
	return nil
}

func (g *Definer) handleDefinitionAssignment(d DecorateStream, assignment *ast.DefinitionAssignment) decshared.DecoratedError {
	if g.localAnnotation == nil {
		return decorated.NewMustHaveAnnotationJustBeforeThisDefinition(assignment)
	}
	if g.localAnnotation.Identifier().Name() != assignment.Identifier().Name() {
		return decorated.NewAnnotationMismatch(g.localAnnotation.Identifier(), assignment)
	}
	name := assignment.Identifier()
	expr := assignment.Expression()
	annotatedType := g.localAnnotation.Type()
	if annotatedType == nil {
		return decorated.NewInternalError(fmt.Errorf("can not have nil in local annotation"))
	}
	g.localAnnotation = nil

	/*
		generator, isGenerator := annotatedType.(dtype.Type)
		if isGenerator {
			g.HackAddGeneratorsCode(d, assignment, generator)

			functionValue, _ := expr.(*ast.FunctionValue)
			if functionValue != nil {
				assembly, isAssembly := functionValue.Expression().(*ast.Asm)

				if isAssembly {
					functionValue, _ := expr.(*ast.FunctionValue)
					return g.HackAddRawCode(d, functionValue, assignment.Identifier(), assembly)
				}
			}
			return nil
		}

	*/

	variableContext := d.NewVariableContext()
	g.localComments = nil
	_, decoratedExpressionErr := decorateDefinition(d, variableContext, name, expr, annotatedType, g.localComments)
	if decoratedExpressionErr != nil {
		return decoratedExpressionErr
	}
	g.localAnnotation = nil
	return nil
}

func (g *Definer) handleAnnotation(declaration *ast.Annotation) decshared.DecoratedError {
	if g.localAnnotation != nil {
		return decorated.NewAlreadyHaveAnnotationForThisName(declaration, nil)
	}
	annotatedType := declaration.AnnotatedType()
	{
		_, declareErr := g.functionAnnotation(declaration.Identifier(), annotatedType)
		if declareErr != nil {
			return declareErr
		}
		return nil
	}
}

func (g *Definer) handleStatement(statement ast.Expression) decshared.DecoratedError {
	switch v := statement.(type) {
	case *ast.AliasStatement:
		return g.handleAliasStatement(v)
	case *ast.CustomTypeStatement:
		return g.handleCustomTypeStatement(v)
	case *ast.Annotation:
		return g.handleAnnotation(v)
	case *ast.DefinitionAssignment:
		return g.handleDefinitionAssignment(g.dectorateStream, v)
	case *ast.Import:
		return g.handleImport(g.dectorateStream, v)
	case *ast.ExternalFunction:
		return g.handleExternalFunction(g.dectorateStream, v)
	case *ast.MultilineComment:
		return g.handleMultilineComment(g.dectorateStream, v)
	case *ast.SingleLineComment:
		return g.handleSinglelineComment(g.dectorateStream, v)
	default:
		return decorated.NewUnknownStatement(token.PositionLength{}, statement)
	}
}

func (g *Definer) firstPass(program *ast.Program) decshared.DecoratedError {
	for _, statement := range program.Statements() {
		err := g.handleStatement(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Definer) Define(program *ast.Program) decshared.DecoratedError {
	err := g.firstPass(program)
	if err != nil {
		return err
	}

	return nil
}
