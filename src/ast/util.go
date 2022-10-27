/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

func LinePaddingBefore(expression Node) int {
	if expression == nil {
		return 0
	}

	isDoubleLineStatement := false

	lines := 0
	if isDoubleLineStatement {
		lines = 2
	}

	return lines
}

func LinePaddingAfter(expression Node) int {
	if expression == nil {
		return 0
	}

	_, isImport := expression.(*Import)
	isSingleLineStatement := isImport

	lines := 1
	if isSingleLineStatement {
		lines = 0
	}

	return lines
}

func ExpectedLinePaddingAfter(expression Node) (int, int) {
	if expression == nil {
		return 0, 0
	}

	_, isImport := expression.(*Import)
	dontCare := isImport

	mustBeSingleLine := false

	lines := 3
	_, wasConstant := expression.(*ConstantDefinition)
	if wasConstant {
		return 1, 3
	}

	if dontCare {
		lines = -1
	} else if mustBeSingleLine {
		lines = 1
	}

	return lines, lines
}

func LinesToInsertBetween(before Expression, now Expression) int {
	previousPadding := LinePaddingAfter(before)
	beforePadding := LinePaddingBefore(now)

	totalPadding := previousPadding + beforePadding

	return totalPadding
}
