package login

import (
	"fmt"
	"time"

	"github.com/dineshengg/matrimony/common/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
)

type Enrolls struct {
	Id        int       `gorm:"column:id"`
	Email     string    `gorm:"column:email"`
	Phone     string    `gorm:"column:phone"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Looking   string    `gorm:"column:looking"`
}

type Profiles struct {
	Id               int            `gorm:"column:id"`
	Matrimonyid      string         `gorm:"column:matrimonyid"`
	FirstName        string         `gorm:"column:firstname"`
	SecondName       string         `gorm:"column:secondname"`
	Email            string         `gorm:"column:email"`
	Phone            string         `gorm:"column:phone"`
	Looking          string         `gorm:"column:looking"`
	DOB              datatypes.Date `gorm:"column:dob"`
	Age              int            `gorm:"column:age"`
	Gender           string         `gorm:"column:gender"`
	Country          string         `gorm:"column:country"`
	Religion         string         `gorm:"column:religion"`
	Caste            string         `gorm:"column:caste"`
	State            string         `gorm:"column:state"`
	City             string         `gorm:"column:city"`
	Language         string         `gorm:"column:language"`
	Password         string         `gorm:"column:password"`
	Hobbies          string         `gorm:"column:hobbies"`
	CreatedAt        time.Time      `gorm:"column:created_at"`
	Status           string         `gorm:"column:status"`
	SubscriptionType string         `gorm:"column:subscription_type"`
	Verified         bool           `gorm:"column:verified"`
	LastActive       time.Time      `gorm:"column:last_active"`
	PrefAgeMin       int            `gorm:"column:pref_age_min"`
	PrefAgeMax       int            `gorm:"column:pref_age_max"`
	PrefReligion     string         `gorm:"column:pref_religion"`
	PrefCaste        string         `gorm:"column:pref_caste"`
	PrefCountry      string         `gorm:"column:pref_country"`
	PrefState        string         `gorm:"column:pref_state"`
	PrefCity         string         `gorm:"column:pref_city"`
	PrefLanguage     string         `gorm:"column:pref_language"`
	UpdatedAt        time.Time      `gorm:"column:updated_at"`
}

func checkIfUserExists(email, phone string) (string, error) {
	//validate the input before quering
	if email == "" && phone == "" {
		return "", fmt.Errorf("email or phone is empty")
	}

	query := `SELECT matrimonyid FROM profiles WHERE email = $1 OR phone = $2`
	var matrimonyid string
	err := utils.GetDB().Raw(query, email, phone).Scan(&matrimonyid).Error
	if err != nil {
		return "", fmt.Errorf("failed to check user existence: %v", err)
	}
	log.Debugf("Matrimony for email and phone - %s", matrimonyid)
	return matrimonyid, nil
}

func createEnrolledUser(email, phone, looking string) (*Enrolls, error) {
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

func createNewProfile(profiles *Profiles) (*Profiles, error) {
	log.Debug("creating full user profile")
	err := utils.GetDB().Raw(`UPDATE profiles SET firstname = $1, secondname = $2, dob = $3, age = $4, gender = $5, country = $6, religion = $7, language = $8, password = $9, status = $10, subscription_type = $11, verified = $12, last_active = $13, pref_age_min = $14, pref_age_max = $15, pref_religion = $16, pref_country = $17, pref_language = $18, updated_at = $19 WHERE matrimonyid = $20 RETURNING id
	`,
		profiles.FirstName,
		profiles.SecondName,
		profiles.DOB,
		profiles.Age,
		profiles.Gender,
		profiles.Country,
		profiles.Religion,
		profiles.Language,
		profiles.Password,
		profiles.Status,
		profiles.SubscriptionType,
		profiles.Verified,
		time.Now(),
		profiles.PrefAgeMin,
		profiles.PrefAgeMax,
		profiles.PrefReligion,
		profiles.PrefCountry,
		profiles.PrefLanguage,
		time.Now(),
		profiles.Matrimonyid,
	).Scan(&profiles.Id).Error
	if err != nil {
		return profiles, fmt.Errorf("failed to create new user profile: %v", err)
	}
	return profiles, nil
}

func checkIfEmailExists(email string) (bool, error) {
	//validate the input before quering
	if email == "" {
		return false, fmt.Errorf("email is empty")
	}

	query := `SELECT COUNT(*) FROM profiles WHERE email = $1`
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

	query := `SELECT COUNT(*) FROM profiles WHERE phone = $1`
	var count int = 0
	err := utils.GetDB().Raw(query, phone).Scan(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check phone no existence: %v", err)
	}
	log.Debugf("User count for phone - %d", count)
	return count > 0, nil
}
