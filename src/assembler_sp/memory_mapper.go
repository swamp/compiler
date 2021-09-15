package assembler_sp

type ZeroMemoryPointer uint32

type StackMemoryMapper struct {
	position     uint32
	maxOctetSize uint
	memory       []byte
}

func NewStackMemoryMapper(maxOctetSize uint) *StackMemoryMapper {
	return &StackMemoryMapper{maxOctetSize: maxOctetSize}
}

func (m *StackMemoryMapper) Allocate(octetSize uint, align uint32, debugString string) TargetStackPosRange {
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

type ZeroMemoryMapper struct {
	position        uint32
	maxOctetSize    uint
	memory          []byte
	maxIndexWritten uint32
}

func ZeroMemoryMapperNew(maxOctetSize uint) *ZeroMemoryMapper {
	return &ZeroMemoryMapper{maxOctetSize: maxOctetSize, memory: make([]byte, maxOctetSize)}
}

func (m *ZeroMemoryMapper) Allocate(octetSize uint, debugString string) SourceZeroMemoryPosRange {
	pos := SourceZeroMemoryPos(m.position)

	m.position += uint32(octetSize)

	return SourceZeroMemoryPosRange{Position: pos, Size: ZeroMemoryRange(octetSize)}
}

func (m *ZeroMemoryMapper) Write(data []byte, debugString string) SourceZeroMemoryPosRange {
	posRange := m.Allocate(uint(len(data)), debugString)
	position := posRange.Position
	endPos := uint32(position) + uint32(len(data))
	if endPos > m.maxIndexWritten {
		m.maxIndexWritten = endPos - 1
	}
	copy(m.memory[position:endPos], data)

	return posRange
}

func (m *ZeroMemoryMapper) Read(pos SourceZeroMemoryPosRange) []byte {
	start := pos.Position
	endPos := uint32(start) + uint32(pos.Size)

	return m.memory[start:endPos]
}

type DynamicMemoryMapper struct {
	position        uint32
	maxOctetSize    uint
	memory          []byte
	maxIndexWritten uint32
}

func DynamicMemoryMapperNew(maxOctetSize uint) *DynamicMemoryMapper {
	return &DynamicMemoryMapper{maxOctetSize: maxOctetSize, memory: make([]byte, maxOctetSize)}
}

func (m *DynamicMemoryMapper) Allocate(octetSize uint, debugString string) SourceDynamicMemoryPosRange {
	pos := SourceDynamicMemoryPos(m.position)

	m.position += uint32(octetSize)

	return SourceDynamicMemoryPosRange{Position: pos, Size: DynamicMemoryRange(octetSize)}
}

func (m *DynamicMemoryMapper) Write(data []byte, debugString string) SourceDynamicMemoryPosRange {
	posRange := m.Allocate(uint(len(data)), debugString)
	position := posRange.Position
	endPos := uint32(position) + uint32(len(data))
	if endPos > m.maxIndexWritten {
		m.maxIndexWritten = endPos - 1
	}
	copy(m.memory[position:endPos], data)

	return posRange
}
