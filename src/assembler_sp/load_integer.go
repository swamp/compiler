package assembler_sp

import (
	"fmt"
)

type LoadInteger struct {
	target   TargetStackPos
	intValue int32
}

func (o *LoadInteger) String() string {
	return fmt.Sprintf("[loadinteger %v <= %v]", o.target, o.intValue)
}
