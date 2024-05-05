package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type app struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	Handler  http.Handler
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	app := &app{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server of %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
