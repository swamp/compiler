/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorator

import (
	"fmt"
	"reflect"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type Definer struct {
	typeRepo          decorated.TypeAddAndReferenceMaker
	localAnnotation   *decorated.AnnotationStatement
	localComments     []decorated.Comment
	localCommentBlock *ast.MultilineComment
	verboseFlag       bool
	decorateStream    DecorateStream
}

func NewDefiner(dectorateStream DecorateStream, typeRepo decorated.TypeAddAndReferenceMaker, debugName string) *Definer {
	g := &Definer{verboseFlag: false, localAnnotation: nil, decorateStream: dectorateStream, typeRepo: typeRepo}
	return g
}

func (g *Definer) createAliasTypeFromType(aliasName *ast.Alias, subType dtype.Type) (*dectype.Alias, decshared.DecoratedError) {
	t, typeErr := g.typeRepo.AddTypeAlias(aliasName, subType, nil)
	//	log.Printf("declaring %v\n", aliasName)
	if typeErr != nil {
		panic(typeErr)
	}
	g.localComments = nil
	return t, nil
}

func (g *Definer) AnnotateConstant(identifier *ast.VariableIdentifier, realType dtype.Type) decshared.DecoratedError {
	g.localAnnotation = decorated.NewAnnotation(identifier, realType)
	return nil
}

func (g *Definer) AnnotateFunc(identifier *ast.VariableIdentifier, funcType dtype.Type) (*decorated.AnnotationStatement, decshared.DecoratedError) {
	g.localAnnotation = decorated.NewAnnotation(identifier, funcType)

	return g.localAnnotation, nil
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

func (g *Definer) functionAnnotation(identifier *ast.VariableIdentifier, constantType ast.Type) (*decorated.AnnotationStatement, decshared.DecoratedError) {
	convertedType, decErr := g.convertAnnotation(identifier, constantType)
	if decErr != nil {
		return nil, decErr
	}
	checkedType, _ := convertedType.(dtype.Type)
	if checkedType != nil {
		g.decorateStream.AddDeclaration(identifier, checkedType)
	}
	return g.AnnotateFunc(identifier, convertedType)
}

func (g *Definer) handleAliasStatement(alias *ast.Alias) (*dectype.Alias, decshared.DecoratedError) {
	referencedType, referencedTypeErr := g.findTypeFromAstType(alias)
	if referencedTypeErr != nil {
		return nil, decorated.NewUnknownTypeAliasType(alias, referencedTypeErr)
	}
	return g.createAliasTypeFromType(alias, referencedType)
}

func (g *Definer) findTypeFromAstType(constantType ast.Type) (dtype.Type, decorated.TypeError) {
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

func (g *Definer) handleCustomTypeStatement(customTypeStatement *ast.CustomTypeStatement) (*dectype.CustomTypeAtom, decshared.DecoratedError) {
	hopefullyCustomType := customTypeStatement.Type()
	customType, convertErr := ConvertWrappedOrNormalCustomTypeStatement(hopefullyCustomType, g.typeRepo, nil)
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

func (g *Definer) handleImport(d DecorateStream, importAst *ast.Import) (*decorated.ImportStatement, decshared.DecoratedError) {
	alias := dectype.MakeSingleModuleName(nil)
	if importAst.Alias() != nil {
		alias = dectype.MakeSingleModuleName(importAst.Alias())
	}
	packageRelative := dectype.MakePackageRelativeModuleName(importAst.ModuleName())
	return d.ImportModule(importAst, packageRelative, alias, importAst.ExposeAll(), g.verboseFlag)
}

func (g *Definer) handleExternalFunction(d DecorateStream, externalFunction *ast.ExternalFunction) (*decorated.ExternalFunctionDeclaration, decshared.DecoratedError) {
	g.localComments = nil
	return d.AddExternalFunction(externalFunction)
}

func (g *Definer) handleSinglelineComment(d DecorateStream, singleLineComment *ast.SingleLineComment) (*decorated.SingleLineComment, decshared.DecoratedError) {
	decoratedComment := decorated.NewSingleLineComment(singleLineComment)
	g.localComments = append(g.localComments, decoratedComment)
	return decoratedComment, nil
}

func (g *Definer) handleMultilineComment(d DecorateStream, multilineComment *ast.MultilineComment) (*decorated.MultilineComment, decshared.DecoratedError) {
	decoratedComment := decorated.NewMultilineComment(multilineComment)
	g.localComments = append(g.localComments, decoratedComment)
	return decoratedComment, nil
}

func (g *Definer) handleNamedFunctionValue(d DecorateStream, assignment *ast.FunctionValueNamedDefinition) (*decorated.NamedFunctionValue, decshared.DecoratedError) {
	if g.localAnnotation == nil {
		return nil, decorated.NewMustHaveAnnotationJustBeforeThisDefinition(assignment)
	}
	if g.localAnnotation.Identifier().Name() != assignment.Identifier().Name() {
		return nil, decorated.NewAnnotationMismatch(g.localAnnotation.Identifier(), assignment)
	}
	name := assignment.Identifier()
	expr := assignment.FunctionValue()
	annotatedType := g.localAnnotation.Type()
	if annotatedType == nil {
		return nil, decorated.NewInternalError(fmt.Errorf("can not have nil in local annotation"))
	}
	variableContext := d.NewVariableContext()
	namedFunctionValue, decoratedExpressionErr := decorateNamedFunctionValue(d, variableContext, name, expr, annotatedType, g.localAnnotation, g.localCommentBlock)
	if decoratedExpressionErr != nil {
		return nil, decoratedExpressionErr
	}
	g.localComments = nil
	g.localAnnotation = nil

	return namedFunctionValue, nil
}

func (g *Definer) handleAnnotation(d DecorateStream, declaration *ast.Annotation) (*decorated.AnnotationStatement, decshared.DecoratedError) {
	if g.localAnnotation != nil {
		return nil, decorated.NewAlreadyHaveAnnotationForThisName(declaration, nil)
	}
	annotatedType := declaration.AnnotatedType()
	g.localCommentBlock = declaration.CommentBlock()

	return g.functionAnnotation(declaration.Identifier(), annotatedType)
}

func (g *Definer) convertStatement(statement ast.Expression) (decorated.Statement, decshared.DecoratedError) {
	switch v := statement.(type) {
	case *ast.Alias:
		return g.handleAliasStatement(v)
	case *ast.CustomTypeStatement:
		return g.handleCustomTypeStatement(v)
	case *ast.Annotation:
		return g.handleAnnotation(g.decorateStream, v)
	case *ast.FunctionValueNamedDefinition:
		return g.handleNamedFunctionValue(g.decorateStream, v)
	case *ast.Import:
		return g.handleImport(g.decorateStream, v)
	case *ast.ExternalFunction:
		return g.handleExternalFunction(g.decorateStream, v)
	case *ast.MultilineComment:
		return g.handleMultilineComment(g.decorateStream, v)
	case *ast.SingleLineComment:
		return g.handleSinglelineComment(g.decorateStream, v)
	default:
		return nil, decorated.NewUnknownStatement(token.SourceFileReference{}, statement)
	}
}

func (g *Definer) firstPass(program *ast.SourceFile) ([]decorated.TypeOrToken, decshared.DecoratedError) {
	var rootNodes []decorated.TypeOrToken

	for _, statement := range program.Statements() {
		convertedStatement, err := g.convertStatement(statement)
		if err != nil {
			return nil, err
		}

		if convertedStatement != nil && !reflect.ValueOf(convertedStatement).IsNil() {
			rootNodes = append(rootNodes, convertedStatement)
		}
	}

	return rootNodes, nil
}

func (g *Definer) Define(program *ast.SourceFile) ([]decorated.TypeOrToken, decshared.DecoratedError) {
	rootNodes, err := g.firstPass(program)
	if err != nil {
		return nil, err
	}

	return rootNodes, nil
}
