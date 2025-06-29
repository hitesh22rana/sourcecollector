package main

import (
	"fmt"
	"strconv"
	"strings"

	"advent-of-code-2024/pkg/utils"
)

func isSafe(vals []int) bool {
	isInc := vals[0] < vals[1]
	for i := 0; i < len(vals)-1; i++ {
		var diff int
		if isInc {
			diff = vals[i+1] - vals[i]
		} else {
			diff = vals[i] - vals[i+1]
		}

		if diff < 1 || diff > 3 {
			return false
		}
	}

	return true
}

func PartOne(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	ans := 0
	for _, line := range data {
		strVals := strings.Split(line, " ")

		vals := make([]int, len(strVals))
		for i, val := range strVals {
			vals[i], _ = strconv.Atoi(val)
		}

		if isSafe(vals) {
			ans++
		}
	}

	return ans
}

func PartTwo(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	ans := 0
	for _, line := range data {
		strVals := strings.Split(line, " ")

		vals := make([]int, len(strVals))
		for i, val := range strVals {
			vals[i], _ = strconv.Atoi(val)
		}

		if isSafe(vals) {
			ans++
			continue
		}

		for i := 0; i < len(vals); i++ {
			newVals := make([]int, len(vals)-1)
			for j := 0; j < len(vals); j++ {
				if j < i {
					newVals[j] = vals[j]
				} else if j > i {
					newVals[j-1] = vals[j]
				}
			}

			if isSafe(newVals) {
				ans++
				break
			}
		}
	}

	return ans
}

func main() {
	var testCases []string

	// First part
	testCases = []string{
		"inputs/day-2/input1.txt",
		"inputs/day-2/input2.txt",
	}

	fmt.Println("Part One")
	for _, tc := range testCases {
		fmt.Println("For ", tc, ":", PartOne(tc))
	}

	// Second part
	testCases = []string{
		"inputs/day-2/input1.txt",
		"inputs/day-2/input2.txt",
	}

	fmt.Println("\nPart Two")
	for _, tc := range testCases {
		fmt.Println("For", tc, ":", PartTwo(tc))
	}
}
