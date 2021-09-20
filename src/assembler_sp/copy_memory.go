package assembler_sp

import (
	"fmt"
)

type CopyMemory struct {
	target TargetStackPos
	source SourceStackPosRange
}

func NewCopyMemory(target TargetStackPos, source SourceStackPosRange) *CopyMemory {
	if source.Size == 0 {
		panic("not allowed copy zero size")
	}
	return &CopyMemory{
		target: target,
		source: source,
	}
}

func (o *CopyMemory) String() string {
	return fmt.Sprintf("[copymemory %v <= %v]", o.target, o.source)
}
