package assembler_sp

import (
	"fmt"

	dectype "github.com/swamp/compiler/src/decorated/types"
)

type ZeroMemoryPointer uint32

type StackMemoryMapper struct {
	position     uint32
	maxOctetSize uint
	memory       []byte
}

func NewStackMemoryMapper(maxOctetSize uint) *StackMemoryMapper {
	return &StackMemoryMapper{maxOctetSize: maxOctetSize}
}

func (m *StackMemoryMapper) AlignUpForMax() {
	rest := m.position % uint32(dectype.Alignof64BitPointer)
	if rest != 0 {
		m.position += uint32(dectype.Alignof64BitPointer) - rest
	}
}

func (m *StackMemoryMapper) Allocate(octetSize uint, align uint32, debugString string) TargetStackPosRange {
	if octetSize == 0 {
		panic(fmt.Errorf("octet size zero is not allowed for allocate stack memory"))
	}
	if align == 0 {
		panic(fmt.Errorf("align zero size is not allowed for allocate stack memory"))
	}
	extra := m.position % align
	if extra != 0 {
		m.position += align - extra
	}
	pos := m.position

	m.position += uint32(octetSize)

	posRange := TargetStackPosRange{
		Pos:  TargetStackPos(pos),
		Size: StackRange(octetSize),
	}
	return posRange
}

func (m *StackMemoryMapper) Set(pos TargetStackPos) {
	m.position = uint32(pos)
}
