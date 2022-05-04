/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"
	"os"

	"github.com/swamp/compiler/src/ast"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

type ParserInterface interface {
	parseExpression(precedence Precedence, startIndentation int) (ast.Expression, parerr.ParseError)
	parseExpressionNormal(startIndentation int) (ast.Expression, parerr.ParseError)
	parseExpressionNormalWithComment(startIndentation int, comment token.Comment) (ast.Expression, parerr.ParseError)
	parseTerm(startIndentation int) (ast.Expression, parerr.ParseError)
}

type ParseStreamImpl struct {
	tokenizer           *tokenize.Tokenizer
	descent             int
	parser              ParserInterface
	disableEnforceStyle bool
	nodes               []ast.Node
	warnings            []parerr.ParseError
	errors              []parerr.ParseError
}

func NewParseStreamImpl(parser ParserInterface, tokenizer *tokenize.Tokenizer, enforceStyle bool) *ParseStreamImpl {
	if parser == nil {
		panic("must have parser")
	}
	p := &ParseStreamImpl{tokenizer: tokenizer, parser: parser, disableEnforceStyle: !enforceStyle}

	return p
}

func (p *ParseStreamImpl) SkipToNextLineWithSameOrLowerIndent(indentation int) (int, tokenize.TokenError) {
	for !p.tokenizer.MaybeEOF() {
		report, tokenErr := p.tokenizer.SkipWhitespaceToNextIndentation()
		if tokenErr != nil {
			return -1, tokenErr
		}

		if report.EndOfFile {
			return -1, tokenize.NewEncounteredEOF()
		}

		if report.NewLineCount == 0 && report.SpacesUntilMaybeNewline == 0 {
			return -1, nil
		}

		if report.ExactIndentation <= indentation {
			return report.ExactIndentation, nil
		}
	}

	return -1, tokenize.NewEncounteredEOF()
}

func (p *ParseStreamImpl) debugInfo(s string) {
	extract := p.tokenizer.DebugInfo()
	fmt.Fprintf(os.Stderr, "*-- %s: (%d) %v\n", s, p.descent, p.tokenizer.ParsingPosition().Position())
	fmt.Fprintf(os.Stderr, "%v\n---\n", extract)
}

func (p *ParseStreamImpl) debugInfoRows(s string, rowCount int) {
	extract := p.tokenizer.DebugInfoLinesWithComment(s, rowCount)
	fmt.Fprintf(os.Stderr, "*-- %s: (%d)\n", s, p.descent)
	fmt.Fprintf(os.Stderr, "%v\n---\n", extract)
}

func (p *ParseStreamImpl) positionLength() token.SourceFileReference {
	return p.sourceFileReference()
}

func (p *ParseStreamImpl) sourceFileReference() token.SourceFileReference {
	pos := p.tokenizer.ParsingPosition()
	reference := p.tokenizer.MakeSourceFileReference(pos)
	return reference
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

func (p *ParseStreamImpl) AddWarning(parseError parerr.ParseError) {
	p.warnings = append(p.warnings, parseError)
}

func (p *ParseStreamImpl) AddError(parseError parerr.ParseError) {
	p.errors = append(p.errors, parseError)
}

func (p *ParseStreamImpl) Warnings() []parerr.ParseError {
	return p.warnings
}

func (p *ParseStreamImpl) Errors() []parerr.ParseError {
	return p.errors
}

func (p *ParseStreamImpl) readVariableIdentifierInternal() (*ast.VariableIdentifier, parerr.ParseError) {
	variableSymbol, variableSymbolErr := p.tokenizer.ParseVariableSymbol()
	if variableSymbolErr != nil {
		return nil, parerr.NewExpectedVariableIdentifierError(variableSymbolErr)
	}
	varIdent := ast.NewVariableIdentifier(variableSymbol)
	return varIdent, nil
}

func (p *ParseStreamImpl) readVariableIdentifier() (*ast.VariableIdentifier, parerr.ParseError) {
	variableIdentifier, variableIdentifierErr := p.readVariableIdentifierInternal()
	if variableIdentifierErr != nil {
		p.tokenizer.DebugPrint("went wrong here")

		return nil, variableIdentifierErr
	}

	return variableIdentifier, nil
}

func (p *ParseStreamImpl) readKeyword() (token.Keyword, parerr.ParseError) {
	pos := p.tokenizer.ParsingPosition()
	ident, err := p.readVariableIdentifierInternal()
	if err != nil {
		return token.Keyword{}, err
	}

	keyword, keywordErr := tokenize.DetectLowercaseKeyword(ident.Symbol())
	if keywordErr != nil {
		return token.Keyword{}, parerr.NewUnknownKeywordError(p.tokenizer.MakeSourceFileReference(pos), keywordErr)
	}

	return keyword, nil
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
			return report, parerr.NewMustBeSpaceOrContinuation(report.PositionLength)
		}

		if report.ExactIndentation == report.PreviousExactIndentation || report.ExactIndentation == report.PreviousExactIndentation+1 {
			return report, nil
		}

		return report, parerr.NewMustBeSpaceOrContinuation(report.PositionLength)
	}

	if report.NewLineCount == 0 || report.CloseIndentation >= report.PreviousCloseIndentation {
		return report, nil
	}

	return report, parerr.NewMustBeSpaceOrContinuation(report.PositionLength)
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

func (p *ParseStreamImpl) maybePipeLeft() bool {
	return p.tokenizer.MaybeString("<|")
}

func (p *ParseStreamImpl) readSpecificParenToken(p2 token.Type) (token.ParenToken, parerr.ParseError) {
	parsedToken, err := p.tokenizer.ReadAnyOperator()
	if err != nil {
		return token.ParenToken{}, err
	}

	if parsedToken.Type() != p2 {
		return token.ParenToken{}, parerr.NewUnknownPrefixInExpression(parsedToken)
	}

	return parsedToken.(token.ParenToken), nil
}

func (p *ParseStreamImpl) readSpecificOperatorToken(p2 token.Type) (token.OperatorToken, parerr.ParseError) {
	parsedToken, err := p.tokenizer.ReadAnyOperator()
	if err != nil {
		return token.OperatorToken{}, err
	}

	if parsedToken.Type() != p2 {
		return token.OperatorToken{}, parerr.NewUnknownPrefixInExpression(parsedToken)
	}

	return parsedToken.(token.OperatorToken), nil
}

func (p *ParseStreamImpl) maybeSpecificKeywordToken(p2 token.Type) (token.Keyword, bool) {
	pos := p.tokenizer.Tell()
	parsedToken, err := p.readSpecificKeywordToken(p2)
	if err != nil {
		p.tokenizer.Seek(pos)
		return token.Keyword{}, false
	}

	return parsedToken, true
}

func (p *ParseStreamImpl) readSpecificKeywordToken(p2 token.Type) (token.Keyword, parerr.ParseError) {
	parsedToken, err := p.readKeyword()
	if err != nil {
		return token.Keyword{}, err
	}

	if parsedToken.Type() != p2 {
		return token.Keyword{}, parerr.NewUnknownPrefixInExpression(parsedToken)
	}

	return parsedToken, nil
}

func (p *ParseStreamImpl) maybeSpecificParenToken(p2 token.Type) (token.ParenToken, bool) {
	pos := p.tokenizer.Tell()

	parenToken, err := p.readSpecificParenToken(p2)
	if err != nil {
		p.tokenizer.Seek(pos)
		return token.ParenToken{}, false
	}

	return parenToken, true
}

func (p *ParseStreamImpl) maybeSpecificOperatorToken(p2 token.Type) (token.OperatorToken, bool) {
	pos := p.tokenizer.Tell()

	operatorToken, err := p.readSpecificOperatorToken(p2)
	if err != nil {
		p.tokenizer.Seek(pos)
		return token.OperatorToken{}, false
	}

	return operatorToken, true
}

func (p *ParseStreamImpl) maybeRightSquareBracket() (token.ParenToken, bool) {
	return p.maybeSpecificParenToken(token.RightSquareBracket)
}

func (p *ParseStreamImpl) maybeComma() (token.OperatorToken, bool) {
	return p.maybeSpecificOperatorToken(token.Comma)
}

func (p *ParseStreamImpl) maybeSpacingAndComma(indentation int) (token.OperatorToken, bool) {
	p.maybeNewLineContinuationHelper(indentation, tokenize.NotAllowedAtAll)
	return p.maybeSpecificOperatorToken(token.Comma)
}

func (p *ParseStreamImpl) maybeAsterisk() (token.OperatorToken, bool) {
	return p.maybeSpecificOperatorToken(token.OperatorMultiply)
}

func (p *ParseStreamImpl) maybeRightArrayBracket() (token.ParenToken, bool) {
	return p.maybeSpecificParenToken(token.RightArrayBracket)
}

func (p *ParseStreamImpl) maybeLeftCurly() (token.ParenToken, bool) {
	return p.maybeSpecificParenToken(token.LeftCurlyBrace)
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

func (p *ParseStreamImpl) maybeLeftParen() (token.ParenToken, bool) {
	return p.maybeSpecificParenToken(token.LeftParen)
}

func (p *ParseStreamImpl) maybeAccessor() bool {
	return p.tokenizer.MaybeAccessor()
}

func (p *ParseStreamImpl) detectTypeIdentifierWithoutScope() bool {
	pos := p.tokenizer.Tell()
	_, typeSymbolErr := p.tokenizer.ParseTypeSymbol()
	accessorFollowing := p.maybeAccessor()
	p.tokenizer.Seek(pos)

	return !accessorFollowing && typeSymbolErr == nil
}

func (p *ParseStreamImpl) detectNormalOrScopedTypeIdentifier() bool {
	pos := p.tokenizer.Tell()

	wasTypeIdentifier := false
	if foundTypeIdentifier, wasTypeSymbol := p.wasTypeIdentifier(); wasTypeSymbol {
		_, typeSymbolErr := parseTypeSymbolWithOptionalModules(p, foundTypeIdentifier)
		if typeSymbolErr != nil {
			wasTypeIdentifier = false
		}
	}

	p.tokenizer.Seek(pos)
	return wasTypeIdentifier
}

func (p *ParseStreamImpl) readVariableIdentifierAssignOrUpdate(indentation int) (*ast.VariableIdentifier, bool, int, parerr.ParseError) {
	ident, identErr := p.readVariableIdentifier()
	if identErr != nil {
		return nil, false, 0, identErr
	}

	_, spaceAfterIdentifierErr := p.eatOneSpaceInternal("space after variableIdentifier assign or update")
	if spaceAfterIdentifierErr != nil {
		return nil, false, 0, parerr.NewExpectedVariableAssignOrRecordUpdate(spaceAfterIdentifierErr)
	}

	wasAssign := p.tokenizer.MaybeRune('=')
	if wasAssign {
		newIndentation, _, spaceAfterAssignErr := p.eatContinuationReturnIndentation(indentation)
		if spaceAfterAssignErr != nil {
			return nil, false, 0, parerr.NewExpectedVariableAssign(spaceAfterAssignErr)
		}
		return ident, true, newIndentation, nil
	}

	eatUpdateErr := p.eatRune('|')
	if eatUpdateErr != nil {
		return nil, false, 0, parerr.NewExpectedRecordUpdate(eatUpdateErr)
	}
	_, spaceAfterUpdateErr := p.eatOneSpaceInternal("space after update |")
	if spaceAfterUpdateErr != nil {
		return nil, false, 0, spaceAfterUpdateErr
	}
	return ident, false, indentation, nil
}

func (p *ParseStreamImpl) readTermToken() (token.Token, parerr.ParseError) {
	pos := p.tokenizer.ParsingPosition()
	t, tErr := p.tokenizer.ReadTermToken()
	if tErr != nil {
		return nil, tErr
	}

	switch t.(type) {
	case token.SpaceToken:
		posLength := p.tokenizer.MakeSourceFileReference(pos)

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

	if _, wasComma := p.maybeComma(); wasComma {
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

func (p *ParseStreamImpl) maybeNewLineContinuationHelper(expectedIndentation int, allowComments tokenize.CommentAllowedType) (bool, token.IndentationReport, parerr.ParseError) {
	pos := p.tokenizer.Tell()
	report, foundPosLengthErr := p.tokenizer.SkipWhitespaceToNextIndentationHelper(allowComments)
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

func (p *ParseStreamImpl) eatNewLineContinuationOrDedent(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	return p.maybeNewLineContinuationHelper(expectedIndentation, tokenize.OwnLine)
}

func (p *ParseStreamImpl) eatOneNewLineContinuationOrDedentAllowComment(expectedIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	return p.maybeNewLineContinuationHelper(expectedIndentation, tokenize.OwnLine)
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
	term, termErr := p.tokenizer.ReadTermTokenOrEndOrSeparator()
	_, isEOF := term.(*tokenize.EndOfFile)
	if isEOF {
		detectedEndOfCallOperator = true
	} else {
		if termErr == nil {
			t := term.Type()
			_, isBinaryOperator := term.(token.OperatorToken)
			if t == token.OperatorUnaryMinus || t == token.OperatorUnaryNot {
				isBinaryOperator = false
			}
			detectedEndOfCallOperator = t == token.Then || t == token.Else || t == token.Of || t == token.RightCurlyBrace || isBinaryOperator || t == token.RightSquareBracket || t == token.RightParen || t == token.Comma || t == token.OperatorPipeRight || t == token.OperatorArrowRight
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

func (p *ParseStreamImpl) addNode(node ast.Node) {
	posLength := node.FetchPositionLength()
	if posLength.Range.Position().Line() == 0 && posLength.Range.Position().Column() == 0 {
		panic("suspicion")
	}
	p.nodes = append(p.nodes, node)
}

// ---------------------------------------------------------------------------------
// MAYBE
// ---------------------------------------------------------------------------------
func (p *ParseStreamImpl) maybeKeywordAlias() (token.Keyword, bool) {
	return p.maybeSpecificKeywordToken(token.Alias)
}

func (p *ParseStreamImpl) maybeKeywordExposing() (token.Keyword, bool) {
	return p.maybeSpecificKeywordToken(token.Exposing)
}

func (p *ParseStreamImpl) maybeKeywordAs() (token.Keyword, bool) {
	return p.maybeSpecificKeywordToken(token.As)
}

func (p *ParseStreamImpl) maybeAssign() bool {
	return p.tokenizer.MaybeAssign()
}

func (p *ParseStreamImpl) maybeEllipsis() bool {
	return p.tokenizer.MaybeString("..")
}

func (p *ParseStreamImpl) maybeNewLine() bool {
	return p.tokenizer.MaybeOneNewLine()
}

func (p *ParseStreamImpl) maybeRightParen() bool {
	return p.tokenizer.MaybeRune(')')
}

func (p *ParseStreamImpl) maybeEmptyParen() bool {
	return p.tokenizer.MaybeString("()")
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
	report, err := p.eatOneSpaceInternal(reason)
	if err != nil {
		_, wasOneSpace := err.(parerr.ExpectedOneSpace)
		if wasOneSpace {
			if report.NewLineCount == 0 && report.SpacesUntilMaybeNewline == 0 {
				return report, err
			}
			p.AddWarning(err)
		}
	}

	return report, nil
}

func (p *ParseStreamImpl) eatOneSpaceInternal(reason string) (token.IndentationReport, parerr.ParseError) {
	pos := p.tokenizer.MakeSourceFileReference(p.tokenizer.ParsingPosition())
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
		spaceErr := parerr.NewExpectedOneSpace(pos)
		return report, spaceErr
	}

	return report, nil
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
			return report.ExactIndentation, report, nil
		}

		if report.NewLineCount == 1 && report.ExactIndentation == indentation+1 {
			return report.ExactIndentation, report, nil
		}
	}

	subErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, indentation, report.CloseIndentation)

	tokenErr := tokenize.NewExpectedIndentationAfterNewLineError(indentation, report.CloseIndentation, subErr)

	return -1, report, parerr.NewExpectedOneSpaceOrExtraIndent(tokenErr)
}

func (p *ParseStreamImpl) eatNewLinesAfterStatement(count int) (token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceAllowCommentsToNextIndentation()
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
		} else if count == -2 {
			if (report.NewLineCount >= 0) && report.IndentationSpaces == 0 {
				return report, nil
			}
		} else {
			if (report.NewLineCount == count) && report.IndentationSpaces == 0 {
				return report, nil
			}
		}

		if report.NewLineCount > 0 && report.IndentationSpaces == 0 {
			err := parerr.NewExpectedNewLineCount(report.PositionLength, count, report.NewLineCount)
			p.AddWarning(err)
			return report, nil
		}
	}

	return report, parerr.NewExpectedNewLineCount(report.PositionLength, count, report.NewLineCount)
}

func (p *ParseStreamImpl) eatCommaSeparatorOrTermination(expectedIndentation int, allowComments tokenize.CommentAllowedType) (bool, token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentationHelper(allowComments)
	if err != nil {
		return false, report, err
	}
	wasComma := p.tokenizer.MaybeRune(',')
	if wasComma {
		wasNextLineOrNoSpaceBeforeComma := tokenize.LegalSameIndentationOrNoSpace(report, expectedIndentation, !p.disableEnforceStyle)
		if !wasNextLineOrNoSpaceBeforeComma {
			return false, report, tokenize.NewUnexpectedIndentationError(report.PositionLength, expectedIndentation, report.CloseIndentation)
		}
		spaceAfterCommaReport, eatErr := p.eatOneSpaceInternal("space after comma")
		if eatErr != nil {
			return true, spaceAfterCommaReport, eatErr
		}

		if p.disableEnforceStyle {
			return true, report, nil
		}

		wasNextLineOrNoSpace := tokenize.LegalSameIndentationOrNoSpace(report, report.ExactIndentation, !p.disableEnforceStyle)
		if wasNextLineOrNoSpace {
			return true, report, nil
		}

		return false, report, parerr.NewExtraSpacing(p.sourceFileReference())
	}

	// It was a termination, that should be handled, but we can still check that the space is correct
	if p.disableEnforceStyle {
		if tokenize.LegalSameIndentationOrOptionalOneSpace(report, expectedIndentation, !p.disableEnforceStyle) {
			return false, report, nil
		}
		fmt.Printf("%d\n", report.SpacesUntilMaybeNewline)
	} else {
		wasNextLineOrOneSpace := tokenize.LegalSameIndentationOrOptionalOneSpace(report, expectedIndentation, !p.disableEnforceStyle)
		if wasNextLineOrOneSpace {
			return false, report, nil
		}
	}

	tokenErr := tokenize.NewUnexpectedIndentationError(report.PositionLength, expectedIndentation, report.CloseIndentation)
	return false, report, parerr.NewExpectedOneSpaceOrExtraIndentCommaSeparator(tokenErr)
}

func (p *ParseStreamImpl) eatOneOrTwoNewLineContinuationOrDedent(currentIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceAllowCommentsToNextIndentation()
	if err != nil {
		return false, report, err
	}

	if p.disableEnforceStyle {
		wasContinuation := report.CloseIndentation >= currentIndentation
		return wasContinuation, report, nil
	}

	if (report.NewLineCount == 1 || report.NewLineCount == 2) && report.ExactIndentation == currentIndentation {
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

func (p *ParseStreamImpl) eatOneSpaceOrDetectEndOfFunctionCallArguments(currentIndentation int) (bool, token.IndentationReport, parerr.ParseError) {
	saveInfo := p.tokenizer.Tell()
	report, err := p.tokenizer.SkipWhitespaceToNextIndentation()
	if err != nil {
		return false, report, err
	}

	if p.disableEnforceStyle {
		wasLegalContinuation := tokenize.LegalContinuationSpaceIndentation(report, currentIndentation, !p.disableEnforceStyle)
		if !wasLegalContinuation { // report.NewLineCount > 0 { // Arguments end even if they are indented
			p.tokenizer.Seek(saveInfo)
			return true, report, nil
		}
		wasEnd := p.detectEndOfCallOperator()
		if wasEnd {
			p.tokenizer.Seek(saveInfo)
		}
		return wasEnd, report, nil
	} else {
		isLegalSpacing := (report.NewLineCount == 1 && report.ExactIndentation == currentIndentation+1) || (report.NewLineCount == 0 && (report.SpacesUntilMaybeNewline == 1 || report.SpacesUntilMaybeNewline == 0))
		if isLegalSpacing {
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

func (p *ParseStreamImpl) eatNewLineAndExactIndentHelper(indentation int, allowComments tokenize.CommentAllowedType) (token.IndentationReport, parerr.ParseError) {
	report, err := p.tokenizer.SkipWhitespaceToNextIndentationHelper(allowComments)
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

func (p *ParseStreamImpl) eatNewLineAndExactIndentAllowComments(indentation int) (token.IndentationReport, parerr.ParseError) {
	return p.eatNewLineAndExactIndentHelper(indentation, tokenize.OwnLine)
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

func (p *ParseStreamImpl) eatNewLineContinuationAllowComment(indentation int) (token.IndentationReport, parerr.ParseError) {
	return p.eatNewLineAndExactIndentAllowComments(indentation + 1)
}

func (p *ParseStreamImpl) eatOperatorUpdate() parerr.ParseError {
	eatUpdateErr := p.eatRune('|')
	return eatUpdateErr
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

func (p *ParseStreamImpl) readRightParen() (token.ParenToken, parerr.ParseError) {
	return p.readSpecificParenToken(token.RightParen)
}

func (p *ParseStreamImpl) eatLeftParen() parerr.ParseError {
	return p.eatRune('(')
}

func (p *ParseStreamImpl) readRightBracket() (token.ParenToken, parerr.ParseError) {
	return p.readSpecificParenToken(token.RightSquareBracket)
}

func (p *ParseStreamImpl) readLeftAngleBracket() (token.ParenToken, parerr.ParseError) {
	return p.readSpecificParenToken(token.LeftAngleBracket)
}

func (p *ParseStreamImpl) readRightAngleBracket() (token.ParenToken, parerr.ParseError) {
	return p.readSpecificParenToken(token.RightAngleBracket)
}

func (p *ParseStreamImpl) readRightCurly() (token.ParenToken, parerr.ParseError) {
	return p.readSpecificParenToken(token.RightCurlyBrace)
}

func (p *ParseStreamImpl) readRightArrayBracket() (token.ParenToken, parerr.ParseError) {
	return p.readSpecificParenToken(token.RightArrayBracket)
}

func (p *ParseStreamImpl) readOf() (token.Keyword, parerr.ParseError) {
	return p.readSpecificKeywordToken(token.Of)
}

func (p *ParseStreamImpl) readThen() (token.Keyword, parerr.ParseError) {
	return p.readSpecificKeywordToken(token.Then)
}

func (p *ParseStreamImpl) readElse() (token.Keyword, parerr.ParseError) {
	return p.readSpecificKeywordToken(token.Else)
}

func (p *ParseStreamImpl) readGuardPipe() (token.GuardToken, parerr.ParseError) {
	someToken, termErr := p.tokenizer.ReadTermToken()
	if termErr != nil {
		return token.GuardToken{}, termErr
	}

	guardToken, wasGuard := someToken.(token.GuardToken)
	if !wasGuard {
		return token.GuardToken{}, parerr.NewInternalError(p.positionLength(), fmt.Errorf("must have guard | token here"))
	}

	return guardToken, nil
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
	if rErr != nil {
		p.AddError(rErr)
	}
	p.descent--
	return r, rErr
}

func (p *ParseStreamImpl) parseExpressionNormal(startIndentation int) (ast.Expression, parerr.ParseError) {
	p.descent++
	r, rErr := p.parser.parseExpressionNormal(startIndentation)
	p.descent--
	return r, rErr
}

func (p *ParseStreamImpl) parseExpressionNormalWithComment(startIndentation int, comment token.Comment) (ast.Expression, parerr.ParseError) {
	p.descent++
	r, rErr := p.parser.parseExpressionNormalWithComment(startIndentation, comment)
	p.descent--
	return r, rErr
}

func (p *ParseStreamImpl) parseTerm(startIndentation int) (ast.Expression, parerr.ParseError) {
	expr, err := p.parser.parseTerm(startIndentation)
	if err != nil {
		p.AddError(err)
	}

	return expr, err
}

func (p *ParseStreamImpl) parseLiteral(startIndentation int) (ast.Literal, parerr.ParseError) {
	return p.parser.parseTerm(startIndentation)
}
