package robot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlace(t *testing.T) {
	robot := &ToyRobot{}

	invalidPlaceArgErr := fmt.Errorf("invalid robot coordinates or direction")

	testCases := map[string]struct {
		x                 int
		y                 int
		direction         string
		expectedX         int
		expectedY         int
		expectedDirection string
		err               error
	}{
		"successful place":                {1, 1, DirectionEastTitle, 1, 1, DirectionEastTitle, nil},
		"failed place - negative x":       {-1, 1, DirectionEastTitle, -1, -1, "", invalidPlaceArgErr},
		"failed place - negative y":       {1, -1, DirectionEastTitle, -1, -1, "", invalidPlaceArgErr},
		"failed place - exceeded x limit": {5, 1, DirectionEastTitle, -1, -1, "", invalidPlaceArgErr},
		"failed place - exceeded y limit": {1, 5, DirectionEastTitle, -1, -1, "", invalidPlaceArgErr},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			robot.Init()

			err := robot.Place(tc.x, tc.y, tc.direction)

			if tc.err == nil {
				require.NoError(t, err, "place command should not return error")
			} else {
				require.Equal(t, tc.err, err)
			}

			require.Equal(t, robot.GetState().X, tc.expectedX, "x should be as expected")
			require.Equal(t, robot.GetState().Y, tc.expectedY, "y should be as expected")
			require.Equal(t, robot.mapDirectionsByValue[robot.GetState().Direction], tc.expectedDirection, "direction should be as expected")
		})
	}

}

func TestReport(t *testing.T) {
	robot := &ToyRobot{}
	robot.Init()

	require.Equal(t, robot.Report(), "Output: -1,-1,")

	robot.Place(1, 2, DirectionNorthTitle)
	require.Equal(t, robot.Report(), "Output: 1,2,NORTH")
}

func TestRight(t *testing.T) {
	robot := &ToyRobot{}
	robot.Init()

	robot.Place(0, 0, DirectionNorthTitle)

	directions := []int{DirectionEast, DirectionSouth, DirectionWest, DirectionNorth}
	for _, d := range directions {
		robot.Right()
		require.Equal(t, robot.GetState().Direction, d, "direction should be as expected")
	}
}

func TestLeft(t *testing.T) {
	robot := &ToyRobot{}
	robot.Init()

	robot.Place(0, 0, DirectionNorthTitle)

	directions := []int{DirectionWest, DirectionSouth, DirectionEast, DirectionNorth}
	for _, d := range directions {
		robot.Left()
		require.Equal(t, robot.GetState().Direction, d, "direction should be as expected")
	}
}

func TestMove(t *testing.T) {
	robot := &ToyRobot{}

	invalidMoveErr := fmt.Errorf("the robot falls off the table")

	testCases := map[string]struct {
		direction string
		expectedX int
		expectedY int
		err       error
	}{
		"successful move - east":  {DirectionEastTitle, 1, 0, nil},
		"successful move - north": {DirectionNorthTitle, 0, 1, nil},
		"failed move - south":     {DirectionSouthTitle, 0, 0, invalidMoveErr},
		"failed move - west":      {DirectionWestTitle, 0, 0, invalidMoveErr},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			robot.Init()

			robot.Place(0, 0, tc.direction)
			err := robot.Move()

			if tc.err == nil {
				require.NoError(t, err, "move command should not return error")
			} else {
				require.Equal(t, tc.err, err)
			}

			require.Equal(t, robot.GetState().X, tc.expectedX, "x should be as expected")
			require.Equal(t, robot.GetState().Y, tc.expectedY, "y should be as expected")
		})
	}

}
