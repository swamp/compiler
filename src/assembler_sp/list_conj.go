package assembler_sp

import (
	"fmt"
)

type ListAppend struct {
	target TargetStackPos
	a      SourceStackPos
	b      SourceStackPos
}

func (o *ListAppend) String() string {
	return fmt.Sprintf("[listappend %v <= %v %v]", o.target, o.a, o.b)
}
