package message

import (
	"fmt"
	"time"

	"github.com/dineshengg/matrimony/common/datasource"
	_ "github.com/lib/pg"
)

func SendMessage(senderID, receiverID int, message string) error {
	// Validate if the sender is subscribed
	query := `
        SELECT COUNT(*) 
        FROM membership 
        WHERE user_id = $1 AND status = 'active'
    `
	DB := datasource.GetDB()
	var count int
	err := DB.QueryRow(query, senderID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to validate membership: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("sender is not subscribed")
	}

	// Insert the message into the database
	query = `
        INSERT INTO messages (sender_id, receiver_id, message, created_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err = DB.Exec(query, senderID, receiverID, message, time.Now())
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}

func GetMessages(userID int) ([]map[string]interface{}, error) {
	query := `
        SELECT sender_id, receiver_id, message, created_at
        FROM messages
        WHERE sender_id = $1 OR receiver_id = $1
        ORDER BY created_at DESC
    `
	DB := datasource.GetDB()
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %v", err)
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var senderID, receiverID int
		var message string
		var createdAt time.Time

		err := rows.Scan(&senderID, &receiverID, &message, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		msg := map[string]interface{}{
			"sender_id":   senderID,
			"receiver_id": receiverID,
			"message":     message,
			"created_at":  createdAt,
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return messages, nil
}
