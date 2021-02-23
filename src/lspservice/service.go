package lspservice

import (
	"fmt"
	"net/url"

	"github.com/piot/go-lsp"
	"github.com/piot/lsp-server/lspserv"

	"github.com/swamp/compiler/src/token"
)

type DecoratedTypeOrToken interface {
	String() string
	// SourceFile() *token.SourceFileURI
	FetchPositionLength() token.Range
}

type DecoratedType interface {
	DecoratedTypeOrToken
	DecoratedHumanReadableString() string
}

type DecoratedToken interface {
	DecoratedTypeOrToken
	DecoratedType() DecoratedType
}

type DecoratedTokenScanner interface {
	FindToken(position token.Position) DecoratedTypeOrToken
}

type Compiler interface {
	Compile(filename string) error
}

type Service struct {
	scanner  DecoratedTokenScanner
	compiler Compiler
}

func NewService(compiler Compiler, scanner DecoratedTokenScanner) *Service {
	return &Service{scanner: scanner, compiler: compiler}
}

func (s *Service) Reset() error {
	return nil
}

func (s *Service) ResetCaches(lock bool) {
}

func (s *Service) ShutDown() {
}

func lspToTokenPosition(position lsp.Position) token.Position {
	return token.MakePosition(position.Line, position.Character)
}

func tokenToLspPosition(position token.Position) lsp.Position {
	return lsp.Position{
		Line:      position.Line(),
		Character: position.Column(),
	}
}

func tokenToLspRange(rangeToken token.Range) *lsp.Range {
	return &lsp.Range{
		Start: tokenToLspPosition(rangeToken.Start()),
		End:   tokenToLspPosition(rangeToken.End()),
	}
}

func (s *Service) HandleHover(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Hover, error) {
	fullUrl, urlErr := url.Parse(string(params.TextDocument.URI))
	if urlErr != nil {
		return nil, urlErr
	}
	s.compiler.Compile(fullUrl.Path)

	tokenPosition := lspToTokenPosition(params.Position)

	decoratedToken := s.scanner.FindToken(tokenPosition)

	if decoratedToken == nil {
		return nil, fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	hover := &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MUKMarkdown,
			Value: "this is **markup** content\n---\nIs this the last line?", //  decoratedToken.String()
		},
		Range: tokenToLspRange(decoratedToken.FetchPositionLength()),
	}

	return hover, nil
}

func (s *Service) HandleGotoDefinition(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleGotoDeclaration(params lsp.DeclarationOptions, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleGotoTypeDefinition(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleGotoImplementation(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleFindReferences(params lsp.ReferenceParams, conn lspserv.Connection) ([]*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleSymbol(params lsp.DocumentSymbolParams, conn lspserv.Connection) ([]*lsp.DocumentSymbol, error) {
	return []*lsp.DocumentSymbol{
		{
			Name:   "name",
			Detail: "String name",
			Kind:   lsp.SKProperty,
			Tags:   nil,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 0,
				},
				End: lsp.Position{
					Line:      0,
					Character: 4,
				},
			},
			SelectionRange: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 0,
				},
				End: lsp.Position{
					Line:      0,
					Character: 4,
				},
			},
			Children: nil,
		},
		{
			Name:   "2",
			Detail: "Int",
			Kind:   lsp.SKNumber,
			Tags:   nil,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 6,
				},
				End: lsp.Position{
					Line:      0,
					Character: 6,
				},
			},
			SelectionRange: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 6,
				},
				End: lsp.Position{
					Line:      0,
					Character: 6,
				},
			},
			Children: nil,
		},
	}, nil
} // Used for outline

func (s *Service) HandleCompletion(params lsp.CompletionParams, conn lspserv.Connection) (*lsp.CompletionList, error) {
	return nil, nil
} // Intellisense when pressing '.'.

func (s *Service) HandleCompletionItemResolve(params lsp.CompletionItem, conn lspserv.Connection) (*lsp.CompletionItem, error) {
	return nil, nil
}

func (s *Service) HandleSignatureHelp(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.SignatureHelp, error) {
	return nil, nil
}

func (s *Service) HandleFormatting(params lsp.DocumentFormattingParams, conn lspserv.Connection) ([]*lsp.TextEdit, error) {
	return nil, nil
}

// HandleRangeFormatting
func (s *Service) HandleHighlights(params lsp.DocumentHighlightParams, conn lspserv.Connection) ([]*lsp.DocumentHighlight, error) {
	return nil, nil
}

func (s *Service) HandleCodeAction(params lsp.CodeActionParams, conn lspserv.Connection) (*lsp.CodeAction, error) {
	return nil, nil
}

func (s *Service) HandleCodeActionResolve(params lsp.CodeAction, conn lspserv.Connection) (*lsp.CodeAction, error) {
	return nil, nil
}

func (s *Service) HandleRename(params lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	return nil, nil
}

func (s *Service) HandleSemanticTokensFull(params lsp.SemanticTokensParams, conn lspserv.Connection) (*lsp.SemanticTokens, error) {
	return &lsp.SemanticTokens{
		ResultId: "",
		Data: []uint{
			0, 0, 5, 0, 1,
			2, 0, 5, 3, 2,
		},
	}, nil
}

func (s *Service) HandleCodeLens(params lsp.CodeLensParams, conn lspserv.Connection) ([]*lsp.CodeLens, error) {
	return []*lsp.CodeLens{
		{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      4,
					Character: 0,
				},
				End: lsp.Position{
					Line:      4,
					Character: 4,
				},
			},
			Command: lsp.Command{
				Title:     "Some Command here",
				Command:   "swamp.somecommand",
				Arguments: nil,
			},
			Data: nil,
		},
	}, nil
}

func (s *Service) HandleCodeLensResolve(params lsp.CodeLens, conn lspserv.Connection) (*lsp.CodeLens, error) {
	return nil, nil
}

func (s *Service) HandleDidChangeWatchedFiles(params lsp.DidChangeWatchedFilesParams, conn lspserv.Connection) error {
	return nil
}
