package main

import (
	"fmt"
	"regexp"
	"io/ioutil"
	"strconv"
	"strings"
)

const (

	// Constants for group system 
	immune = 0
	infection = 1

)

type group struct {
	system string
	units int
	hitpoints int
	attackDamage int
	attackType string
	initiative int
	weaknesses []string
	immunities []string
}

func main() {

	groups := []group{}

	input, _ := ioutil.ReadFile("input.txt")
	sectionsRegex := regexp.MustCompile("(?s)(Immune System|Infection):(?:(\\r\\n)|\\n)(.*?(?:(\\r\\n)|\\n){2})")
	sections := sectionsRegex.FindAllStringSubmatch(string(input), -1)

	for _, section := range sections {

		subsectionsRegex := regexp.MustCompile("(\\d*) units each with (\\d*) hit points (\\((.*)\\) |)with an attack that does (\\d*) (.*) damage at initiative (\\d*)(\r\n?|\n)?")
		subsections := subsectionsRegex.FindAllStringSubmatch(section[3], -1)

		for _, subsection := range subsections {

			units, _ := strconv.Atoi(subsection[2])
			hitpoints, _ := strconv.Atoi(subsection[3])
			attackDamage, _ := strconv.Atoi(subsection[6])
			initiative, _ := strconv.Atoi(subsection[8])

			newGroup := group {
				system: section[1],
				units: units,
				hitpoints: hitpoints,
				attackDamage: attackDamage,
				attackType: subsection[7],
				initiative: initiative,
			}

			if subsection[4] != "" {
				abilities := strings.Split(subsection[4], ";")
				for _, ability := range abilities {
					abilitySpecsRegex := regexp.MustCompile("(immune|weak) to (.*)")
					abilitySpecs := abilitySpecsRegex.FindStringSubmatch(ability)
					abilitySpecsSplit := strings.Split(strings.TrimSpace(abilitySpecs[2]), ",")
					if abilitySpecs[1] == "immune" {
						newGroup.immunities = append(newGroup.immunities, abilitySpecsSplit...)
					} else if abilitySpecs[1] == "weak" {
						newGroup.weaknesses = append(newGroup.weaknesses, abilitySpecsSplit...)
					}
				}
			}

			groups = append(groups, newGroup)

		}

	}
	
	fmt.Println(groups)

}