package main

import (
	"flag"

	"github.com/olivere/env"
)

var (
	accessKey   *string
	secretKey   *string
	esUrl       *string
	sniff       *bool
	region      *string
	port        *string
	identifier  *string
	domain      *string
	dbUser      *string
	dbPassword  *string
	dbName      *string
	dbHost      *string
	dbPort      *string
	dbParseTime *bool
)

func init() {
	port = flag.String("port", "8080", "server port to listen on")
	identifier = flag.String("identifier", "GoRent", "unique identifier")
	domain = flag.String("domain", "localhost", "domain name (e.g. example.com)")
	dbUser = flag.String("dbuser", "", "database user")
	dbPassword = flag.String("dbpassword", "", "database password")
	dbName = flag.String("dbname", "", "database name")
	dbHost = flag.String("dbhost", "", "database host")
	dbPort = flag.String("dbport", "", "database port")
	dbParseTime = flag.Bool("dbparsetime", true, "database parse time option")
	accessKey = flag.String("access-key", env.String("", "AWS_ACCESS_KEY", "AWS_ACCESS_KEY_ID"), "Access Key ID")
	secretKey = flag.String("secret-key", env.String("", "AWS_SECRET_KEY", "AWS_SECRET_ACCESS_KEY"), "Secret access key")
	esUrl = flag.String("esUrl", "", "Elasticsearch URL")
	sniff = flag.Bool("sniff", false, "Enable or disable sniffing")
	region = flag.String("region", "", "AWS Region name")
	flag.Parse()
}
