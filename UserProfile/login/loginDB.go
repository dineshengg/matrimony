package login

import (
	"fmt"
	"time"

	"github.com/dineshengg/matrimony/common/utils"
)

func checkIfUserExists(email, phone string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1 OR phone = $2`
	var count int
	err := utils.GetDB().Exec(query, email, phone).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %v", err)
	}
	return count > 0, nil
}

func createUser(email, phone, hashedPassword string) error {
	query := `INSERT INTO users (email, phone, password, created_at) VALUES ($1, $2, $3, $4)`
	err := utils.GetDB().Exec(query, email, phone, hashedPassword, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}
