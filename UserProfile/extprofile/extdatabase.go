package extprofile

import (
	"fmt"
	"strings"

	"github.com/dineshengg/matrimony/common/utils"
)

// UpdateUserProfile updates specific attributes of a user profile
func UpdateUserProfile(id int, updates map[string]interface{}) error {
	// Build the SQL query dynamically based on the provided updates

	var setClauses []string
	var args []interface{}
	argIndex := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	// Ensure there are updates to apply
	if len(setClauses) == 0 {
		return fmt.Errorf("no updates provided")
	}

	// Add the ID as the last argument
	args = append(args, id)

	// Construct the final SQL query
	query := fmt.Sprintf("UPDATE profiles SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIndex)

	// Execute the query
	DB := utils.GetDB()
	_, err := DB.Database.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user profile: %v", err)
	}

	return nil
}

// GetUserPreferenceMatches retrieves profiles that match user preferences
func GetUserPreferenceMatches(ageMin, ageMax int, gender, maritalStatus, state, location, religion, caste, language, color, workingStatus, company string, salaryMin, salaryMax int) ([]map[string]interface{}, error) {
	query := `
        SELECT id, name, age, gender, marital_status, state, location, religion, caste, language, color, working_status, salary, company
        FROM profiles
        WHERE age BETWEEN $1 AND $2
        AND gender = $3
        AND marital_status = $4
        AND state = $5
        AND location = $6
        AND religion = $7
        AND caste = $8
        AND language = $9
        AND color = $10
        AND working_status = $11
        AND salary BETWEEN $12 AND $13
        AND company = $14
    `

	DB := utils.GetDB()
	rows, err := DB.Database.Query(query, ageMin, ageMax, gender, maritalStatus, state, location, religion, caste, language, color, workingStatus, salaryMin, salaryMax, company)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user preference matches: %v", err)
	}
	defer rows.Close()

	var matches []map[string]interface{}
	for rows.Next() {
		var id int
		var name, gender, maritalStatus, state, location, religion, caste, language, color, workingStatus, company string
		var age, salary int

		err := rows.Scan(&id, &name, &age, &gender, &maritalStatus, &state, &location, &religion, &caste, &language, &color, &workingStatus, &salary, &company)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		match := map[string]interface{}{
			"id":             id,
			"name":           name,
			"age":            age,
			"gender":         gender,
			"marital_status": maritalStatus,
			"state":          state,
			"location":       location,
			"religion":       religion,
			"caste":          caste,
			"language":       language,
			"color":          color,
			"working_status": workingStatus,
			"salary":         salary,
			"company":        company,
		}
		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return matches, nil
}
