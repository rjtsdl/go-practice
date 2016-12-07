The Go language sample doesn't differentiate the underlining platform, it relies on Go's cross platform compiling ability to achieve that. The sample depends on the mssql driver "github.com/denisenkom/go-mssqldb".

Here is the instruction to run the sample

* Properly installed Go in your operating system
* Properly set up GOPATH
* Run 'go get github.com/denisenkom/go-mssqldb' for once
* Properly set the connection string in file sql-sample.go LN #13
* Run 'go run sql-sample.go'

Cheers
