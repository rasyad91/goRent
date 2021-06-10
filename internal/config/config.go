package config

import (
	"goRent/internal/model"
	"log"

	"github.com/alexedwards/scs/v2"
	awsS3 "github.com/aws/aws-sdk-go/aws/session"
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
	AWSS3Session  *awsS3.Session
}

const (
	DateLayout = "02-01-2006"
)
