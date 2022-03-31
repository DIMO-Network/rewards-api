package services

type LevelInfo struct {
	Level    int
	MinWeeks int
	Points   int
}

var levelInfos = []*LevelInfo{
	{Level: 1, MinWeeks: 0, Points: 0},
	{Level: 2, MinWeeks: 4, Points: 1000},
	{Level: 3, MinWeeks: 21, Points: 2000},
	{Level: 4, MinWeeks: 36, Points: 3000},
}

func GetLevel(connStreak int) *LevelInfo {
	levelInfo := levelInfos[0]
	for _, nextLevelInfo := range levelInfos {
		if connStreak < nextLevelInfo.MinWeeks {
			break
		}
		levelInfo = nextLevelInfo
	}
	return levelInfo
}

type StreakInput struct {
	ConnectedThisWeek           bool
	ExistingConnectionStreak    int
	ExistingDisconnectionStreak int
}

type StreakOutput struct {
	DisconnectionStreak int
	ConnectionStreak    int
	Points              int
}

func ComputeStreak(i StreakInput) StreakOutput {
	if i.ConnectedThisWeek {
		connStreak := i.ExistingConnectionStreak + 1
		return StreakOutput{
			ConnectionStreak:    connStreak,
			DisconnectionStreak: 0,
			Points:              GetLevel(connStreak).Points,
		}
	}

	connStreak := i.ExistingConnectionStreak
	discStreak := i.ExistingDisconnectionStreak + 1
	if discStreak%3 == 0 {
		levelIndex := GetLevel(connStreak).Level - 1
		if levelIndex > 0 {
			levelIndex--
		}
		connStreak = levelInfos[levelIndex].MinWeeks
	}
	return StreakOutput{
		ConnectionStreak:    connStreak,
		DisconnectionStreak: discStreak,
		Points:              0,
	}
}
