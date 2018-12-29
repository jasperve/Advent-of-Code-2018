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

var board map[coordinate]int
var players map[coordinate]*player

func main() {

	file, _ := os.Open("input.txt")

	board = make(map[coordinate]int)
	players = make(map[coordinate]*player)

	input := bufio.NewScanner(file)
	for y := 0; input.Scan(); y++ {
		for x, u := range input.Text() {
			switch u {
			case 35: // WALL FOUND
				board[coordinate{x,y}] = wall
				continue
			case 69: // ElF FOUND
				players[coordinate{x,y}] = &player{ class: elf, attackPower: 3, hitPoints: 200 }
			case 71: // GOBLIN FOUND
				players[coordinate{x,y}] = &player{ class: goblin, attackPower: 3, hitPoints: 200 }
			}
			board[coordinate{x,y}] = empty

		}
	}

	for {

		var playersKeys []coordinate
		for k := range players {
			playersKeys = append(playersKeys, k)
		}
		sort.Sort(byXY(playersKeys))
	
		for _, k := range playersKeys {
			
			targets := findTargets(k)
			fmt.Println(targets)

			break

		}

		break

	}

}


func findTargets(startLocation coordinate) []*node {
	
	neighbourNodes := []coordinate{ coordinate{x: 0, y: -1}, coordinate{x: -1, y: 0}, coordinate{x: 1, y: 0}, coordinate{x: 0, y: 1}}
	startNode := node { location: startLocation, steps: 0 }
	openList := []*node{}
	openList = append(openList, &startNode)
	closedList := make(map[coordinate]*node)
	targetList := []*node{}

	for len(openList) > 0 {

		currentNode := openList[0]
		openList = append([]*node{}, openList[1:]...)
		closedList[ coordinate { currentNode.location.x, currentNode.location.y } ] = currentNode

		// If we bump into a wall
		if board[currentNode.location] == wall {
			continue
		}

		// If we find a location that has a player on it
		if _, ok := players[currentNode.location]; ok {
			if currentNode.location == startNode.location {
				//Our own position
			} else if players[currentNode.location].class != players[startNode.location].class {
				// Opponent found (treas as wall)
				currentNode.target = players[currentNode.location]
				targetList = append(targetList, currentNode)
				continue			

			} else if players[currentNode.location].class == players[startNode.location].class {
				
				// Player of same class found (treat as wall)
				continue
			}			
		}

		NEIGHBOURLOOP:
		for _, neighbourNode := range neighbourNodes {

			// If the coordinate has already been marked as closed
			if foundNode, ok := closedList[ coordinate { currentNode.location.x + neighbourNode.x, currentNode.location.y + neighbourNode.y }]; ok {
				if currentNode.steps + 1 < foundNode.steps {
					foundNode.parent = currentNode
					foundNode.steps = currentNode.steps + 1
				}
				continue
			}

			for o := 0; o < len(openList); o++ {
				if openList[o].location.x == currentNode.location.x + neighbourNode.x && openList[o].location.y == currentNode.location.y + neighbourNode.y {

					// Check if this coordinate has been reached before with more steps. If so update the coordinate in the open list
					if currentNode.steps + 1 < openList[o].steps {
						openList[o].parent = currentNode
						openList[o].steps = currentNode.steps + 1
					}
					continue NEIGHBOURLOOP
				}
			}

			newNode := node { 
				parent: currentNode,
				location: coordinate { currentNode.location.x + neighbourNode.x, currentNode.location.y + neighbourNode.y },
				steps: currentNode.steps + 1,
			}

			openList = append(openList, &newNode)

		}

	}

	if len(targetList) > 0 {
		
		sort.Sort(byStepsAndXY(targetList))
		

	if len
	fmt.Println(targetList[0])
	/*parent := closedList[coordinate{22,9}].parent
	for parent != nil {
		fmt.Println(parent.location.x, parent.location.y)
		parent = parent.parent
	}*/

	return nil

}