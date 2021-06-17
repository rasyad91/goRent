package main

import (
	"flag"

	"github.com/olivere/env"
)

var (
	accessKey *string
	secretKey *string
	esUrl     *string
	sniff     *bool
	region    *string

	port       *string
	identifier *string
	domain     *string

	dbDialect   *string
	dbUser      *string
	dbPassword  *string
	dbName      *string
	dbHost      *string
	dbPort      *string
	dbParseTime *bool

	mailhost     *string
	mailport     *int
	mailUsername *string
	mailPassword *string
)

func init() {
	// server flags
	port = flag.String("port", "8080", "server port to listen on")
	identifier = flag.String("identifier", "GoRent", "unique identifier")
	domain = flag.String("domain", "localhost", "domain name (e.g. example.com)")
	// db flags
	dbUser = flag.String("dbuser", "", "database user")
	dbPassword = flag.String("dbpassword", "", "database password")
	dbName = flag.String("dbname", "", "database name")
	dbHost = flag.String("dbhost", "", "database host")
	dbPort = flag.String("dbport", "", "database port")
	dbDialect = flag.String("dbdialect", "mysql", "type of database, eg. postgres, mysql, maria")

	dbParseTime = flag.Bool("dbparsetime", true, "database parse time option")
	// aws flags
	accessKey = flag.String("access-key", env.String("", "AWS_ACCESS_KEY", "AWS_ACCESS_KEY_ID"), "Access Key ID")
	secretKey = flag.String("secret-key", env.String("", "AWS_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"), "Secret access key")
	esUrl = flag.String("esUrl", "", "Elasticsearch URL")
	sniff = flag.Bool("sniff", false, "Enable or disable sniffing")
	region = flag.String("region", "", "AWS Region name")

	// email flags
	mailhost = flag.String("mailhost", "localhost", "mail host")
	mailport = flag.Int("mailport", 1025, "mail port")
	mailUsername = flag.String("mailuser", "", "mailuser")
	mailPassword = flag.String("mailpassword", "", "mailpassword")

	flag.Parse()
}
