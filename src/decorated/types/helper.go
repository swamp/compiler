/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package dectype

import "github.com/swamp/compiler/src/decorated/dtype"

func LocalTypesToArgumentNames(types []*LocalTypeDefinition) []*dtype.LocalTypeName {
	var argumentNames []*dtype.LocalTypeName
	for _, localType := range types {
		argumentNames = append(argumentNames, localType.identifier.LocalTypeName())
	}
	return argumentNames
}
