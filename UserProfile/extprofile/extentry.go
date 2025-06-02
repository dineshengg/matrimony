package extprofile

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getUserPreferenceMatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the user preferences
	var preferences struct {
		AgeMin        int    `json:"age_min"`
		AgeMax        int    `json:"age_max"`
		Gender        string `json:"gender"`
		MaritalStatus string `json:"marital_status"`
		State         string `json:"state"`
		Location      string `json:"location"`
		Religion      string `json:"religion"`
		Caste         string `json:"caste"`
		Language      string `json:"language"`
		Color         string `json:"color"`
		WorkingStatus string `json:"working_status"`
		Company       string `json:"company"`
		SalaryMin     int    `json:"salary_min"`
		SalaryMax     int    `json:"salary_max"`
	}

	// Parse the JSON body from the request
	err := json.NewDecoder(r.Body).Decode(&preferences)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Call the GetUserPreferenceMatches function with the parsed preferences
	matches, err := GetUserPreferenceMatches(
		preferences.AgeMin,
		preferences.AgeMax,
		preferences.Gender,
		preferences.MaritalStatus,
		preferences.State,
		preferences.Location,
		preferences.Religion,
		preferences.Caste,
		preferences.Language,
		preferences.Color,
		preferences.WorkingStatus,
		preferences.Company,
		preferences.SalaryMin,
		preferences.SalaryMax,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch matches: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the matches in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}

func updateUserProfileHandler(w http.ResponseWriter, r *http.Request) {

	// {
	// 	"id": 1,
	// 	"updates": {
	// 	  "region": "North America",
	// 	  "caste": "Brahmin",
	// 	  "state": "California"
	// 	}
	// }

	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var payload struct {
		ID      int                    `json:"id"`
		Updates map[string]interface{} `json:"updates"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Call the UpdateUserProfile function
	err = UpdateUserProfile(payload.ID, payload.Updates)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user profile: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User profile updated successfully!")
}

func ExtProfileRoutingFunctions(mux *http.ServeMux) {
	mux.HandleFunc("/get-user-preference-matches", getUserPreferenceMatchesHandler)
	mux.HandleFunc("/update-user-profile", updateUserProfileHandler)
}
