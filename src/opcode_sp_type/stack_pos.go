/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package opcode_sp_type

import "fmt"

type StackPosition uint32

func (r StackPosition) String() string {
	return fmt.Sprintf("%v", uint32(r))
}

type StackRange uint16

type SourceStackRange StackRange

type TargetStackPosition StackPosition

func (r TargetStackPosition) String() string {
	return fmt.Sprintf("%v<-", uint32(r))
}

type SourceStackPosition StackPosition

func (r SourceStackPosition) String() string {
	return fmt.Sprintf("%v", uint32(r))
}

type StackPositionRange struct {
	Position StackPosition
	Range    StackRange
}

func (r StackPositionRange) String() string {
	return fmt.Sprintf("%v:%v", r.Position, r.Range)
}

type SourceStackPositionRange struct {
	Position SourceStackPosition
	Range    SourceStackRange
}

func (r SourceStackPositionRange) String() string {
	return fmt.Sprintf("%v:%v", r.Position, r.Range)
}

type TargetStackPositionRange StackPositionRange

func (r TargetStackPositionRange) String() string {
	return fmt.Sprintf("%v:%v<-", r.Position, r.Range)
}

type TargetFieldOffset StackRange

func (r TargetFieldOffset) String() string {
	return fmt.Sprintf("#%v", uint32(r))
}
