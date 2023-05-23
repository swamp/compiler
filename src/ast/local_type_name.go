package ast

import (
	"fmt"

	"github.com/swamp/compiler/src/token"
)

type LocalTypeName struct {
	ident *VariableIdentifier `debug:"true"`
}

func (l *LocalTypeName) Name() string {
	return l.ident.Name()
}

func (l *LocalTypeName) Identifier() *VariableIdentifier {
	return l.ident
}

func (l *LocalTypeName) FetchPositionLength() token.SourceFileReference {
	return l.ident.FetchPositionLength()
}

func (l *LocalTypeName) String() string {
	return fmt.Sprintf("[LocalTypeName %v]", l.ident.Name())
}

func NewLocalTypeName(ident *VariableIdentifier) *LocalTypeName {
	return &LocalTypeName{ident: ident}
}
