package doc

import (
	"sort"

	decorated "github.com/swamp/compiler/src/decorated/expression"
	"github.com/swamp/compiler/src/loader"
)

type DocPackage struct {
	foundPackage       *loader.Package
	environmentModules []*decorated.Module
	sharedModules      []*decorated.Module
	modules            []*decorated.Module
}

type DocRoot struct {
	packages []*DocPackage
}

func FilterOutDocRoot(packages []*loader.Package) *DocRoot {
	documentedModules := make(map[string]*decorated.Module)
	packageMap := make(map[string]*loader.Package)
	packageKeys := make([]string, 0, len(packages))
	for _, foundPackage := range packages {
		packageMap[foundPackage.Name()] = foundPackage
	}

	sort.Strings(packageKeys)

	documentRoot := &DocRoot{}
	for _, compiledPackage := range packages {
		docPackage := &DocPackage{
			foundPackage: compiledPackage,
		}

		moduleKeysInPackage := make([]string, 0, len(compiledPackage.AllModules()))
		modulesInPackageMap := make(map[string]*decorated.Module)
		for _, foundModule := range compiledPackage.AllModules() {
			key := foundModule.FullyQualifiedModuleName().String()
			if _, alreadyDocumented := documentedModules[key]; alreadyDocumented {
				continue
			}

			filteredTypes := filterTypes(foundModule.LocalTypes().AllTypes())

			filteredFunctions, filteredConstants := filterDefinitions(foundModule.LocalDefinitions().Definitions())
			if len(filteredFunctions) == 0 && len(filteredConstants) == 0 && len(filteredTypes) == 0 {
				continue
			}

			modulesInPackageMap[key] = foundModule
			moduleKeysInPackage = append(moduleKeysInPackage, key)
			documentedModules[key] = foundModule
		}

		sort.Strings(moduleKeysInPackage)

		for _, moduleKeyInPackage := range moduleKeysInPackage {
			moduleToDocument := modulesInPackageMap[moduleKeyInPackage]

			switch moduleToDocument.ModuleType() {
			case decorated.ModuleTypeNormal:
				docPackage.modules = append(docPackage.modules, moduleToDocument)
			case decorated.ModuleTypeFromPath:
				docPackage.sharedModules = append(docPackage.sharedModules, moduleToDocument)
			case decorated.ModuleTypeFromEnvironment:
				docPackage.environmentModules = append(docPackage.environmentModules, moduleToDocument)
			}
		}

		if len(moduleKeysInPackage) > 0 {
			documentRoot.packages = append(documentRoot.packages, docPackage)
		}
	}

	return documentRoot
}
