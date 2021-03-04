package lspservice

import (
	"fmt"
	"log"
	"os"
	"reflect"

	swampcompiler "github.com/swamp/compiler/src/compiler"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/token"
)

type LspImpl struct {
	world         *loader.World
	documentCache *DocumentCache
}

func (l *LspImpl) NewLspImp(fallbackProvider loader.DocumentProvider) *LspImpl {
	return &LspImpl{
		world:         nil,
		documentCache: NewDocumentCache(fallbackProvider),
	}
}

func (l *LspImpl) Compile(filename string) error {
	const enforceStyle = true

	const verboseFlag = false

	world, module, err := swampcompiler.CompileMainFindLibraryRoot(filename, l.documentCache, enforceStyle, verboseFlag)
	if err != nil {
		return err
	}

	if module == nil {
		return fmt.Errorf("module can not be nil!")
	}
	fmt.Fprintf(os.Stderr, "COMPILE DONE!\n")

	l.world = world

	return nil
}

func findModuleFromSourceFile(world *loader.World, sourceFileURI token.DocumentURI) (*decorated.Module, error) {
	localFilePath, convertErr := sourceFileURI.ToLocalFilePath()
	if convertErr != nil {
		return nil, convertErr
	}

	foundModule := world.FindModuleFromAbsoluteFilePath(loader.LocalFileSystemPath(localFilePath))
	if foundModule == nil {
		return nil, fmt.Errorf("couldn't find module for path %v", localFilePath)
	}

	return foundModule, nil
}

func (l *LspImpl) RootTokens(sourceFile token.DocumentURI) []decorated.TypeOrToken {
	if reflect.ValueOf(l.world).IsNil() {
		log.Printf("world is nil, probably compilation didn't work")
		return nil
	}
	module, moduleErr := findModuleFromSourceFile(l.world, sourceFile)
	if moduleErr != nil {
		log.Printf("could not find source file %v\n", sourceFile)
		return nil
	}
	var tokens []decorated.TypeOrToken
	for _, node := range module.RootNodes() {
		tokens = append(tokens, node.(decorated.TypeOrToken))
	}

	return tokens
}

func (l *LspImpl) FindToken(sourceFile token.DocumentURI, position token.Position) decorated.TypeOrToken {
	if reflect.ValueOf(l.world).IsNil() {
		log.Printf("world is nil, probably compilation didn't work")
		return nil
	}
	module, moduleErr := findModuleFromSourceFile(l.world, sourceFile)
	if moduleErr != nil {
		return nil
	}
	tokens := module.Nodes()

	smallestRange := token.MakeRange(
		token.MakePosition(0, 0),
		token.MakePosition(9999999, 0))

	var bestToken decorated.TypeOrToken

	for _, decoratedToken := range tokens {
		if decoratedToken == nil {
			panic("can not be nil")
		}
		if reflect.ValueOf(decoratedToken).IsNil() {
			panic("bad here")
		}
		// log.Printf("checking node:%v '%v'\n", decoratedToken.FetchPositionLength(), decoratedToken.String())
		foundRange := decoratedToken.FetchPositionLength().Range
		if foundRange.Contains(position) {
			if foundRange.SmallerThan(smallestRange) {
				smallestRange = foundRange
				bestToken = decoratedToken
			}
		}
	}

	if bestToken == nil {
		log.Printf("FindToken: couldn't find anything at %v %v\n", sourceFile, position)
	} else {
		log.Printf("FindToken: best is: %T\n", bestToken)
	}

	return bestToken
}
