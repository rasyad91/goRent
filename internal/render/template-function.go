package render

import (
	"fmt"
	"goRent/internal/config"
	"goRent/internal/model"
	"html/template"
	"math"
	"time"
)

var function = template.FuncMap{
	"shortDate":        ShortDate,
	"iterate":          Iterate,
	"floatToInt":       FloatToInt,
	"substract":        Substract,
	"unprocessedRents": UnprocessedRents,
	"processedRents":   ProcessedRents,
	"multiply":         Multiply,
	"add":              Add,
	"totalCostInCart":  TotalCostInCart,
	"format2DP":        Format2DP,
}

// returns a slice of ints, starting at 1 going to count
func Iterate(count int) []int {
	var items []int

	for i := 1; i <= count; i++ {
		items = append(items, i)
	}
	return items
}

// Short date returns time in DD-MM-YYYY format
func ShortDate(t time.Time) string {
	return t.Format(config.DateLayout)
}

func FloatToInt(f float32) int {
	y := math.Ceil(float64(f))
	return int(y)
}

func Substract(x, y int) int {
	return x - y
}

func Add(x ...float32) float32 {
	var result float32
	for _, v := range x {
		result = result + v
	}
	return result
}

func UnprocessedRents(rents []model.Rent) []model.Rent {
	unprocessed := []model.Rent{}
	for _, r := range rents {
		if !r.Processed {
			unprocessed = append(unprocessed, r)
		}
	}
	return unprocessed
}

func ProcessedRents(rents []model.Rent) []model.Rent {
	unprocessed := []model.Rent{}
	for _, r := range rents {
		if r.Processed {
			unprocessed = append(unprocessed, r)
		}
	}
	return unprocessed
}

func Multiply(x int, y float32) float32 {
	result := float32(x) * y
	return result
}

func TotalCostInCart(rents []model.Rent) float32 {
	var total float32
	for _, x := range rents {
		total = total + x.TotalCost
	}
	return total
}

func Format2DP(x float32) string {
	return fmt.Sprintf("%.2f", x)
}
