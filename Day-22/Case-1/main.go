package main 

import(
	"fmt"
	"sort"
	"math"
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

	neither = 0
	torch = 1
	gear = 2
)

type region struct {
	class int
	index int
	erosion int
}

type coordinate struct {
    parent *coordinate 
	y int
	x int
	priority int
	stepsTaken int
	stepsToGo int
	equipment int
}
type byPriority []coordinate
func (c byPriority) Len() int {
	return len(c)
}
func (c byPriority) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byPriority) Less(i, j int) bool {
	return c[i].priority < c[j].priority
}

var cave map[int]map[int]*region

func main() {

	cave = make(map[int]map[int]*region)
	createCave(beginY, beginX, endY, endX)
	//printCave()
	fmt.Println("Total risk is:", calculateRisc())

}


func createCave(fromY int, fromX int, tillY int, tillX int) {

	for y := fromY; y <= tillY; y++ {

		if _, ok := cave[y]; !ok { 
			cave[y] = make(map[int]*region)
		}

		if y == fromY {
			for x := fromX; x <= tillX; x++ {
				cave[y][x] = &region{}
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
			cave[y][fromX] = &region{}
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



func findRoute(beginY int, beginX int, targetY int, targetX int) int {
	
	startCoordinate := coordinate{ y: 0, x:0, equipment: torch }

	openList := []coordinate{}
	openList = append(openList, startCoordinate)
	closedList := []coordinate{}

	for len(openList) > 0 {

		sort.Sort(byPriority(openList))
		
		//Get the heighest priority coordinate from the open list and add it to the closed list
		currentCoordinate := openList[0]
		openList = append([]coordinate{}, openList[1:]...)
		closedList = append(closedList, currentCoordinate)

		for y := -1; y <= 1; y++ {
			XLOOP:
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {
					
					penalty := 0
					switch cave[currentCoordinate.y+y][currentCoordinate.x+x].class {
					case rocky:
						if currentCoordinate.equipment == neither { penalty = 7 }
					case wet:
						if currentCoordinate.equipment == torch { penalty = 7 }
					case narrow:
						if currentCoordinate.equipment == gear { penalty = 7 }
					}
					
					if currentCoordinate.y + y == targetY && currentCoordinate.x + x == targetX {
						
						// ROUTE FOUND 

						/*route := []coordinate{}
						for currentCoordinate.parent != nil {
							route = append(route, currentCoordinate)
							currentCoordinate = *currentCoordinate.parent
						}
						return route*/

					}

					for c := 0; c < len(closedList); c++ {
						if closedList[c].y == currentCoordinate.y + y && closedList[c].x == currentCoordinate.x + x {

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

					stepsToGo := int(math.Abs(float64((currentCoordinate.y+y)-targetY))) + int(math.Abs(float64((currentCoordinate.x+x)-targetX)))

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
