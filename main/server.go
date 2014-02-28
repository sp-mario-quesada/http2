package main

import (
	"crypto/tls"
	"flag"
	"github.com/jxck/http2"
	"github.com/jxck/logger"
	"io"
	"log"
	"net/http"
	"os"
)

var verbose bool
var loglevel int

func init() {
	log.SetFlags(log.Lshortfile)
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.IntVar(&loglevel, "l", 0, "log level (1 ERR, 2 WARNING, 3 INFO, 4 DEBUG)")
	f.Parse(os.Args[1:])
	for 0 < f.NArg() {
		f.Parse(f.Args()[1:])
	}
	logger.LogLevel(loglevel)
}

type Hello struct{}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func main() {
	// params
	addr := ":" + os.Args[1]
	cert := "keys/cert.pem"
	key := "keys/key.pem"

	var handler http.Handler = &Hello{}
	// var handler http.Handler = http.FileServer(http.Dir("/tmp"))

	// setup TLS config
	config := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{http2.Version},
	}

	// setup Server
	server := &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		TLSConfig:      config,
		TLSNextProto:   http2.TLSNextProto,
	}

	log.Printf("server starts at localhost%v", addr)
	log.Println(server.ListenAndServeTLS(cert, key))
}
