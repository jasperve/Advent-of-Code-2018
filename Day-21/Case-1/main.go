package main

import (
	"fmt"
)

func main() {

	register := []int{0,0,0,0,0,0}

	LOOP1:
	for {

		register[3] = register[5] | 65536
		register[5] = 9010242

		for {

			register[1] = register[3] & 255
			register[5] += register[1]
			register[5] = ((register[5] & 16777215) * 65899) & 16777215

			if 256 > register[3] {
				break LOOP1
			} else {

				register[1] = 0

				LOOP3:
				for {
				
					register[4] = register[1] + 1 
					register[4] *= 256

					if register[4] > register[3] {
						register[3] = register[1]
						break LOOP3
					} else {
						register[1]++
					}
							
				}

			}

		}

	}

	fmt.Println(register[5])

}