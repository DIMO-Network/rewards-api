package tasks

import "time"

var firstWeekStart = time.Date(2022, time.March, 14, 5, 0, 0, 0, time.UTC)

var weekDuration = 7 * 24 * time.Hour

func GetWeekNum(calculationTime time.Time) int {
	sinceStart := calculationTime.Sub(firstWeekStart)
	weekNum := int(sinceStart.Truncate(weekDuration)/weekDuration) - 1
	return weekNum
}
