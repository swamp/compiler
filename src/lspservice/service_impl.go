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
	module, moduleErr := findModuleFromSourceFile(l.world, sourceFile)
	if moduleErr != nil {
		return nil
	}
	var tokens []decorated.TypeOrToken
	for _, node := range module.RootNodes() {
		tokens = append(tokens, node.(decorated.TypeOrToken))
	}

	return tokens
}

func (l *LspImpl) FindToken(sourceFile token.DocumentURI, position token.Position) decorated.TypeOrToken {
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
