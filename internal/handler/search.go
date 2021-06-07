package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/olivere/elastic/aws"
	// "github.com/olivere/elastic"
	"github.com/olivere/elastic/v7"

	aws "github.com/olivere/elastic/aws/v4"
)

func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) SearchResult(w http.ResponseWriter, r *http.Request) {

	awsSigningFn := awsSigning(m.App.AwsAccessKey, m.App.AwsSecretKey, m.App.AwsRegion)
	awsClient, err := awsCreateClient(m.App.AwsUrl, m.App.AwsSniff, awsSigningFn)

	trialMultiSearchQuery(awsClient)

	fmt.Printf("\n\n\n\n\n")

	if err != nil {
		m.App.Error.Println(err)
	}

	var data = make(map[string]interface{})
	var product []model.Product
	// var urlquery string
	x := r.URL.Query()
	fmt.Println(x)

	searchkeywords := strings.ToLower(url.QueryEscape(x["q"][0])) //hockey+sticks

	if len(searchkeywords) == 0 {
		product = searchEmptyQuery(awsClient)
	} else {
		product = searchQuery(awsClient, searchkeywords)
	}

	data["product"] = product

	if err := render.Template(w, r, "searchresult.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func awsSigning(awsAccesKey, awsSecretKey, awsRegoin string) *http.Client {

	signingClient := aws.NewV4SigningClient(credentials.NewStaticCredentials(
		awsAccesKey,
		awsSecretKey,
		"",
	), awsRegoin)

	return signingClient

}

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

func searchQuery(client *elastic.Client, searchKewords string) []model.Product {

	stringQuery := elastic.NewQueryStringQuery(searchKewords)

	searchResult, err := client.Search().
		Index("sample_product_list"). // search in index "tweets"
		Query(stringQuery).           // specify the query
		Pretty(true).                 // pretty print request and response JSON
		Do(context.Background())      // execute

	if err != nil {
		fmt.Println("error from search", err)
	}

	var product []model.Product

	for _, hit := range searchResult.Hits.Hits {

		var t model.Product

		if err := json.Unmarshal(hit.Source, &t); err != nil {
			// log.Errorf("ERROR UNMARSHALLING ES SUGGESTION RESPONSE: %v", err)
			continue
		}
		if err != nil {
			// Deserialization failed
			fmt.Println("error unmarshaling json", err)
		}
		product = append(product, t)
	}
	fmt.Println(product)
	return product

}

func searchEmptyQuery(client *elastic.Client) []model.Product {

	searchResult, err := client.Search().
		Index("sample_product_list"). // search in index "tweets"
		Query(elastic.NewMatchAllQuery()).
		Pretty(true).            // pretty print request and response JSON
		Do(context.Background()) // execute

	if err != nil {
		fmt.Println("error from search", err)
	}

	var product []model.Product

	for _, hit := range searchResult.Hits.Hits {

		var t model.Product

		if err := json.Unmarshal(hit.Source, &t); err != nil {
			// log.Errorf("ERROR UNMARSHALLING ES SUGGESTION RESPONSE: %v", err)
			continue
		}
		if err != nil {
			// Deserialization failed
			fmt.Println("error unmarshaling json", err)

		}
		product = append(product, t)

	}

	fmt.Println(product)
	return product

}

func trialMultiSearchQuery(client *elastic.Client) {

	sreq3 := elastic.NewTermQuery("brand_name", "nike")
	sreq4 := elastic.NewRangeQuery("price").From(0).To(49.99)
	sreq5 := elastic.NewRangeQuery("rating").From(0).To(0)

	query := elastic.NewBoolQuery().Must(sreq3, sreq4, sreq5)

	searchResult, err := client.Search().
		Index("sample_product_list").
		Type("sampleproducttype"). // search in type
		Query(query).
		Do(context.Background()) // execute

	if err != nil {
		fmt.Println("error from multi search request", err)
	}

	var product []model.Product

	sres := searchResult

	for _, hit := range sres.Hits.Hits {

		var t model.Product

		if err := json.Unmarshal(hit.Source, &t); err != nil {
			// log.Errorf("ERROR UNMARSHALLING ES SUGGESTION RESPONSE: %v", err)
			continue
		}
		if err != nil {
			fmt.Println("error unmarshaling json", err)

		}
		product = append(product, t)

	}

	// return product
	fmt.Println(product)

}
