package processor

import (
	"alvinlucillo/toy-robot-challenge/internal/robot"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	MessageInvalidCommand       = "> Invalid command. Enter HELP for usage."
	MessageNotPlaced            = "> Oops. Robot not yet placed. Enter a PLACE command first."
	MessageInvalidPlace         = "> Invalid use of PLACE. Enter HELP for usage."
	MessageNotPlacedOutOfBounds = "> Robot not placed. It'll fall off the table."
	MessageNotMovedOutOfBounds  = "> Robot not moved. It'll fall off the table."
)

type StdinProcessor struct {
	source io.Reader
	robot  robot.Robot
	logger Logger
}

func (p *StdinProcessor) Init(src io.Reader, robot robot.Robot, logger Logger) {
	p.source = src
	p.robot = robot
	p.logger = logger
}

func (p *StdinProcessor) Process() error {

	p.robot.Init()

	p.logger.Println("Welcome to the Toy Robot Challenge Program! 🤖")
	p.logger.Println("Enter HELP for more details")
	p.logger.Println("\nEnter your commands below:")

	helpText := `Usage: COMMAND [ARGS]
    Commands:
     PLACE x,y,z  Places the robot on the table in position (x,y) facing z direction 
            - Args:   x - the X coordinate (valid values: 0-4)
                      y - the Y coordinate (valid values: 0-4)
                      z - the direction the robot is facing (valid values: NORTH,EAST,SOUTH,WEST)
            - Example: PLACE 1,2,NORTH
     LEFT         Rotates the robot 90 degrees to the left
     RIGHT        Rotates the robot 90 degrees to the right
     MOVE         Moves the robot one unit forward
     REPORT       Prints the robot's location (X,Y) and direction it's facing`

	scanner := bufio.NewScanner(p.source)

	commandsMap := map[string]bool{
		robot.CommandPlace:  true,
		robot.CommandLeft:   true,
		robot.CommandRight:  true,
		robot.CommandMove:   true,
		robot.CommandReport: true,
		robot.CommandHelp:   true,
	}

	directionsMap := map[string]bool{
		robot.DirectionNorthTitle: true,
		robot.DirectionEastTitle:  true,
		robot.DirectionSouthTitle: true,
		robot.DirectionWestTitle:  true,
	}

	// Processes each command until the end
	for {
		if scanner.Scan() {
			commandParts := strings.Split(strings.TrimSpace(scanner.Text()), " ")

			if found := commandsMap[commandParts[0]]; !found {
				p.logger.Println(MessageInvalidCommand)
				continue
			}

			if commandParts[0] != robot.CommandPlace && commandParts[0] != robot.CommandHelp {
				if !p.robot.IsPlaced() {
					p.logger.Println(MessageNotPlaced)
					continue
				}
			}

			switch commandParts[0] {
			case robot.CommandHelp:
				p.logger.Println(helpText)
			case robot.CommandPlace:
				if len(commandParts) != 2 {
					p.logger.Println(MessageInvalidPlace)
					continue
				}

				placeParts := strings.Split(commandParts[1], ",")
				if len(placeParts) != 3 {
					p.logger.Println(MessageInvalidPlace)
					continue
				}

				xValue, err := strconv.Atoi(placeParts[0])
				if err != nil {
					p.logger.Println(MessageInvalidPlace)
					continue
				}

				yValue, err := strconv.Atoi(placeParts[1])
				if err != nil {
					p.logger.Println(MessageInvalidPlace)
					continue
				}

				if _, found := directionsMap[placeParts[2]]; !found {
					p.logger.Println(MessageInvalidPlace)
					continue
				}

				if err := p.robot.Place(xValue, yValue, placeParts[2]); err != nil {
					p.logger.Println(MessageNotPlacedOutOfBounds)
				}
			case robot.CommandReport:
				p.logger.Println(fmt.Sprintf("> %s", p.robot.Report()))
			case robot.CommandMove:
				if err := p.robot.Move(); err != nil {
					p.logger.Println(MessageNotMovedOutOfBounds)
				}
			case robot.CommandLeft:
				p.robot.Left()
			case robot.CommandRight:
				p.robot.Right()
			}

		} else {
			if scanner.Err() != nil {
				return scanner.Err()
			}
			break
		}
	}

	return nil
}
