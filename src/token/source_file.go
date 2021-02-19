package token

import (
	"fmt"
)

type SourceFile struct {
	name string
}

func MakeSourceFile(name string) *SourceFile {
	return &SourceFile{name: name}
}

func (s *SourceFile) String() string {
	return s.name
}

func (s *SourceFile) ReferenceString() string {
	return fmt.Sprintf("%v:", s.name)
}

func (s *SourceFile) ReferenceWithPositionString(pos Position) string {
	return fmt.Sprintf("%v:%d:%d:", s.name, pos.line+1, pos.column+1)
}
