package model

func Levels() []Level                   { return []Level{LevelLight, LevelBalanced, LevelStrong} }
func ParseLevel(v string) (Level, bool) { l := Level(v); return l, l.Valid() }
func LevelWeight(l Level) int {
	switch l {
	case LevelLight:
		return 1
	case LevelBalanced:
		return 2
	case LevelStrong:
		return 3
	}
	return 0
}
func Stronger(a, b Level) Level {
	if LevelWeight(a) >= LevelWeight(b) {
		return a
	}
	return b
}
