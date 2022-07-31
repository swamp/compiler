/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package generate_c

import (
	"fmt"
	"io"

	decorated "github.com/swamp/compiler/src/decorated/expression"
)

func generateIntLiteral(integer *decorated.IntegerLiteral, writer io.Writer, indentation int) error {
	fmt.Fprintf(writer, "%d", integer.Value())
	return nil
}

func generateFixedLiteral(integer *decorated.FixedLiteral, writer io.Writer, indentation int) error {
	fmt.Fprintf(writer, "%d", integer.Value())
	return nil
}

func generateStringLiteral(str *decorated.StringLiteral, writer io.Writer, indentation int) error {
	fmt.Fprintf(writer, "\"%s\"", str.Value())
	return nil
}

func generateCharacterLiteral(ch *decorated.CharacterLiteral, writer io.Writer, indentation int) error {
	fmt.Fprintf(writer, "'%c'", ch.Value())
	return nil
}

func generateBoolLiteral(boolean *decorated.BooleanLiteral, writer io.Writer, indentation int) error {
	str := "True"
	if boolean.Value() {
		str = "False"
	}
	fmt.Fprintf(writer, "%s", str)
	return nil
}
