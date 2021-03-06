package main

import (
	"fmt"
	"sort"
)

const (
	rocky  = 0
	wet    = 1
	narrow = 2

	caveDepth = 11739 //7740
	beginY    = 0
	beginX    = 0
	endY      = 1000
	endX      = 100
	targetY   = 718 //763
	targetX   = 11 //12

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
	y                 int
	x                 int
	stepsTaken        int
	equipment		  int
}
type byStepsTaken []coordinate

func (c byStepsTaken) Len() int {
	return len(c)
}
func (c byStepsTaken) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byStepsTaken) Less(i, j int) bool {
	return c[i].stepsTaken < c[j].stepsTaken
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
	fmt.Println("Number of minutes is:", findRoute())

}

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

func findRoute() int {

	startCoordinate := coordinate{y: 0, x: 0, stepsTaken: 0, equipment: torch }
	openList := []coordinate{}
	openList = append(openList, startCoordinate)
	closedList := make(map[string]*coordinate)
	
	for len(openList) > 0 {
	
		sort.Sort(byStepsTaken(openList))
		currentCoordinate := openList[0]
		openList = append([]coordinate{}, openList[1:]...)
		
		if currentCoordinate.y == targetY && currentCoordinate.x == targetX && currentCoordinate.equipment == torch {
			return currentCoordinate.stepsTaken
		}

		if _, ok := closedList[fmt.Sprintf("%v-%v-%v", currentCoordinate.y, currentCoordinate.x, currentCoordinate.equipment)]; ok {
			continue
		}
		
		closedList[fmt.Sprintf("%v-%v-%v", currentCoordinate.y, currentCoordinate.x, currentCoordinate.equipment)] = &currentCoordinate

		for e := 0; e <=2; e++ {
			if (cave[currentCoordinate.y][currentCoordinate.x].class == rocky && (e == torch || e == gear)) ||
			(cave[currentCoordinate.y][currentCoordinate.x].class == wet && (e == neither || e == gear)) ||
			(cave[currentCoordinate.y][currentCoordinate.x].class == narrow && (e == neither || e == torch)) {
			
				newCoordinate := coordinate{
					y:                 currentCoordinate.y,
					x:                 currentCoordinate.x,
					stepsTaken:        currentCoordinate.stepsTaken + 7,
					equipment:		   e,
				}
				openList = append(openList, newCoordinate)	
			
			}
		}

		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {

					if _, ok := cave[currentCoordinate.y+y][currentCoordinate.x+x]; !ok {
						continue
					}

					if _, ok := closedList[fmt.Sprintf("%v-%v-%v", currentCoordinate.y+y, currentCoordinate.x+x, currentCoordinate.equipment)]; ok {
						continue
					}
	
					if (cave[currentCoordinate.y+y][currentCoordinate.x+x].class == rocky && (currentCoordinate.equipment == torch || currentCoordinate.equipment == gear)) ||
					(cave[currentCoordinate.y+y][currentCoordinate.x+x].class == wet && (currentCoordinate.equipment == neither || currentCoordinate.equipment == gear)) ||
					(cave[currentCoordinate.y+y][currentCoordinate.x+x].class == narrow && (currentCoordinate.equipment == neither || currentCoordinate.equipment == torch)) {

						newCoordinate := coordinate{
							y:                 currentCoordinate.y + y,
							x:                 currentCoordinate.x + x,
							stepsTaken:        currentCoordinate.stepsTaken + 1,
							equipment:		   currentCoordinate.equipment,
						}
						openList = append(openList, newCoordinate)
					
					}

				}
			}
		}

	}

	return -1

}