/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

type Precedence int

const (
	NONE Precedence = iota
	LOWEST
	UPDATE
	ASSIGN
	PIPE
	ANDOR
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	HIGHEST
)
