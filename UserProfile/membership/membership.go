package profile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dineshengg/matrimony/common/utility"
)

func MembershipRoutingFunctions(mux *http.ServeMux) {
	// Membership endpoints
	mux.HandleFunc("/create-membership", CreateMembershipHandler)
	mux.HandleFunc("/get-membership", GetMembershipHandler)
	mux.HandleFunc("/update-membership", UpdateMembershipHandler)
	mux.HandleFunc("/delete-membership", DeleteMembershipHandler)
}

// CreateMembershipHandler handles the creation of a user's membership
func CreateMembershipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		UserID         int    `json:"user_id"`
		MembershipType string `json:"membership_type"`
		DurationDays   int    `json:"duration_days"` // Duration in days
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, payload.DurationDays)

	err = CreateMembership(payload.UserID, payload.MembershipType, startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create membership: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Membership created successfully!")
}

// GetMembershipHandler handles retrieving a user's membership
func GetMembershipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	membership, err := GetMembership(utility.Atoi(userID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get membership: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(membership)
}

// UpdateMembershipHandler handles updating a user's membership
func UpdateMembershipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		UserID  int                    `json:"user_id"`
		Updates map[string]interface{} `json:"updates"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	err = UpdateMembership(payload.UserID, payload.Updates)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update membership: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Membership updated successfully!")
}

// DeleteMembershipHandler handles deleting a user's membership
func DeleteMembershipHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	err := DeleteMembership(utility.Atoi(userID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete membership: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Membership deleted successfully!")
}
