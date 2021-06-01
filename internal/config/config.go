package config

import (
	"goRent/internal/driver/mysqlDriver"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	DB            *mysqlDriver.DB
	Session       *scs.SessionManager
	InProduction  bool
	Domain        string
	PreferenceMap map[string]string
	Version       string
	Identifier    string
	Info          *log.Logger
	Error         *log.Logger
}
