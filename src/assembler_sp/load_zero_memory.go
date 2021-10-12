package assembler_sp

import (
	"fmt"
)

type LoadZeroMemoryPointer struct {
	target           TargetStackPos
	sourceZeroMemory SourceDynamicMemoryPos
}

func (o *LoadZeroMemoryPointer) String() string {
	return fmt.Sprintf("[loadzeromem %v <= %v]", o.target, o.sourceZeroMemory)
}
