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

				players[i] += nextNumber + getNeighbour(&currentMarble,-7).number
				getNeighbour(&currentMarble,-8).right = getNeighbour(&currentMarble,-6)
				getNeighbour(&currentMarble,-6).left = getNeighbour(&currentMarble,-8)
				currentMarble = *getNeighbour(&currentMarble,-6)

			} else {

				newMarble := marble {
					number: nextNumber,
					left: getNeighbour(&currentMarble, 1),
					right: getNeighbour(&currentMarble, 2),
				}

				getNeighbour(&currentMarble, 2).left = &newMarble
				getNeighbour(&currentMarble, 1).right = &newMarble
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


func getNeighbour(current *marble, offset int) *marble {
	var actual *marble
	if offset < 0 {
		actual = getNeighbour(current.left, offset+1)
	} else if offset > 0 {
		actual = getNeighbour(current.right, offset-1)
	} else if offset == 0 {
		actual = current
	}
	return actual
}