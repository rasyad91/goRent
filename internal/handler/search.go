package handler

import (
	"fmt"
	"goRent/internal/render"
	"html"
	"net/http"
	"strings"
)

func (m *Repository) Search(w http.ResponseWriter, r *http.Request) {

	m.CreateProductList()

	if err := render.Template(w, r, "home.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) SearchResult(w http.ResponseWriter, r *http.Request) {
	x := r.URL.Query()
	fmt.Println(x)

	searchKW := r.FormValue("searchtext")

	sr := html.EscapeString(searchKW)
	srtL := strings.ToLower((sr))
	_ = srtL

	url := `http://localhost:9200/product_list/_search?q=` +
		`{
		"query": {
		  "fuzzy": {
			"q": {
			  "value":` + x["q"][0] + `,
			  "fuzziness": "AUTO",
			  "max_expansions": 50,
			  "prefix_length": 0,
			  "transpositions": true,
			  "rewrite": "constant_score"
			}
		  }
		}
	  }`

	fmt.Println("this is the query url", url)

	// 	response, err := http.Get(url)

	// 	if err != nil {
	// 		fmt.Printf("The HTTP request failed with error %s \n", err)
	// 	} else {
	// 		fmt.Println(url)
	// 		data, _ := ioutil.ReadAll(response.Body)
	// 		fmt.Println("[GET] API validation from client dashboard [My API] made:", response.StatusCode)

	// 		_ = data // fmt.Println(string(data))
	// 		response.Body.Close()

	// 		if response.StatusCode == 200 {

	// 			MyHandler(w, req, apiKey)
	// 			me["apiKey"] = "Your API Key has been validated. Thank you"
	// 			parseData := struct {
	// 				APIKey string
	// 				Data   map[string]string
	// 			}{
	// 				apiKey, me,
	// 			}
	// 			tpl.ExecuteTemplate(w, "userprofile.gohtml", parseData) //replace nil as data
	// 			return

	// 		}

	// m.CreateProductList()

	// if r.Method == http.MethodPost {
	// 	searchKW := r.FormValue("searchtext")

	// 	sr := html.EscapeString(searchKW)
	// 	srtL := strings.ToLower((sr))
	// 	_ = srtL

	// 	//can we store this value into middleware so the value can be passed to the next page?

	// 	http.Redirect(w, r, "/searchresult", http.StatusSeeOther)
	// }

	if err := render.Template(w, r, "searchresult.page.html", &render.TemplateData{
		Data: nil,
	}); err != nil {
		m.App.Error.Println(err)
	}

}
