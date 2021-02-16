/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package tokenize

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type TokenError interface {
	FetchPositionLength() token.PositionLength
	Error() string
}

type SubError struct {
	SubErr TokenError
}

func NewSubError(err TokenError) SubError {
	if err == nil {
		panic("new sub err")
	}
	return SubError{SubErr: err}
}

func (e SubError) FetchPositionLength() token.PositionLength {
	return e.SubErr.FetchPositionLength()
}

type StandardTokenError struct {
	posLength token.PositionLength
}

func (e StandardTokenError) FetchPositionLength() token.PositionLength {
	return e.posLength
}

type UnexpectedEatTokenError struct {
	StandardTokenError
	requiredRune    rune
	encounteredRune rune
}

func NewUnexpectedEatTokenError(posLength token.PositionLength, requiredRune rune, encounteredRune rune) UnexpectedEatTokenError {
	return UnexpectedEatTokenError{StandardTokenError: StandardTokenError{posLength}, requiredRune: requiredRune, encounteredRune: encounteredRune}
}

func (e UnexpectedEatTokenError) Error() string {
	return fmt.Sprintf("unexpected rune. required %v, but encountered %v", string(e.requiredRune), string(e.encounteredRune))
}

type InternalError struct {
	err error
}

func NewInternalError(err error) InternalError {
	return InternalError{err: err}
}

func (e InternalError) Error() string {
	return fmt.Sprintf("tokenize internal error %v", e.err)
}

func (e InternalError) FetchPositionLength() token.PositionLength {
	return token.NewPositionLength(token.NewPositionTopLeft(), 0, -1)
}

type ExpectedVariableSymbolError struct {
	StandardTokenError
	encountered string
}

func NewExpectedVariableSymbolError(posLength token.PositionLength, encountered string) ExpectedVariableSymbolError {
	return ExpectedVariableSymbolError{StandardTokenError: StandardTokenError{posLength}, encountered: encountered}
}

func (e ExpectedVariableSymbolError) Error() string {
	return fmt.Sprintf("expected variable symbol but encountered %v", e.encountered)
}

type ExpectedTypeSymbolError struct {
	StandardTokenError
	encountered string
}

func NewExpectedTypeSymbolError(posLength token.PositionLength, encountered string) ExpectedTypeSymbolError {
	return ExpectedTypeSymbolError{StandardTokenError: StandardTokenError{posLength}, encountered: encountered}
}

func (e ExpectedTypeSymbolError) Error() string {
	return fmt.Sprintf("expected type symbol but encountered %v", e.encountered)
}

type EncounteredEOF struct {
	StandardTokenError
}

func NewEncounteredEOF() EncounteredEOF {
	return EncounteredEOF{}
}

func (e EncounteredEOF) Error() string {
	return fmt.Sprintf("EOF")
}

type ExpectedNewLineError struct {
	eatError TokenError
}

func NewExpectedNewLineError(eatError TokenError) ExpectedNewLineError {
	return ExpectedNewLineError{eatError: eatError}
}

func (e ExpectedNewLineError) Error() string {
	return fmt.Sprintf("expected newline ")
}

func (e ExpectedNewLineError) FetchPositionLength() token.PositionLength {
	return e.eatError.FetchPositionLength()
}

type ExpectedNewLineAndIndentationError struct {
	SubError
	expectedIndentation    int
	encounteredIndentation int
}

func NewExpectedNewLineAndIndentationError(expectedIndentation int, encounteredIndentation int, err TokenError) ExpectedNewLineAndIndentationError {
	return ExpectedNewLineAndIndentationError{SubError: NewSubError(err), expectedIndentation: expectedIndentation, encounteredIndentation: encounteredIndentation}
}

func (e ExpectedNewLineAndIndentationError) Error() string {
	return fmt.Sprintf("expected newline and indentation %v, but encountered %v", e.expectedIndentation, e.encounteredIndentation)
}

type ExpectedIndentationAfterNewLineError struct {
	SubError
	expectedIndentation    int
	encounteredIndentation int
}

func NewExpectedIndentationAfterNewLineError(expectedIndentation int, encounteredIndentation int, err TokenError) ExpectedIndentationAfterNewLineError {
	return ExpectedIndentationAfterNewLineError{SubError: NewSubError(err), expectedIndentation: expectedIndentation, encounteredIndentation: encounteredIndentation}
}

func (e ExpectedIndentationAfterNewLineError) Error() string {
	return fmt.Sprintf("expected indentation %v (%v spaces), but encountered %v (%v)", e.expectedIndentation, e.expectedIndentation*SpacesForIndentation, e.encounteredIndentation, e.SubError)
}

type IllegalIndentationError struct {
	posLength         token.PositionLength
	encounteredSpaces int
	multiples         int
}

func NewIllegalIndentationError(posLength token.PositionLength, encounteredSpaces int, multiples int) IllegalIndentationError {
	return IllegalIndentationError{posLength: posLength, encounteredSpaces: encounteredSpaces, multiples: multiples}
}

func (e IllegalIndentationError) Error() string {
	return fmt.Sprintf("illegal indentation, found %d spaces. Must be multiple of %d", e.encounteredSpaces, e.multiples)
}

func (e IllegalIndentationError) FetchPositionLength() token.PositionLength {
	return e.posLength
}

type UnexpectedIndentationError struct {
	posLength              token.PositionLength
	requiredIndentation    int
	encounteredIndentation int
}

func NewUnexpectedIndentationError(posLength token.PositionLength, requiredIndentation int, encounteredIndentation int) UnexpectedIndentationError {
	return UnexpectedIndentationError{posLength: posLength, requiredIndentation: requiredIndentation, encounteredIndentation: encounteredIndentation}
}

func (e UnexpectedIndentationError) Error() string {
	return fmt.Sprintf("unexpected indentation, wanted %d but encountered %d", e.requiredIndentation, e.encounteredIndentation)
}

func (e UnexpectedIndentationError) FetchPositionLength() token.PositionLength {
	return e.posLength
}

type ExpectedOneSpaceError struct {
	SubError
}

func NewExpectedOneSpaceError(err TokenError) ExpectedOneSpaceError {
	return ExpectedOneSpaceError{SubError: NewSubError(err)}
}

func (e ExpectedOneSpaceError) Error() string {
	return fmt.Sprintf("expected one single space %v", e.SubError)
}

type ExpectedIndentationError struct {
	SubError
}

func NewExpectedIndentationError(err TokenError) ExpectedIndentationError {
	return ExpectedIndentationError{SubError: NewSubError(err)}
}

func (e ExpectedIndentationError) Error() string {
	return fmt.Sprintf("expected exact indentation ")
}

type IllegalCharacterError struct {
	StandardTokenError
	encountered rune
}

func NewIllegalCharacterError(posLength token.PositionLength, encountered rune) IllegalCharacterError {
	return IllegalCharacterError{StandardTokenError: StandardTokenError{posLength}, encountered: encountered}
}

func (e IllegalCharacterError) Error() string {
	return fmt.Sprintf("illegal character %d", e.encountered)
}

type TrailingSpaceError struct {
	StandardTokenError
}

func NewTrailingSpaceError(posLength token.PositionLength) TrailingSpaceError {
	return TrailingSpaceError{StandardTokenError: StandardTokenError{posLength}}
}

func (e TrailingSpaceError) Error() string {
	return fmt.Sprintf("illegal trailing space error %v", e.StandardTokenError)
}

type CommentNotAllowedHereError struct {
	StandardTokenError
	subError error
}

func NewCommentNotAllowedHereError(posLength token.PositionLength, subError error) CommentNotAllowedHereError {
	return CommentNotAllowedHereError{StandardTokenError: StandardTokenError{posLength}, subError: subError}
}

func (e CommentNotAllowedHereError) Error() string {
	return fmt.Sprintf("not allowed to have comment here %v %v", e.StandardTokenError, e.subError)
}
