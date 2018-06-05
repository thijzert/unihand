Unihand is a standalone server application for character information in the Unicode Han Database, published by the Unicode Consortium.

Links:

* https://www.unicode.org/reports/tr38/
* https://www.unicode.org/Public/UNIDATA/Unihan.zip

Building
--------
Compiling this application requires a recent version of the Go programming language. With it, run:

    go build cmd/server/server.go

Usage
-----
Download the Unihan.zip file to your local system. After compiling, run a local server with this command:

    cmd/server/server -zip /path/to/Unihan.zip -listen :8978

This will open a local HTTP server on port 8978.
Query it by sending it HTTP requests like this:

    curl http://localhost:8978/character/91CF

This will provide information on Unicode character U+91CF.

License
-------
UNIHANd and its source code are available under the terms of the BSD 3-clause license. Find out what that means here: https://www.tldrlegal.com/l/bsd3

The Unihan dataset is published by the Unicode Consortium under certain terms and conditions. Read about them here: https://www.unicode.org/copyright.html

