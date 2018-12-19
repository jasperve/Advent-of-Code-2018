package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
)

const (
	open = 0
	trees = 1
	lumberyard = 2

	numMinutes = 1000
)

func main() {

	grid := make(map[int]map[int]int)
	
	file, err := os.Open("input.txt")
	if err != nil { log.Fatalln("Unable to open input file") }

	y := 0
	input := bufio.NewScanner(file)
	for input.Scan() {
		row := make(map[int]int)
		for x, u := range input.Text() {
			switch u {
			case 35: 
				row[x] = lumberyard
			case 46:
				row[x] = open
			case 124:
				row[x] = trees
			}
		}
		grid[y] = row
		y++
	}

	oldResult := 0
	diffAtMinute500 := -1
	stopAtMinute := -1

	for m := 1; m <= numMinutes; m++ {
		
		newGrid := make(map[int]map[int]int)

		// Fill NEWGRID with the correct values
		for y := 0; y < len(grid); y++ {
			row := make(map[int]int)
			for x := 0; x < len(grid[y]); x++ {
				
				numLumberyard, numTrees := getSurroundings(y, x, grid)

				if grid[y][x] == open {
					if numTrees >= 3 { 
						row[x] = trees
					} else {
						row[x] = open
					}
				} else if grid[y][x] == trees {
					if numLumberyard >= 3 {
						row[x] = lumberyard
					} else {
						row[x] = trees
					}
				} else if grid[y][x] == lumberyard {
					if numLumberyard >= 1 && numTrees >= 1 {
						row[x] = lumberyard
					} else {
						row[x] = open
					}
				}

			}
			newGrid[y] = row
		}
		
		// Copy contents of NEW GRID to GRID
		for y := 0; y < len(newGrid); y++ {
			for x := 0; x < len(newGrid[y]); x++ {
				grid[y][x] = newGrid[y][x]
			}
		}

		totalNumLumberyard := 0
		totalNumTrees := 0
	
		for y := 0; y < len(grid); y++ {
			for x := 0; x < len(grid[y]); x++ {
				if grid[y][x] == lumberyard { totalNumLumberyard++ }
				if grid[y][x] == trees { totalNumTrees++ }
		/*		if grid[y][x] == lumberyard { fmt.Printf("#") }
				if grid[y][x] == trees { fmt.Printf("|") }
				if grid[y][x] == open { fmt.Printf(".") }*/
			}
			//fmt.Printf("\n")
		}
	
		result := totalNumLumberyard * totalNumTrees
		
		if m == 10 {
			fmt.Println("After 10 minutes:", result)
		}

		if oldResult - result == diffAtMinute500 { 
			stopAtMinute = m + (1000000000 - m)%(m-500)
		}
		
		if m == stopAtMinute {
			fmt.Println("After 1000000000 minutes:", result)
			break
		}

		if m == 500 { diffAtMinute500 = oldResult - result }
		oldResult = result

	}

}


func getSurroundings(y int, x int, grid map[int]map[int]int) (numLumberyard int, numTrees int) {
	
	for subY := -1; subY <= 1; subY++ {
		for subX := -1; subX <= 1; subX++ {
			if (subY == 0 && subX == 0) || y + subY < 0 || x + subX < 0 { continue }
			if grid[y+subY][x+subX] == lumberyard { numLumberyard++ }
			if grid[y+subY][x+subX] == trees { numTrees++ }
		}
	}
	return numLumberyard, numTrees

}
