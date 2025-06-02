package message

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dineshengg/matrimony/common/utility"
)

func InterestRoutingFunctions(mux *http.ServeMux) {
	// Existing endpoints
	mux.HandleFunc("/send-message", SendMessageHandler)
	mux.HandleFunc("/get-messages", GetMessagesHandler)
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		SenderID   int    `json:"sender_id"`
		ReceiverID int    `json:"receiver_id"`
		Message    string `json:"message"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	err = SendMessage(payload.SenderID, payload.ReceiverID, payload.Message)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Message sent successfully!")
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	messages, err := GetMessages(utility.Atoi(userID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch messages: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
