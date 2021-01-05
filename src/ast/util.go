/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package ast

func LinePaddingBefore(expression Expression) int {
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


func LinePaddingAfter(expression Expression) int {
	if expression == nil {
		return 0
	}

	_, isExternalFn := expression.(*ExternalFunction)
	_, isImport := expression.(*Import)
	_, isAnnotation := expression.(*Annotation)
	isSingleLineStatement := isExternalFn || isImport || isAnnotation

	lines := 1
	if isSingleLineStatement {
		lines = 0
	}

	return lines
}



func ExpectedLinePaddingAfter(expression Expression) int {
	if expression == nil {
		return 0
	}

	_, isExternalFn := expression.(*ExternalFunction)
	_, isImport := expression.(*Import)
	dontCare := isExternalFn || isImport

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
    beforePadding :=  LinePaddingBefore(now)

    totalPadding := previousPadding + beforePadding

    return totalPadding
}
