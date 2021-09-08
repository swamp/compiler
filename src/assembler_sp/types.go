package assembler_sp

type ConstantPosAndRange struct {
	Pos  uint32
	Size uint32
}

func (s ConstantPosAndRange) getPosition() uint32 {
	return s.Pos
}

func (s ConstantPosAndRange) getSize() uint32 {
	return s.Size
}

type StackPos uint32

type StackRange uint32

type StackPosOffset uint32

type StackPosAndRange struct {
	Pos  StackPos
	Size StackRange
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

type TargetStackPosRange StackPosAndRange

type SourceStackPosRange StackPosAndRange

type SourceStackPosRangeCompound StackPosAndRange

type TargetStackPos StackPos

type SourceStackPos StackPos

type SourceStackPosOffsetRange StackPosOffsetAndRange

type SourceStackRange StackRange
