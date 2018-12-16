package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"math"
	"sort"
	//"time"
)

const (
	EMPTY = 0
	WALL = 1
	ELF = 2
	GOBLIN = 3
)

type object struct {
	y int
	x int
	class int
	attackPower int
	hitPoints int
	opponent *object
}
type byDistance []object
func (o byDistance) Len() int {
	return len(o)
}
func (o byDistance) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
func (o byDistance) Less(i, j int) bool {
	iDistance := int(math.Abs(float64(o[i].x-o[i].opponent.x))) + int(math.Abs(float64(o[i].y-o[i].opponent.y)))
	jDistance := int(math.Abs(float64(o[j].x-o[j].opponent.x))) + int(math.Abs(float64(o[j].y-o[j].opponent.y)))
	return iDistance < jDistance
}

// Coordinate used to keep track of the shortest route with a sorting function to sort by priority to follow the steps
type coordinate struct {
    parent *coordinate 
	y int
	x int
	priority int
	stepsTaken int
	stepsToGo int
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

/*type byPosition [][]*object

func (c byPosition) Len() int {
	return len(c)
}
func (c byPosition) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byPosition) Less(i, j int) bool {
	return c[i][len(c)-1].y < c[j][len(c)-1].y || (c[i][len(c)-1].y == c[j][len(c)-1].y && c[i][len(c)-1].x < c[j][len(c)-1].x)
}*/

var objects map[int]map[int]*object

func main() {

	grid := [][]*object{}
	players := []object{}

	file, err := os.Open("input.txt")
	if err != nil { log.Fatalln("Unable to open input file") }

	input := bufio.NewScanner(file)
	for y := 0; input.Scan(); y++ {
		
		row := []*object{}
		for x, u := range input.Text() {
			
			newObject := object{ x: x, y: y, class: EMPTY }
			switch u {
			case 35: 
				newObject.class = WALL
			case 69:
				newObject.class = ELF
				newObject.attackPower = 3
				newObject.hitPoints = 200
				players = append(players, newObject)
			case 71:
				newObject.class = GOBLIN
				newObject.attackPower = 3
				newObject.hitPoints = 200
				players = append(players, newObject)
			}
			row = append(row, &newObject)
			
		}
		grid = append(grid, row)

	}

	//TODO : CREATE A LIST OF OPPONENTS SORTED BY LOCATION FROM FIGHTER
	// While there are still Goblins and Elfs
	for {

		for p := 0; p < len(players); p++ {

			opponents := []object{}

			for e := 0; e < len(players); e++ {
				if players[e].class != players[p].class {
					opponents = append(opponents, players[e])
				}
			}				

			for o := 0; o < len(opponents); o++ {
				opponents[o].opponent = &players[p]
			}

			sort.Sort(byDistance(opponents))

			shortestRoute := []coordinate{}

			for o := 0; o < len(opponents); o++ {

				shortestDistance := int(math.Abs(float64(opponents[o].x-players[p].x))) + int(math.Abs(float64(opponents[o].y-players[p].y)))

				if shortestDistance < len(shortestRoute) || len(shortestRoute) == 0 {
					fmt.Printf("Finding shortest route from %v, %v towards %v, %v\n", players[p].x, players[p].y, opponents[o].x, opponents[o].y)
					shortestRoute = findShortestRoute(players[p], opponents[o], grid)
				}
			}				

			for sR := 0; sR < len(shortestRoute); sR++ {
				fmt.Printf("x: %v, y: %v\n", shortestRoute[sR].x, shortestRoute[sR].y)
			}

		}

		return
	
	}

	//fmt.Printf("Finding shortest route from %v, %v towards %v, %v\n", players[0].x, players[0].y, players[27].x, players[27].y)
	//findShortestRoute(players[0], players[27], grid)


}



func findShortestRoute(startPosition object, targetPosition object, grid [][]*object) (route []coordinate) {

	startCoordinate := coordinate { x: startPosition.x, y: startPosition.y, priority: 0, stepsTaken: 0, stepsToGo: 0 }
	endCoordinate := coordinate { x: targetPosition.x, y: targetPosition.y }
	
	openList := []coordinate{}
	openList = append(openList, startCoordinate)
	closedList := make(map[int]map[int]coordinate)

	INFINITYLOOP:
	for len(openList) > 0 {

		sort.Sort(byPriority(openList))
		
		//Get the heighest priority coordinate from the open list and add it to the closed list
		currentCoordinate := openList[0]
		openList = append([]coordinate{}, openList[1:]...)

		if _, ok := closedList[currentCoordinate.y]; !ok {
			closedList[currentCoordinate.y] = make(map[int]coordinate)
		}
		closedList[currentCoordinate.y][currentCoordinate.x] = currentCoordinate

		for y := -1; y <= 1; y++ {
			XLOOP:
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {
					
					if currentCoordinate.x + x == endCoordinate.x && currentCoordinate.y + y == endCoordinate.y {
						fmt.Println("TARGET FOUND")
						fmt.Println("SHORTEST ROUTE BACKWARDS:")
						for currentCoordinate.parent != nil {
							route = append(route, currentCoordinate)
							currentCoordinate = *currentCoordinate.parent
						}
						break INFINITYLOOP

					}
			
					// IF a WALL, ELF or GOBLIN is found stop processing this coordinate
					if grid[currentCoordinate.y + y][currentCoordinate.x + x].class == WALL ||
					   grid[currentCoordinate.y + y][currentCoordinate.x + x].class == ELF ||
					   grid[currentCoordinate.y + y][currentCoordinate.x + x].class == GOBLIN {
						continue XLOOP
					}

					if _, ok := closedList[currentCoordinate.y + y][currentCoordinate.x +x]; ok {
						//fmt.Println("IN CLOSED LIST SO IGNORE")
						continue XLOOP
					}

					stepsToGo := int(math.Abs(float64((currentCoordinate.x+x)-endCoordinate.x))) + int(math.Abs(float64((currentCoordinate.y+y)-endCoordinate.y)))

					for o := 0; o < len(openList); o++ {
						if openList[o].x == x && openList[o].y == y {

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

		//time.Sleep(time.Second)
		
	}

	return route

}