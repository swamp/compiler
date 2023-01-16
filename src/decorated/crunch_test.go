/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"github.com/swamp/compiler/src/parser"
	"log"
	"testing"

	"github.com/swamp/compiler/src/verbosity"
)

func TestCrunch(t *testing.T) {
	rootModule, mErr := CreateDefaultRootModule(true)
	if parser.IsCompileErr(mErr) {
		t.Fatal(mErr)
	}
	const verboseFlag = verbosity.None
	if verboseFlag >= verbosity.None {
		log.Printf("module\n%v\n", rootModule)
	}
}
