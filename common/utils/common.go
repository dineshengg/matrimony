package utils

import (
	"time"

	"gorm.io/datatypes"
)

// CalculateAge calculates age based on the provided date of birth
func CalculateAge(dob datatypes.Date) int {
	now := time.Now()
	dobTime := time.Time(dob)

	age := now.Year() - dobTime.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.Month() < dobTime.Month() || (now.Month() == dobTime.Month() && now.Day() < dobTime.Day()) {
		age--
	}

	return age
}
