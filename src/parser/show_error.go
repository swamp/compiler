/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/swamp/compiler/src/coloring"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/pathutil"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

func ShowSourceCode(tokenizer *tokenize.Tokenizer, highlightLine int,
	highlightColumn int, posLength token.PositionLength) error {
	const beforeAndAfterCount = 3
	startRow := highlightLine - beforeAndAfterCount
	if startRow < 0 {
		startRow = 0
	}

	endRow := highlightLine + beforeAndAfterCount
	rowCount := endRow - startRow + 1

	rows := tokenizer.ExtractStrings(startRow, rowCount)
	const useColor = true

	for rowIndex, row := range rows {
		actualRow := rowIndex + startRow
		displayRow := row
		if useColor {
			coloredRow, coloringErr := coloring.SyntaxColor(row)
			if coloringErr == nil {
				displayRow = coloredRow
			} else {
				return tokenize.NewInternalError(coloringErr)
			}
		}
		fmt.Printf("%v\n", displayRow)
		if actualRow == highlightLine {
			skipSpaces := highlightColumn
			indentString := strings.Repeat(" ", skipSpaces)
			repeatCount := posLength.RuneWidth()
			if repeatCount < 1 {
				repeatCount = 1
			}
			underlineString := color.HiYellowString(strings.Repeat("^", repeatCount))
			fmt.Printf("%v%v\n", indentString, underlineString)
		}
	}

	return nil
}

func doubleLine(colorer coloring.Colorer, indentation int) {
	colorer.NewLine(0)
	colorer.NewLine(indentation)
}

func userInstruction(s string, indentation int, colorer coloring.Colorer) {
	colorer.NewLine(0)
	colorer.NewLine(indentation)
	colorer.UserInstruction(s)
	colorer.NewLine(0)
	colorer.NewLine(indentation + 1)
}

func showUnmatchingError(unmatching *decorated.UnMatchingTypesError, indentation int, colorer coloring.Colorer) {
	userInstruction("I expected to get type:", indentation, colorer)

	indentation++
	doubleLine(colorer, indentation)

	ColorTypeWithAtom(unmatching.ExpectedType, indentation, false, colorer)
	indentation--
	newDoubleLine(indentation, colorer)
	userInstruction("but detected:", indentation, colorer)
	indentation++
	newDoubleLine(indentation, colorer)
	ColorType(unmatching.HasType, indentation, false, colorer)
}

func showWrongNumberOfArgumentsInFunctionValue(e *decorated.WrongNumberOfArgumentsInFunctionValue, indentation int, colorer coloring.Colorer) {
	userInstruction("wrong number of parameters in function value compared to annotation:", indentation, colorer)
	indentation++
	colorFunctionParametersWithAlias(e.EncounteredArgumentTypes(), indentation, false, colorer)
	indentation--
	userInstruction("but expected:", indentation, colorer)
	colorFunctionParametersWithAlias(e.ExpectedFunctionType().FunctionParameterTypes(), indentation+1, false, colorer)
	// ColorTypesWithAtom(e.ExpectedFunctionType().FunctionParameterTypes(), indentation, true, colorer)
}

func showAllItemsInListMustHaveSameType(e *decorated.EveryItemInThelistMustHaveTheSameType, indentation int, colorer coloring.Colorer) {
	userInstruction("all items in a list must have the same type. You started with the type:", indentation, colorer)
	indentation++
	ColorTypeWithAtom(e.ExpectedType, indentation, false, colorer)
	indentation--
	newDoubleLine(indentation, colorer)
	userInstruction("but I encountered:", indentation, colorer)
	indentation++
	ColorType(e.ActualType, indentation, false, colorer)
}

func showOneSpaceAfterRecordTypeColon(e *parerr.OneSpaceAfterRecordTypeColon, indentation int, colorer coloring.Colorer) {
	userInstruction("you must have a space after the ':' in a record type field:", indentation, colorer)
}

func ShowError(tokenizer *tokenize.Tokenizer, filename string, parserError parerr.ParseError,
	verbose bool, errorAsWarning bool) tokenize.TokenError {
	if parserError == nil {
		panic("parserError is nil. internal error.")
	}
	moduleErr, isModuleErr := parserError.(*decorated.ModuleError)
	if isModuleErr {
		parserError = moduleErr.WrappedError()
	}

	parseAliasErr, _ := parserError.(*parerr.ParseAliasError)
	if parseAliasErr != nil {
		parserError = parseAliasErr.Unwrap()
	}

	posLength := parserError.FetchPositionLength()
	highlightLine := posLength.Position().Line()
	highlightColumn := posLength.Position().Column()

	messageError := parserError.Error()
	severityString := "error"
	colorToUse := color.New(color.FgHiRed)
	if errorAsWarning {
		colorToUse = color.New(color.FgHiGreen)
	}

	pathToShow := pathutil.TryToMakeRelativePath(filename)

	errorString := colorToUse.Sprintf("%v:%d:%d: %v: %v", pathToShow, highlightLine+1, highlightColumn+1,
		severityString, messageError)
	fmt.Printf("\n%v\n", errorString)

	color.NoColor = false
	ShowSourceCode(tokenizer, highlightLine, highlightColumn, posLength)

	const useColor = true
	var colorer coloring.Colorer
	if useColor {
		colorer = coloring.NewColorerWithColor()
	} else {
		colorer = coloring.NewColorerWithoutColor()
	}

	indentation := 1
	colorer.NewLine(indentation)

	switch e := parserError.(type) {
	case *decorated.UnMatchingFunctionReturnTypesInFunctionValue:
		showUnmatchingError(&e.UnMatchingTypesError, indentation, colorer)
	case *decorated.WrongNumberOfArgumentsInFunctionValue:
		showWrongNumberOfArgumentsInFunctionValue(e, indentation, colorer)
	case *decorated.EveryItemInThelistMustHaveTheSameType:
		showAllItemsInListMustHaveSameType(e, indentation, colorer)
	case *parerr.OneSpaceAfterRecordTypeColon:
		showOneSpaceAfterRecordTypeColon(e, indentation, colorer)
	default:
		fmt.Printf("internal: I have no good description for error %T\n", e)
	}

	fmt.Println(colorer.String())

	return nil
}
