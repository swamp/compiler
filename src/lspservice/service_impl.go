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

func (l *LspImpl) FindToken(position token.Position) decorated.DecoratedTypeOrToken {
	if l.module == nil {
		return nil
	}
	tokens := l.module.Nodes()

	smallestRange := token.MakeRange(
		token.MakePosition(0, 0),
		token.MakePosition(9999999, 0))

	var bestToken decorated.DecoratedTypeOrToken

	for _, decoratedToken := range tokens {
		log.Printf("checking node:%v '%v'\n", decoratedToken.FetchPositionLength(), decoratedToken.String())
		foundRange := decoratedToken.FetchPositionLength().Range
		if foundRange.Contains(position) {
			if foundRange.SmallerThan(smallestRange) {
				smallestRange = foundRange
				bestToken = decoratedToken
			}
		}
	}

	return bestToken
}
