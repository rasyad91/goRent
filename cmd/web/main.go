package main

import (
	"context"
	"fmt"
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/handler"
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

	dbPort      = "3306"
	dbName      = "goRent"
	dbUser      = "root"
	dbPassword  = "" //insert own password
	dbHost      = "127.0.0.1"
	dbParseTime = true
)

var (
	session *scs.SessionManager
)

func main() {

	f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		println("fail to open log file: ", err)
	}
	defer f.Close()

	// Customized loggers
	infoLog := log.New(io.MultiWriter(f, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(io.MultiWriter(f, os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbParseTime,
	)

	infoLog.Printf("Connecting to DB: %s...\n", dsn)
	db, err := mysqlDriver.Connect(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Printf("Successfully connected to DB: %s\n", dsn)
	defer db.SQL.Close()

	// session
	log.Printf("Initializing session manager....")
	session = scs.New()
	session.Store = mysqlstore.New(db.SQL)
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	// session.Cookie.Name = fmt.Sprintf("gbsession_id_%s", *identifier)
	session.Cookie.SameSite = http.SameSiteLaxMode
	// session.Cookie.Secure = *inProduction

	r := handler.NewMySQLHandler(db, errorLog, infoLog)
	handler.New(r)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		Handler:     routes(),
		IdleTimeout: idleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// start server on seperate thread
	go func() {
		infoLog.Printf("Listening on port:%d\n", port)
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	// blocks code, waits for stop to initiate
	<-stop

	infoLog.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		errorLog.Fatalln(err)
	}
	infoLog.Println("Server shut down")
}
