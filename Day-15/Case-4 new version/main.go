package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"math"
	"sort"
)

const (
	EMPTY = 0
	WALL = 1
	ELF = 2
	GOBLIN = 3
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

type byEndPointRoute [][]coordinate
func (c byEndPointRoute) Len() int {
	return len(c)
}
func (c byEndPointRoute) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byEndPointRoute) Less(i, j int) bool {
	if len(c[i]) < len(c[j]) || len(c[i]) == 0 {
		return true
	} else if len(c[i]) == len(c[j]) && len(c[j]) > 0 {
		return c[i][0].y < c[j][0].y || (c[i][0].y == c[j][0].y && c[i][0].x < c[j][0].x)
	} else {
		return false
	}
}

var objects map[int]map[int]object

func main() {

	grid := [][]object{}
	players := []player{}
	
	file, err := os.Open("input.txt")
	if err != nil { log.Fatalln("Unable to open input file") }

	input := bufio.NewScanner(file)
	for y := 0; input.Scan(); y++ {
		
		row := []object{}
	
		for x, u := range input.Text() {
			
			newObject := object{ x: x, y: y, class: EMPTY }
			
			switch u {
			case 35:
				newObject.class = WALL
			case 69:
				newPlayer := player{ x: x, y: y, class: ELF, attackPower: 3, hitPoints: 200 }
				players = append(players, newPlayer)
			case 71:
				newPlayer := player{ x: x, y: y, class: GOBLIN, attackPower: 3, hitPoints: 200 }
				players = append(players, newPlayer)
			}

			row = append(row, newObject)
			
		}
		grid = append(grid, row)

	}

	round := 0

	//TODO : CREATE A LIST OF OPPONENTS SORTED BY LOCATION FROM FIGHTER
	// While there are still Goblins and Elfs
	for {

		playerCanAttack := true

		PLAYERSLOOP:
		for p := 0; p < len(players); p++ {

			opponents := make([]player, len(players))
			copy(opponents, players)
			sort.Sort(byXY(opponents))

			lenShortestRoute := -1
			routes := [][]coordinate{}

			OPPONENTSLOOP:
			for o := 0; o < len(opponents); o++ {

				if opponents[o].class == players[p].class { continue }

				// Calculate the MANHATTAN DISTANCE towards OPPONENT. If longer then SHORTEST ROUTE -> CONTINUE with NEXT OPPONENT
				shortestDistance := int(math.Abs(float64(opponents[o].x-players[p].x))) + int(math.Abs(float64(opponents[o].y-players[p].y)))
				if shortestDistance <= lenShortestRoute || lenShortestRoute == -1 {
					
					// Start to FIND a ROUTE towards this OPPONENT
					route, routeFound := findRoute(players[p], opponents[o], grid, players)

					// NO ROUTE has been found towards this OPPONENT
					if !routeFound { continue OPPONENTSLOOP }

					// If this is the SHORTEST ROUTE found so far OR the FIRST ROUTE found so far -> ADD ROUTE to list and update shortest route length
					if len(route) <= lenShortestRoute || lenShortestRoute == -1 {
						routes = append(routes, route)
						lenShortestRoute = len(route)
					}

				}
			}				

			// If no ROUTE has been found towards ANY OPPONENT and PLAYER is also NOT STANDING NEXT TO AN OPPONENT -> CONTINUE with NEXT PLAYER
			if len(routes) == 0 { continue PLAYERSLOOP }
			

			// Remove ROUTES that were LONGER then the SHORTEST ROUTE found
			for r := len(routes)-1; r >= 0; r-- {
				if len(routes[r]) > lenShortestRoute { 
					routes = append(routes[:r], routes[r+1:]...)
				}
			}

			// If PLAYER is STANDING NEXT TO AN OPPONENT -> PLAYER can ATTACK || ELSE || If PLAYER can TAKE STEPS -> PLAYER moves 1 POSITION towards OPPONENT
			if len(routes[0]) == 0 {
				playerCanAttack = true
			} else {
				players[p].x = routes[0][len(routes[0])-1].x
				players[p].y = routes[0][len(routes[0])-1].y
			}

			// IF PLAYER can attack -> CALL ATTACK FUNCTION -> IF RETURNED INDEX is -1 NO PLAYER was KILLED / IF RETURNED INDEX != -1 a PLAYER was KILLED and INDEX should be AMENDED
			if playerCanAttack { 
				if indexKilledPlayer := attackPlayer(players[p], &players); indexKilledPlayer != -1 {
					if indexKilledPlayer <= p { p-- }
				}
			}

			elfsFound := false
			goblinsFound := false

			totalHitPointsLeft := 0
				for m := 0; m < len(players); m++ {
				if players[m].class == ELF { elfsFound = true }
				if players[m].class == GOBLIN { goblinsFound = true }
				totalHitPointsLeft += players[m].hitPoints
			}
		
			if !elfsFound || !goblinsFound {

				if p == len(players) - 1 {
					round++
				}
				fmt.Println("vv Board after round", round, "vv")
				printBoard(grid, players)
				fmt.Println("hitpoints left", totalHitPointsLeft)
				fmt.Println("fot a total of", totalHitPointsLeft * round)
				return
	
			}
			
		}
		/*
		printBoard(grid, players)
		time.Sleep(2000*time.Millisecond)
		fmt.Println()*/

		// Sort the remaining players by X-VALUE and Y-VALUE
		sort.Sort(byXY(players))
		round++

	}
	
}


func findRoute(startPosition player, targetPosition player, grid [][]object, players []player) (shortestRoute []coordinate, routeFound bool) {

	routes := [][]coordinate{}	
	lenShortestRoute := -1

	startCoordinate := coordinate { x: startPosition.x, y: startPosition.y, priority: 0, stepsTaken: 0, stepsToGo: 0 }
	endCoordinate := coordinate { x: targetPosition.x, y: targetPosition.y }
	
	openList := []coordinate{}
	openList = append(openList, startCoordinate)
	closedList := make(map[int]map[int]coordinate)

	for len(openList) > 0 {

		sort.Sort(byPriority(openList))
		
		//Get the heighest priority coordinate from the open list and add it to the closed list
		currentCoordinate := openList[0]
		openList = append([]coordinate{}, openList[1:]...)

		if _, ok := closedList[currentCoordinate.y]; !ok {
			closedList[currentCoordinate.y] = make(map[int]coordinate)
		}

		closedList[currentCoordinate.y][currentCoordinate.x] = currentCoordinate

		// If the current COORDINATE already has a LONGER ROUTE then the SHORTEST ROUTE
		if currentCoordinate.priority > lenShortestRoute && lenShortestRoute != -1 { continue}

		for y := -1; y <= 1; y++ {
			XLOOP:
			for x := -1; x <= 1; x++ {
				if (y == -1 && x == 0) || (y == 0 && (x == -1 || x == 1)) || (y == 1 && x == 0) {
					
					if currentCoordinate.x + x == endCoordinate.x && currentCoordinate.y + y == endCoordinate.y {
						route := []coordinate{}
						for currentCoordinate.parent != nil {
							route = append(route, currentCoordinate)
							currentCoordinate = *currentCoordinate.parent
						}
						lenShortestRoute = len(route)
						routes = append(routes, route)
						routeFound = true
						continue XLOOP
					}
			
					// If a WALL, ELF or GOBLIN is found stop processing this coordinate
					if grid[currentCoordinate.y + y][currentCoordinate.x + x].class == WALL {
						continue XLOOP
					}

					// If a player is found standing in the way
					for p := 0; p < len(players); p++ {
						if players[p].y == currentCoordinate.y + y && players[p].x == currentCoordinate.x + x {
							continue XLOOP
						}
					}

					// If the coordinate has already been marked as closed
					if _, ok := closedList[currentCoordinate.y + y][currentCoordinate.x + x]; ok {
						continue XLOOP
					}

					stepsToGo := int(math.Abs(float64((currentCoordinate.x+x)-endCoordinate.x))) + int(math.Abs(float64((currentCoordinate.y+y)-endCoordinate.y)))

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

					priority := currentCoordinate.stepsTaken + 1 + stepsToGo
					if y == -1 { priority = priority + 4 }
					if y == 0 && x == -1 { priority = priority + 3 }
					if y == 0 && x == 1 { priority = priority + 2 }
					if y == -1 { priority = priority + 1 }

					newCoordinate := coordinate {
						parent: &currentCoordinate,
						x: currentCoordinate.x + x, 
						y: currentCoordinate.y + y, 
						priority : priority,
						stepsTaken: currentCoordinate.stepsTaken + 1,
						stepsToGo: stepsToGo,
					}

					openList = append(openList, newCoordinate)
					
				}
			}
		}

	}

	// REMOVE ALL ROUTES THAT ARE TOO LONG
	for r:= len(routes)-1; r >= 0; r-- {
		if len(routes[r]) > lenShortestRoute {
			routes = append(routes[:r], routes[r+1:]...)
		}
	}

	if len(routes) > 1 {
		sort.Sort(byEndPointRoute(routes))
		return routes[0], routeFound
	} else if len(routes) == 1 {
		return routes[0], routeFound
	} else {
		return []coordinate{}, routeFound
	}

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
					*players = append((*players)[:p], (*players)[p+1:]...)
					return p
				}
			}
		}

	}

	return -1

}


func printBoard(grid [][]object, players []player) {

	for y := 0; y < len(grid); y++ {
		XLOOP:
		for x := 0; x < len(grid[y]); x++ {

			for p := 0; p < len(players); p++ {
				if players[p].y == y && players[p].x == x {
					switch players[p].class {
					case ELF:
						fmt.Printf("E")
					case GOBLIN:
						fmt.Printf("G")
					}
					continue XLOOP
				}
			}

			switch grid[y][x].class {
			case EMPTY:
				fmt.Printf(".")
			case WALL:
				fmt.Printf("#")
			}

		}
		fmt.Printf("\n")
	}

}