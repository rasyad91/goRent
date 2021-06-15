package handler

import (
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

		ID:          12,  //change values
		OwnerID:     1,   //change values
		OwnerName:   "1", //change values
		Brand:       "Disney",
		Category:    "DVD",                                                                                                                                                                                                                      //change values
		Title:       "Lion king movie cd",                                                                                                                                                                                                       //change values
		Description: "You know, one day our children wont even know what a CD player is. Rent this, take this time to educate them on what a CD players is. Besides, lion king is very good show about children. #ensureCDplayerGetsPassedDown", //change values
		Price:       3.99,                                                                                                                                                                                                                       //change values
		Images:      []string{"https://wooteam-productslist.s3.ap-southeast-1.amazonaws.com/product_list/images/2021-06-14_15-00-13_1.jpg"},                                                                                                     //change values
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

func DeleteProductsElasticUserID(r *http.Request, client *elastic.Client, i int) {

	// Delete all documents by sandrae
	q := elastic.NewTermQuery("owner_id", i)
	res, err := client.DeleteByQuery().
		Index("product_list").
		Query(q).
		Pretty(true).
		Do(r.Context())
	if err != nil {
		// return fmt.Errorf("error deleting products belonging to ")
		fmt.Println("error from delting product from index cause user deleted account", err)
	}
	if res == nil {
		fmt.Println("expected response != nil; got: %v", res)
	}

}
