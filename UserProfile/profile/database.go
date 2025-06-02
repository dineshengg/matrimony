package profile

import (
	"fmt"

	"github.com/dineshengg/matrimony/common/utils"
)

// CreateProfile inserts a new profile into the utils
func CreateProfile(firstname, secondname, email string, phone int, gender, dob, looking, religion, country, language, createdAt string, verified int) error {
	query := `
        INSERT INTO profiles (firstname, secondname, email, phone, gender, dob, looking, religion, country, language, createdat, verified)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `
	DB := utils.GetDB()
	_, err := DB.Exec(query, firstname, secondname, email, phone, gender, dob, looking, religion, country, language, createdAt, verified)
	if err != nil {
		return fmt.Errorf("failed to create profile: %v", err)
	}
	return nil
}

// GetProfile retrieves a profile by ID
func GetProfile(id int) (string, string, string, int, string, string, string, string, string, string, string, int, error) {
	var firstname, secondname, email, gender, dob, looking, religion, country, language, createdat string
	var phone, verified int
	query := `SELECT firstname, secondname, email, phone, gender, dob, looking, religion, country, language, createdat, verified FROM profiles WHERE id = $1`
	DB := utils.GetDB()
	err := DB.Exec(query, id).Scan(&firstname, &secondname, &email, &phone, &gender, &dob, &looking, &religion, &country, &language, &createdat, &verified)
	if err != nil {
		return "", "", "", 0, "", "", "", "", "", "", "", 0, fmt.Errorf("failed to get profile: %v", err)
	}
	return firstname, secondname, email, phone, gender, dob, looking, religion, country, language, createdat, verified, nil
}

// UpdateProfile updates an existing profile
func UpdateProfile(id int, firstname, secondname, email string, phone int, gender, dob, looking, religion, country, language string, verified int) error {
	query := `
        UPDATE profiles 
        SET firstname = $1, secondname = $2, email = $3, phone = $4, gender = $5, dob = $6, looking = $7, religion = $8, country = $9, language = $10, verified = $11
        WHERE id = $12
    `
	DB := utils.GetDB()
	_, err := DB.Exec(query, firstname, secondname, email, phone, gender, dob, looking, religion, country, language, verified, id)
	if err != nil {
		return fmt.Errorf("failed to update profile: %v", err)
	}
	return nil
}

// DeleteProfile deletes a profile by ID
func DeleteProfile(id int) error {
	query := `DELETE FROM profiles WHERE id = $1`
	DB := utils.GetDB()
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %v", err)
	}
	return nil
}

//CRUD preference table

// CreatePreference inserts a new preference for a user
func CreatePreference(userID int, gender, religion, caste, language, state, country, workingStatus string, salaryMin, salaryMax int, maritalStatus string) error {
	query := `
        INSERT INTO preference (user_id, gender, religion, caste, language, state, country, working_status, salary_min, salary_max, marital_status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	DB := utils.GetDB()
	_, err := DB.Exec(query, userID, gender, religion, caste, language, state, country, workingStatus, salaryMin, salaryMax, maritalStatus)
	if err != nil {
		return fmt.Errorf("failed to create preference: %v", err)
	}
	return nil
}

// GetPreference retrieves a user's preference by user ID
func GetPreference(userID int) (map[string]interface{}, error) {
	query := `
        SELECT gender, religion, caste, language, state, country, working_status, salary_min, salary_max, marital_status
        FROM preference
        WHERE user_id = $1
    `
	DB := utils.GetDB()
	row := DB.Exec(query, userID)

	var gender, religion, caste, language, state, country, workingStatus, maritalStatus string
	var salaryMin, salaryMax int

	err := row.Scan(&gender, &religion, &caste, &language, &state, &country, &workingStatus, &salaryMin, &salaryMax, &maritalStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to get preference: %v", err)
	}

	preference := map[string]interface{}{
		"gender":         gender,
		"religion":       religion,
		"caste":          caste,
		"language":       language,
		"state":          state,
		"country":        country,
		"working_status": workingStatus,
		"salary_min":     salaryMin,
		"salary_max":     salaryMax,
		"marital_status": maritalStatus,
	}

	return preference, nil
}

// UpdatePreference updates a user's preference
func UpdatePreference(userID int, updates map[string]interface{}) error {
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

	// Add the user ID as the last argument
	args = append(args, userID)

	// Construct the final SQL query
	query := fmt.Sprintf("UPDATE preference SET %s WHERE user_id = $%d", setClauses, userID)

	// Execute the query
	DB := utils.GetDB()
	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update preference: %v", err)
	}

	return nil
}

// DeletePreference deletes a user's preference by user ID
func DeletePreference(userID int) error {
	query := `DELETE FROM preference WHERE user_id = $1`
	DB := utils.GetDB()
	_, err := DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete preference: %v", err)
	}
	return nil
}
