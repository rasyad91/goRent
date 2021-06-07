package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"goRent/internal/config"
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/handler"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
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
	idleTimeout     = 5 * time.Minute  // idleTimeout for server
	shutdownTimeout = 10 * time.Second // shutdown timeout before connections are cancelled
)

var (
	session *scs.SessionManager
	app     *config.AppConfig
)

func init() {
	gob.Register(model.User{})
	gob.Register(model.Product{})
	gob.Register([]model.Product{})
	gob.Register([]model.Rent{})
}

// TODO clean up main() --- Rasyad
func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

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
	app.Session = session
	// session.Cookie.Secure = *inProduction
	app.Info.Printf("Session manager initialized")

	app.Info.Printf("Connecting to AWS elasticsearch client....")
	//initialise AWS elastisearch client
	client, err := newAWSClient()
	if err != nil {
		app.Error.Fatal(err)
	}
	app.AWSClient = client
	app.Info.Printf("AWS elasticsearch client connected")

	app.Info.Printf("Initializing handlers ...")
	r := handler.NewMySQLHandler(db, app)
	handler.New(r)
	helper.New(app)
	render.New(app)
	app.Info.Printf("Handlers initialized...")

	app.Info.Printf("Listen for mail channel...")
	listenForMail()
	defer close(app.MailChan)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%s", *port),
		Handler:     routes(),
		IdleTimeout: idleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

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
