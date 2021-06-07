package main

import (
	"context"
	"encoding/gob"
	"flag"
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

	// "github.com/olivere/elastic"

	"github.com/olivere/env"
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

	port := flag.String("port", "8080", "server port to listen on")
	identifier := flag.String("identifier", "GoRent", "unique identifier")
	domain := flag.String("domain", "localhost", "domain name (e.g. example.com)")
	dbUser := flag.String("dbuser", "", "database user")
	dbPassword := flag.String("dbpassword", "", "database password")
	dbName := flag.String("dbname", "", "database name")
	dbHost := flag.String("dbhost", "", "database host")
	dbPort := flag.String("dbport", "", "database port")
	dbParseTime := flag.Bool("dbparsetime", true, "database parse time option")
	accessKey := flag.String("access-key", env.String("", "AWS_ACCESS_KEY", "AWS_ACCESS_KEY_ID"), "Access Key ID")
	secretKey := flag.String("secret-key", env.String("", "AWS_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"), "Secret access key")
	esUrl := flag.String("esUrl", "", "Elasticsearch URL")
	sniff := flag.Bool("sniff", false, "Enable or disable sniffing")
	region := flag.String("region", "", "AWS Region name")

	flag.Parse()
	log.SetFlags(0)

	f, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		println("fail to open log file: ", err)
	}
	defer f.Close()

	app = &config.AppConfig{}
	app.Domain = *domain

	app.AwsAccessKey = *accessKey
	app.AwsSecretKey = *secretKey
	app.AwsUrl = *esUrl
	app.AwsSniff = *sniff
	app.AwsRegion = *region

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

	r := handler.NewMySQLHandler(db, app)
	handler.New(r)

	helper.New(app)
	render.New(app)

	listenForMail()
	defer close(app.MailChan)

	server := &http.Server{
		Addr:        fmt.Sprintf(":%s", *port),
		Handler:     routes(),
		IdleTimeout: idleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	//start of elastic search codes
	if *esUrl == "" {
		// log.Fatal("please specify a URL with -url")
		app.Error.Println("please specify a URL with -url")
	}
	if *accessKey == "" {
		// log.Fatal("missing -access-key or AWS_ACCESS_KEY environment variable")
		app.Error.Println("missing -access-key or AWS_ACCESS_KEY environment variable")
	}
	if *secretKey == "" {
		// log.Fatal("missing -secret-key or AWS_SECRET_KEY environment variable")
		app.Error.Println("missing -secret-key or AWS_SECRET_KEY environment variable")

	}
	if *region == "" {
		// log.Fatal("please specify an AWS region with -region")
		app.Error.Println("please specify an AWS region with -region")

	}

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
