package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func DashboardRoutingFunctions(mux *http.ServeMux) {
	mux.HandleFunc("/get-users-by-time-range", GetUsersByTimeRangeHandler)
}

// GetUsersByTimeRange retrieves users created within a specific time range
func GetUsersByTimeRangeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the "range" query parameter
	timeRange := r.URL.Query().Get("range")
	if timeRange == "" {
		http.Error(w, "Missing 'range' query parameter", http.StatusBadRequest)
		return
	}

	// Calculate the start time based on the range
	var startTime time.Time
	switch timeRange {
	case "1d":
		startTime = time.Now().AddDate(0, 0, -1) // Last 1 day
	case "1w":
		startTime = time.Now().AddDate(0, 0, -7) // Last 1 week
	case "1m":
		startTime = time.Now().AddDate(0, -1, 0) // Last 1 month
	default:
		http.Error(w, "Invalid 'range' value. Use '1d', '1w', or '1m'.", http.StatusBadRequest)
		return
	}

	// Query the database
	users, err := GetUsersByTimeRange(startTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch users: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the users in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
