/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

type ParserInterface interface {
	parseExpression(precedence Precedence, startIndentation int) (ast.Expression, parerr.ParseError)
	parseExpressionNormal(startIndentation int) (ast.Expression, parerr.ParseError)
	parseTerm(startIndentation int) (ast.Expression, parerr.ParseError)
}

type ParseStreamImpl struct {
	tokenizer           *tokenize.Tokenizer
	descent             int
	parser              ParserInterface
	disableEnforceStyle bool
}

func NewParseStreamImpl(parser ParserInterface, tokenizer *tokenize.Tokenizer, enforceStyle bool) *ParseStreamImpl {
	if parser == nil {
		panic("must have parser")
	}
	p := &ParseStreamImpl{tokenizer: tokenizer, parser: parser, disableEnforceStyle: !enforceStyle}
	return p
}

func (p *ParseStreamImpl) debugInfo(s string) {
	extract := p.tokenizer.DebugInfo()
	fmt.Printf("*-- %s: (%d) %v\n", s, p.descent, p.tokenizer.ParsingPosition().Position())
	fmt.Printf("%v\n---\n", extract)
}

func (p *ParseStreamImpl) debugInfoRows(s string, rowCount int) {
	extract := p.tokenizer.DebugInfoLinesWithComment(s, rowCount)
	fmt.Printf("*-- %s: (%d)\n", s, p.descent)
	fmt.Printf("%v\n---\n", extract)
}

func (p *ParseStreamImpl) positionLength() token.PositionLength {
	pos := p.tokenizer.ParsingPosition().Position()
	return token.NewPositionLength(pos, 1, pos.Column()/2)
}

func (p *ParseStreamImpl) getPrecedenceFromType(tType token.Type) Precedence {
	if p, ok := precedences[tType]; ok {
		return p
	}

	return LOWEST
}

func (p *ParseStreamImpl) getPrecedenceFromToken(t token.Token) Precedence {
	return p.getPrecedenceFromType(t.Type())
}

func (p *ParseStreamImpl) readOperatorToken() (token.Token, parerr.ParseError) {
	return p.tokenizer.ParseOperator()
}

func (p *ParseStreamImpl) readTypeIdentifier() (*ast.TypeIdentifier, parerr.ParseError) {
	typeSymbol, typeSymbolErr := p.tokenizer.ParseTypeSymbol()
	if typeSymbolErr != nil {
		return nil, parerr.NewExpectedTypeIdentifierError(typeSymbolErr)
	}
	typeIdent := ast.NewTypeIdentifier(typeSymbol)
	return typeIdent, nil
}


func (p *ParseStreamImpl) addWarning(description string, length token.PositionLength) {
	color.Yellow("%v:%d:%d: %v: %v", p.tokenizer.RelativeFilename(), length.Position().Line()+1,
		length.Position().Column()+1, "Warning", description)
}

func (p *ParseStreamImpl) readVariableIdentifier() (*ast.VariableIdentifier, parerr.ParseError) {
	variableSymbol, variableSymbolErr := p.tokenizer.ParseVariableSymbol()
	if variableSymbolErr != nil {
		return nil, parerr.NewExpectedVariableIdentifierError(variableSymbolErr)
	}
	varIdent := ast.NewVariableIdentifier(variableSymbol)
	return varIdent, nil
}

func (p *ParseStreamImpl) skipMaybeSpaceAndSameIndentationOrContinuation() (token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return report, err
	}

	if !p.disableEnforceStyle {
		if report.NewLineCount == 0 {
			if report.SpacesUntilMaybeNewline == 1 || report.SpacesUntilMaybeNewline == 0 {
				return report, nil
			}
			return report, parerr.NewInternalError(report.PositionLength, fmt.Errorf("must be space or continuation"))
		}

		if report.ExactIndentation == report.PreviousExactIndentation || report.ExactIndentation == report.PreviousExactIndentation+1 {
			return report, nil
		}

		return report, parerr.NewInternalError(report.PositionLength, fmt.Errorf("must be space or continuation3"))

	}
	if report.NewLineCount == 0 || report.CloseIndentation >= report.PreviousCloseIndentation {
		return report, nil
	}

	return report, parerr.NewInternalError(report.PositionLength, fmt.Errorf("must be space or continuation2"))
}

func (p *ParseStreamImpl) maybeOneSpace() (bool, token.IndentationReport, parerr.ParseError) {
	save := p.tokenizer.Tell()
	report, foundPosLengthErr := p.tokenizer.SkipWhitespaceToNextIndentation()
	if foundPosLengthErr != nil {
		return false, report, foundPosLengthErr
	}

	if p.disableEnforceStyle {
		if tokenize.LegalContinuationSpace(report, !p.disableEnforceStyle) {
			return true, report, nil
		}
	}

	if report.SpacesUntilMaybeNewline == 1 {
		return true, report, nil
	}

	p.tokenizer.Seek(save)
	return false, report, nil
}

func (p *ParseStreamImpl) eatContinuationOrNoSpaceInternal() (bool, token.IndentationReport, parerr.ParseError) {
	report, foundPosLengthErr := p.tokenizer.SkipWhitespaceToNextIndentation()
	if foundPosLengthErr != nil {
		return false, report, foundPosLengthErr
	}

	if report.SpacesUntilMaybeNewline == 0 || report.SpacesUntilMaybeNewline == 1 {
		return true, report, nil
	}

	if tokenize.LegalContinuationSpace(report, !p.disableEnforceStyle) {
		return true, report, nil
	}
	return false, report, parerr.NewExpectedContinuationLineOrOneSpace(report.PositionLength)
}

func (p *ParseStreamImpl) maybeColon() bool {
	return p.tokenizer.MaybeRune(':')
}

func (p *ParseStreamImpl) maybeComma() bool {
	return p.tokenizer.MaybeRune(',')
}

func (p *ParseStreamImpl) maybePipeLeft() bool {
	return p.tokenizer.MaybeString("<|")
}

func (p *ParseStreamImpl) maybeRightBracket() bool {
	return p.tokenizer.MaybeRune(']')
}

func (p *ParseStreamImpl) maybeRightCurly() bool {
	return p.tokenizer.MaybeRune('}')
}

func (p *ParseStreamImpl) maybeRightArrow() bool {
	return p.tokenizer.MaybeString("->")
}

func (p *ParseStreamImpl) maybeOneSpaceAndRightArrow() bool {
	pos := p.tokenizer.Tell()
	report, _ := p.tokenizer.SkipWhitespaceToNextIndentation()
	if !p.disableEnforceStyle {
		if report.NewLineCount != 0 || report.SpacesUntilMaybeNewline != 1 {
			p.tokenizer.Seek(pos)
			return false
		}
	}
	foundArrow := p.tokenizer.MaybeString("->")
	if !foundArrow {
		p.tokenizer.Seek(pos)
	}

	return foundArrow
}

func (p *ParseStreamImpl) maybeLeftParen() bool {
	return p.tokenizer.MaybeRune('(')
}

func (p *ParseStreamImpl) maybeLeftCurly() bool {
	return p.tokenizer.MaybeString("{")
}

func (p *ParseStreamImpl) maybeAccessor() bool {
	return p.tokenizer.MaybeAccessor()
}

func (p *ParseStreamImpl) readVariableIdentifierAssignOrUpdate() (*ast.VariableIdentifier, bool, parerr.ParseError) {
	ident, identErr := p.readVariableIdentifier()
	if identErr != nil {
		return nil, false, identErr
	}

	_, spaceAfterIdentifierErr := p.eatOneSpace("space after variableIdentifier assign or update")
	if spaceAfterIdentifierErr != nil {
		return nil, false, parerr.NewExpectedVariableAssignOrRecordUpdate(spaceAfterIdentifierErr)
	}

	wasAssign := p.tokenizer.MaybeRune('=')
	if wasAssign {
		_, spaceAfterAssignErr := p.eatOneSpace(" space after assign =")
		if spaceAfterAssignErr != nil {
			return nil, false, parerr.NewExpectedVariableAssign(spaceAfterAssignErr)
		}
		return ident, true, nil
	}

	eatUpdateErr := p.eatRune('|')
	if eatUpdateErr != nil {
		return nil, false, parerr.NewExpectedRecordUpdate(eatUpdateErr)
	}
	_, spaceAfterUpdateErr := p.eatOneSpace("space after update |")
	if spaceAfterUpdateErr != nil {
		return nil, false, spaceAfterUpdateErr
	}
	return ident, false, nil
}

func (p *ParseStreamImpl) guessToken() (token.Token, parerr.ParseError) {
	pos := p.tokenizer.ParsingPosition()
	t, tErr := p.tokenizer.GuessNext()
	if tErr != nil {
		return nil, tErr
	}

	switch t.(type) {
	case token.SpaceToken:
		posLength := p.tokenizer.MakePositionLength(pos)
		return nil, parerr.NewExtraSpacing(posLength)
	}

	return t, nil
}

// ---------------------------------------------------------------------------------
// DETECT
// ---------------------------------------------------------------------------------

func (p *ParseStreamImpl) detectOneSpaceAndTermination() bool {
	pos := p.tokenizer.Tell()

	if p.maybeNewLine() {
		p.tokenizer.Seek(pos)
		return true
	}

	if p.maybeComma() {
		p.tokenizer.Seek(pos)
		return true
	}

	if p.maybeRightParen() {
		p.tokenizer.Seek(pos)
		return true
	}

	_, eatOneSpaceErr := p.eatOneSpace("eat one space detectOneSpaceAndTermination")
	if eatOneSpaceErr != nil {
		p.tokenizer.Seek(pos)
		return true
	}

	if p.detectOneSpaceAndRightArrow() {
		p.tokenizer.Seek(pos)
		return true
	}

	if p.maybeRightCurly() {
		p.tokenizer.Seek(pos)
		return true
	}

	if p.maybeRightArrow() {
		p.tokenizer.Seek(pos)
		return true
	}

	p.tokenizer.Seek(pos)
	return false
}

func (p *ParseStreamImpl) maybeNewLineContinuation(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	pos := p.tokenizer.Tell()
	report, foundPosLengthErr := p.tokenizer.SkipWhitespaceToNextIndentation()
	if foundPosLengthErr != nil {
		return false, report, nil
	}

	if report.NewLineCount == 0 {
		return false, report, nil
	}
	if p.disableEnforceStyle {
		if report.CloseIndentation < expectedIndentation {
			p.tokenizer.Seek(pos)
			return false, report, nil
		}
		if report.CloseIndentation == expectedIndentation {
			return true, report, nil
		}
		if tokenize.LegalContinuationSpace(report, !p.disableEnforceStyle) {
			return true, report, nil
		} else {
			p.tokenizer.Seek(pos)
			return false, report, nil
		}
	} else {
		if report.ExactIndentation == -1 {
			return false, report, parerr.NewExpectedIndentationError(report.PositionLength, expectedIndentation)
		}

		if report.ExactIndentation > expectedIndentation {
			p.tokenizer.Seek(pos)
			return false, report, parerr.NewExpectedIndentationError(report.PositionLength, expectedIndentation)
		}
		if report.ExactIndentation < expectedIndentation {
			p.tokenizer.Seek(pos)
			return false, report, nil
		}
	}

	return true, report, nil
}

func (p *ParseStreamImpl) maybeNewLineContinuationWithExtraEmptyLine(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	save := p.tokenizer.Tell()
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return false, report, nil
	}

	if p.disableEnforceStyle {
		if report.CloseIndentation == expectedIndentation {
			return true, report, nil
		}
		if tokenize.LegalContinuationSpace(report, !p.disableEnforceStyle) {
			return true, report, nil
		}
	} else {
		if report.ExactIndentation == expectedIndentation && report.NewLineCount == 2 {
			return true, report, nil
		}
	}

	p.tokenizer.Seek(save)
	return false, report, nil
}

func (p *ParseStreamImpl) maybeUpToOneLogicalSpaceForContinuation(currentIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	save := p.tokenizer.Tell()
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		p.tokenizer.Seek(save)
		return false, report, err
	}

	if p.disableEnforceStyle {
		if tokenize.LegalContinuationSpaceIndentation(report, currentIndentation, !p.disableEnforceStyle) {
			return true, report, nil
		}
	} else {
		if report.ExactIndentation == currentIndentation+1 {
			return true, report, nil
		}
	}

	if report.SpacesUntilMaybeNewline == 0 || report.SpacesUntilMaybeNewline == 1 {
		return true, report, nil
	}

	p.tokenizer.Seek(save)
	return false, report, nil
}

func (p *ParseStreamImpl) detectEndOfCallOperator() bool {
	save := p.tokenizer.Tell()
	detectedEndOfCallOperator := false
	term, termErr := p.guessToken()
	_, isEOF := term.(*tokenize.EndOfFile)
	if isEOF {
		detectedEndOfCallOperator = true
	} else {
		if termErr == nil {
			t := term.Type()
			_, isOperator := term.(token.OperatorToken)
			detectedEndOfCallOperator = t == token.Else || t == token.RightCurlyBrace || isOperator || t == token.RightBracket || t == token.RightParen || t == token.Comma
		}
	}
	p.tokenizer.Seek(save)

	return detectedEndOfCallOperator
}

func (p *ParseStreamImpl) detectOneSpaceAndRightArrow() bool {
	pos := p.tokenizer.Tell()
	_, err := p.eatOneSpace("detect before right arrow")
	if err != nil {
		p.tokenizer.Seek(pos)
		return false
	}
	return p.tokenizer.DetectString("->")
}

func (p *ParseStreamImpl) detectAssign() bool {
	return p.tokenizer.DetectRune('=')
}

func (p *ParseStreamImpl) detectNewLine() bool {
	save := p.tokenizer.Tell()
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	p.tokenizer.Seek(save)
	if err != nil {
		return false
	}

	return report.NewLineCount > 0
}

func (p *ParseStreamImpl) detectNewLineOrSpace() (bool, bool) {
	save := p.tokenizer.Tell()
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	p.tokenizer.Seek(save)
	if err != nil {
		return false, false
	}

	if report.NewLineCount > 0 {
		return true, true
	}

	if report.SpacesUntilMaybeNewline == 1 {
		return true, false
	}

	return false, false
}

// ---------------------------------------------------------------------------------
// WAS
// ---------------------------------------------------------------------------------
func (p *ParseStreamImpl) wasDefaultSymbol() (token.RuneToken, bool) {
	return p.tokenizer.WasDefaultSymbol()
}

func (p *ParseStreamImpl) wasVariableIdentifier() (*ast.VariableIdentifier, bool) {
	variableSymbol, parseErr := p.tokenizer.ParseVariableSymbol()
	if parseErr != nil {
		return nil, false
	}
	varIdent := ast.NewVariableIdentifier(variableSymbol)
	return varIdent, true
}

func (p *ParseStreamImpl) wasTypeIdentifier() (*ast.TypeIdentifier, bool) {
	typeSymbol, parseErr := p.tokenizer.ParseTypeSymbol()
	if parseErr != nil {
		return nil, false
	}
	typeIdent := ast.NewTypeIdentifier(typeSymbol)
	return typeIdent, true
}

// ---------------------------------------------------------------------------------
// MAYBE
// ---------------------------------------------------------------------------------
func (p *ParseStreamImpl) maybeKeywordAlias() bool {
	variableIdentifier, wasVariable := p.wasVariableIdentifier()
	if !wasVariable {
		return false
	}
	return variableIdentifier.Name() == "alias"
}

func (p *ParseStreamImpl) maybeKeywordExposing() bool {
	variableIdentifier, wasVariable := p.wasVariableIdentifier()
	if !wasVariable {
		return false
	}
	return variableIdentifier.Name() == "exposing"
}

func (p *ParseStreamImpl) maybeKeywordAs() bool {
	variableIdentifier, wasVariable := p.wasVariableIdentifier()
	if !wasVariable {
		return false
	}
	return variableIdentifier.Name() == "as"
}


func (p *ParseStreamImpl) maybeAssign() bool {
	return p.tokenizer.MaybeAssign()
}

func (p *ParseStreamImpl) maybeNewLine() bool {
	return p.tokenizer.MaybeOneNewLine()
}

func (p *ParseStreamImpl) maybeRightParen() bool {
	return p.tokenizer.MaybeRune(')')
}

// ---------------------------------------------------------------------------------
// EAT
// ---------------------------------------------------------------------------------
func (p *ParseStreamImpl) eatRune(r rune) tokenize.TokenError {
	return p.tokenizer.EatRune(r)
}
func (p *ParseStreamImpl) eatString(s string) tokenize.TokenError {
	return p.tokenizer.EatString(s)
}

func (p *ParseStreamImpl) eatOneSpace(reason string) (token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return report, err
	}

	if p.disableEnforceStyle {
		if tokenize.LegalContinuationSpace(report, !p.disableEnforceStyle) {
			return report, nil
		} else {
			subErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, report.PreviousCloseIndentation, report.CloseIndentation)
			return report, parerr.NewExpectedOneSpaceOrExtraIndent(subErr)
		}
	}

	if report.SpacesUntilMaybeNewline != 1 {
		return report, parerr.NewInternalError(report.PositionLength, fmt.Errorf("expected one space %v", reason))
	}

	return report, nil
}


func (p *ParseStreamImpl) eatNewLinesAfterStatement(count int) (token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return report, err
	}

	if report.EndOfFile {
		return report, nil
	}

	if p.disableEnforceStyle {
		if report.NewLineCount > 0 {
			return report, nil
		}
	} else {
		if count == -1 {
			if (report.NewLineCount >= 1) && report.IndentationSpaces == 0 {
				return report, nil
			}
		} else {
			if (report.NewLineCount == count) && report.IndentationSpaces == 0 {
				return report, nil
			}
		}
	}

	return report, parerr.NewInternalError(report.PositionLength, fmt.Errorf("wrong exact number of line %v expected %v", report.NewLineCount, count))
}



func (p *ParseStreamImpl) eatCommaSeparatorOrTermination(expectedIndentation int, allowComments bool) (bool, token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentationHelper(allowComments)
	if err != nil {
		return false, report, err
	}
	wasComma := p.tokenizer.MaybeRune(',')
	if wasComma {
		spaceAfterCommaReport, eatErr := p.eatOneSpace("space after comma")
		if eatErr != nil {
			return true, spaceAfterCommaReport, eatErr
		}

		if p.disableEnforceStyle {
			return true, report, nil
		}

		wasNextLineOrNoSpace := tokenize.LegalSameIndentationOrNoSpace(report, expectedIndentation, !p.disableEnforceStyle)
		if wasNextLineOrNoSpace {
			return true, report, nil
		}

		return false, report, parerr.NewExtraSpacing(p.positionLength())
	}

	// It was a termination, that should be handled, but we can still check that the space is correct
	if p.disableEnforceStyle {
		if tokenize.LegalSameIndentationOrOptionalOneSpace(report, expectedIndentation, !p.disableEnforceStyle) {
			return false, report, nil
		}
		fmt.Printf("%d\n", report.SpacesUntilMaybeNewline)
		p.debugInfo("wasnt a space indentation")
	} else {
		wasNextLineOrOneSpace := tokenize.LegalSameIndentationOrOptionalOneSpace(report, expectedIndentation, !p.disableEnforceStyle)
		if wasNextLineOrOneSpace {
			return false, report, nil
		}
	}

	tokenErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, expectedIndentation, report.CloseIndentation)
	return false, report, parerr.NewExpectedOneSpaceOrExtraIndentCommaSeparator(tokenErr)
}

func (p *ParseStreamImpl) eatContinuationReturnIndentation(indentation int) (int, token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return -1, report, err
	}

	if p.disableEnforceStyle {
		if report.NewLineCount == 0 {
			return -1, report, nil
		}
		if tokenize.LegalContinuationSpace(report, !p.disableEnforceStyle) {
			return report.CloseIndentation, report, nil
		}
	} else {
		if report.SpacesUntilMaybeNewline == 1 {
			return indentation, report, nil
		}

		if report.NewLineCount == 1 && report.ExactIndentation == indentation+1 {
			return report.ExactIndentation, report, nil
		}
	}

	subErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, indentation, report.CloseIndentation)

	tokenErr := tokenize.NewExpectedIndentationAfterNewLineError(indentation, report.CloseIndentation, subErr)

	return -1, report, parerr.NewExpectedOneSpaceOrExtraIndent(tokenErr)
}

func (p *ParseStreamImpl) eatTwoNewLineContinuationOrDedent(currentIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return false, report, err
	}

	if p.disableEnforceStyle {
		wasContinuation := report.CloseIndentation >= currentIndentation
		return wasContinuation, report, nil
	}

	if report.NewLineCount == 2 && report.ExactIndentation == currentIndentation {
		wasContinuation := report.ExactIndentation == currentIndentation
		return wasContinuation, report, nil
	}

	if report.NewLineCount == 1 && report.ExactIndentation == currentIndentation-1 {
		wasContinuation := report.ExactIndentation == currentIndentation
		return wasContinuation, report, nil
	}

	subErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, currentIndentation, report.CloseIndentation)

	return false, report, parerr.NewExpectedOneSpaceOrExtraIndent(subErr)
}

func (p *ParseStreamImpl) eatArgumentSpaceOrDetectEndOfArguments(currentIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	saveInfo := p.tokenizer.Tell()
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return false, report, err
	}

	if p.disableEnforceStyle {
		if report.NewLineCount > 0 { // Arguments end even if they are indented
			p.tokenizer.Seek(saveInfo)
			return true, report, nil
		}
		wasEnd := p.detectEndOfCallOperator()
		if wasEnd {
			p.tokenizer.Seek(saveInfo)
		}
		return wasEnd, report, nil
	} else {
		if report.NewLineCount == 1 && report.ExactIndentation == currentIndentation+1 {
			return true, report, nil
		}
		if report.NewLineCount == 0 && (report.SpacesUntilMaybeNewline == 1 || report.SpacesUntilMaybeNewline == 0) {
			wasEnd := p.detectEndOfCallOperator()
			if wasEnd {
				p.tokenizer.Seek(saveInfo)
			}
			return wasEnd, report, nil
		}
		if report.NewLineCount > 0 { // Arguments end even if they are indented
			p.tokenizer.Seek(saveInfo)
			return true, report, nil
		}
	}

	p.tokenizer.Seek(saveInfo)
	subErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, currentIndentation, report.CloseIndentation)
	return false, report, parerr.NewExpectedOneSpaceOrExtraIndentArgument(subErr)
}

func (p *ParseStreamImpl) eatNewLineAndExactIndent(indentation int) (token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return report, err
	}

	if p.disableEnforceStyle {
		return report, nil
	}

	if report.NewLineCount == 1 && report.ExactIndentation == indentation {
		return report, nil
	}

	subErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, indentation, report.ExactIndentation)
	return report, subErr
}

func (p *ParseStreamImpl) eatNewLineAndExactIndentExtraLine(indentation int) (token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return report, err
	}

	if p.disableEnforceStyle {
		if tokenize.LegalContinuationSpace(report, !p.disableEnforceStyle) {
			return report, nil
		}
	}

	if indentation == report.CloseIndentation {
		return report, nil
	}

	subErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, indentation, report.CloseIndentation)

	return report, parerr.NewExpectedOneSpaceOrExtraIndent(subErr)
}


func (p *ParseStreamImpl) eatNewLineContinuation(indentation int) (token.IndentationReport, parerr.ParseError) {
	return p.eatNewLineAndExactIndent(indentation + 1)
}

func (p *ParseStreamImpl) eatOperatorUpdate() parerr.ParseError {
	eatUpdateErr := p.eatRune('|')
	return eatUpdateErr
}

func (p *ParseStreamImpl) eatBlockSpacing(isBlock bool, requestedIndentation int) (token.IndentationReport, parerr.ParseError) {
	if isBlock {
		comment, eatOneLineAndIndentErr := p.eatNewLineAndExactIndent(requestedIndentation)
		if eatOneLineAndIndentErr != nil {
			return comment, parerr.NewExpectedBlockSpacing(eatOneLineAndIndentErr)
		}
		return comment, nil
	}
	return p.eatOneSpace("space after block spacing")
}

func (p *ParseStreamImpl) eatBlockSpacingOneExtraLine(isBlock bool, requestedIndentation int) (token.IndentationReport, parerr.ParseError) {
	if isBlock {
		comment, eatOneLineAndIndentErr := p.eatNewLineAndExactIndentExtraLine(requestedIndentation)
		if eatOneLineAndIndentErr != nil {
			return comment, parerr.NewExpectedBlockSpacing(eatOneLineAndIndentErr)
		}
		return comment, nil
	}
	return p.eatOneSpace("eatBlockSpacingOneExtraLine")
}

func (p *ParseStreamImpl) eatRightArrow() parerr.ParseError {
	err := p.eatString("->")
	if err != nil {
		return parerr.NewExpectedRightArrowError(err)
	}
	return nil
}

func (p *ParseStreamImpl) eatRightParen() parerr.ParseError {
	return p.eatRune(')')
}

func (p *ParseStreamImpl) eatRightCurly() parerr.ParseError {
	return p.eatRune('}')
}

func (p *ParseStreamImpl) eatRightBracket() parerr.ParseError {
	return p.eatRune(']')
}

func (p *ParseStreamImpl) eatOf() parerr.ParseError {
	return p.eatString("of")
}

func (p *ParseStreamImpl) eatThen() parerr.ParseError {
	return p.eatString("then")
}

func (p *ParseStreamImpl) eatElse() parerr.ParseError {
	return p.eatString("else")
}

func (p *ParseStreamImpl) eatAccessor() parerr.ParseError {
	return p.eatRune('.')
}

func (p *ParseStreamImpl) eatColon() parerr.ParseError {
	return p.eatRune(':')
}

func (p *ParseStreamImpl) eatAssign() parerr.ParseError {
	return p.eatRune('=')
}

func (p *ParseStreamImpl) eatIn() parerr.ParseError {
	return p.eatString("in")
}

// ---------------------------------------------------------------------------------
// PARSER
// ---------------------------------------------------------------------------------
func (p *ParseStreamImpl) parseExpression(precedence Precedence, startIndentation int) (ast.Expression, parerr.ParseError) {
	p.descent++
	r, rErr := p.parser.parseExpression(precedence, startIndentation)
	p.descent--
	return r, rErr
}

func (p *ParseStreamImpl) parseExpressionNormal(startIndentation int) (ast.Expression, parerr.ParseError) {
	p.descent++
	r, rErr := p.parser.parseExpressionNormal(startIndentation)
	p.descent--
	return r, rErr
}

func (p *ParseStreamImpl) parseTerm(startIndentation int) (ast.Expression, parerr.ParseError) {
	return p.parser.parseTerm(startIndentation)
}
