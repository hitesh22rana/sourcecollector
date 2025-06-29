package main

import (
	"fmt"
	"strconv"
	"strings"

	"advent-of-code-2024/pkg/utils"
)

func concatenate(a int, b int) int {
	multiplyFactor := 1
	c := b
	for c > 0 {
		multiplyFactor *= 10
		c /= 10
	}

	return a*multiplyFactor + b
}

func isPossibleWithThreeOperators(originalValue int, indx int, currentValue int, testValues []int) bool {
	if indx == len(testValues) {
		return originalValue == currentValue
	}

	return isPossibleWithThreeOperators(originalValue, indx+1, currentValue+testValues[indx], testValues) || isPossibleWithThreeOperators(originalValue, indx+1, currentValue*testValues[indx], testValues) || isPossibleWithThreeOperators(originalValue, indx+1, concatenate(currentValue, testValues[indx]), testValues)
}

func isPossibleWithTwoOperators(originalValue int, indx int, currentValue int, testValues []int) bool {
	if indx == len(testValues) {
		return originalValue == currentValue
	}

	return isPossibleWithTwoOperators(originalValue, indx+1, currentValue+testValues[indx], testValues) || isPossibleWithTwoOperators(originalValue, indx+1, currentValue*testValues[indx], testValues)
}

func PartTwo(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	ans := 0
	for _, line := range data {
		vals := strings.Split(line, ": ")

		originalValue, _ := strconv.Atoi(vals[0])
		testValuesStr := strings.Split(vals[1], " ")

		testValues := make([]int, len(testValuesStr))
		for i, v := range testValuesStr {
			testValues[i], _ = strconv.Atoi(v)
		}

		if isPossibleWithThreeOperators(originalValue, 0, 0, testValues) {
			ans += originalValue
		}
	}

	return ans
}

func PartOne(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	ans := 0
	for _, line := range data {
		vals := strings.Split(line, ": ")

		originalValue, _ := strconv.Atoi(vals[0])
		testValuesStr := strings.Split(vals[1], " ")

		testValues := make([]int, len(testValuesStr))
		for i, v := range testValuesStr {
			testValues[i], _ = strconv.Atoi(v)
		}

		if isPossibleWithTwoOperators(originalValue, 0, 0, testValues) {
			ans += originalValue
		}
	}

	return ans
}

func main() {
	var testCases []string

	// First part
	testCases = []string{
		"inputs/day-7/input1.txt",
		"inputs/day-7/input2.txt",
	}

	fmt.Println("Part One")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartOne(tc))
	}

	// Second part
	testCases = []string{
		"inputs/day-7/input1.txt",
		"inputs/day-7/input2.txt",
	}

	fmt.Println("\nPart Two")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartTwo(tc))
	}
}
