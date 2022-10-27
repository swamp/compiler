/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package codewriter

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/coloring"
)

func WriteRecordType(recordType *ast.Record, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("{")
	colorer.OneSpace()
	for index, field := range recordType.Fields() {
		if index > 0 {
			colorer.NewLine(indentation)
			colorer.KeywordString(",")
			colorer.OneSpace()
		}
		colorer.RecordField(field.VariableIdentifier().Symbol())
		colorer.OneSpace()
		colorer.OperatorString(":")
		colorer.OneSpace()
		WriteType(field.Type(), colorer, true, indentation)
	}

	colorer.NewLine(indentation)
	colorer.KeywordString("}")
}

func writeCustomType(customType *ast.CustomType, colorer coloring.Colorer, indentation int) {
	// writeTypeIdentifier(customType.Parameter(), colorer)
	// colorer.NewLine(indentation+1)
	// colorer.OneSpace()
	for index, variant := range customType.Variants() {
		if index > 0 {
			colorer.NewLine(indentation)
			colorer.KeywordString("|")
			colorer.OneSpace()
		}
		writeTypeIdentifier(variant.TypeIdentifier(), colorer)
		hasParams := len(variant.Types()) > 0
		if !hasParams {
			continue
		}
		colorer.OneSpace()
		for paramIndex, variantParam := range variant.Types() {
			if paramIndex > 0 {
				colorer.OneSpace()
			}
			WriteType(variantParam, colorer, true, indentation)
		}
	}
}

func writeLocalType(localType *ast.LocalType, colorer coloring.Colorer, indentation int) {
	colorer.LocalType(localType.TypeParameter().Identifier().Symbol())
}

func writeUnmanagedType(unmanagedType *ast.UnmanagedType, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("Unmanaged")
	colorer.OperatorString("<")
	colorer.TypeSymbol(unmanagedType.NativeLanguageTypeName().Symbol())
	colorer.OperatorString(">")
}

func WriteFunctionParameterTypes(functionParameters []ast.Type, colorer coloring.Colorer, indentation int) {
	for index, partType := range functionParameters {
		if index > 0 {
			colorer.OneSpace()
			colorer.RightArrow()
			colorer.OneSpace()
		}
		WriteType(partType, colorer, true, indentation)
	}
}

func writeFunctionType(functionType *ast.FunctionType, colorer coloring.Colorer, addParen bool, indentation int) {
	if addParen {
		colorer.OperatorString("(")
	}
	WriteFunctionParameterTypes(functionType.FunctionParameters(), colorer, indentation)

	if addParen {
		colorer.OperatorString(")")
	}
}

func writeModuleReference(moduleReference *ast.ModuleReference, colorer coloring.Colorer) {
	for index, part := range moduleReference.Parts() {
		if index > 0 {
			colorer.OperatorString(".")
		}

		colorer.ModuleReference(part.TypeIdentifier().Symbol())
	}
	colorer.OperatorString(".")
}

func writeScopedTypeIdentifier(typeIdentifier *ast.TypeIdentifierScoped, colorer coloring.Colorer) {
	moduleReference := typeIdentifier.ModuleReference()
	writeModuleReference(moduleReference, colorer)
	colorer.TypeSymbol(typeIdentifier.Symbol().Symbol())
}

func writeScopedTypeResolver(typeIdentifier *ast.TypeIdentifierScoped, colorer coloring.Colorer) {
	moduleReference := typeIdentifier.ModuleReference()
	writeModuleReference(moduleReference, colorer)
	colorer.TypeGeneratorName(typeIdentifier.Symbol().Symbol())
}

func writeTypeIdentifier(typeIdentifier *ast.TypeIdentifier, colorer coloring.Colorer) {
	colorer.TypeSymbol(typeIdentifier.Symbol())
}

func writeTypeResolver(typeIdentifier *ast.TypeIdentifier, colorer coloring.Colorer) {
	colorer.TypeGeneratorName(typeIdentifier.Symbol())
}

func writeSomeTypeIdentifier(typeIdentifier ast.TypeIdentifierNormalOrScoped, colorer coloring.Colorer) {
	scoped, wasScoped := typeIdentifier.(*ast.TypeIdentifierScoped)
	if wasScoped {
		writeScopedTypeIdentifier(scoped, colorer)
	} else {
		writeTypeIdentifier(typeIdentifier.(*ast.TypeIdentifier), colorer)
	}
}

func writeTypeReference(typeReference *ast.TypeReference, colorer coloring.Colorer) {
	hasArguments := len(typeReference.Arguments()) > 0
	if hasArguments {
		writeTypeResolver(typeReference.TypeIdentifier(), colorer)
	} else {
		writeTypeIdentifier(typeReference.TypeIdentifier(), colorer)
	}
	if !hasArguments {
		return
	}
	colorer.OneSpace()
	for _, field := range typeReference.Arguments() {
		WriteType(field, colorer, true, 0)
	}
}

func writeScopedTypeReference(typeReference *ast.TypeReferenceScoped, colorer coloring.Colorer) {
	hasArguments := len(typeReference.Arguments()) > 0
	if hasArguments {
		writeScopedTypeResolver(typeReference.TypeResolver(), colorer)
	} else {
		writeSomeTypeIdentifier(typeReference.TypeResolver(), colorer)
	}
	if !hasArguments {
		return
	}
	colorer.OneSpace()
	for _, field := range typeReference.Arguments() {
		WriteType(field, colorer, true, 0)
	}
}

func WriteAliasStatement(alias *ast.Alias, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("type alias ")
	colorer.AliasNameSymbol(alias.Identifier().Symbol())
	colorer.OneSpace()
	colorer.KeywordString("=")
	colorer.NewLine(indentation + 1)
	WriteType(alias.ReferencedType(), colorer, false, indentation+1)
}

func writeFunctionValue(value *ast.FunctionValue, colorer coloring.Colorer, indentation int) {
	WriteExpression(value.Expression(), colorer, indentation)
}

func writeSingleLineComment(comment *ast.SingleLineComment, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("x--")
	colorer.KeywordString(comment.Value())
}

func writeMultiLineComment(comment *ast.MultilineComment, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("y{-")
	colorer.NewLine(indentation)
	colorer.KeywordString(comment.Value())
	colorer.NewLine(indentation)
	colorer.KeywordString("-}")
}

func writeCase(caseExpression *ast.CaseForCustomType, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("case")
	colorer.OneSpace()
	WriteExpression(caseExpression.Test(), colorer, 0)
	colorer.OneSpace()
	colorer.KeywordString("of")
	colorer.NewLine(indentation + 1)

	for index, consequence := range caseExpression.Consequences() {
		if index > 0 {
			colorer.NewLine(0)
			colorer.NewLine(indentation + 1)
		}
		colorer.TypeSymbol(consequence.Identifier().Symbol())
		colorer.OneSpace()
		if len(consequence.Arguments()) > 0 {
			writeFunctionIdentifiers(consequence.Arguments(), colorer)
			colorer.OneSpace()
		}
		colorer.RightArrow()
		colorer.NewLine(indentation + 2)
		WriteExpression(consequence.Expression(), colorer, 0)
	}

	// colorer.NewLine(indentation)
}

func writeGuard(guardExpression *ast.GuardExpression, colorer coloring.Colorer, indentation int) {
	for index, item := range guardExpression.Items() {
		if index > 0 {
			colorer.NewLine(indentation + 1)
		}
		colorer.KeywordString("|")
		colorer.OneSpace()

		WriteExpression(item.Condition, colorer, 0)
		colorer.OneSpace()
		colorer.RightArrow()
		colorer.OneSpace()
		WriteExpression(item.Consequence, colorer, 0)
	}

	colorer.NewLine(indentation + 1)
	colorer.KeywordString("|")
	colorer.OneSpace()
	colorer.OperatorString("_")
	colorer.OneSpace()
	colorer.RightArrow()
	colorer.OneSpace()
	WriteExpression(guardExpression.Default().Consequence, colorer, 0)
}

func writeGetVariable(identifier *ast.VariableIdentifier, colorer coloring.Colorer, indentation int) {
	colorer.VariableSymbol(identifier.Symbol())
}

func writeGetScopedVariable(identifier *ast.VariableIdentifierScoped, colorer coloring.Colorer, indentation int) {
	moduleReference := identifier.ModuleReference()
	if moduleReference != nil {
		writeModuleReference(moduleReference, colorer)
	}
	colorer.VariableSymbol(identifier.AstVariableReference().Symbol())
}

func writeIntegerLiteral(identifier *ast.IntegerLiteral, colorer coloring.Colorer, indentation int) {
	colorer.NumberLiteral(identifier.Token)
}

func writeListLiteral(listLiteral *ast.ListLiteral, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("[")
	colorer.OneSpace()
	for index, expression := range listLiteral.Expressions() {
		if index > 0 {
			colorer.OperatorString(", ")
		}
		WriteExpression(expression, colorer, indentation+1)
	}
	colorer.OneSpace()
	colorer.KeywordString("]")
}

func writeRecordLiteral(recordLiteral *ast.RecordLiteral, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("{")
	colorer.OneSpace()

	if recordLiteral.TemplateExpression() != nil {
		WriteExpression(recordLiteral.TemplateExpression(), colorer, indentation)
		colorer.OneSpace()
		colorer.OperatorString("|")
		colorer.OneSpace()
	}
	for index, assignment := range recordLiteral.ParseOrderedAssignments() {
		if index > 0 {
			colorer.NewLine(indentation)
			colorer.OperatorString(", ")
		}
		colorer.RecordField(assignment.Identifier().Symbol())
		colorer.OneSpace()
		colorer.OperatorString("=")
		colorer.OneSpace()
		WriteExpression(assignment.Expression(), colorer, indentation+1)
	}
	colorer.NewLine(indentation)
	colorer.KeywordString("}")
}

func writeStringLiteral(stringLiteral *ast.StringLiteral, colorer coloring.Colorer, indentation int) {
	colorer.StringLiteral(stringLiteral.Token)
}

func writeBooleanLiteral(boolean *ast.BooleanLiteral, colorer coloring.Colorer, indentation int) {
	colorer.BooleanLiteral(boolean.Token())
}

func writeLet(letExpression *ast.Let, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("let")

	for index, assignment := range letExpression.Assignments() {
		if index >= 0 {
			colorer.NewLine(indentation + 1)
		}

		colorer.VariableSymbol(assignment.Identifiers()[0].Symbol())
		colorer.OneSpace()
		colorer.OperatorString("=")
		colorer.OneSpace()
		WriteExpression(assignment.Expression(), colorer, indentation+1)
	}

	colorer.NewLine(indentation)
	colorer.KeywordString("in")
	colorer.NewLine(indentation)
	WriteExpression(letExpression.Consequence(), colorer, indentation)
}

func writeIf(letExpression *ast.IfExpression, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("if")
	colorer.OneSpace()
	WriteExpression(letExpression.Condition(), colorer, 0)
	colorer.OneSpace()
	colorer.KeywordString("then")
	colorer.NewLine(indentation + 1)
	WriteExpression(letExpression.Consequence(), colorer, indentation+1)
	colorer.NewLine(indentation)
	colorer.KeywordString("else")
	colorer.NewLine(indentation + 1)
	WriteExpression(letExpression.Alternative(), colorer, indentation+1)
}

func writeBinaryOperator(binaryOperator *ast.BinaryOperator, colorer coloring.Colorer, indentation int) {
	WriteExpression(binaryOperator.Left(), colorer, 0)
	colorer.OneSpace()
	colorer.Operator(binaryOperator.OperatorToken())
	colorer.OneSpace()
	WriteExpression(binaryOperator.Right(), colorer, 0)
}

func writeUnaryOperator(unaryOperator *ast.UnaryExpression, colorer coloring.Colorer, indentation int) {
	colorer.Operator(unaryOperator.OperatorToken())
	colorer.OneSpace()
	WriteExpression(unaryOperator.Left(), colorer, 0)
}

func writeLookups(lookups *ast.Lookups, colorer coloring.Colorer, indentation int) {
	writeGetVariable(lookups.ContextIdentifier(), colorer, indentation)
	for _, variable := range lookups.FieldNames() {
		colorer.OperatorString(".")
		writeGetVariable(variable, colorer, indentation)
	}
}

/* TODO:
func writeTypeParameters(typeParameters []*ast.TypeParameter, colorer coloring.Colorer) {
	for index, param := range typeParameters {
		if index > 0 {
			colorer.OneSpace()
		}
		colorer.LocalType(param.Parameter().Symbol())
	}
}
*/

func writeConstructorCall(constructorCall *ast.ConstructorCall, colorer coloring.Colorer, indentation int) {
	colorer.OperatorString("constructorcall")
	writeSomeTypeIdentifier(constructorCall.TypeReference().SomeTypeIdentifier(), colorer)
	hasArguments := len(constructorCall.Arguments()) > 0
	if hasArguments {
		colorer.OneSpace()
	}
	colorer.OperatorString("constructorcall after")
	for index, argument := range constructorCall.Arguments() {
		if index > 0 {
			colorer.OneSpace()
		}
		WriteExpression(argument, colorer, indentation)
	}
}

func functionArgumentsCanBeRightPiped(arguments []ast.Expression) (*ast.FunctionCall, int) {
	lastArgumentIndex := len(arguments) - 1
	lastArgument := arguments[lastArgumentIndex]
	lastArgumentFunctionCall, _ := lastArgument.(*ast.FunctionCall)
	if lastArgumentFunctionCall == nil {
		lastArgumentIndex = -1
	}

	return lastArgumentFunctionCall, lastArgumentIndex
}

func writeFunctionCall(functionCall *ast.FunctionCall, colorer coloring.Colorer, indentation int) {
	nextFunctionCall, ignoreIndex := functionArgumentsCanBeRightPiped(functionCall.Arguments())
	if nextFunctionCall != nil {
		writeFunctionCall(nextFunctionCall, colorer, indentation)
		colorer.NewLine(indentation + 1)
		colorer.OperatorString("|>")
		colorer.OneSpace()
	}

	WriteExpression(functionCall.FunctionExpression(), colorer, indentation)
	for index, argument := range functionCall.Arguments() {
		if ignoreIndex == index {
			continue
		}
		colorer.OneSpace()
		writeExpressionAsTerm(argument, colorer, indentation)
	}
}

func WriteCustomTypeStatement(customTypeStatement *ast.CustomType, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("type ")
	writeTypeIdentifier(customTypeStatement.Identifier(), colorer)
	colorer.NewLine(indentation + 1)
	colorer.OperatorString("=")
	colorer.OneSpace()
	writeCustomType(customTypeStatement, colorer, indentation+1)
}

func writeFunctionIdentifiers(parameters []*ast.VariableIdentifier, colorer coloring.Colorer) {
	for index, parameter := range parameters {
		if index > 0 {
			colorer.OneSpace()
		}
		colorer.Parameter(parameter.Symbol())
	}
}

func writeFunctionValueParameters(parameters []*ast.FunctionParameter, colorer coloring.Colorer) {
	for index, parameter := range parameters {
		if index > 0 {
			colorer.OneSpace()
		}
		colorer.Parameter(parameter.Identifier().Symbol())
		colorer.OneSpace()
		colorer.OperatorString(":")
		colorer.OneSpace()
		WriteType(parameter.Type(), colorer, false, 0)
	}
}

func writeDefinitionAssignment(definition *ast.FunctionValueNamedDefinition, colorer coloring.Colorer, indentation int) {
	colorer.Definition(definition.Identifier().Symbol())
	colorer.OneSpace()

	functionValue := definition.FunctionValue()
	writeFunctionValueParameters(functionValue.Parameters(), colorer)
	colorer.OneSpace()

	colorer.OperatorString("=")
	colorer.NewLine(indentation + 1)
	WriteExpression(definition.FunctionValue(), colorer, indentation+1)
}

func writeDefinitionAssignmentConstant(definition *ast.ConstantDefinition, colorer coloring.Colorer, indentation int) {
	colorer.Definition(definition.Identifier().Symbol())
	colorer.OneSpace()
	colorer.OperatorString("=")
	colorer.NewLine(indentation + 1)
	WriteExpression(definition.Expression(), colorer, indentation+1)
}

func writeImport(importStatement *ast.Import, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("import")
	colorer.OneSpace()
	for index, pathPart := range importStatement.ModuleName().Parts() {
		if index > 0 {
			colorer.OperatorString(".")
		}

		colorer.ModuleReference(pathPart.TypeIdentifier().Symbol())
	}
}

func WriteType(astType ast.Type, colorer coloring.Colorer, addParen bool, indentation int) {
	switch t := astType.(type) {
	case *ast.Alias:
		{
			WriteAliasStatement(t, colorer, indentation)
		}
	case *ast.Record:
		{
			WriteRecordType(t, colorer, indentation)
		}
	case *ast.TypeReference:
		{
			writeTypeReference(t, colorer)
		}
	case *ast.FunctionType:
		{
			writeFunctionType(t, colorer, addParen, indentation)
		}
	case *ast.CustomType:
		{
			WriteCustomTypeStatement(t, colorer, indentation)
		}
	case *ast.LocalType:
		{
			writeLocalType(t, colorer, indentation)
		}
	case *ast.UnmanagedType:
		{
			writeUnmanagedType(t, colorer, indentation)
		}
	case *ast.TypeReferenceScoped:
		{
			writeScopedTypeReference(t, colorer)
		}
	default:
		panic(fmt.Errorf("what is this type %T", t))
	}
}

func isLeaf(expression ast.Expression) bool {
	switch expression.(type) {
	case *ast.BooleanLiteral:
		return true
	case *ast.StringLiteral:
		return true
	case *ast.VariableIdentifier:
		return true
	case *ast.IntegerLiteral:
		return true
	}
	return false
}

func writeExpressionAsTerm(expression ast.Expression, colorer coloring.Colorer, indentation int) {
	isLeaf := isLeaf(expression)
	if !isLeaf {
		colorer.OperatorString("(")
	}
	WriteExpression(expression, colorer, indentation)
	if !isLeaf {
		colorer.OperatorString(")")
	}
}

func WriteExpression(expression ast.Expression, colorer coloring.Colorer, indentation int) {
	switch t := expression.(type) {
	case *ast.BinaryOperator:
		{
			writeBinaryOperator(t, colorer, indentation)
		}
	case *ast.BooleanLiteral:
		{
			writeBooleanLiteral(t, colorer, indentation)
		}
	case *ast.CaseForCustomType:
		{
			writeCase(t, colorer, indentation)
		}
	case *ast.SingleLineComment:
		{
			writeSingleLineComment(t, colorer, indentation)
		}
	case *ast.MultilineComment:
		{
			writeMultiLineComment(t, colorer, indentation)
		}
	case *ast.ConstructorCall:
		{
			writeConstructorCall(t, colorer, indentation)
		}
	case *ast.FunctionCall:
		{
			writeFunctionCall(t, colorer, indentation)
		}
	case *ast.FunctionValue:
		{
			writeFunctionValue(t, colorer, indentation)
		}
	case *ast.VariableIdentifier:
		{
			writeGetVariable(t, colorer, indentation)
		}
	case *ast.VariableIdentifierScoped:
		{
			writeGetScopedVariable(t, colorer, indentation)
		}
	case *ast.TypeIdentifier:
		{
			writeTypeIdentifier(t, colorer)
		}
	case *ast.IntegerLiteral:
		{
			writeIntegerLiteral(t, colorer, indentation)
		}
	case *ast.Let:
		{
			writeLet(t, colorer, indentation)
		}
	case *ast.ListLiteral:
		{
			writeListLiteral(t, colorer, indentation)
		}
	case *ast.IfExpression:
		{
			writeIf(t, colorer, indentation)
		}
	case *ast.GuardExpression:
		{
			writeGuard(t, colorer, indentation)
		}
	case *ast.RecordLiteral:
		{
			writeRecordLiteral(t, colorer, indentation)
		}
	case *ast.StringLiteral:
		{
			writeStringLiteral(t, colorer, indentation)
		}
	case *ast.UnaryExpression:
		{
			writeUnaryOperator(t, colorer, indentation)
		}
	case *ast.Lookups:
		writeLookups(t, colorer, indentation)
	default:
		panic(fmt.Errorf("what is this expression %T", expression))
	}
}

func WriteCode(program *ast.SourceFile, useColor bool) (string, error) {
	var colorer coloring.Colorer
	if useColor {
		colorer = coloring.NewColorerWithColor()
	} else {
		colorer = coloring.NewColorerWithoutColor()
	}

	return WriteCodeUsingColorer(program, colorer, 0)
}

func WriteStatementUsingColorer(expression ast.Expression, colorer coloring.Colorer, indentation int) error {
	switch t := expression.(type) {
	case *ast.Alias:
		WriteAliasStatement(t, colorer, indentation)
	case *ast.FunctionValueNamedDefinition:
		writeDefinitionAssignment(t, colorer, indentation)

	case *ast.Import:
		writeImport(t, colorer, indentation)
	case *ast.CustomType:
		{
			WriteCustomTypeStatement(t, colorer, indentation)
		}
	case *ast.ConstantDefinition:
		writeDefinitionAssignmentConstant(t, colorer, indentation)
	default:
		panic(fmt.Errorf("what is this statement %T", t))
	}

	return nil
}

func WriteCodeUsingColorer(program *ast.SourceFile, colorer coloring.Colorer, indentation int) (string, error) {
	var lastStatement ast.Expression
	for index, expression := range program.Statements() {
		if index > 0 {
			colorer.NewLine(0)
			numberOfLinesToPad := ast.LinesToInsertBetween(lastStatement, expression)

			if numberOfLinesToPad >= 2 {
				colorer.NewLine(indentation)
				colorer.NewLine(indentation)
			}
		}

		if err := WriteStatementUsingColorer(expression, colorer, indentation); err != nil {
			return "", err
		}
		lastStatement = expression
	}

	return colorer.String(), nil
}
