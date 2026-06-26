package utils

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// IsValidISO6346 validates a container number according to ISO 6346 standard
func IsValidISO6346(containerNumber string) bool {
	containerNumber = strings.ToUpper(strings.ReplaceAll(containerNumber, " ", ""))

	// 1. Basic format check: 4 letters followed by 7 digits
	re := regexp.MustCompile(`^[A-Z]{3}[UJZ][0-9]{7}$`)
	if !re.MatchString(containerNumber) {
		return false
	}

	// 2. Check digit calculation
	sum := 0
	for i := 0; i < 10; i++ {
		char := containerNumber[i]
		var val int

		if char >= 'A' && char <= 'Z' {
			// Letter values: A=10, B=12, ..., Z=38 (skipping multiples of 11)
			val = int(char) - 55 // A is 65 in ASCII, 65-55 = 10
			
			// Adjust for skipped multiples of 11 (11, 22, 33)
			if val >= 11 {
				val++
			}
			if val >= 22 {
				val++
			}
			if val >= 33 {
				val++
			}
		} else {
			// Digit values: 0-9
			val, _ = strconv.Atoi(string(char))
		}

		// Multiply by 2^i
		sum += val * int(math.Pow(2, float64(i)))
	}

	// The check digit is the remainder of (sum / 11)
	// If the remainder is 10, the check digit is 0
	checkDigit := sum % 11
	if checkDigit == 10 {
		checkDigit = 0
	}

	lastDigit, _ := strconv.Atoi(string(containerNumber[10]))
	return checkDigit == lastDigit
}
