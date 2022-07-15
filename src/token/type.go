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
	OperatorRemainder
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
	OperatorBitwiseShiftLeft
	OperatorBitwiseShiftRight
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
	RightSquareBracket
	LeftSquareBracket
	LeftAngleBracket
	RightAngleBracket
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
	ExternalFunction
	Let
	In
	VariableSymbol
	ResourceNameSymbol
	TypeIdSymbol
	TypeSymbol
	Import
	BooleanType
	NumberInteger
	NumberFixed
	StringConstant
	StringInterpolationTupleConstant
	StringInterpolationStringConstant
	RuneConstant
	CommentConstant
	NewLine
	Space
	EOF
)
