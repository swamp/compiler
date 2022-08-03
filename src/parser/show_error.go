/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package parser

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/swamp/compiler/src/verbosity"

	"github.com/swamp/compiler/src/coloring"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	parerr "github.com/swamp/compiler/src/parser/errors"
	"github.com/swamp/compiler/src/pathutil"
	"github.com/swamp/compiler/src/token"
	"github.com/swamp/compiler/src/tokenize"
)

func ShowSourceCode(tokenizer *tokenize.Tokenizer, highlightLine int,
	highlightColumn int, posLength token.SourceFileReference) error {
	const beforeAndAfterCount = 3
	startRow := highlightLine - beforeAndAfterCount
	if startRow < 0 {
		startRow = 0
	}

	endRow := highlightLine + beforeAndAfterCount
	rowCount := endRow - startRow + 1

	if tokenizer == nil {
		return fmt.Errorf("we need a tokenizer to show this error")
	}

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
		fmt.Fprintf(os.Stderr, "%v\n", displayRow)
		if actualRow == highlightLine {
			skipSpaces := highlightColumn
			indentString := strings.Repeat(" ", skipSpaces)
			repeatCount := posLength.Range.RuneWidth()
			if repeatCount < 1 {
				repeatCount = 1
			}
			underlineString := color.HiYellowString(strings.Repeat("^", repeatCount))
			fmt.Fprintf(os.Stderr, "%v%v\n", indentString, underlineString)
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
	// ColorTypesWithAtom(e.ExpectedFunctionType().ParameterTypes(), indentation, true, colorer)
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

type ReportAsSeverity uint8

const (
	ReportAsSeverityNote ReportAsSeverity = iota
	ReportAsSeverityInfo
	ReportAsSeverityWarning
	ReportAsSeverityError
)

func ShowError(tokenizer *tokenize.Tokenizer, filename string, parserError parerr.ParseError,
	verbose verbosity.Verbosity, errorAsWarning ReportAsSeverity) tokenize.TokenError {
	if parserError == nil {
		panic("parserError is nil. internal error.")
	}
	_, isModuleErr := parserError.(*decorated.ModuleError)
	if isModuleErr {
		panic("can not have multi errors")
	}

	parseAliasErr, _ := parserError.(*parerr.ParseAliasError)
	if parseAliasErr != nil {
		panic("can not have multi errors")
	}
	_, wasMulti := parserError.(*tokenize.MultiErrors)
	if wasMulti {
		panic("can not have multi errors")
	}

	_, wasMultiParErr := parserError.(parerr.MultiError)
	if wasMultiParErr {
		panic("can not have multi errors")
	}

	posLength := parserError.FetchPositionLength()
	highlightLine := posLength.Range.Position().Line()
	highlightColumn := posLength.Range.Position().Column()

	messageError := parserError.Error()
	severityString := "Error"
	colorToUse := color.New(color.FgRed)
	switch errorAsWarning {
	case ReportAsSeverityWarning:
		severityString = "Warning"
		colorToUse = color.New(color.FgHiYellow)
	case ReportAsSeverityNote:
		severityString = "Note"
		colorToUse = color.New(color.FgHiBlue)
	}

	pathToShow := pathutil.TryToMakeRelativePath(filename)

	coloredErrorMessage := colorToUse.Sprintf("%v: %v", severityString, messageError)

	errorString := fmt.Sprintf("%v:%d:%d: %v", pathToShow, highlightLine+1, highlightColumn+1,
		coloredErrorMessage)
	fmt.Fprintf(os.Stderr, "%v %T\n", errorString, parserError)

	color.NoColor = false

	if errorAsWarning == ReportAsSeverityError {
		ShowSourceCode(tokenizer, highlightLine, highlightColumn, posLength)
	}

	const useColor = true
	var colorer coloring.Colorer
	if useColor {
		colorer = coloring.NewColorerWithColor()
	} else {
		colorer = coloring.NewColorerWithoutColor()
	}

	indentation := 1
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
		// log.Printf("internal: I have no good description for error %T\n", e)
	}

	return nil
}

func HighestSeverity(err error) ReportAsSeverity {
	highestError := ReportAsSeverityNote
	if err == nil {
		return highestError
	}

	moduleErr, wasModuleErr := err.(*decorated.ModuleError)
	if wasModuleErr {
		return HighestSeverity(moduleErr.WrappedError())
	}

	switch t := err.(type) {
	case *decorated.MultiErrors:
		for _, subErr := range t.Errors() {
			detectedError := HighestSeverity(subErr)
			if detectedError > highestError {
				highestError = detectedError
			}
		}
		return highestError
	case parerr.MultiError:
		for _, subErr := range t.Errors() {
			detectedError := HighestSeverity(subErr)
			if detectedError > highestError {
				highestError = detectedError
			}
		}
		return highestError
	case *tokenize.MultiErrors:
		for _, subErr := range t.Errors() {
			detectedError := HighestSeverity(subErr)
			if detectedError > highestError {
				highestError = detectedError
			}
		}
		return highestError
	}

	parserErr, wasParserErr := err.(parerr.ParseError)
	if wasParserErr {
		return TypeOfWarning(parserErr)
	}

	log.Printf("unknown err %v %T", err, err)
	return ReportAsSeverityError
}

func IsCompileError(parseError parerr.ParseError) bool {
	return HighestSeverity(parseError) == ReportAsSeverityError
}

func IsCompileErr(parseError error) bool {
	return HighestSeverity(parseError) == ReportAsSeverityError
}

func TypeOfWarning(parserError parerr.ParseError) ReportAsSeverity {
	switch parserError.(type) {
	case parerr.ExpectedOneSpace:
		return ReportAsSeverityWarning
	case parerr.UnexpectedImportAlias:
		return ReportAsSeverityNote
	case parerr.TooManyDepths:
		return ReportAsSeverityWarning
	case *decorated.UnusedWarning:
		return ReportAsSeverityNote
	case *decorated.UnusedTypeWarning:
		return ReportAsSeverityNote
	case tokenize.LineIsLongerThanRecommendedError:
		return ReportAsSeverityNote
	case tokenize.LineIsTooLongError:
		return ReportAsSeverityWarning
	case tokenize.LineCountIsMoreThanRecommendedError:
		return ReportAsSeverityWarning
	}

	return ReportAsSeverityError
}

func ShowWarningOrError(tokenizer *tokenize.Tokenizer, parserError parerr.ParseError) ReportAsSeverity {
	moduleErr, isModuleErr := parserError.(*decorated.ModuleError)
	if isModuleErr {
		parserError = moduleErr.WrappedError()
	}
	parseAliasErr, _ := parserError.(*parerr.ParseAliasError)
	if parseAliasErr != nil {
		parserError = parseAliasErr.Unwrap()
	}

	multi, wasMulti := parserError.(*tokenize.MultiErrors)
	if wasMulti {
		highestError := ReportAsSeverityNote
		for _, tokenizeErr := range multi.Errors() {
			detectedError := ShowWarningOrError(tokenizer, tokenizeErr)
			if detectedError > highestError {
				highestError = detectedError
			}
		}
		return highestError
	}

	parErrMulti, wasParErrMulti := parserError.(parerr.MultiError)
	if wasParErrMulti {
		highestError := ReportAsSeverityNote
		for _, tokenizeErr := range parErrMulti.Errors() {
			detectedError := ShowWarningOrError(tokenizer, tokenizeErr)
			if detectedError > highestError {
				highestError = detectedError
			}
		}
		return highestError
	}

	decoratedMultiErr, wasDecoratedMultiErr := parserError.(*decorated.MultiErrors)
	if wasDecoratedMultiErr {
		highestError := ReportAsSeverityNote
		for _, tokenizeErr := range decoratedMultiErr.Errors() {
			detectedError := ShowWarningOrError(tokenizer, tokenizeErr)
			if detectedError > highestError {
				highestError = detectedError
			}
		}
		return highestError
	}
	showAsWarning := TypeOfWarning(parserError)
	localPath := ""
	if parserError.FetchPositionLength().Document != nil {
		localPath, _ = parserError.FetchPositionLength().Document.Uri.ToLocalFilePath()
	}
	ShowError(tokenizer, localPath, parserError, verbosity.High, showAsWarning)

	return showAsWarning
}

func ShowAsError(tokenizer *tokenize.Tokenizer, parserError parerr.ParseError) ReportAsSeverity {
	localPath := ""
	if parserError.FetchPositionLength().Document != nil {
		localPath, _ = parserError.FetchPositionLength().Document.Uri.ToLocalFilePath()
	}

	ShowError(tokenizer, localPath, parserError, verbosity.High, ReportAsSeverityError)

	return ReportAsSeverityError
}
