package handler

import (
	"fmt"
	"goRent/internal/form"
	"goRent/internal/model"
	"goRent/internal/render"
	"net/http"
	"strings"
	"sync"
)

type Cat struct {
	data map[string]int
	l    sync.Mutex
}

var singleCat Cat = Cat{}
var doubleCat Cat = Cat{}

// rentohQuery takes in a user input from the search page, processes the query into arrays and then fire them off as different go routines to check against Rentoh's brains
func rentohQuery(s string) {

	var foundCategory string
	var wg sync.WaitGroup

	catCh := make(chan model.RentohKeyword)
	catSCh := make(chan model.RentohKeyword)

	rentohRawQuery := strings.Replace(s, "+", ",", -1)
	rentohArray := strings.Split(rentohRawQuery, ",")

	if len(rentohArray) < 2 {

		for i := 0; i < len(rentohArray); i++ {
			wg.Add(1)
			go singleCat.checkSingleCategory(i, rentohArray[i], catSCh, &wg)
		}

		go func() {
			wg.Wait()
			close(catSCh)
		}()
		for categoryName := range catSCh {
			if categoryName.Index != -1 {
				foundCategory = categoryName.Keyword
				fmt.Println("this is the name of the category found", foundCategory)
				break
			}
		}
		return
	}

	for i := 0; i < len(rentohArray)-1; i++ {
		wg.Add(1)
		go doubleCat.checkDoubleCategory(i, rentohArray[i], rentohArray[i+1], catCh, &wg)
	}
	go func() {
		wg.Wait()
		close(catCh)
	}()

	for categoryName := range catCh {
		if categoryName.Index != -1 {
			foundCategory = categoryName.Keyword
			fmt.Println("this is the name of the category [DOUBLE] FOUND", foundCategory)
			return
		}
	}

	for i := 0; i < len(rentohArray); i++ {
		wg.Add(1)
		go singleCat.checkSingleCategory(i, rentohArray[i], catSCh, &wg)
	}

	go func() {
		wg.Wait()
		close(catSCh)
	}()
	for categoryName := range catSCh {
		if categoryName.Index != -1 {
			foundCategory = categoryName.Keyword
			fmt.Println("this is the name of the category [SINGLE] FOUND", foundCategory)
			break
		}
	}
}

// checkSingleCategory works hand in hand with rentohQuery to check against Rentoh's brain that contain the single worded categories
func (m *Cat) checkSingleCategory(i int, s string, catCh chan model.RentohKeyword, wg *sync.WaitGroup) {

	defer wg.Done()
	m.l.Lock()
	defer m.l.Unlock()

	fmt.Println("these are the cat strings from SINGLE", s)

	if k, ok := m.data[s]; ok {
		m.data[s] = k + 1
		catCh <- model.RentohKeyword{Index: i, Keyword: s}

	} else {
		catCh <- model.RentohKeyword{Index: -1, Keyword: ""}
	}

}

// checkDoubleCategory works hand in hand with rentohQuery to check against Rentoh's brain that contain the double worded categories
func (m *Cat) checkDoubleCategory(i int, s1, s2 string, catCh chan model.RentohKeyword, wg *sync.WaitGroup) {

	defer wg.Done()
	m.l.Lock()
	defer m.l.Unlock()

	catString := s1 + " " + s2
	catString2 := s2 + " " + s1

	fmt.Println("these are the cat strings from DOUBLE", catString)
	fmt.Println("these are the cat strings from DOUBLE", catString2)

	if k, ok := m.data[catString]; ok {

		m.data[catString] = k + 1
		catCh <- model.RentohKeyword{Index: i, Keyword: catString}

	} else if k, ok := m.data[catString2]; ok {
		m.data[catString2] = k + 1
		catCh <- model.RentohKeyword{Index: i, Keyword: catString2}
	} else {
		catCh <- model.RentohKeyword{Index: -1, Keyword: ""}

	}

}

// SearchTrend is responsible for creating and storing the categories and displaying on SearchTrend data page.
func (m *Repository) SearchTrend(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	categoriesSlice := []model.SearchTrends{}

	for k, v := range singleCat.data {
		categoriesSlice = append(categoriesSlice, model.SearchTrends{CategoryName: k, Count: v})

	}

	for k, v := range doubleCat.data {
		categoriesSlice = append(categoriesSlice, model.SearchTrends{CategoryName: k, Count: v})

	}

	SortArrayCategory(categoriesSlice)
	data["Categories"] = categoriesSlice

	if err := render.Template(w, r, "searchtrend.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}

}

// CreateCategoryDataBase is responsible for teaching, guiding and helping Rentoh become smarter. (Rentoh's brain)
func (m *Repository) CreateCategoryDataBase(w http.ResponseWriter, r *http.Request) {

	singleCat.data = make(map[string]int)
	doubleCat.data = make(map[string]int)

	singleCat.data["studio"] = 0
	singleCat.data["condo"] = 0
	singleCat.data["hdb"] = 0
	singleCat.data["office"] = 0
	singleCat.data["shops"] = 0
	singleCat.data["hockey"] = 0
	singleCat.data["book"] = 0
	singleCat.data["books"] = 0
	singleCat.data["rugby"] = 0
	singleCat.data["wedding"] = 0
	singleCat.data["dvd"] = 0
	singleCat.data["computer"] = 0
	singleCat.data["computers"] = 0
	singleCat.data["laptop"] = 0
	singleCat.data["laptops"] = 0
	singleCat.data["vacuum"] = 0
	singleCat.data["car"] = 0
	singleCat.data["piano"] = 0
	singleCat.data["kayak"] = 0
	singleCat.data["surfboard"] = 0
	singleCat.data["bridesmaids"] = 0
	singleCat.data["kitchen"] = 0
	singleCat.data["bike"] = 0
	singleCat.data["bicycles"] = 0
	singleCat.data["kiosk"] = 0
	singleCat.data["camera"] = 0
	singleCat.data["tripod"] = 0
	singleCat.data["tents"] = 0
	singleCat.data["clothes"] = 0
	singleCat.data["guitar"] = 0
	singleCat.data["movies"] = 0
	singleCat.data["boat"] = 0
	singleCat.data["textbook"] = 0
	singleCat.data["textbooks"] = 0
	singleCat.data["plants"] = 0
	singleCat.data["party"] = 0
	singleCat.data["jeans"] = 0
	singleCat.data["drones"] = 0
	singleCat.data["prams"] = 0
	doubleCat.data["washing machines"] = 0
	doubleCat.data["washing machine"] = 0
	doubleCat.data["wedding dress"] = 0
	doubleCat.data["movie projector"] = 0
	doubleCat.data["sewing machine"] = 0
	doubleCat.data["empty room"] = 0
	doubleCat.data["wedding gown"] = 0
	doubleCat.data["party tent"] = 0
	doubleCat.data["moving gear"] = 0
	doubleCat.data["camera lens"] = 0
	doubleCat.data["camera equipment"] = 0
	doubleCat.data["lawn mower"] = 0
	doubleCat.data["parking space"] = 0
	doubleCat.data["camera tripod"] = 0
	doubleCat.data["baking equipment"] = 0
	doubleCat.data["video games"] = 0
	doubleCat.data["waffle machine"] = 0
	doubleCat.data["baking machine"] = 0
	doubleCat.data["storage space"] = 0
	doubleCat.data["office space"] = 0
	doubleCat.data["campsite rental"] = 0
	doubleCat.data["fishing rental"] = 0
	doubleCat.data["conferences space"] = 0
	doubleCat.data["tech rental"] = 0
	doubleCat.data["arcade games"] = 0
	doubleCat.data["power tools"] = 0
	doubleCat.data["av equipment"] = 0
	doubleCat.data["board games"] = 0
	doubleCat.data["camping space"] = 0
	doubleCat.data["camping equipment"] = 0
	doubleCat.data["hockey sticks"] = 0
	doubleCat.data["hockey equipment"] = 0
	doubleCat.data["golf carts"] = 0
	doubleCat.data["bounce houses"] = 0
	doubleCat.data["party dresses"] = 0
	doubleCat.data["format wear"] = 0
	doubleCat.data["night gowns"] = 0
	doubleCat.data["baby gear"] = 0
	doubleCat.data["photography gear"] = 0
	doubleCat.data["musical instruments"] = 0
	doubleCat.data["sporting goods"] = 0
	doubleCat.data["fishing gear"] = 0
	doubleCat.data["sewing machine"] = 0
	doubleCat.data["christmas trees"] = 0

	http.Redirect(w, r, "/v1/user/searchtrend", http.StatusSeeOther)

}
