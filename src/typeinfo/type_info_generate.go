/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package typeinfo

import (
	"log"

	"github.com/swamp/compiler/src/decorated/dtype"
	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/loader"
	"github.com/swamp/compiler/src/resourceid"
	"github.com/swamp/compiler/src/verbosity"
)

type TypeLookup interface {
	Lookup(d dtype.Type) (int, error)
}

func generateModuleToChunk(module *decorated.Module, chunk *Chunk, verboseFlag verbosity.Verbosity) error {
	for _, exposedDef := range module.LocalDefinitions().Definitions() {
		exposedType := exposedDef.Expression().Type()
		if verboseFlag >= verbosity.Low {
			log.Printf("generateModuleToChunkTypeInfo definition: %s\n", exposedType.String())
		}
		convertedType, err := chunk.ConsumeType(exposedType)
		if err != nil {
			continue
		}
		if verboseFlag >= verbosity.Low {
			log.Printf("converted: %v\n", convertedType)
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

func GenerateModule(module *decorated.Module) ([]byte, TypeLookup, resourceid.ResourceNameLookup, error) {
	const verboseFlag = verbosity.None

	chunk := &Chunk{}
	resourceLookup := resourceid.NewResourceNameLookupImpl()

	if err := generateModuleToChunk(module, chunk, verboseFlag); err != nil {
		return nil, nil, nil, err
	}

	octets, octetsErr := ChunkToOctets(chunk)
	if octetsErr != nil {
		return nil, nil, nil, octetsErr
	}

	return octets, chunk, resourceLookup, nil
}

func GeneratePackageToChunk(world *loader.Package, chunk *Chunk) error {
	const verboseFlag = verbosity.None
	for _, module := range world.AllModules() {
		if err := generateModuleToChunk(module, chunk, verboseFlag); err != nil {
			return err
		}
	}

	return nil
}
