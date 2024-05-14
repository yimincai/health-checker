package utils

import (
	"fmt"
	"time"
)

func TimeFormat(t time.Time) string {
	months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	weekdays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	chineseFormat := fmt.Sprintf("%d/%s/%d %s %02d:%02d:%02d",
		t.Year(),
		months[t.Month()-1],
		t.Day(),
		weekdays[t.Weekday()],
		t.Hour(),
		t.Minute(),
		t.Second())

	return chineseFormat
}

func MonthToInt(m time.Month) int {
	// Map month to integer representation
	monthMap := map[time.Month]int{
		time.January:   1,
		time.February:  2,
		time.March:     3,
		time.April:     4,
		time.May:       5,
		time.June:      6,
		time.July:      7,
		time.August:    8,
		time.September: 9,
		time.October:   10,
		time.November:  11,
		time.December:  12,
	}

	// Retrieve integer representation from the map
	if val, ok := monthMap[m]; ok {
		return val
	}

	panic("Invalid month")
}

func IsVaildateDate(y, m, d int) bool {
	//check year
	if y < 1 || y > 9999 {
		return false
	}

	//check month
	if m < 1 || m > 12 {
		return false
	}

	//check day
	if d < 1 || d > 31 {
		return false
	}

	//check day in month
	if d > 30 && (m == 4 || m == 6 || m == 9 || m == 11) {
		return false
	}

	if m == 2 {
		if y%400 == 0 || (y%4 == 0 && y%100 != 0) {
			if d > 29 {
				return false
			}
		} else {
			if d > 28 {
				return false
			}
		}
	}

	return true
}
