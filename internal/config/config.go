package config

import (
	"goRent/internal/model"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/olivere/elastic/v7"
)

type AppConfig struct {
	Session       *scs.SessionManager
	Domain        string
	PreferenceMap map[string]string
	Info          *log.Logger
	Error         *log.Logger
	MailChan      chan model.MailData
	AWSClient     *elastic.Client
}

const (
	DateLayout = "02-01-2006"
)
