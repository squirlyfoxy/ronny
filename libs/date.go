/*
	RON LIBRARYES

	Date, used by function execution

	Developed by: @squirlyfoxy
*/

package libs

import "time"

type date struct {
	Year  int
	Month int
	Day   int
}

func NewDate(year int, month int, day int) date {
	return date{year, month, day}
}

func CurrentDate() date {
	y := time.Now().Year()
	m := time.Now().Month()
	d := time.Now().Day()
	return NewDate(y, int(m), d)
}

func (d date) CalculateAge() int {
	months := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	current := CurrentDate()

	if d.Day > current.Day {
		current.Day += current.Day + months[current.Month-1]
		current.Month += current.Month - 1
	}
	if d.Month > current.Month {
		current.Month += current.Month + 12
		current.Year += current.Year - 1
	}
	return current.Year - d.Year
}
