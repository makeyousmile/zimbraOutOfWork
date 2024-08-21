package main

import (
	"time"
)

// GetNextWeekendDates возвращает даты ближайших выходных (суббота и воскресенье)
func GetNextWeekendDates() (time.Time, time.Time) {
	now := time.Now()

	// Определяем сколько дней осталось до субботы
	daysUntilSaturday := (time.Saturday - now.Weekday() + 7) % 7
	if daysUntilSaturday == 0 {
		daysUntilSaturday = 7
	}

	// Определяем сколько дней осталось до воскресенья
	daysUntilSunday := daysUntilSaturday + 1

	// Получаем даты субботы и воскресенья
	saturday := now.AddDate(0, 0, int(daysUntilSaturday))
	sunday := now.AddDate(0, 0, int(daysUntilSunday))

	return saturday, sunday
}
