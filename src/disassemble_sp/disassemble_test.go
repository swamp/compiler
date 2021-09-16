/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package swampdisasm_sp

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSomething(t *testing.T) {
	s := "030501010003060201000d030506030701010103080201010e04070801000203040a"

	octets, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	stringLines := Disassemble(octets)
	output := fmt.Sprintf("%v", stringLines)

	const expectedOutput = `[00: get 5, 1, [#0] 05: get 6, 2, [#0] 0a: add 3,5,6 0e: get 7, 1, [#1] 13: get 8, 2, [#1] 18: sub 4,7,8 1c: crs 0 [3 4] 21: ret]`

	fmt.Println(output)

	if output != expectedOutput {
		t.Errorf("disassemble produced wrong output. expected\n%s\nbut received\n%s\n", expectedOutput, output)
	}
}
