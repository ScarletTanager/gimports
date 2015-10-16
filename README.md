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
