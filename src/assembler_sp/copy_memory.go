package assembler_sp

import (
	"fmt"
)

type CopyMemory struct {
	target TargetStackPos
	source SourceStackPosRange
}

func (o *CopyMemory) String() string {
	return fmt.Sprintf("[copymemory %v <= %v %v]", o.target, o.source)
}
