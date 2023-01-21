/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package semantic

import (
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/token"
	"log"
)

func generateNodesLeaf(expanded *decorated.ExpandedNode, builder *SemanticBuilder) error {
	for _, childNode := range expanded.Children() {
		if childNode.HasChildren() {
			generateNodesLeaf(childNode, builder)
		} else {
			if err := addSemanticToken(childNode.TypeOrToken(), builder); err != nil {
				log.Printf("was problem with token %v %v", childNode, err)
				return err
			}
		}
	}

	return nil
}

func generateNodesHelper(expandedRootNodes []*decorated.ExpandedNode, debugDocument *token.SourceFileDocument) (*SemanticBuilder, error) {
	builder := NewSemanticBuilder(debugDocument)
	for _, foundToken := range expandedRootNodes {
		if err := generateNodesLeaf(foundToken, builder); err != nil {
			return nil, err
		}
	}
	return builder, nil
}

func GenerateTokensEncodedValues(allTokens []*decorated.ExpandedNode, debugDocument *token.SourceFileDocument) ([]uint, error) {
	builder, err := generateNodesHelper(allTokens, debugDocument)
	if err != nil {
		return nil, err
	}

	return builder.EncodedValues(), nil
}

func GenerateTokensNodes(allTokens []*decorated.ExpandedNode, debugDocument *token.SourceFileDocument) ([]SemanticNode, error) {
	builder, err := generateNodesHelper(allTokens, debugDocument)
	if err != nil {
		return nil, err
	}

	return builder.nodes, nil
}
