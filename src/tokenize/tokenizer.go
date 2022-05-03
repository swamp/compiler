/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"
	"os"
	"strings"

	"github.com/swamp/compiler/src/runestream"
	"github.com/swamp/compiler/src/token"
)

type TokenizerError struct {
	err      error
	position token.SourceFileReference
}

func (f TokenizerError) FetchPositionLength() token.SourceFileReference {
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
	document              *token.SourceFileDocument
	lastTokenWasDelimiter bool
	lastReport            token.IndentationReport
	enforceStyleGuide     bool
}

func verifyOctets(octets []byte, relativeFilename string) TokenError {
	pos := token.NewPositionTopLeft()
	fileDocument := token.MakeSourceFileDocumentFromLocalPath(relativeFilename)
	if len(relativeFilename) == 0 {
		panic("must have relative filename")
	}
	for _, octet := range octets {
		r := rune(octet)
		if r != 0 && r != 10 && (r < 32 || r > 126) {
			posLength := token.SourceFileReference{
				Range:    token.NewPositionLength(pos, 1),
				Document: nil,
			}
			return NewUnexpectedEatTokenError(posLength, ' ', r)
		}
		if r == '\n' || r == 0 {
			const maxColumn = 120
			const recommendedMaxColumn = 110

			sourceFileReference := token.SourceFileReference{
				Range:    token.MakeRange(pos, pos),
				Document: fileDocument,
			}
			if pos.Column() > maxColumn {
				fmt.Fprintf(os.Stderr, "%v: Warning: line is too long (%v of max %v).\n", sourceFileReference.ToStandardReferenceString(),
					pos.Column(), maxColumn)
			} else if pos.Column() > recommendedMaxColumn {
				fmt.Fprintf(os.Stderr, "%v: Note: exceeds recommended line length (%v of recommended %v).\n", sourceFileReference.ToStandardReferenceString(),
					pos.Column(), recommendedMaxColumn)
				/*
					fmt.Fprintf(os.Stderr, "%s Warning: %v\n", warning.FetchPositionLength().ToStandardReferenceString(), warning.Warning())

				*/
			}
		}
		pos = nextPosition(pos, r)
	}
	return nil
}

// NewTokenizerInternal :
func NewTokenizerInternal(r *runestream.RuneReader, exactWhitespace bool) (*Tokenizer, TokenError) {
	t := &Tokenizer{
		r:                     r,
		document:              token.MakeSourceFileDocumentFromLocalPath(r.RelativeFilename()),
		position:              token.NewPositionToken(token.NewPositionTopLeft(), 0),
		lastTokenWasDelimiter: true,
		enforceStyleGuide:     exactWhitespace,
	}

	return t, nil
}

// NewTokenizerInternalWithStartPosition :
func NewTokenizerInternalWithStartPosition(r *runestream.RuneReader, position token.Position, exactWhitespace bool) (*Tokenizer, TokenError) {
	t := &Tokenizer{
		r:                     r,
		document:              token.MakeSourceFileDocumentFromLocalPath(r.RelativeFilename()),
		position:              token.NewPositionToken(position, 0),
		lastTokenWasDelimiter: true,
		enforceStyleGuide:     exactWhitespace,
	}

	return t, nil
}

// NewTokenizer :
func NewTokenizer(r *runestream.RuneReader, exactWhitespace bool) (*Tokenizer, TokenError) {
	verifyErr := verifyOctets(r.Octets(), r.RelativeFilename())
	if verifyErr != nil {
		return nil, verifyErr
	}

	return NewTokenizerInternal(r, exactWhitespace)
}

func (t *Tokenizer) SourceFile() *token.SourceFileURI {
	return token.MakeSourceFileURI(t.r.RelativeFilename())
}

func (t *Tokenizer) Document() *token.SourceFileDocument {
	return t.document
}

func (t *Tokenizer) MakeRangeMinusOne(pos token.PositionToken) token.Range {
	endPos := t.position.Position()
	if endPos.Column() == 0 {
		//	panic("can not go back")
	} else {
		endPos = endPos.PreviousColumn()
	}
	return token.NewPositionLengthFromEndPosition(pos.Position(), endPos)
}

func (t *Tokenizer) MakeSourceFileReference(pos token.PositionToken) token.SourceFileReference {
	tokenRange := t.MakeRangeMinusOne(pos)
	return token.SourceFileReference{
		Range:    tokenRange,
		Document: t.document,
	}
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
	fmt.Fprintf(os.Stderr, t.DebugInfoWithComment(s))
}

func nextPosition(pos token.Position, ch rune) token.Position {
	if ch == '\n' {
		pos = pos.NewLine()
	} else {
		pos = pos.NextColumn()
	}
	return pos
}

func (t *Tokenizer) reversePositionHelper(pos token.Position, ch rune) token.Position {
	if ch == '\n' {
		column, detectedIndentationSpaces := t.r.DetectCurrentLineLength()
		pos = token.MakePosition(pos.Line()-1, column, pos.OctetOffset()+column)
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

func LegalOneSpaceOrSameIndentation(report token.IndentationReport, indentation int, enforceStyle bool) (bool, TokenError) {
	if enforceStyle {
		if report.NewLineCount == 1 && report.ExactIndentation == indentation {
			return true, nil
		}
		if report.SpacesUntilMaybeNewline == 1 && report.NewLineCount == 0 {
			return false, nil
		}
	} else {
		if report.CloseIndentation > indentation {
			return true, nil
		}
		if report.NewLineCount == 0 && report.SpacesUntilMaybeNewline > 0 {
			return false, nil
		}
	}

	return false, NewUnexpectedIndentationError(report.PositionLength, 0, 0)
}

func LegalOneSpaceOrNewLineIndentation(report token.IndentationReport, indentation int, enforceStyle bool) (bool, TokenError) {
	if enforceStyle {
		if report.NewLineCount == 1 && report.ExactIndentation == indentation+1 {
			return true, nil
		}
		if report.SpacesUntilMaybeNewline == 1 && report.NewLineCount == 0 {
			return false, nil
		}
	} else {
		if report.CloseIndentation > indentation {
			return true, nil
		}
		if report.NewLineCount == 0 && report.SpacesUntilMaybeNewline > 0 {
			return false, nil
		}
	}

	return false, NewUnexpectedIndentationError(report.PositionLength, 0, 0)
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
	if report.NewLineCount == 0 && (report.SpacesUntilMaybeNewline == 1 || report.SpacesUntilMaybeNewline == 0) {
		return true
	}

	if enforceStyle {
		if report.NewLineCount == 1 && (report.ExactIndentation == indentation || report.ExactIndentation == indentation+1) {
			return true
		} else {
			return false
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
			return report, NewExpectedOneSpaceError(NewInternalError(fmt.Errorf("need one space")))
		}
	}

	if LegalContinuationSpace(report, t.enforceStyleGuide) {
		return report, nil
	}

	return report, nil
}

func (t *Tokenizer) SkipWhitespaceToNextIndentation() (token.IndentationReport, TokenError) {
	return t.SkipWhitespaceToNextIndentationHelper(NotAllowedAtAll)
}

func (t *Tokenizer) SkipWhitespaceAllowCommentsToNextIndentation() (token.IndentationReport, TokenError) {
	return t.SkipWhitespaceToNextIndentationHelper(SameLine)
}

type CommentAllowedType int

const (
	SameLine CommentAllowedType = iota
	OwnLine
	NotAllowedAtAll
)

func (t *Tokenizer) SkipWhitespaceToNextIndentationHelper(allowComments CommentAllowedType) (token.IndentationReport, TokenError) {
	var comments []token.Comment

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
					trailingPosLength := token.SourceFileReference{
						Range:    token.NewPositionLength(startPos.Position(), 1),
						Document: t.document,
					}
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

			if allowComments != NotAllowedAtAll {
				previous := token.NewPositionToken(t.position.Position().PreviousColumn(), t.position.Indentation())

				comment, found, err := t.checkComment(ch, previous)
				if err != nil {
					return token.IndentationReport{}, err
				}
				if found {
					if allowComments == OwnLine {
						if newLineCount == 0 {
							trailingPosLength := token.SourceFileReference{
								Range:    token.NewPositionLength(startPos.Position(), 1),
								Document: t.document,
							}
							return token.IndentationReport{}, NewCommentNotAllowedHereError(trailingPosLength, fmt.Errorf("not allowed to have comment on same line"))
						}
					}
					comments = append(comments, comment)
					detectedIndentationSpaces = 0 // t.lastReport.IndentationSpaces
					startPos = t.position
					spacesUntilMaybeNewline = 0
					hasTrailingSpaces = false
					closeIndentation = t.lastReport.CloseIndentation
					exactIndentation = t.lastReport.ExactIndentation
					if newLineCount > 0 {
						newLineCount-- // discard newline after a comment that was preceded with a new line
					}

					continue
				}
			}
			endOfFile := ch == 0
			t.unreadRune()
			if newLineCount > 0 {
				spacesUntilMaybeNewline = -1
			}

			newReport := token.IndentationReport{
				IndentationSpaces:         detectedIndentationSpaces,
				CloseIndentation:          closeIndentation,
				ExactIndentation:          exactIndentation,
				Comments:                  token.MakeCommentBlock(comments),
				NewLineCount:              newLineCount,
				StartPos:                  startPos,
				PositionLength:            t.MakeSourceFileReference(startPos),
				TrailingSpacesFound:       hasTrailingSpaces,
				SpacesUntilMaybeNewline:   spacesUntilMaybeNewline,
				PreviousCloseIndentation:  t.lastReport.CloseIndentation,
				PreviousExactIndentation:  t.lastReport.ExactIndentation,
				PreviousIndentationSpaces: t.lastReport.IndentationSpaces,
				EndOfFile:                 endOfFile,
			}

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
	if t.r == nil {
		return []string{"rune reader is nil"}
	}
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

func (t *Tokenizer) skipSpaces() {
	for {
		r := t.nextRune()
		if !isIndentation(r) || r == 0 {
			t.unreadRune()
			break
		}
	}
}

func (t *Tokenizer) ParseCharacter(startPosition token.PositionToken) (token.CharacterToken, TokenError) {
	ch := t.nextRune()
	if ch == 0 {
		return token.CharacterToken{}, NewEncounteredEOF()
	}

	if ch == '\\' {
		// Escape character
		ch = t.nextRune()
		if ch == 0 {
			return token.CharacterToken{}, NewEncounteredEOF()
		}
	}
	terminator := t.nextRune()
	if terminator != '\'' {
		return token.CharacterToken{}, NewUnexpectedEatTokenError(t.MakeSourceFileReference(startPosition), '\'', terminator)
	}
	posLen := t.MakeSourceFileReference(startPosition)
	return token.NewCharacterToken("'"+string(ch)+"'", ch, posLen), nil
}

func (t *Tokenizer) ParseString(startStringRune rune, startPosition token.PositionToken) (token.StringToken, TokenError) {
	var a string
	raw := string(startStringRune)
	for {
		ch := t.nextRune()
		raw += string(ch)
		if ch == startStringRune {
			break
		}
		if ch == 0 {
			return token.StringToken{}, NewEncounteredEOF()
		}

		if ch == '\\' {
			next := t.nextRune()
			if next == '\n' || next == '\r' {
				t.skipSpaces()
				continue
			} else {
				t.unreadRune()
			}
		}

		if ch == '\n' || ch == '\r' {
			// we ignore new line (LF) in normal string literals. See verbatim strings (triple quote strings) for other behavior.
			continue
		}

		a += string(ch)
	}
	posLen := t.MakeSourceFileReference(startPosition)
	return token.NewStringToken(raw, a, posLen), nil
}

func (t *Tokenizer) isTriple(ch rune, startStringRune rune) (bool, error) {
	if ch == startStringRune {
		if t.nextRune() == startStringRune {
			if t.nextRune() == startStringRune {
				return true, nil
			} else {
				t.unreadRune()
				t.unreadRune()
			}
		} else {
			t.unreadRune()
		}
	} else if ch == 0 {
		return false, fmt.Errorf("unexpected end while finding end of triple string")
	}

	return false, nil
}

func (t *Tokenizer) parseTripleString(startStringRune rune, startPosition token.PositionToken) (token.StringToken, TokenError) {
	var a string
	raw := string(startStringRune + startStringRune + startStringRune)
	for {
		ch := t.nextRune()
		raw += string(ch)
		if ch == 0 {
			return token.StringToken{}, NewEncounteredEOF()
		}

		wasEnd, err := t.isTriple(ch, startStringRune)
		if err != nil {
			return token.StringToken{}, NewUnexpectedEatTokenError(t.MakeSourceFileReference(startPosition), '\'', ' ')
		}

		if wasEnd {
			break
		}
		a += string(ch)
	}
	posLen := t.MakeSourceFileReference(startPosition)
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
	multilineCommentToken, err := t.ReadStringUntilEndOfMultilineComment(positionToken)
	if err != nil {
		return token.MultiLineCommentToken{}, err
	}
	return multilineCommentToken, nil
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
	singlePart := token.MultiLineCommentPart{
		SourceFileReference: t.MakeSourceFileReference(positionToken),
		RawString:           "--" + s,
		CommentString:       s,
	}
	return token.NewMultiLineCommentToken([]token.MultiLineCommentPart{singlePart}, documentationComment)
}

func (t *Tokenizer) ReadStringUntilEndOfMultilineComment(pos token.PositionToken) (token.MultiLineCommentToken, TokenError) {
	firstCh := t.nextRune()
	if firstCh == 0 {
		return token.MultiLineCommentToken{}, NewInternalError(fmt.Errorf("unexpected end of file"))
	}
	documentationComment := false
	if firstCh == '|' {
		documentationComment = true
	} else {
		t.unreadRune()
	}

	s := ""
	var parts []token.MultiLineCommentPart

	for {
		r := t.nextRune()
		if r == '\n' {
			t.unreadRune()
			sourceFileReference := t.MakeSourceFileReference(pos)
			part := token.MultiLineCommentPart{
				SourceFileReference: sourceFileReference,
				RawString:           s,
				CommentString:       s,
			}
			t.nextRune()
			pos = t.ParsingPosition()
			parts = append(parts, part)
			s = ""
		} else if r == '-' {
			if t.nextRune() == '}' {
				sourceFileReference := t.MakeSourceFileReference(pos)
				part := token.MultiLineCommentPart{
					SourceFileReference: sourceFileReference,
					RawString:           s,
					CommentString:       s,
				}
				parts = append(parts, part)
				break
			} else {
				t.unreadRune()
			}
		} else if r == 0 {
			return token.MultiLineCommentToken{}, NewInternalError(fmt.Errorf("unexpected end of file"))
		}

		s += string(r)
	}

	return token.NewMultiLineCommentToken(parts, documentationComment), nil
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

func (t *Tokenizer) ReadEndOrSeparatorToken() (token.Token, TokenError) {
	posToken := t.position
	r := t.nextRune()
	singleCharLength := t.MakeSourceFileReference(posToken)
	if r == ')' {
		return token.NewParenToken(string(r), token.RightParen, singleCharLength, " )R "), nil
	} else if r == '}' {
		return token.NewParenToken(string(r), token.RightCurlyBrace, singleCharLength, " } "), nil
	} else if r == ']' {
		return token.NewParenToken(string(r), token.RightSquareBracket, singleCharLength, " ] "), nil
	} else if r == ',' {
		return token.NewOperatorToken(token.Comma, singleCharLength, string(r), ","), nil
	} else if r == '|' {
		r := t.nextRune()
		if r == ']' {
			return token.NewParenToken(string(r), token.RightArrayBracket, singleCharLength, "|]"), nil
		} else if r == '>' {
			return token.NewOperatorToken(token.OperatorPipeRight, singleCharLength, "", "|>"), nil
		} else {
			t.unreadRune()
		}
	} else if r == '<' {
		r := t.nextRune()
		if r == '|' {
			return token.NewOperatorToken(token.OperatorPipeLeft, singleCharLength, "", "<|"), nil
		} else {
			t.unreadRune()
			return token.NewParenToken(string(r), token.LeftAngleBracket, singleCharLength, "<"), nil
		}
	} else if r == '-' {
		n := t.nextRune()
		if n == '>' {
			return token.NewOperatorToken(token.OperatorArrowRight, singleCharLength, "", "->"), nil
		}
		t.unreadRune()
	} else if r == '>' {
		return token.NewParenToken(string(r), token.RightAngleBracket, singleCharLength, ">"), nil
	}
	t.unreadRune()
	return nil, NewNotEndToken(singleCharLength, r)
}

func (t *Tokenizer) ReadTermTokenOrEndOrSeparator() (token.Token, error) {
	tokenFound, err := t.ReadEndOrSeparatorToken()
	if err == nil {
		return tokenFound, nil
	}
	return t.readTermToken()
}

func (t *Tokenizer) ReadOpenOperatorToken(r rune, singleCharLength token.SourceFileReference) (token.Token, TokenError) {
	posToken := t.position
	if r == '(' {
		return token.NewParenToken(string(r), token.LeftParen, singleCharLength, " L( "), nil
	} else if r == '{' {
		nch := t.nextRune()
		if nch == '-' {
			return t.ReadMultilineComment(posToken)
		}
		t.unreadRune()
		return token.NewParenToken(string(r), token.LeftCurlyBrace, singleCharLength, " { "), nil
	} else if r == '[' {
		nch := t.nextRune()
		if nch == '|' {
			return token.NewParenToken(string(r), token.LeftArrayBracket, singleCharLength, " [| "), nil
		} else {
			t.unreadRune()
		}
		return token.NewParenToken(string(r), token.LeftSquareBracket, singleCharLength, " [ "), nil
	} else if r == '*' {
		return token.NewOperatorToken(token.OperatorMultiply, singleCharLength, "*", "*"), nil
	}

	return nil, NewNotAnOpenOperatorError(t.MakeSourceFileReference(posToken), r)
}

func (t *Tokenizer) readTermToken() (token.Token, TokenError) {
	posToken := t.position
	r := t.nextRune()
	singleCharLength := t.MakeSourceFileReference(posToken)
	if r == 0 {
		return &EndOfFile{}, nil
	}
	if isNewLine(r) {
		return token.NewLineDelimiter(t.MakeSourceFileReference(posToken)), nil
	}
	t.lastTokenWasDelimiter = false
	if r == '%' || r == '$' {
		n := t.nextRune()
		if n == '"' {
			if r == '%' {
				return t.ParseStringInterpolationTuple('"', posToken)
			} else {
				return t.ParseStringInterpolationString('"', posToken)
			}
		} else {
			t.unreadRune()
		}
	}

	if isLetter(r) {
		t.unreadRune()
		return t.parseAnySymbol(posToken)
	} else if isDigit(r) {
		t.unreadRune()
		return t.ParseNumber("")
	} else if r == '@' {
		return t.parseResourceName(posToken)
	} else if r == '$' {
		return t.parseTypeId(posToken)
	} else if r == '\'' {
		return t.ParseCharacter(posToken)
	} else if isStartString(r) {
		if wasTriple, _ := t.isTriple(r, r); wasTriple {
			return t.parseTripleString(r, posToken)
		}
		return t.ParseString(r, posToken)
	} else if isUnaryOperator(r) {
		t.unreadRune()
		return t.ParseUnaryOperator()
	} else if r == '|' {
		next := t.nextRune()
		t.unreadRune()
		if next == ' ' {
			return token.NewGuardToken(singleCharLength, string(r), " guard "), nil
		}
		return nil, NewUnexpectedEatTokenError(singleCharLength, ' ', ' ')
	} else if r == '-' {
		return token.NewOperatorToken(token.OperatorUnaryMinus, singleCharLength, string(r), "unary-"), nil
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
		return token.NewSpaceToken(t.MakeSourceFileReference(t.position), r), nil
	} else {
		if foundOperator, operatorErr := t.ReadOpenOperatorToken(r, singleCharLength); operatorErr == nil {
			return foundOperator, nil
		}
	}
	return nil, NewUnexpectedEatTokenError(t.MakeSourceFileReference(posToken), r, ' ')
}

func (t *Tokenizer) ReadTermToken() (token.Token, TokenError) {
	startPos := t.position
	token, err := t.readTermToken()
	if err != nil {
		return nil, TokenizerError{err: err, position: t.MakeSourceFileReference(startPos)}
	}
	return token, nil
}

func (e *EndOfFile) String() string {
	return "EOF"
}
