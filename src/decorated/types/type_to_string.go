/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
)

func TypesToHumanReadable(types []dtype.Type) string {
	s := ""
	for index, t := range types {
		if index > 0 {
			s += ","
		}
		s += t.HumanReadable()
	}

	return s
}

func TypesToHumanReadableWithinBrackets(types []dtype.Type) string {
	if len(types) == 0 {
		return ""
	}

	return fmt.Sprintf("<%v>", TypesToHumanReadable(types))
}

func TypesToString(types []dtype.Type) string {
	s := ""
	for index, t := range types {
		if index > 0 {
			s += ","
		}
		s += t.String()
	}

	return s
}

func TypesToStringSuffix(types []dtype.Type) string {
	if len(types) == 0 {
		return ""
	}
	s := " ["
	for index, t := range types {
		if index > 0 {
			s += ","
		}
		s += t.String()
	}
	s += "]"

	return s
}

func TypeArgumentsToString(types []*dtype.LocalTypeName) string {
	s := ""
	for index, t := range types {
		if index > 0 {
			s += ","
		}
		s += t.String()
	}

	return s
}

func TypeParametersSuffix(types []dtype.Type) string {
	if len(types) == 0 {
		return ""
	}

	return fmt.Sprintf("<%s>", TypesToString(types))
}

func TypeParametersHumanReadableSuffix(types []dtype.Type) string {
	if len(types) == 0 {
		return ""
	}

	return fmt.Sprintf(" %s", TypesToHumanReadable(types))
}
