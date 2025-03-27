package main

import (
	"fmt"
	"time"
)

func main() {
	startDateStr := "2023-06-01"
	endDateStr := "2023-07-30"
	hours := []string{"20:00:00"}

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)

	// Loop over the range of dates
	for currentDate := startDate; currentDate.Before(endDate) || currentDate.Equal(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		for _, hour := range hours {
			hourTime, _ := time.Parse("15:04:05", hour)
			combinedDateTime := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), hourTime.Hour(), hourTime.Minute(), hourTime.Second(), 0, currentDate.Location())

			// Create objects response example
			fmt.Println(combinedDateTime.Format("2006-01-02T15:04:05.000Z"))
		}
	}
}
