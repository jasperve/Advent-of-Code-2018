package main

import (
	"fmt"
	"github.com/fatih/color"
//	"math"
//	"sort"
)

//11109, 731, 9 == 1008
const (
	rocky  = 0
	wet    = 1
	narrow = 2

	caveDepth = 11109 //7740 //510 //7740
	beginY    = 0
	beginX    = 0
	endY      = 800
	endX      = 80
	targetY   = 731 //10 //763
	targetX   = 9  //10 //12

	neither = 0
	torch   = 1
	gear    = 2
)

type region struct {
	class   int
	index   int
	erosion int
}

type coordinate struct {
	parent     *coordinate
	y          int
	x          int
	priority   int
	stepsTaken int
	stepsToGo  int
	equipment  int
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

type byLength [][]coordinate

func (c byLength) Len() int {
	return len(c)
}
func (c byLength) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byLength) Less(i, j int) bool {
	return len(c[i]) < len(c[j])
}

var cave map[int]map[int]*region

func main() {

	cave = make(map[int]map[int]*region)
	createCave(beginY, beginX, endY, endX)
	fmt.Println("Total risk is:", calculateRisc())
	findRoute()
	//printCave(findRoute())

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
				if y == targetY && x == targetX {
					cave[y][x].class = rocky
				} else {
					cave[y][x].class = cave[y][x].erosion % 3
				}
			}
		} else if y > fromY && fromX <= tillX {
			cave[y][fromX] = &region{}
			if y > beginY && fromX == beginX {
				cave[y][fromX].index = y * 48271
			} else if y > beginY && fromX > beginX {
				cave[y][fromX].index = cave[y-1][fromX].erosion * cave[y][fromX-1].erosion
			}
			cave[y][fromX].erosion = (cave[y][fromX].index + caveDepth) % 20183
			cave[y][fromX].class = cave[y][fromX].erosion % 3
		}

	}

	if fromY < tillY || fromX < tillX {
		createCave(fromY+1, fromX+1, tillY, tillX)
	}

}

func calculateRisc() (result int) {
	for y := beginY; y <= targetY; y++ {
		for x := beginX; x <= targetX; x++ {
			if y != targetY || x != targetX {
				result += cave[y][x].class
			}
		}
	}
	return result
}

func printCave(route []coordinate) {

	c := color.New(color.FgRed).Add(color.Underline)

	for y := beginY; y <= endY; y++ {
		for x := beginX; x <= endX; x++ {

			partOfRoute := false

			for r := 0; r < len(route); r++ {
				if route[r].y == y && route[r].x == x {
					partOfRoute = true
				}
			}

			switch cave[y][x].class {
			case rocky:
				if !partOfRoute {
					fmt.Printf(".")
				} else {
					c.Printf(".")
				}
			case wet:
				if !partOfRoute {
					fmt.Printf("=")
				} else {
					c.Printf("=")
				}
			case narrow:
				if !partOfRoute {
					fmt.Printf("|")
				} else {
					c.Printf("|")
				}
			}
		}
		fmt.Printf("\n")
	}
}

func findRoute() []coordinate {

	//routes := [][]coordinate{}

	startCoordinate := coordinate{y: 0, x: 0, equipment: torch}
	openList := []coordinate{}
	openList = append(openList, startCoordinate)
	closedList := []coordinate{}

	for len(openList) > 0 {

		//sort.Sort(byPriority(openList))

		//Get the heighest priority coordinate from the open list and add it to the closed list
		currentCoordinate := openList[0]
		openList = append([]coordinate{}, openList[1:]...)
		closedList = append(closedList, currentCoordinate)

		/*if currentCoordinate.y == targetY && currentCoordinate.x == targetX {

			// ROUTE FOUND
			fmt.Println(currentCoordinate.stepsTaken)
			route := []coordinate{}
			for currentCoordinate.parent != nil {
				route = append(route, currentCoordinate)
				currentCoordinate = *currentCoordinate.parent
			}
			routes = append(routes, route)
			return route

		}*/

		for y := -1; y <= 1; y++ {
		XLOOP:
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {

					if _, ok := cave[currentCoordinate.y+y][currentCoordinate.x+x]; !ok {
						continue XLOOP
					}

					equipment := currentCoordinate.equipment
					penalty := 0

					switch cave[currentCoordinate.y+y][currentCoordinate.x+x].class {
					case rocky:
						
						if currentCoordinate.equipment == neither {
							if cave[currentCoordinate.y][currentCoordinate.x].class == wet {
								equipment = gear
								if currentCoordinate.y+y == targetY && currentCoordinate.x+x == targetX {
									penalty += 7
								}
							} else if cave[currentCoordinate.y][currentCoordinate.x].class == narrow {
								equipment = torch
							}
							penalty += 7
						} else if currentCoordinate.equipment == gear && currentCoordinate.y+y == targetY && currentCoordinate.x+x == targetX {
							equipment = torch
							penalty += 7
						}
					case wet:
						if currentCoordinate.equipment == torch {
							if cave[currentCoordinate.y][currentCoordinate.x].class == rocky {
								equipment = gear
							} else if cave[currentCoordinate.y][currentCoordinate.x].class == narrow {
								equipment = neither
							}
							penalty += 7
						}
					case narrow:
						if currentCoordinate.equipment == gear {
							if cave[currentCoordinate.y][currentCoordinate.x].class == rocky {
								equipment = torch
							} else if cave[currentCoordinate.y][currentCoordinate.x].class == wet {
								equipment = neither
							}
							penalty += 7
						}
					}

					//stepsToGo := int(math.Abs(float64((currentCoordinate.y+y)-targetY))) + int(math.Abs(float64((currentCoordinate.x+x)-targetX)))

					for o := 0; o < len(openList); o++ {
						if openList[o].y == currentCoordinate.y+y && openList[o].x == currentCoordinate.x+x {
							// Check if this coordinate has been reached before with more steps. If so update the coordinate in the open list
							if currentCoordinate.stepsTaken + 1 + penalty < openList[o].stepsTaken {
								openList[o].parent = &currentCoordinate
								//openList[o].priority = currentCoordinate.stepsTaken + 1 + penalty + stepsToGo
								openList[o].stepsTaken = currentCoordinate.stepsTaken + 1 + penalty
								//openList[o].stepsToGo = stepsToGo
								openList[o].equipment = equipment

							}
							continue XLOOP
						}
					}

					for c := 0; c < len(closedList); c++ {
						if closedList[c].y == currentCoordinate.y+y && closedList[c].x == currentCoordinate.x+x {
							// Check if this coordinate has been reached before with more steps. If so update the coordinate and re-add it to the open list
							
							if currentCoordinate.stepsTaken + 1 + penalty < closedList[c].stepsTaken {
								closedList[c].parent = &currentCoordinate
								//closedList[c].priority = currentCoordinate.stepsTaken + 1 + penalty + stepsToGo
								closedList[c].stepsTaken = currentCoordinate.stepsTaken + 1 + penalty
								//closedList[c].stepsToGo = stepsToGo
								closedList[c].equipment = equipment
								openList = append(openList, closedList[c])
								closedList = append(closedList[:c], closedList[c+1:]...)
							}
							continue XLOOP
						}
					}

					newCoordinate := coordinate{
						parent:     &currentCoordinate,
						y:          currentCoordinate.y+y,
						x:          currentCoordinate.x+x,
						//priority:   currentCoordinate.stepsTaken + 1 + penalty + stepsToGo,
						stepsTaken: currentCoordinate.stepsTaken + 1 + penalty,
						//stepsToGo:  stepsToGo,
						equipment:  equipment,
					}

					openList = append(openList, newCoordinate)

				}
			}
		}

	}
	
	fmt.Println(len(closedList))
	for c := 0; c < len(closedList); c++ {
		if closedList[c].y == targetY && closedList[c].x == targetX {
			fmt.Println("hier", closedList[c].stepsTaken)
/*
			route := []coordinate{}
			for closedList[c].parent != nil {
				route = append(route, closedList[c])
				closedList[c] = *closedList[c].parent
			}
			return route
*/		
		}
	}

	return []coordinate{}

}
