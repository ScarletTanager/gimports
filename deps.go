package main

import (
	// "go/token"
	"strings"
)

/*
 * This file contains the structures and methods used to store
 * the dependency mappings between packages, files, and
 * function/method calls.
 */

type DepSet struct {
	// Key is the package name
	deps map[string][]*Dependency
}

func NewDepSet() *DepSet {
	return &DepSet{
		make(map[string][]*Dependency),
	}
}

func (d *DepSet) AddImport(pkg, file string) {
	d.Add(&Dependency{
		file: file,
		pkg:  pkg,
	})
}

func (d *DepSet) Add(dep *Dependency) {
	pname := resolveAlias(d, dep.Package())
	dep.pkg = pname
	for _, x := range d.Deps(dep.Package()) {
		if x.File() == dep.File() {
			return
		}
	}

	d.deps[dep.Package()] = append(d.deps[dep.Package()], dep)
}

func (d *DepSet) Deps(pkg string) []*Dependency {
	return d.deps[resolveAlias(d, pkg)]
}

/*
 * Utility function to resolve a qualifier to the name of
 * of an imported package.  E.g. if we encounter "ioutil.ReadDir()",
 * we resolve ioutil to "io/ioutil".
 */
func resolveAlias(d *DepSet, qual string) string {
	suffix := "/" + qual
	for ip := range d.deps {
		if strings.HasSuffix(ip, suffix) {
			return ip
		}
	}

	return qual
}

func (d *DepSet) Calls(pkg, file string) []Call {
	for _, dep := range d.Deps(pkg) {
		if dep.File() == file {
			return dep.Calls()
		}
	}

	return nil
}

func (d *DepSet) Packages() []string {
	var pkgs []string
	for k := range d.deps {
		pkgs = append(pkgs, k)
	}
	return pkgs
}

/*
 * Merges ds into d.  Does not deduplicate.
 */
func (d *DepSet) Merge(ds *DepSet) {
	/* For now, we'll use the slow but easy way */
	for _, pkg := range ds.Packages() {
		for _, dep := range ds.Deps(pkg) {
			d.Add(dep)
		}
	}
}

/*
 * Adds a call to the list for the Dependency between
 * the containing file and the call's qualifier (presumably
 * the package).  Does not check for uniqueness (yet).
 */
func (d *DepSet) AddCall(c Call, file string) {
	for _, dep := range d.Deps(c.Qual) {
		if dep.File() == file {
			dep.Add(c)
			return
		}
	}

	/*
	 * If we're here, we didn't find the appropriate
	 * Dependency, so let's create it.
	 */
	newDep := &Dependency{
		file: file,
		pkg:  c.Qual,
	}

	newDep.Add(c)

	d.Add(newDep)
}

/*
 * AddCall(), but if checkImports == true, then we assume that
 * we've already populated our imports and only add the call if
 * the qualifier (or the package it aliases) exists as a key.
 */
func (d *DepSet) AddPackageCall(c Call, file string, checkImports bool) {
	if checkImports {
		if d.deps[resolveAlias(d, c.Qual)] == nil {
			return
		}
	}
	d.AddCall(c, file)
}

/*
 * Each Dependency represents the unitary dependency from a file
 * to a package.
 */
type Dependency struct {
	file  string
	pkg   string
	calls []Call
}

func (d *Dependency) Add(c Call) {
	for _, call := range d.Calls() {
		if (call.Qual == c.Qual) &&
			(call.Sel == c.Sel) &&
			(call.Line == c.Line) {
			return
		}
	}
	// This is an easy place for a bug to creep in...
	d.calls = append(d.Calls(), c)
}

func (d *Dependency) Calls() []Call {
	return d.calls
}

func (d *Dependency) File() string {
	return d.file
}

func (d *Dependency) Package() string {
	return d.pkg
}

/*
 * Each actual invocation of a function/method/identifier
 * from the package within a file is represented by a Call.
 */
type Call struct {
	Qual string
	Sel  string
	Line int
}
