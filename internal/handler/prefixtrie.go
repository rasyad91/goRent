package handler

import (
	"fmt"
)

func init() {

	// productList := InitProductsTrie()

}

type Trie struct { //for Tries Data Structure
	root *Node
}

type Node struct { //for Tries Data structure
	children [39]*Node //prefix tries only have array size of 26 (26 alphabets), but we have an array of 39 sizes. We can effectively store an asrrotement of values from a-z and 0-9 in ANY placement
	isEnd    bool
}

// type DBrepo struct {
// 	*sql.DB
// }

// func NewRepo(Conn *sql.DB) repository.DatabaseRepo {
// 	return &DBrepo{
// 		DB: Conn,
// 	}
// }

// var RentalProductsList []string //global variable

func (m *Repository) CreateProductList() {

	// var (
	// 	title string
	// 	brand string
	// )

	products, err := m.DB.GetAllProducts()
	if err != nil {
		m.App.Error.Println(err)
	}

	for _, v := range products {
		//do something with title brand
		fmt.Println(v)
	}

	fmt.Println("this should show all the products in the products table from mysql")
	fmt.Println("RentalProductsList")

}

func InitProductsTrie() *Trie {
	result := &Trie{root: &Node{}}
	return result
}

func (t *Trie) GetSuggestion2(query string, total int, c chan []string) { //edit
	// fmt.Println(query)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovery from Panic. Please DO NOT use any other characters apart from small case alphabets a-z and numbers 0-9 thank you", r)
			fmt.Println("Please resume running the program")
		}
	}()

	var result []string
	//move to next position node from the searching character
	currentNode := t.root //starts from root.

	r := []rune(query)
	for i := 0; i < len(query); i++ {
		switch r[i] {
		case 32: //represents value space on the keyboard
			r[i] = 123
			break
		case 39:
			r[i] = 124 //represents value ' on the keyboard
			break
		case 45: // represents value - on the keyboard
			r[i] = 125
			break
		case 48: //represents value 0 on the keyboard
			r[i] = 126
			break
		case 49: //represents value 1 on the keyboard
			r[i] = 127
			break
		case 50: //represents value 2 on the keyboard
			r[i] = 128
			break
		case 51: //represents value 3 on the keyboard
			r[i] = 129
			break
		case 52: //represents value 4 on the keyboard
			r[i] = 130
			break
		case 53: //represents value 5 on the keyboard
			r[i] = 131
			break
		case 54: //represents value 6 on the keyboard
			r[i] = 132
			break
		case 55: //represents value 7 on the keyboard
			r[i] = 133
			break
		case 56: //represents value 8 on the keyboard
			r[i] = 134
			break
		case 57: //represents value 9 on the keyboard
			r[i] = 135
			break
		}
		charIndex := r[i] - 'a'
		if currentNode.children[charIndex] == nil {
			// return result
			c <- result
			return

		}
		currentNode = currentNode.children[charIndex] //helps to move to the next Node.
	}

	if currentNode.isEnd && isLastNode(currentNode) { //if the current Node is the end

		result = append(result, query)
		c <- result
		return

	}

	wordList := []string{}

	if currentNode.isEnd {
		wordList = append(wordList, query)
		total--
	}

	if !isLastNode(currentNode) { //this returns a true if the current Node is pointing to other children nodes. So if no more children then correct; just return false
		_, result = Suggestion(query, wordList, total, currentNode)
	}

	c <- result
	return
}

func Suggestion(prefix string, wordList []string, repeat int, currentNode *Node) (int, []string) { //edit

	if isLastNode(currentNode) { //if there are children nodes
		if currentNode.isEnd && len(wordList) < 1 { //your current word can be an end; and it can also point to other children
			wordList = append(wordList, prefix)
		}
		return repeat, wordList
	}

	for i := 0; i < 39; i++ {
		if repeat < 1 {
			return repeat, wordList
		}
		nt := currentNode
		if currentNode.children[i] != nil { //this checks if there are any children nodes have any values.
			newl := ConvItoS(i) //function that gives int value but comes back with alphabetical representation of our customized tree struct
			prefix += string(newl)
			nt = nt.children[i]

			if nt.isEnd {
				wordList = append(wordList, prefix)
				repeat--
			}
			repeat, wordList = Suggestion(prefix, wordList, repeat, nt)
			prefix = prefix[0 : len(prefix)-1]
		}
	}

	return repeat, wordList

}

func isLastNode(nextNode *Node) bool { //this function checks for his friends.

	for i := 0; i < 39; i++ {
		if nextNode.children[i] != nil {
			return false
		}
	}
	return true
}

func ConvItoS(i int) string {

	array := [39]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", " ", "'", "-", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	return array[i]

}
