package utils

import (
	"encoding/json"
	"fmt"

	"bytes"

	log "github.com/sirupsen/logrus"
)

//Logic -
//Events are used to send the email immediately to users if not will be sent as per scheduled email job

const (
	//profile created related email events mapping
	WelcomeEmail                  = 1
	PasswordResetEmail            = 2
	VerificationEmail             = 3
	PreferenceMatchEmail          = 4
	ProfileUpdateEmail            = 5
	ReportSuspiciousActivityEmail = 6
	//password related email events mapping
	PasswordInvalidEmail = 7
	//profile related email events mapping
	ProfileVerifiedEmail = 8
	ProfileRejectedEmail = 9
	//interest and message related email events mapping
	InterestReceivedEmail = 10
	MessageReceivedEmail  = 11
	InterestAcceptedEmail = 12
	InterestRejectedEmail = 13
	//profile related email events mapping
	ProfileDeletedEmail     = 14
	ProfileDeactivatedEmail = 15
	ProfileReactivatedEmail = 16
	ProfileSuspendedEmail   = 17
	ProfileGotMarriedEmail  = 18
)

type PostData struct {
	Email string `json:"email"`
	Matid string `json:"matid"`
}

func init() {
	log.Debug("initializing email client utility for sending email events to server")
	setup()
}

func setup() {
	log.Debug("Setting up email client utility")
	//TODO - make sure email server is running and reachable if not panic and exit
}

func IsEmailClientInitialized() bool {
	return true
}

func LogEmailEvent(event int, data ...string) error {
	log.Debugf("Email Event: %d", event)
	switch event {
	case WelcomeEmail:
		if len(data) < 2 {
			return fmt.Errorf("insufficient data for welcome email event")
		}

		log.Info("Sending welcome email to user - ", data[0])
		url := "http://0.0.0.0:8000/jobs/welcomeemail/"
		payload := PostData{Email: data[0], Matid: data[1]}
		postBody, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Error marshalling JSON: %v", err)
			return fmt.Errorf("Error marshalling JSON: %v", err)
		}
		requestBody := bytes.NewBuffer(postBody)

		HttpPost(url, "application/json", requestBody)
		return nil
	case PasswordResetEmail:
		//http client to send password reset to email server

		log.Info("Sending password reset email to user", data[0])
		if len(data) < 2 {
			return fmt.Errorf("insufficient data for welcome email event")
		}

		url := "http://0.0.0.0:8000/jobs/passwordresetemail/"
		payload := PostData{Email: data[0], Matid: data[1]}
		postBody, err := json.Marshal(payload)
		if err != nil {
			log.Fatalf("Error marshalling JSON: %v", err)
			return fmt.Errorf("Error marshalling JSON: %v", err)
		}
		requestBody := bytes.NewBuffer(postBody)

		HttpPost(url, "application/json", requestBody)
		return nil
	case VerificationEmail:
		log.Info("Sending verification email to user")
	case PreferenceMatchEmail:
		log.Info("Sending preference match email to user")
	case ProfileUpdateEmail:
		log.Info("Sending profile update email to user")
	case ReportSuspiciousActivityEmail:
		log.Info("Sending report suspicious activity email to user")
	case PasswordInvalidEmail:
		log.Info("Sending password invalid email to user")
	case ProfileVerifiedEmail:
		log.Info("Sending profile verified email to user")
	case ProfileRejectedEmail:
		log.Info("Sending profile rejected email to user")
	case InterestReceivedEmail:
		log.Info("Sending interest received email to user")
	case MessageReceivedEmail:
		log.Info("Sending message received email to user")
	case InterestAcceptedEmail:
		log.Info("Sending interest accepted email to user")
	case InterestRejectedEmail:
		log.Info("Sending interest rejected email to user")
	case ProfileDeletedEmail:
		log.Info("Sending profile deleted email to user")
	case ProfileDeactivatedEmail:
		log.Info("Sending profile deactivated email to user")
	case ProfileReactivatedEmail:
		log.Info("Sending profile reactivated email to user")
	case ProfileSuspendedEmail:
		log.Info("Sending profile suspended email to user")
	case ProfileGotMarriedEmail:
		log.Info("Sending profile got married email to user")
	case 19: // Placeholder for any future email event
		log.Info("Sending future email event to user")
	default:
		log.Warnf("Unknown email event: %d", event)
	}
	return nil
}
