package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"math"
	"sort"
	"github.com/fatih/color"
//	"time"
)

const (
	empty = 0
	wall = 1
	elf = 2
	goblin = 3
)

type player struct {
	x int
	y int
	class int
	attackPower int
	hitPoints int
	opponent *object
}

type byXY []player
func (p byXY) Len() int {
	return len(p)
}
func (p byXY) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p byXY) Less(i, j int) bool {
	return p[i].y < p[j].y || (p[i].y == p[j].y && p[i].x < p[j].x)
}

type byHitPoints []player
func (p byHitPoints) Len() int {
	return len(p)
}
func (p byHitPoints) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p byHitPoints) Less(i, j int) bool {
	return p[i].hitPoints < p[j].hitPoints
}

type object struct {
	y int
	x int
	class int
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

type byLenAndEndPoint [][]coordinate
func (c byLenAndEndPoint) Len() int {
	return len(c)
}
func (c byLenAndEndPoint) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byLenAndEndPoint) Less(i, j int) bool {
	if len(c[i]) < len(c[j]) { 
		return true
	} else if len(c[i]) == len(c[j]) && (c[i][len(c[i])-1].y < c[j][len(c[j])-1].y || (c[i][len(c[i])-1].y == c[j][len(c[j])-1].y && c[i][len(c[i])-1].x < c[j][len(c[j])-1].x)) {
		return true
	} else {
		return false
	}
}

type byLenAndBeginPoint [][]coordinate
func (c byLenAndBeginPoint) Len() int {
	return len(c)
}
func (c byLenAndBeginPoint) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byLenAndBeginPoint) Less(i, j int) bool {
	if len(c[i]) < len(c[j]) { 
		return true
	} else if len(c[i]) == len(c[j]) && (c[i][0].y < c[j][0].y || (c[i][0].y == c[j][0].y && c[i][0].x < c[j][0].x)) {
		return true
	} else {
		return false
	}
}

var objects map[int]map[int]object

func main() {

	grid := [][]object{}
	playersOriginal := []player{}
	
	file, err := os.Open("input.txt")
	if err != nil { log.Fatalln("Unable to open input file") }

	input := bufio.NewScanner(file)
	for y := 0; input.Scan(); y++ {
		
		row := []object{}
	
		for x, u := range input.Text() {
			
			newObject := object{ x: x, y: y, class: empty }
			
			switch u {
			case 35:
				newObject.class = wall
			case 69:
				newPlayer := player{ x: x, y: y, class: elf, attackPower: 3, hitPoints: 200 }
				playersOriginal = append(playersOriginal, newPlayer)
			case 71:
				newPlayer := player{ x: x, y: y, class: goblin, attackPower: 3, hitPoints: 200 }
				playersOriginal = append(playersOriginal, newPlayer)
			}

			row = append(row, newObject)
			
		}
		grid = append(grid, row)

	}

	attackPower := 3

	//TODO : CREATE A LIST OF OPPONENTS SORTED BY LOCATION FROM FIGHTER
	// While there are still Goblins and Elfs
	for {

		attackPower++

		players := make([]player, len(playersOriginal))
		
		for p := 0; p < len(playersOriginal); p++ {
			if playersOriginal[p].class == elf {
				playersOriginal[p].attackPower++
			}
		}
		
		fmt.Println("increasing attack power to", attackPower)

		copy(players, playersOriginal)

		round := 0

		OUTER:
		for {

			PLAYERSLOOP:
			for p := 0; p < len(players); p++ {

				opponents := make([]player, len(players))
				copy(opponents, players)
				sort.Sort(byXY(opponents))

				routes := [][]coordinate{}
				nextToOpponent := false

				RANGELOOP:
				for y := -1; y <= 1; y++ {
					for x := -1; x <= 1; x++ {
						if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {

							startCoordinate := coordinate { x: players[p].x+x, y: players[p].y+y, priority: 0, stepsTaken: 0, stepsToGo: 0 }
							var shortestRoute []coordinate

							OPPONENTSLOOP:
							for o := 0; o < len(opponents); o++ {

								if opponents[o].class == players[p].class { continue OPPONENTSLOOP }
								
								if opponents[o].x == players[p].x+x && opponents[o].y == players[p].y+y { 
									nextToOpponent = true
									break RANGELOOP	
								}

								// Calculate the MANHATTAN DISTANCE towards OPPONENT. If longer then SHORTEST ROUTE -> CONTINUE with NEXT OPPONENT
								shortestDistance := int(math.Abs(float64(opponents[o].x-players[p].x + x))) + int(math.Abs(float64(opponents[o].y-players[p].y + y)))
								if shortestRoute == nil || shortestDistance <= len(shortestRoute) + 10 {
									
									targetCoordinate := coordinate { x: opponents[o].x, y: opponents[o].y }

									// Start to FIND a ROUTE towards this OPPONENT
									route := findRoute(startCoordinate, targetCoordinate, grid, players)

									/*fmt.Println("from y: ", players[p].y+y, ",x:", players[p].x+x, " VV ")
									for t:= 0; t < len(route); t++ {
										fmt.Println(route[t].y, route[t].x)
									}*/

									// NO ROUTE has been found towards this OPPONENT
									if route == nil { continue OPPONENTSLOOP }

									// If this is the SHORTEST ROUTE found so far OR the FIRST ROUTE found so far -> UPDATE shortestRoute
									if shortestRoute == nil || len(route) <= len(shortestRoute) {
										shortestRoute = route
									}

								}
							}

							if shortestRoute != nil {
								shortestRoute = append([]coordinate { startCoordinate }, shortestRoute...)
								routes = append(routes, shortestRoute)
							}
						
						}
					}
				}

				// If no ROUTE has been found towards ANY OPPONENT and PLAYER is also NOT STANDING NEXT TO AN OPPONENT
				if !nextToOpponent && len(routes) != 0 { //continue PLAYERSLOOP }

					sort.Sort(byLenAndEndPoint(routes))

					for r := len(routes)-1; r >= 0; r-- {
						if routes[r][len(routes[r])-1].y != routes[0][len(routes[0])-1].y && routes[r][len(routes[r])-1].x != routes[0][len(routes[0])-1].x {
							routes = append(routes[:r], routes[r+1:]..
						}
					}

					sort.Sort(byLenAndBeginPoint(routes))

					//printBoard(grid, players, routes[0])
					//bufio.NewReader(os.Stdin).ReadBytes('\n') 

					players[p].x = routes[0][0].x
					players[p].y = routes[0][0].y

				}

				// IF PLAYER can attack -> CALL ATTACK FUNCTION -> IF RETURNED INDEX is -1 NO PLAYER was KILLED / IF RETURNED INDEX != -1 a PLAYER was KILLED and INDEX should be AMENDED
				if nextToOpponent == true || (len(routes) > 0 && len(routes[0]) == 1) {

					indexKilledPlayer := attackPlayer(players[p], &players) 
					
					if indexKilledPlayer == -2 { 
						break OUTER 
					} else if indexKilledPlayer != -1 {
						if indexKilledPlayer <= p { p-- }
					}

					elfsFound := false
					goblinsFound := false

					totalHitPointsLeft := 0
						for m := 0; m < len(players); m++ {
						if players[m].class == elf { elfsFound = true }
						if players[m].class == goblin { goblinsFound = true }
						totalHitPointsLeft += players[m].hitPoints
					}
				
					if !elfsFound || !goblinsFound {
						if p == len(players) - 1 {
							round++
						}
						fmt.Println("vv Board after round", round, "vv")
						fmt.Println("hitpoints left", totalHitPointsLeft)
						fmt.Println("fot a total of", totalHitPointsLeft * round)
						return
					}

					continue PLAYERSLOOP

				}


			}

			sort.Sort(byXY(players))
			round++

		}
	}
	
}


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


func attackPlayer(attacker player, players *[]player) int {

	neighbours := []player{}
	
	for y := -1; y <= 1; y++ {
		XLOOP:
		for x := -1; x <= 1; x++ {
			if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {
				
				for p := 0; p < len(*players); p++ {
					if (*players)[p].x == attacker.x + x && (*players)[p].y == attacker.y + y {

						// If a PLAYER is found in LOCATION but it has the SAME CLASS as the ATTACK -> CONTINUE
						if (*players)[p].class == attacker.class { continue XLOOP }
						neighbours = append(neighbours, (*players)[p])
						
					}
				}
			}
		}
	}

	if len(neighbours) != 0 { 

		sort.Sort(byHitPoints(neighbours))

		for n := len(neighbours)-1; n >= 0; n-- {
			if neighbours[n].hitPoints != neighbours[0].hitPoints {
				neighbours = append(neighbours[:n], neighbours[n+1:]...)
			}
		}

		sort.Sort(byXY(neighbours))

		for p := 0; p < len(*players); p++ {
			if (*players)[p].x == neighbours[0].x && (*players)[p].y == neighbours[0].y {

				// Attack THIS PLAYER and if HITPOINTS < 0 -> REMOVE PLAYER FROM LIST
				(*players)[p].hitPoints -= attacker.attackPower
				if (*players)[p].hitPoints < 0 {

					if (*players)[p].class == elf {
						fmt.Println("A ELF JUST DIED")
						return -2
					}
					*players = append((*players)[:p], (*players)[p+1:]...)
					return p
				}
			}
		}

	}

	return -1

}


func printBoard(grid [][]object, players []player, route []coordinate) {

	r := color.New(color.FgRed).Add(color.Underline)
	b := color.New(color.FgBlue).Add(color.Underline)
	c := color.New(color.FgGreen).Add(color.Underline)
	p := color.New(color.FgMagenta).Add(color.Underline)
	
	for y := 0; y < len(grid); y++ {
		XLOOP:
		for x := 0; x < len(grid[y]); x++ {

			for p := 0; p < len(players); p++ {
				if players[p].y == y && players[p].x == x {
					switch players[p].class {
					case elf:
						b.Printf("E")
					case goblin:
						r.Printf("G")
					}
					continue XLOOP
				}
			}

			for r := 0; r < len(route); r++ {
				if route[r].y == y && route[r].x == x {
					p.Printf("*")
					continue XLOOP
				}
			}

			switch grid[y][x].class {
			case empty:
				fmt.Printf(".")
			case wall:
				c.Printf("#")
			}

		}
		fmt.Printf("\n")
	}

}