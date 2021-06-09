package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsS3 "github.com/aws/aws-sdk-go/aws/session"
)

// var app *config.AppConfig

func NewAWSSession() (*awsS3.Session, error) {

	sess, err := awsS3.NewSession(&aws.Config{
		Region: aws.String(*region),
		Credentials: credentials.NewStaticCredentials(
			*accessKey, // id
			*secretKey, // secret
			""),        // token can be left blank for now
	})

	return sess, err
}
