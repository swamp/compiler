package lspservice_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/piot/lsp-server/lspserv"

	"github.com/swamp/compiler/src/lspservice"
)

type ReadWriteCloserTester struct {
	reader io.Reader
	buffer *bytes.Buffer
}

func NewReadWriteCloserTester(s string) *ReadWriteCloserTester {
	return &ReadWriteCloserTester{
		reader: strings.NewReader(s),
		buffer: bytes.NewBuffer(nil),
	}
}

func (r *ReadWriteCloserTester) Read(b []byte) (int, error) {
	return r.reader.Read(b)
}

func (r *ReadWriteCloserTester) Write(b []byte) (int, error) {
	return r.buffer.Write(b)
}

func (r *ReadWriteCloserTester) Close() error {
	return nil
}

func (r *ReadWriteCloserTester) Result() string {
	return r.buffer.String()
}

func convertToCrLf(s string) string {
	r := strings.NewReplacer("\r\n", "\r\n", "\n", "\r\n").Replace(s)

	return r
}

func convertFromCrLf(s string) string {
	r := strings.NewReplacer("\n", "\n", "\r\n", "\n").Replace(s)

	return r
}

func runCommand(service lspserv.Service, linuxEndings string) (string, error) {
	linuxEndingsTrimmed := strings.TrimSpace(linuxEndings)

	octetSize := len(linuxEndingsTrimmed)
	linuxEndingsTrimmed = fmt.Sprintf("Content-Length: %d\n\n", octetSize) + linuxEndingsTrimmed
	s := convertToCrLf(linuxEndingsTrimmed)

	rwc := NewReadWriteCloserTester(s)
	const logOutput = false
	service.RunUntilClose(rwc, logOutput)

	r := rwc.Result()
	r = convertFromCrLf(r)

	//	fmt.Fprintf(os.Stderr, "\nreceived '%s'\n", r)

	return r, nil
}

func testHelperWithTesting(t *testing.T, document string, s string) string {
	lspService := &lspservice.LspImpl{}
	service := lspservice.NewService(lspService, lspService)

	initializeCommand := `{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"processId":22206,"clientInfo":{"name":"Visual Studio Code","version":"1.53.2"},"locale":"en-us","rootPath":null,"rootUri":null,"capabilities":{"workspace":{"applyEdit":true,"workspaceEdit":{"documentChanges":true,"resourceOperations":["create","rename","delete"],"failureHandling":"textOnlyTransactional","normalizesLineEndings":true,"changeAnnotationSupport":{"groupsOnLabel":true}},"didChangeConfiguration":{"dynamicRegistration":true},"didChangeWatchedFiles":{"dynamicRegistration":true},"symbol":{"dynamicRegistration":true,"symbolKind":{"valueSet":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26]},"tagSupport":{"valueSet":[1]}},"codeLens":{"refreshSupport":true},"executeCommand":{"dynamicRegistration":true},"configuration":true,"workspaceFolders":true,"semanticTokens":{"refreshSupport":true},"fileOperations":{"dynamicRegistration":true,"didCreate":true,"didRename":true,"didDelete":true,"willCreate":true,"willRename":true,"willDelete":true}},"textDocument":{"publishDiagnostics":{"relatedInformation":true,"versionSupport":false,"tagSupport":{"valueSet":[1,2]},"codeDescriptionSupport":true,"dataSupport":true},"synchronization":{"dynamicRegistration":true,"willSave":true,"willSaveWaitUntil":true,"didSave":true},"completion":{"dynamicRegistration":true,"contextSupport":true,"completionItem":{"snippetSupport":true,"commitCharactersSupport":true,"documentationFormat":["markdown","plaintext"],"deprecatedSupport":true,"preselectSupport":true,"tagSupport":{"valueSet":[1]},"insertReplaceSupport":true,"resolveSupport":{"properties":["documentation","detail","additionalTextEdits"]},"insertTextModeSupport":{"valueSet":[1,2]}},"completionItemKind":{"valueSet":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25]}},"hover":{"dynamicRegistration":true,"contentFormat":["markdown","plaintext"]},"signatureHelp":{"dynamicRegistration":true,"signatureInformation":{"documentationFormat":["markdown","plaintext"],"parameterInformation":{"labelOffsetSupport":true},"activeParameterSupport":true},"contextSupport":true},"definition":{"dynamicRegistration":true,"linkSupport":true},"references":{"dynamicRegistration":true},"documentHighlight":{"dynamicRegistration":true},"documentSymbol":{"dynamicRegistration":true,"symbolKind":{"valueSet":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26]},"hierarchicalDocumentSymbolSupport":true,"tagSupport":{"valueSet":[1]},"labelSupport":true},"codeAction":{"dynamicRegistration":true,"isPreferredSupport":true,"disabledSupport":true,"dataSupport":true,"resolveSupport":{"properties":["edit"]},"codeActionLiteralSupport":{"codeActionKind":{"valueSet":["","quickfix","refactor","refactor.extract","refactor.inline","refactor.rewrite","source","source.organizeImports"]}},"honorsChangeAnnotations":false},"codeLens":{"dynamicRegistration":true},"formatting":{"dynamicRegistration":true},"rangeFormatting":{"dynamicRegistration":true},"onTypeFormatting":{"dynamicRegistration":true},"rename":{"dynamicRegistration":true,"prepareSupport":true,"prepareSupportDefaultBehavior":1,"honorsChangeAnnotations":true},"documentLink":{"dynamicRegistration":true,"tooltipSupport":true},"typeDefinition":{"dynamicRegistration":true,"linkSupport":true},"implementation":{"dynamicRegistration":true,"linkSupport":true},"colorProvider":{"dynamicRegistration":true},"foldingRange":{"dynamicRegistration":true,"rangeLimit":5000,"lineFoldingOnly":true},"declaration":{"dynamicRegistration":true,"linkSupport":true},"selectionRange":{"dynamicRegistration":true},"callHierarchy":{"dynamicRegistration":true},"semanticTokens":{"dynamicRegistration":true,"tokenTypes":["namespace","type","class","enum","interface","struct","typeParameter","parameter","variable","property","enumMember","event","function","method","macro","keyword","modifier","comment","string","number","regexp","operator"],"tokenModifiers":["declaration","definition","readonly","static","deprecated","abstract","async","modification","documentation","defaultLibrary"],"formats":["relative"],"requests":{"range":true,"full":{"delta":true}},"multilineTokenSupport":false,"overlappingTokenSupport":false},"linkedEditingRange":{"dynamicRegistration":true}},"window":{"showMessage":{"messageActionItem":{"additionalPropertiesSupport":true}},"showDocument":{"support":true},"workDoneProgress":true},"general":{"regularExpressions":{"engine":"ECMAScript","version":"ES2020"},"markdown":{"parser":"marked","version":"1.1.0"}}},"trace":"verbose","workspaceFolders":null}}`

	lspservService := lspserv.NewService(service)

	_, initializeErr := runCommand(lspservService, initializeCommand)
	if initializeErr != nil {
		t.Fatal(initializeErr)
	}

	didOpenNotification := fmt.Sprintf(`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"%v","languageId":"swamp","version":1,"text":""}}}`, document)
	_, didOpenErr := runCommand(lspservService, didOpenNotification)
	if didOpenErr != nil {
		t.Fatal(didOpenErr)
	}
	result, err := runCommand(lspservService, s)
	if err != nil {
		t.Fatal(err)
	}

	return result
}

func testHelperWithTestingStrings(t *testing.T, document string, cmds []string, expectedResult string) string {
	last := ""
	for _, cmd := range cmds {
		last = testHelperWithTesting(t, document, cmd)
	}

	if last != expectedResult {
		t.Errorf("mismatch. expected '%v' but received '%v'", expectedResult, last)
	}

	return last
}

func testHelperWithTestingString(t *testing.T, cmd string, expectedResult string) string {
	last := testHelperWithTesting(t, "file:///home/peter/test.swamp", cmd)

	if last != expectedResult {
		t.Errorf("mismatch. expected '%v' but received '%v'", expectedResult, last)
	}

	return last
}

func testHelperWithTestingStringDoc(t *testing.T, document string, cmd string, expectedResult string) string {
	last := testHelperWithTesting(t, document, cmd)

	if last != expectedResult {
		t.Errorf("mismatch. expected '%v' but received '%v'", expectedResult, last)
	}

	return last
}

func TestHover(t *testing.T) {
	//nolint: lll
	testHelperWithTestingStrings(t, "file:///home/peter/test.swamp", []string{`{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///home/peter/test.swamp","languageId":"swamp","version":1,"text":"test : Int -> Int\ntest a =\n    4\n"}}}`, `
{"jsonrpc":"2.0","id":1,"method":"textDocument/hover","params":{"textDocument":{"uri":"file:///home/peter/test.swamp"},"position":{"line":8,"character":8}}}
`}, `Content-Length: 156

{"id":1,"result":{"contents":{"kind":"markdown","value":"$test"},"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":3}}},"jsonrpc":"2.0"}`)
}

func TestDocumentSymbolOutline(t *testing.T) {
	testHelperWithTestingString(t, `{"jsonrpc":"2.0","id":4,"method":"textDocument/documentSymbol","params":{"textDocument":{"uri":"file:///home/peter/test.swamp"}}}`, ``)
}

func TestDocumentSemanticSymbols(t *testing.T) {
	testHelperWithTestingString(t, `{"jsonrpc":"2.0","id":2,"method":"textDocument/semanticTokens/full","params":{"textDocument":{"uri":"file:///home/peter/own/hackman/swamp/gameplay/Main.swamp"}}}`, ``)
}

func TestDocumentGotoDefinition(t *testing.T) {
	testHelperWithTestingString(t, `{"jsonrpc":"2.0","id":6,"method":"textDocument/definition","params":{"textDocument":{"uri":"file:///home/peter/test.swamp"},"position":{"line":0,"character":9}}}`, ``)
}

func TestDocumentGotoDefinition2(t *testing.T) {
	testHelperWithTestingString(t, `{"jsonrpc":"2.0","id":6,"method":"textDocument/definition","params":{"textDocument":{"uri":"file:///home/peter/test.swamp"},"position":{"line":13,"character":18}}}`, ``)
}

func TestDocumentGotoDefinition3(t *testing.T) {
	testHelperWithTestingString(t, `{"jsonrpc":"2.0","id":6,"method":"textDocument/definition","params":{"textDocument":{"uri":"file:///home/peter/test.swamp"},"position":{"line":13,"character":30}}}`, ``)
}

func TestDocumentEditingRange(t *testing.T) {
	testHelperWithTestingString(t, `{"jsonrpc":"2.0","id":2,"method":"textDocument/linkedEditingRange","params":{"textDocument":{"uri":"file:///home/peter/own/hackman/swamp/gameplay/Main.swamp"},"position":{"line":0,"character":0}}}`, ``)
}

func TestDocumentCodeLens(t *testing.T) {
	testHelperWithTestingStringDoc(t, "file:///home/peter/own/hackman/swamp/gameplay/Main.swamp",
		`{"jsonrpc":"2.0","id":1,"method":"textDocument/codeLens","params":{"textDocument":{"uri":"file:///home/peter/own/hackman/swamp/gameplay/Main.swamp"}}}`, ``)
}

func TestDocumentSemanticSymbols2(t *testing.T) {
	testHelperWithTestingStringDoc(t, "file:///home/peter/own/hackman/swamp/gameplay/Main.swamp",
		`{"jsonrpc":"2.0","id":2,"method":"textDocument/semanticTokens/full","params":{"textDocument":{"uri":"file:///home/peter/own/hackman/swamp/gameplay/Main.swamp"}}}`, ``)
}

func TestDocumentSemanticSymbols4(t *testing.T) {
	testHelperWithTestingStringDoc(t, "file:///home/peter/own/hackman/swamp/gameplayshared/Main.swamp",
		`{"jsonrpc":"2.0","id":2,"method":"textDocument/semanticTokens/full","params":{"textDocument":{"uri":"file:///home/peter/own/hackman/swamp/gameplayshared/Main.swamp"}}}`, ``)
}

func TestLinkedEditingRange(t *testing.T) {
	testHelperWithTestingStringDoc(t, "file:///home/peter/own/hackman/swamp/gameplayshared/Main.swamp",
		`{"jsonrpc":"2.0","id":1,"method":"textDocument/linkedEditingRange","params":{"textDocument":{"uri":"file:///home/peter/own/hackman/swamp/gameplayshared/Main.swamp"},"position":{"line":0,"character":0}}}`, ``)
}

func TestDocumentSemanticSymbols3(t *testing.T) {
	testHelperWithTestingStringDoc(t, "file:///home/peter/own/hackman/swamp/gameplay/MazeLoad.swamp",
		`{"jsonrpc":"2.0","id":2,"method":"textDocument/semanticTokens/full","params":{"textDocument":{"uri":"file:///home/peter/own/hackman/swamp/gameplay/MazeLoad.swamp"}}}`, ``)
}

func TestDocumentSemanticSymbols5(t *testing.T) {
	testHelperWithTestingStringDoc(t, "file:///home/peter/own/turmoil/src/local_packages/turmoil/Collide2.swamp",
		`{"jsonrpc":"2.0","id":2,"method":"textDocument/semanticTokens/full","params":{"textDocument":{"uri":"file:///home/peter/own/turmoil/src/local_packages/turmoil/Collide2.swamp"}}}`, ``)
}
