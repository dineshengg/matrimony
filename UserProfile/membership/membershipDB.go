package profile

import (
	"fmt"
	"time"

	"github.com/dineshengg/matrimony/common/utils"
)

// CreateMembership creates a new membership for a user
func CreateMembership(userID int, membershipType string, startDate, endDate time.Time) error {
	query := `
        INSERT INTO membership (user_id, membership_type, start_date, end_date, status)
        VALUES ($1, $2, $3, $4, 'active')
    `
	DB := utils.GetDB()
	_, err := DB.Database.Exec(query, userID, membershipType, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to create membership: %v", err)
	}
	return nil
}

// GetMembership retrieves a user's membership by user ID
func GetMembership(userID int) (map[string]interface{}, error) {
	query := `
        SELECT membership_type, start_date, end_date, status
        FROM membership
        WHERE user_id = $1
    `
	DB := utils.GetDB()
	row := DB.QueryRow(query, userID)

	var membershipType, status string
	var startDate, endDate time.Time

	err := row.Scan(&membershipType, &startDate, &endDate, &status)
	if err != nil {
		return nil, fmt.Errorf("failed to get membership: %v", err)
	}

	membership := map[string]interface{}{
		"membership_type": membershipType,
		"start_date":      startDate,
		"end_date":        endDate,
		"status":          status,
	}

	return membership, nil
}

// UpdateMembership updates a user's membership
func UpdateMembership(userID int, updates map[string]interface{}) error {
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
	query := fmt.Sprintf("UPDATE membership SET %s WHERE user_id = $%d", setClauses, argIndex)

	// Execute the query
	DB := utils.GetDB()
	_, err := DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update membership: %v", err)
	}

	return nil
}

// DeleteMembership deletes a user's membership by user ID
func DeleteMembership(userID int) error {
	query := `DELETE FROM membership WHERE user_id = $1`
	DB := utils.GetDB()
	_, err := DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete membership: %v", err)
	}
	return nil
}
