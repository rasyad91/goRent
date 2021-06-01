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
}
