/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parerr

import (
	"fmt"

	"github.com/swamp/compiler/src/ast"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

type ParseError interface {
	FetchPositionLength() token.SourceFileReference
	Error() string
}

type SubError struct {
	SubErr tokenize.TokenError
}

func NewSubError(err tokenize.TokenError) SubError {
	if err == nil {
		panic("must have sub error")
	}
	return SubError{SubErr: err}
}

func (e SubError) FetchPositionLength() token.SourceFileReference {
	return e.SubErr.FetchPositionLength()
}

type SubParserError struct {
	SubErr ParseError
}

func NewSubParserError(err ParseError) SubParserError {
	return SubParserError{SubErr: err}
}

func (e SubParserError) FetchPositionLength() token.SourceFileReference {
	return e.SubErr.FetchPositionLength()
}

func (e SubParserError) Unwrap() ParseError {
	return e.SubErr
}

type ExpectedTypeOrParenError struct {
	token.Token
}

type ExpectedNewLineCount struct {
	reference   token.SourceFileReference
	expected    int
	encountered int
}

func NewExpectedNewLineCount(reference token.SourceFileReference, expected int, encountered int) ExpectedNewLineCount {
	return ExpectedNewLineCount{reference: reference, expected: expected, encountered: encountered}
}

func (e ExpectedNewLineCount) Error() string {
	return fmt.Sprintf("wrong exact number of line %v expected %v", e.encountered, e.expected)
}

func (e ExpectedNewLineCount) FetchPositionLength() token.SourceFileReference {
	return e.reference
}

type UnexpectedImportAlias struct {
	importStatement *ast.Import
}

func NewUnexpectedImportAlias(importStatement *ast.Import) UnexpectedImportAlias {
	return UnexpectedImportAlias{importStatement: importStatement}
}

func (e UnexpectedImportAlias) Error() string {
	return fmt.Sprintf("it is advised to use the last part of the import as alias. `import Some.Long.Name as Name` (%v)", e.importStatement.Alias().Name())
}

func (e UnexpectedImportAlias) FetchPositionLength() token.SourceFileReference {
	return e.importStatement.FetchPositionLength()
}

func NewExpectedTypeOrParenError(t token.Token) ExpectedTypeOrParenError {
	return ExpectedTypeOrParenError{Token: t}
}

func (e ExpectedTypeOrParenError) Error() string {
	return fmt.Sprintf("typeterm: expected type symbol or (. %v", e.Token)
}

type InternalError struct {
	token.SourceFileReference
	internalError error
}

func NewInternalError(t token.SourceFileReference, internalError error) InternalError {
	return InternalError{SourceFileReference: t, internalError: internalError}
}

func (e InternalError) Error() string {
	return fmt.Sprintf("parser internal error %v", e.internalError)
}

func (e InternalError) FetchPositionLength() token.SourceFileReference {
	return token.SourceFileReference{}
}

type MustBeSpaceOrContinuation struct {
	token.SourceFileReference
}

func NewMustBeSpaceOrContinuation(t token.SourceFileReference) MustBeSpaceOrContinuation {
	return MustBeSpaceOrContinuation{SourceFileReference: t}
}

func (e MustBeSpaceOrContinuation) Error() string {
	return fmt.Sprintf("must be space or continuation %v", e.SourceFileReference)
}

func (e MustBeSpaceOrContinuation) FetchPositionLength() token.SourceFileReference {
	return token.SourceFileReference{}
}

type NotATermError struct {
	token.Token
}

func NewNotATermError(t token.Token) NotATermError {
	if _, isEOF := t.(*tokenize.EndOfFile); isEOF {
		t = nil
	}
	return NotATermError{Token: t}
}

func (e NotATermError) Error() string {
	return fmt.Sprintf("not a term error %v", e.Token)
}

type ExpectedCaseConsequenceSymbolError struct {
	SubError
}

func NewExpectedCaseConsequenceSymbolError(err tokenize.TokenError) ExpectedCaseConsequenceSymbolError {
	return ExpectedCaseConsequenceSymbolError{SubError: SubError{err}}
}

func (e ExpectedCaseConsequenceSymbolError) Error() string {
	return fmt.Sprintf("wasn't a case symbol %v", e.SubError)
}

type ExpectedTwoLinesAfterStatement struct {
	subError ParseError
}

func NewExpectedTwoLinesAfterStatement(err ParseError) ExpectedTwoLinesAfterStatement {
	return ExpectedTwoLinesAfterStatement{subError: err}
}

func (e ExpectedTwoLinesAfterStatement) Error() string {
	return fmt.Sprintf("needs exactly two empty lines after each statement %v", e.subError)
}

func (e ExpectedTwoLinesAfterStatement) FetchPositionLength() token.SourceFileReference {
	return e.subError.FetchPositionLength()
}

type ExpectedIndentationError struct {
	posLength           token.SourceFileReference
	expectedIndentation int
}

func NewExpectedIndentationError(posLength token.SourceFileReference, expectedIndentation int) ExpectedIndentationError {
	return ExpectedIndentationError{posLength: posLength, expectedIndentation: expectedIndentation}
}

func (e ExpectedIndentationError) Error() string {
	return fmt.Sprintf("surprised by the indentation %v", e.posLength)
}

func (e ExpectedIndentationError) FetchPositionLength() token.SourceFileReference {
	return e.posLength
}

type ExpectedTypeReferenceError struct {
	posLength token.SourceFileReference
}

func NewExpectedTypeReferenceError(posLength token.SourceFileReference) ExpectedTypeReferenceError {
	return ExpectedTypeReferenceError{posLength: posLength}
}

func (e ExpectedTypeReferenceError) Error() string {
	return fmt.Sprintf("expected a type reference here %v", e.posLength)
}

func (e ExpectedTypeReferenceError) FetchPositionLength() token.SourceFileReference {
	return e.posLength
}

type ExpectedRightArrowError struct {
	SubError
}

func NewExpectedRightArrowError(err tokenize.TokenError) ExpectedRightArrowError {
	return ExpectedRightArrowError{SubError: SubError{err}}
}

func (e ExpectedRightArrowError) Error() string {
	return fmt.Sprintf("expected a -> %v", e.SubError)
}

type CaseConsequenceExpectedVariableOrRightArrow struct {
	SubParserError
}

func NewCaseConsequenceExpectedVariableOrRightArrow(err ParseError) CaseConsequenceExpectedVariableOrRightArrow {
	return CaseConsequenceExpectedVariableOrRightArrow{SubParserError: NewSubParserError(err)}
}

func (e CaseConsequenceExpectedVariableOrRightArrow) Error() string {
	return fmt.Sprintf("case consequences expected variable or '->' %v", e.SubParserError)
}

type TypeMustBeFollowedByTypeArgumentOrEqualError struct {
	subError ParseError
}

func NewTypeMustBeFollowedByTypeArgumentOrEqualError(subError ParseError) TypeMustBeFollowedByTypeArgumentOrEqualError {
	return TypeMustBeFollowedByTypeArgumentOrEqualError{subError: subError}
}

func (e TypeMustBeFollowedByTypeArgumentOrEqualError) Error() string {
	return fmt.Sprintf("type can only be followed by type parameters or = (%v)", e.subError)
}

func (e TypeMustBeFollowedByTypeArgumentOrEqualError) FetchPositionLength() token.SourceFileReference {
	return e.subError.FetchPositionLength()
}

type MustHaveAtLeastOneParameterError struct {
	posLength token.SourceFileReference
}

func NewMustHaveAtLeastOneParameterError(posLength token.SourceFileReference) MustHaveAtLeastOneParameterError {
	return MustHaveAtLeastOneParameterError{posLength: posLength}
}

func (e MustHaveAtLeastOneParameterError) Error() string {
	return fmt.Sprintf("must have at least one parameter %v", e.posLength)
}

func (e MustHaveAtLeastOneParameterError) FetchPositionLength() token.SourceFileReference {
	return e.posLength
}

type ImportMustHaveUppercaseIdentifierError struct {
	subError ParseError
}

func NewImportMustHaveUppercaseIdentifierError(subError ParseError) ImportMustHaveUppercaseIdentifierError {
	return ImportMustHaveUppercaseIdentifierError{subError: subError}
}

func (e ImportMustHaveUppercaseIdentifierError) Error() string {
	return fmt.Sprintf("import must be followed bu uppercase identifier (%v)", e.subError)
}

func (e ImportMustHaveUppercaseIdentifierError) FetchPositionLength() token.SourceFileReference {
	return e.subError.FetchPositionLength()
}

type ImportMustHaveUppercasePathError struct {
	subError ParseError
}

func NewImportMustHaveUppercasePathError(subError ParseError) ImportMustHaveUppercasePathError {
	return ImportMustHaveUppercasePathError{subError: subError}
}

func (e ImportMustHaveUppercasePathError) Error() string {
	return fmt.Sprintf("import must be followed bu uppercase identifier (%v)", e.subError)
}

func (e ImportMustHaveUppercasePathError) FetchPositionLength() token.SourceFileReference {
	return e.subError.FetchPositionLength()
}

type UnexpectedEndOfFileError struct {
	token.Token
}

func NewUnexpectedEndOfFileError(t token.Token) UnexpectedEndOfFileError {
	return UnexpectedEndOfFileError{Token: t}
}

func (e UnexpectedEndOfFileError) Error() string {
	return fmt.Sprintf("unexpected end of file %v", e.Token)
}

type ExpectedTypeIdentifierError struct {
	SubError
}

func NewExpectedTypeIdentifierError(subError tokenize.TokenError) ExpectedTypeIdentifierError {
	return ExpectedTypeIdentifierError{SubError: NewSubError(subError)}
}

func (e ExpectedTypeIdentifierError) Error() string {
	return fmt.Sprintf("expected a type identifier (start with uppercase) %v", e.SubError)
}

type ExpectedVariableIdentifierError struct {
	SubError
}

func NewExpectedVariableIdentifierError(err tokenize.TokenError) ExpectedVariableIdentifierError {
	return ExpectedVariableIdentifierError{SubError: SubError{err}}
}

func (e ExpectedVariableIdentifierError) Error() string {
	return fmt.Sprintf("expected a variable identifier (start with lowercase) %v", e.SubError)
}

type ExpectedSpacingAfterAnnotationOrDefinition struct {
	SubError
}

func NewExpectedSpacingAfterAnnotationOrDefinition(err tokenize.TokenError) ExpectedSpacingAfterAnnotationOrDefinition {
	return ExpectedSpacingAfterAnnotationOrDefinition{SubError: SubError{err}}
}

func (e ExpectedSpacingAfterAnnotationOrDefinition) Error() string {
	return fmt.Sprintf("expected ' ' after annotation or definition %v", e.SubError)
}

type ExpectedVariableAssignOrRecordUpdate struct {
	SubError
}

func NewExpectedVariableAssignOrRecordUpdate(err tokenize.TokenError) ExpectedVariableAssignOrRecordUpdate {
	return ExpectedVariableAssignOrRecordUpdate{SubError: SubError{err}}
}

func (e ExpectedVariableAssignOrRecordUpdate) Error() string {
	return fmt.Sprintf("expected an identifier followed by '=' or a variable followed by record update | %v", e.SubError)
}

type ExpectedVariableAssign struct {
	SubError
}

func NewExpectedVariableAssign(err tokenize.TokenError) ExpectedVariableAssign {
	return ExpectedVariableAssign{SubError: SubError{err}}
}

func (e ExpectedVariableAssign) Error() string {
	return fmt.Sprintf("expected an identifier followed by '=' %v", e.SubError)
}

type ExpectedBlockSpacing struct {
	SubError
}

func NewExpectedBlockSpacing(err tokenize.TokenError) ExpectedBlockSpacing {
	return ExpectedBlockSpacing{SubError: SubError{err}}
}

func (e ExpectedBlockSpacing) Error() string {
	return fmt.Sprintf("expected a block spacing %v", e.SubError)
}

type ExpectedContinuationLineOrOneSpace struct {
	posLength token.SourceFileReference
}

func NewExpectedContinuationLineOrOneSpace(posLength token.SourceFileReference) ExpectedContinuationLineOrOneSpace {
	return ExpectedContinuationLineOrOneSpace{posLength: posLength}
}

func (e ExpectedContinuationLineOrOneSpace) Error() string {
	return fmt.Sprintf("expected a new line with extra indentation or just a space (%v)", e.posLength)
}

func (e ExpectedContinuationLineOrOneSpace) FetchPositionLength() token.SourceFileReference {
	return e.posLength
}

type ExpectedRecordUpdate struct {
	SubError
}

func NewExpectedRecordUpdate(err tokenize.TokenError) ExpectedRecordUpdate {
	return ExpectedRecordUpdate{SubError: SubError{err}}
}

func (e ExpectedRecordUpdate) Error() string {
	return fmt.Sprintf("expected a variable followed by record update | %v", e.SubError)
}

type LeftPartOfPipeMustBeFunctionCallError struct {
	token.Token
}

func NewLeftPartOfPipeMustBeFunctionCallError(t token.Token) LeftPartOfPipeMustBeFunctionCallError {
	return LeftPartOfPipeMustBeFunctionCallError{Token: t}
}

func (e LeftPartOfPipeMustBeFunctionCallError) Error() string {
	return fmt.Sprintf("left part must be a function call %v", e.Token)
}

type RightPartOfPipeMustBeFunctionCallError struct {
	token.Token
}

func NewRightPartOfPipeMustBeFunctionCallError(t token.Token) RightPartOfPipeMustBeFunctionCallError {
	return RightPartOfPipeMustBeFunctionCallError{Token: t}
}

func (e RightPartOfPipeMustBeFunctionCallError) Error() string {
	return fmt.Sprintf("right part must be a function call %v", e.Token)
}

type ExpectedElseKeyword struct {
	SubError
}

func NewExpectedElseKeyword(eatError tokenize.TokenError) ExpectedElseKeyword {
	return ExpectedElseKeyword{SubError: SubError{eatError}}
}

func (e ExpectedElseKeyword) Error() string {
	return fmt.Sprintf("expected else keyword %v", e.SubError)
}

type ExpectedInKeyword struct {
	sourceReference token.SourceFileReference
}

func NewExpectedInKeyword(sourceReference token.SourceFileReference) ExpectedInKeyword {
	return ExpectedInKeyword{sourceReference: sourceReference}
}

func (e ExpectedInKeyword) Error() string {
	return fmt.Sprintf("expected IN keyword %v", e.sourceReference)
}

func (e ExpectedInKeyword) FetchPositionLength() token.SourceFileReference {
	return e.sourceReference
}

type ExpectedUniqueLetIdentifier struct {
	sourceReference token.SourceFileReference
}

func NewExpectedUniqueLetIdentifier(sourceReference token.SourceFileReference) ExpectedUniqueLetIdentifier {
	return ExpectedUniqueLetIdentifier{sourceReference: sourceReference}
}

func (e ExpectedUniqueLetIdentifier) Error() string {
	return fmt.Sprintf("expected unique let identifier %v", e.sourceReference)
}

func (e ExpectedUniqueLetIdentifier) FetchPositionLength() token.SourceFileReference {
	return e.sourceReference
}

type MissingElseExpression struct {
	other ParseError
}

func NewMissingElseExpression(other ParseError) MissingElseExpression {
	return MissingElseExpression{other: other}
}

func (e MissingElseExpression) Error() string {
	return fmt.Sprintf("missing else expression %v", e.other)
}

func (e MissingElseExpression) FetchPositionLength() token.SourceFileReference {
	return e.other.FetchPositionLength()
}

type UnknownStatement struct {
	token.Token
}

func NewUnknownStatement(t token.Token) UnknownStatement {
	return UnknownStatement{Token: t}
}

func (e UnknownStatement) Error() string {
	return fmt.Sprintf("unknown statement %v %T", e.Token, e.Token)
}

type UnknownPrefixInExpression struct {
	token.Token
}

func NewUnknownPrefixInExpression(t token.Token) UnknownPrefixInExpression {
	return UnknownPrefixInExpression{Token: t}
}

func (e UnknownPrefixInExpression) Error() string {
	return fmt.Sprintf("unknown prefix in expression %v %T", e.Token, e.Token)
}

type ExtraSpacing struct {
	posLength token.SourceFileReference
}

func NewExtraSpacing(posLength token.SourceFileReference) ExtraSpacing {
	return ExtraSpacing{posLength: posLength}
}

func (e ExtraSpacing) Error() string {
	return fmt.Sprintf("extra spacing")
}

func (e ExtraSpacing) FetchPositionLength() token.SourceFileReference {
	return e.posLength
}

type ExpectedOneSpaceOrExtraIndent struct {
	SubError
}

func NewExpectedOneSpaceOrExtraIndent(subError tokenize.TokenError) ExpectedOneSpaceOrExtraIndent {
	return ExpectedOneSpaceOrExtraIndent{SubError: NewSubError(subError)}
}

func (e ExpectedOneSpaceOrExtraIndent) Error() string {
	return fmt.Sprintf("expected either a single space or a new indented block %v", e.SubError)
}

type ExpectedOneSpaceAfterComma struct {
	SubError
}

func NewExpectedOneSpaceAfterComma(subError tokenize.TokenError) ExpectedOneSpaceAfterComma {
	return ExpectedOneSpaceAfterComma{SubError: NewSubError(subError)}
}

func (e ExpectedOneSpaceAfterComma) Error() string {
	return fmt.Sprintf("expected a single space after ',' (%v)", e.SubError)
}

type ExpectedOneSpaceAfterBinaryOperator struct {
	SubError
}

func NewExpectedOneSpaceAfterBinaryOperator(subError tokenize.TokenError) ExpectedOneSpaceAfterBinaryOperator {
	return ExpectedOneSpaceAfterBinaryOperator{SubError: NewSubError(subError)}
}

func (e ExpectedOneSpaceAfterBinaryOperator) Error() string {
	return fmt.Sprintf("expected a single space after binary operator (%v)", e.SubError)
}

type ExpectedOneSpaceAfterVariableAndBeforeAssign struct {
	SubError
}

func NewExpectedOneSpaceAfterVariableAndBeforeAssign(subError tokenize.TokenError) ExpectedOneSpaceAfterVariableAndBeforeAssign {
	return ExpectedOneSpaceAfterVariableAndBeforeAssign{SubError: NewSubError(subError)}
}

func (e ExpectedOneSpaceAfterVariableAndBeforeAssign) Error() string {
	return fmt.Sprintf("expected a single space after variable and before '=' (%v)", e.SubError)
}

type ExpectedOneSpace struct {
	sourceFileReference token.SourceFileReference
}

func NewExpectedOneSpace(sourceFileReference token.SourceFileReference) ExpectedOneSpace {
	return ExpectedOneSpace{sourceFileReference: sourceFileReference}
}

func (e ExpectedOneSpace) Error() string {
	return fmt.Sprintf("expected a single space here (%v)", e.sourceFileReference)
}

func (e ExpectedOneSpace) FetchPositionLength() token.SourceFileReference {
	return e.sourceFileReference
}

type ExpectedOneSpaceAfterAssign struct {
	SubError
}

func NewExpectedOneSpaceAfterAssign(subError tokenize.TokenError) ExpectedOneSpaceAfterAssign {
	return ExpectedOneSpaceAfterAssign{SubError: NewSubError(subError)}
}

func (e ExpectedOneSpaceAfterAssign) Error() string {
	return fmt.Sprintf("expected a single space after = (%v)", e.SubError)
}

type ExpectedOneSpaceOrExtraIndentCommaSeparator struct {
	SubError
}

func NewExpectedOneSpaceOrExtraIndentCommaSeparator(subError tokenize.TokenError) ExpectedOneSpaceOrExtraIndentCommaSeparator {
	return ExpectedOneSpaceOrExtraIndentCommaSeparator{SubError: NewSubError(subError)}
}

func (e ExpectedOneSpaceOrExtraIndentCommaSeparator) Error() string {
	return fmt.Sprintf("expected ',' and space or space and termination (%v)", e.SubError)
}

type ExpectedOneSpaceOrExtraIndentArgument struct {
	SubError
}

func NewExpectedOneSpaceOrExtraIndentArgument(subError tokenize.TokenError) ExpectedOneSpaceOrExtraIndentArgument {
	return ExpectedOneSpaceOrExtraIndentArgument{SubError: NewSubError(subError)}
}

func (e ExpectedOneSpaceOrExtraIndentArgument) Error() string {
	return fmt.Sprintf("expected either a single space or a new indented block and a call argument %v", e.SubError)
}

type LetInConsequenceOnSameColumn struct {
	SubError
}

func NewLetInConsequenceOnSameColumn(subError tokenize.TokenError) LetInConsequenceOnSameColumn {
	return LetInConsequenceOnSameColumn{SubError: NewSubError(subError)}
}

func (e LetInConsequenceOnSameColumn) Error() string {
	return fmt.Sprintf("the in consequence should be aligned on the same column as 'let' %v", e.SubError)
}

type OneSpaceAfterRecordTypeColon struct {
	SubError
}

func NewOneSpaceAfterRecordTypeColon(subError tokenize.TokenError) OneSpaceAfterRecordTypeColon {
	return OneSpaceAfterRecordTypeColon{SubError: NewSubError(subError)}
}

func (e OneSpaceAfterRecordTypeColon) Error() string {
	return fmt.Sprintf("one space after record type field colon %v", e.SubError)
}

type ParseAliasError struct {
	SubParserError
}

func NewParseAliasError(subError ParseError) ParseAliasError {
	return ParseAliasError{SubParserError: NewSubParserError(subError)}
}

func (e ParseAliasError) Error() string {
	return fmt.Sprintf("problem parsing alias %v", e.SubParserError)
}

type ExpectedDefaultLastError struct {
	expression ast.Expression
}

func NewExpectedDefaultLastError(expression ast.Expression) ExpectedDefaultLastError {
	return ExpectedDefaultLastError{expression: expression}
}

func (e ExpectedDefaultLastError) Error() string {
	return fmt.Sprintf("default '_' must come last of the conditions. %v", e.expression)
}

func (e ExpectedDefaultLastError) FetchPositionLength() token.SourceFileReference {
	return e.expression.FetchPositionLength()
}

type MustHaveDefaultInConditionsError struct {
	expression ast.Expression
}

func NewMustHaveDefaultInConditionsError(expression ast.Expression) MustHaveDefaultInConditionsError {
	return MustHaveDefaultInConditionsError{expression: expression}
}

func (e MustHaveDefaultInConditionsError) Error() string {
	return fmt.Sprintf("must include a default '_' in conditions. %v", e.expression)
}

func (e MustHaveDefaultInConditionsError) FetchPositionLength() token.SourceFileReference {
	return e.expression.FetchPositionLength()
}
