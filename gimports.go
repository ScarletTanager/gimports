package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

var showCalls bool // List calls into the target package?
var debug bool

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

func init() {
	flag.BoolVar(&showCalls, "calls", false, "lists calls into the target package")
	flag.BoolVar(&debug, "d", false, "enable debug mode")
}

/*
 * MAIN
 */
func main() {
	flag.Parse()
	var pkgMatch string
	dset := NewDepSet()
	if len(flag.Args()) < 1 {
		fmt.Println("Usage: gimports [-calls] <dirname> [<package>]")
		os.Exit(1)
	}
	dirname := flag.Arg(0)
	if len(flag.Args()) > 1 {
		pkgMatch = flag.Arg(1)
	}

	// imps := make(map[string][]PkgUse)

	dirlist := getGoDirs(dirname)

	for _, d := range dirlist {
		dset.Merge(getImports(d))
	}

	// for p, fs := range imps {
	var packages sort.StringSlice = dset.Packages()

	for _, pkg := range packages {
		if pkgMatch != "" {
			if pkg != pkgMatch {
				continue
			}
		}
		deps := dset.Deps(pkg)
		fmt.Printf("IMPORTERS of %s (%d results):\n", pkg, len(deps))
		for _, dep := range deps {
			fmt.Printf("\t%s\n", dep.File())
			if showCalls {
				for _, c := range dep.Calls() {
					// locfields := strings.Split(fset.Position(c.Location).String(), ":")
					fmt.Printf("\t\tLine: %d\t%s.%s\n",
						// locfields[1],
						c.Line,
						c.Qual,
						c.Sel)
				}
			}
		}
	}
	// 	fmt.Printf("IMPORTERS OF %s (%d results):\n", p, len(fs))
	// 	for _, f := range fs {
	// 		fmt.Printf("\t%s\n", f)
	// 	}
	// }
}
