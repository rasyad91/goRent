package main

import (
	"context"
	"flag"
	"fmt"
	"goRent/internal/config"
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/handler"
	"goRent/internal/helper"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

const (
	port            = 8080             //http port to listen to
	idleTimeout     = 5 * time.Minute  // idleTimeout for server
	shutdownTimeout = 10 * time.Second // shutdown timeout before connections are cancelled
)

var (
	session *scs.SessionManager
	app     *config.AppConfig
)

// TODO clean up main() --- Rasyad
func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	port := flag.String("port", "8080", "server port to listen on")
	identifier := flag.String("identifier", "GoRent", "unique identifier")
	domain := flag.String("domain", "localhost", "domain name (e.g. example.com)")
	dbUser := flag.String("dbuser", "", "database user")
	dbPassword := flag.String("dbpassword", "", "database password")
	dbName := flag.String("dbname", "", "database name")
	dbHost := flag.String("dbhost", "", "database host")
	dbPort := flag.String("dbport", "", "database port")
	dbParseTime := flag.Bool("dbparsetime", true, "database parse time option")
	flag.Parse()

	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		println("fail to open log file: ", err)
	}
	defer f.Close()

	app = &config.AppConfig{}
	app.Domain = *domain

	// Customized loggers
	app.Info = log.New(io.MultiWriter(f, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	app.Error = log.New(io.MultiWriter(f, os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		*dbUser, *dbPassword, *dbHost, *dbPort, *dbName, *dbParseTime,
	)

	app.Info.Printf("Connecting to DB: %s...\n", dsn)
	db, err := mysqlDriver.Connect(dsn)
	if err != nil {
		app.Error.Fatal(err)
	}
	app.Info.Printf("Successfully connected to DB: %s\n", dsn)
	defer db.SQL.Close()

	// session
	app.Info.Printf("Initializing session manager....")
	session = scs.New()
	session.Store = mysqlstore.New(db.SQL)
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Name = fmt.Sprintf("gbsession_id_%s", *identifier)
	session.Cookie.SameSite = http.SameSiteLaxMode
	// session.Cookie.Secure = *inProduction

	r := handler.NewMySQLHandler(db, app)
	handler.New(r)

	helper.NewHelpers(app)

	listenForMail()
	defer close(app.MailChan)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%s", *port),
		Handler:     routes(),
		IdleTimeout: idleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start server on seperate thread
	go func() {
		app.Info.Printf("Listening on port:%s\n", *port)
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	// blocks code, waits for stop to initiate
	<-stop

	app.Info.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		app.Error.Fatalln(err)
	}
	app.Info.Println("Server shut down")
}
