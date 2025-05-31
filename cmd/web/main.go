package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// application struct, for app wide dependencies for webapp
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// default is :4000, value is stored in addr variable
	addr := flag.String("addr", ":4000", "HTTP network address")

	// flag.Parse() to parse cl flag
	flag.Parse()

	// third one is the flags to indicate what additional information to include (local date and time).
	// Note that the flags are joined using the bitwise OR operator |.
	// use the log.Lshortfile flag to include the relevant file name and line number.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom error Log logger in
	// the event of any problems
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// The value returned from the flag.String() function is a pointer to the flag value, not the value itself.
	// Note: any error returned by http.ListenAndServe is always non-nil
	infoLog.Printf("Starting server on port %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

// can save logs inside certain files
// f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
// if err != nil {
// log.Fatal(err)
// }
// defer f.Close()
// infoLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)
