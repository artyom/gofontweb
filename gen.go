//go:build generate

package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomediumitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"golang.org/x/image/font/gofont/gosmallcapsitalic"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if _, err := exec.LookPath("woff2_compress"); err != nil {
		return fmt.Errorf("%w\nyou can find woff2_compress source at https://github.com/google/woff2", err)
	}
	td, err := os.MkdirTemp("", "gofonts-convert-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(td)
	dstdir := filepath.FromSlash("internal/data")
	if err := os.MkdirAll(dstdir, 0777); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dstdir, "pkg.go"), []byte(pkgFile), 0666); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dstdir, "LICENSE.txt"), []byte(license), 0666); err != nil {
		return err
	}
	cssBuf := bytes.NewBufferString("/* see LICENSE.txt for Copyright details */\n")
	for _, r := range [...]struct {
		name    string
		ttfData []byte
	}{
		{"Go-Regular", goregular.TTF},
		{"Go-Italic", goitalic.TTF},
		{"Go-Medium", gomedium.TTF},
		{"Go-Medium-Italic", gomediumitalic.TTF},
		{"Go-Bold", gobold.TTF},
		{"Go-Bold-Italic", gobolditalic.TTF},
		{"Go-Smallcaps", gosmallcaps.TTF},
		{"Go-Smallcaps-Italic", gosmallcapsitalic.TTF},
		{"Go-Mono", gomono.TTF},
		{"Go-Mono-Italic", gomonoitalic.TTF},
		{"Go-Mono-Bold", gomonobold.TTF},
		{"Go-Mono-Bold-Italic", gomonobolditalic.TTF},
	} {
		if err := convertFont(r.ttfData, r.name, td, dstdir); err != nil {
			return err
		}
		info := struct{ Family, LocalFont, FileName, Style, Weight string }{
			Family: "Go", LocalFont: repl.Replace(r.name), FileName: r.name + ".woff2",
			Style: "normal", Weight: "normal",
		}
		switch {
		case strings.Contains(r.name, "Smallcaps"):
			info.Family = "Go Smallcaps"
		case strings.Contains(r.name, "Medium"):
			info.Family = "Go Medium"
			info.Weight = "500"
		case strings.Contains(r.name, "Mono"):
			info.Family = "Go Mono"
		}
		if strings.Contains(r.name, "Bold") {
			info.Weight = "600"
		}
		if strings.Contains(r.name, "Italic") {
			info.Style = "italic"
		}
		if err := tpl.Execute(cssBuf, info); err != nil {
			return err
		}
	}
	return os.WriteFile(filepath.Join(dstdir, "go.css"), cssBuf.Bytes(), 0666)
}

func convertFont(data []byte, name, workdir, dstdir string) error {
	ttfName := filepath.Join(workdir, name+".ttf")
	if err := os.WriteFile(ttfName, data, 0666); err != nil {
		return err
	}
	defer os.Remove(ttfName)
	cmd := exec.Command("woff2_compress", filepath.Base(ttfName))
	cmd.Dir = workdir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	woffFile, err := os.Open(filepath.Join(workdir, name+".woff2"))
	if err != nil {
		return err
	}
	defer woffFile.Close()
	defer os.Remove(woffFile.Name())
	dst, err := os.Create(filepath.Join(dstdir, name+".woff2"))
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, woffFile); err != nil {
		os.Remove(dst.Name())
		return err
	}
	return dst.Close()
}

var tpl = template.Must(template.New("").Parse(`@font-face {font-family:{{.Family | printf "%q"}};
src:local({{.LocalFont | printf "%q"}}),url({{.FileName}}) format("woff2");
font-style:{{.Style}};font-weight:{{.Weight}};}
`)).Option("missingkey=error")

var repl = strings.NewReplacer("-", " ")

//go:embed LICENSE.txt
var license string

const pkgFile = `// Code generated with gen.go DO NOT EDIT.

package data

import "embed"

//go:embed go.css LICENSE.txt *.woff2
var FS embed.FS
`
