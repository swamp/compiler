/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package typeinfo

import (
	"fmt"

	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/loader"
)

type TypeLookup interface {
	Lookup(d dtype.Type) (int, error)
}

func generateModuleToChunk(module *decorated.Module, chunk *Chunk, verboseFlag bool) error {
	for _, exposedDef := range module.Definitions().Definitions() {
		exposedType := exposedDef.Expression().Type()
		if verboseFlag {
			fmt.Printf("generateModuleToChunkTypeInfo definition: %s\n", exposedType.DecoratedName())
		}
		convertedType, err := chunk.ConsumeType(exposedType)
		if err != nil {
			continue
		}
		if verboseFlag {
			fmt.Printf("converted: %v\n", convertedType)
		}
	}

	return nil
}

func ChunkToOctets(chunk *Chunk) ([]byte, error) {
	octets, serializeErr := SerializeToOctets(chunk)
	if serializeErr != nil {
		return nil, serializeErr
	}

	return octets, nil
}

func GenerateModule(module *decorated.Module) ([]byte, TypeLookup, error) {
	const verboseFlag = true

	chunk := &Chunk{}

	if err := generateModuleToChunk(module, chunk, verboseFlag); err != nil {
		return nil, nil, err
	}

	octets, octetsErr := ChunkToOctets(chunk)
	if octetsErr != nil {
		return nil, nil, octetsErr
	}

	return octets, chunk, nil
}

func Generate(world *loader.World) (*Chunk, error) {
	const verboseFlag = false

	chunk := &Chunk{}

	for _, module := range world.AllModules() {
		if err := generateModuleToChunk(module, chunk, verboseFlag); err != nil {
			return nil, err
		}
	}


	return chunk, nil
}
