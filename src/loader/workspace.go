package loader

import (
	"log"

	decorated "github.com/swamp/compiler/src/decorated/expression"
)

type Workspace struct {
	rootDirectory LocalFileSystemRoot
	projects      map[string]Project
}

func NewWorkspace(rootDirectory LocalFileSystemRoot) *Workspace {
	return &Workspace{rootDirectory: rootDirectory, projects: make(map[string]Project)}
}

func (w *Workspace) FindProjectFromRootDirectory(root LocalFileSystemRoot) Project {
	foundProject := w.projects[string(root)]

	return foundProject
}

func (w *Workspace) AddProject(root LocalFileSystemRoot, project Project) {
	w.projects[string(root)] = project
}

func (w *Workspace) AddPackage(p *Package) {
	log.Printf("adding package %v", p.root)
	w.AddProject(p.Root(), p)
}

func (w *Workspace) AllPackages() []*Package {
	var packages []*Package
	for _, project := range w.projects {
		foundPackage, wasPackage := project.(*Package)
		if wasPackage {
			packages = append(packages, foundPackage)
		}
	}
	return packages
}

func (w *Workspace) AddOrReplacePackage(p *Package) {
	existingPackage := w.FindProject(p.root)
	if existingPackage != nil {
		log.Printf("remove for overwrite package %v", p.Root())
		delete(w.projects, string(p.root))
	}

	w.AddPackage(p)
}

func (w *Workspace) FindProject(root LocalFileSystemRoot) *Package {
	p := w.FindProjectFromRootDirectory(root)
	if p != nil {
		foundPackage, wasPackage := p.(*Package)
		if wasPackage {
			return foundPackage
		}
	}

	return nil
}

func (w *Workspace) FindModuleFromSourceFile(path LocalFileSystemPath) (*decorated.Module, *Package) {
	for _, foundPackage := range w.AllPackages() {
		foundModule := foundPackage.FindModuleFromAbsoluteFilePath(path)
		if foundModule != nil {
			return foundModule, foundPackage
		}
	}

	return nil, nil
}
