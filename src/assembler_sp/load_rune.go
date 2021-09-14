package assembler_sp

import (
	"fmt"
)

type LoadRune struct {
	target TargetStackPos
	rune   uint8
}

func (o *LoadRune) String() string {
	return fmt.Sprintf("[loadrune %v <= %v %v]", o.target, o.rune)
}
