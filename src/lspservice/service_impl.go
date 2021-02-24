package lspservice

import (
	"fmt"
	"log"

	swampcompiler "github.com/swamp/compiler/src/compiler"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/token"
)

type LspImpl struct {
	world  *loader.World
	module *decorated.Module
}

func (l *LspImpl) Compile(filename string) error {
	const enforceStyle = true
	const verboseFlag = false
	world, module, err := swampcompiler.CompileFile(filename, enforceStyle, verboseFlag)
	if err != nil {
		return err
	}
	if module == nil {
		return fmt.Errorf("module can not be nil!")
	}
	l.world = world
	l.module = module
	return nil
}

func (l *LspImpl) FindToken(position token.Position) DecoratedTypeOrToken {
	if l.module == nil {
		return nil
	}
	allNodes := l.module.Nodes()
	for _, node := range allNodes {
		log.Printf("checking node:%v '%v'\n", node.FetchPositionLength(), node.String())
		if node.FetchPositionLength().Contains(position) {
			return node
		}
	}
	return nil
}
