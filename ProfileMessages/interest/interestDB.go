package interest

import (
	"fmt"
	"time"

	"github.com/dineshengg/matrimony/common/datasource"
)

func SendInterest(senderID, receiverID int, senderName string, senderAge int, receiverName string, receiverAge int) error {
	query := `
        INSERT INTO interests (sender_id, receiver_id, sender_name, sender_age, receiver_name, receiver_age, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, 'pending', $7)
    `
	DB := datasource.GetDB()
	_, err := DB.Exec(query, senderID, receiverID, senderName, senderAge, receiverName, receiverAge, time.Now())
	if err != nil {
		return fmt.Errorf("failed to send interest: %v", err)
	}
	return nil
}

// UpdateInterestStatus updates the status of an interest (accepted/rejected)
func UpdateInterestStatus(senderID, receiverID int, status string) error {
	query := `
        UPDATE interests
        SET status = $1
        WHERE sender_id = $2 AND receiver_id = $3
    `
	DB := datasource.GetDB()
	_, err := DB.Exec(query, status, senderID, receiverID)
	if err != nil {
		return fmt.Errorf("failed to update interest status: %v", err)
	}
	return nil
}

func GetInterests(userID int, interestType string) ([]map[string]interface{}, error) {
	var query string
	if interestType == "received" {
		query = `
            SELECT sender_id, sender_name, sender_age, status, created_at
            FROM interests
            WHERE receiver_id = $1
        `
	} else if interestType == "sent" {
		query = `
            SELECT receiver_id, receiver_name, receiver_age, status, created_at
            FROM interests
            WHERE sender_id = $1
        `
	} else {
		return nil, fmt.Errorf("invalid interest type: %s", interestType)
	}

	DB := datasource.GetDB()
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch interests: %v", err)
	}
	defer rows.Close()

	var interests []map[string]interface{}
	for rows.Next() {
		var otherUserID int
		var name string
		var age int
		var status string
		var createdAt time.Time

		err := rows.Scan(&otherUserID, &name, &age, &status, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		interest := map[string]interface{}{
			"user_id":    otherUserID,
			"name":       name,
			"age":        age,
			"status":     status,
			"created_at": createdAt,
		}
		interests = append(interests, interest)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return interests, nil
}
