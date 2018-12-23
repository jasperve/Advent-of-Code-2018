package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

const (
	EMPTY  = 0
	WALL   = 1
	ELF    = 2
	GOBLIN = 3
)

type coordinate struct {
	x           int
	y           int
	class       int
	attackPower int
	hitPoints   int
	step        int
}

type byPosition [][]*coordinate

func (c byPosition) Len() int {
	return len(c)
}
func (c byPosition) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byPosition) Less(i, j int) bool {
	return c[i][len(c)-1].y < c[j][len(c)-1].y || (c[i][len(c)-1].y == c[j][len(c)-1].y && c[i][len(c)-1].x < c[j][len(c)-1].x)
}

var coordinates map[int]map[int]*coordinate

func main() {

	coordinates = make(map[int]map[int]*coordinate)
	players := []*coordinate{}

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln("Unable to open input file")
	}

	y := 0
	input := bufio.NewScanner(file)
	for input.Scan() {

		row := make(map[int]*coordinate)
		for x, u := range input.Text() {

			newCoordinate := coordinate{x: x, y: y, class: EMPTY}

			switch u {
			case 35:
				newCoordinate.class = WALL
			case 69:
				newCoordinate.class = ELF
				newCoordinate.attackPower = 3
				newCoordinate.hitPoints = 200
				players = append(players, &newCoordinate)
			case 71:
				newCoordinate.class = GOBLIN
				newCoordinate.attackPower = 3
				newCoordinate.hitPoints = 200
				players = append(players, &newCoordinate)
			}

			row[x] = &newCoordinate

		}
		coordinates[y] = row
		y++

	}

	// While there are still Goblins and Elfs
	for {

		for p := 0; p < len(players); p++ {

			for o := 0; o < len(players); o++ {

				if players[p].class == players[o].class {
					continue
				}

				//From this point on we need to calculate which opponent is closest
				fmt.Printf("Calculating routes for %v on: %v, %v towards %v on %v, %v\n", players[p].class, players[p].x, players[p].y, players[o].class, players[o].x, players[o].y)

				leastSteps := 0
				possibleRoutes := [][]*coordinate{}
				calculateRoutes(players[p], players[p], players[o], &[]*coordinate{}, []*coordinate{}, &possibleRoutes, 0, &leastSteps)

				if len(possibleRoutes) > 0 {

					//Remove all the routes that are longer then the minimal route
					for pR := 0; pR < len(possibleRoutes); pR++ {
						if len(possibleRoutes[pR]) > leastSteps {
							possibleRoutes = append(possibleRoutes[:pR], possibleRoutes[pR+1:]...)
						}
					}

					//Sort the remaining possible routes by endpoint
					sort.Sort(byPosition(possibleRoutes))

					if len(possibleRoutes[0]) > 1 {
						fmt.Printf("Player on position x: %v, y: %v is moving to x: %v, y: %v\n\n", players[p].x, players[p].y, possibleRoutes[0][1].x, possibleRoutes[0][1].y)
						players[p].x = possibleRoutes[0][1].x
						players[p].y = possibleRoutes[0][1].y
						coordinates[players[p].y][players[p].x], coordinates[players[p].y][players[p].x] = possibleRoutes[0][1], players[p]
					}

				} else {
					fmt.Println("NO MORE ROUTES FOUND!")
					return
				}

				// Reset steps for each coordinate
				for i := 0; i < len(coordinates); i++ {
					for j := 0; j < len(coordinates[i]); j++ {
						coordinates[i][j].step = 0
					}
				}

			}

		}

		return

	}

}

func calculateRoutes(currentPosition *coordinate, startPosition *coordinate, targetPosition *coordinate, passedPositions *[]*coordinate, currentRoute []*coordinate, possibleRoutes *[][]*coordinate, step int, leastSteps *int) {

	*passedPositions = append(*passedPositions, currentPosition)

	if currentPosition.step == 0 || step <= currentPosition.step {
		currentPosition.step = step
	} else {
		return
	}

	if *leastSteps != 0 {
		if len(currentRoute)+int(math.Abs(float64(currentPosition.x-targetPosition.x)))+int(math.Abs(float64(currentPosition.y-targetPosition.y))) > *leastSteps {
			return
		}
	}

	if currentPosition == targetPosition {
		if len(currentRoute) < *leastSteps || *leastSteps == 0 {
			*leastSteps = len(currentRoute)
		}
		*possibleRoutes = append(*possibleRoutes, currentRoute)
		return
	}

	currentRoute = append(currentRoute, currentPosition)

	if len(currentRoute) > 1 && (currentPosition.class == WALL || currentPosition.class == ELF || currentPosition.class == GOBLIN) {
		return
	}

	directionPriority := []*coordinate{}
	if currentPosition.x < targetPosition.x && currentPosition.y < targetPosition.y {
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x+1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y+1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y-1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x-1])
	} else if currentPosition.x > targetPosition.x && currentPosition.y < targetPosition.y {
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x-1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y+1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y-1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x+1])
	} else if currentPosition.x == targetPosition.x && currentPosition.y < targetPosition.y {
		directionPriority = append(directionPriority, coordinates[currentPosition.y+1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x-1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x+1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y-1][currentPosition.x])
	} else if currentPosition.x == targetPosition.x && currentPosition.y > targetPosition.y {
		directionPriority = append(directionPriority, coordinates[currentPosition.y-1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x-1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x+1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y+1][currentPosition.x])
	} else if currentPosition.x < targetPosition.x && currentPosition.y > targetPosition.y {
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x+1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y-1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y+1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x-1])
	} else if currentPosition.x > targetPosition.x && currentPosition.y > targetPosition.y {
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x-1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y-1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y+1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x+1])
	} else if currentPosition.x < targetPosition.x && currentPosition.y == targetPosition.y {
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x-1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y-1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y+1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x+1])
	} else if currentPosition.x > targetPosition.x && currentPosition.y == targetPosition.y {
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x-1])
		directionPriority = append(directionPriority, coordinates[currentPosition.y-1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y+1][currentPosition.x])
		directionPriority = append(directionPriority, coordinates[currentPosition.y][currentPosition.x+1])
	}

	for dP := 0; dP < len(directionPriority); dP++ {
		nextPositionFound := false
		for cR := 0; cR < len(currentRoute); cR++ {
			if directionPriority[dP] == currentRoute[cR] {
				nextPositionFound = true
			}
		}
		for pP := 0; pP < len(*passedPositions); pP++ {
			if directionPriority[dP] == (*passedPositions)[pP] && step >= (*passedPositions)[pP].step {
				nextPositionFound = true
			}
		}
		if !nextPositionFound {
			calculateRoutes(directionPriority[dP], startPosition, targetPosition, passedPositions, currentRoute, possibleRoutes, step+1, leastSteps)
		}
	}

}
