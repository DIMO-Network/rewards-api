package points

type EngineType int

const (
	ICE EngineType = iota
	PHEV
	EV
)

type LevelInfo struct {
	MinWeeks int
	Points   int
}

var levelInfos = []LevelInfo{
	{MinWeeks: 0, Points: 0},
	{MinWeeks: 4, Points: 1000},
	{MinWeeks: 21, Points: 2000},
	{MinWeeks: 36, Points: 3000},
}

func findLevel(weeksConnected int) int {
	lastLevel := 0
	for level, levelInfo := range levelInfos {
		if weeksConnected < levelInfo.MinWeeks {
			break
		}
		lastLevel = level
	}
	return lastLevel
}

type Input struct {
	ExistingDisconnectionStreak int
	ExistingConnectionStreak    int
	ConnectedThisWeek           bool
}

type Output struct {
	DisconnectionStreak int
	ConnectionStreak    int
	Points              int
}

func ComputeStreak(i Input) Output {
	if i.ConnectedThisWeek {
		connStreak := i.ExistingConnectionStreak + 1
		return Output{
			DisconnectionStreak: 0,
			ConnectionStreak:    connStreak,
			Points:              levelInfos[findLevel(connStreak)].Points,
		}
	}

	connStreak := i.ExistingConnectionStreak
	discStreak := i.ExistingDisconnectionStreak + 1
	if discStreak == 3 {
		level := findLevel(connStreak)
		if level > 0 {
			level--
		}
		connStreak = levelInfos[level].MinWeeks
		discStreak = 0
	}
	return Output{
		ConnectionStreak:    connStreak,
		DisconnectionStreak: discStreak,
		Points:              0,
	}
}
