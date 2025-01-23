package filetype

const (
	Map    byte = 1
	World  byte = 2
	Player byte = 3
)

func String(t byte) string {
	switch t {
	default:
		return "None"
	case Map:
		return "Map"
	case World:
		return "World"
	case Player:
		return "Player"
	}
}
