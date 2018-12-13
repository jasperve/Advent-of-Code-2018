package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	NORTH = 360
	EAST  = 90
	SOUTH = 180
	WEST  = 270
)

const (
	LEFT     = -2
	STRAIGHT = -1
	RIGHT    = 0
)

type cart struct {
	x         int
	y         int
	direction int
	lastTurn  int
}

type byPosition []cart

func (s byPosition) Len() int {
	return len(s)
}
func (s byPosition) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPosition) Less(i, j int) bool {
	return s[i].y < s[j].y || (s[i].y == s[j].y && s[i].x < s[j].x)
}

func main() {

	tracks := make(map[int]map[int]rune)
	carts := []cart{}

	file, _ := os.Open("input.txt")
	input := bufio.NewScanner(file)

	//Create a 2-dimensional map with the track layout and a slice with the cart locations and directions
	for lineCount := 0; input.Scan(); lineCount++ {
		tracks[lineCount] = make(map[int]rune)
		for columnCount, track := range input.Text() {

			if track == 60 || track == 62 || track == 94 || track == 118 {

				newCart := cart{x: columnCount, y: lineCount}

				switch track {
				case 60:
					newCart.direction = WEST
					track = 45
				case 62:
					newCart.direction = EAST
					track = 45
				case 94:
					newCart.direction = NORTH
					track = 124
				case 118:
					newCart.direction = SOUTH
					track = 124
				}

				carts = append(carts, newCart)

			}

			tracks[lineCount][columnCount] = track

		}
	}

	for {

		//Move each cart 1 position
		sort.Sort(byPosition(carts))
		for c := 0; c < len(carts); c++ {

			var track rune
			switch carts[c].direction {
			case NORTH:
				track = tracks[carts[c].y-1][carts[c].x]
				carts[c].y--
			case EAST:
				track = tracks[carts[c].y][carts[c].x+1]
				carts[c].x++
			case SOUTH:
				track = tracks[carts[c].y+1][carts[c].x]
				carts[c].y++
			case WEST:
				track = tracks[carts[c].y][carts[c].x-1]
				carts[c].x--
			}

			switch track {
			case 47: // Track with shape /
				switch carts[c].direction {
				case NORTH:
					carts[c].direction = EAST
				case EAST:
					carts[c].direction = NORTH
				case SOUTH:
					carts[c].direction = WEST
				case WEST:
					carts[c].direction = SOUTH
				}
			case 92: // Track with shape \
				switch carts[c].direction {
				case NORTH:
					carts[c].direction = WEST
				case EAST:
					carts[c].direction = SOUTH
				case SOUTH:
					carts[c].direction = EAST
				case WEST:
					carts[c].direction = NORTH
				}
			case 43: // Track with shape +
				switch carts[c].lastTurn {
				case LEFT:
					carts[c].lastTurn = STRAIGHT
				case STRAIGHT:
					switch carts[c].direction {
					case NORTH:
						carts[c].direction = EAST
					case EAST:
						carts[c].direction = SOUTH
					case SOUTH:
						carts[c].direction = WEST
					case WEST:
						carts[c].direction = NORTH
					}
					carts[c].lastTurn = RIGHT
				case RIGHT:
					switch carts[c].direction {
					case NORTH:
						carts[c].direction = WEST
					case EAST:
						carts[c].direction = NORTH
					case SOUTH:
						carts[c].direction = EAST
					case WEST:
						carts[c].direction = SOUTH
					}
					carts[c].lastTurn = LEFT
				}
			}

			//Check after every movement for possible crashes
			for sC := 0; sC < len(carts); sC++ {
				if carts[c] != carts[sC] && carts[c].x == carts[sC].x && carts[c].y == carts[sC].y {
					fmt.Printf("Crash detected at %v, %v\n", carts[c].x, carts[c].y)
					if sC > c {
						carts = append(carts[:sC], carts[sC+1:]...)
						carts = append(carts[:c], carts[c+1:]...)
						c = c - 1
					} else if c > sC {
						carts = append(carts[:c], carts[c+1:]...)
						carts = append(carts[:sC], carts[sC+1:]...)
						c = c - 2
					}
					break
				}
			}

		}

		//With only 1 cart left
		if len(carts) == 1 {
			fmt.Printf("last remaining cart came to a standstill at: %v, %v\n", carts[0].x, carts[0].y)
			return
		}

	}

}
