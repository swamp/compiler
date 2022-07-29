package dectype

import "github.com/swamp/compiler/src/decorated/dtype"

func LocalTypesToArgumentNames(types []dtype.Type) []*dtype.TypeArgumentName {
	var argumentNames []*dtype.TypeArgumentName
	for _, genericType := range types {
		localType, wasLocalType := genericType.(*LocalType)
		if wasLocalType {
			argumentNames = append(argumentNames, dtype.NewTypeArgumentName(localType.Identifier().Identifier()))
		}
	}
	return argumentNames
}
