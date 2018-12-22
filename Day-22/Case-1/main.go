package main 

import(
	"fmt"
)

const (
	rocky = 0
	wet = 1
	narrow = 2

	caveDepth = 7740
	beginY = 0
	beginX = 0
	endY = 2000
	endX = 2000
	targetY = 763
	targetX = 12

)

type coordinate struct {
	class int
	index int
	erosion int
}

var cave map[int]map[int]*coordinate

func main() {

	cave = make(map[int]map[int]*coordinate)
	createCave(beginY, beginX, endY, endX)
	//printCave()
	fmt.Println("Total risk is:", calculateRisc())

}


func createCave(fromY int, fromX int, tillY int, tillX int) {

	for y := fromY; y <= tillY; y++ {

		if _, ok := cave[y]; !ok { 
			cave[y] = make(map[int]*coordinate)
		}

		if y == fromY {
			for x := fromX; x <= tillX; x++ {
				cave[y][x] = &coordinate{}
				if y == beginY && x == beginX {
					cave[y][x].index = 0
				} else if y == beginY && x > beginX { 
					cave[y][x].index = x * 16807
				} else {
					cave[y][x].index = cave[y-1][x].erosion * cave[y][x-1].erosion
				}
				cave[y][x].erosion = (cave[y][x].index + caveDepth) % 20183
				cave[y][x].class = cave[y][x].erosion % 3
			}
		} else if y > fromY && fromX <= tillX {
			cave[y][fromX] = &coordinate{}
			if y > beginY && fromX == beginX { 
				cave[y][fromX].index = y * 48271
			} else if y > beginY && fromX > beginX {
				cave[y][fromX].index = cave[y-1][fromX].erosion * cave[y][fromX-1].erosion
			}
			cave[y][fromX].erosion = (cave[y][fromX].index + caveDepth) % 20183
			cave[y][fromX].class =  cave[y][fromX].erosion % 3 
		}

	}

	if fromY < tillY || fromX < tillX {
		createCave(fromY+1, fromX+1, tillY, tillX)
	}

}


func calculateRisc() (result int) {
	for y := beginY; y <= targetY; y++ {
		for x := beginX; x <= targetX; x++ {
			if y != targetY || x != targetX { result += cave[y][x].class }
		}
	}
	return result
}


/*
func findRoute(startCoordinate coordinate, targetCoordinate coordinate, grid [][]object, players []player) []coordinate {
	
	openList := []coordinate{}
	openList = append(openList, startCoordinate)
	closedList := make(map[int]map[int]coordinate)

	LISTLOOP:
	for len(openList) > 0 {

		sort.Sort(byPriority(openList))
		
		//Get the heighest priority coordinate from the open list and add it to the closed list
		currentCoordinate := openList[0]
		openList = append([]coordinate{}, openList[1:]...)

		if _, ok := closedList[currentCoordinate.y]; !ok {
			closedList[currentCoordinate.y] = make(map[int]coordinate)
		}
		closedList[currentCoordinate.y][currentCoordinate.x] = currentCoordinate

		// If a wall is found stop processing this coordinate
		if grid[currentCoordinate.y][currentCoordinate.x].class == wall { continue LISTLOOP }

		// If a PLAYER is found standing in the way
		for p := 0; p < len(players); p++ {
			if players[p].y == currentCoordinate.y && players[p].x == currentCoordinate.x { continue LISTLOOP }
		}

		for y := -1; y <= 1; y++ {
			XLOOP:
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {
					
					if currentCoordinate.x + x == targetCoordinate.x && currentCoordinate.y + y == targetCoordinate.y {
						route := []coordinate{}
						for currentCoordinate.parent != nil {
							route = append(route, currentCoordinate)
							currentCoordinate = *currentCoordinate.parent
						}
						return route
					}

					// If the coordinate has already been marked as closed
					if _, ok := closedList[currentCoordinate.y + y][currentCoordinate.x + x]; ok {
						continue XLOOP
					}

					stepsToGo := int(math.Abs(float64((currentCoordinate.x+x)-targetCoordinate.x))) + int(math.Abs(float64((currentCoordinate.y+y)-targetCoordinate.y)))

					for o := 0; o < len(openList); o++ {
						if openList[o].x == currentCoordinate.x + x && openList[o].y == currentCoordinate.y + y {

							// Check if this coordinate has been reached before with more steps. If so update the coordinate in the open list
							if currentCoordinate.stepsTaken + stepsToGo < openList[o].stepsTaken + openList[o].stepsToGo {
								openList[o].parent = &currentCoordinate
								openList[o].priority = currentCoordinate.stepsTaken + 1 + stepsToGo
								openList[o].stepsTaken = currentCoordinate.stepsTaken + 1
								openList[o].stepsToGo = stepsToGo
								
							}
							continue XLOOP
															
						}
					}

					newCoordinate := coordinate {
						parent: &currentCoordinate,
						x: currentCoordinate.x + x, 
						y: currentCoordinate.y + y, 
						priority : currentCoordinate.stepsTaken + 1 + stepsToGo,
						stepsTaken: currentCoordinate.stepsTaken + 1,
						stepsToGo: stepsToGo,
					}

					openList = append(openList, newCoordinate)
					
				}
			}
		}

	}

	return nil

}
*/

func printCave() {
	for y := beginY; y <= endY; y++ {
		for x := beginX; x <= endX; x++ {
			switch cave[y][x].class {
			case rocky:
				fmt.Printf(".")
			case wet:
				fmt.Printf("=")
			case narrow:
				fmt.Printf("|")
			}
		}
		fmt.Printf("\n")
	}
}