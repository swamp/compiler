package assembler_sp

import (
	"fmt"
)

type LoadBool struct {
	target  TargetStackPos
	boolean bool
}

func (o *LoadBool) String() string {
	return fmt.Sprintf("[loadbool %v <= %v]", o.target, o.boolean)
}
