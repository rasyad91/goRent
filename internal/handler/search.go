package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	// "github.com/olivere/elastic/aws"
	// "github.com/olivere/elastic"
	"github.com/olivere/elastic/v7"
)

func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) SearchResult(w http.ResponseWriter, r *http.Request) {

	var data = make(map[string]interface{})
	var product []model.Product
	//var priceFilter bool
	// var urlquery string
	x := r.URL.Query()
	fmt.Println(x)

	searchkeywords := strings.ToLower(url.QueryEscape(x["q"][0])) //hockey+sticks

	//check if minprice/ map exists

	_, okMin := x["minprice"]
	_, okMax := x["maxprice"]
	err := fmt.Errorf("")
	fmt.Println(err)

	// owner_id := 1
	// imagequery(m.App.AWSClient, owner_id)

	if okMin && okMax {
		//please use multiseach instead.
		fmt.Println("multi search functon got fired")
		product, err = trialMultiSearchQuery(m.App.AWSClient, x["minprice"][0], x["maxprice"][0], searchkeywords)
		if err != nil {
			m.App.Error.Println(err)
		}

	} else {

		if len(searchkeywords) == 0 {
			product = searchEmptyQuery(m.App.AWSClient)
		} else {
			fmt.Println("single function got fired")
			product = searchQuery(m.App.AWSClient, searchkeywords)
		}

	}

	data["product"] = product
	data["query"] = searchkeywords

	if err := render.Template(w, r, "searchresult.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}

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

func trialMultiSearchQuery(client *elastic.Client, min, max, searchKeywords string) ([]model.Product, error) {

	var product []model.Product

	minPrice, err := strconv.ParseFloat(min, 32)
	if err != nil {
		return product, err
	}

	maxPrice, err := strconv.ParseFloat(max, 32)
	if err != nil {
		return product, err
	}

	stringQuery := elastic.NewQueryStringQuery(searchKeywords)

	// sreq3 := elastic.NewTermQuery("brand_name", "nike")

	// searchResult, err := client.Search().
	// 	Index("sample_product_list"). // search in index "tweets"
	// 	Query(stringQuery).           // specify the query
	// 	Pretty(true).                 // pretty print request and response JSON
	// 	Do(context.Background())      // execute

	priceQuery := elastic.NewRangeQuery("price").From(minPrice).To(maxPrice)
	// sreq5 := elastic.NewRangeQuery("rating").From(0).To(0)

	query := elastic.NewBoolQuery().Must(stringQuery, priceQuery)

	searchResult, err := client.Search().
		Index("sample_product_list").
		// Type("sampleproducttype"). // search in type
		Query(query).
		Do(context.Background()) // execute

	if err != nil {
		fmt.Println("error from multi search request", err)
	}

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

	fmt.Println(product)
	return product, nil
}

//for rachel.
func imagequery(client *elastic.Client, owner_id int) ([]model.Product, error) {

	var product []model.Product

	req := elastic.NewTermQuery("owner_id", owner_id) //maybe change "1"  to ownder_id , received frum functoin

	searchResult, err := client.Search().
		Index("sample_product_list"). // search in index "tweets"
		Query(req).                   // specify the query
		Pretty(true).                 // pretty print request and response JSON
		Do(context.Background())      // execute

	if err != nil {
		fmt.Println("error from querying image", err)
	}

	for _, hit := range searchResult.Hits.Hits {

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

	fmt.Println("these are the products image query.")
	return product, nil
}
