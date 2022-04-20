package libs

import (
	"strconv"
	"strings"
	"time"
)

type date struct {
	Year  int
	Month int
	Day   int
}

func NewDate(year int, month int, day int) date {
	return date{year, month, day}
}

func ToDate(_d string) date {
	//DD-MM-YYYY
	d := strings.Split(_d, "-")
	ye, _ := strconv.Atoi(d[2])
	mo, _ := strconv.Atoi(d[1])
	da, _ := strconv.Atoi(d[0])
	return NewDate(ye, mo, da)
}

func CurrentDate() date {
	y := time.Now().Year()
	m := time.Now().Month()
	d := time.Now().Day()
	return NewDate(y, int(m), d)
}

func (d date) CalculateAge() int {
	current := CurrentDate()
	int age = current.Year - d.Year
	if current.Month > d.Month {
		age++
	}

	return age
}
