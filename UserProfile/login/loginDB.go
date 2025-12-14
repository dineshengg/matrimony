package login

import (
	"fmt"
	"time"

	"github.com/beevik/guid"
	"github.com/dineshengg/matrimony/common/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type Enrolls struct {
	Id          int       `gorm:"column:id"`
	Matrimonyid string    `gorm:"column:matrimonyid"`
	Email       string    `gorm:"column:email"`
	Phone       string    `gorm:"column:phone"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	Looking     string    `gorm:"column:looking"`
}

type Profiles struct {
	Id          int            `gorm:"column:id"`
	Matrimonyid string         `gorm:"column:matrimonyid"`
	FirstName   string         `gorm:"column:firstname"`
	SecondName  string         `gorm:"column:secondname"`
	Email       string         `gorm:"column:email"`
	Phone       string         `gorm:"column:phone"`
	Looking     string         `gorm:"column:looking"`
	DOB         datatypes.Date `gorm:"column:dob"`
	Gender      string         `gorm:"column:gender"`
	Country     string         `gorm:"column:country"`
	Religion    string         `gorm:"column:religion"`
	Language    string         `gorm:"column:language"`
	Password    string         `gorm:"column:password"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
}

func resetPassword(email string) error {
	//validate the input before quering
	if email == "" {
		return fmt.Errorf("email is empty")
	}

	var query string
	query = `SELECT times FROM forgot WHERE email = $1`
	var count int = 0
	err := utils.GetDB().Raw(query, email).Scan(&count).Error
	if err != nil {
		return fmt.Errorf("failed to query forgot password existence: %v", err)
	}
	log.Debugf("Forgot password request count for email - %d", count)
	if count == 0 {
		//insert new record
		query = `INSERT INTO forgot (email, reset_at, times, guid) VALUES ($1, $2, $3, $4)`
		err := utils.GetDB().Exec(query, email, time.Now(), 1, guid.NewString()).Error
		if err != nil {
			return fmt.Errorf("failed to insert forgot password record: %v", err)
		}
	} else {
		//update existing record
		query = `UPDATE forgot SET reset_at = $1, guid = $2, times = times + 1 WHERE email = $3`
		err := utils.GetDB().Exec(query, time.Now(), guid.NewString(), email).Error
		if err != nil {
			return fmt.Errorf("failed to update forgot password record: %v", err)
		}
	}
	return nil
}

func checkIfUserExists(email, phone string) (bool, error) {
	//validate the input before quering
	if email == "" && phone == "" {
		return false, fmt.Errorf("email or phone is empty")
	}

	query := `SELECT COUNT(*) FROM profiles WHERE email = $1 OR phone = $2`
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

func createFullProfile(profiles *Profiles) (*Profiles, error) {
	log.Debug("creating full user profile")
	err := utils.GetDB().Create(&profiles).Error
	if err != nil {
		return profiles, fmt.Errorf("failed to create full user profile: %v", err)
	}
	return profiles, nil
}

func checkIfEmailExists(email string) (bool, error) {
	//validate the input before quering
	if email == "" {
		return false, fmt.Errorf("email is empty")
	}

	query := `SELECT COUNT(*) FROM enrolls WHERE email = $1`
	var count int = 0
	err := utils.GetDB().Raw(query, email).Scan(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %v", err)
	}
	log.Debugf("User count for email - %d", count)
	return count > 0, nil
}
func checkIfPhoneExists(phone string) (bool, error) {
	//validate the input before quering
	if phone == "" {
		return false, fmt.Errorf("email or phone is empty")
	}

	query := `SELECT COUNT(*) FROM enrolls WHERE phone = $1`
	var count int = 0
	err := utils.GetDB().Raw(query, phone).Scan(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check phone no existence: %v", err)
	}
	log.Debugf("User count for phone - %d", count)
	return count > 0, nil
}
