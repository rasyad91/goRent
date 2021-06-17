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

// rent + toh  -> budget version of a brain
var singleCatMap = make(map[string]int)
var doubleCatMap = make(map[string]int)

func rentohQuery(s string) {

	var foundCategory string
	var wg sync.WaitGroup
	var mutex sync.Mutex

	catCh := make(chan model.RentohKeyword)
	catSCh := make(chan model.RentohKeyword)

	rentohRawQuery := strings.Replace(s, "+", ",", -1)
	rentohArray := strings.Split(rentohRawQuery, ",")

	if len(rentohArray) < 2 {

		for i := 0; i < len(rentohArray); i++ {
			wg.Add(1)
			go checkSingleCategory(i, rentohArray[i], catSCh, &wg, &mutex)
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

		if len(foundCategory) == 0 {
			//store unknown keyword somewhese else
		}
		return
	}

	for i := 0; i < len(rentohArray)-1; i++ {
		wg.Add(1)
		go checkDoubleCategory(i, rentohArray[i], rentohArray[i+1], catCh, &wg, &mutex)
	}
	go func() {
		wg.Wait()
		close(catCh)
	}()

	for categoryName := range catCh {
		if categoryName.Index != -1 {
			foundCategory = categoryName.Keyword
			//process keywords for left and right of index
			fmt.Println("this is the name of the category [DOUBLE] FOUND", foundCategory)
			return
		}
	}

	for i := 0; i < len(rentohArray); i++ {
		wg.Add(1)
		go checkSingleCategory(i, rentohArray[i], catSCh, &wg, &mutex)
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

func checkSingleCategory(i int, s string, catCh chan model.RentohKeyword, wg *sync.WaitGroup, mutex *sync.Mutex) {

	defer wg.Done()
	// var singleCatMap = make(map[string]int)
	mutex.Lock()

	fmt.Println("these are the cat strings from SINGLE", s)

	if k, ok := singleCatMap[s]; ok {
		singleCatMap[s] = k + 1
		catCh <- model.RentohKeyword{i, s}

	} else {
		catCh <- model.RentohKeyword{-1, ""}
	}
	mutex.Unlock()

}

func checkDoubleCategory(i int, s1, s2 string, catCh chan model.RentohKeyword, wg *sync.WaitGroup, mutex *sync.Mutex) {

	defer wg.Done()
	mutex.Lock()

	catString := s1 + " " + s2
	catString2 := s2 + " " + s1

	fmt.Println("these are the cat strings from DOUBLE", catString)
	fmt.Println("these are the cat strings from DOUBLE", catString2)

	if k, ok := doubleCatMap[catString]; ok {

		doubleCatMap[catString] = k + 1
		catCh <- model.RentohKeyword{i, catString}

	} else if k, ok := doubleCatMap[catString2]; ok {
		doubleCatMap[catString2] = k + 1
		catCh <- model.RentohKeyword{i, catString2}

	} else {
		catCh <- model.RentohKeyword{-1, ""}

	}
	mutex.Unlock()

}

func (m *Repository) SearchTrend(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})

	categoriesSlice := []model.SearchTrends{}

	for k, v := range singleCatMap {
		categoriesSlice = append(categoriesSlice, model.SearchTrends{CategoryName: k, Count: v})

	}

	for k, v := range doubleCatMap {
		categoriesSlice = append(categoriesSlice, model.SearchTrends{CategoryName: k, Count: v})

	}

	SortArrayCategory(categoriesSlice)
	data["Categories"] = categoriesSlice
	// fmt.Println("these are all the categories", categoriesSlice)

	if err := render.Template(w, r, "searchtrend.page.html", &render.TemplateData{
		Data: data,
		Form: &form.Form{},
	}); err != nil {
		m.App.Error.Println(err)
	}

}

func (m *Repository) CreateCategoryDataBase(w http.ResponseWriter, r *http.Request) {

	singleCatMap["studio"] = 0
	singleCatMap["condo"] = 0
	singleCatMap["hdb"] = 0
	singleCatMap["office"] = 0
	singleCatMap["shops"] = 0
	singleCatMap["hockey"] = 0
	singleCatMap["book"] = 0
	singleCatMap["books"] = 0
	singleCatMap["rugby"] = 0
	singleCatMap["wedding"] = 0
	singleCatMap["dvd"] = 0
	singleCatMap["computer"] = 0
	singleCatMap["computers"] = 0
	singleCatMap["laptop"] = 0
	singleCatMap["laptops"] = 0
	singleCatMap["vacuum"] = 0
	singleCatMap["car"] = 0
	singleCatMap["piano"] = 0
	singleCatMap["kayak"] = 0
	singleCatMap["surfboard"] = 0
	singleCatMap["bridesmaids"] = 0
	singleCatMap["kitchen"] = 0
	singleCatMap["bike"] = 0
	singleCatMap["bicycles"] = 0
	singleCatMap["kiosk"] = 0
	singleCatMap["camera"] = 0
	singleCatMap["tripod"] = 0
	singleCatMap["tents"] = 0
	singleCatMap["clothes"] = 0
	singleCatMap["guitar"] = 0
	singleCatMap["movies"] = 0
	singleCatMap["boat"] = 0
	singleCatMap["textbook"] = 0
	singleCatMap["textbooks"] = 0
	singleCatMap["plants"] = 0
	singleCatMap["party"] = 0
	singleCatMap["jeans"] = 0
	singleCatMap["drones"] = 0
	singleCatMap["prams"] = 0
	doubleCatMap["washing machines"] = 0
	doubleCatMap["washing machine"] = 0
	doubleCatMap["wedding dress"] = 0
	doubleCatMap["movie projector"] = 0
	doubleCatMap["sewing machine"] = 0
	doubleCatMap["empty room"] = 0
	doubleCatMap["wedding gown"] = 0
	doubleCatMap["party tent"] = 0
	doubleCatMap["moving gear"] = 0
	doubleCatMap["camera lens"] = 0
	doubleCatMap["camera equipment"] = 0
	doubleCatMap["lawn mower"] = 0
	doubleCatMap["parking space"] = 0
	doubleCatMap["camera tripod"] = 0
	doubleCatMap["baking equipment"] = 0
	doubleCatMap["video games"] = 0
	doubleCatMap["waffle machine"] = 0
	doubleCatMap["baking machine"] = 0
	doubleCatMap["storage space"] = 0
	doubleCatMap["office space"] = 0
	doubleCatMap["campsite rental"] = 0
	doubleCatMap["fishing rental"] = 0
	doubleCatMap["conferences space"] = 0
	doubleCatMap["tech rental"] = 0
	doubleCatMap["arcade games"] = 0
	doubleCatMap["power tools"] = 0
	doubleCatMap["av equipment"] = 0
	doubleCatMap["board games"] = 0
	doubleCatMap["camping space"] = 0
	doubleCatMap["camping equipment"] = 0
	doubleCatMap["hockey sticks"] = 0
	doubleCatMap["hockey equipment"] = 0
	doubleCatMap["golf carts"] = 0
	doubleCatMap["bounce houses"] = 0
	doubleCatMap["party dresses"] = 0
	doubleCatMap["format wear"] = 0
	doubleCatMap["night gowns"] = 0
	doubleCatMap["baby gear"] = 0
	doubleCatMap["photography gear"] = 0
	doubleCatMap["musical instruments"] = 0
	doubleCatMap["sporting goods"] = 0
	doubleCatMap["fishing gear"] = 0
	doubleCatMap["sewing machine"] = 0
	doubleCatMap["christmas trees"] = 0

	http.Redirect(w, r, "/v1/user/searchtrend", http.StatusSeeOther)

}
