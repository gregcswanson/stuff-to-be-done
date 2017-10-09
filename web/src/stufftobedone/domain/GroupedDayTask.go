package domain

type GroupedDayTask struct {
	DateAsInt int
	Display string
	DayTasks []DayTask
}

