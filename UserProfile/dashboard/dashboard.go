package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	TimeRange1Day   = "1d"
	TimeRange1Week  = "1w"
	TimeRange1Month = "1m"
)

func DashboardRoutingFunctions(mux *http.ServeMux) {
	mux.HandleFunc("/get-users-by-time-range", GetUsersByTimeRangeHandler)
}

// GetUsersByTimeRangeHandler retrieves users created within a specific time range
func GetUsersByTimeRangeHandler(w http.ResponseWriter, r *http.Request) {
	// Restrict to GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log the request URL and parameters
	log.Printf("Received request: %s, Query: %v", r.URL.Path, r.URL.Query())

	// Parse the "range" query parameter
	timeRange := r.URL.Query().Get("range")
	if timeRange == "" {
		http.Error(w, "Missing 'range' query parameter", http.StatusBadRequest)
		return
	}

	// Calculate the start time based on the range
	startTime, err := calculateStartTime(timeRange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query the database
	users, err := GetUsersByTimeRange(startTime)
	if err != nil {
		log.Printf("Error fetching users: %v", err) // Log the error
		http.Error(w, fmt.Sprintf("Failed to fetch users: %v", err), http.StatusInternalServerError)
		return
	}

	// Handle case where no users are found
	if len(users) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "No users found"}`))
		return
	}

	// Respond with the users in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Error encoding response: %v", err) // Log the error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// calculateStartTime calculates the start time based on the given range
func calculateStartTime(timeRange string) (time.Time, error) {
	switch timeRange {
	case TimeRange1Day:
		return time.Now().AddDate(0, 0, -1), nil // Last 1 day
	case TimeRange1Week:
		return time.Now().AddDate(0, 0, -7), nil // Last 1 week
	case TimeRange1Month:
		return time.Now().AddDate(0, -1, 0), nil // Last 1 month
	default:
		return time.Time{}, fmt.Errorf("Invalid 'range' value. Use '%s', '%s', or '%s'.", TimeRange1Day, TimeRange1Week, TimeRange1Month)
	}
}
