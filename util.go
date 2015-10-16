package main

import (
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type SourceFile struct {
	Name    string // Short name of file
	Path    string // Full path relative to the starting directory
	Imports []string
	Package string
}

func getImports(dir string) map[string][]string {
	fset := token.NewFileSet()
	files := getGoFiles(dir)
	imps := make(map[string][]string)
	for _, file := range files {
		// Parse the file containing this very example
		// but stop after processing the imports.
		pf, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
		if err != nil {
			// fmt.Println(err)
		}

		// Build a map with the imports mapping to the files
		for _, s := range pf.Imports {
			imps[s.Path.Value] = append(imps[s.Path.Value], file)
		}
	}

	return imps
}

/*
 * Populates the list of import paths for the specified file
 */
func calculateImports(path string) {

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
