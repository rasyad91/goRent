package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"goRent/internal/config"
	"goRent/internal/driver"
	"goRent/internal/handler"
	"goRent/internal/helper"
	"goRent/internal/model"
	"goRent/internal/render"
	"goRent/internal/repository/mysql"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
)

const (
	idleTimeout       = 1 * time.Minute // idleTimeout for server
	readTimeout       = 3 * time.Second // readTimeout for server
	writeTimeout      = 3 * time.Second // writeTimeout for server
	readHeaderTimeout = 3 * time.Second // readHeaderTimeout for server

	sessionLifetimeTimeout = 24 * time.Hour // readHeaderTimeout for server

	shutdownTimeout = 5 * time.Second // shutdown timeout before connections are cancelled
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

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	app = &config.AppConfig{}
	app.Domain = *domain

	f, err := setLogs()
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t",
		*dbUser, *dbPassword, *dbHost, *dbPort, *dbName, *dbParseTime,
	)

	app.Info.Printf("Connecting to DB: %s...\n", dsn)
	db, err := driver.Connect(dsn, *dbDialect)
	if err != nil {
		app.Error.Fatal(err)
	}
	app.Info.Printf("Successfully connected to DB: %s\n", dsn)
	defer db.Close()

	// session
	setSession(db)
	// aws
	if err := setAWS(); err != nil {
		app.Error.Fatal(err)
	}

	app.Info.Printf("Initializing handlers ...")
	dbRepo := mysql.NewRepo(db)
	handlerRepo := handler.NewRepo(dbRepo, app)

	handler.New(handlerRepo)
	helper.New(app)
	render.New(app)
	app.Info.Printf("Handlers initialized...")

	app.MailChan = make(chan model.MailData)
	app.Info.Printf("Listen for mail channel...")
	listenForMail()
	defer close(app.MailChan)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", *port),
		Handler:           routes(),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
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
	close(stop)

	app.Info.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		app.Error.Fatalln(err)
	}
	app.Info.Println("Server shut down")
}

func setLogs() (io.WriteCloser, error) {
	abs, err := filepath.Abs("./log/server.log")
	if err != nil {
		return nil, fmt.Errorf("fail to get absolute path log file: %v", err)
	}
	f, err := os.OpenFile(abs, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("fail to open log file: %v", err)
	}
	// Customized loggers
	app.Info = log.New(io.MultiWriter(f, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	app.Error = log.New(io.MultiWriter(f, os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return f, nil
}

func setSession(db *sql.DB) {
	app.Info.Printf("Initializing session manager....")
	session = scs.New()
	session.Store = mysqlstore.New(db)
	session.Lifetime = sessionLifetimeTimeout
	session.Cookie.Persist = true
	session.Cookie.Name = fmt.Sprintf("gbsession_id_%s", *identifier)
	session.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = session
	app.Info.Printf("Session manager initialized")
}

func setAWS() error {
	app.Info.Printf("Connecting to AWS elasticsearch client....")
	//initialise AWS elastisearch client
	client, err := newAWSClient()
	if err != nil {
		return err
	}
	app.AWSClient = client
	app.Info.Printf("AWS elasticsearch client connected")

	app.Info.Printf("Connecting to AWS S3 client session....")
	awsS3Session, err := NewAWSSession()
	if err != nil {
		return err
	}
	app.AWSS3Session = awsS3Session
	app.Info.Printf("AWS S3 Session Established")
	return nil
}
