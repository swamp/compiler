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

func writeRecordType(recordType *ast.Record, colorer coloring.Colorer, indentation int) {
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
		writeType(field.Type(), colorer, true, indentation)
	}

	colorer.NewLine(indentation)
	colorer.KeywordString("}")
}

func writeCustomType(customType *ast.CustomType, colorer coloring.Colorer, indentation int) {
	// writeTypeIdentifier(customType.Identifier(), colorer)
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
			writeType(variantParam, colorer, true, indentation)
		}
	}
}

func writeLocalType(localType *ast.LocalType, colorer coloring.Colorer, indentation int) {
	colorer.LocalType(localType.TypeParameter().Identifier().Symbol())
}

func WriteFunctionParameterTypes(functionParameters []ast.Type, colorer coloring.Colorer, indentation int) {
	for index, partType := range functionParameters {
		if index > 0 {
			colorer.OneSpace()
			colorer.RightArrow()
			colorer.OneSpace()
		}
		writeType(partType, colorer, true, indentation)
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

func writeTypeIdentifier(typeIdentifier *ast.TypeIdentifier, colorer coloring.Colorer) {
	moduleReference := typeIdentifier.ModuleReference()
	if moduleReference != nil {
		writeModuleReference(moduleReference, colorer)
	}
	colorer.TypeSymbol(typeIdentifier.Symbol())
}

func writeTypeReference(typeReference *ast.TypeReference, colorer coloring.Colorer) {
	writeTypeIdentifier(typeReference.TypeResolver(), colorer)
	hasArguments := len(typeReference.Arguments()) > 0
	if !hasArguments {
		return
	}
	colorer.OneSpace()
	for _, field := range typeReference.Arguments() {
		writeType(field, colorer, true, 0)
	}
}

func writeAliasStatement(alias *ast.AliasStatement, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("type alias ")
	colorer.AliasNameSymbol(alias.TypeIdentifier().Symbol())
	colorer.OneSpace()
}

func writeFunctionValue(value *ast.FunctionValue, colorer coloring.Colorer, indentation int) {
	writeExpression(value.Expression(), colorer, indentation)
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

func writeCase(caseExpression *ast.CaseCustomType, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("case")
	colorer.OneSpace()
	writeExpression(caseExpression.Test(), colorer, 0)
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
			writeFunctionValueParameters(consequence.Arguments(), colorer)
			colorer.OneSpace()
		}
		colorer.RightArrow()
		colorer.NewLine(indentation + 2)
		writeExpression(consequence.Expression(), colorer, 0)
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

		writeExpression(item.Condition, colorer, 0)
		colorer.OneSpace()
		colorer.RightArrow()
		colorer.OneSpace()
		writeExpression(item.Consequence, colorer, 0)
	}

	colorer.NewLine(indentation + 1)
	colorer.KeywordString("|")
	colorer.OneSpace()
	colorer.OperatorString("_")
	colorer.OneSpace()
	colorer.RightArrow()
	colorer.OneSpace()
	writeExpression(guardExpression.DefaultExpression(), colorer, 0)
}

func writeGetVariable(identifier *ast.VariableIdentifier, colorer coloring.Colorer, indentation int) {
	moduleReference := identifier.ModuleReference()
	if moduleReference != nil {
		writeModuleReference(moduleReference, colorer)
	}
	colorer.VariableSymbol(identifier.Symbol())
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
		writeExpression(expression, colorer, indentation+1)
	}
	colorer.OneSpace()
	colorer.KeywordString("]")
}

func writeRecordLiteral(recordLiteral *ast.RecordLiteral, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("{")
	colorer.OneSpace()

	if recordLiteral.TemplateRecord() != nil {
		colorer.VariableSymbol(recordLiteral.TemplateRecord().Symbol())
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
		writeExpression(assignment.Expression(), colorer, indentation+1)
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
	colorer.NewLine(indentation + 1)

	for index, assignment := range letExpression.Assignments() {
		if index > 0 {
			colorer.NewLine(0)
			colorer.NewLine(indentation + 1)
		}

		colorer.VariableSymbol(assignment.Identifier().Symbol())
		colorer.OneSpace()
		colorer.OperatorString("=")
		colorer.NewLine(indentation + 2)
		writeExpression(assignment.Expression(), colorer, indentation+1)
	}

	colorer.NewLine(indentation)
	colorer.KeywordString("in")
	colorer.NewLine(indentation)
	writeExpression(letExpression.Consequence(), colorer, indentation)
}

func writeIf(letExpression *ast.IfExpression, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("if")
	colorer.OneSpace()
	writeExpression(letExpression.Condition(), colorer, 0)
	colorer.OneSpace()
	colorer.KeywordString("then")
	colorer.NewLine(indentation + 1)
	writeExpression(letExpression.Consequence(), colorer, indentation+1)
	colorer.NewLine(indentation)
	colorer.KeywordString("else")
	colorer.NewLine(indentation + 1)
	writeExpression(letExpression.Alternative(), colorer, indentation+1)
}

func writeAsm(asm *ast.Asm, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString(asm.Asm())
}

func writeBinaryOperator(binaryOperator *ast.BinaryOperator, colorer coloring.Colorer, indentation int) {
	writeExpression(binaryOperator.Left(), colorer, 0)
	colorer.OneSpace()
	colorer.Operator(binaryOperator.OperatorToken())
	colorer.OneSpace()
	writeExpression(binaryOperator.Right(), colorer, 0)
}

func writeUnaryOperator(unaryOperator *ast.UnaryExpression, colorer coloring.Colorer, indentation int) {
	colorer.Operator(unaryOperator.OperatorToken())
	colorer.OneSpace()
	writeExpression(unaryOperator.Left(), colorer, 0)
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
		colorer.LocalType(param.Identifier().Symbol())
	}
}
*/

func writeConstructorCall(constructorCall *ast.ConstructorCall, colorer coloring.Colorer, indentation int) {
	writeTypeIdentifier(constructorCall.TypeIdentifier(), colorer)
	hasArguments := len(constructorCall.Arguments()) > 0
	if hasArguments {
		colorer.OneSpace()
	}
	for index, argument := range constructorCall.Arguments() {
		if index > 0 {
			colorer.OneSpace()
		}
		writeExpression(argument, colorer, indentation)
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

	writeExpression(functionCall.FunctionExpression(), colorer, indentation)
	for index, argument := range functionCall.Arguments() {
		if ignoreIndex == index {
			continue
		}
		colorer.OneSpace()
		writeExpressionAsTerm(argument, colorer, indentation)
	}
}

func writeCustomTypeStatement(customTypeStatement *ast.CustomTypeStatement, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("type ")
	writeTypeIdentifier(customTypeStatement.TypeIdentifier(), colorer)
	colorer.OneSpace()
	/*
		colorer.OperatorString("=")
		colorer.NewLine(indentation+1)
		writeType(customTypeStatement.Type(), colorer,true,  indentation+1)

	*/
}

func writeFunctionValueParameters(parameters []*ast.VariableIdentifier, colorer coloring.Colorer) {
	for index, parameter := range parameters {
		if index > 0 {
			colorer.OneSpace()
		}
		colorer.Parameter(parameter.Symbol())
	}
}

func writeDefinitionAssignment(definition *ast.DefinitionAssignment, colorer coloring.Colorer, indentation int) {
	colorer.Definition(definition.Identifier().Symbol())
	colorer.OneSpace()

	functionValue, isFunctionValue := definition.Expression().(*ast.FunctionValue)
	if isFunctionValue {
		writeFunctionValueParameters(functionValue.Parameters(), colorer)
		colorer.OneSpace()
	}

	colorer.OperatorString("=")
	colorer.NewLine(indentation + 1)
	writeExpression(definition.Expression(), colorer, indentation+1)
}

func writeExternalFunction(externalFunction *ast.ExternalFunction, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("__externalfn")
	colorer.OneSpace()
	colorer.KeywordString(externalFunction.ExternalFunction())
	colorer.OneSpace()
	colorer.KeywordString(fmt.Sprintf("%v", externalFunction.ParameterCount()))
}

func writeImport(importStatement *ast.Import, colorer coloring.Colorer, indentation int) {
	colorer.KeywordString("import")
	colorer.OneSpace()
	for index, pathPart := range importStatement.Path() {
		if index > 0 {
			colorer.OperatorString(".")
		}

		colorer.ModuleReference(pathPart.Symbol())
	}
}

func writeAnnotation(annotation *ast.Annotation, colorer coloring.Colorer, indentation int) {
	colorer.Definition(annotation.Identifier().Symbol())
	colorer.OneSpace()
	colorer.OperatorString(":")
	colorer.OneSpace()
	writeType(annotation.AnnotatedType(), colorer, false, indentation)
}

func writeType(astType ast.Type, colorer coloring.Colorer, addParen bool, indentation int) {
	switch t := astType.(type) {
	case *ast.Record:
		{
			writeRecordType(t, colorer, indentation)
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
			writeCustomType(t, colorer, indentation)
		}
	case *ast.LocalType:
		{
			writeLocalType(t, colorer, indentation)
		}
	default:
		panic(fmt.Errorf(">>> what is this type %T\n", t))
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
	writeExpression(expression, colorer, indentation)
	if !isLeaf {
		colorer.OperatorString(")")
	}
}

func writeExpression(expression ast.Expression, colorer coloring.Colorer, indentation int) {
	switch t := expression.(type) {
	case *ast.Asm:
		{
			writeAsm(t, colorer, indentation)
		}
	case *ast.BinaryOperator:
		{
			writeBinaryOperator(t, colorer, indentation)
		}
	case *ast.BooleanLiteral:
		{
			writeBooleanLiteral(t, colorer, indentation)
		}
	case *ast.CaseCustomType:
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
		panic(fmt.Errorf(">>> what is this expression %T\n", expression))
	}
}

func WriteCode(program *ast.Program, useColor bool) (string, error) {
	var colorer coloring.Colorer
	if useColor {
		colorer = coloring.NewColorerWithColor()
	} else {
		colorer = coloring.NewColorerWithoutColor()
	}
	var lastStatement ast.Expression
	for index, expression := range program.Statements() {
		if index > 0 {
			colorer.NewLine(0)
			numberOfLinesToPad := ast.LinesToInsertBetween(lastStatement, expression)

			if numberOfLinesToPad >= 2 {
				colorer.NewLine(0)
				colorer.NewLine(0)
			}
		}
		switch t := expression.(type) {
		case *ast.AliasStatement:
			writeAliasStatement(t, colorer, 0)
		case *ast.DefinitionAssignment:
			writeDefinitionAssignment(t, colorer, 0)
		case *ast.ExternalFunction:
			writeExternalFunction(t, colorer, 0)
		case *ast.Import:
			writeImport(t, colorer, 0)
		case *ast.Annotation:
			writeAnnotation(t, colorer, 0)
		case *ast.CustomTypeStatement:
			{
				writeCustomTypeStatement(t, colorer, 0)
			}

		default:
			return "", fmt.Errorf(">>> what is this statement %T\n", t)
		}

		lastStatement = expression
	}

	return colorer.String(), nil
}
