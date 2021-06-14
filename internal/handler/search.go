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

	form := form.New(r.PostForm)
	form.Required("minprice", "maxprice")

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
		fmt.Println("call search Query")
		product = searchQuery(m.App.AWSClient, searchkeywords)
		fmt.Println("single function got fired")
	}

	data["product"] = product
	data["query"] = searchkeywords

	if err := render.Template(w, r, "searchresult.page.html", &render.TemplateData{
		Data: data,
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func searchQuery(client *elastic.Client, searchKeywords ...string) []model.Product {

	var searchResult *elastic.SearchResult
	var query elastic.Query
	fmt.Println("length of searchKeywords", len(searchKeywords))
	fmt.Printf("whats in searchKeywords %#v", searchKeywords)

	if searchKeywords[0] != "" {
		query = elastic.NewQueryStringQuery(searchKeywords[0])
	} else {
		query = elastic.NewMatchAllQuery()
		fmt.Printf("query is matchallquery %#v", query)
	}

	searchResult, err := client.Search().
		// Index("sample_product_list"). // search in index "tweets"
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

	for i, v := range products {
		fmt.Printf("%d: %v\n", i, v.Title)
	}
	return products

}

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
		Index("product_list").
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
		products = append(products, t)
	}

	for i, v := range products {
		fmt.Printf("%d: %v\n", i, v.Title)
	}
	return products, nil
}

//for rachel.
// func imagequery(client *elastic.Client, owner_id int) ([]model.Product, error) {
// 	var product []model.Product
// 	req := elastic.NewTermQuery("owner_id", owner_id) //maybe change "1"  to ownder_id , received frum functoin
// 	searchResult, err := client.Search().
// 		Index("sample_product_list"). // search in index "tweets"
// 		Query(req).                   // specify the query
// 		Pretty(true).                 // pretty print request and response JSON
// 		Do(context.Background())      // execute

// 	if err != nil {
// 		fmt.Println("error from querying image", err)
// 	}

// 	for _, hit := range searchResult.Hits.Hits {

// 		var t model.Product

// 		if err := json.Unmarshal(hit.Source, &t); err != nil {
// 			// log.Errorf("ERROR UNMARSHALLING ES SUGGESTION RESPONSE: %v", err)
// 			continue
// 		}
// 		if err != nil {
// 			fmt.Println("error unmarshaling json", err)

// 		}
// 		product = append(product, t)

// 	}

// 	fmt.Println("these are the products image query.")
// 	return product, nil
// }
