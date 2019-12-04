package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql" // blank identifier alias, underscore stops compiler throwing and error
	"github.com/kwhitlock/lets-go-book/pkg/models/mysql"
)

// application struct
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	// setup a channel to read from, of size 1
	killSignal := make(chan os.Signal, 1)
	// use os.signal.Notify to send a notification based on the type of os signal.
	signal.Notify(killSignal, os.Interrupt)

	addr := flag.String("addr", ":4000", "HTTP network address")
	staticDir := flag.String("static", "./ui/static/", "Directory where static files are located.")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// init a new template cache
	infoLog.Printf("Initializing a new template in memory cache")
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// initialise a new instance of the application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
	}

	// Initialize a new http.Server struct, so we can set a custom logger
	// for error log handling.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(*staticDir),
	}

	// anonymous function
	go func() {
		infoLog.Printf("Starting server on %s", *addr)
		err = srv.ListenAndServe()
		errorLog.Fatal(err)
	}()

	// read off the killSignal channel
	<-killSignal
	fmt.Println("Thanks for using the app.")
	// time.Sleep(30 * time.Second)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
