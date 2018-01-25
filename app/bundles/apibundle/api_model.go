package apibundle

type APIAccessLevel int

const (
	LevelCheck 		APIAccessLevel = 1
	LevelCharge		APIAccessLevel = 2
	LevelEdit	 	APIAccessLevel = 3
)
