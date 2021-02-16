/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package deccy

import (
	"fmt"
	"testing"
)

func TestCrunch(t *testing.T) {
	aliasModules, importModules, mErr := CreateDefaultRootModule(true)
	if mErr != nil {
		t.Fatal(mErr)
	}
	const verboseFlag bool = true
	if verboseFlag {
		fmt.Printf("module\n%v\n%v\n", aliasModules, importModules)
	}
}
