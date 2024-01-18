package robot

import "fmt"

type ToyRobot struct {
	state                RobotState
	mapDirectionsByTitle map[string]int
	mapDirectionsByValue map[int]string
}

func (t *ToyRobot) Place(x, y int, direction string) error {
	d := t.mapDirectionsByTitle[direction]

	if (x >= 0 && x <= 4) && (y >= 0 && y <= 4) && d != 0 {
		t.state.X = x
		t.state.Y = y
		t.state.Direction = d
		t.state.IsPlaced = true

		return nil
	} else {
		return fmt.Errorf("invalid robot coordinates or direction")
	}
}

func (t *ToyRobot) Report() string {
	return fmt.Sprintf("Output: %v,%v,%v", t.state.X, t.state.Y, t.mapDirectionsByValue[t.state.Direction])
}

func (t *ToyRobot) Left() {
	t.changeDirection(-1)
}

func (t *ToyRobot) Right() {
	t.changeDirection(1)
}

// remainder determines the new direction
// e.g., 2 (i.e., EAST) + 1 (i.e., move 90 deg to right) = 3 % 4 = 3 (SOUTH)
// e.g., 1 (i.e., NORTH) - 1 (i.e., move 90 deg to left) = 0 % 4 = 0 = 4 (WEST)
func (t *ToyRobot) changeDirection(value int) {
	d := (t.state.Direction + value) % 4

	if d == 0 {
		d = 4
	}

	t.state.Direction = d
}

func (t *ToyRobot) IsPlaced() bool {
	return t.state.IsPlaced
}

func (t *ToyRobot) GetState() RobotState {
	return t.state
}

func (t *ToyRobot) Move() error {
	// moving to north means adding one unit to the y coordinate
	// moving to south means subtracting one from the y coordinate
	// etc.
	directionMoveMap := map[int]struct {
		coordinate string
		value      int
	}{
		DirectionNorth: {"y", 1},
		DirectionEast:  {"x", 1},
		DirectionSouth: {"y", -1},
		DirectionWest:  {"x", -1},
	}

	// temporary variable (state) contains modified state
	direction := directionMoveMap[t.state.Direction]
	state := t.state
	if direction.coordinate == "x" {
		state.X += direction.value
	} else {
		state.Y += direction.value
	}

	// only apply the temporary state if the robot doesn't fall off the table
	if state.X >= 0 && state.X <= 4 && state.Y >= 0 && state.Y <= 4 {
		t.state = state
		return nil
	} else {
		return fmt.Errorf("the robot falls off the table")
	}
}

func (t *ToyRobot) Init() {
	t.state = RobotState{
		X:         -1,
		Y:         -1,
		Direction: -1,
	}

	t.mapDirectionsByTitle = map[string]int{
		DirectionNorthTitle: DirectionNorth,
		DirectionEastTitle:  DirectionEast,
		DirectionSouthTitle: DirectionSouth,
		DirectionWestTitle:  DirectionWest,
	}

	t.mapDirectionsByValue = map[int]string{
		DirectionNorth: DirectionNorthTitle,
		DirectionEast:  DirectionEastTitle,
		DirectionSouth: DirectionSouthTitle,
		DirectionWest:  DirectionWestTitle,
	}
}
