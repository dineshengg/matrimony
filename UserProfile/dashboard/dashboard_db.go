package dashboard

import (
	"fmt"
	"time"

	"github.com/dineshengg/matrimony/common/utils"
)

func GetUsersByTimeRange(startTime time.Time) ([]map[string]interface{}, error) {
	query := `
        SELECT id, name, email, gender, age, phone, created_at
        FROM profiles
        WHERE created_at >= $1
    `
	DB := utils.GetDB()
	rows, err := DB.Database.Query(query, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id, age, phone int
		var name, email, gender string
		var createdAt time.Time

		err := rows.Scan(&id, &name, &email, &gender, &age, &phone, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		user := map[string]interface{}{
			"id":         id,
			"name":       name,
			"email":      email,
			"gender":     gender,
			"age":        age,
			"phone":      phone,
			"created_at": createdAt,
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return users, nil
}
