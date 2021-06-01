package config

import (
	"goRent/internal/driver/mysqlDriver"
	"goRent/internal/model"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	DB            *mysqlDriver.DB
	Session       *scs.SessionManager
	Domain        string
	PreferenceMap map[string]string
	Info          *log.Logger
	Error         *log.Logger
	MailChan      chan model.MailData
}
