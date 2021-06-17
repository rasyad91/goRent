package handler

import (
	"context"
	"fmt"
	"goRent/internal/model"
	"net/http"
	"strconv"

	// "github.com/olivere/elastic/aws"
	// "github.com/olivere/elastic"
	"github.com/olivere/elastic/v7"
)

func ManualProductFix(r *http.Request, client *elastic.Client, product model.Product) {

	var editedProduct = model.Product{

		ID:          10,  //change values
		OwnerID:     1,   //change values
		OwnerName:   "1", //change values
		Brand:       "HP",
		Category:    "Laptops",                                                                                                                                                                                                                                                       //change values
		Title:       "school computer laptops",                                                                                                                                                                                                                                       //change values
		Description: "All these computer laptops are availale for rent. They are sanizited frequently and definitely after they are returned to us. All 12 of them are good for you to host your own programming classes. If you need more quantity, please feel free to ask or PM.", //change values
		Price:       99.99,                                                                                                                                                                                                                                                           //change values
		Images:      []string{"https://wooteam-productslist.s3.ap-southeast-1.amazonaws.com/product_list/images/2021-06-14_15-00-13_1.jpg"},                                                                                                                                          //change values
	}

	_, err := client.Index().
		Index("product_list").
		Type("product").
		Id(strconv.Itoa(editedProduct.ID)).
		BodyJson(editedProduct).
		// Do(context.Background())
		Do(r.Context())

	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("\n\nSuccesfully patched the product ID %d, Brand %s, Title %s\n\n", editedProduct.ID, editedProduct.Brand, editedProduct.Title)

}

func ManualDeleteProductsElastic(r *http.Request, client *elastic.Client, i int) {

	// Delete all documents by sandrae
	fmt.Println("DeleteProductsElasticUserID fn got triggered. please check for products index")
	q := elastic.NewTermQuery("ID", i)
	res, err := client.DeleteByQuery().
		Index("product_list").
		Query(q).
		Pretty(true).
		Do(r.Context())
	if err != nil {
		fmt.Println("error from delting product from index cause user deleted account", err)
	}
	if res == nil {
		fmt.Printf("\nexpected response != nil; got: %v\n", res)
	}

}

func DeleteProductsElasticUserID(r *http.Request, client *elastic.Client, s string) error {
	fmt.Println("owner_id value passed it", s)

	owner_id_int, err := strconv.Atoi(s)

	if err != nil {
		return err
	}
	q := elastic.NewTermQuery("owner_id", owner_id_int)
	res, err := client.DeleteByQuery().
		Index("product_list").
		Query(q).
		Pretty(true).
		Do(r.Context())
	if err != nil {
		return err
	}
	if res == nil {
		return err
	}

	return nil

}

func ManualUpdateViaDoc(r *http.Request, client *elastic.Client) {
	// client := setupTestClient(t) // , SetTraceLog(log.New(os.Stdout, "", 0)))

	fmt.Println("update function got called")
	doc, err := client.Get().
		Index("product_list").Id("10").
		Do(context.Background())

	res, err := client.Update().
		Index("product_list").Id(doc.Id).
		Doc(map[string]interface{}{"brand_name": "alvin"}).
		IfSeqNo(*doc.SeqNo).
		IfPrimaryTerm(*doc.PrimaryTerm).
		FetchSource(true).
		Do(r.Context())

	if err != nil {
		fmt.Println(err)
	}
	if res == nil {
		fmt.Println(err)
	}
	if res.GetResult == nil {
		fmt.Println(err)
	}
	// return nil
}

func ReviewUpdateViaDoc(r *http.Request, client *elastic.Client, i int, f float32) error {

	elastic_id := strconv.Itoa(i)
	fmt.Println("review update function got called")
	doc, _ := client.Get().
		Index("product_list").Id(elastic_id).
		Do(r.Context())

	res, err := client.Update().
		Index("product_list").Id(doc.Id).
		Doc(map[string]interface{}{"rating": f}).
		IfSeqNo(*doc.SeqNo).
		IfPrimaryTerm(*doc.PrimaryTerm).
		FetchSource(true).
		Do(r.Context())

	if err != nil {
		return err
	}
	if res == nil {
		return err
	}
	if res.GetResult == nil {
		return err
	}
	return nil
}
