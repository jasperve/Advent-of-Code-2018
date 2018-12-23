package main

import (
	"fmt"
	"io/ioutil"
)

type coordinate struct {
	y          int
	x          int
	stepsTaken int
}

const (
	wall = 0
	room = 1
	door = 2
)

var grid map[int]map[int]int

func main() {

	grid = make(map[int]map[int]int)

	input, _ := ioutil.ReadFile("input.txt")

	row := make(map[int]int)
	row[0] = room
	grid[0] = row

	fillGrid(0, 0, input)

	//determine grid size
	minX, maxX, minY, maxY := 4294967295, 0, 4294967295, 0
	for y, row := range grid {
		for x := range row {
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	maxNumDoors, numRooms := 0, 0

	distances := findDistances(coordinate{y: 0, x: 0, stepsTaken: 0})
	for _, distance := range distances {
		if distance.stepsTaken/2 > maxNumDoors {
			maxNumDoors = distance.stepsTaken / 2
		}
		if distance.stepsTaken/2 >= 1000 && grid[distance.y][distance.x] == room {
			numRooms++
		}
	}

	fmt.Println("Longest distance: ", maxNumDoors)
	fmt.Println("Numbers of rooms which are at least 1000 doors away: ", numRooms)

}

func fillGrid(origY int, origX int, input []byte) {

	y, x := origY, origX

	for c := 0; c < len(input); c++ {

		switch input[c] {

		case 40: // (

			numBracketsOpened := 1

			for b := c + 1; b < len(input); b++ {
				if input[b] == 40 {
					numBracketsOpened++
				}
				if input[b] == 41 {
					numBracketsOpened--
				}
				if numBracketsOpened == 0 {
					fillGrid(y, x, input[c+1:b])
					c = b
					break
				}
			}

		case 78: // NORTH
			y--
			if _, ok := grid[y]; !ok {
				row := make(map[int]int)
				grid[y] = row
			}
			grid[y][x] = door
			y--
			if _, ok := grid[y]; !ok {
				row := make(map[int]int)
				grid[y] = row
			}
			grid[y][x] = room

		case 69: // EAST
			x++
			grid[y][x] = door
			x++
			grid[y][x] = room

		case 83: // SOUTH
			y++
			if _, ok := grid[y]; !ok {
				row := make(map[int]int)
				grid[y] = row
			}
			grid[y][x] = door
			y++
			if _, ok := grid[y]; !ok {
				row := make(map[int]int)
				grid[y] = row
			}
			grid[y][x] = room

		case 87: // WEST
			x--
			grid[y][x] = door
			x--
			grid[y][x] = room

		case 124: // |
			y = origY
			x = origX
		}

	}

}

func findDistances(startCoordinate coordinate) (closedList []coordinate) {

	openList := []coordinate{}
	openList = append(openList, startCoordinate)

	for len(openList) > 0 {

		//Get the heighest priority coordinate from the open list and add it to the closed list
		currentCoordinate := openList[0]
		openList = append([]coordinate{}, openList[1:]...)

		if grid[currentCoordinate.y][currentCoordinate.x] != wall {
			closedList = append(closedList, currentCoordinate)
		}

		for y := -1; y <= 1; y++ {
		XLOOP:
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {

					for c := 0; c < len(closedList); c++ {
						if closedList[c].x == currentCoordinate.x+x && closedList[c].y == currentCoordinate.y+y {
							continue XLOOP
						}
					}

					if grid[currentCoordinate.y+y][currentCoordinate.x+x] == wall {
						continue
					}

					newCoordinate := coordinate{
						x:          currentCoordinate.x + x,
						y:          currentCoordinate.y + y,
						stepsTaken: currentCoordinate.stepsTaken + 1,
					}

					openList = append(openList, newCoordinate)

				}
			}
		}

	}

	return closedList

}
