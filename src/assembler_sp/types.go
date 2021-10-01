package assembler_sp

import (
	"fmt"
)

type ConstantPosAndRange struct {
	Pos         uint32
	Size        uint32
	DebugString string
}

func (s ConstantPosAndRange) getPosition() uint32 {
	return s.Pos
}

func (s ConstantPosAndRange) getSize() uint32 {
	return s.Size
}

type StackPos uint32

type StackRange uint16

type StackItemSize uint16

type ZeroMemoryRange uint16

type StackPosOffset uint32

type MemoryAlign uint8

type StackPosAndRange struct {
	Pos         StackPos
	Size        StackRange
	DebugString string
}

func (s StackPosAndRange) getPosition() StackPos {
	return s.Pos
}

func (s StackPosAndRange) getSize() StackRange {
	return s.Size
}

type StackPosOffsetAndRange struct {
	Pos  StackPosOffset
	Size StackRange
}

func (s StackPosOffsetAndRange) getOffset() StackPosOffset {
	return s.Pos
}

func (s StackPosOffsetAndRange) getSize() StackRange {
	return s.Size
}

type TargetStackPosRange struct {
	Pos  TargetStackPos
	Size StackRange
}

type FieldRanges []SourceStackPosRange

type SourceStackPosRangeCompound StackPosAndRange

type TargetStackPos StackPos

func (t TargetStackPos) String() string {
	return fmt.Sprintf("targetPos: %04X", uint32(t))
}

type TargetFieldOffset uint16

type VariableArgumentPosSize struct {
	Offset uint16
	Size   uint16
}

type VariableArgumentPosSizeAlign struct {
	Offset uint16
	Size   uint16
	Align  uint8
}

type SourceStackPosAndRangeToLocalOffset struct {
	PosRange     SourceStackPosRange
	TargetOffset TargetFieldOffset
}

type SourceStackPos StackPos

type SourceStackPosOffsetRange StackPosOffsetAndRange

type SourceStackRange StackRange

type SourceStackPosRange struct {
	Pos  SourceStackPos
	Size SourceStackRange
}

type DynamicMemoryPos uint32

type SourceDynamicMemoryPos uint32

func (t SourceDynamicMemoryPos) String() string {
	return fmt.Sprintf("DynPos %04X", uint32(t))
}

type DynamicMemoryRange uint16

type SourceDynamicMemoryPosRange struct {
	Position SourceDynamicMemoryPos
	Size     DynamicMemoryRange
}

func (t SourceDynamicMemoryPosRange) String() string {
	return fmt.Sprintf("%v:%v", t.Position, t.Size)
}
