package main

import (

	"sort"
	"fmt"
	"bufio"
	"time"
	"log"
	"os"
	"regexp"
	"strconv"

)

const timeFormat = "2006-01-02 15:04"
var timestampRegex *regexp.Regexp

type byTimestamp []string
func (s byTimestamp) Len() int {
	return len(s)
}
func (s byTimestamp) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byTimestamp) Less(i, j int) bool {
	iTime, err := time.Parse(timeFormat, timestampRegex.FindStringSubmatch(s[i])[1])
	if err != nil { log.Fatalln(err) }
	jTime, err := time.Parse(timeFormat, timestampRegex.FindStringSubmatch(s[j])[1])
	if err != nil { log.Fatalln(err) }
	return iTime.Before(jTime)
}

func main() {

	timestampRegex = regexp.MustCompile("\\[(.*)\\]")
	numberRegex := regexp.MustCompile("#(\\d*)")
	minutesRegex := regexp.MustCompile("\\d{2}:(\\d{2})")

	lines := []string{}

	file, err := os.Open("input.txt")
	if err != nil {	log.Fatalln(err) }

	input := bufio.NewScanner(file)

	for input.Scan() {
		lines = append(lines, input.Text())
	}

	sort.Sort(byTimestamp(lines))

	guards := map[string]map[int]int{}
	var activeGuard string

	for i := 0; i < len(lines); i++ {

		number := numberRegex.FindStringSubmatch(lines[i])
		if len(number) > 0 {

			if guards[number[1]] == nil { guards[number[1]] = make(map[int]int)	}
			activeGuard = number[1]

		} else {

			from, _ := strconv.Atoi(minutesRegex.FindStringSubmatch(lines[i])[1])
			till, _ := strconv.Atoi(minutesRegex.FindStringSubmatch(lines[i+1])[1])

			for from < till {
				guards[activeGuard][from] = guards[activeGuard][from] + 1
				from++
			}

			i++

		}

	}

	var pickedGuard string
	var pickedMinute int
	var timesSleepiestMinute int


	for i, guard := range guards {

		for minute, times := range guard {
			if times > timesSleepiestMinute {
				pickedGuard = i
				pickedMinute = minute
				timesSleepiestMinute = times
			}
		}

	}

	intPickedGuard, _ := strconv.Atoi(pickedGuard)

	fmt.Printf("The guard picked was number %v and minute was %v. So the result is %v!", pickedGuard, pickedMinute, intPickedGuard*pickedMinute)

}