package handler

import (
	"goRent/internal/config"
	"goRent/internal/model"
	"time"
)

func ListDatesFromRents(rents []model.Rent) []string {
	dates := []string{}
	for _, r := range rents {
		x := r.StartDate
		start := r.StartDate.Format(config.DateLayout)
		end := r.EndDate.Format(config.DateLayout)
		dates = append(dates, start)
		for start != end {
			x = x.AddDate(0, 0, 1)
			start = x.Format(config.DateLayout)
			dates = append(dates, start)
		}
	}
	return dates
}

func ListDates(start, end string) ([]string, error) {
	dates := []string{}
	x, err := time.Parse(config.DateLayout, start)
	if err != nil {
		return nil, err
	}
	dates = append(dates, start)
	for start != end {
		x = x.AddDate(0, 0, 1)
		start = x.Format(config.DateLayout)
		dates = append(dates, start)
	}
	return dates, nil
}

// includes compares a slice with string/s and returns true, if the instance of the string/s is in the slice
func Includes(s1 []string, s2 ...string) bool {
	for _, v1 := range s1 {
		for _, v2 := range s2 {
			if v1 == v2 {
				return true
			}
		}
	}
	return false
}
