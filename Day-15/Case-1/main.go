package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
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
	class int
	attackPower int
	hitPoints int
}

func main() {

	coordinates := make(map[int]map[int]*coordinate)
	elfs := []*coordinate{}
	goblins := []*coordinate{}
	
	file, err := os.Open("input-test.txt")
	if err != nil { log.Fatalln("Unable to open input file") }

	y := 0
	input := bufio.NewScanner(file)
	for input.Scan() {

		row := make(map[int]*coordinate)
		for x, u := range input.Text() {

			newCoordinate := coordinate{ x: x, y :y, class: empty }

			switch u {
			case 35: 
				newCoordinate.class = wall
			case 69:
				newCoordinate.class = elf
				newCoordinate.attackPower = 3
				newCoordinate.hitPoints = 200
				elfs = append(elfs, &newCoordinate)
			case 71:
				newCoordinate.class = goblin
				newCoordinate.attackPower = 3
				newCoordinate.hitPoints = 200
				goblins = append(goblins, &newCoordinate)
			}
			
			row[x] = &newCoordinate
			
		}
		coordinates[y] = row
		y++

	}

	// While there are still Goblins and Elfs
	for {

		for y := 0; y < len(coordinates); y++ {

			for x := 0; x < len(coordinates[y]); x++ {

				if coordinates[y][x].class == goblin {
					fmt.Println("goblin found")
				}
				if coordinates[y][x].class == elf {
					fmt.Println("elf found")
				}

			}

		}

		fmt.Println(len(elfs), len(goblins))
		
		break

	}

}