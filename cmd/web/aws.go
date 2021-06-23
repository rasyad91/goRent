package main

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsS3 "github.com/aws/aws-sdk-go/aws/session"

	"github.com/olivere/elastic/v7"

	aws4 "github.com/olivere/elastic/aws/v4"
)

// NewAWSSession cretes a new session for the usage of AWS S3
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

// NewAWSSession cretes a new client for the usage of AWS Elasticsearch
func newAWSClient() (*elastic.Client, error) {

	//start of elastic search codes
	if *esUrl == "" || *accessKey == "" || *secretKey == "" || *region == "" {
		// log.Fatal("please specify a URL with -url")
		app.Error.Println("Missing required flags")
	}

	awsSigningFn := awsSigning(*accessKey, *secretKey, *region)
	awsClient, err := awsCreateClient(*esUrl, *sniff, awsSigningFn)

	return awsClient, err
}

// awsSigning works together with func newAWSClient to sign requests before sending to AWS Elasticsearch
func awsSigning(awsAccesKey, awsSecretKey, awsRegoin string) *http.Client {
	signingClient := aws4.NewV4SigningClient(credentials.NewStaticCredentials(
		awsAccesKey,
		awsSecretKey,
		"",
	), awsRegoin)
	return signingClient
}

// awsCreateClient works together with func newAWSClient to create a client instance with AWS Elasticsearch
func awsCreateClient(url string, sniff bool, signingClient *http.Client) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(sniff),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(signingClient),
	)
	if err != nil {
		// log.Fatal(err)
		return client, err
	}
	return client, nil
}
