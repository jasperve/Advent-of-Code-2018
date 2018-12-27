package main

import (
	"fmt"
	"sort"
)

//11109, 731, 9 == 1008
const (
	rocky  = 0
	wet    = 1
	narrow = 2

	caveDepth = 510 //7740
	beginY    = 0
	beginX    = 0
	endY      = 12 //800 //12
	endX      = 12 //50 //12
	targetY   = 10 //763
	targetX   = 10 //12

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
	possibleEquipment map[int]struct{}
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
	findRoute()

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
				} else if y == targetY && x == targetX {
					cave[y][x].index = 0
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

func findRoute() []coordinate {

	startCoordinate := coordinate{y: 0, x: 0, possibleEquipment: map[int]struct{}{ torch: struct{}{} } }
	openList := []coordinate{}
	openList = append(openList, startCoordinate)
	closedList := []coordinate{}

	for len(openList) > 0 {

		sort.Sort(byStepsTaken(openList))

		//Get the heighest priority coordinate from the open list and add it to the closed list
		currentCoordinate := openList[0]
		openList = append([]coordinate{}, openList[1:]...)
		closedList = append(closedList, currentCoordinate)

		for y := -1; y <= 1; y++ {
		XLOOP:
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {

					if _, ok := cave[currentCoordinate.y+y][currentCoordinate.x+x]; !ok {
						continue XLOOP
					}

					for c := 0; c < len(closedList); c++ {
						if closedList[c].y == currentCoordinate.y+y && closedList[c].x == currentCoordinate.x+x {
							continue XLOOP
						}
					}

					possibleEquipment := map[int]struct{}{}
					penalty := 0

					
					if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 1", currentCoordinate.possibleEquipment)	}

					switch cave[currentCoordinate.y+y][currentCoordinate.x+x].class {
					case rocky:
						if cave[currentCoordinate.y][currentCoordinate.x].class == rocky {
							possibleEquipment = currentCoordinate.possibleEquipment
							if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 10")	}
						} else if cave[currentCoordinate.y][currentCoordinate.x].class == wet {
							if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 2")	}
							possibleEquipment[gear] = struct{}{}
							if !possibleEquipmentContains(currentCoordinate.possibleEquipment, gear) {
								penalty += 7
							}
						} else if cave[currentCoordinate.y][currentCoordinate.x].class == narrow {
							if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 3")	}
							possibleEquipment[torch] = struct{}{}
							if !possibleEquipmentContains(currentCoordinate.possibleEquipment, torch) {
								penalty += 7
							}
						}
						if !possibleEquipmentContains(possibleEquipment, torch) && currentCoordinate.y+y == targetY && currentCoordinate.x+x == targetX {
							if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 4")	}
							possibleEquipment[torch] = struct{}{}
							penalty += 7
						}
					case wet:
						if cave[currentCoordinate.y][currentCoordinate.x].class == wet {
							possibleEquipment = currentCoordinate.possibleEquipment
						} else if cave[currentCoordinate.y][currentCoordinate.x].class == rocky {
							if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 5")	}
							possibleEquipment[gear] = struct{}{}
							if !possibleEquipmentContains(currentCoordinate.possibleEquipment, gear) {
								penalty += 7
							}
						} else if cave[currentCoordinate.y][currentCoordinate.x].class == narrow {
							if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 6")	}
							possibleEquipment[neither] = struct{}{}
							if !possibleEquipmentContains(currentCoordinate.possibleEquipment, neither) {
								penalty += 7
							}
						}
					case narrow:
						if cave[currentCoordinate.y][currentCoordinate.x].class == narrow {
							possibleEquipment = currentCoordinate.possibleEquipment
						} else if cave[currentCoordinate.y][currentCoordinate.x].class == rocky {
							if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 7")	}
							possibleEquipment[torch] = struct{}{}
							if !possibleEquipmentContains(currentCoordinate.possibleEquipment, torch) {
								penalty += 7
							}
						} else if cave[currentCoordinate.y][currentCoordinate.x].class == wet {
							possibleEquipment[neither] = struct{}{}
							if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier: 8")	}
							if !possibleEquipmentContains(currentCoordinate.possibleEquipment, neither) {
								penalty += 7
							}
						}
					}

					for o := 0; o < len(openList); o++ {
						if openList[o].y == currentCoordinate.y+y && openList[o].x == currentCoordinate.x+x {
							if currentCoordinate.stepsTaken+1+penalty < openList[o].stepsTaken {
								openList[o].stepsTaken = currentCoordinate.stepsTaken + 1 + penalty
								openList[o].possibleEquipment = possibleEquipment
							} else if currentCoordinate.stepsTaken+1+penalty == openList[o].stepsTaken {
								if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier dan?")	}
								for k, v := range possibleEquipment {
									if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println("hier blaat:", k)	}
									openList[o].possibleEquipment[k] = v
								}
							}
							continue XLOOP
						}
					}

					newCoordinate := coordinate{
						y:                 currentCoordinate.y + y,
						x:                 currentCoordinate.x + x,
						stepsTaken:        currentCoordinate.stepsTaken + 1 + penalty,
						possibleEquipment: possibleEquipment,
					}

					if currentCoordinate.y+y == 11 &&currentCoordinate.x+x == 10 { 	fmt.Println(newCoordinate)	}
					openList = append(openList, newCoordinate)

				}
			}
		}

	}

	
	for y := 0; y <= endY; y++ {
		for x := 0; x <= endX; x++ {
			for c := 0; c < len(closedList); c++ {
				if closedList[c].y == y && closedList[c].x == x {
					fmt.Printf("%v (", closedList[c].stepsTaken)
					for k := range closedList[c].possibleEquipment {
						fmt.Printf("%v", k)
					}
					fmt.Printf(")\t")
				}
			}
		}
		fmt.Printf("\n")
	}
	
	for c := 0; c < len(closedList); c++ {
		if closedList[c].y == 11 && closedList[c].x == 10 {
			fmt.Printf("HIER: %v (", closedList[c].stepsTaken)
			for k := range closedList[c].possibleEquipment {
				fmt.Printf("%v", k)
			}
			fmt.Printf(")\t")
		}
	}



	for c := 0; c < len(closedList); c++ {
		if closedList[c].y == targetY && closedList[c].x == targetX {
			fmt.Println("steps needed:", closedList[c].stepsTaken)
		}
	}

	return []coordinate{}

}

func possibleEquipmentContains(currentEquipment map[int]struct{}, equipment int) bool {
	if _, ok := currentEquipment[equipment]; ok {
		return true
	}
	return false
}
