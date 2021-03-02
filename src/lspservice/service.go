package lspservice

import (
	"fmt"
	"log"
	"net/url"

	"github.com/piot/go-lsp"
	"github.com/piot/lsp-server/lspserv"

	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	dectype "github.com/swamp/compiler/src/decorated/types"
	"github.com/swamp/compiler/src/token"
)

type DecoratedTokenScanner interface {
	FindToken(documentURI token.DocumentURI, position token.Position) decorated.TypeOrToken
	RootTokens(documentURI token.DocumentURI) []decorated.TypeOrToken
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

func toDocumentURI(URI lsp.DocumentURI) token.DocumentURI {
	return token.DocumentURI(URI)
}

func lspToTokenPosition(position lsp.Position) token.Position {
	return token.MakePosition(position.Line, position.Character)
}

func sourceFileReferenceToLocation(reference token.SourceFileReference) *lsp.Location {
	return &lsp.Location{
		URI:   lsp.DocumentURI(reference.Document.Uri),
		Range: *tokenToLspRange(reference.Range),
	}
}

func tokenToLspPosition(position token.Position) lsp.Position {
	return lsp.Position{
		Line:      position.Line(),
		Character: position.Column(),
	}
}

func tokenToLspRange(rangeToken token.Range) *lsp.Range {
	exclusiveEndPosition := lsp.Position{
		Line:      rangeToken.End().Line(),
		Character: rangeToken.End().Column() + 1,
	}

	return &lsp.Range{
		Start: tokenToLspPosition(rangeToken.Start()),
		End:   exclusiveEndPosition,
	}
}

func (s *Service) HandleHover(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Hover, error) {
	tokenPosition := lspToTokenPosition(params.Position)
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	decoratedToken := s.scanner.FindToken(sourceFileURI, tokenPosition)

	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	tokenType, isItAType := decoratedToken.(dtype.Type)

	log.Printf("the token is %T and type:%v\n", decoratedToken, isItAType)

	var codeSignature string
	if isItAType {
		codeSignature = tokenType.HumanReadable()
	} else {
		normalToken, _ := decoratedToken.(decorated.Token)
		codeSignature = normalToken.Type().HumanReadable()
	}

	showString := fmt.Sprintf("```swamp\n%v\n```", codeSignature)

	hover := &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.MUKMarkdown,
			Value: showString,
		},
		Range: tokenToLspRange(decoratedToken.FetchPositionLength().Range),
	}

	return hover, nil
}

func (s *Service) HandleGotoDefinition(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	tokenPosition := lspToTokenPosition(params.Position)
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	decoratedToken := s.scanner.FindToken(sourceFileURI, tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	var sourceFileReference token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.Import:
		sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterReference:
		sourceFileReference = t.ParameterRef().FetchPositionLength()
	case *decorated.LetVariableReference:
		sourceFileReference = t.LetVariable().FetchPositionLength()
	case *decorated.FunctionReference:
		sourceFileReference = t.FunctionValue().FetchPositionLength()
	case *decorated.FunctionValue:
		sourceFileReference = t.FetchPositionLength()
	case *decorated.FunctionParameterDefinition:
		sourceFileReference = t.FetchPositionLength()
	case *decorated.Let:
		sourceFileReference = t.FetchPositionLength()
	case *decorated.LetVariable:
		sourceFileReference = t.FetchPositionLength()
	case *decorated.FunctionCall:
		sourceFileReference = t.FetchPositionLength()
	// TYPES
	case *dectype.FunctionTypeReference:
		sourceFileReference = t.FunctionAtom().FetchPositionLength()
	case *dectype.TypeReference:
		sourceFileReference = t.Next().FetchPositionLength()
	}

	if sourceFileReference.Document == nil {
		log.Printf("couldn't go to definition for %T\n", decoratedToken)
		return nil, nil
	}

	location := sourceFileReferenceToLocation(sourceFileReference)

	return location, nil
}

func (s *Service) HandleLinkedEditingRange(params lsp.LinkedEditingRangeParams, conn lspserv.Connection) (*lsp.LinkedEditingRanges, error) {
	tokenPosition := lspToTokenPosition(params.Position)
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	decoratedToken := s.scanner.FindToken(sourceFileURI, tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return &lsp.LinkedEditingRanges{
			Ranges:      []lsp.Range{},
			WordPattern: nil,
		}, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	documentURI := token.MakeDocumentURI(string(params.TextDocument.URI))
	sourceFileReferences := findAllLinkedSymbolsInDocument(decoratedToken, documentURI)

	var renameRanges []lsp.Range

	for _, ref := range sourceFileReferences {
		renameRanges = append(renameRanges, *tokenToLspRange(ref.Range))
	}

	return &lsp.LinkedEditingRanges{
		Ranges:      renameRanges,
		WordPattern: nil,
	}, nil
}

func (s *Service) HandleGotoDeclaration(params lsp.DeclarationOptions, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, fmt.Errorf("concept of go to declaration not in the Swamp language")
}

func (s *Service) HandleGotoTypeDefinition(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, nil
}

func (s *Service) HandleGotoImplementation(params lsp.TextDocumentPositionParams, conn lspserv.Connection) (*lsp.Location, error) {
	return nil, fmt.Errorf("concept of go to implementation not in the Swamp language")
}

func findReferences(uri lsp.DocumentURI, position lsp.Position, scanner DecoratedTokenScanner) ([]token.SourceFileReference, error) {
	tokenPosition := lspToTokenPosition(position)
	sourceFileURI := toDocumentURI(uri)
	decoratedToken := scanner.FindToken(sourceFileURI, tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	var sourceFileReferences []token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.Import:
		// sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterDefinition:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.LetVariable:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.LetAssignment:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	case *decorated.FunctionValue:
		for _, ref := range t.References() {
			sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
		}
	}

	return sourceFileReferences, nil
}

func (s *Service) HandleFindReferences(params lsp.ReferenceParams, conn lspserv.Connection) ([]*lsp.Location, error) {
	sourceFileReferences, err := findReferences(params.TextDocument.URI, params.Position, s.scanner)
	if err != nil {
		return nil, err
	}

	var locations []*lsp.Location

	for _, reference := range sourceFileReferences {
		location := sourceFileReferenceToLocation(reference)
		locations = append(locations, location)
	}

	return locations, nil
}

func functionParametersToDocumentSymbols(parameters []*decorated.FunctionParameterDefinition) []lsp.DocumentSymbol {
	var symbols []lsp.DocumentSymbol

	for _, param := range parameters {
		symbol := lsp.DocumentSymbol{
			Name:           param.Identifier().Name(),
			Detail:         param.Type().HumanReadable(),
			Kind:           lsp.SKVariable,
			Tags:           nil,
			Range:          *tokenToLspRange(param.FetchPositionLength().Range),
			SelectionRange: *tokenToLspRange(param.FetchPositionLength().Range),
			Children:       nil,
		}
		symbols = append(symbols, symbol)
	}

	return symbols
}

func convertRootTokenToOutlineSymbol(rootToken decorated.TypeOrToken) *lsp.DocumentSymbol {
	switch t := rootToken.(type) {
	case *decorated.FunctionValue:
		return &lsp.DocumentSymbol{
			Name:           t.AstFunctionValue().DebugFunctionIdentifier().Name(),
			Detail:         t.Type().HumanReadable(),
			Kind:           lsp.SKFunction,
			Tags:           nil,
			Range:          *tokenToLspRange(t.FetchPositionLength().Range),
			SelectionRange: *tokenToLspRange(t.FetchPositionLength().Range),
			Children:       functionParametersToDocumentSymbols(t.Parameters()),
		}
	}
	return nil
}

func (s *Service) HandleSymbol(params lsp.DocumentSymbolParams, conn lspserv.Connection) ([]*lsp.DocumentSymbol, error) {
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	rootTokens := s.scanner.RootTokens(sourceFileURI)

	var symbols []*lsp.DocumentSymbol

	for _, rootToken := range rootTokens {
		documentSymbol := convertRootTokenToOutlineSymbol(rootToken)
		if documentSymbol == nil {
			continue
		}
		symbols = append(symbols, documentSymbol)
	}
	return symbols, nil
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

func findAllLinkedSymbolsInDocument(decoratedToken decorated.TypeOrToken, filterDocument token.DocumentURI) []token.SourceFileReference {
	var sourceFileReferences []token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.Import:
		// sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterDefinition:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.FetchPositionLength())
		}
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionParameterReference:
		return findAllLinkedSymbolsInDocument(t.ParameterRef(), filterDocument)

	case *decorated.LetVariable:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.FetchPositionLength())
		}
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.LetVariableReference:
		return findAllLinkedSymbolsInDocument(t.LetVariable(), filterDocument)

		/*
			case *decorated.LetAssignment:
				for _, ref := range t.References() {
					if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
						sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
					}
				}

		*/

	case *decorated.FunctionValue:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.AstFunctionValue().DebugFunctionIdentifier().FetchPositionLength())
		}
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionReference:
		return findAllLinkedSymbolsInDocument(t.FunctionValue(), filterDocument)
	}

	return sourceFileReferences
}

func findLinkedSymbolsInDocument(decoratedToken decorated.TypeOrToken, filterDocument token.DocumentURI) []token.SourceFileReference {
	var sourceFileReferences []token.SourceFileReference

	switch t := decoratedToken.(type) {
	case *decorated.Import:
		// sourceFileReference = token.MakeSourceFileReference(token.MakeSourceFileDocumentFromURI(t.Module().DocumentURI()), token.MakeRange(token.MakePosition(0, 0), token.MakePosition(0, 0)))
	case *decorated.FunctionParameterDefinition:
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionParameterReference:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.ParameterRef().FetchPositionLength())
		}
	case *decorated.LetVariable:
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.LetVariableReference:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.LetVariable().FetchPositionLength())
		}
	case *decorated.LetAssignment:
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionValue:
		for _, ref := range t.References() {
			if ref.FetchPositionLength().Document.EqualTo(filterDocument) {
				sourceFileReferences = append(sourceFileReferences, ref.FetchPositionLength())
			}
		}
	case *decorated.FunctionReference:
		if t.FetchPositionLength().Document.EqualTo(filterDocument) {
			sourceFileReferences = append(sourceFileReferences, t.FunctionValue().FetchPositionLength())
		}
	}

	return sourceFileReferences
}

// HandleHighlights :
func (s *Service) HandleHighlights(params lsp.DocumentHighlightParams,
	conn lspserv.Connection) ([]*lsp.DocumentHighlight, error) {
	tokenPosition := lspToTokenPosition(params.Position)
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	decoratedToken := s.scanner.FindToken(sourceFileURI, tokenPosition)
	if decoratedToken == nil {
		log.Printf("couldn't find a token at %v\n", tokenPosition)
		return nil, nil // fmt.Errorf("couldn't find a token at %v", tokenPosition)
	}

	documentURI := token.MakeDocumentURI(string(params.TextDocument.URI))
	sourceFileReferences := findLinkedSymbolsInDocument(decoratedToken, documentURI)

	var highlights []*lsp.DocumentHighlight

	for _, reference := range sourceFileReferences {
		highlight := &lsp.DocumentHighlight{
			Range: *tokenToLspRange(reference.Range),
			Kind:  1, // Read only
		}

		highlights = append(highlights, highlight)
	}

	return highlights, nil
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

/*
"namespace",
"type",
"class",
"enum",
"interface",
"struct",
"typeParameter",
"parameter",
"variable",
"property",
"enumMember",
"event",
"function",
"method",
"macro",
"keyword",
"modifier",
"comment",
"string",
"number",
"regexp",
"operator",
},
tokenModifiers: []string{
"declaration",
"definition",
"readonly",
"static",
"deprecated",
"abstract",
"async",
"modification",
"documentation",
"defaultLibrary",

*/

func (s *Service) HandleSemanticTokensFull(params lsp.SemanticTokensParams, conn lspserv.Connection) (*lsp.SemanticTokens, error) {
	sourceFileURI := toDocumentURI(params.TextDocument.URI)
	allTokens := s.scanner.RootTokens(sourceFileURI)
	builder := NewSemanticBuilder()
	for _, foundToken := range allTokens {
		if err := addSemanticToken(foundToken, builder); err != nil {
			return nil, err
		}
	}
	return &lsp.SemanticTokens{
		ResultId: "",
		Data:     builder.EncodedValues(),
	}, nil
}

func (s *Service) HandleCodeLens(params lsp.CodeLensParams, conn lspserv.Connection) ([]*lsp.CodeLens, error) {
	var codeLenses []*lsp.CodeLens

	for _, rootToken := range s.scanner.RootTokens(toDocumentURI(params.TextDocument.URI)) {
		switch t := rootToken.(type) {
		case *decorated.FunctionValue:
			textToDisplay := fmt.Sprintf("%d references", len(t.References()))
			if len(t.References()) == 0 {
				textToDisplay = "no references"
			}
			codeLens := &lsp.CodeLens{
				Range: *tokenToLspRange(t.Annotation().FetchPositionLength().Range),
				Command: lsp.Command{
					Title:     textToDisplay,
					Command:   "",
					Arguments: nil,
				},
				Data: nil,
			}
			codeLenses = append(codeLenses, codeLens)
		}
	}

	return codeLenses, nil
}

func (s *Service) HandleCodeLensResolve(params lsp.CodeLens, conn lspserv.Connection) (*lsp.CodeLens, error) {
	return nil, nil
}

func (s *Service) HandleDidChangeWatchedFiles(params lsp.DidChangeWatchedFilesParams, conn lspserv.Connection) error {
	return nil
}

func (s *Service) HandleDidOpen(params lsp.DidOpenTextDocumentParams, conn lspserv.Connection) error {
	fullUrl, urlErr := url.Parse(string(params.TextDocument.URI))
	if urlErr != nil {
		return urlErr
	}
	compileErr := s.compiler.Compile(fullUrl.Path)
	if compileErr != nil {
		log.Printf("couldn't compile it:%v\n", compileErr)
	}
	return nil
}
