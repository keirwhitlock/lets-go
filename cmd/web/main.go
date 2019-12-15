package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql" // blank identifier alias, underscore stops compiler throwing and error
	"github.com/golangcollege/sessions"
	"github.com/kwhitlock/lets-go-book/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

// application struct
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	users         *mysql.UserModel
}

func main() {

	// setup a channel to read from, of size 1
	killSignal := make(chan os.Signal, 1)
	// use os.signal.Notify to send a notification based on the type of os signal.
	signal.Notify(killSignal, os.Interrupt)

	addr := flag.String("addr", ":4000", "HTTP network address")
	staticDir := flag.String("static", "./ui/static/", "Directory where static files are located.")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret session key")

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

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	// initialise a new instance of the application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		session:  session,
		snippets: &mysql.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Initialize a new http.Server struct, so we can set a custom logger
	// for error log handling.
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(*staticDir),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// anonymous function
	go func() {
		infoLog.Printf("Starting server on %s", *addr)
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
