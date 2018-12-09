package main

import "fmt"

type marble struct {
	number int
	left *marble
	right *marble
}

const maxNumber = 7162600 //25 sample
const numPlayers = 438 //9 sample

func main() {

	players := make(map[int]int)

	for i := 1; i <= numPlayers; i++ {
		players[i] = 0
	}

	marble0 := marble {
		number: 0,
	}
	marble0.left = &marble0
	marble0.right = &marble0

	currentMarble := marble0
	nextNumber := 1

	for nextNumber <= maxNumber {

		for i := 1; i <= numPlayers && nextNumber <= maxNumber; i++ {

			if nextNumber%23 == 0 {

				players[i] += nextNumber + currentMarble.getNeighbour(-7).number
				currentMarble.getNeighbour(-8).right = currentMarble.getNeighbour(-6)
				currentMarble.getNeighbour(-6).left = currentMarble.getNeighbour(-8)
				currentMarble = *currentMarble.getNeighbour(-6)

			} else {

				newMarble := marble {
					number: nextNumber,
					left: currentMarble.getNeighbour(1),
					right: currentMarble.getNeighbour(2),
				}

				currentMarble.getNeighbour(2).left = &newMarble
				currentMarble.getNeighbour(1).right = &newMarble
				currentMarble = newMarble

			}

			nextNumber++

		}

	}

	maxScore := 0

	for _, v := range players {
		if v > maxScore {
			maxScore = v
		}
	}

	fmt.Println(maxScore)

}

func (current *marble) getNeighbour(offset int) *marble {
	var actual *marble
	if offset < 0 {
		actual = current.left.getNeighbour(offset + 1)
	} else if offset > 0 {
		actual = current.right.getNeighbour(offset - 1)
	} else if offset == 0 {
		actual = current
	}
	return actual
}