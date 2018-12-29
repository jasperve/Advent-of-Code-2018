package main

import (
	"fmt"
	"sort"
)

const (
	caveDepth = 7740
	beginY    = 0
	beginX    = 0
	targetY   = 763
	targetX   = 12
	endY      = 1000
	endX      = 100

	rocky  = 0
	wet    = 1
	narrow = 2

	torch = 0
	gear = 1
	neither = 2
)

type region struct {
	class   int
	index   int
	erosion int
}

type coordinate struct {
	y                 int
	x                 int
	stepsTaken        int
	equipment		  int
}

type coordinateSlice []*coordinate
func (c coordinateSlice) Len() int {
	return len(c)
}
func (c coordinateSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c coordinateSlice) Less(i, j int) bool {
	return c[i].stepsTaken < c[j].stepsTaken
}

var cave map[int]map[int]*region

func main() {

	cave = make(map[int]map[int]*region)
	createCave(beginY, beginX, endY, endX)
	fmt.Println("Number of minutes is:", findRoute())

}


// Fills the map with the correct region structs
func createCave(fromY int, fromX int, tillY int, tillX int) {

	for y := fromY; y <= tillY; y++ {

		if _, ok := cave[y]; !ok {
			cave[y] = make(map[int]*region)
		}

		if y == fromY {
			for x := fromX; x <= tillX; x++ {
				cave[y][x] = &region{}
				if (y == beginY && x == beginX) || (y == targetY && x == targetX) {
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
			if (y == beginY && fromX == beginX) || (y == targetY && fromX == targetX) {
				cave[y][fromX].index = 0
			} else if y > beginY && fromX == beginX {
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


// Tries to find the quickest route through the cave
func findRoute() int {

	startCoordinate := coordinate{y: 0, x: 0, stepsTaken: 0, equipment: torch }
	openList := make(map[string]*coordinate)
	openList[fmt.Sprintf("%v-%v-%v", startCoordinate.y, startCoordinate.x, startCoordinate.equipment)] = &startCoordinate
	closedList := make(map[string]*coordinate)

	for len(openList) > 0 {

		openListSlice := make(coordinateSlice, 0, len(openList))
		for _, c := range openList {
			openListSlice = append(openListSlice, c)
		}       
		sort.Sort(openListSlice)

		for ols := 0; ols < len(openListSlice); ols++ {
			fmt.Println(openListSlice[ols].y, openListSlice[ols].x, openListSlice[ols].stepsTaken, openListSlice[ols].equipment)
		}

		currentCoordinate := openListSlice[0]
		delete(openList, fmt.Sprintf("%v-%v-%v-%v", currentCoordinate.y, currentCoordinate.x, currentCoordinate.stepsTaken, currentCoordinate.equipment))
	
		if currentCoordinate.y == targetY && currentCoordinate.x == targetX && currentCoordinate.equipment == torch {
			return currentCoordinate.stepsTaken
		}

		if _, ok := closedList[fmt.Sprintf("%v-%v-%v", currentCoordinate.y, currentCoordinate.x, currentCoordinate.equipment)]; ok {
			continue
		}

		closedList[fmt.Sprintf("%v-%v-%v", currentCoordinate.y, currentCoordinate.x, currentCoordinate.equipment)] = currentCoordinate

		for equipment := 0; equipment <=2; equipment++ {
			if (cave[currentCoordinate.y][currentCoordinate.x].class == rocky && (equipment == torch || equipment == gear)) ||
				(cave[currentCoordinate.y][currentCoordinate.x].class == wet && (equipment == neither || equipment == gear)) ||
				(cave[currentCoordinate.y][currentCoordinate.x].class == narrow && (equipment == neither || equipment == torch)) {
			
				if _, ok := closedList[fmt.Sprintf("%v-%v-%v", currentCoordinate.y, currentCoordinate.x, equipment)]; ok {
					continue
				}

				newCoordinate := coordinate{
					y:                 currentCoordinate.y,
					x:                 currentCoordinate.x,
					stepsTaken:        currentCoordinate.stepsTaken + 7,
					equipment:		   equipment,
				}

				if v, ok := openList[fmt.Sprintf("%v-%v-%v-%v", currentCoordinate.y, currentCoordinate.x, currentCoordinate.stepsTaken + 7, equipment)]; ok {
					if v.stepsTaken < currentCoordinate.stepsTaken + 7 {
						continue
					}
				}
				openList[fmt.Sprintf("%v-%v-%v-%v", currentCoordinate.y, currentCoordinate.x, currentCoordinate.stepsTaken + 7, equipment)] = &newCoordinate	
			}

		}

		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {

					if _, ok := cave[currentCoordinate.y+y][currentCoordinate.x+x]; !ok {
						continue
					}

					if (cave[currentCoordinate.y+y][currentCoordinate.x+x].class == rocky && (currentCoordinate.equipment == torch || currentCoordinate.equipment == gear)) ||
					(cave[currentCoordinate.y+y][currentCoordinate.x+x].class == wet && (currentCoordinate.equipment == neither || currentCoordinate.equipment == gear)) ||
					(cave[currentCoordinate.y+y][currentCoordinate.x+x].class == narrow && (currentCoordinate.equipment == neither || currentCoordinate.equipment == torch)) {

						if _, ok := closedList[fmt.Sprintf("%v-%v-%v", currentCoordinate.y + y, currentCoordinate.x + x, currentCoordinate.equipment)]; ok {
							continue
						}

						newCoordinate := coordinate{
							y:                 currentCoordinate.y + y,
							x:                 currentCoordinate.x + x,
							stepsTaken:        currentCoordinate.stepsTaken + 1,
							equipment:		   currentCoordinate.equipment,
						}
						if v, ok := openList[fmt.Sprintf("%v-%v-%v-%v", currentCoordinate.y + y, currentCoordinate.x + x, currentCoordinate.stepsTaken + 1, currentCoordinate.equipment)]; ok {
							if v.stepsTaken < currentCoordinate.stepsTaken + 1 {
								continue
							}
						}
						openList[fmt.Sprintf("%v-%v-%v-%v", currentCoordinate.y + y, currentCoordinate.x + x, currentCoordinate.stepsTaken + 1, currentCoordinate.equipment)] = &newCoordinate	
						
					}

				}
			}
		}

	}

	return -1

}