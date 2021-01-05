/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"
	"strings"

	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/token"
)

type TokenizerError struct {
	err      error
	position token.PositionLength
}

func (f TokenizerError) FetchPositionLength() token.PositionLength {
	return f.position
}

func (f TokenizerError) Error() string {
	return fmt.Sprintf("%v at %v", f.err, f.position)
}

// Tokenizer :
type Tokenizer struct {
	r                     *runestream.RuneReader
	oldPosition           token.PositionToken
	position              token.PositionToken
	lastTokenWasDelimiter bool
	lastReport            token.IndentationReport
	enforceStyleGuide     bool
}

func verifyOctets(octets []byte) TokenError {
	pos := token.NewPositionTopLeft()
	for _, octet := range octets {
		r := rune(octet)
		if r != 0 && r != 10 && (r < 32 || r > 126) {
			posLength := token.NewPositionLength(pos, 1, -1)
			return NewUnexpectedEatTokenError(posLength, ' ', r)
		}
		pos = nextPosition(pos, r)
	}
	return nil
}

// NewTokenizer :
func NewTokenizer(r *runestream.RuneReader, exactWhitespace bool) (*Tokenizer, TokenError) {
	t := &Tokenizer{r: r,
		position:              token.NewPositionToken(token.NewPositionTopLeft(), 0),
		lastTokenWasDelimiter: true,
		enforceStyleGuide:     exactWhitespace}
	verifyErr := verifyOctets(r.Octets())
	if verifyErr != nil {
		return t, verifyErr
	}
	return t, nil
}

func (t *Tokenizer) MakePositionLength(pos token.PositionToken) token.PositionLength {
	return token.NewPositionLength(pos.Position(), t.position.Position().Column()-pos.Position().Column(), pos.Indentation())
}

const SpacesForIndentation = 4

func (t *Tokenizer) ParsingPosition() token.PositionToken {
	return t.position
}

func (t *Tokenizer) RelativeFilename() string {
	return t.r.RelativeFilename()
}

func (t *Tokenizer) internalDebugInfo(rowCount int) string {
	pos := t.ParsingPosition().Position()
	focusStartColumn := pos.Column()
	lines := t.ExtractStrings(pos.Line(), rowCount)
	if len(lines) == 0 {
		return fmt.Sprintf(">> strange, no lines (%v)", pos)
	}
	line := lines[0]
	if len(line) == 0 {
		return ">> strange, line is empty"
	}

	if focusStartColumn > len(line) {
		fmt.Printf("!!!!!!!!!!!!!! cant be right focusCol:%v lineLength:%v\n", focusStartColumn, len(line))
	}

	startColumn := 0
	endColumn := len(line) - 1

	if rowCount == 1 {
		const lookAround = 3
		startColumn = focusStartColumn - lookAround
		if startColumn < 0 {
			startColumn = 0
		}
		endColumn = focusStartColumn + lookAround
		if endColumn >= len(line) {
			endColumn = len(line) - 1
		}
	}

	extract := line[startColumn:endColumn+1] + "↵"
	underlineIndex := focusStartColumn - startColumn
	if underlineIndex > len(extract) {
		fmt.Printf("!!!!!!!!!!!!!! we have a problem startCol:%v endCol:%v focusCol:%v\n", startColumn, endColumn, focusStartColumn)
	}
	prefix := strings.Repeat(" ", underlineIndex)
	showLine := fmt.Sprintf("%v\n%v^", extract, prefix)
	if len(lines) > 1 {
		showLine += "\n" + lines[1]
	}
	return showLine
}

func (t *Tokenizer) DebugInfo() string {
	return t.internalDebugInfo(2)
}

func (t *Tokenizer) DebugInfoWithComment(s string) string {
	debug := t.DebugInfo()
	return fmt.Sprintf("---- %v:\n%v\n--#--\n", s, debug)
}

func (t *Tokenizer) DebugInfoLinesWithComment(s string, rowCount int) string {
	debug := t.internalDebugInfo(rowCount)
	return fmt.Sprintf("---- %v:\n%v\n--#--\n", s, debug)
}

func (t *Tokenizer) DebugPrint(s string) {
	fmt.Print(t.DebugInfoWithComment(s))
}

func nextPosition(pos token.Position, ch rune) token.Position {
	if ch == '\n' {
		pos = pos.NextLine()
		pos = pos.FirstColumn()
	} else {
		pos = pos.NextColumn()
	}
	return pos
}

func (t *Tokenizer) reversePositionHelper(pos token.Position, ch rune) token.Position {
	if ch == '\n' {
		column, detectedIndentationSpaces := t.r.DetectCurrentColumn()
		pos = token.MakePosition(pos.Line()-1, column)
		t.lastReport.IndentationSpaces = detectedIndentationSpaces
		t.lastReport.CloseIndentation = detectedIndentationSpaces / SpacesForIndentation
		if (detectedIndentationSpaces % SpacesForIndentation) == 0 {
			t.lastReport.ExactIndentation = t.lastReport.CloseIndentation
		} else {
			t.lastReport.ExactIndentation = -1
		}
	} else {
		pos = pos.PreviousColumn()
	}
	return pos
}

func (t *Tokenizer) reversePosition(ch rune) {
	pos := t.position.Position()
	pos = t.reversePositionHelper(pos, ch)
	t.position = token.NewPositionToken(pos, (pos.Column())/SpacesForIndentation)
}

func (t *Tokenizer) updatePosition(ch rune) {
	pos := t.position.Position()
	pos = nextPosition(pos, ch)
	t.position = token.NewPositionToken(pos, (pos.Column())/SpacesForIndentation)
}

func LegalContinuationSpaceIndentation(report token.IndentationReport, indentation int, enforceStyle bool) bool {
	if enforceStyle {
		if report.ExactIndentation > indentation+1 {
			return true
		}
	} else {
		if report.CloseIndentation > indentation {
			return true
		}
	}

	if report.NewLineCount == 0 && report.SpacesUntilMaybeNewline >= 0 {
		return true
	}

	return false
}

func LegalContinuationSpace(report token.IndentationReport, enforceStyle bool) bool {
	if enforceStyle {
		if report.SpacesUntilMaybeNewline == 1 {
			return true
		}
		if report.ExactIndentation == report.PreviousExactIndentation+1 {
			return true
		}
		return false
	} else {
		if report.NewLineCount == 0 && report.SpacesUntilMaybeNewline == 1 {
			return true
		}

		if report.CloseIndentation == report.PreviousCloseIndentation+1 {
			return true
		}
		return false
	}

	return false
}

func LegalSameIndentationOrNoSpace(report token.IndentationReport, indentation int, enforceStyle bool) bool {
	if report.NewLineCount == 0 && report.SpacesUntilMaybeNewline == 0 {
		return true
	}

	if enforceStyle {
		if report.NewLineCount == 1 && report.ExactIndentation == indentation {
			return true
		}
	} else {
		if report.NewLineCount > 0 && report.CloseIndentation == indentation {
			return true
		}
	}

	return false
}

func LegalSameIndentationOrOptionalOneSpace(report token.IndentationReport, indentation int, enforceStyle bool) bool {
	if report.NewLineCount == 0 && report.SpacesUntilMaybeNewline == 1 {
		return true
	}

	if enforceStyle {
		if report.NewLineCount == 1 && report.ExactIndentation == indentation {
			return true
		}
	} else {
		if report.NewLineCount == 0 && (report.SpacesUntilMaybeNewline == 1 || report.SpacesUntilMaybeNewline == 0) {
			return true
		}
		if report.NewLineCount > 0 && report.CloseIndentation == indentation {
			return true
		}
	}

	return false
}

func IsDedent(report token.IndentationReport, enforceStyle bool) bool {
	if enforceStyle {
		if report.ExactIndentation < report.PreviousExactIndentation {
			return true
		}
	} else {
		if report.CloseIndentation > report.PreviousCloseIndentation {
			return true
		}
	}

	return false
}

func (t *Tokenizer) EatOneSpace() (token.IndentationReport, TokenError) {
	report, err := t.SkipWhitespaceToNextIndentation()
	if err != nil {
		return report, err
	}

	if t.enforceStyleGuide {
		if report.SpacesUntilMaybeNewline != 1 {
			return report, NewExpectedOneSpaceError(err)
		}
	}

	if LegalContinuationSpace(report, t.enforceStyleGuide) {
		return report, nil
	}

	return report, nil
}

func (t *Tokenizer) SkipWhitespaceToNextIndentation() (token.IndentationReport, TokenError) {
	const disallowComments = false

	return t.SkipWhitespaceToNextIndentationHelper(disallowComments)
}

func (t *Tokenizer) SkipWhitespaceAllowCommentsToNextIndentation() (token.IndentationReport, TokenError) {
	const allowComments = true

	return t.SkipWhitespaceToNextIndentationHelper(allowComments)
}

func (t *Tokenizer) SkipWhitespaceToNextIndentationHelper(allowComments bool) (token.IndentationReport, TokenError) {
	var comments []token.CommentToken

	detectedIndentationSpaces := 0 // t.lastReport.IndentationSpaces
	newLineCount := 0
	startPos := t.position
	spacesUntilMaybeNewline := 0
	hasTrailingSpaces := false
	closeIndentation := t.lastReport.CloseIndentation
	exactIndentation := t.lastReport.ExactIndentation

	for {
		ch := t.nextRune()
		if ch == '\n' {
			newLineCount++
			if detectedIndentationSpaces > 0 {
				hasTrailingSpaces = true
			}
			if spacesUntilMaybeNewline > 0 || detectedIndentationSpaces > 0 {
				if t.enforceStyleGuide {
					trailingPosLength := token.NewPositionLength(startPos.Position(), 1, startPos.Indentation())
					return token.IndentationReport{}, NewTrailingSpaceError(trailingPosLength)
				}
			}
			detectedIndentationSpaces = 0
			spacesUntilMaybeNewline = 0
		} else if isIndentation(ch) {
			if newLineCount > 0 {
				detectedIndentationSpaces++
			}
			if newLineCount == 0 {
				spacesUntilMaybeNewline++
			}
		} else {
			if ch == 0 { // treat end of file as return
				newLineCount++
				if detectedIndentationSpaces > 0 {
					hasTrailingSpaces = true
				}
				detectedIndentationSpaces = 0
				spacesUntilMaybeNewline = 0
			}

			if newLineCount > 0 {
				exactIndentation = -1
				closeIndentation = detectedIndentationSpaces / SpacesForIndentation
				if (detectedIndentationSpaces % SpacesForIndentation) == 0 {
					exactIndentation = closeIndentation
				}
			}

			if allowComments {
				comment, found, err := t.checkComment(ch, t.position)
				if err != nil {
					return token.IndentationReport{}, err
				}
				if found {
					comments = append(comments, comment)
					detectedIndentationSpaces = 0 // t.lastReport.IndentationSpaces
					// newLineCount = 0  // keep number of lines
					startPos = t.position
					spacesUntilMaybeNewline = 0
					hasTrailingSpaces = false
					closeIndentation = t.lastReport.CloseIndentation
					exactIndentation = t.lastReport.ExactIndentation
					continue
				}
			}
			endOfFile := ch == 0
			t.unreadRune()
			if newLineCount > 0 {
				spacesUntilMaybeNewline = -1
			}

			newReport := token.IndentationReport{IndentationSpaces: detectedIndentationSpaces,
				CloseIndentation:          closeIndentation,
				ExactIndentation:          exactIndentation,
				Comments:                  token.MakeCommentBlock(comments),
				NewLineCount:              newLineCount,
				StartPos:                  startPos,
				PositionLength:            token.NewPositionLength(startPos.Position(), 1, startPos.Indentation()),
				TrailingSpacesFound:       hasTrailingSpaces,
				SpacesUntilMaybeNewline:   spacesUntilMaybeNewline,
				PreviousCloseIndentation:  t.lastReport.CloseIndentation,
				PreviousExactIndentation:  t.lastReport.ExactIndentation,
				PreviousIndentationSpaces: t.lastReport.IndentationSpaces,
				EndOfFile:                 endOfFile}

			t.lastReport = newReport
			return newReport, nil
		}
	}
}

func (t *Tokenizer) Tell() int {
	return t.r.Tell()
}

func (t *Tokenizer) Seek(pos int) {
	deltaPos := pos - t.r.Tell()
	if deltaPos > 0 {
		panic("this cant be good")
	}
	count := -deltaPos

	for i := 0; i < count; i++ {
		t.unreadRune()
	}
}

func (t *Tokenizer) ExtractStrings(startRow int, rowCount int) []string {
	octets := t.r.Octets()
	row := 0
	stringRow := ""
	var foundRows []string
	for _, octet := range octets {
		r := rune(octet)
		if r == '\n' || r == 0 {
			if row >= startRow && row < startRow+rowCount {
				foundRows = append(foundRows, stringRow)
			}
			row++
			stringRow = ""
		} else {
			stringRow += string(r)
		}
	}

	return foundRows
}

func AddInvisibleCharacters(s string) string {
	stringRow := ""
	for _, r := range s {
		if r == '\n' || r == 0 {
			if r == 0 {
				stringRow += "␃"
			} else if r == 10 {
				stringRow += "␤"
			}

			stringRow = ""
		} else {
			if (r < 32) || (r > 126) {
				stringRow += "⛔"
			}
			stringRow += string(r)
		}
	}
	return stringRow
}

func (t *Tokenizer) nextRune() rune {
	t.oldPosition = t.position
	ch := t.r.Read()
	t.updatePosition(ch)
	return ch
}

func (t *Tokenizer) unreadRune() {
	ch := t.r.Unread()
	t.reversePosition(ch)
}

func (t *Tokenizer) ParseString(startStringRune rune, startPosition token.PositionToken) (token.StringToken, error) {
	var a string
	raw := string(startStringRune)
	for {
		ch := t.nextRune()
		raw += string(ch)
		if ch == startStringRune {
			break
		}
		if ch == 0 {
			return token.StringToken{}, fmt.Errorf("unexpected end while finding end of string")
		}
		a += string(ch)
	}
	posLen := t.MakePositionLength(startPosition)
	return token.NewStringToken(raw, a, posLen), nil
}

func (t *Tokenizer) ReadStringUntilEndOfLine() string {
	s := ""
	for {
		r := t.nextRune()
		if isNewLine(r) || r == 0 {
			t.unreadRune()
			break
		}
		s += string(r)
	}
	return s
}

func (t *Tokenizer) ReadMultilineComment(positionToken token.PositionToken) (token.MultiLineCommentToken, TokenError) {
	s, documentationComment, err := t.ReadStringUntilEndOfMultilineComment()
	if err != nil {
		return token.MultiLineCommentToken{}, err
	}
	return token.NewMultiLineCommentToken("{-"+s, s, documentationComment, t.MakePositionLength(positionToken)), nil
}

func (t *Tokenizer) ReadSingleLineComment(positionToken token.PositionToken) token.MultiLineCommentToken {
	firstCh := t.nextRune()
	documentationComment := false
	if firstCh == '|' {
		documentationComment = true
	} else {
		t.unreadRune()
	}
	s := t.ReadStringUntilEndOfLine()
	return token.NewMultiLineCommentToken("--"+s, s, documentationComment, t.MakePositionLength(positionToken))
}

func (t *Tokenizer) ReadStringUntilEndOfMultilineComment() (string, bool, TokenError) {
	firstCh := t.nextRune()
	if firstCh == 0 {
		return "", false, NewInternalError(fmt.Errorf("unexpected end of file"))
	}
	documentationComment := false
	if firstCh == '|' {
		documentationComment = true
	} else {
		t.unreadRune()
	}
	s := ""
	for {
		r := t.nextRune()
		if r == '-' {
			if t.nextRune() == '}' {
				break
			} else {
				t.unreadRune()
			}
		} else if r == 0 {
			return "", false, NewInternalError(fmt.Errorf("unexpected end of file"))
		}

		s += string(r)
	}

	return s, documentationComment, nil
}

func (t *Tokenizer) ParseStartingKeyword() (token.Token, TokenError) {
	r := t.nextRune()
	if r == '_' {
		nextRune := t.nextRune()
		if nextRune == '_' {
			return t.ParseSpecialKeyword(t.position)
		}
		return nil, NewInternalError(fmt.Errorf("unknown starting keyword"))
	}

	t.unreadRune()
	return t.ParseVariableSymbol()
}

func (t *Tokenizer) internalGuessNext() (token.Token, error) {
	posToken := t.position
	r := t.nextRune()
	singleCharLength := t.MakePositionLength(posToken)
	if r == 0 {
		return &EndOfFile{}, nil
	}
	if isNewLine(r) {
		return token.NewLineDelimiter(t.MakePositionLength(posToken)), nil
	}
	t.lastTokenWasDelimiter = false
	if isLetter(r) {
		t.unreadRune()
		return t.parseAnySymbol(posToken)
	} else if isDigit(r) {
		t.unreadRune()
		return t.ParseNumber("")
	} else if isStartString(r) {
		return t.ParseString(r, posToken)
	} else if isOperator(r) {
		t.unreadRune()
		return t.ParseOperator()
	} else if r == ')' {
		return token.NewParenToken(string(r), token.RightParen, singleCharLength, " )R "), nil
	} else if r == '(' {
		return token.NewParenToken(string(r), token.LeftParen, singleCharLength, " L( "), nil
	} else if r == '{' {
		nch := t.nextRune()
		if nch == '-' {
			return t.ReadMultilineComment(posToken)
		}
		t.unreadRune()

		return token.NewParenToken(string(r), token.LeftCurlyBrace, singleCharLength, " { "), nil
	} else if r == '}' {
		return token.NewParenToken(string(r), token.RightCurlyBrace, singleCharLength, " } "), nil
	} else if r == '[' {
		return token.NewParenToken(string(r), token.LeftBracket, singleCharLength, " [ "), nil
	} else if r == ']' {
		return token.NewParenToken(string(r), token.RightBracket, singleCharLength, " ] "), nil
	} else if r == ',' {
		return token.NewParenToken(string(r), token.Comma, singleCharLength, ","), nil
	} else if r == '\\' {
		return token.NewLambdaToken(singleCharLength, " lambda "), nil
	} else if r == '_' {
		nextRune := t.nextRune()
		if nextRune == '_' {
			return t.ParseSpecialKeyword(t.position)
		}
		t.unreadRune()
		return token.NewTypeSymbolToken("_", singleCharLength, -1), nil
	} else if r == 0 {
		return nil, nil
	} else if r == ' ' {
		return token.NewSpaceToken(t.MakePositionLength(t.position), r), nil
	}
	return nil, fmt.Errorf("unknown rune '%c' %v", r, r)
}

func (t *Tokenizer) GuessNext() (token.Token, TokenError) {
	startPos := t.position
	token, err := t.internalGuessNext()
	if err != nil {
		return nil, TokenizerError{err: err, position: t.MakePositionLength(startPos)}
	}
	return token, nil
}

func (e *EndOfFile) String() string {
	return "EOF"
}
