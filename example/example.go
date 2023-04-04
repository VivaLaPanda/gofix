package example

import (
	"errors"
	"fmt"
)

// This file is an example of a chunk of logic which contains a few functions and a main that calls those functions
// Some of those functions call other functions, some are leaf functions
// We will have unit tests for each functions. Some functions will fail, some will pass
// This will be used to test a separate Python script which uses GPT-4 to target and fix the failing functions

func add(a int, b int) int { return a + b }

func subtract(a int, b int) int { return a - b }

// multiply uses add to multiply two numbers by repeated addition
func multiply(a int, b int) int {
	var result int
	for i := 0; i < b; i++ {
		result = add(result, a)
	}
	return result
}

// divide uses subtract to divide two numbers by repeated subtraction
func divide(a int, b int) int {
	var result int = 0
	for a >= b {
		a = subtract(a, b)
		result++
	}
	return result
}

// multiplyString takes a number and a string and returns the string repeated that many times
func multiplyString(a int, b string) string {
	var result string
	for i := 0; i < a; i++ {
		result += b
	}
	return result
}

func ConvertRomanToDecimal(roman string) (int, error) {
	romanNumerals := map[string]int{"I": 1, "V": 5, "X": 10, "L": 50, "C": 100, "D": 500, "M": 1000}

	// check for invalid Roman numeral
	for i := 0; i < len(roman); i++ {
		if _, ok := romanNumerals[string(roman[i])]; !ok {
			return 0, errors.New("invalid Roman numeral")
		}
	}

	total := 0
	for i := 0; i < len(roman); i++ {
		currentValue := romanNumerals[string(roman[i])]
		if i+1 < len(roman) {
			nextValue := romanNumerals[string(roman[i+1])]
			if currentValue < nextValue {
				currentValue = -currentValue
			}
		}
		total += currentValue
	}
	return total, nil
}

func main() {
	fmt.Println(add(1, 2))
	fmt.Println(subtract(1, 2))
	fmt.Println(multiply(1, 2))
	fmt.Println(divide(1, 2))
	fmt.Println(multiplyString(1, "2"))
}
