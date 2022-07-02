// Package gofontweb provides embedded [Go fonts] in woff2 format.
//
// # Usage
//
// You can plug it into your web server code like this:
//
//	mux := http.NewServeMux()
//	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(gofontweb.FS()))))
//
// Then in your HTML link to the “go.css” file to set up all the supported font faces:
//
//	<link rel=stylesheet href=/assets/go.css>
//
// This stylesheet declares the following font families:
//
//   - Go
//   - Go Medium
//   - Go Smallcaps
//   - Go Mono
//
// You can find the full example in “example” subdirectory.
//
// [Go fonts]: https://go.dev/blog/go-fonts
package gofontweb

import (
	"io/fs"

	"artyom.dev/gofontweb/internal/data"
)

// FS returns an embedded file system holding all the woff2 files.
// The “go.css” file in this filesystem configures supported font faces.
func FS() fs.FS { return data.FS }

//go:generate go run gen.go
