package assembler_sp

import (
	"fmt"
)

type SetEnum struct {
	target    TargetStackPos
	enumIndex uint8
}

func (o *SetEnum) String() string {
	return fmt.Sprintf("[setenum %v <= %v]", o.target, o.enumIndex)
}
