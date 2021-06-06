package config

import (
	"goRent/internal/model"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	Session       *scs.SessionManager
	Domain        string
	PreferenceMap map[string]string
	Info          *log.Logger
	Error         *log.Logger
	MailChan      chan model.MailData

	AwsAccessKey string
	AwsSecretKey string
	AwsUrl       string
	AwsSniff     bool
	AwsRegion    string
}

const (
	DateLayout = "02-01-2006"
)
