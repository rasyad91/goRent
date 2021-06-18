package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/olivere/elastic/v7"
)

//  Search is the handler that displays the product search page
func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {
	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}
}

// SearchResult handler takes in a user's seach query and fires requests to elasticsearch as well as Rentoh
func (m *Repository) SearchResult(w http.ResponseWriter, r *http.Request) {

	var data = make(map[string]interface{})
	var product []model.Product
	x := r.URL.Query()

	searchkeywords := strings.ToLower(url.QueryEscape(x["q"][0])) //hockey+sticks

	//call rentoh
	go rentohQuery(searchkeywords)

	_, okMin := x["minprice"]
	_, okMax := x["maxprice"]
	err := fmt.Errorf("")
	_ = err

	form := form.New(r.PostForm)
	form.Required("minprice", "maxprice")

	if okMin && okMax {
		//please use multiseach instead.
		m.App.Info.Println("multi search functon got fired")
		product, err = trialMultiSearchQuery(m.App.AWSClient, x["minprice"][0], x["maxprice"][0], searchkeywords)
		if err != nil {
			m.App.Error.Println(err)
		}
	} else {
		m.App.Info.Println("single function got fired")
		product = searchQuery(m.App.AWSClient, searchkeywords)
	}

	data["product"] = product
	data["query"] = searchkeywords

	if err := render.Template(w, r, "searchresult.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}

}

// searchQuery takes in a santizied query string from the user as an argument and returns results in []model.Product.
func searchQuery(client *elastic.Client, searchKeywords ...string) []model.Product {

	var searchResult *elastic.SearchResult
	var query elastic.Query

	if searchKeywords[0] != "" {
		query = elastic.NewQueryStringQuery(searchKeywords[0])
	} else {
		query = elastic.NewMatchAllQuery()
	}

	searchResult, err := client.Search().
		Index("product_list"). // search in index "tweets"
		Query(query).          // specify the query
		Pretty(true).          // pretty print request and response JSON
		From(0).
		Size(20).
		Do(context.Background()) // execute
	if err != nil {
		fmt.Println("error from search", err)
	}

	var products []model.Product
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
		products = append(products, t)
	}

	return products

}

// trialMultiSearchQuery takes in a santizied query string with other arguments (from filter) from the user as an argument and returns results.
func trialMultiSearchQuery(client *elastic.Client, min, max string, searchKeywords ...string) ([]model.Product, error) {

	var products []model.Product
	minPrice, err := strconv.ParseFloat(min, 32)
	if err != nil {
		return nil, err
	}
	maxPrice, err := strconv.ParseFloat(max, 32)
	if err != nil {
		return nil, err
	}
	var stringQuery elastic.Query
	if searchKeywords[0] != "" {
		stringQuery = elastic.NewQueryStringQuery(searchKeywords[0])
	} else {
		stringQuery = elastic.NewMatchAllQuery()
		fmt.Printf("query is matchallquery %#v", stringQuery)
	}

	priceQuery := elastic.NewRangeQuery("price").From(minPrice).To(maxPrice)

	query := elastic.NewBoolQuery().Must(stringQuery, priceQuery)

	searchResult, err := client.Search().
		Index("product_list").
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
		products = append(products, t)
	}

	return products, nil
}
