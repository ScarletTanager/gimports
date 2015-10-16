package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/*
 * Add a dir to the list of dirs known to contain *.go files
 * Will extend the slice if add would exceed capacity
 */
func addDirs(to, from []string) []string {
	m := len(to)
	n := m + len(from)

	// if necessary, reallocate
	if n > cap(to) {
		newSlice := make([]string, n)
		copy(newSlice, to)
		to = newSlice
	}
	to = to[0:n]
	copy(to[m:n], from)
	return to
}

/*
 * Returns a list of directories containing *.go files
 */
func getGoDirs(dirname string) []string {
	isGoDir := false // does the current dir directly contain .go source files?
	var dlist []string

	/* Get the file info */
	fi, err := os.Lstat(dirname)
	if err != nil {
		return nil
	}

	if fi.IsDir() == true { /* Is a directory */
		// dlist = make([]string, 1)
		files, err := ioutil.ReadDir(dirname) // Get contents of the current dir
		if err != nil {
			return nil
		}

		/*
		 * Process contents of current directory
		 */
		for _, f := range files {
			fullpath := dirname + "/" + f.Name()
			if f.IsDir() == true { // If the file is a dir, then process that...
				dlist = addDirs(dlist, getGoDirs(fullpath))
			} else {
				if strings.HasSuffix(f.Name(), ".go") == true {
					isGoDir = true
				}
			}
		}

		if isGoDir == true {
			dlist = addDirs(dlist, []string{dirname})
		}
	} /* End of is a directory */

	return dlist
}

func main() {
	var pkgMatch string
	if len(os.Args) < 2 {
		fmt.Println("Usage: gimports <dirname> [<package>]")
		os.Exit(1)
	}
	dirname := os.Args[1]
	if len(os.Args) > 2 {
		pkgMatch = "\"" + os.Args[2] + "\""
	}

	imps := make(map[string][]string)

	dirlist := getGoDirs(dirname)
	for _, d := range dirlist {
		for pkg, files := range getImports(d) {
			imps[pkg] = append(imps[pkg], files...)
		}
	}

	for p, fs := range imps {
		if pkgMatch != "" {
			if p != pkgMatch {
				continue
			}
		}
		fmt.Printf("IMPORTERS OF %s (%d results):\n", p, len(fs))
		for _, f := range fs {
			fmt.Printf("\t%s\n", f)
		}
	}
}
