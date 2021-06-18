package config

import (
	"goRent/internal/model"
	"log"

	"github.com/alexedwards/scs/v2"
	awsS3 "github.com/aws/aws-sdk-go/aws/session"
	"github.com/olivere/elastic/v7"
)

// AppConfig is a global struct
type AppConfig struct {
	Session       *scs.SessionManager
	Domain        string
	Production    bool
	PreferenceMap map[string]string
	Info          *log.Logger
	Error         *log.Logger
	MailChan      chan model.MailData
	AWSClient     *elastic.Client
	AWSS3Session  *awsS3.Session
}

const (
	DateLayout           = "02-01-2006"
	AWSProductBucket     = "wooteam-productslist/product_list/images/"
	AWSProfileBucketLink = "wooteam-productslist/profile_images/"
	AWSProductImageLink  = "https://wooteam-productslist.s3.ap-southeast-1.amazonaws.com/product_list/images/"
	AWSProfileImageLink  = "https://wooteam-productslist.s3.ap-southeast-1.amazonaws.com/profile_images/"
)
