package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

var fset *token.FileSet

func getImports(dir string) *DepSet {
	var pf *ast.File
	var err error
	fset = token.NewFileSet()
	files := getGoFiles(dir)
	dset := NewDepSet()

	for _, file := range files {
		if showCalls == true {
			pf, err = parser.ParseFile(fset, file, nil, 0)
			// Build a map with the imports mapping to the files
			for _, s := range pf.Imports {
				pname := strings.Trim(s.Path.Value, "\"")
				if debug {
					fmt.Printf("Import found: %s imports %s\n", file, pname)
				}
				dset.AddImport(pname, file)
			}

			if debug {
				for _, pkg := range dset.Packages() {
					fmt.Printf("Imported: %s\n", pkg)
				}
			}

			// ast.Inspect(pf, parseCalls)
			ast.Walk(finder{
				find: findCalls,
				dset: dset,
			}, pf)
		} else {
			pf, err = parser.ParseFile(fset, file, nil, parser.ImportsOnly)
			// Build a map with the imports mapping to the files
			for _, s := range pf.Imports {
				pname := strings.Trim(s.Path.Value, "\"")
				dset.AddImport(pname, file)
			}
		}
		if err != nil {
			// fmt.Println(err)
		}

		/*
		 * We'll need to AND the imports and the keys in the
		 * package calls map to trim out keys from the latter
		 * which are not actually packages.
		 * We also need to figure out how to account for methods -
		 * since the qualifier (X node) in those cases is the
		 * method receiver, not the package.
		 */
	}

	return dset
}

func getGoFiles(dir string) []string {
	var flist []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".go") == true {
			fullpath := filepath.Join(dir, f.Name())
			flist = append(flist, fullpath)
		}
	}

	return flist
}

/*
 * Below here is the plumbing to use with ast.Walk()
 */

/*
 * findCalls() parses the node passed in and, if it's a
 * SelectorExpr (which any call to a function/method will be)
 * populates the map passed in as the second argument.
 * This is the function used for finder.find().
 *
 * Functions will have the package as the X node,
 * whereas methods will have the object identifier as the
 * X.  We do not differentiate in this function (there
 * really isn't a good way without access to the list of
 * imports and identifiers, so that has to be done higher up
 * the stack).
 */
func findCalls(n ast.Node, f *finder) bool {
	switch x := n.(type) {
	case *ast.File:
		locfields := strings.Split(fset.Position(n.Pos()).String(), ":")
		f.currentFile = locfields[0]
	case *ast.SelectorExpr:
		switch y := x.X.(type) {
		case *ast.Ident:
			locfields := strings.Split(fset.Position(y.NamePos).String(), ":")
			ln, err := strconv.Atoi(locfields[1])
			if err != nil {
				ln = -1
			}
			f.dset.AddPackageCall(Call{
				Qual: y.Name,
				Sel:  x.Sel.Name,
				// Location: y.NamePos,
				Line: ln,
			}, f.currentFile, true)
		default:
			return true
		} // END switch y :=
	default:
		return true
	} // END switch x :=
	return true
}

/*
 * finder is a struct which implements the ast.Visitor interface.
 */
type finder struct {
	find        func(ast.Node, *finder) bool
	dset        *DepSet
	currentFile string
}

// func (f *finder) Calls() []Call {
// return f.calls
// }

func (f *finder) DepSet() *DepSet {
	return f.dset
}

func (f finder) Visit(node ast.Node) ast.Visitor {
	if f.find(node, &f) {
		return f
	}
	return nil
}
