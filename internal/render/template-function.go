package render

import (
	"html/template"
	"math"
	"time"
)

var function = template.FuncMap{
	"shortDate":  ShortDate,
	"iterate":    Iterate,
	"floatToInt": FloatToInt,
	"substract":  Substract,
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
	return t.Format("2006-01-02")
}

func FloatToInt(f float32) int {
	y := math.Ceil(float64(f))
	return int(y)
}

func Substract(x, y int) int {
	return x - y
}
