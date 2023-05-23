/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"log"
	"testing"

	"github.com/swamp/compiler/src/parser"

	"github.com/swamp/compiler/src/verbosity"
)

func xTestCrunch(t *testing.T) {
	rootModule, mErr := CreateDefaultRootModule(true)
	if parser.IsCompileErr(mErr) {
		t.Fatal(mErr)
	}
	const verboseFlag = verbosity.None
	if verboseFlag >= verbosity.None {
		log.Printf("module\n%v\n", rootModule)
	}
}
