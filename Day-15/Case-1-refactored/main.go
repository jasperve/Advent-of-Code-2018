package main 

import (
	"fmt"
	"os"
	"bufio"
	"sort"
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
	
		for i, k := range playersKeys {

			if _, ok := players[k]; !ok || players[k].hitPoints <= 0 {
				continue
			}

			route := locateTarget(k)
			
			if len(route) > 1 {
				players[route[1]] = players[k]
				delete(players, k)
			}

			if len(route) == 1 || len(route) == 2 {

				if len(route) == 2 {
					k = route[1]
				}
				
				if attackTarget(k) {

					elfsFound := false
					goblinsFound := false

					totalHitPoints := 0
					for _, p := range players {
						if p.class == elf && p.hitPoints > 0 { elfsFound = true }
						if p.class == goblin && p.hitPoints > 0 { goblinsFound = true }
						if p.hitPoints > 0 { totalHitPoints += p.hitPoints }
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
		if _, ok := players[currentNode.location]; ok && currentNode.location != startNode.location && players[currentNode.location].hitPoints > 0 {
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

		for _, neighbour := range neighbours {

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
func attackTarget(attackerLocation coordinate) bool {

	targetLocations := []coordinate{}

	for _, neighbour := range neighbours {
		if _, ok := players[ coordinate{ attackerLocation.x + neighbour.x, attackerLocation.y + neighbour.y } ]; ok &&
			players[attackerLocation].class != players[ coordinate{ attackerLocation.x + neighbour.x, attackerLocation.y + neighbour.y } ].class &&
			players[ coordinate{ attackerLocation.x + neighbour.x, attackerLocation.y + neighbour.y } ].hitPoints > 0 { 
		
			targetLocations = append(targetLocations, coordinate{ attackerLocation.x + neighbour.x, attackerLocation.y + neighbour.y })
		
		}
	}

	if len(targetLocations) != 0 { 

		sort.Sort(byHitPointsAndXY(targetLocations))
		players[targetLocations[0]].hitPoints -= players[attackerLocation].attackPower
		if players[targetLocations[0]].hitPoints <= 0 {
			return true
		}

	}

	return false

}


func printCave() {

	for y := 0; y <= maxY; y++ {
		for x := 0; x <=maxX; x++ {
			if p, ok := players[ coordinate{ x, y }]; ok {
				if p.class == elf && p.hitPoints > 0 {
					fmt.Printf("E")
				} else if p.class == goblin && p.hitPoints > 0 {
					fmt.Printf("G")
				} else {
					fmt.Printf(".")
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