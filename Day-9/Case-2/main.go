package main

import "fmt"

type stone struct {
	number int
	left *stone
	right *stone
}

const maxNumber = 7162600 //25 sample
const numPlayers = 438 //9 sample

func main() {

	players := make(map[int]int)

	for i := 1; i <= numPlayers; i++ {
		players[i] = 0
	}

	stone0 := stone {
		number: 0,
	}
	stone0.left = &stone0
	stone0.right = &stone0

	currentMarble := stone0
	nextNumber := 1

	for nextNumber <= maxNumber {

		for i := 1; i <= numPlayers && nextNumber <= maxNumber; i++ {

			if nextNumber%23 == 0 {

				players[i] += nextNumber + currentMarble.left.left.left.left.left.left.left.number
				currentMarble.left.left.left.left.left.left.left.left.right = currentMarble.left.left.left.left.left.left
				currentMarble.left.left.left.left.left.left.left = currentMarble.left.left.left.left.left.left.left.left
				currentMarble = *currentMarble.left.left.left.left.left.left

			} else {

				newStone := stone {
					number: nextNumber,
					left: currentMarble.right,
					right: currentMarble.right.right,
				}

				currentMarble.right.right.left = &newStone
				currentMarble.right.right = &newStone
				currentMarble = newStone

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