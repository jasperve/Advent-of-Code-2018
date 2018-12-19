package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	//"math"
)

const (
	EMPTY = 0
	WALL = 1
	ELF = 2
	GOBLIN = 3
)

type coordinate struct {
	x int
	y int
	class int
	attackPower int
	hitPoints int
}

var coordinates map[int]map[int]*coordinate

func main() {

	coordinates = make(map[int]map[int]*coordinate)
	elfs := []*coordinate{}
	goblins := []*coordinate{}
	
	file, err := os.Open("input-test2.txt")
	if err != nil { log.Fatalln("Unable to open input file") }

	y := 0
	input := bufio.NewScanner(file)
	for input.Scan() {

		row := make(map[int]*coordinate)
		for x, u := range input.Text() {

			newCoordinate := coordinate{ x: x, y :y, class: EMPTY }

			switch u {
			case 35: 
				newCoordinate.class = WALL
			case 69:
				newCoordinate.class = ELF
				newCoordinate.attackPower = 3
				newCoordinate.hitPoints = 200
				elfs = append(elfs, &newCoordinate)
			case 71:
				newCoordinate.class = GOBLIN
				newCoordinate.attackPower = 3
				newCoordinate.hitPoints = 200
				goblins = append(goblins, &newCoordinate)
			}
			
			row[x] = &newCoordinate
			
		}
		coordinates[y] = row
		y++

	}

	// While there are still Goblins and Elfs
	for {

		for y := 0; y < len(coordinates); y++ {

			for x := 0; x < len(coordinates[y]); x++ {

				if coordinates[y][x].class == GOBLIN || coordinates[y][x].class == ELF {
		
					var opponents []*coordinate

					if coordinates[y][x].class == GOBLIN { 
						opponents = elfs 
					} else if coordinates[y][x].class == ELF { 
						opponents = goblins 
					}
					
					for o := 0; o < len(opponents); o++ {

						//From this point on we need to calculate which opponent is closest
						fmt.Printf("Calculating routes for goblin/elf on: %v, %v towards goblin/elf on %v, %v\n", coordinates[y][x].x, coordinates[y][x].y, opponents[o].x, opponents[o].y) 

						leastSteps := 0
						possibleRoutes := [][]*coordinate{}
						calculateRoutes(coordinates[y][x], coordinates[y][x], opponents[o], &[]*coordinate{}, []*coordinate{}, &possibleRoutes, &leastSteps)
					
						var shortestRoute int
						for pR := 0; pR < len(possibleRoutes); pR++ {
							if len(possibleRoutes[pR]) < shortestRoute || shortestRoute == 0 {
								shortestRoute = len(possibleRoutes[pR])
							}		
						}
					
						fmt.Println("Shortest route: ", shortestRoute)

						amountSameRoutes := 0

						for pR := 0; pR < len(possibleRoutes); pR++ {
							if len(possibleRoutes[pR]) == shortestRoute {
								fmt.Println("-------------------------")
								amountSameRoutes++
								for pRR := 0; pRR < len(possibleRoutes[pR]); pRR++ {
									fmt.Printf("%v, %v\n", possibleRoutes[pR][pRR].x, possibleRoutes[pR][pRR].y)
								}			
							}		
						}

						if(amountSameRoutes > 1) { fmt.Println(amountSameRoutes)}

					}

				}

			}

		}

		fmt.Println(len(elfs), len(goblins))
		
		break

	}

}

func calculateRoutes(currentPosition *coordinate, startPosition *coordinate, targetPosition *coordinate, passedPositions *[]*coordinate, currentRoute []*coordinate, possibleRoutes *[][]*coordinate, leastSteps *int) {

	*passedPositions = append(*passedPositions, currentPosition)

	/*if *leastSteps != 0 {
		if len(currentRoute) + int(math.Abs(float64(startPosition.x-targetPosition.x))) + int(math.Abs(float64(startPosition.y-targetPosition.y))) > *leastSteps { 
			return 
		}
	}*/

	if currentPosition == targetPosition {
		/*if len(currentRoute) < *leastSteps || *leastSteps == 0 {
			*leastSteps = len(currentRoute)
		}*/
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
			if directionPriority[dP] == currentRoute[cR] { nextPositionFound = true }
		}
		for pP := 0; pP < len(*passedPositions); pP++ {
			if directionPriority[dP] == (*passedPositions)[pP] { nextPositionFound = true }
		}
		if nextPositionFound == false {
			calculateRoutes(directionPriority[dP], startPosition, targetPosition, passedPositions, currentRoute, possibleRoutes, leastSteps)
		}
	}

}