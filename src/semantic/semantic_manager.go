/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package semantic

import (
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"log"
)

func generateNodesHelper(allTokens []decorated.TypeOrToken) (*SemanticBuilder, error) {
	builder := NewSemanticBuilder()
	for _, foundToken := range allTokens {
		if err := addSemanticToken(foundToken, builder); err != nil {
			log.Printf("was problem with token %v %v", foundToken, err)
			return nil, err
		}
	}
	return builder, nil
}

func GenerateTokensEncodedValues(allTokens []decorated.TypeOrToken) ([]uint, error) {
	builder, err := generateNodesHelper(allTokens)
	if err != nil {
		return nil, err
	}

	return builder.EncodedValues(), nil
}

func GenerateTokensNodes(allTokens []decorated.TypeOrToken) ([]SemanticNode, error) {
	builder, err := generateNodesHelper(allTokens)
	if err != nil {
		return nil, err
	}

	return builder.nodes, nil
}
