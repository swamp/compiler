/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package decorated

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/decorated/decshared"
	"github.com/swamp/compiler/src/decorated/dtype"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type UnMatchingBinaryOperatorTypes struct {
	typeA    dtype.Type
	typeB    dtype.Type
	operator *ast.BinaryOperator
}

func NewUnMatchingBinaryOperatorTypes(operator *ast.BinaryOperator, typeA dtype.Type, typeB dtype.Type) *UnMatchingBinaryOperatorTypes {
	return &UnMatchingBinaryOperatorTypes{operator: operator, typeA: typeA, typeB: typeB}
}

func (e *UnMatchingBinaryOperatorTypes) Error() string {
	return fmt.Sprintf("binary operator \n%v resolves to %v\n and \n %v resolves to %v", e.operator.Left(), e.typeA.HumanReadable(), e.operator.Right(), e.typeB.HumanReadable())
}

func (e *UnMatchingBinaryOperatorTypes) FetchPositionLength() token.PositionLength {
	return e.operator.Token().FetchPositionLength()
}

type UnMatchingArithmeticOperatorTypes struct {
	typeA    DecoratedExpression
	typeB    DecoratedExpression
	operator *ast.BinaryOperator
}

func NewUnMatchingArithmeticOperatorTypes(operator *ast.BinaryOperator, typeA DecoratedExpression, typeB DecoratedExpression) *UnMatchingArithmeticOperatorTypes {
	return &UnMatchingArithmeticOperatorTypes{operator: operator, typeA: typeA, typeB: typeB}
}

func (e *UnMatchingArithmeticOperatorTypes) Error() string {
	return fmt.Sprintf("arithmetic operator %v has different types %v, %v", e.operator, e.typeA, e.typeB)
}

func (e *UnMatchingArithmeticOperatorTypes) FetchPositionLength() token.PositionLength {
	return e.typeA.FetchPositionAndLength()
}

type TypeNotFound struct {
	requestedType string
}

func NewTypeNotFound(requestedType string) *TypeNotFound {
	return &TypeNotFound{requestedType: requestedType}
}

func (e *TypeNotFound) Error() string {
	return fmt.Sprintf("type %v not found", e.requestedType)
}

func (e *TypeNotFound) FetchPositionLength() token.PositionLength {
	return token.PositionLength{}
}

type UnmatchingBitwiseOperatorTypes struct {
	typeA    DecoratedExpression
	typeB    DecoratedExpression
	operator *ast.BinaryOperator
}

func NewUnmatchingBitwiseOperatorTypes(operator *ast.BinaryOperator, typeA DecoratedExpression, typeB DecoratedExpression) *UnmatchingBitwiseOperatorTypes {
	return &UnmatchingBitwiseOperatorTypes{operator: operator, typeA: typeA, typeB: typeB}
}

func (e *UnmatchingBitwiseOperatorTypes) Error() string {
	return fmt.Sprintf("bitwise operator %v has different types %v, %v", e.operator, e.typeA, e.typeB)
}

func (e *UnmatchingBitwiseOperatorTypes) FetchPositionLength() token.PositionLength {
	return e.typeA.FetchPositionAndLength()
}

type UnMatchingBooleanOperatorTypes struct {
	typeA    DecoratedExpression
	typeB    DecoratedExpression
	operator *ast.BinaryOperator
}

func NewUnMatchingBooleanOperatorTypes(operator *ast.BinaryOperator, typeA DecoratedExpression, typeB DecoratedExpression) *UnMatchingBooleanOperatorTypes {
	return &UnMatchingBooleanOperatorTypes{operator: operator, typeA: typeA, typeB: typeB}
}

func (e *UnMatchingBooleanOperatorTypes) Error() string {
	return fmt.Sprintf("boolean operator %v has different types %v, %v", e.operator, e.typeA, e.typeB)
}

func (e *UnMatchingBooleanOperatorTypes) FetchPositionLength() token.PositionLength {
	return e.typeA.FetchPositionAndLength()
}

type UnknownBinaryOperator struct {
	typeA    DecoratedExpression
	typeB    DecoratedExpression
	operator *ast.BinaryOperator
}

func NewUnknownBinaryOperator(operator *ast.BinaryOperator, typeA DecoratedExpression, typeB DecoratedExpression) *UnknownBinaryOperator {
	return &UnknownBinaryOperator{operator: operator, typeA: typeA, typeB: typeB}
}

func (e *UnknownBinaryOperator) Error() string {
	return fmt.Sprintf("unknown binary operator %v  %v, %v", e.operator, e.typeA, e.typeB)
}

func (e *UnknownBinaryOperator) FetchPositionLength() token.PositionLength {
	return e.typeA.FetchPositionAndLength()
}

type LogicalOperatorLeftMustBeBoolean struct {
	typeA    DecoratedExpression
	operator *LogicalOperator
}

func NewLogicalOperatorLeftMustBeBoolean(operator *LogicalOperator, typeA DecoratedExpression) *LogicalOperatorLeftMustBeBoolean {
	return &LogicalOperatorLeftMustBeBoolean{operator: operator, typeA: typeA}
}

func (e *LogicalOperatorLeftMustBeBoolean) Error() string {
	return fmt.Sprintf("logical operator %v  must have left boolean %v", e.operator, e.typeA)
}

func (e *LogicalOperatorLeftMustBeBoolean) FetchPositionLength() token.PositionLength {
	return e.typeA.FetchPositionAndLength()
}

type LogicalOperatorsMustBeBoolean struct {
	typeA    DecoratedExpression
	typeB    DecoratedExpression
	operator *ast.BinaryOperator
}

func NewLogicalOperatorsMustBeBoolean(operator *ast.BinaryOperator, typeA DecoratedExpression, typeB DecoratedExpression) *LogicalOperatorsMustBeBoolean {
	return &LogicalOperatorsMustBeBoolean{operator: operator, typeA: typeA, typeB: typeB}
}

func (e *LogicalOperatorsMustBeBoolean) Error() string {
	return fmt.Sprintf("logical operator %v  must have booleans %v vs %v", e.operator, e.typeA, e.typeB)
}

func (e *LogicalOperatorsMustBeBoolean) FetchPositionLength() token.PositionLength {
	return e.typeA.FetchPositionAndLength()
}

type LogicalOperatorRightMustBeBoolean struct {
	typeA    DecoratedExpression
	operator *LogicalOperator
}

func NewLogicalOperatorRightMustBeBoolean(operator *LogicalOperator, typeA DecoratedExpression) *LogicalOperatorRightMustBeBoolean {
	return &LogicalOperatorRightMustBeBoolean{operator: operator, typeA: typeA}
}

func (e *LogicalOperatorRightMustBeBoolean) Error() string {
	return fmt.Sprintf("logical operator %v  must have right boolean %v", e.operator, e.typeA)
}

func (e *LogicalOperatorRightMustBeBoolean) FetchPositionLength() token.PositionLength {
	return e.typeA.FetchPositionAndLength()
}

type MustBeCustomType struct {
	typeA DecoratedExpression
}

func NewMustBeCustomType(typeA DecoratedExpression) *MustBeCustomType {
	return &MustBeCustomType{typeA: typeA}
}

func (e *MustBeCustomType) Error() string {
	return fmt.Sprintf("must be a custom type %v", e.typeA)
}

func (e *MustBeCustomType) FetchPositionLength() token.PositionLength {
	return e.typeA.FetchPositionAndLength()
}

type CaseCouldNotFindCustomVariantType struct {
	caseExpression *ast.Case
	consequence    *ast.CaseConsequence
}

func NewCaseCouldNotFindCustomVariantType(caseExpression *ast.Case, consequence *ast.CaseConsequence) *CaseCouldNotFindCustomVariantType {
	return &CaseCouldNotFindCustomVariantType{consequence: consequence, caseExpression: caseExpression}
}

func (e *CaseCouldNotFindCustomVariantType) Error() string {
	return fmt.Sprintf("couldn't find custom variant in case consequence %v (%v)", e.consequence, e.caseExpression)
}

func (e *CaseCouldNotFindCustomVariantType) FetchPositionLength() token.PositionLength {
	return e.consequence.Identifier().Symbol().FetchPositionLength()
}

type UnMatchingTypesError struct {
	ExpectedType dtype.Type
	HasType      dtype.Type
}

func (e UnMatchingTypesError) Error() string {
	return fmt.Sprintf("%v\n %v\n", e.HasType, e.ExpectedType)
}

type UnMatchingTypesExpression struct {
	UnMatchingTypesError
	expression ast.Expression
	err        error
}

func NewUnMatchingTypes(expression ast.Expression, expectedType dtype.Type, hasType dtype.Type, err error) *UnMatchingTypesExpression {
	if expectedType == nil || hasType == nil {
		panic("not allowed")
	}
	return &UnMatchingTypesExpression{
		expression:           expression,
		UnMatchingTypesError: UnMatchingTypesError{ExpectedType: expectedType, HasType: hasType},
		err:                  err,
	}
}

func (e *UnMatchingTypesExpression) Error() string {
	expectedAtom, _ := e.ExpectedType.Resolve()
	hasAtom, _ := e.HasType.Resolve()
	return fmt.Sprintf("unmatching types %v\nvs\n%v\n %v\n%v", expectedAtom.AtomName(), hasAtom.AtomName(), e.UnMatchingTypesError.Error(), e.err)
}

func (e *UnMatchingTypesExpression) FetchPositionLength() token.PositionLength {
	return e.expression.PositionLength()
}

type UnMatchingFunctionReturnTypesInFunctionValue struct {
	UnMatchingTypesExpression
	expression ast.Expression
	fn         *ast.FunctionValue
	err        error
}

func NewUnMatchingFunctionReturnTypesInFunctionValue(fn *ast.FunctionValue, expression ast.Expression, expectedType dtype.Type, hasType dtype.Type, err error) *UnMatchingFunctionReturnTypesInFunctionValue {
	return &UnMatchingFunctionReturnTypesInFunctionValue{
		fn: fn, expression: expression,
		UnMatchingTypesExpression: UnMatchingTypesExpression{
			expression:           expression,
			UnMatchingTypesError: UnMatchingTypesError{ExpectedType: expectedType, HasType: hasType},
		}, err: err,
	}
}

func (e *UnMatchingFunctionReturnTypesInFunctionValue) Error() string {
	return fmt.Sprintf("unmatching function return types")
}

func (e *UnMatchingFunctionReturnTypesInFunctionValue) FetchPositionLength() token.PositionLength {
	return e.fn.Expression().PositionLength()
}

type FunctionArgumentTypeMismatch struct {
	UnMatchingTypesError
	expression        ast.Expression
	err               error
	decoratedArgument DecoratedExpression
	symbolToken       token.PositionLength
}

func NewFunctionArgumentTypeMismatch(symbolToken token.PositionLength, decoratedArgument DecoratedExpression, expression ast.Expression, expectedType dtype.Type, hasType dtype.Type, err error) *FunctionArgumentTypeMismatch {
	return &FunctionArgumentTypeMismatch{
		symbolToken: symbolToken, decoratedArgument: decoratedArgument,
		expression:           expression,
		UnMatchingTypesError: UnMatchingTypesError{ExpectedType: expectedType, HasType: hasType}, err: err,
	}
}

func (e *FunctionArgumentTypeMismatch) Error() string {
	return fmt.Sprintf("function argument type mismatch\n\n%v\n\n in expression %v\n%v\n%v\n", e.UnMatchingTypesError.Error(), e.expression, e.decoratedArgument, e.err)
}

func (e *FunctionArgumentTypeMismatch) FetchPositionLength() token.PositionLength {
	return e.symbolToken.FetchPositionLength()
}

type RecordLiteralFieldTypeMismatch struct {
	UnMatchingTypesError
	field      *dectype.RecordField
	assignment *ast.RecordLiteralFieldAssignment
	err        error
}

func NewRecordLiteralFieldTypeMismatch(assignment *ast.RecordLiteralFieldAssignment, field *dectype.RecordField, encounteredType dtype.Type, err error) *RecordLiteralFieldTypeMismatch {
	return &RecordLiteralFieldTypeMismatch{
		field: field, assignment: assignment,
		UnMatchingTypesError: UnMatchingTypesError{ExpectedType: field.Type(), HasType: encounteredType}, err: err,
	}
}

func (e *RecordLiteralFieldTypeMismatch) Error() string {
	return fmt.Sprintf("record literal type field mismatch %v  %v %v\n", e.field, e.UnMatchingTypesError.Error(), e.err)
}

func (e *RecordLiteralFieldTypeMismatch) FetchPositionLength() token.PositionLength {
	return e.assignment.Expression().PositionLength()
}

type NewRecordLiteralFieldNotInType struct {
	recordType *dectype.RecordAtom
	assignment *ast.RecordLiteralFieldAssignment
}

func NewNewRecordLiteralFieldNotInType(assignment *ast.RecordLiteralFieldAssignment, recordType *dectype.RecordAtom) *NewRecordLiteralFieldNotInType {
	return &NewRecordLiteralFieldNotInType{recordType: recordType, assignment: assignment}
}

func (e *NewRecordLiteralFieldNotInType) Error() string {
	return fmt.Sprintf("record literal field not in type %v %v\n", e.assignment, e.recordType)
}

func (e *NewRecordLiteralFieldNotInType) FetchPositionLength() token.PositionLength {
	return e.assignment.Identifier().PositionLength()
}

type ConstructorArgumentTypeMismatch struct {
	expectedType dtype.Type
	hasType      dtype.Type
	call         *ast.ConstructorCall
	err          error
}

func NewConstructorArgumentTypeMismatch(call *ast.ConstructorCall, expectedType dtype.Type, hasType dtype.Type, err error) *ConstructorArgumentTypeMismatch {
	return &ConstructorArgumentTypeMismatch{call: call, expectedType: expectedType, hasType: hasType, err: err}
}

func (e *ConstructorArgumentTypeMismatch) Error() string {
	return fmt.Sprintf("expected constructor argument %v but has %v in expression %v\n\n%v but received %v\n%v\n", e.expectedType, e.hasType, e.call, e.expectedType, e.hasType, e.err)
}

func (e *ConstructorArgumentTypeMismatch) FetchPositionLength() token.PositionLength {
	return e.call.TypeIdentifier().Symbol().FetchPositionLength()
}

type ExpectedCustomTypeVariantConstructor struct {
	call *ast.ConstructorCall
	err  error
}

func NewExpectedCustomTypeVariantConstructor(call *ast.ConstructorCall) *ExpectedCustomTypeVariantConstructor {
	return &ExpectedCustomTypeVariantConstructor{call: call}
}

func (e *ExpectedCustomTypeVariantConstructor) Error() string {
	return fmt.Sprintf("expected constructor %v", e.call.TypeIdentifier())
}

func (e *ExpectedCustomTypeVariantConstructor) FetchPositionLength() token.PositionLength {
	return e.call.TypeIdentifier().Symbol().FetchPositionLength()
}

type WrongTypeForRecordConstructorField struct {
	UnMatchingTypesError
	call               *ast.ConstructorCall
	recordField        *dectype.RecordField
	expectedExpression DecoratedExpression
	err                error
}

func NewWrongTypeForRecordConstructorField(recordField *dectype.RecordField, expectedExpression DecoratedExpression,
	call *ast.ConstructorCall, err error) *WrongTypeForRecordConstructorField {
	return &WrongTypeForRecordConstructorField{
		UnMatchingTypesError: UnMatchingTypesError{ExpectedType: expectedExpression.Type(), HasType: recordField.Type()},
		call:                 call, recordField: recordField, expectedExpression: expectedExpression, err: err,
	}
}

func (e *WrongTypeForRecordConstructorField) Error() string {
	return fmt.Sprintf("wrong type for constructor field %v\nexpected\n%v\nvs\n%v", e.recordField.Name(), e.UnMatchingTypesError.ExpectedType.HumanReadable(), e.UnMatchingTypesError.HasType.HumanReadable())
}

func (e *WrongTypeForRecordConstructorField) FetchPositionLength() token.PositionLength {
	return e.call.TypeIdentifier().Symbol().FetchPositionLength()
}

type UnhandledCustomTypeVariants struct {
	unhandledVariants []*dectype.CustomTypeVariant
	caseExpression    *ast.Case
}

func NewUnhandledCustomTypeVariants(caseExpression *ast.Case, unhandledVariants []*dectype.CustomTypeVariant) *UnhandledCustomTypeVariants {
	return &UnhandledCustomTypeVariants{unhandledVariants: unhandledVariants, caseExpression: caseExpression}
}

func (e *UnhandledCustomTypeVariants) Error() string {
	return fmt.Sprintf("unhandled custom type variants %v", e.unhandledVariants)
}

func (e *UnhandledCustomTypeVariants) FetchPositionLength() token.PositionLength {
	return e.caseExpression.Keyword().FetchPositionLength()
}

type AlreadyHandledCustomTypeVariant struct {
	unhandledVariant *dectype.CustomTypeVariant
	caseExpression   *ast.Case
	consequence      *ast.CaseConsequence
}

func NewAlreadyHandledCustomTypeVariant(caseExpression *ast.Case, consequence *ast.CaseConsequence, unhandledVariant *dectype.CustomTypeVariant) *AlreadyHandledCustomTypeVariant {
	return &AlreadyHandledCustomTypeVariant{unhandledVariant: unhandledVariant, caseExpression: caseExpression, consequence: consequence}
}

func (e *AlreadyHandledCustomTypeVariant) Error() string {
	return fmt.Sprintf("already handled variant %v", e.unhandledVariant)
}

func (e *AlreadyHandledCustomTypeVariant) FetchPositionLength() token.PositionLength {
	return e.consequence.Identifier().Symbol().FetchPositionLength()
}

type ExpectedFunctionType struct {
	expected    dtype.Type
	encountered ast.Expression
}

func NewExpectedFunctionType(expected dtype.Type, encountered ast.Expression) *ExpectedFunctionType {
	return &ExpectedFunctionType{expected: expected, encountered: encountered}
}

func (e *ExpectedFunctionType) Error() string {
	return fmt.Sprintf("supposed to be function %v %T", e.expected, e.expected)
}

func (e *ExpectedFunctionType) FetchPositionLength() token.PositionLength {
	return e.encountered.PositionLength()
}

type ExpectedFunctionTypeForCall struct {
	encountered DecoratedExpression
}

func NewExpectedFunctionTypeForCall(encountered DecoratedExpression) *ExpectedFunctionTypeForCall {
	return &ExpectedFunctionTypeForCall{encountered: encountered}
}

func (e *ExpectedFunctionTypeForCall) Error() string {
	return fmt.Sprintf("this must be a function type that you are calling, but was something completely else %v", e.encountered)
}

func (e *ExpectedFunctionTypeForCall) FetchPositionLength() token.PositionLength {
	return e.encountered.FetchPositionAndLength()
}

type FunctionCallTypeMismatch struct {
	expected    *dectype.FunctionAtom
	encountered *dectype.FunctionAtom
	call        *ast.FunctionCall
}

func NewFunctionCallTypeMismatch(call *ast.FunctionCall, expected *dectype.FunctionAtom, encountered *dectype.FunctionAtom) *FunctionCallTypeMismatch {
	return &FunctionCallTypeMismatch{call: call, expected: expected, encountered: encountered}
}

func (e *FunctionCallTypeMismatch) Error() string {
	return fmt.Sprintf("mismatch function type %v vs %v\n%v\n", e.expected, e.encountered, e.call)
}

func (e *FunctionCallTypeMismatch) FetchPositionLength() token.PositionLength {
	return e.call.PositionLength()
}

type CouldNotSmashFunctions struct {
	expected    *dectype.FunctionAtom
	encountered *dectype.FunctionAtom
	call        *ast.FunctionCall
	err         error
}

func NewCouldNotSmashFunctions(call *ast.FunctionCall, expected *dectype.FunctionAtom, encountered *dectype.FunctionAtom, err error) *CouldNotSmashFunctions {
	return &CouldNotSmashFunctions{call: call, expected: expected, encountered: encountered, err: err}
}

func (e *CouldNotSmashFunctions) Error() string {
	return fmt.Sprintf("could not smash function types \n%v\n and\n%v\n%v\n%v", e.expected.HumanReadable(), e.encountered.HumanReadable(), e.call, e.err)
}

func (e *CouldNotSmashFunctions) FetchPositionLength() token.PositionLength {
	return e.call.PositionLength()
}

type CaseWrongParameterCountInCustomTypeVariant struct {
	unhandledVariant *dectype.CustomTypeVariant
	caseExpression   *ast.Case
	consequence      *ast.CaseConsequence
}

type ExtraFunctionArguments struct {
	expected    []dtype.Type
	encountered []dtype.Type
	posLength   token.PositionLength
}

func NewExtraFunctionArguments(posLength token.PositionLength, expected []dtype.Type, encountered []dtype.Type) *ExtraFunctionArguments {
	return &ExtraFunctionArguments{posLength: posLength, expected: expected, encountered: encountered}
}

func (e *ExtraFunctionArguments) Error() string {
	return fmt.Sprintf("you can not define more arguments than what is expected of you. function:\n%v\nEncountered:\n%v\n", e.expected, e.encountered)
}

func (e *ExtraFunctionArguments) FetchPositionLength() token.PositionLength {
	return e.posLength
}

func NewCaseWrongParameterCountInCustomTypeVariant(caseExpression *ast.Case, consequence *ast.CaseConsequence, unhandledVariant *dectype.CustomTypeVariant) *CaseWrongParameterCountInCustomTypeVariant {
	return &CaseWrongParameterCountInCustomTypeVariant{unhandledVariant: unhandledVariant, caseExpression: caseExpression, consequence: consequence}
}

func (e *CaseWrongParameterCountInCustomTypeVariant) Error() string {
	return fmt.Sprintf("wrong parameter count in custom type variant %v (%v)", e.consequence, e.caseExpression)
}

func (e *CaseWrongParameterCountInCustomTypeVariant) FetchPositionLength() token.PositionLength {
	return e.consequence.Identifier().Symbol().FetchPositionLength()
}

type YouCanOnlySetFieldInRecordOnce struct {
	recordType *dectype.RecordAtom
	fieldName  *ast.VariableIdentifier
}

func NewYouCanOnlySetFieldInRecordOnce(recordType *dectype.RecordAtom, fieldName *ast.VariableIdentifier) *YouCanOnlySetFieldInRecordOnce {
	return &YouCanOnlySetFieldInRecordOnce{recordType: recordType, fieldName: fieldName}
}

func (e *YouCanOnlySetFieldInRecordOnce) Error() string {
	return fmt.Sprintf("you can only set field in record once %v %v", e.fieldName, e.recordType)
}

func (e *YouCanOnlySetFieldInRecordOnce) FetchPositionLength() token.PositionLength {
	return e.fieldName.Symbol().FetchPositionLength()
}

type WrongNumberOfArgumentsInFunctionValue struct {
	argumentTypes       []dtype.Type
	encounteredFunction *ast.FunctionValue
	expectedFunction    *dectype.FunctionAtom
}

func NewWrongNumberOfArgumentsInFunctionValue(encounteredFunction *ast.FunctionValue,
	expectedFunction *dectype.FunctionAtom, argumentTypes []dtype.Type) *WrongNumberOfArgumentsInFunctionValue {
	return &WrongNumberOfArgumentsInFunctionValue{
		encounteredFunction: encounteredFunction,
		argumentTypes:       argumentTypes,
		expectedFunction:    expectedFunction,
	}
}

func (e *WrongNumberOfArgumentsInFunctionValue) ExpectedFunctionType() *dectype.FunctionAtom {
	return e.expectedFunction
}

func (e *WrongNumberOfArgumentsInFunctionValue) EncounteredFunctionValue() *ast.FunctionValue {
	return e.encounteredFunction
}

func (e *WrongNumberOfArgumentsInFunctionValue) EncounteredArgumentTypes() []dtype.Type {
	return e.argumentTypes
}

func (e *WrongNumberOfArgumentsInFunctionValue) Error() string {
	return fmt.Sprintf("wrong number of arguments in functionvalue")
}

func (e *WrongNumberOfArgumentsInFunctionValue) FetchPositionLength() token.PositionLength {
	return e.encounteredFunction.PositionLength()
}

type IfTestMustHaveBooleanType struct {
	ifTestExpression DecoratedExpression
	ifExpression     *ast.IfExpression
}

func NewIfTestMustHaveBooleanType(ifExpression *ast.IfExpression, ifTestExpression DecoratedExpression) *IfTestMustHaveBooleanType {
	return &IfTestMustHaveBooleanType{ifExpression: ifExpression, ifTestExpression: ifTestExpression}
}

func (e *IfTestMustHaveBooleanType) Error() string {
	return fmt.Sprintf("if test must have Bool type %v", e.ifExpression)
}

func (e *IfTestMustHaveBooleanType) FetchPositionLength() token.PositionLength {
	return e.ifTestExpression.FetchPositionAndLength()
}

type IfConsequenceAndAlternativeMustHaveSameType struct {
	ifExpression  *ast.IfExpression
	consequence   DecoratedExpression
	alternative   DecoratedExpression
	compatibleErr error
}

func NewIfConsequenceAndAlternativeMustHaveSameType(ifExpression *ast.IfExpression, consequence DecoratedExpression, alternative DecoratedExpression, compatibleErr error) *IfConsequenceAndAlternativeMustHaveSameType {
	return &IfConsequenceAndAlternativeMustHaveSameType{ifExpression: ifExpression, consequence: consequence, alternative: alternative, compatibleErr: compatibleErr}
}

func (e *IfConsequenceAndAlternativeMustHaveSameType) Error() string {
	return fmt.Sprintf("if: consequence and alternative must have same type %v\nvs\n%v\n", e.consequence.Type().HumanReadable(), e.alternative.Type().HumanReadable())
}

func (e *IfConsequenceAndAlternativeMustHaveSameType) FetchPositionLength() token.PositionLength {
	return e.consequence.FetchPositionAndLength()
}

type EveryItemInThelistMustHaveTheSameType struct {
	list             *ast.ListLiteral
	faultyExpression ast.Expression
	ExpectedType     dtype.Type
	ActualType       dtype.Type
	compatibleErr    error
}

func NewEveryItemInThelistMustHaveTheSameType(list *ast.ListLiteral, faultyExpression ast.Expression, expectedType dtype.Type, actualType dtype.Type, compatibleErr error) *EveryItemInThelistMustHaveTheSameType {
	return &EveryItemInThelistMustHaveTheSameType{list: list, faultyExpression: faultyExpression, ExpectedType: expectedType, ActualType: actualType, compatibleErr: compatibleErr}
}

func (e *EveryItemInThelistMustHaveTheSameType) Error() string {
	return fmt.Sprintf("every item in the list must have same type %v", e.list)
}

func (e *EveryItemInThelistMustHaveTheSameType) FetchPositionLength() token.PositionLength {
	return e.faultyExpression.PositionLength()
}

type CouldNotFindDefinitionOrTypeForIdentifier struct {
	ident *ast.VariableIdentifier
	// context *dectype.VariableContext
}

func NewCouldNotFindDefinitionOrTypeForIdentifier(ident *ast.VariableIdentifier) *CouldNotFindDefinitionOrTypeForIdentifier {
	return &CouldNotFindDefinitionOrTypeForIdentifier{ident: ident}
}

func (e *CouldNotFindDefinitionOrTypeForIdentifier) Error() string {
	return fmt.Sprintf("could not find %v", e.ident)
}

func (e *CouldNotFindDefinitionOrTypeForIdentifier) FetchPositionLength() token.PositionLength {
	return e.ident.Symbol().FetchPositionLength()
}

type CouldNotFindTypeForTypeIdentifier struct {
	ident *ast.TypeIdentifier
	// context *dectype.VariableContext
}

func NewCouldNotFindTypeForTypeIdentifier(ident *ast.TypeIdentifier) *CouldNotFindTypeForTypeIdentifier {
	return &CouldNotFindTypeForTypeIdentifier{ident: ident}
}

func (e *CouldNotFindTypeForTypeIdentifier) Error() string {
	return fmt.Sprintf("could not find %v", e.ident)
}

func (e *CouldNotFindTypeForTypeIdentifier) FetchPositionLength() token.PositionLength {
	return e.ident.Symbol().FetchPositionLength()
}

type CouldNotFindIdentifierInLookups struct {
	lookups *ast.Lookups
}

func NewCouldNotFindIdentifierInLookups(lookups *ast.Lookups) *CouldNotFindIdentifierInLookups {
	return &CouldNotFindIdentifierInLookups{lookups: lookups}
}

func (e *CouldNotFindIdentifierInLookups) Error() string {
	return fmt.Sprintf("could not find %v", e.lookups)
}

func (e *CouldNotFindIdentifierInLookups) FetchPositionLength() token.PositionLength {
	return e.lookups.ContextIdentifier().Symbol().FetchPositionLength()
}

type CouldNotFindFieldInLookup struct {
	lookups        *ast.Lookups
	variable       *ast.VariableIdentifier
	typeToLookupIn dtype.Type
}

func NewCouldNotFindFieldInLookup(lookups *ast.Lookups, variable *ast.VariableIdentifier, typeToLookupIn dtype.Type) *CouldNotFindFieldInLookup {
	return &CouldNotFindFieldInLookup{lookups: lookups, variable: variable, typeToLookupIn: typeToLookupIn}
}

func (e *CouldNotFindFieldInLookup) Error() string {
	return fmt.Sprintf("could not find %v", e.lookups)
}

func (e *CouldNotFindFieldInLookup) FetchPositionLength() token.PositionLength {
	return e.variable.Symbol().FetchPositionLength()
}

type MustHaveAnnotationJustBeforeThisDefinition struct {
	assignment *ast.DefinitionAssignment
}

func NewMustHaveAnnotationJustBeforeThisDefinition(assignment *ast.DefinitionAssignment) *MustHaveAnnotationJustBeforeThisDefinition {
	return &MustHaveAnnotationJustBeforeThisDefinition{assignment: assignment}
}

func (e *MustHaveAnnotationJustBeforeThisDefinition) Error() string {
	return fmt.Sprintf("must have annotation before this definition %v", e.assignment)
}

func (e *MustHaveAnnotationJustBeforeThisDefinition) FetchPositionLength() token.PositionLength {
	return e.assignment.Identifier().Symbol().FetchPositionLength()
}

type AlreadyHaveAnnotationForThisName struct {
	annotation         *ast.Annotation
	previousAnnotation *ast.Annotation
}

func NewAlreadyHaveAnnotationForThisName(annotation *ast.Annotation, previousAnnotation *ast.Annotation) *AlreadyHaveAnnotationForThisName {
	return &AlreadyHaveAnnotationForThisName{annotation: annotation, previousAnnotation: previousAnnotation}
}

func (e *AlreadyHaveAnnotationForThisName) Error() string {
	return fmt.Sprintf("already have annotation for this name %v", e.annotation)
}

func (e *AlreadyHaveAnnotationForThisName) FetchPositionLength() token.PositionLength {
	return e.annotation.Identifier().Symbol().FetchPositionLength()
}

type UnknownStatement struct {
	statement ast.Expression
	posLength token.PositionLength
}

func NewUnknownStatement(posLength token.PositionLength, statement ast.Expression) *UnknownStatement {
	return &UnknownStatement{statement: statement}
}

func (e *UnknownStatement) Error() string {
	return fmt.Sprintf("unknown statement %v %T", e.statement, e.statement)
}

func (e *UnknownStatement) FetchPositionLength() token.PositionLength {
	return e.posLength
}

type UnknownVariable struct {
	ident *ast.VariableIdentifier
}

func NewUnknownVariable(ident *ast.VariableIdentifier) *UnknownVariable {
	return &UnknownVariable{ident: ident}
}

func (e *UnknownVariable) Error() string {
	return fmt.Sprintf("unknown variable '%v'", e.ident)
}

func (e *UnknownVariable) FetchPositionLength() token.PositionLength {
	return e.ident.PositionLength()
}

type AnnotationMismatch struct {
	assignment           *ast.DefinitionAssignment
	annotationIdentifier *ast.VariableIdentifier
}

func NewAnnotationMismatch(annotationIdentifier *ast.VariableIdentifier, assignment *ast.DefinitionAssignment) *AnnotationMismatch {
	return &AnnotationMismatch{assignment: assignment, annotationIdentifier: annotationIdentifier}
}

func (e *AnnotationMismatch) Error() string {
	return fmt.Sprintf("annotation mismatch %v", e.assignment)
}

func (e *AnnotationMismatch) FetchPositionLength() token.PositionLength {
	return e.assignment.Identifier().Symbol().FetchPositionLength()
}

type ModuleError struct {
	err        decshared.DecoratedError
	sourceFile string
}

func NewModuleError(sourceFile string, err decshared.DecoratedError) *ModuleError {
	return &ModuleError{sourceFile: sourceFile, err: err}
}

func (e *ModuleError) WrappedError() decshared.DecoratedError {
	return e.err
}

func (e *ModuleError) Error() string {
	return fmt.Sprintf("module error '%v': %v", e.sourceFile, e.err)
}

func (e *ModuleError) FetchPositionLength() token.PositionLength {
	return e.err.FetchPositionLength()
}

type InternalError struct {
	err error
}

func NewInternalError(err error) *InternalError {
	return &InternalError{err: err}
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("decoration internal error %v", e.err)
}

func (e *InternalError) FetchPositionLength() token.PositionLength {
	return token.PositionLength{}
}

type UnknownAnnotationTypeReference struct {
	annotationIdentifier *ast.VariableIdentifier
	err                  dectype.DecoratedTypeError
}

func NewUnknownAnnotationTypeReference(annotationIdentifier *ast.VariableIdentifier, err dectype.DecoratedTypeError) *UnknownAnnotationTypeReference {
	return &UnknownAnnotationTypeReference{err: err, annotationIdentifier: annotationIdentifier}
}

func (e *UnknownAnnotationTypeReference) Error() string {
	return fmt.Sprintf("unknown type in annotation %v %v", e.annotationIdentifier, e.err)
}

func (e *UnknownAnnotationTypeReference) FetchPositionLength() token.PositionLength {
	return e.annotationIdentifier.Symbol().FetchPositionLength()
}

type UnknownTypeAliasType struct {
	alias *ast.AliasStatement
	err   dectype.DecoratedTypeError
}

func NewUnknownTypeAliasType(alias *ast.AliasStatement, err dectype.DecoratedTypeError) *UnknownTypeAliasType {
	return &UnknownTypeAliasType{err: err, alias: alias}
}

func (e *UnknownTypeAliasType) Error() string {
	return fmt.Sprintf("unknown type in alias %v %v", e.alias, e.err)
}

func (e *UnknownTypeAliasType) FetchPositionLength() token.PositionLength {
	return e.alias.TypeIdentifier().Symbol().FetchPositionLength()
}

type UnknownTypeInCustomTypeVariant struct {
	variant *ast.CustomTypeVariant
	err     dectype.DecoratedTypeError
}

func NewUnknownTypeInCustomTypeVariant(variant *ast.CustomTypeVariant, err dectype.DecoratedTypeError) *UnknownTypeInCustomTypeVariant {
	return &UnknownTypeInCustomTypeVariant{err: err, variant: variant}
}

func (e *UnknownTypeInCustomTypeVariant) Error() string {
	return fmt.Sprintf("unknown type in custom type variant %v %v", e.variant, e.err)
}

func (e *UnknownTypeInCustomTypeVariant) FetchPositionLength() token.PositionLength {
	return e.variant.TypeIdentifier().Symbol().FetchPositionLength()
}
