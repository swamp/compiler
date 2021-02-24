/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"github.com/swamp/compiler/src/token"
)

type EndOfFile struct {
	token.SourceFileReference
}

func (e *EndOfFile) Position() token.Position {
	return token.Position{}
}

func (e *EndOfFile) FetchPositionLength() token.SourceFileReference {
	return e.SourceFileReference
}

func (e *EndOfFile) Type() token.Type {
	return token.EOF
}

func (t *Tokenizer) ParseUnaryOperator() (token.Token, TokenError) {
	startPosition := t.position

	ch := t.nextRune()
	if ch == 0 {
		return nil, NewEncounteredEOF()
	}
	raw := string(ch)
	debugString := ""
	var operatorType token.Type
	switch ch {
	case '-':
		next := t.nextRune()
		if next == ' ' {
			t.unreadRune()
			t.unreadRune()
			return nil, NewUnexpectedEatTokenError(t.MakeSourceFileReference(startPosition), 'x', ' ')
		} else if next == '>' {
			return nil, NewUnexpectedEatTokenError(t.MakeSourceFileReference(startPosition), 'y', ' ')
		}
		t.unreadRune()
		operatorType = token.OperatorUnaryMinus
	case '!':
		operatorType = token.OperatorUnaryNot

	}
	return token.NewOperatorToken(operatorType, t.MakeSourceFileReference(startPosition), raw, debugString), nil
}

func (t *Tokenizer) ParseOperator() (token.Token, TokenError) {
	startPosition := t.position
	ch := t.nextRune()
	if ch == 0 {
		return nil, NewEncounteredEOF()
	}
	raw := string(ch)
	debugString := ""
	var operatorType token.Type
	switch ch {
	case '>':
		nch := t.nextRune()
		if nch == '=' {
			operatorType = token.OperatorGreaterOrEqual
			raw += string(nch)
		} else {
			operatorType = token.OperatorGreater
			t.unreadRune()
		}
	case '<':
		nch := t.nextRune()
		if nch == '=' || nch == '|' {
			raw += string(nch)
			if nch == '|' {
				operatorType = token.OperatorPipeLeft
			} else {
				operatorType = token.OperatorLessOrEqual
			}
		} else {
			operatorType = token.OperatorLess
			t.unreadRune()
		}
	case '&':
		nch := t.nextRune()
		if nch == '&' {
			raw += string(nch)
			debugString = "AND"
			operatorType = token.OperatorAnd
		} else {
			operatorType = token.OperatorBitwiseAnd
			t.unreadRune()
		}
	case '|':
		nch := t.nextRune()
		if nch == '>' {
			raw += string(nch)
			operatorType = token.OperatorPipeRight
		} else if nch == '|' {
			raw += string(nch)
			debugString = "OR"
			operatorType = token.OperatorOr
		} else if nch == ']' {
			raw += string(nch)
			debugString = "|]"
			operatorType = token.RightArrayBracket
		} else {
			operatorType = token.OperatorUpdate
			t.unreadRune()
		}
	case '-':
		nch := t.nextRune()
		if nch == '>' {
			operatorType = token.OperatorArrowRight
			raw += string(nch)
		} else if nch == '-' {
			commentString := t.ReadStringUntilEndOfLine()
			return token.NewSingleLineCommentToken("--"+commentString, commentString, false, t.MakeSourceFileReference(startPosition)), nil
		} else {
			if isDigit(nch) {
				t.unreadRune()
				return t.ParseNumber("-")
			}
			operatorType = token.OperatorMinus
			t.unreadRune()
		}
	case '=':
		nch := t.nextRune()
		if nch == '=' {
			raw += string(nch)
			operatorType = token.OperatorEqual
		} else {
			t.unreadRune()
			operatorType = token.OperatorAssign
		}
	case '!':
		nch := t.nextRune()
		if nch == '=' {
			raw += string(nch)
			operatorType = token.OperatorNotEqual
		} else {
			t.unreadRune()
			operatorType = token.OperatorUnaryNot
		}
	default:
		switch ch {
		case '+':
			nch := t.nextRune()
			if nch == '+' {
				operatorType = token.OperatorAppend
				raw += string(nch)
			} else {
				t.unreadRune()
				operatorType = token.OperatorPlus
			}
		case ':':
			nch := t.nextRune()
			if nch == ':' {
				operatorType = token.OperatorCons
			} else {
				t.unreadRune()
				operatorType = token.Colon
			}
		case '*':
			operatorType = token.OperatorMultiply
		case '.':
			operatorType = token.Accessor
		case '/':
			operatorType = token.OperatorDivide
		case '^':
			operatorType = token.OperatorBitwiseXor
		case '~':
			operatorType = token.OperatorBitwiseNot
		default:
			return nil, NewUnexpectedEatTokenError(t.MakeSourceFileReference(startPosition), ' ', ch)
		}
	}
	if debugString == "" {
		debugString = raw
	}
	return token.NewOperatorToken(operatorType, t.MakeSourceFileReference(startPosition), raw, debugString), nil
}
