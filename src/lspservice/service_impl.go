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
	workspace     *loader.Workspace
	documentCache *DocumentCache
}

func NewLspImpl(fallbackProvider loader.DocumentProvider) *LspImpl {
	return &LspImpl{
		workspace:     loader.NewWorkspace(loader.LocalFileSystemRoot("")),
		documentCache: NewDocumentCache(fallbackProvider),
	}
}

func (l *LspImpl) Compile(filename string) (*decorated.Module, error) {
	const enforceStyle = true

	const verboseFlag = false

	world, module, err := swampcompiler.CompileMainFindLibraryRoot(filename, l.documentCache, enforceStyle, verboseFlag)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, fmt.Errorf("module can not be nil!")
	}
	fmt.Fprintf(os.Stderr, "COMPILE DONE!\n")
	l.workspace.AddOrReplacePackage(world)

	return module, nil
}

func findModuleFromSourceFile(world *loader.Package, sourceFileURI token.DocumentURI) (*decorated.Module, error) {
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

func (l *LspImpl) AllModules() []*decorated.Module {
	var allModules []*decorated.Module
	for _, foundPackage := range l.workspace.AllPackages() {
		allModules = append(allModules, foundPackage.AllModules()...)
	}

	return allModules
}

func (l *LspImpl) FindModuleHelper(sourceFile token.DocumentURI) *decorated.Module {
	localPath, err := sourceFile.ToLocalFilePath()
	if err != nil {
		return nil
	}
	module, _ := l.workspace.FindModuleFromSourceFile(loader.LocalFileSystemPath(localPath))
	if module == nil {
		log.Printf("could not find source file %v\n", sourceFile)
		return nil
	}

	return module
}

func (l *LspImpl) RootTokens(sourceFile token.DocumentURI) []decorated.TypeOrToken {
	module := l.FindModuleHelper(sourceFile)
	if module == nil {
		return nil
	}
	var tokens []decorated.TypeOrToken
	for _, node := range module.RootNodes() {
		tokens = append(tokens, node.(decorated.TypeOrToken))
	}

	return tokens
}

func (l *LspImpl) FindToken(sourceFile token.DocumentURI, position token.Position) decorated.TypeOrToken {
	module := l.FindModuleHelper(sourceFile)
	if module == nil {
		return nil
	}

	tokens := module.Nodes()

	smallestRange := token.MakeRange(
		token.MakePosition(0, 0, -1),
		token.MakePosition(9999999, 0, -1))

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
			// log.Printf("considered %v %T", foundRange, decoratedToken)
			if foundRange.SmallerThan(smallestRange) {
				smallestRange = foundRange
				bestToken = decoratedToken
			}
		}
	}

	if bestToken == nil {
		log.Printf("FindToken: couldn't find anything at %v %v\n", sourceFile, position)
	} else {
		log.Printf("FindToken: best is: %T %v\n", bestToken, bestToken.FetchPositionLength().ToCompleteReferenceString())
	}

	return bestToken
}

func (l *LspImpl) GetDocument(localFilePath LocalFileSystemPath, newVersion DocumentVersion) (*InMemoryDocument, error) {
	inMemoryDocument, err := l.documentCache.GetDocumentByVersion(localFilePath, newVersion-1)
	if err != nil {
		return nil, err
	}
	return inMemoryDocument, nil
}
