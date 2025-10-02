// Package calendar provides utilities for generating and manipulating monthly and weekly calendars.
package calendar

import (
	"time"
)

type CalendarCell struct {
	Day            int
	Month          int
	Year           int
	IsCurrentMonth bool
	IsToday        bool
	Events         []Event
}

type Event struct{}

type CalendarMonth struct {
	Year  int
	Month string
	Cells [][]CalendarCell
}

func PrevYearMonth(year int, month time.Month) (prevYear int, prevMonth time.Month) {
	if month == time.January {
		return year - 1, time.December
	}
	return year, month - 1
}

func NextYearMonth(year int, month time.Month) (nextYear int, nextMonth time.Month) {
	if month == time.December {
		return year + 1, time.January
	}
	return year, month + 1
}

func IsToday(day, month, year int) bool {
	today := time.Now()
	return day == today.Day() && month == int(today.Month()) && year == today.Year()
}

func NewCalendarCell(day, month, year int, isCurrentMonth bool) CalendarCell {
	return CalendarCell{
		Day:            day,
		Month:          month,
		Year:           year,
		IsCurrentMonth: isCurrentMonth,
		IsToday:        IsToday(day, month, year),
		Events:         nil,
	}
}

func BuildWeek(startDay, month, year int, isCurrentMonth bool) []CalendarCell {
	week := []CalendarCell{}
	for i := 0; i < 7; i++ {
		cell := NewCalendarCell(startDay+i, month, year, isCurrentMonth)
		week = append(week, cell)
	}
	return week
}

func CalendarGen(year int, month time.Month) [][]CalendarCell {
	// Days in the current month, previous month, and next month
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	daysInPrevMonth := time.Date(year, month, 0, 0, 0, 0, 0, time.UTC).Day()
	prevYear, prevMonth := PrevYearMonth(year, month)
	nextYear, nextMonth := NextYearMonth(year, month)

	// Calculate the day of the week on which the month begins (0=Sunday)
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	weekDay := int(firstDay.Weekday())

	var weeks [][]CalendarCell

	// ---- FirstWeek ----
	week := []CalendarCell{}

	// Days of the previous month to fill in the first week
	firstDayInCalendar := daysInPrevMonth - weekDay + 1
	for i := 0; i < weekDay; i++ {
		day := firstDayInCalendar + i
		cell := NewCalendarCell(day, int(prevMonth), prevYear, false)
		week = append(week, cell)
	}
	// Days of the current month in the first week
	for i := weekDay; i < 7; i++ {
		day := i - weekDay + 1
		cell := NewCalendarCell(day, int(month), year, true)
		week = append(week, cell)
	}
	weeks = append(weeks, week)

	// ---- Full weeks of the current month ----
	currentDay := 7 - weekDay + 1
	for currentDay <= daysInMonth {
		week := []CalendarCell{}
		for i := 0; i < 7; i++ {
			var cell CalendarCell
			if currentDay <= daysInMonth {
				cell = NewCalendarCell(currentDay, int(month), year, true)
			} else {
				// Finally, fill in the days of the following month.
				fillerDay := currentDay - daysInMonth
				cell = NewCalendarCell(fillerDay, int(nextMonth), nextYear, false)
			}
			week = append(week, cell)
			currentDay++
		}
		weeks = append(weeks, week)
	}

	return weeks
}
