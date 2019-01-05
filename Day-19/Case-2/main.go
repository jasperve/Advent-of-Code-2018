package main

import (
	"fmt"
)

func main() {

	reg := []int{0, 0, 0, 10550400, 10551306, 1}

	for reg[5] <= reg[4] {
		if reg[4]%reg[5] == 0 {
			reg[0] += reg[5]
		}
		reg[5]++
	}

	fmt.Println(reg)

}
