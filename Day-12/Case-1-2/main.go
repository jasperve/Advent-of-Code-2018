package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type pot struct {
	id    int
	plant bool
	left  *pot
	right *pot
}

type note struct {
	id     []bool
	result bool
}

func main() {

	//Create a slice for the pots and make sure there are already 4 empty pots to the left for checking
	pots := []pot{pot{id: -4}, pot{id: -3}, pot{id: -2}, pot{id: -1}}

	notes := []note{}

	lineNumber := 1
	file, _ := os.Open("input.txt")
	input := bufio.NewScanner(file)
	for input.Scan() {
		if lineNumber == 1 {

			//Read the initial state into a slice of pots
			lineSplit := strings.Split(strings.TrimSpace(input.Text()), "initial state: ")
			for i, v := range lineSplit[1] {
				plant := false
				if string(v) == "#" {
					plant = true
				}
				newPot := pot{id: i, plant: plant}
				pots = append(pots, newPot)
			}

		} else if lineNumber > 2 {

			line := strings.Split(strings.TrimSpace(input.Text()), " => ")
			result := false
			id := []bool{false, false, false, false, false}
			if line[1] == "#" {
				result = true
			}
			for i, v := range line[0] {
				if string(v) == "#" {
					id[i] = true
				}
			}
			notes = append(notes, note{id: id, result: result})
		}
		lineNumber++
	}

	//Add 4 pots to the end of the slice
	pots = append(pots, pot{id: len(pots) - 4}, pot{id: len(pots) - 3}, pot{id: len(pots) - 2}, pot{id: len(pots) - 1})
	potsAdded := true

	total := 0
	diff := 0
	diffCounter := 0
	g := 0

	//For each generation
	for {

		fmt.Println(g)

		if potsAdded {
			//Assign pots in slice to each other
			for p := 0; p < len(pots); p++ {
				if p > 0 {
					pots[p].left = &pots[p-1]
				}
				if p < len(pots)-1 {
					pots[p].right = &pots[p+1]
				}
			}
			potsAdded = false
		}

		nextGPlant := []bool{}

		//For each pot
		for p := 0; p < len(pots); p++ {

			match := false

			//for each pattern which might possible match
		NOTES:
			for n := 0; n < len(notes); n++ {
				for i := 0; i < len(notes[n].id); i++ {
					//If the pot exists it should not match the boolean. If it doesn't it should only match if the boolean is false
					if (pots[p].getPot(i-2) != nil && pots[p].getPot(i-2).plant != notes[n].id[i]) || (pots[p].getPot(i-2) == nil && notes[n].id[i]) {
						continue NOTES
					}
				}
				match = notes[n].result
				break
			}

			nextGPlant = append(nextGPlant, match)

		}

		//For each pot set the new plant value
		for p := 0; p < len(pots); p++ {
			pots[p].plant = nextGPlant[p]
		}

		//Check for the first pot with a plant. If its within the first 4 pots add some extra empty pots
		for p := 0; p < len(pots) && p < 4; p++ {
			if pots[p].plant == true {
				for i := 1; i <= 4-p; i++ {
					pots = append([]pot{pot{id: pots[p].id - 1*p - 1*i}}, pots...)
					potsAdded = true
				}
			}
		}

		//Check for the last pot with a plant. If its within the last 4 pots add some extra empty pots
		for p := len(pots) - 1; p >= 0 && p > len(pots)-5; p-- {
			if pots[p].plant == true {
				numPots := len(pots)
				for i := 1; i <= 4-(numPots-1-p); i++ {
					pots = append(pots, pot{id: pots[p].id + 1*i + numPots - 1 - p})
					potsAdded = true
				}
			}
		}

		//Calculate the totals
		subTotal := 0
		for p := 0; p < len(pots); p++ {
			if pots[p].plant == true {
				subTotal += pots[p].id
			}
		}

		if subTotal-total == diff {
			diffCounter++
		} else {
			diff = subTotal - total
		}

		total = subTotal

		if diffCounter == 20 {
			break
		}

		g++

	}

	fmt.Printf("after 50 billion times: %v", (50000000000-g-1)*diff+total)

}

func (current *pot) getPot(offset int) *pot {
	var actual *pot
	if offset < 0 {
		if current.left != nil {
			actual = current.left.getPot(offset + 1)
		} else {
			actual = nil
		}
	} else if offset > 0 {
		if current.right != nil {
			actual = current.right.getPot(offset - 1)
		} else {
			actual = nil
		}
	} else if offset == 0 {
		actual = current
	}
	return actual
}
