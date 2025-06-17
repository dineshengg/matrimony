package login

import (
	"fmt"
	"time"

	"github.com/dineshengg/matrimony/common/utils"
	log "github.com/sirupsen/logrus"
)

type Enrolls struct {
	Id          int       `gorm:"column:id"`
	Matrimonyid string    `gorm:"column:matrimonyid"`
	Email       string    `gorm:"column:email"`
	Phone       string    `gorm:"column:phone"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	Looking     string    `gorm:"column:looking"`
}

func checkIfUserExists(email, phone string) (bool, error) {
	//validate the input before quering
	if email == "" && phone == "" {
		return false, fmt.Errorf("email or phone is empty")
	}

	query := `SELECT COUNT(*) FROM enrolls WHERE email = $1 OR phone = $2`
	var count int = 0
	err := utils.GetDB().Raw(query, email, phone).Scan(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %v", err)
	}
	log.Debugf("User count for email and phone - %d", count)
	return count > 0, nil
}

func createUser(email, phone, looking string) (*Enrolls, error) {
	//check if user, phone and hashed password are not empty
	if email == "" || phone == "" {
		return nil, fmt.Errorf("email, phone, and hashed password must not be empty")
	}
	//query := `INSERT INTO enroll (email, phone, looking) VALUES ($1, $2, $3)`
	enrolls := Enrolls{
		Email:   email,
		Phone:   phone,
		Looking: looking,
	}
	err := utils.GetDB().Create(&enrolls).Error
	if err != nil {
		return &enrolls, fmt.Errorf("failed to create user: %v", err)
	}
	return &enrolls, nil
}
