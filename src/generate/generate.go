/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate

import (
	"fmt"
	"reflect"

	asmcompile "github.com/swamp/assembler/compiler"
	assembler "github.com/swamp/assembler/lib"
	swampopcodeinst "github.com/swamp/opcodes/instruction"
	swamppack "github.com/swamp/pack/lib"

	decorator "github.com/swamp/compiler/src/decorated/convert"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/typeinfo"
)

type Function struct {
	name           *decorated.FullyQualifiedVariableName
	signature      swamppack.TypeRef
	parameterCount uint
	variableCount  uint
	constants      []*assembler.Constant
	opcodes        []byte
}

type ExternalFunction struct {
	name           *decorated.FullyQualifiedVariableName
	signature      swamppack.TypeRef
	parameterCount uint
}

func Pack(functions []*Function, externalFunctions []*ExternalFunction, typeInfoPayload []byte, lookup typeinfo.TypeLookup) ([]byte, error) {
	constantRepo := swamppack.NewConstantRepo()

	for _, externalFunction := range externalFunctions {
		constantRepo.AddExternalFunction(externalFunction.name.ResolveToString(), externalFunction.parameterCount)
	}

	for _, declareFunction := range functions {
		constantRepo.AddFunctionDeclaration(declareFunction.name.ResolveToString(), declareFunction.signature, declareFunction.parameterCount)
	}

	for _, function := range functions {
		var packConstants []*swamppack.Constant
		for _, subConstant := range function.constants {
			var packConstant *swamppack.Constant
			switch subConstant.ConstantType() {
			case assembler.ConstantTypeInteger:
				packConstant = constantRepo.AddInteger(subConstant.IntegerValue())
			case assembler.ConstantTypeString:
				packConstant = constantRepo.AddString(subConstant.StringValue())
			case assembler.ConstantTypeBoolean:
				packConstant = constantRepo.AddBoolean(subConstant.BooleanValue())
			case assembler.ConstantTypeFunction:
				refConstant, functionRefErr := constantRepo.AddFunctionReference(subConstant.FunctionReferenceFullyQualifiedName())
				if functionRefErr != nil {
					return nil, functionRefErr
				}
				packConstant = refConstant
			case assembler.ConstantTypeFunctionExternal:
				refConstant, functionRefErr := constantRepo.AddExternalFunctionReference(subConstant.FunctionReferenceFullyQualifiedName())
				if functionRefErr != nil {
					return nil, functionRefErr
				}
				packConstant = refConstant
			default:
				return nil, fmt.Errorf("not handled constanttype %v", subConstant)
			}
			if packConstant == nil {
				return nil, fmt.Errorf("internal error: not handled constanttype %v", subConstant)
			}

			packConstants = append(packConstants, packConstant)
		}
		constantRepo.AddFunction(function.name.ResolveToString(), function.signature, function.parameterCount, function.variableCount, packConstants, function.opcodes)
	}
	octets, packErr := swamppack.Pack(constantRepo, typeInfoPayload)
	if packErr != nil {
		return nil, packErr
	}
	return octets, nil
}

func NewFunction(fullyQualifiedName *decorated.FullyQualifiedVariableName, signature swamppack.TypeRef, parameterCount uint, variableCount uint, constants []*assembler.Constant, opcodes []byte) *Function {
	f := &Function{name: fullyQualifiedName, signature: signature, parameterCount: parameterCount, variableCount: variableCount, constants: constants, opcodes: opcodes}
	return f
}

func NewExternalFunction(fullyQualifiedName *decorated.FullyQualifiedVariableName, signature swamppack.TypeRef, parameterCount uint) *ExternalFunction {
	f := &ExternalFunction{name: fullyQualifiedName, signature: signature, parameterCount: parameterCount}
	return f
}

func (f *Function) String() string {
	return fmt.Sprintf("[function %v %v %v %v]", f.name, f.signature, f.parameterCount, f.constants)
}

func (f *Function) Opcodes() []byte {
	return f.opcodes
}

type Generator struct {
	code *assembler.Code
}

func NewGenerator() *Generator {
	return &Generator{code: assembler.NewCode()}
}

func arithmeticToBinaryOperatorType(operatorType decorated.ArithmeticOperatorType) swampopcodeinst.BinaryOperatorType {
	switch operatorType {
	case decorated.ArithmeticPlus:
		return swampopcodeinst.BinaryOperatorArithmeticPlus
	case decorated.ArithmeticMinus:
		return swampopcodeinst.BinaryOperatorArithmeticMinus
	case decorated.ArithmeticMultiply:
		return swampopcodeinst.BinaryOperatorArithmeticMultiply
	case decorated.ArithmeticDivide:
		return swampopcodeinst.BinaryOperatorArithmeticDivide
	case decorated.ArithmeticAppend:
		return swampopcodeinst.BinaryOperatorArithmeticAppend
	case decorated.ArithmeticFixedMultiply:
		return swampopcodeinst.BinaryOperatorArithmeticFixedMultiply
	case decorated.ArithmeticFixedDivide:
		return swampopcodeinst.BinaryOperatorArithmeticFixedDivide
	}

	panic("unknown binary operator")
}

func bitwiseToUnaryOperatorType(operatorType decorated.BitwiseUnaryOperatorType) swampopcodeinst.UnaryOperatorType {
	switch operatorType {
	case decorated.BitwiseUnaryNot:
		return swampopcodeinst.UnaryOperatorBitwiseNot
	}
	panic("illegal unaryoperator")
}

func logicalToUnaryOperatorType(operatorType decorated.LogicalUnaryOperatorType) swampopcodeinst.UnaryOperatorType {
	switch operatorType {
	case decorated.LogicalUnaryNot:
		return swampopcodeinst.UnaryOperatorNot
	}
	panic("illegal unaryoperator")
}

func bitwiseToBinaryOperatorType(operatorType decorated.BitwiseOperatorType) swampopcodeinst.BinaryOperatorType {
	switch operatorType {
	case decorated.BitwiseAnd:
		return swampopcodeinst.BinaryOperatorBitwiseAnd
	case decorated.BitwiseOr:
		return swampopcodeinst.BinaryOperatorBitwiseOr
	case decorated.BitwiseXor:
		return swampopcodeinst.BinaryOperatorBitwiseXor
	}

	return 0
}

func generateListAppend(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ArithmeticOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	// leftVar := context.AllocateTempVariable("arit-left")
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), context, definitions, "list-append-left")
	if leftErr != nil {
		return leftErr
	}

	// rightVar := context.AllocateTempVariable("arit-right")
	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), context, definitions, "list-append-right")
	if rightErr != nil {
		return rightErr
	}
	code.ListAppend(target, leftVar, rightVar)
	context.FreeVariableIfNeeded(leftVar)
	context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generateStringAppend(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ArithmeticOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), context, definitions, "string-append-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), context, definitions, "string-append-right")
	if rightErr != nil {
		return rightErr
	}
	code.StringAppend(target, leftVar, rightVar)
	context.FreeVariableIfNeeded(leftVar)
	context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generateAsm(code *assembler.Code, target assembler.TargetVariable, asm *decorated.AsmConstant, context *assembler.Context, definitions *decorator.VariableContext) error {
	compileErr := asmcompile.CompileToCodeAndContext(asm.Asm().Asm(), code, context)
	return compileErr
}

func generateListCons(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ConsOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	// leftVar := context.AllocateTempVariable("arit-left")
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), context, definitions, "cons-left")
	if leftErr != nil {
		return leftErr
	}

	// rightVar := context.AllocateTempVariable("arit-right")
	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), context, definitions, "cons-right")
	if rightErr != nil {
		return rightErr
	}

	code.ListConj(target, leftVar, rightVar)
	context.FreeVariableIfNeeded(leftVar)
	context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generateArithmetic(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ArithmeticOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), context, definitions, "arith-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), context, definitions, "arit-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := arithmeticToBinaryOperatorType(operator.OperatorType())
	code.BinaryOperator(target, leftVar, rightVar, opcodeBinaryOperator)
	context.FreeVariableIfNeeded(leftVar)
	context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generateUnaryBitwise(code *assembler.Code, target assembler.TargetVariable, operator *decorated.BitwiseUnaryOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), context, definitions, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := bitwiseToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target, leftVar, opcodeUnaryOperatorType)
	context.FreeVariableIfNeeded(leftVar)
	return nil
}

func generateUnaryLogical(code *assembler.Code, target assembler.TargetVariable, operator *decorated.LogicalUnaryOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), context, definitions, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := logicalToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target, leftVar, opcodeUnaryOperatorType)
	context.FreeVariableIfNeeded(leftVar)
	return nil
}

func generateBitwise(code *assembler.Code, target assembler.TargetVariable, operator *decorated.BitwiseOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), context, definitions, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), context, definitions, "bitwise-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := bitwiseToBinaryOperatorType(operator.OperatorType())
	code.BinaryOperator(target, leftVar, rightVar, opcodeBinaryOperator)
	context.FreeVariableIfNeeded(leftVar)
	context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generateLogical(code *assembler.Code, target assembler.TargetVariable, operator *decorated.LogicalOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	leftErr := generateExpression(code, target, operator.Left(), context, definitions)
	if leftErr != nil {
		return leftErr
	}

	codeAlternative := assembler.NewCode()
	rightErr := generateExpression(codeAlternative, target, operator.Right(), context, definitions)
	if rightErr != nil {
		return rightErr
	}
	afterLabel := codeAlternative.Label(nil, "after-alternative")

	if operator.OperatorType() == decorated.LogicalAnd {
		code.BranchFalse(target, afterLabel)
	} else if operator.OperatorType() == decorated.LogicalOr {
		code.BranchTrue(target, afterLabel)
	}
	code.Copy(codeAlternative)

	return nil
}

func booleanToBinaryOperatorType(operatorType decorated.BooleanOperatorType) swampopcodeinst.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return swampopcodeinst.BinaryOperatorBooleanEqual
	case decorated.BooleanNotEqual:
		return swampopcodeinst.BinaryOperatorBooleanNotEqual
	case decorated.BooleanLess:
		return swampopcodeinst.BinaryOperatorBooleanLess
	case decorated.BooleanLessOrEqual:
		return swampopcodeinst.BinaryOperatorBooleanLessOrEqual
	case decorated.BooleanGreater:
		return swampopcodeinst.BinaryOperatorBooleanGreater
	case decorated.BooleanGreaterOrEqual:
		return swampopcodeinst.BinaryOperatorBooleanGreaterOrEqual
	}

	return 0
}

func generateBoolean(code *assembler.Code, target assembler.TargetVariable, operator *decorated.BooleanOperator, context *assembler.Context, definitions *decorator.VariableContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), context, definitions, "bool-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), context, definitions, "bool-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := booleanToBinaryOperatorType(operator.OperatorType())
	code.BinaryOperator(target, leftVar, rightVar, opcodeBinaryOperator)
	context.FreeVariableIfNeeded(leftVar)
	context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generateLet(code *assembler.Code, target assembler.TargetVariable, let *decorated.Let, context *assembler.Context, definitions *decorator.VariableContext) error {
	for _, assignment := range let.Assignments() {
		varName := assembler.NewVariableName(assignment.Name().Name())
		targetVar := context.AllocateVariable(varName)
		genErr := generateExpression(code, targetVar, assignment.Expression(), context, definitions)
		if genErr != nil {
			return genErr
		}
	}

	codeErr := generateExpression(code, target, let.Consequence(), context, definitions)
	if codeErr != nil {
		return codeErr
	}

	return nil
}

func generateLookups(code *assembler.Code, target assembler.TargetVariable, lookups *decorated.Lookups, context *assembler.Context) error {
	variableName := assembler.NewVariableName(lookups.Variable().Identifier().Name())
	a := context.FindVariable(variableName)
	if a == nil {
		return fmt.Errorf("couldn't find name %v", lookups.Variable())
	}
	var structLookups []uint8
	for _, indexLookups := range lookups.LookupFields() {
		structLookups = append(structLookups, uint8(indexLookups.Index()))
	}
	code.Lookups(target, a, structLookups)
	return nil
}

func generateIf(code *assembler.Code, target assembler.TargetVariable, ifExpr *decorated.If, context *assembler.Context, definitions *decorator.VariableContext) error {
	conditionVar, testErr := generateExpressionWithSourceVar(code, ifExpr.Condition(), context, definitions, "if-condition")
	if testErr != nil {
		return testErr
	}

	consequenceCode := assembler.NewCode()
	consequenceContext := context.MakeScopeContext()
	consErr := generateExpression(consequenceCode, target, ifExpr.Consequence(), consequenceContext, definitions)
	if consErr != nil {
		return consErr
	}
	consequenceContext.Free()

	alternativeCode := assembler.NewCode()

	alternativeLabel := alternativeCode.Label(nil, "if-alternative")
	alternativeContext := context.MakeScopeContext()
	altErr := generateExpression(alternativeCode, target, ifExpr.Alternative(), alternativeContext, definitions)
	if altErr != nil {
		return altErr
	}
	endLabel := alternativeCode.Label(nil, "if-end")
	alternativeContext.Free()

	code.BranchFalse(conditionVar, alternativeLabel)
	context.FreeVariableIfNeeded(conditionVar)
	consequenceCode.Jump(endLabel)
	code.Copy(consequenceCode)
	code.Copy(alternativeCode)

	return nil
}

func generateCase(code *assembler.Code, target assembler.TargetVariable, caseExpr *decorated.Case, context *assembler.Context, definitions *decorator.VariableContext) error {
	testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), context, definitions, "cast-test")
	if testErr != nil {
		return testErr
	}

	var consequences []*assembler.CaseConsequence
	var consequencesCodes []*assembler.Code

	for _, consequence := range caseExpr.Consequences() {
		consequenceContext := context.MakeScopeContext()

		consequencesCode := assembler.NewCode()

		var parameters []assembler.SourceVariable
		for _, param := range consequence.Parameters() {
			consequenceLabelVariableName := assembler.NewVariableName(param.Identifier().Name())
			paramVariable := consequenceContext.AllocateVariable(consequenceLabelVariableName)
			parameters = append(parameters, paramVariable)
		}
		labelVariableName := assembler.NewVariableName(consequence.Identifier().Name())
		caseLabel := consequencesCode.Label(labelVariableName, "case")
		caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), consequenceContext, definitions)
		if caseExprErr != nil {
			return caseExprErr
		}
		asmConsequence := assembler.NewCaseConsequence(uint8(consequence.InternalIndex()), parameters, caseLabel)
		consequences = append(consequences, asmConsequence)

		consequencesCodes = append(consequencesCodes, consequencesCode)

		consequenceContext.Free()
	}

	var defaultCase *assembler.CaseConsequence
	if caseExpr.DefaultCase() != nil {
		consequencesCode := assembler.NewCode()
		defaultContext := context.MakeScopeContext()

		decoratedDefault := caseExpr.DefaultCase()
		defaultLabel := consequencesCode.Label(nil, "default")
		caseExprErr := generateExpression(consequencesCode, target, decoratedDefault, defaultContext, definitions)
		if caseExprErr != nil {
			return caseExprErr
		}
		defaultCase = assembler.NewCaseConsequence(0xff, nil, defaultLabel)
		consequencesCodes = append(consequencesCodes, consequencesCode)
		//		endLabel := consequencesBlockCode.Label(nil, "if-end")
		defaultContext.Free()
	}

	consequencesBlockCode := assembler.NewCode()

	lastConsequnce := consequencesCodes[len(consequencesCodes)-1]

	labelVariableEndName := assembler.NewVariableName("case end")
	endLabel := lastConsequnce.Label(labelVariableEndName, "caseend")
	for index, consequenceCode := range consequencesCodes {
		if index != len(consequencesCodes)-1 {
			consequenceCode.Jump(endLabel)
		}
	}

	for _, consequenceCode := range consequencesCodes {
		consequencesBlockCode.Copy(consequenceCode)
	}

	code.Case(target, testVar, consequences, defaultCase)
	context.FreeVariableIfNeeded(testVar)
	code.Copy(consequencesBlockCode)

	return nil
}

func generateStringLiteral(code *assembler.Code, target assembler.TargetVariable, str *decorated.StringLiteral, context *assembler.Context) error {
	constant := context.Constants().AllocateStringConstant(str.Value())
	code.CopyVariable(target, constant)
	return nil
}

func generateIntLiteral(code *assembler.Code, target assembler.TargetVariable, integer *decorated.IntegerLiteral, context *assembler.Context) error {
	constant := context.Constants().AllocateIntegerConstant(integer.Value())
	code.CopyVariable(target, constant)
	return nil
}

func generateFixedLiteral(code *assembler.Code, target assembler.TargetVariable, fixed *decorated.FixedLiteral, context *assembler.Context) error {
	constant := context.Constants().AllocateIntegerConstant(fixed.Value())
	code.CopyVariable(target, constant)
	return nil
}

func generateResourceNameLiteral(code *assembler.Code, target assembler.TargetVariable, resourceName *decorated.ResourceNameLiteral, context *assembler.Context) error {
	constant := context.Constants().AllocateResourceNameConstant(resourceName.Value())
	code.CopyVariable(target, constant)
	return nil
}


func generateGetVariable(code *assembler.Code, target assembler.TargetVariable, getVar *decorated.GetVariableOrReferenceFunction, context *assembler.Context) error {
	varName := assembler.NewVariableName(getVar.Identifier().Name())
	variable := context.FindVariable(varName)
	code.CopyVariable(target, variable)
	return nil
}

func generateCustomTypeVariantConstructor(code *assembler.Code, target assembler.TargetVariable, constructor *decorated.CustomTypeVariantConstructor, context *assembler.Context, definitions *decorator.VariableContext) error {
	var arguments []assembler.SourceVariable
	for _, arg := range constructor.Arguments() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, context, definitions, "customTypeVariantArgs")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg)
	}

	code.CreateEnum(target, constructor.CustomTypeVariantIndex(), arguments)

	return nil
}

func generateCurry(code *assembler.Code, target assembler.TargetVariable, call *decorated.CurryFunction, context *assembler.Context, definitions *decorator.VariableContext) error {
	var arguments []assembler.SourceVariable
	for _, arg := range call.ArgumentsToSave() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, context, definitions, "sourceSave")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg)
	}
	functionRegister, functionGenErr := generateExpressionWithSourceVar(code, call.FunctionValue(), context, definitions, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}
	code.Curry(target, functionRegister, arguments)

	return nil
}

func generateFunctionCall(code *assembler.Code, target assembler.TargetVariable, call *decorated.FunctionCall, context *assembler.Context, definitions *decorator.VariableContext) error {
	var arguments []assembler.SourceVariable
	for _, arg := range call.Arguments() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, context, definitions, "arg")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg)
	}

	fn := call.FunctionValue()

	functionRegister, functionGenErr := generateExpressionWithSourceVar(code, fn, context, definitions, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}
	code.Call(target, functionRegister, arguments)

	return nil
}

func generateRecurCall(code *assembler.Code, call *decorated.RecurCall, context *assembler.Context, definitions *decorator.VariableContext) error {
	var arguments []assembler.SourceVariable
	for _, arg := range call.Arguments() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, context, definitions, "recurarg")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg)
	}

	code.Recur(arguments)

	return nil
}

func generateBoolConstant(code *assembler.Code, target assembler.TargetVariable, test *decorated.BooleanLiteral, context *assembler.Context) error {
	constant := context.Constants().AllocateBooleanConstant(test.Value())
	code.CopyVariable(target, constant)
	return nil
}

func generateRecordSortedAssignments(code *assembler.Code, target assembler.TargetVariable, sortedAssignments []*decorated.RecordLiteralAssignment, context *assembler.Context, definitions *decorator.VariableContext) error {
	variables := make([]assembler.SourceVariable, len(sortedAssignments))
	for index, assignment := range sortedAssignments {
		debugName := fmt.Sprintf("assign%v", assignment.FieldName())
		assignmentVar, genErr := generateExpressionWithSourceVar(code, assignment.Expression(), context, definitions, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = assignmentVar
	}
	code.Constructor(target, variables)

	return nil
}

func generateRecordLiteral(code *assembler.Code, target assembler.TargetVariable, record *decorated.RecordLiteral, context *assembler.Context, definitions *decorator.VariableContext) error {
	if record.RecordTemplate() != nil {
		structToCopyVar, genErr := generateExpressionWithSourceVar(code, record.RecordTemplate(), context, definitions, "gopher")
		if genErr != nil {
			return genErr
		}
		var updateFields []assembler.UpdateField
		for _, assignment := range record.SortedAssignments() {
			debugName := fmt.Sprintf("update%v", assignment.FieldName())
			assignmentVar, genErr := generateExpressionWithSourceVar(code, assignment.Expression(), context, definitions, debugName)
			if genErr != nil {
				return genErr
			}
			field := assembler.UpdateField{TargetField: uint8(assignment.Index()), Source: assignmentVar}
			updateFields = append(updateFields, field)
		}
		code.UpdateStruct(target, structToCopyVar, updateFields)
	} else {
		return generateRecordSortedAssignments(code, target, record.SortedAssignments(), context, definitions)
	}
	return nil
}

func generateList(code *assembler.Code, target assembler.TargetVariable, list *decorated.ListLiteral, context *assembler.Context, definitions *decorator.VariableContext) error {
	variables := make([]assembler.SourceVariable, len(list.Expressions()))
	for index, expr := range list.Expressions() {
		debugName := fmt.Sprintf("listliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, context, definitions, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar
	}
	code.ListLiteral(target, variables)
	return nil
}

/*
	foundNamedExpression := definitions.FindNamedDecoratedExpression(getVar.Identifier())
	if foundNamedExpression == nil {
		expr := definitions.FindNamedDecoratedExpression(getVar.Identifier())
		if expr != nil {
			typeRef, typeErr := lookup.Lookup(expr.Type())
			if typeErr != nil {
				return nil, typeErr
			}
			funcAnnotation := assembler.NewFunctionAnnotation(uint8(typeRef))
			funcName := assembler.NewFunctionName(expr.FullyQualifiedName())
			constant, _ := c.AddFunctionReferenceConstant(funcName, funcAnnotation)
			return constant, nil
		}
		return nil, fmt.Errorf("didn't find named expression %v (%v)", getVar.Identifier(), definitions)
	}

	typeRef, typeErr := lookup.Lookup(foundNamedExpression.Type())
	if typeErr != nil {
		return nil, typeErr
	}
	funcAnnotation := assembler.NewFunctionAnnotation(uint8(typeRef))
	funcName := assembler.NewFunctionName(foundNamedExpression.FullyQualifiedName())
	foundFunc, foundFuncErr := c.AddFunctionReferenceConstant(funcName, funcAnnotation)
	if foundFuncErr != nil {
		return nil, foundFuncErr
	}
	if foundFunc != nil {
		return foundFunc, nil
	}
	return nil, fmt.Errorf("couldn't find variable or constant %v", getVar)
*/

func generateExpressionWithSourceVar(code *assembler.Code, expr decorated.DecoratedExpression, c *assembler.Context, definitions *decorator.VariableContext, debugName string) (assembler.SourceVariable, error) {
	getVar, _ := expr.(*decorated.GetVariableOrReferenceFunction)
	if getVar != nil {
		ident := getVar.Identifier()
		getVarName := assembler.NewVariableName(ident.Name())
		foundVar := c.FindVariable(getVarName)
		if foundVar != nil {
			return foundVar, nil
		}
		foundNamedExpression := definitions.FindNamedDecoratedExpression(ident)
		if foundNamedExpression == nil {
			return nil, fmt.Errorf("sorry, I don't know what '%v' is", ident)
		}
		fullyQualifiedName := foundNamedExpression.FullyQualifiedName()
		refConstant, _ := c.Constants().AllocateFunctionReferenceConstant(fullyQualifiedName)
		return refConstant, nil
	}

	stringConstant, _ := expr.(*decorated.StringLiteral)
	if stringConstant != nil {
		constant := c.Constants().AllocateStringConstant(stringConstant.Value())
		return constant, nil
	}
	intConstant, _ := expr.(*decorated.IntegerLiteral)
	if intConstant != nil {
		constant := c.Constants().AllocateIntegerConstant(intConstant.Value())
		return constant, nil
	}

	booleanConstant, _ := expr.(*decorated.BooleanLiteral)
	if booleanConstant != nil {
		constant := c.Constants().AllocateBooleanConstant(booleanConstant.Value())
		return constant, nil
	}

	newVar := c.AllocateTempVariable(debugName)
	if genErr := generateExpression(code, newVar, expr, c, definitions); genErr != nil {
		return nil, genErr
	}

	return newVar, nil
}

func generateExpression(code *assembler.Code, target assembler.TargetVariable, expr decorated.DecoratedExpression, context *assembler.Context, definitions *decorator.VariableContext) error {
	switch e := expr.(type) {
	case *decorated.Let:
		return generateLet(code, target, e, context, definitions)

	case *decorated.ArithmeticOperator:
		if e.Left().Type().DecoratedName() == "List" {
			return generateListAppend(code, target, e, context, definitions)
		} else if e.Left().Type().DecoratedName() == "String" {
			return generateStringAppend(code, target, e, context, definitions)
		} else {
			return generateArithmetic(code, target, e, context, definitions)
		}

	case *decorated.BitwiseOperator:
		return generateBitwise(code, target, e, context, definitions)

	case *decorated.BitwiseUnaryOperator:
		return generateUnaryBitwise(code, target, e, context, definitions)

	case *decorated.LogicalUnaryOperator:
		return generateUnaryLogical(code, target, e, context, definitions)

	case *decorated.LogicalOperator:
		return generateLogical(code, target, e, context, definitions)

	case *decorated.BooleanOperator:
		return generateBoolean(code, target, e, context, definitions)

	case *decorated.Lookups:
		return generateLookups(code, target, e, context)

	case *decorated.Case:
		return generateCase(code, target, e, context, definitions)

	case *decorated.RecordLiteral:
		return generateRecordLiteral(code, target, e, context, definitions)

	case *decorated.If:
		return generateIf(code, target, e, context, definitions)

	case *decorated.StringLiteral:
		return generateStringLiteral(code, target, e, context)

	case *decorated.IntegerLiteral:
		return generateIntLiteral(code, target, e, context)

	case *decorated.FixedLiteral:
		return generateFixedLiteral(code, target, e, context)

	case *decorated.ResourceNameLiteral:
		return generateResourceNameLiteral(code, target, e, context)

	case *decorated.BooleanLiteral:
		return generateBoolConstant(code, target, e, context)

	case *decorated.ListLiteral:
		return generateList(code, target, e, context, definitions)

	case *decorated.FunctionCall:
		return generateFunctionCall(code, target, e, context, definitions)

	case *decorated.RecurCall:
		return generateRecurCall(code, e, context, definitions)

	case *decorated.CurryFunction:
		return generateCurry(code, target, e, context, definitions)

	case *decorated.CustomTypeVariantConstructor:
		return generateCustomTypeVariantConstructor(code, target, e, context, definitions)

	case *decorated.GetVariableOrReferenceFunction:
		return generateGetVariable(code, target, e, context)

	case *decorated.ConsOperator:
		return generateListCons(code, target, e, context, definitions)

	case *decorated.AsmConstant:
		return generateAsm(code, target, e, context, definitions)

	case *decorated.RecordConstructorRecord:
		return generateExpression(code, target, e.Expression(), context, definitions)

	case *decorated.RecordConstructor:
		return generateRecordSortedAssignments(code, target, e.SortedAssignments(), context, definitions)

	}

	return fmt.Errorf("generate: unknown node %v %v %v %v", expr, reflect.TypeOf(expr), context, definitions)
}

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedVariableName, f *decorated.FunctionValue, root *assembler.FunctionRootContext, definitions *decorator.VariableContext, lookup typeinfo.TypeLookup, verboseFlag bool) (*Function, error) {
	code := assembler.NewCode()
	funcContext := root.ScopeContext()
	tempVar := root.ReturnVariable()
	for _, parameter := range f.Parameters() {
		paramVarName := assembler.NewVariableName(parameter.Identifier().Name())
		funcContext.AllocateKeepParameterVariable(paramVarName)
	}
	genErr := generateExpression(code, tempVar, f.Expression(), funcContext, definitions)
	if genErr != nil {
		return nil, genErr
	}
	code.Return()
	opcodes, resolveErr := code.Resolve(root, verboseFlag)
	if resolveErr != nil {
		return nil, resolveErr
	}
	if verboseFlag {
		code.PrintOut()
	}

	parameterTypes, _ := f.ForcedFunctionType().ParameterAndReturn()
	parameterCount := uint(len(parameterTypes))

	signature, lookupErr := lookup.Lookup(f.Type())
	if lookupErr != nil {
		signature = -1
		return nil, lookupErr
	}
	functionConstant := NewFunction(fullyQualifiedVariableName, swamppack.TypeRef(signature), parameterCount, uint(root.Layouter().HighestUsedRegisterValue()), root.Constants().Constants(), opcodes)
	return functionConstant, nil
}

func (g *Generator) GenerateAllLocalDefinedFunctions(module *decorated.Module, definitions *decorator.VariableContext,
	lookup typeinfo.TypeLookup, verboseFlag bool) ([]*Function, error) {
	var functionConstants []*Function

	for _, named := range module.Definitions().Definitions() {
		unknownType := named.Expression()

		fullyQualifiedName := module.FullyQualifiedName(named.Identifier())

		maybeFunction, _ := unknownType.(*decorated.FunctionValue)
		if maybeFunction != nil {
			if verboseFlag {
				fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
			}
			rootContext := assembler.NewFunctionRootContext()

			functionConstant, genFuncErr := generateFunction(fullyQualifiedName, maybeFunction, rootContext, definitions, lookup, verboseFlag)
			if genFuncErr != nil {
				return nil, genFuncErr
			}

			if functionConstant == nil {
				panic(fmt.Sprintf("problem %v\n", maybeFunction))
			}

			functionConstants = append(functionConstants, functionConstant)
		} else {
			return nil, fmt.Errorf("generate: unknown type %T", unknownType)
		}
	}

	return functionConstants, nil
}
