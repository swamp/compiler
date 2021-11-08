/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

func LinePaddingBefore(expression Node) int {
	if expression == nil {
		return 0
	}

	_, isAnnotation := expression.(*Annotation)
	isDoubleLineStatement := isAnnotation

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
	_, isAnnotation := expression.(*Annotation)
	isSingleLineStatement := isImport || isAnnotation

	lines := 1
	if isSingleLineStatement {
		lines = 0
	}

	return lines
}

func ExpectedLinePaddingAfter(expression Node) int {
	if expression == nil {
		return 0
	}

	_, isImport := expression.(*Import)
	dontCare := isImport

	_, isAnnotation := expression.(*Annotation)
	mustBeSingleLine := isAnnotation

	lines := 3

	if dontCare {
		lines = -1
	}

	if mustBeSingleLine {
		lines = 1
	}

	return lines
}

func LinesToInsertBetween(before Expression, now Expression) int {
	previousPadding := LinePaddingAfter(before)
	beforePadding := LinePaddingBefore(now)

	totalPadding := previousPadding + beforePadding

	return totalPadding
}
