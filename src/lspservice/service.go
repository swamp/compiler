package lspservice

import (
	"fmt"

	"github.com/piot/go-lsp"
	"github.com/piot/lsp-server/lspserv"

	"github.com/swamp/compiler/src/token"
)

type DecoratedTypeOrToken interface {
	HumanReadableString() string
	SourceFile() *token.SourceFile
	Range() token.Range
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

func (s *Service) RunForever() {
	lspserv.RunForever("", s)
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
	s.compiler.Compile(string(params.TextDocument.URI))

	tokenPosition := lspToTokenPosition(params.Position)

	decoratedToken := s.scanner.FindToken(tokenPosition)

	if decoratedToken == nil {
		return nil, fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	hover := &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MUKMarkdown,
			Value: fmt.Sprintf("%v", decoratedToken.HumanReadableString()),
		},
		Range: tokenToLspRange(decoratedToken.Range()),
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
	return nil, nil
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
	return nil, nil
}

func (s *Service) HandleCodeLens(params lsp.CodeLensParams, conn lspserv.Connection) ([]*lsp.CodeLens, error) {
	return nil, nil
}

func (s *Service) HandleCodeLensResolve(params lsp.CodeLens, conn lspserv.Connection) (*lsp.CodeLens, error) {
	return nil, nil
}

func (s *Service) HandleDidChangeWatchedFiles(params lsp.DidChangeWatchedFilesParams, conn lspserv.Connection) error {
	return nil
}

func RunForever() {
	service := &Service{}
	lspserv.RunForever("", service)
}
