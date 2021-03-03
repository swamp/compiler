/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

type Type uint

const (
	OperatorPlus Type = iota
	OperatorMinus
	OperatorUnaryMinus
	OperatorMultiply
	OperatorDivide
	OperatorPipeRight
	OperatorPipeLeft
	OperatorArrowRight
	OperatorAnd
	OperatorOr
	OperatorGreater
	OperatorGreaterOrEqual
	OperatorLess
	OperatorLessOrEqual
	OperatorEqual
	OperatorNotEqual
	OperatorUnaryNot
	OperatorBitwiseOr
	OperatorBitwiseAnd
	OperatorBitwiseXor
	OperatorBitwiseNot
	OperatorAppend
	OperatorCons
	Colon
	OperatorAssign
	Guard
	OperatorUpdate
	Accessor
	RightParen
	LeftParen
	RightCurlyBrace
	LeftCurlyBrace
	RightBracket
	LeftBracket
	RightArrayBracket
	LeftArrayBracket
	Comma
	If
	Then
	Else
	Case
	Of
	Alias
	As
	Exposing
	TypeDef
	Asm
	ExternalFunction
	Let
	In
	VariableSymbol
	ResourceNameSymbol
	TypeIdSymbol
	TypeSymbol
	Lambda
	Import
	BooleanType
	NumberInteger
	NumberFixed
	StringConstant
	RuneConstant
	CommentConstant
	NewLine
	Space
	EOF
)
