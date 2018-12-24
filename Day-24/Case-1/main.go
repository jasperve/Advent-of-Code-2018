package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type group struct {
	system       string
	units        int
	hitpoints    int
	attackDamage int
	attackType   string
	initiative   int
	weaknesses   map[string]bool
	immunities   map[string]bool
}

type byEffectivePowerAndInitiative []*group

func (c byEffectivePowerAndInitiative) Len() int {
	return len(c)
}
func (c byEffectivePowerAndInitiative) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byEffectivePowerAndInitiative) Less(i, j int) bool {
	return c[i].units*c[i].attackDamage > c[j].units*c[j].attackDamage || (c[i].units*c[i].attackDamage == c[j].units*c[j].attackDamage && c[i].initiative > c[j].initiative)
}

type action struct {
	attacker *group
	defender *group
	damage   int
}

type byInitiativeAttacker []action

func (c byInitiativeAttacker) Len() int {
	return len(c)
}
func (c byInitiativeAttacker) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c byInitiativeAttacker) Less(i, j int) bool {
	return c[i].attacker.initiative > c[j].attacker.initiative
}

func main() {

	groups := []*group{}

	input, _ := ioutil.ReadFile("input.txt")
	sectionsRegex := regexp.MustCompile("(?s)(Immune System|Infection):(?:(\\r\\n)|\\n)(.*?(?:(\\r\\n)|\\n){2})")
	sections := sectionsRegex.FindAllStringSubmatch(string(input), -1)

	for _, section := range sections {

		subsectionsRegex := regexp.MustCompile("(\\d*) units each with (\\d*) hit points (?:\\((.*)\\) |)with an attack that does (\\d*) (.*) damage at initiative (\\d*)(?:(\\r\\n)|\\n)?")
		subsections := subsectionsRegex.FindAllStringSubmatch(section[3], -1)

		for _, subsection := range subsections {

			units, _ := strconv.Atoi(subsection[1])
			hitpoints, _ := strconv.Atoi(subsection[2])
			attackDamage, _ := strconv.Atoi(subsection[4])
			initiative, _ := strconv.Atoi(subsection[6])

			newGroup := group{
				system:       section[1],
				units:        units,
				hitpoints:    hitpoints,
				attackDamage: attackDamage,
				attackType:   subsection[5],
				initiative:   initiative,
			}

			newGroup.immunities = make(map[string]bool)
			newGroup.weaknesses = make(map[string]bool)

			if subsection[3] != "" {
				abilities := strings.Split(subsection[3], ";")
				for _, ability := range abilities {
					abilitySpecsRegex := regexp.MustCompile("(immune|weak) to (.*)")
					abilitySpecs := abilitySpecsRegex.FindStringSubmatch(ability)
					abilitySpecsSplit := strings.Split(abilitySpecs[2], ",")
					for a := 0; a < len(abilitySpecsSplit); a++ {
						if abilitySpecs[1] == "immune" {
							newGroup.immunities[strings.TrimSpace(abilitySpecsSplit[a])] = true
						} else if abilitySpecs[1] == "weak" {
							newGroup.weaknesses[strings.TrimSpace(abilitySpecsSplit[a])] = true
						}
					}
				}
			}

			groups = append(groups, &newGroup)

		}

	}

	for {

		actions := []action{}

		sort.Sort(byEffectivePowerAndInitiative(groups))
		for a := 0; a < len(groups); a++ {

			pickedGroupDamage := -1
			var pickedGroup *group

			DEFENDERLOOP:
			for d := 0; d < len(groups); d++ {
				if groups[a].system == groups[d].system {
					continue
				}
				for a := 0; a < len(actions); a++ {
					if groups[d] == actions[a].defender {
						continue DEFENDERLOOP
					}
				}

				damage := groups[a].units * groups[a].attackDamage
				if _, ok := groups[d].immunities[groups[a].attackType]; ok {
					damage = 0
				}
				if _, ok := groups[d].weaknesses[groups[a].attackType]; ok {
					damage *= 2
				}

				if damage > pickedGroupDamage ||
					(damage == pickedGroupDamage && groups[d].units*groups[d].attackDamage > pickedGroup.units*pickedGroup.attackDamage) ||
					(damage == pickedGroupDamage && groups[d].units*groups[d].attackDamage == pickedGroup.units*pickedGroup.attackDamage && groups[d].initiative > pickedGroup.initiative) {

					pickedGroup = groups[d]
					pickedGroupDamage = damage
				}

			}

			if pickedGroupDamage != -1 {
				actions = append(actions, action{attacker: groups[a], defender: pickedGroup})
			}

		}

		sort.Sort(byInitiativeAttacker(actions))
		for a := 0; a < len(actions); a++ {

			if actions[a].attacker.units <= 0 { continue }

			damage := actions[a].attacker.units * actions[a].attacker.attackDamage
			if _, ok := actions[a].defender.immunities[actions[a].attacker.attackType]; ok {
				damage = 0
			}
			if _, ok := actions[a].defender.weaknesses[actions[a].attacker.attackType]; ok {
				damage *= 2
			}
			
			amountUnitsHealthy, amountUnitsWounded := ((actions[a].defender.units*actions[a].defender.hitpoints)-damage)/actions[a].defender.hitpoints,
				((actions[a].defender.units*actions[a].defender.hitpoints)-damage)%actions[a].defender.hitpoints

			if amountUnitsWounded > 0 {
				amountUnitsHealthy++
			}

			actions[a].defender.units = amountUnitsHealthy
			if actions[a].defender.units <= 0 {
				for g := 0; g < len(groups); g++ {
					if groups[g] == actions[a].defender {
						groups = append(groups[:g], groups[g+1:]...)
					}
				}
			}

		}

		// Check if there are still remaining units and already count remaining units
		immuneSystemFound := false
		infectionFound := false
		amountUnitsLeft := 0

		for g := 0; g < len(groups); g++ {
			if groups[g].system == "Immune System" {
				immuneSystemFound = true
			} else if groups[g].system == "Infection" {
				infectionFound = true
			}
			amountUnitsLeft += groups[g].units
		}

		if !immuneSystemFound || !infectionFound {
			fmt.Println("Amount units left:", amountUnitsLeft, "immune groups found:", immuneSystemFound, "infection groups found:", infectionFound)
			break
		}

	}

}
