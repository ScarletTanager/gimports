# gimports
gimports is a simple tool to operate as something of an inverse of "go list".  It reports by package which source files in a dir tree import that package.

For example, if you want to know which go source files under a directory tree import the `errors` package from the standard library, you could run:

    $ gimports path/to/src/tree errors
    
You would get output something like this:

    IMPORTERS OF "errors":
	../openssl/examples/evpexpl/main.go
	../openssl/ssl/httpsclient.go
	../webster/webster.go

If you omit the package name, then `gimports` will report the imports for all files under the tree by package.

For more detailed reporting, you can call `gimports -calls <dir> [<package>]` - this will trigger the per-file reporting of all function calls by package (or for the specific package if named on the command line):

    IMPORTERS of net/http (2 results):
    	../webster/webster.go
    		Line: 25	http.Response
    		Line: 76	http.Response
    		Line: 84	http.Client
    	../webster/webster_types.go
    		Line: 9	http.Client
