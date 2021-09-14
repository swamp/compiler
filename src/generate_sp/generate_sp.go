/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_sp

import (
	"fmt"

	asmcompile "github.com/swamp/assembler/compiler"
	"github.com/swamp/compiler/src/assembler_sp"
	decorator "github.com/swamp/compiler/src/decorated/convert"
	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/instruction_sp"
	"github.com/swamp/compiler/src/typeinfo"
	"github.com/swamp/compiler/src/verbosity"
	swamppack "github.com/swamp/pack/lib"
)

type AnyPosAndRange interface {
	getPosition() uint32
	getSize() uint32
}

type Context struct {
	startMemoryConstants *StartMemoryConstants
	constants            *assembler_sp.Constants
	functionVariables    *assembler_sp.FunctionVariables
	stackMemory          *assembler_sp.StackMemoryMapper
}

func (c *Context) StartMemoryConstants() *StartMemoryConstants {
	return c.startMemoryConstants
}

func (c *Context) Constants() *assembler_sp.Constants {
	return c.constants
}

type generateContext struct {
	context     *Context
	definitions *decorator.VariableContext
	lookup      typeinfo.TypeLookup
}

type Function struct {
	name      *decorated.FullyQualifiedPackageVariableName
	signature swamppack.TypeRef
	opcodes   []byte
}

type ExternalFunction struct {
	name           *decorated.FullyQualifiedPackageVariableName
	signature      swamppack.TypeRef
	parameterCount uint
}

func NewFunction(fullyQualifiedName *decorated.FullyQualifiedPackageVariableName, signature swamppack.TypeRef,
	opcodes []byte) *Function {
	f := &Function{
		name: fullyQualifiedName, signature: signature, opcodes: opcodes,
	}

	return f
}

func NewExternalFunction(fullyQualifiedName *decorated.FullyQualifiedPackageVariableName,
	signature swamppack.TypeRef, parameterCount uint) *ExternalFunction {
	f := &ExternalFunction{name: fullyQualifiedName, signature: signature, parameterCount: parameterCount}

	return f
}

func (f *Function) String() string {
	return fmt.Sprintf("[function %v %v %v %v]", f.name, f.signature)
}

func (f *Function) Opcodes() []byte {
	return f.opcodes
}

type Generator struct {
	code *assembler_sp.Code
}

func NewGenerator() *Generator {
	return &Generator{code: assembler_sp.NewCode()}
}

func arithmeticToBinaryOperatorType(operatorType decorated.ArithmeticOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.ArithmeticPlus:
		return instruction_sp.BinaryOperatorArithmeticIntPlus
	case decorated.ArithmeticCons:
		panic("cons not handled here")
	case decorated.ArithmeticMinus:
		return instruction_sp.BinaryOperatorArithmeticIntMinus
	case decorated.ArithmeticMultiply:
		return instruction_sp.BinaryOperatorArithmeticIntMultiply
	case decorated.ArithmeticDivide:
		return instruction_sp.BinaryOperatorArithmeticIntDivide
	case decorated.ArithmeticAppend:
		return instruction_sp.BinaryOperatorArithmeticListAppend
	case decorated.ArithmeticFixedMultiply:
		return instruction_sp.BinaryOperatorArithmeticFixedMultiply
	case decorated.ArithmeticFixedDivide:
		return instruction_sp.BinaryOperatorArithmeticFixedDivide
	}

	panic("unknown binary operator")
}

func bitwiseToUnaryOperatorType(operatorType decorated.BitwiseUnaryOperatorType) instruction_sp.UnaryOperatorType {
	switch operatorType {
	case decorated.BitwiseUnaryNot:
		return instruction_sp.UnaryOperatorBitwiseNot
	}

	panic("illegal unaryoperator")
}

func logicalToUnaryOperatorType(operatorType decorated.LogicalUnaryOperatorType) instruction_sp.UnaryOperatorType {
	switch operatorType {
	case decorated.LogicalUnaryNot:
		return instruction_sp.UnaryOperatorNot
	}

	panic("illegal unaryoperator")
}

func arithmeticToUnaryOperatorType(operatorType decorated.ArithmeticUnaryOperatorType) instruction_sp.UnaryOperatorType {
	switch operatorType {
	case decorated.ArithmeticUnaryMinus:
		return instruction_sp.UnaryOperatorNegate
	}

	panic("illegal unaryoperator")
}

func bitwiseToBinaryOperatorType(operatorType decorated.BitwiseOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BitwiseAnd:
		return instruction_sp.BinaryOperatorBitwiseIntAnd
	case decorated.BitwiseOr:
		return instruction_sp.BinaryOperatorBitwiseIntOr
	case decorated.BitwiseXor:
		return instruction_sp.BinaryOperatorBitwiseIntXor
	case decorated.BitwiseNot:
		return 0
		// return opcode_sp.BinaryOperatorBitwiseIntNot
	}

	return 0
}

func generateListAppend(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
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

	code.ListAppend(target.Pos, leftVar.Pos, rightVar.Pos)

	return nil
}

func generateStringAppend(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "string-append-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "string-append-right")
	if rightErr != nil {
		return rightErr
	}

	code.StringAppend(target.Pos, leftVar.Pos, rightVar.Pos)

	return nil
}

func generateAsm(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, asm *decorated.AsmConstant, context *assembler_sp.Context) error {
	compileErr := asmcompile.CompileToCodeAndContext(asm.Asm().Asm(), code, context)
	return compileErr
}

func generateListCons(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ConsOperator, genContext *generateContext) error {
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

	code.ListConj(target.Pos, leftVar.Pos, rightVar.Pos)

	return nil
}

func generateArithmetic(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "arith-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "arit-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := arithmeticToBinaryOperatorType(operator.OperatorType())
	code.BinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)

	return nil
}

func generatePipeLeft(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.PipeLeftOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.GenerateLeft(), genContext)
	if leftErr != nil {
		return leftErr
	}
	return nil
}

func generatePipeRight(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.PipeRightOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.GenerateRight(), genContext)
	if leftErr != nil {
		return leftErr
	}
	return nil
}

func generateUnaryBitwise(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.BitwiseUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := bitwiseToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target.Pos, leftVar.Pos, opcodeUnaryOperatorType)

	return nil
}

func generateUnaryLogical(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.LogicalUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := logicalToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target.Pos, leftVar.Pos, opcodeUnaryOperatorType)
	return nil
}

func generateUnaryArithmetic(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.ArithmeticUnaryOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}
	opcodeUnaryOperatorType := arithmeticToUnaryOperatorType(operator.OperatorType())
	code.UnaryOperator(target.Pos, leftVar.Pos, opcodeUnaryOperatorType)

	return nil
}

func generateBitwise(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.BitwiseOperator, genContext *generateContext) error {
	leftVar, leftErr := generateExpressionWithSourceVar(code, operator.Left(), genContext, "bitwise-left")
	if leftErr != nil {
		return leftErr
	}

	rightVar, rightErr := generateExpressionWithSourceVar(code, operator.Right(), genContext, "bitwise-right")
	if rightErr != nil {
		return rightErr
	}

	opcodeBinaryOperator := bitwiseToBinaryOperatorType(operator.OperatorType())
	code.BinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)

	return nil
}

func generateLogical(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.LogicalOperator, genContext *generateContext) error {
	leftErr := generateExpression(code, target, operator.Left(), genContext)
	if leftErr != nil {
		return leftErr
	}

	codeAlternative := assembler_sp.NewCode()
	rightErr := generateExpression(codeAlternative, target, operator.Right(), genContext)
	if rightErr != nil {
		return rightErr
	}
	afterLabel := codeAlternative.Label(nil, "after-alternative")

	if operator.OperatorType() == decorated.LogicalAnd {
		code.BranchFalse(targetToSourceStackPosRange(target).Pos, afterLabel)
	} else if operator.OperatorType() == decorated.LogicalOr {
		code.BranchTrue(targetToSourceStackPosRange(target).Pos, afterLabel)
	}
	code.Copy(codeAlternative)

	return nil
}

func booleanToBinaryIntOperatorType(operatorType decorated.BooleanOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return instruction_sp.BinaryOperatorBooleanIntEqual
	case decorated.BooleanNotEqual:
		return instruction_sp.BinaryOperatorBooleanIntNotEqual
	case decorated.BooleanLess:
		return instruction_sp.BinaryOperatorBooleanIntLess
	case decorated.BooleanLessOrEqual:
		return instruction_sp.BinaryOperatorBooleanIntLessOrEqual
	case decorated.BooleanGreater:
		return instruction_sp.BinaryOperatorBooleanIntGreater
	case decorated.BooleanGreaterOrEqual:
		return instruction_sp.BinaryOperatorBooleanIntGreaterOrEqual
	}

	return 0
}

func booleanToBinaryValueOperatorType(operatorType decorated.BooleanOperatorType) instruction_sp.BinaryOperatorType {
	switch operatorType {
	case decorated.BooleanEqual:
		return instruction_sp.BinaryOperatorBooleanValueEqual
	case decorated.BooleanNotEqual:
		return instruction_sp.BinaryOperatorBooleanValueNotEqual
	}

	return 0
}

func generateBoolean(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, operator *decorated.BooleanOperator, genContext *generateContext) error {
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

	code.BinaryOperator(target.Pos, leftVar.Pos, rightVar.Pos, opcodeBinaryOperator)

	return nil
}

func generateLet(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, let *decorated.Let, genContext *generateContext) error {
	for _, assignment := range let.Assignments() {
		sourceVar, sourceErr := generateExpressionWithSourceVar(code, assignment.Expression(), genContext, "let source")
		if sourceErr != nil {
			return sourceErr
		}

		if len(assignment.LetVariables()) == 1 {
			firstVar := assignment.LetVariables()[0]
			genContext.context.functionVariables.DefineVariable(firstVar.Name().Name(), sourceVar)
		} else {
		}
	}

	codeErr := generateExpression(code, target, let.Consequence(), genContext)
	if codeErr != nil {
		return codeErr
	}

	return nil
}

func generateLookups(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, lookups *decorated.RecordLookups,
	genContext *generateContext) error {
	startOfStruct, err := generateExpressionWithSourceVar(code, lookups.Expression(), genContext, "lookups")
	if err != nil {
		return err
	}

	indexOffset := uint(0)

	var lastLookup decorated.LookupField
	for _, indexLookup := range lookups.LookupFields() {
		indexOffset += indexLookup.MemoryOffset()
		lastLookup = indexLookup
	}

	sourcePosRange := assembler_sp.SourceStackPosRange{
		Pos:  assembler_sp.SourceStackPos(uint(startOfStruct.Pos) + indexOffset),
		Size: assembler_sp.SourceStackRange(lastLookup.MemorySize()),
	}

	code.CopyMemory(target.Pos, sourcePosRange)

	return nil
}

func generateIf(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, ifExpr *decorated.If, genContext *generateContext) error {
	conditionVar, testErr := generateExpressionWithSourceVar(code, ifExpr.Condition(), genContext, "if-condition")
	if testErr != nil {
		return testErr
	}

	consequenceCode := assembler_sp.NewCode()
	consequenceContext2 := *genContext
	consequenceContext2.context = genContext.context.MakeScopeContext()

	consErr := generateExpression(consequenceCode, target, ifExpr.Consequence(), &consequenceContext2)
	if consErr != nil {
		return consErr
	}

	consequenceContext2.context.Free()

	alternativeCode := assembler_sp.NewCode()
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

	consequenceCode.Jump(endLabel)
	code.Copy(consequenceCode)
	code.Copy(alternativeCode)

	return nil
}

func generateGuard(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, guardExpr *decorated.Guard, genContext *generateContext) error {
	type codeItem struct {
		ConditionVariable     assembler_sp.SourceStackPosRange
		ConditionCode         *assembler_sp.Code
		ConsequenceCode       *assembler_sp.Code
		EndOfConsequenceLabel *assembler_sp.Label
	}

	defaultCode := assembler_sp.NewCode()
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
		conditionCode := assembler_sp.NewCode()
		conditionCodeContext := *genContext
		conditionCodeContext.context = genContext.context.MakeScopeContext()

		conditionVar, testErr := generateExpressionWithSourceVar(conditionCode,
			item.Condition(), &conditionCodeContext, "guard-condition")
		if testErr != nil {
			return testErr
		}

		consequenceCode := assembler_sp.NewCode()
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

		code.Copy(codeItem.ConsequenceCode)
	}

	code.Copy(defaultCode)

	return nil
}

func generateCaseCustomType(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, caseExpr *decorated.CaseCustomType, genContext *generateContext) error {
	testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), genContext, "cast-test")
	if testErr != nil {
		return testErr
	}

	var consequences []*assembler_sp.CaseConsequence

	var consequencesCodes []*assembler_sp.Code

	for _, consequence := range caseExpr.Consequences() {
		consequenceContext := *genContext
		consequenceContext.context = genContext.context.MakeScopeContext()

		consequencesCode := assembler_sp.NewCode()

		var parameters []assembler_sp.SourceVariable

		for _, param := range consequence.Parameters() {
			consequenceLabelVariableName := assembler_sp.NewVariableName(param.Identifier().Name())
			paramVariable := consequenceContext.context.AllocateVariable(consequenceLabelVariableName)
			parameters = append(parameters, paramVariable)
		}

		labelVariableName := assembler_sp.NewVariableName(
			consequence.VariantReference().AstIdentifier().SomeTypeIdentifier().Name())

		caseLabel := consequencesCode.Label(labelVariableName, "case")

		caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), &consequenceContext)
		if caseExprErr != nil {
			return caseExprErr
		}

		asmConsequence := assembler_sp.NewCaseConsequence(uint8(consequence.InternalIndex()), parameters, caseLabel)

		consequences = append(consequences, asmConsequence)

		consequencesCodes = append(consequencesCodes, consequencesCode)

		consequenceContext.context.Free()
	}

	var defaultCase *assembler_sp.CaseConsequence

	if caseExpr.DefaultCase() != nil {
		consequencesCode := assembler_sp.NewCode()
		defaultContext := *genContext
		defaultContext.context = genContext.context.MakeScopeContext()

		decoratedDefault := caseExpr.DefaultCase()
		defaultLabel := consequencesCode.Label(nil, "default")
		caseExprErr := generateExpression(consequencesCode, target, decoratedDefault, &defaultContext)
		if caseExprErr != nil {
			return caseExprErr
		}
		defaultCase = assembler_sp.NewCaseConsequence(0xff, nil, defaultLabel)
		consequencesCodes = append(consequencesCodes, consequencesCode)
		//		endLabel := consequencesBlockCode.Label(nil, "if-end")
		defaultContext.context.Free()
	}

	consequencesBlockCode := assembler_sp.NewCode()

	lastConsequnce := consequencesCodes[len(consequencesCodes)-1]

	labelVariableEndName := assembler_sp.NewVariableName("case end")
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

	code.Copy(consequencesBlockCode)

	return nil
}

func generateCasePatternMatching(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, caseExpr *decorated.CaseForPatternMatching, genContext *generateContext) error {
	testVar, testErr := generateExpressionWithSourceVar(code, caseExpr.Test(), genContext, "cast-test")
	if testErr != nil {
		return testErr
	}

	var consequences []*assembler_sp.CaseConsequencePatternMatching

	var consequencesCodes []*assembler_sp.Code

	for _, consequence := range caseExpr.Consequences() {
		consequenceContext := *genContext
		consequenceContext.context = genContext.context.MakeScopeContext()

		consequencesCode := assembler_sp.NewCode()

		literalVariable, literalVariableErr := generateExpressionWithSourceVar(consequencesCode,
			consequence.Literal(), genContext, "literal")
		if literalVariableErr != nil {
			return literalVariableErr
		}

		labelVariableName := assembler_sp.NewVariableName("a1")
		caseLabel := consequencesCode.Label(labelVariableName, "case")

		caseExprErr := generateExpression(consequencesCode, target, consequence.Expression(), &consequenceContext)
		if caseExprErr != nil {
			return caseExprErr
		}

		asmConsequence := assembler_sp.NewCaseConsequencePatternMatching(literalVariable, caseLabel)
		consequences = append(consequences, asmConsequence)

		consequencesCodes = append(consequencesCodes, consequencesCode)

		consequenceContext.context.Free()
	}

	var defaultCase *assembler_sp.CaseConsequencePatternMatching

	if caseExpr.DefaultCase() != nil {
		consequencesCode := assembler_sp.NewCode()
		defaultContext := *genContext
		defaultContext.context = genContext.context.MakeScopeContext()

		decoratedDefault := caseExpr.DefaultCase()
		defaultLabel := consequencesCode.Label(nil, "default")
		caseExprErr := generateExpression(consequencesCode, target, decoratedDefault, &defaultContext)
		if caseExprErr != nil {
			return caseExprErr
		}
		defaultCase = assembler_sp.NewCaseConsequencePatternMatching(nil, defaultLabel)
		consequencesCodes = append(consequencesCodes, consequencesCode)
		//		endLabel := consequencesBlockCode.Label(nil, "if-end")
		defaultContext.context.Free()
	}

	consequencesBlockCode := assembler_sp.NewCode()

	lastConsequnce := consequencesCodes[len(consequencesCodes)-1]

	labelVariableEndName := assembler_sp.NewVariableName("case end")
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

	code.Copy(consequencesBlockCode)

	return nil
}

func generateStringLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.StringLiteral,
	context *assembler_sp.Context) error {
	constant := context.Constants().AllocateStringConstant(str.Value())
	code.LoadZeroMemoryPointer(target.Pos, constant.PosRange().Position)
	return nil
}

func generateCharacterLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, str *decorated.CharacterLiteral,
	context *assembler_sp.Context) error {
	code.LoadRune(target.Pos, uint8(str.Value()))
	return nil
}

func generateTypeIdLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, typeId *decorated.TypeIdLiteral,
	genContext *generateContext) error {
	return nil
}

func generateIntLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, integer *decorated.IntegerLiteral,
	context *assembler_sp.Context) error {
	code.LoadInteger(target.Pos, integer.Value())
	return nil
}

func generateFixedLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, fixed *decorated.FixedLiteral,
	context *assembler_sp.Context) error {
	code.LoadInteger(target.Pos, fixed.Value())
	return nil
}

func generateResourceNameLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	resourceName *decorated.ResourceNameLiteral, context *assembler_sp.Context) error {
	return nil
}

func generateConstant(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	constant *decorated.Constant, context *generateContext) error {
	return generateExpression(code, target, constant.Expression(), context)
}

func generateLocalFunctionParameterReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.FunctionParameterReference, context *assembler_sp.Context) error {
	varName := assembler_sp.NewVariableName(getVar.Identifier().Name())
	variable := context.FindVariable(varName)
	code.CopyConstantPointer(target, variable)
	return nil
}

func generateLocalConsequenceParameterReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.CaseConsequenceParameterReference, context *assembler_sp.Context) error {
	varName := assembler_sp.NewVariableName(getVar.Identifier().Name())
	variable := context.FindVariable(varName)
	code.CopyVariable(target, variable)
	return nil
}

func generateLetVariableReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.LetVariableReference, context *assembler_sp.Context) error {
	varName := assembler_sp.NewVariableName(getVar.LetVariable().Name().Name())
	variable := context.FindVariable(varName)
	code.CopyVariable(target, variable)
	return nil
}

func generateCustomTypeVariantConstructor(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	constructor *decorated.CustomTypeVariantConstructor, genContext *generateContext) error {
	var arguments []assembler_sp.SourceVariable
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

func generateCurry(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, call *decorated.CurryFunction,
	genContext *generateContext) error {
	var arguments []assembler_sp.SourceVariable

	for _, arg := range call.ArgumentsToSave() {
		argReg, argRegErr := generateExpressionWithSourceVar(code, arg, genContext, "sourceSave")
		if argRegErr != nil {
			return argRegErr
		}
		arguments = append(arguments, argReg)
	}

	functionRegister, functionGenErr := generateExpressionWithSourceVar(code,
		call.FunctionValue(), genContext, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}

	indexIntoTypeInformationChunk, lookupErr := genContext.lookup.Lookup(call.Type())
	if lookupErr != nil {
		return lookupErr
	}

	code.Curry(target, uint16(indexIntoTypeInformationChunk), functionRegister, arguments)

	return nil
}

func getMemorySizeAndAlignment(p dtype.Type) (uint, uint32) {
	unaliased := dectype.Unalias(p)
	switch t := unaliased.(type) {
	case *dectype.RecordAtom:
	}
}

func generateFunctionCall(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, call *decorated.FunctionCall,
	genContext *generateContext) error {
	functionType := dectype.Unalias(call.FunctionExpression().Type())
	functionAtom, wasFunctionAtom := functionType.(*dectype.FunctionAtom)

	if !wasFunctionAtom {
		return fmt.Errorf("this is not a function atom %T", functionType)
	}

	fn := call.FunctionExpression()
	functionRegister, functionGenErr := generateExpressionWithSourceVar(code, fn, genContext, "functioncall")
	if functionGenErr != nil {
		return functionGenErr
	}

	var arguments []assembler_sp.TargetStackPosRange
	for index, arg := range call.Arguments() {
		memorySize, alignment := getMemorySizeAndAlignment(arg.Type())
		arguments[index] = genContext.context.stackMemory.Allocate(memorySize, alignment, fmt.Sprintf("arg %d", index))
	}

	for index, arg := range call.Arguments() {
		functionArgType := functionAtom.FunctionParameterTypes()[index]
		functionArgTypeUnalias := dectype.Unalias(functionArgType)

		argReg := arguments[index]
		argRegErr := generateExpression(code, argReg, arg, genContext)
		if argRegErr != nil {
			return argRegErr
		}

		isAny := dectype.IsAny(functionArgTypeUnalias)
		if isAny { // arg.NeedsTypeId() {
			/*
				constant, err := generateTypeIdConstant(arg.Type(), genContext)
				if err != nil {
					return err
				}

				tempAnyConstructor := genContext.context.AllocateTempVariable("anyConstructor")
				code.Constructor(tempAnyConstructor, []assembler_sp.SourceVariable{constant, argReg})

				argReg = tempAnyConstructor

				tempVariables = append(tempVariables, tempAnyConstructor)

			*/
		}

		arguments = append(arguments, argReg)
	}

	code.Call(functionRegister.Pos, arguments[0].Pos)

	return nil
}

func generateRecurCall(code *assembler_sp.Code, call *decorated.RecurCall, genContext *generateContext) error {
	code.Recur()

	return nil
}

func generateBoolConstant(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	test *decorated.BooleanLiteral, context *assembler_sp.Context) error {
	constant := context.Constants().AllocateBooleanConstant(test.Value())
	code.CopyVariable(target, constant)
	return nil
}

func generateRecordSortedAssignments(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	sortedAssignments []*decorated.RecordLiteralAssignment, genContext *generateContext) error {
	variables := make([]assembler_sp.SourceVariable, len(sortedAssignments))
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

func generateRecordLiteral(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	record *decorated.RecordLiteral, genContext *generateContext) error {
	if record.RecordTemplate() != nil {
		structToCopyVar, genErr := generateExpressionWithSourceVar(code, record.RecordTemplate(), genContext, "gopher")
		if genErr != nil {
			return genErr
		}
		var updateFields []assembler_sp.UpdateField
		for _, assignment := range record.SortedAssignments() {
			debugName := fmt.Sprintf("update%v", assignment.FieldName())
			assignmentVar, genErr := generateExpressionWithSourceVar(code, assignment.Expression(), genContext, debugName)
			if genErr != nil {
				return genErr
			}
			field := assembler_sp.UpdateField{TargetField: uint8(assignment.Index()), Source: assignmentVar}
			updateFields = append(updateFields, field)
		}
		code.UpdateStruct(target, structToCopyVar, updateFields)
	} else {
		return generateRecordSortedAssignments(code, target, record.SortedAssignments(), genContext)
	}
	return nil
}

func generateList(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	list *decorated.ListLiteral, genContext *generateContext) error {
	variables := make([]assembler_sp.SourceStackPos, len(list.Expressions()))
	for index, expr := range list.Expressions() {
		debugName := fmt.Sprintf("listliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar.Pos
	}
	primitive, _ := list.Type().(*dectype.PrimitiveAtom)
	itemSize, _ := getMemorySizeAndAlignment(primitive.Next())
	code.ListLiteral(target.Pos, variables, assembler_sp.StackRange(itemSize))
	return nil
}

func generateTuple(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	tupleLiteral *decorated.TupleLiteral, genContext *generateContext) error {
	variables := make([]assembler_sp.SourceStackPos, len(tupleLiteral.Expressions()))

	tuplePointer := genContext.context.stackMemory.Allocate(tupleLiteral.TupleType().MemorySize(),
		tupleLiteral.TupleType().MemoryAlignment(), "tuple")
	for index, expr := range tupleLiteral.Expressions() {
		tupleIndex := tupleLiteral.TupleType().
		debugName := fmt.Sprintf("tupleliteral%v", index)
		exprVar, genErr := generateExpression(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar
	}

	return nil
}

func generateArray(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	array *decorated.ArrayLiteral, genContext *generateContext) error {
	variables := make([]assembler_sp.SourceVariable, len(array.Expressions()))
	for index, expr := range array.Expressions() {
		debugName := fmt.Sprintf("arrayliteral%v", index)
		exprVar, genErr := generateExpressionWithSourceVar(code, expr, genContext, debugName)
		if genErr != nil {
			return genErr
		}
		variables[index] = exprVar
	}
	code.CreateArray(target, variables)
	return nil
}

const (
	PointerSize  = 8
	PointerAlign = 8
)

func targetToSourceStackPosRange(functionPointer assembler_sp.TargetStackPosRange) assembler_sp.SourceStackPosRange {
	sourcePosRange := assembler_sp.SourceStackPosRange{
		Pos:  assembler_sp.SourceStackPos(functionPointer.Pos),
		Size: assembler_sp.SourceStackRange(functionPointer.Size),
	}

	return sourcePosRange
}

func constantToSourceStackPosRange(code *assembler_sp.Code, stackMemory *assembler_sp.StackMemoryMapper, constant *assembler_sp.Constant) (assembler_sp.SourceStackPosRange, error) {
	functionPointer := stackMemory.Allocate(PointerSize, PointerAlign, "functionReference")
	code.LoadZeroMemoryPointer(functionPointer.Pos, constant.PosRange().Position)

	return targetToSourceStackPosRange(functionPointer), nil
}

func handleFunctionReference(code *assembler_sp.Code,
	t *decorated.FunctionReference,
	stackMemory *assembler_sp.StackMemoryMapper,
	constants *assembler_sp.Constants) (assembler_sp.SourceStackPosRange, error) {
	ident := t.Identifier()
	functionReferenceName := assembler_sp.VariableName(ident.Name())
	foundConstant := constants.FindFunction(functionReferenceName)
	if foundConstant == nil {
		return assembler_sp.SourceStackPosRange{}, fmt.Errorf("couldn't find it %v", t)
	}

	return constantToSourceStackPosRange(code, stackMemory, foundConstant)
}

func generateFunctionReference(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange,
	getVar *decorated.FunctionReference, context *assembler_sp.Context) error {
	varName := assembler_sp.VariableName(getVar.Identifier().Name())
	functionConstant := context.Constants().FindFunction(varName)
	code.LoadZeroMemoryPointer(target.Pos, functionConstant.PosRange().Position)
	return nil
}

const (
	SizeofSwampInt  = 4
	SizeofSwampRune = 2
	SizeofSwampBool = 1
)

func generateExpressionWithSourceVar(code *assembler_sp.Code, expr decorated.Expression,
	genContext *generateContext, debugName string) (assembler_sp.SourceStackPosRange, error) {
	switch t := expr.(type) {
	case *decorated.StringLiteral:
		constant := genContext.context.Constants().AllocateStringConstant(t.Value())
		return constantToSourceStackPosRange(code, genContext.context.stackMemory, constant)
	case *decorated.IntegerLiteral:
		{
			intStorage := genContext.context.stackMemory.Allocate(SizeofSwampInt, SizeofSwampInt, "intLiteral")
			code.LoadInteger(intStorage.Pos, t.Value())
			return targetToSourceStackPosRange(intStorage), nil
		}
	case *decorated.CharacterLiteral:
		{
			runeStorage := genContext.context.stackMemory.Allocate(SizeofSwampRune, SizeofSwampRune, "runeLiteral")
			code.LoadRune(runeStorage.Pos, uint8(t.Value()))
			return targetToSourceStackPosRange(runeStorage), nil
		}
	case *decorated.BooleanLiteral:
		{
			boolStorage := genContext.context.stackMemory.Allocate(SizeofSwampBool, SizeofSwampBool, "boolLiteral")
			code.LoadBool(boolStorage.Pos, t.Value())
			return targetToSourceStackPosRange(boolStorage), nil
		}
	case *decorated.LetVariableReference:
		letVariableReferenceName := t.LetVariable().Name().Name()
		return genContext.context.functionVariables.FindVariable(letVariableReferenceName)
	case *decorated.FunctionParameterReference:
		parameterReferenceName := t.Identifier().Name()
		return genContext.context.functionVariables.FindVariable(parameterReferenceName)
	case *decorated.FunctionReference:
		return handleFunctionReference(code, t, genContext.context.stackMemory, genContext.context.constants)
	}

	return assembler_sp.SourceStackPosRange{}, nil
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

func generateExpression(code *assembler_sp.Code, target assembler_sp.TargetStackPosRange, expr decorated.Expression, genContext *generateContext) error {
	switch e := expr.(type) {
	case *decorated.Let:
		return generateLet(code, target, e, genContext)

	case *decorated.ArithmeticOperator:
		{
			leftPrimitive, _ := dectype.UnReference(e.Left().Type()).(*dectype.PrimitiveAtom)
			switch {
			case isListLike(e.Left().Type()) && e.OperatorType() == decorated.ArithmeticAppend:
				return generateListAppend(code, target, e, genContext)
			case leftPrimitive != nil && leftPrimitive.AtomName() == "String" && e.OperatorType() == decorated.ArithmeticAppend:
				return generateStringAppend(code, target, e, genContext)
			case isIntLike(e.Left().Type()):
				return generateArithmetic(code, target, e, genContext)
			default:
				return fmt.Errorf("cant generate arithmetic for type: %v <-> %v (%v)",
					e.Left().Type(), e.Right().Type(), e.OperatorType())
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

	case *decorated.Constant:
		return generateConstant(code, target, e, genContext)

	case *decorated.ConstantReference:
		return generateExpression(code, target, e.Constant(), genContext)

	case *decorated.FunctionParameterReference:
		return generateLocalFunctionParameterReference(code, target, e, genContext.context)

	case *decorated.CaseConsequenceParameterReference:
		return generateLocalConsequenceParameterReference(code, target, e, genContext.context)

	case *decorated.ConsOperator:
		return generateListCons(code, target, e, genContext)

	case *decorated.AsmConstant:
		return generateAsm(code, target, e, genContext.context)

	case *decorated.RecordConstructorFromRecord:
		return generateExpression(code, target, e.Expression(), genContext)

	case *decorated.RecordConstructorFromParameters:
		return generateRecordSortedAssignments(code, target, e.SortedAssignments(), genContext)
	}

	return fmt.Errorf("generate: unknown node %T %v %v", expr, expr, genContext)
}

func generateFunction(fullyQualifiedVariableName *decorated.FullyQualifiedPackageVariableName, f *decorated.FunctionValue, root *assembler_sp.FunctionRootContext, definitions *decorator.VariableContext, lookup typeinfo.TypeLookup, verboseFlag verbosity.Verbosity) (*Function, error) {
	code := assembler_sp.NewCode()
	funcContext := root.ScopeContext()
	tempVar := root.ReturnVariable()

	for _, parameter := range f.Parameters() {
		paramVarName := assembler_sp.NewVariableName(parameter.Identifier().Name())
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

	code.Return(0)

	opcodes, resolveErr := code.Resolve(root, verboseFlag >= verbosity.Mid)
	if resolveErr != nil {
		return nil, resolveErr
	}

	if verboseFlag >= verbosity.Mid {
		code.PrintOut()
	}

	parameterTypes, _ := f.ForcedFunctionType().ParameterAndReturn()
	parameterCount := uint(len(parameterTypes))

	signature, lookupErr := lookup.Lookup(f.Type())
	if lookupErr != nil {
		return nil, lookupErr
	}

	functionConstant := NewFunction(fullyQualifiedVariableName, swamppack.TypeRef(signature),
		parameterCount, root.Constants().Constants(), opcodes)

	return functionConstant, nil
}

func (g *Generator) GenerateAllLocalDefinedFunctions(module *decorated.Module, definitions *decorator.VariableContext,
	lookup typeinfo.TypeLookup, verboseFlag verbosity.Verbosity) ([]*Function, error) {
	var functionConstants []*Function

	for _, named := range module.LocalDefinitions().Definitions() {
		unknownType := named.Expression()

		fullyQualifiedName := module.FullyQualifiedName(named.Identifier())

		maybeFunction, _ := unknownType.(*decorated.FunctionValue)
		if maybeFunction != nil {
			if verboseFlag >= verbosity.Mid {
				fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
			}

			rootContext := assembler_sp.NewFunctionRootContext()

			functionConstant, genFuncErr := generateFunction(fullyQualifiedName, maybeFunction,
				rootContext, definitions, lookup, verboseFlag)
			if genFuncErr != nil {
				return nil, genFuncErr
			}

			if functionConstant == nil {
				panic(fmt.Sprintf("problem %v\n", maybeFunction))
			}

			functionConstants = append(functionConstants, functionConstant)
		} else {
			maybeConstant, _ := unknownType.(*decorated.Constant)
			if maybeConstant != nil {
				if verboseFlag >= verbosity.Mid {
					fmt.Printf("--------------------------- GenerateAllLocalDefinedFunctions function %v --------------------------\n", fullyQualifiedName)
				}
			} else {
				return nil, fmt.Errorf("generate: unknown type %T", unknownType)
			}
		}
	}

	return functionConstants, nil
}
