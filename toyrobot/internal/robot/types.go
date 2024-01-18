package robot

const (
	_ int = iota
	DirectionNorth
	DirectionEast
	DirectionSouth
	DirectionWest
)

const (
	CommandMove   = "MOVE"
	CommandPlace  = "PLACE"
	CommandLeft   = "LEFT"
	CommandRight  = "RIGHT"
	CommandReport = "REPORT"
	CommandHelp   = "HELP"

	DirectionNorthTitle = "NORTH"
	DirectionEastTitle  = "EAST"
	DirectionSouthTitle = "SOUTH"
	DirectionWestTitle  = "WEST"
)

type RobotState struct {
	X         int  // x coordinate
	Y         int  // y coordinate
	Direction int  // NORTH, EAST, SOUTH, WEST
	IsPlaced  bool // robot placed on the table or not
}

type Robot interface {
	Init()
	Place(x, y int, direction string) error
	Move() error
	Left()
	Right()
	Report() string
	IsPlaced() bool
	GetState() RobotState
}
