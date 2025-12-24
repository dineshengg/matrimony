package forgot

import (
	"fmt"
	"time"

	"github.com/beevik/guid"
	"github.com/dineshengg/matrimony/common/utils"
	log "github.com/sirupsen/logrus"
)

func checkIfUserExistsAlready(email, phone string) (string, error) {
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

func resetPassword(email, matid string) error {
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
		query = `INSERT INTO forgot (email, reset_at, times, guid, matrimonyid) VALUES ($1, $2, $3, $4, $5)`
		err := utils.GetDB().Exec(query, email, time.Now().Add(24*time.Hour), 1, guid.NewString(), matid).Error
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
