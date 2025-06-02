package interest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dineshengg/matrimony/common/utility"
)

func InterestRoutingFunctions(mux *http.ServeMux) {
	// Interest endpoints
	mux.HandleFunc("/send-interest", SendInterestHandler)
	mux.HandleFunc("/update-interest-status", UpdateInterestStatusHandler)
	mux.HandleFunc("/get-interests", GetInterestsHandler)
}

func SendInterestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		SenderID     int    `json:"sender_id"`
		ReceiverID   int    `json:"receiver_id"`
		SenderName   string `json:"sender_name"`
		SenderAge    int    `json:"sender_age"`
		ReceiverName string `json:"receiver_name"`
		ReceiverAge  int    `json:"receiver_age"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	err = SendInterest(payload.SenderID, payload.ReceiverID, payload.SenderName, payload.SenderAge, payload.ReceiverName, payload.ReceiverAge)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send interest: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Interest sent successfully!")
}

// UpdateInterestStatusHandler handles updating the status of an interest
func UpdateInterestStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		SenderID   int    `json:"sender_id"`
		ReceiverID int    `json:"receiver_id"`
		Status     string `json:"status"` // 'accepted' or 'rejected'
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	err = UpdateInterestStatus(payload.SenderID, payload.ReceiverID, payload.Status)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update interest status: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Interest status updated successfully!")
}

func GetInterestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	interestType := r.URL.Query().Get("type") // 'received' or 'sent'

	if userID == "" || interestType == "" {
		http.Error(w, "Missing user_id or type parameter", http.StatusBadRequest)
		return
	}

	interests, err := GetInterests(utility.Atoi(userID), interestType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch interests: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interests)
}
