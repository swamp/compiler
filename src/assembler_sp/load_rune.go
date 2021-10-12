package assembler_sp

import (
	"fmt"

	"github.com/swamp/compiler/src/instruction_sp"
)

type LoadRune struct {
	target TargetStackPos
	rune   instruction_sp.ShortRune
}

func (o *LoadRune) String() string {
	return fmt.Sprintf("[loadrune %v <= %v]", o.target, o.rune)
}
