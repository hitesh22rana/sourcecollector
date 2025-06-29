package main

import (
	"fmt"
	"strconv"

	"advent-of-code-2024/pkg/utils"
)

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func extractNumber(i *int, corruptedMemory *string) int {
	num := 0
	for num < 1000 && isNumber((*corruptedMemory)[*i]) {
		val, _ := strconv.Atoi(string((*corruptedMemory)[*i]))
		num = 10*num + val
		*i++
	}

	if 1 <= num && num < 1000 {
		return num
	}

	return -1
}

func PartTwo(inputFile string) int {
	data, err := utils.ReadInput(inputFile)
	if err != nil {
		panic(err)
	}

	ans := 0
	isEnabled := true
	for indx := 0; indx < len(data); indx++ {
		corruptedMemory := data[indx]
		for i := 0; i < len(corruptedMemory)-7; i++ {
			if corruptedMemory[i:i+4] == "do()" {
				isEnabled = true
			}

			if corruptedMemory[i:i+7] == "don't()" {
				isEnabled = false
			}

			if isEnabled && corruptedMemory[i:i+3] == "mul" {
				if corruptedMemory[i+3] != '(' {
					continue
				}

				i += 4
				num1 := extractNumber(&i, &corruptedMemory)
				if num1 == -1 {
					continue
				}

				if corruptedMemory[i] == ',' {
					i += 1
					num2 := extractNumber(&i, &corruptedMemory)
					if num2 == -1 {
						continue
					}

					if corruptedMemory[i] == ')' {
						ans += num1 * num2
					}
				}
			}
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
	for indx := 0; indx < len(data); indx++ {
		corruptedMemory := data[indx]
		for i := 0; i < len(corruptedMemory)-3; i++ {
			if corruptedMemory[i:i+3] == "mul" {
				if corruptedMemory[i+3] != '(' {
					continue
				}

				i += 4
				num1 := extractNumber(&i, &corruptedMemory)
				if num1 == -1 {
					continue
				}

				if corruptedMemory[i] == ',' {
					i += 1
					num2 := extractNumber(&i, &corruptedMemory)
					if num2 == -1 {
						continue
					}

					if corruptedMemory[i] == ')' {
						ans += num1 * num2
					}
				}
			}
		}
	}

	return ans
}

func main() {
	var testCases []string

	// First part
	testCases = []string{
		"inputs/day-3/input1.txt",
		"inputs/day-3/input2.txt",
	}

	fmt.Println("Part One")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartOne(tc))
	}

	// Second part
	testCases = []string{
		"inputs/day-3/input1.txt",
		"inputs/day-3/input2.txt",
	}

	fmt.Println("\nPart Two")
	for _, tc := range testCases {
		fmt.Printf("Test case %s: %d\n", tc, PartTwo(tc))
	}
}
