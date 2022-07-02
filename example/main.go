package main

import (
	"io"
	"log"
	"net"
	"net/http"

	"artyom.dev/gofontweb"
)

func main() {
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(gofontweb.FS()))))
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Printf("navigate to http://%s/", ln.Addr().String())
	srv := &http.Server{Handler: mux}
	return srv.Serve(ln)
}

func hello(w http.ResponseWriter, r *http.Request) { io.WriteString(w, indexPage) }

const indexPage = `<!doctype html><title>Go fonts</title>
<link rel=stylesheet href=/assets/go.css>
<style>:root{font-family:Go}code{font-family:"Go Mono"}</style>
<p>You should see this text in Go font.
<ul><li><i>Italic</i>
<li><b>Bold</b></ul>
<p>Some code too: <code>if err != nil {...}</code>.
`
