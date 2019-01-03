package main 

import (
	"fmt"
	"os"
	"bufio"
	"sort"
	//"github.com/pkg/profile"
)

const (
	empty = 0
	wall = 1
	elf = 2
	goblin = 3
)

type coordinate struct {
	x int
	y int
}

type byXY []coordinate
func (c byXY) Len() int {
	return len(c)
}
func (c byXY) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byXY) Less(i, j int) bool {
	return c[i].y < c[j].y || (c[i].y == c[j].y && c[i].x < c[j].x)
}

type byHitPointsAndXY []coordinate
func (c byHitPointsAndXY) Len() int {
	return len(c)
}
func (c byHitPointsAndXY) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byHitPointsAndXY) Less(i, j int) bool {
	return players[c[i]].hitPoints < players[c[j]].hitPoints || ((players[c[i]].hitPoints == players[c[j]].hitPoints) && c[i].y < c[j].y || (c[i].y == c[j].y && c[i].x < c[j].x))
}

type player struct {
	class int
	attackPower int
	hitPoints int
}

type node struct {
	parent *node
	location coordinate
	steps int
	target *player
}

type byStepsAndXY []*node
func (c byStepsAndXY) Len() int {
	return len(c)
}
func (c byStepsAndXY) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byStepsAndXY) Less(i, j int) bool {
	return c[i].steps < c[j].steps || (c[i].steps == c[j].steps && (c[i].location.y < c[j].location.y || (c[i].location.y == c[j].location.y && c[i].location.x < c[j].location.x)))
}


var cave map[coordinate]int
var players map[coordinate]*player
var neighbours = []coordinate{ coordinate{x: 0, y: -1}, coordinate{x: -1, y: 0}, coordinate{x: 1, y: 0}, coordinate{x: 0, y: 1}}

var maxX, maxY int

func main() {

	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	cave = make(map[coordinate]int)
	players = make(map[coordinate]*player)

	analyseCave("input-ivo.txt")

	round := 0
	for {

		var playersKeys []coordinate
		for k := range players {
			playersKeys = append(playersKeys, k)
		}
		sort.Sort(byXY(playersKeys))
	
		for i := 0; i < len(playersKeys); i++ {

			route := locateTarget(playersKeys[i])
			
			if len(route) > 1 {
				players[route[1]] = players[playersKeys[i]]
				delete(players, playersKeys[i])
			}

			if len(route) == 1 || len(route) == 2 {

				if len(route) == 2 {
					playersKeys[i] = route[1]
				}
				
				killedCoordinate := attackTarget(playersKeys[i])
				if killedCoordinate.x != 0 && killedCoordinate.y != 0 {

					// If a player has been killed make sure we remove the key from the list
					for iPK, pK := range playersKeys {
						if pK.x == killedCoordinate.x && pK.y == killedCoordinate.y {
							playersKeys = append(playersKeys[:iPK], playersKeys[iPK+1:]...)
							if iPK <= i {
								i--
							}
						}
					}					

					elfsFound := false
					goblinsFound := false

					totalHitPoints := 0
					for _, p := range players {
						if p.class == elf { elfsFound = true }
						if p.class == goblin { goblinsFound = true }
						totalHitPoints += p.hitPoints
					}
					
					if !elfsFound || !goblinsFound {

						if i == len(playersKeys) - 1 {
							round++
						}

						fmt.Println("hitpoints left", totalHitPoints)
						fmt.Println("rounds done:", round)
						fmt.Println("hitpoints * round: ", totalHitPoints*round)
						return
					}

				}

			}

		}

		round++
		
	}

}


func analyseCave(fileLocation string) {

	file, _ := os.Open(fileLocation)
	input := bufio.NewScanner(file)
	for y := 0; input.Scan(); y++ {
		if y > maxY {
			maxY = y
		}
		for x, u := range input.Text() {
			if x > maxX {
				maxX = x
			}
			switch u {
			case 35:
				cave[coordinate{x,y}] = wall
				continue
			case 69:
				players[coordinate{x,y}] = &player{ class: elf, attackPower: 3, hitPoints: 200 }
			case 71:
				players[coordinate{x,y}] = &player{ class: goblin, attackPower: 3, hitPoints: 200 }
			}
			cave[coordinate{x,y}] = empty

		}
	}

}


func locateTarget(startLocation coordinate) (route []coordinate) {
	
	startNode := node { location: startLocation, steps: 0 }
	openList := []*node{}
	openList = append(openList, &startNode)
	closedList := make(map[coordinate]*node)
	shortestRoute := 0
	targetList := []*node{}
	openListPlace := 0

	for openListPlace < len(openList) {

		currentNode := openList[openListPlace]
		openListPlace++
		//openList = append([]*node{}, openList[1:]...)
		// If the coordinate is already in the closed list
		if _, ok := closedList[ coordinate { currentNode.location.x, currentNode.location.y }]; ok {
			continue
		}

		closedList[ coordinate { currentNode.location.x, currentNode.location.y } ] = currentNode

		// If the number of steps is higher then the shortest route found or we bump into a wall
		if (currentNode.steps > shortestRoute && shortestRoute != 0) || cave[currentNode.location] == wall {
			continue
		}

		// If we find a location that has a player on it
		if _, ok := players[currentNode.location]; ok && currentNode.location != startNode.location {
			if players[currentNode.location].class != players[startNode.location].class {
				currentNode.target = players[currentNode.location]
				currentNode.location = currentNode.parent.location
				targetList = append(targetList, currentNode)

				if currentNode.steps < shortestRoute {
					shortestRoute = currentNode.steps
				}
			}
			continue
		}

		//NEIGHBOURLOOP:
		for _, neighbour := range neighbours {
		
			// If the coordiante is already in the open list
			/*for o := 0; o < len(openList); o++ {
				if openList[o].location.x == currentNode.location.x + neighbour.x && openList[o].location.y == currentNode.location.y + neighbour.y {
					continue NEIGHBOURLOOP
				}
			}*/

			newNode := node { 
				parent: currentNode,
				location: coordinate { currentNode.location.x + neighbour.x, currentNode.location.y + neighbour.y },
				steps: currentNode.steps + 1,
			}

			openList = append(openList, &newNode)

		}

	}

	if len(targetList) > 0 {
		sort.Sort(byStepsAndXY(targetList))
		parent := targetList[0].parent
		for parent != nil {
			route = append([]coordinate{ coordinate{ parent.location.x, parent.location.y } }, route...)
			parent = parent.parent
		}
	}
		
	return route

}


// Returns the coordinate of the killed target
func attackTarget(attackerLocation coordinate) coordinate {

	targetLocations := []coordinate{}

	for _, neighbour := range neighbours {
		if _, ok := players[ coordinate{ attackerLocation.x + neighbour.x, attackerLocation.y + neighbour.y } ]; ok {
			if players[attackerLocation].class != players[ coordinate{ attackerLocation.x + neighbour.x, attackerLocation.y + neighbour.y } ].class { 
				targetLocations = append(targetLocations, coordinate{ attackerLocation.x + neighbour.x, attackerLocation.y + neighbour.y })
			}
		}
	}

	if len(targetLocations) != 0 { 

		sort.Sort(byHitPointsAndXY(targetLocations))
		players[targetLocations[0]].hitPoints -= players[attackerLocation].attackPower
		if players[targetLocations[0]].hitPoints <= 0 {
			delete(players, targetLocations[0])
			return targetLocations[0]
		}

	}

	return coordinate{}

}


func printCave() {

	for y := 0; y <= maxY; y++ {
		for x := 0; x <=maxX; x++ {
			if p, ok := players[ coordinate{ x, y }]; ok {
				if p.class == elf {
					fmt.Printf("E")
				} else if p.class == goblin {
					fmt.Printf("G")
				}
			} else {
				if cave[ coordinate{ x, y }] == wall {
					fmt.Printf("#")
				} else if cave[ coordinate{ x, y }] == empty {
					fmt.Printf(".")
				}
			}
		}
		fmt.Printf("\n")
	}

}