/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package token

import "fmt"

// GuardToken :
type GuardToken struct {
	SourceFileReference
	raw         string
	debugString string
}

func NewGuardToken(startPosition SourceFileReference, raw string, debugString string) GuardToken {
	return GuardToken{SourceFileReference: startPosition, raw: raw, debugString: debugString}
}

func (s GuardToken) Type() Type {
	return Guard
}

func (s GuardToken) String() string {
	return s.debugString
}

func (s GuardToken) Raw() string {
	return s.raw
}

func (s GuardToken) DebugString() string {
	return fmt.Sprintf("[guard %s]", s.debugString)
}

func (s GuardToken) FetchPositionLength() SourceFileReference {
	return s.SourceFileReference
}
