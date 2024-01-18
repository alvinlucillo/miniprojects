package processor

import (
	ro "alvinlucillo/toy-robot-challenge/internal/robot"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockLogger struct {
	logs []string
}

func (m *MockLogger) Println(args ...interface{}) {
	var log string
	for _, arg := range args {
		log += fmt.Sprint(arg)
	}

	m.logs = append(m.logs, log)
}

func TestProcess(t *testing.T) {
	testCases := map[string]struct {
		commands       []string
		expectedOutput []string
		expectedState  ro.RobotState
	}{
		"successful process - 1": {[]string{"PLACE 0,0,NORTH", "MOVE", "REPORT"}, []string{"> Output: 0,1,NORTH"},
			ro.RobotState{
				X:         0,
				Y:         1,
				Direction: ro.DirectionNorth,
				IsPlaced:  true,
			}},
		"successful process - 2": {[]string{"PLACE 0,0,NORTH", "LEFT", "REPORT"}, []string{"> Output: 0,0,WEST"},
			ro.RobotState{
				X:         0,
				Y:         0,
				Direction: ro.DirectionWest,
				IsPlaced:  true,
			}},
		"successful process - 3": {[]string{"PLACE 1,2,EAST", "MOVE", "MOVE", "LEFT", "MOVE", "REPORT"}, []string{"> Output: 3,3,NORTH"},
			ro.RobotState{
				X:         3,
				Y:         3,
				Direction: ro.DirectionNorth,
				IsPlaced:  true,
			}},
		"successful process - 4": {[]string{"PLACE 3,3,WEST", "MOVE", "RIGHT", "MOVE", "MOVE", "LEFT", "REPORT"}, []string{MessageNotMovedOutOfBounds, "> Output: 2,4,WEST"},
			ro.RobotState{
				X:         2,
				Y:         4,
				Direction: ro.DirectionWest,
				IsPlaced:  true,
			}},
		"failed process - place args out of bounds 1": {[]string{"PLACE -1,3,WEST"}, []string{MessageNotPlacedOutOfBounds},
			ro.RobotState{
				X:         -1,
				Y:         -1,
				Direction: -1,
				IsPlaced:  false,
			}},
		"failed process - place args out of bounds 2": {[]string{"PLACE 1,5,WEST"}, []string{MessageNotPlacedOutOfBounds},
			ro.RobotState{
				X:         -1,
				Y:         -1,
				Direction: -1,
				IsPlaced:  false,
			}},
		"failed process - move out of bounds": {[]string{"PLACE 4,3,EAST", "MOVE"}, []string{MessageNotMovedOutOfBounds},
			ro.RobotState{
				X:         4,
				Y:         3,
				Direction: ro.DirectionEast,
				IsPlaced:  true,
			}},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {

			robot := &ro.ToyRobot{}
			logger := &MockLogger{}
			processor := &StdinProcessor{}
			stdin := strings.NewReader(strings.Join(tc.commands, "\n"))

			processor.Init(stdin, robot, logger)
			processor.Process()

			require.Equal(t, tc.expectedState, robot.GetState(), "expected state should be equal")
			// first 3 logs (introduction) are ignored
			require.Equal(t, len(tc.expectedOutput), len(logger.logs)-3, "expected log count should be equal")

			for i, expectedOutput := range tc.expectedOutput {
				require.Equal(t, expectedOutput, logger.logs[i+3])
			}
		})
	}

}
