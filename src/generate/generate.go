/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate

import (
	"fmt"

	asmcompile "github.com/swamp/assembler/compiler"
	assembler "github.com/swamp/assembler/lib"
	swampopcodeinst "github.com/swamp/opcodes/instruction"
	swamppack "github.com/swamp/pack/lib"

	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/typeinfo"
)

type generateContext struct {
	context     *assembler.Context
	definitions *decorator.VariableContext
	lookup      typeinfo.TypeLookup
}

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
			case assembler.ConstantTypeResourceName:
				packConstant = constantRepo.AddResourceName(subConstant.StringValue())
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
		return swampopcodeinst.BinaryOperatorArithmeticIntPlus
	case decorated.ArithmeticMinus:
		return swampopcodeinst.BinaryOperatorArithmeticIntMinus
	case decorated.ArithmeticMultiply:
		return swampopcodeinst.BinaryOperatorArithmeticIntMultiply
	case decorated.ArithmeticDivide:
		return swampopcodeinst.BinaryOperatorArithmeticIntDivide
	case decorated.ArithmeticAppend:
		return swampopcodeinst.BinaryOperatorArithmeticListAppend
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

func arithmeticToUnaryOperatorType(operatorType decorated.ArithmeticUnaryOperatorType) swampopcodeinst.UnaryOperatorType {
	switch operatorType {
	case decorated.ArithmeticUnaryMinus:
		return swampopcodeinst.UnaryOperatorNegate
	}
	panic("illegal unaryoperator")
}

func bitwiseToBinaryOperatorType(operatorType decorated.BitwiseOperatorType) swampopcodeinst.BinaryOperatorType {
	switch operatorType {
	case decorated.BitwiseAnd:
		return swampopcodeinst.BinaryOperatorBitwiseIntAnd
	case decorated.BitwiseOr:
		return swampopcodeinst.BinaryOperatorBitwiseIntOr
	case decorated.BitwiseXor:
		return swampopcodeinst.BinaryOperatorBitwiseIntXor
	}

	return 0
}

func generateListAppend(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
	// leftVar := context.AllocateTempVariable("arit-left")
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "list-append-left")
	if leftErr != nil {
		return leftErr
	}

	// rightVar := context.AllocateTempVariable("arit-right")
	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "list-append-right")
	if rightErr != nil {
		return rightErr
	}

	code.ListAppend(target, leftVar, rightVar)
	genContext.context.FreeVariableIfNeeded(leftVar)
	genContext.context.FreeVariableIfNeeded(rightVar)

	return nil
}

func generateStringAppend(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "string-append-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "string-append-right")
	if rightErr != nil {
		return rightErr
	}

	code.StringAppend(target, leftVar, rightVar)
	genContext.context.FreeVariableIfNeeded(leftVar)
	genContext.context.FreeVariableIfNeeded(rightVar)

	return nil
}

func generateAsm(code *assembler.Code, target assembler.TargetVariable, asm *decorated.AsmConstant, context *assembler.Context) error {
	compileErr := asmcompile.CompileToCodeAndContext(asm.Asm().Asm(), code, context)
	return compileErr
}

func generateListCons(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ConsOperator, genContext *generateContext) error {
	// leftVar := context.AllocateTempVariable("arit-left")
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "cons-left")
	if leftErr != nil {
		return leftErr
	}

	// rightVar := context.AllocateTempVariable("arit-right")
	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "cons-right")
	if rightErr != nil {
		return rightErr
	}

	code.ListConj(target, leftVar, rightVar)
	genContext.context.FreeVariableIfNeeded(leftVar)
	genContext.context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generateArithmetic(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "arith-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "arit-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := arithmeticToBinaryOperatorType(operator.OperatorType())
	code.BinaryOperator(target, leftVar, rightVar, opcodeBinaryOperator)
	genContext.context.FreeVariableIfNeeded(leftVar)
	genContext.context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generatePipeLeft(code *assembler.Code, target assembler.TargetVariable, operator *decorated.PipeLeftOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.GenerateLeft(), genContext)
	if leftErr != nil {
		return leftErr
	}
	return nil
}

func generatePipeRight(code *assembler.Code, target assembler.TargetVariable, operator *decorated.PipeRightOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.GenerateRight(), genContext)
	if leftErr != nil {
		return leftErr
	}
	return nil
}

func generateUnaryBitwise(code *assembler.Code, target assembler.TargetVariable, operator *decorated.BitwiseUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := bitwiseToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target, leftVar, opcodeUnaryOperatorType)
	genContext.context.FreeVariableIfNeeded(leftVar)
	return nil
}

func generateUnaryLogical(code *assembler.Code, target assembler.TargetVariable, operator *decorated.LogicalUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := logicalToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target, leftVar, opcodeUnaryOperatorType)
	genContext.context.FreeVariableIfNeeded(leftVar)
	return nil
}

func generateUnaryArithmetic(code *assembler.Code, target assembler.TargetVariable, operator *decorated.ArithmeticUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := arithmeticToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target, leftVar, opcodeUnaryOperatorType)
	genContext.context.FreeVariableIfNeeded(leftVar)
	return nil
}

func generateBitwise(code *assembler.Code, target assembler.TargetVariable, operator *decorated.BitwiseOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "bitwise-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := bitwiseToBinaryOperatorType(operator.OperatorType())
	code.BinaryOperator(target, leftVar, rightVar, opcodeBinaryOperator)
	genContext.context.FreeVariableIfNeeded(leftVar)
	genContext.context.FreeVariableIfNeeded(rightVar)
	return nil
}

func generateLogical(code *assembler.Code, target assembler.TargetVariable, operator *decorated.LogicalOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.Left(), genContext)
	if leftErr != nil {
		return leftErr
	}

	codeAlternative := assembler.NewCode()
	rightErr := generateExpression(codeAlternative, target, operator.Right(), genContext)
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

func booleanToBinaryIntOperatorType(operatorType decorated.BooleanOperatorType) swampopcodeinst.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return swampopcodeinst.BinaryOperatorBooleanIntEqual
	case decorated.BooleanNotEqual:
		return swampopcodeinst.BinaryOperatorBooleanIntNotEqual
	case decorated.BooleanLess:
		return swampopcodeinst.BinaryOperatorBooleanIntLess
	case decorated.BooleanLessOrEqual:
		return swampopcodeinst.BinaryOperatorBooleanIntLessOrEqual
	case decorated.BooleanGreater:
		return swampopcodeinst.BinaryOperatorBooleanIntGreater
	case decorated.BooleanGreaterOrEqual:
		return swampopcodeinst.BinaryOperatorBooleanIntGreaterOrEqual
	}

	return 0
}

func booleanToBinaryValueOperatorType(operatorType decorated.BooleanOperatorType) swampopcodeinst.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return swampopcodeinst.BinaryOperatorBooleanValueEqual
	case decorated.BooleanNotEqual:
		return swampopcodeinst.BinaryOperatorBooleanValueNotEqual
	}

	return 0
}

func generateBoolean(code *assembler.Code, target assembler.TargetVariable, operator *decorated.BooleanOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bool-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "bool-right")
	if rightErr != nil {
		return rightErr
	}

	unaliasedTypeLeft := dectype.UnaliasWithResolveInvoker(operator.Left().Type())
	foundPrimitive, _ := unaliasedTypeLeft.(*dectype.PrimitiveAtom)

	opcodeBinaryOperator := booleanToBinaryIntOperatorType(operator.OperatorType())
	if foundPrimitive == nil || foundPrimitive.AtomName() != "Int" {
		opcodeBinaryOperator = booleanToBinaryValueOperatorType(operator.OperatorType())
	}

	code.BinaryOperator(target, leftVar, rightVar, opcodeBinaryOperator)
	genContext.context.FreeVariableIfNeeded(leftVar)
	genContext.context.FreeVariableIfNeeded(rightVar)

	return nil
}

func generateLet(code *assembler.Code, target assembler.TargetVariable, let *decorated.Let, genContext *generateContext) error {
	for _, assignment := range let.Assignments() {
		if len(assignment.LetVariables()) == 1 {
			varName := assembler.NewVariableName(assignment.LetVariables()[0].Name().Name())
			targetVar := genContext.context.AllocateVariable(varName)
			genErr := generateExpression(code, targetVar, assignment.Expression(), genContext)
			if genErr != nil {
				return genErr
			}
		} else {
			sourceVar, sourceErr := generateExpressionWithSourceVar(code, assignment.Expression(), genContext, "tuple split")
			if sourceErr != nil {
				return sourceErr
			}

			var targetVariables []assembler.TargetVariable

			for _, letVariable := range assignment.LetVariables() {
				varName := assembler.NewVariableName(letVariable.Name().Name())
				targetVar := genContext.context.AllocateVariable(varName)
				targetVariables = append(targetVariables, targetVar)
			}
			code.StructSplit(sourceVar, targetVariables)
		}
	}

	codeErr := generateExpression(code, target, let.Consequence(), genContext)
	if codeErr != nil {
		return codeErr
	}

	return nil
}

func generateLookups(code *assembler.Code, target assembler.TargetVariable, lookups *decorated.RecordLookups, genContext *generateContext) error {
	sourceVariable, err := generateExpressionWithSourceVar(code, lookups.Expression(), genContext, "lookups")
	if err != nil {
		return err
	}

	var structLookups []uint8
	for _, indexLookups := range lookups.LookupFields() {
		structLookups = append(structLookups, uint8(indexLookups.Index()))
	}
	code.Lookups(target, sourceVariable, structLookups)

	return nil
}

func generateIf(code *assembler.Code, target assembler.TargetVariable, ifExpr *decorated.If, genContext *generateContext) error {
	conditionVar, testErr := generateExpressionWithSourceVar(code, ifExpr.Condition(), genContext, "if-condition")
	if testErr != nil {
		return testErr
	}

	consequenceCode := assembler.NewCode()
	consequenceContext2 := *genContext
	consequenceContext2.context = genContext.context.MakeScopeContext()
	consErr := generateExpression(consequenceCode, target, ifExpr.Consequence(), &consequenceContext2)
	if consErr != nil {
		return consErr
	}
	consequenceContext2.context.Free()

	alternativeCode := assembler.NewCode()
	alternativeLabel := alternativeCode.Label(nil, "if-alternative")
	alternativeContext2 := *genContext
	alternativeContext2.context = genContext.context.MakeScopeContext()
	altErr := generateExpression(alternativeCode, target, ifExpr.Alternative(), &alternativeContext2)
	if altErr != nil {
		return altErr
	}
	endLabel := alternativeCode.Label(nil, "if-end")
	alternativeContext2.context.Free()

	code.BranchFalse(conditionVar, alternativeLabel)
	genContext.context.FreeVariableIfNeeded(conditionVar)
	consequenceCode.Jump(endLabel)
	code.Copy(consequenceCode)
	code.Copy(alternativeCode)

	return nil
}

func generateGuard(code *assembler.Code, target assembler.TargetVariable, guardExpr *decorated.Guard, genContext *generateContext) error {
	type codeItem struct {
		ConditionVariable     assembler.SourceVariable
		ConditionCode         *assembler.Code
		ConsequenceCode       *assembler.Code
		EndOfConsequenceLabel *assembler.Label
	}

	defaultCode := assembler.NewCode()
	// defaultLabel := defaultCode.Label(nil, "guard-default")
	defaultContext := *genContext
	defaultContext.context = genContext.context.MakeScopeContext()

	altErr := generateExpression(defaultCode, target, guardExpr.DefaultGuard().Expression(), &defaultContext)
	if altErr != nil {
		return altErr
	}
	endLabel := defaultCode.Label(nil, "guard-end")

	var codeItems []codeItem
	for _, item := range guardExpr.Items() {
		conditionCode := assembler.NewCode()
		conditionCodeContext := *genContext
		conditionCodeContext.context = genContext.context.MakeScopeContext()
		conditionVar, testErr := generateExpressionWithSourceVar(conditionCode, item.Condition(), &conditionCodeContext, "guard-condition")
		if testErr != nil {
			return testErr
		}

		consequenceCode := assembler.NewCode()
		consequenceContext := *genContext
		consequenceContext.context = genContext.context.MakeScopeContext()
		consErr := generateExpression(consequenceCode, target, item.Expression(), &consequenceContext)
		if consErr != nil {
			return consErr
		}
		consequenceCode.Jump(endLabel)
		endOfConsequenceLabel := consequenceCode.Label(nil, "guard-end")
		consequenceContext.context.Free()
		codeItem := codeItem{
			ConditionCode: conditionCode, ConditionVariable: conditionVar, ConsequenceCode: consequenceCode,
			EndOfConsequenceLabel: endOfConsequenceLabel,
		}
		codeItems = append(codeItems, codeItem)
	}

	for _, codeItem := range codeItems {
		code.Copy(codeItem.ConditionCode)
		code.BranchFalse(codeItem.ConditionVariable, codeItem.EndOfConsequenceLabel)
		genContext.context.FreeVariableIfNeeded(codeItem.ConditionVariable)
		code.Copy(codeItem.ConsequenceCode)
	}

	code.Copy(defaultCode)

	return nil
}

func generateCaseCustomType(code *assembler.Code, target assembler.TargetVariable, caseExpr *decorated.CaseCustomType, genContext *generateContext) error {
	testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), genContext, "cast-test")
	if testErr != nil {
		return testErr
	}

	var consequences []*assembler.CaseConsequence
	var consequencesCodes []*assembler.Code

	for _, consequence := range caseExpr.Consequences() {
		consequenceContext := *genContext
		consequenceContext.context = genContext.context.MakeScopeContext()

		consequencesCode := assembler.NewCode()

		var parameters []assembler.SourceVariable
		for _, param := range consequence.Parameters() {
			consequenceLabelVariableName := assembler.NewVariableName(param.Identifier().Name())
			paramVariable := consequenceContext.context.AllocateVariable(consequenceLabelVariableName)
			parameters = append(parameters, paramVariable)
		}
		labelVariableName := assembler.NewVariableName(consequence.VariantReference().AstIdentifier().SomeTypeIdentifier().Name())
		caseLabel := consequencesCode.Label(labelVariableName, "case")
		caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), &consequenceContext)
		if caseExprErr != nil {
			return caseExprErr
		}
		asmConsequence := assembler.NewCaseConsequence(uint8(consequence.InternalIndex()), parameters, caseLabel)
		consequences = append(consequences, asmConsequence)

		consequencesCodes = append(consequencesCodes, consequencesCode)

		consequenceContext.context.Free()
	}

	var defaultCase *assembler.CaseConsequence
	if caseExpr.DefaultCase() != nil {
		consequencesCode := assembler.NewCode()
		defaultContext := *genContext
		defaultContext.context = genContext.context.MakeScopeContext()

		decoratedDefault := caseExpr.DefaultCase()
		defaultLabel := consequencesCode.Label(nil, "default")
		caseExprErr := generateExpression(consequencesCode, target, decoratedDefault, &defaultContext)
		if caseExprErr != nil {
			return caseExprErr
		}
		defaultCase = assembler.NewCaseConsequence(0xff, nil, defaultLabel)
		consequencesCodes = append(consequencesCodes, consequencesCode)
		//		endLabel := consequencesBlockCode.Label(nil, "if-end")
		defaultContext.context.Free()
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

	code.Case(testVar, consequences, defaultCase)
	genContext.context.FreeVariableIfNeeded(testVar)
	code.Copy(consequencesBlockCode)

	return nil
}

func generateCasePatternMatching(code *assembler.Code, target assembler.TargetVariable, caseExpr *decorated.CaseForPatternMatching, genContext *generateContext) error {
	testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), genContext, "cast-test")
	if testErr != nil {
		return testErr
	}

	var consequences []*assembler.CaseConsequencePatternMatching
	var consequencesCodes []*assembler.Code

	for _, consequence := range caseExpr.Consequences() {
		consequenceContext := *genContext
		consequenceContext.context = genContext.context.MakeScopeContext()

		consequencesCode := assembler.NewCode()

		literalVariable, literalVariableErr := generateExpressionWithSourceVar(consequencesCode, consequence.Literal(), genContext, "literal")
		if literalVariableErr != nil {
			return literalVariableErr
		}

		labelVariableName := assembler.NewVariableName("a1")
		caseLabel := consequencesCode.Label(labelVariableName, "case")
		caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), &consequenceContext)
		if caseExprErr != nil {
			return caseExprErr
		}
		asmConsequence := assembler.NewCaseConsequencePatternMatching(literalVariable, caseLabel)
		consequences = append(consequences, asmConsequence)

		consequencesCodes = append(consequencesCodes, consequencesCode)

		consequenceContext.context.Free()
	}

	var defaultCase *assembler.CaseConsequencePatternMatching
	if caseExpr.DefaultCase() != nil {
		consequencesCode := assembler.NewCode()
		defaultContext := *genContext
		defaultContext.context = genContext.context.MakeScopeContext()

		decoratedDefault := caseExpr.DefaultCase()
		defaultLabel := consequencesCode.Label(nil, "default")
		caseExprErr := generateExpression(consequencesCode, target, decoratedDefault, &defaultContext)
		if caseExprErr != nil {
			return caseExprErr
		}
		defaultCase = assembler.NewCaseConsequencePatternMatching(nil, defaultLabel)
		consequencesCodes = append(consequencesCodes, consequencesCode)
		//		endLabel := consequencesBlockCode.Label(nil, "if-end")
		defaultContext.context.Free()
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

	code.CasePatternMatching(testVar, consequences, defaultCase)
	genContext.context.FreeVariableIfNeeded(testVar)
	code.Copy(consequencesBlockCode)

	return nil
}

func generateStringLiteral(code *assembler.Code, target assembler.TargetVariable, str *decorated.StringLiteral, context *assembler.Context) error {
	constant := context.Constants().AllocateStringConstant(str.Value())
	code.CopyVariable(target, constant)
	return nil
}

func generateCharacterLiteral(code *assembler.Code, target assembler.TargetVariable, str *decorated.CharacterLiteral, context *assembler.Context) error {
	constant := context.Constants().AllocateIntegerConstant(str.Value())
	code.CopyVariable(target, constant)
	return nil
}

func generateTypeIdLiteral(code *assembler.Code, target assembler.TargetVariable, typeId *decorated.TypeIdLiteral, genContext *generateContext) error {
	indexIntoTypeInformationChunk, err := genContext.lookup.Lookup(typeId.ContainedType())
	if err != nil {
		return err
	}
	constant := genContext.context.Constants().AllocateIntegerConstant(int32(indexIntoTypeInformationChunk))
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

func generateFunctionReference(code *assembler.Code, target assembler.TargetVariable, getVar *decorated.FunctionReference, context *assembler.Context) error {
	varName := assembler.NewVariableName(getVar.Identifier().Name())
	variable := context.FindVariable(varName)
	code.CopyVariable(target, variable)
	return nil
}

func generateConstant(code *assembler.Code, target assembler.TargetVariable, constant *decorated.Constant, context *generateContext) error {
	return generateExpression(code, target, constant.Expression(), context)
}

func generateLocalFunctionParameterReference(code *assembler.Code, target assembler.TargetVariable, getVar *decorated.FunctionParameterReference, context *assembler.Context) error {
	varName := assembler.NewVariableName(getVar.Identifier().Name())
	variable := context.FindVariable(varName)
	code.CopyVariable(target, variable)
	return nil
}

func generateLocalConsequenceParameterReference(code *assembler.Code, target assembler.TargetVariable, getVar *decorated.CaseConsequenceParameterReference, context *assembler.Context) error {
	varName := assembler.NewVariableName(getVar.Identifier().Name())
	variable := context.FindVariable(varName)
	code.CopyVariable(target, variable)
	return nil
}

func generateLetVariableReference(code *assembler.Code, target assembler.TargetVariable, getVar *decorated.LetVariableReference, context *assembler.Context) error {
	varName := assembler.NewVariableName(getVar.LetVariable().Name().Name())
	variable := context.FindVariable(varName)
	code.CopyVariable(target, variable)
	return nil
}

func generateCustomTypeVariantConstructor(code *assembler.Code, target assembler.TargetVariable, constructor *decorated.CustomTypeVariantConstructor, genContext *generateContext) error {
	var arguments []assembler.SourceVariable
	for _, arg := range constructor.Arguments() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, genContext, "customTypeVariantArgs")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg)
	}

	code.CreateEnum(target, constructor.CustomTypeVariantIndex(), arguments)

	return nil
}

func generateCurry(code *assembler.Code, target assembler.TargetVariable, call *decorated.CurryFunction, genContext *generateContext) error {
	var arguments []assembler.SourceVariable
	for _, arg := range call.ArgumentsToSave() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, genContext, "sourceSave")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg)
	}
	functionRegister, functionGenErr := generateExpressionWithSourceVar(code, call.FunctionValue(), genContext, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}
	code.Curry(target, functionRegister, arguments)

	return nil
}

func generateFunctionCall(code *assembler.Code, target assembler.TargetVariable, call *decorated.FunctionCall, genContext *generateContext) error {
	var arguments []assembler.SourceVariable
	for _, arg := range call.Arguments() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, genContext, "arg")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg)
	}

	fn := call.FunctionExpression()

	functionRegister, functionGenErr := generateExpressionWithSourceVar(code, fn, genContext, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}
	code.Call(target, functionRegister, arguments)

	return nil
}

func generateRecurCall(code *assembler.Code, call *decorated.RecurCall, genContext *generateContext) error {
	var arguments []assembler.SourceVariable
	for _, arg := range call.Arguments() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, genContext, "recurarg")
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

func generateRecordSortedAssignments(code *assembler.Code, target assembler.TargetVariable, sortedAssignments []*decorated.RecordLiteralAssignment, genContext *generateContext) error {
	variables := make([]assembler.SourceVariable, len(sortedAssignments))
	for index, assignment := range sortedAssignments {
		debugName := fmt.Sprintf("assign%v", assignment.FieldName())
		assignmentVar, genErr := generateExpressionWithSourceVar(code, assignment.Expression(), genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = assignmentVar
	}
	code.Constructor(target, variables)

	return nil
}

func generateRecordLiteral(code *assembler.Code, target assembler.TargetVariable, record *decorated.RecordLiteral, genContext *generateContext) error {
	if record.RecordTemplate() != nil {
		structToCopyVar, genErr := generateExpressionWithSourceVar(code, record.RecordTemplate(), genContext, "gopher")
		if genErr != nil {
			return genErr
		}
		var updateFields []assembler.UpdateField
		for _, assignment := range record.SortedAssignments() {
			debugName := fmt.Sprintf("update%v", assignment.FieldName())
			assignmentVar, genErr := generateExpressionWithSourceVar(code, assignment.Expression(), genContext, debugName)
			if genErr != nil {
				return genErr
			}
			field := assembler.UpdateField{TargetField: uint8(assignment.Index()), Source: assignmentVar}
			updateFields = append(updateFields, field)
		}
		code.UpdateStruct(target, structToCopyVar, updateFields)
	} else {
		return generateRecordSortedAssignments(code, target, record.SortedAssignments(), genContext)
	}
	return nil
}

func generateList(code *assembler.Code, target assembler.TargetVariable, list *decorated.ListLiteral, genContext *generateContext) error {
	variables := make([]assembler.SourceVariable, len(list.Expressions()))
	for index, expr := range list.Expressions() {
		debugName := fmt.Sprintf("listliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar
	}
	code.ListLiteral(target, variables)
	return nil
}

func generateTuple(code *assembler.Code, target assembler.TargetVariable, tupleLiteral *decorated.TupleLiteral, genContext *generateContext) error {
	variables := make([]assembler.SourceVariable, len(tupleLiteral.Expressions()))
	for index, expr := range tupleLiteral.Expressions() {
		debugName := fmt.Sprintf("tupleliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar
	}
	code.Constructor(target, variables)
	return nil
}

func generateArray(code *assembler.Code, target assembler.TargetVariable, array *decorated.ArrayLiteral, genContext *generateContext) error {
	variables := make([]assembler.SourceVariable, len(array.Expressions()))
	for index, expr := range array.Expressions() {
		debugName := fmt.Sprintf("arrayliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar
	}
	code.Constructor(target, variables)
	return nil
}

func generateExpressionWithSourceVar(code *assembler.Code, expr decorated.Expression, genContext *generateContext, debugName string) (assembler.SourceVariable, error) {
	switch t := expr.(type) {
	case *decorated.StringLiteral:
		constant := genContext.context.Constants().AllocateStringConstant(t.Value())
		return constant, nil
	case *decorated.IntegerLiteral:
		constant := genContext.context.Constants().AllocateIntegerConstant(t.Value())
		return constant, nil
	case *decorated.CharacterLiteral:
		constant := genContext.context.Constants().AllocateIntegerConstant(t.Value())
		return constant, nil
	case *decorated.BooleanLiteral:
		constant := genContext.context.Constants().AllocateBooleanConstant(t.Value())
		return constant, nil
	case *decorated.LetVariableReference:
		parameterReferenceName := assembler.NewVariableName(t.LetVariable().Name().Name())
		return genContext.context.FindVariable(parameterReferenceName), nil
	case *decorated.FunctionParameterReference:
		parameterReferenceName := assembler.NewVariableName(t.Identifier().Name())
		return genContext.context.FindVariable(parameterReferenceName), nil
	case *decorated.FunctionReference:
		ident := t.Identifier()
		functionReferenceName := assembler.NewVariableName(ident.Name())
		foundVar := genContext.context.FindVariable(functionReferenceName)
		if foundVar != nil {
			return foundVar, nil
		}
		foundNamedExpression := genContext.definitions.FindScopedNamedDecoratedExpressionScopedOrNormal(ident)
		if foundNamedExpression == nil {
			return nil, fmt.Errorf("sorry, I don't know what '%v' is %v", ident, ident.FetchPositionLength())
		}
		fullyQualifiedName := foundNamedExpression.FullyQualifiedName()
		refConstant, _ := genContext.context.Constants().AllocateFunctionReferenceConstant(fullyQualifiedName)
		return refConstant, nil
	}

	newVar := genContext.context.AllocateTempVariable(debugName)
	if genErr := generateExpression(code, newVar, expr, genContext); genErr != nil {
		return nil, genErr
	}

	return newVar, nil
}

func isIntLike(typeToCheck dtype.Type) bool {
	unaliasType := dectype.UnaliasWithResolveInvoker(typeToCheck)
	primitiveAtom, _ := unaliasType.(*dectype.PrimitiveAtom)
	if primitiveAtom == nil {
		return false
	}
	name := primitiveAtom.AtomName()

	return name == "Int" || name == "Fixed" || name == "Char"
}

func isListLike(typeToCheck dtype.Type) bool {
	unaliasType := dectype.UnaliasWithResolveInvoker(typeToCheck)
	primitiveAtom, _ := unaliasType.(*dectype.PrimitiveAtom)
	if primitiveAtom == nil {
		return false
	}
	name := primitiveAtom.PrimitiveName().Name()
	return name == "List"
}

func generateExpression(code *assembler.Code, target assembler.TargetVariable, expr decorated.Expression, genContext *generateContext) error {
	//	log.Printf("gen expr:%T (%v)\n", expr, expr)
	switch e := expr.(type) {
	case *decorated.Let:
		return generateLet(code, target, e, genContext)

	case *decorated.ArithmeticOperator:
		{
			if isListLike(e.Left().Type()) && e.OperatorType() == decorated.ArithmeticAppend {
				return generateListAppend(code, target, e, genContext)
			} else if e.Left().Type().DecoratedName() == "String" && e.OperatorType() == decorated.ArithmeticAppend {
				return generateStringAppend(code, target, e, genContext)
			} else if isIntLike(e.Left().Type()) {
				return generateArithmetic(code, target, e, genContext)
			} else {
				return fmt.Errorf("Cant generate arithmetic for type:%v <-> %v (%v)", e.Left().Type(), e.Right().Type(), e.OperatorType())
			}
		}

	case *decorated.BitwiseOperator:
		return generateBitwise(code, target, e, genContext)

	case *decorated.BitwiseUnaryOperator:
		return generateUnaryBitwise(code, target, e, genContext)

	case *decorated.LogicalUnaryOperator:
		return generateUnaryLogical(code, target, e, genContext)

	case *decorated.ArithmeticUnaryOperator:
		return generateUnaryArithmetic(code, target, e, genContext)

	case *decorated.LogicalOperator:
		return generateLogical(code, target, e, genContext)

	case *decorated.BooleanOperator:
		return generateBoolean(code, target, e, genContext)

	case *decorated.PipeLeftOperator:
		return generatePipeLeft(code, target, e, genContext)

	case *decorated.PipeRightOperator:
		return generatePipeRight(code, target, e, genContext)

	case *decorated.RecordLookups:
		return generateLookups(code, target, e, genContext)

	case *decorated.CaseCustomType:
		return generateCaseCustomType(code, target, e, genContext)

	case *decorated.CaseForPatternMatching:
		return generateCasePatternMatching(code, target, e, genContext)

	case *decorated.RecordLiteral:
		return generateRecordLiteral(code, target, e, genContext)

	case *decorated.If:
		return generateIf(code, target, e, genContext)

	case *decorated.Guard:
		return generateGuard(code, target, e, genContext)

	case *decorated.StringLiteral:
		return generateStringLiteral(code, target, e, genContext.context)

	case *decorated.CharacterLiteral:
		return generateCharacterLiteral(code, target, e, genContext.context)

	case *decorated.TypeIdLiteral:
		return generateTypeIdLiteral(code, target, e, genContext)

	case *decorated.IntegerLiteral:
		return generateIntLiteral(code, target, e, genContext.context)

	case *decorated.FixedLiteral:
		return generateFixedLiteral(code, target, e, genContext.context)

	case *decorated.ResourceNameLiteral:
		return generateResourceNameLiteral(code, target, e, genContext.context)

	case *decorated.BooleanLiteral:
		return generateBoolConstant(code, target, e, genContext.context)

	case *decorated.ListLiteral:
		return generateList(code, target, e, genContext)

	case *decorated.TupleLiteral:
		return generateTuple(code, target, e, genContext)

	case *decorated.ArrayLiteral:
		return generateArray(code, target, e, genContext)

	case *decorated.FunctionCall:
		return generateFunctionCall(code, target, e, genContext)

	case *decorated.RecurCall:
		return generateRecurCall(code, e, genContext)

	case *decorated.CurryFunction:
		return generateCurry(code, target, e, genContext)

	case *decorated.StringInterpolation:
		return generateExpression(code, target, e.Expression(), genContext)

	case *decorated.CustomTypeVariantConstructor:
		return generateCustomTypeVariantConstructor(code, target, e, genContext)

	case *decorated.FunctionReference:
		return generateFunctionReference(code, target, e, genContext.context)

	case *decorated.Constant:
		return generateConstant(code, target, e, genContext)

	case *decorated.FunctionParameterReference:
		return generateLocalFunctionParameterReference(code, target, e, genContext.context)

	case *decorated.CaseConsequenceParameterReference:
		return generateLocalConsequenceParameterReference(code, target, e, genContext.context)

	case *decorated.LetVariableReference:
		return generateLetVariableReference(code, target, e, genContext.context)

	case *decorated.ConsOperator:
		return generateListCons(code, target, e, genContext)

	case *decorated.AsmConstant:
		return generateAsm(code, target, e, genContext.context)

	case *decorated.RecordConstructorRecord:
		return generateExpression(code, target, e.Expression(), genContext)

	case *decorated.RecordConstructor:
		return generateRecordSortedAssignments(code, target, e.SortedAssignments(), genContext)
	}

	return fmt.Errorf("generate: unknown node %T %v %v", expr, expr, genContext)
}

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedVariableName, f *decorated.FunctionValue, root *assembler.FunctionRootContext, definitions *decorator.VariableContext, lookup typeinfo.TypeLookup, verboseFlag bool) (*Function, error) {
	code := assembler.NewCode()
	funcContext := root.ScopeContext()
	tempVar := root.ReturnVariable()
	for _, parameter := range f.Parameters() {
		paramVarName := assembler.NewVariableName(parameter.Identifier().Name())
		funcContext.AllocateKeepParameterVariable(paramVarName)
	}
	genContext := &generateContext{
		context:     funcContext,
		definitions: definitions,
		lookup:      lookup,
	}
	genErr := generateExpression(code, tempVar, f.Expression(), genContext)
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
