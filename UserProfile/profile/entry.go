package profile

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dineshengg/matrimony/common/utility"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type Profile struct {
	Id          int
	Matrimonyid string
	Firstname   string
	Second      string
	Email       string
	Phone       string
	Gender      string
	DOB         time.Time
	Looking     string
	Religion    string
	Country     string
	Language    string
	CreatedAt   int64
	Verified    int
}

func ProfileRoutingFunctions(router *routing.RouteGroup) {
	// Start HTTP server
	//profile table
	dashboardGrp := router.Group("/dashboard")
	dashboardGrp.Post("/create-profile", createProfileHandler)
	dashboardGrp.Get("/get-profile", getProfileHandler)
	dashboardGrp.Put("/update-profile", updateProfileHandler)
	dashboardGrp.Delete("/delete-profile", deleteProfileHandler)

	//preference table
	dashboardGrp.Post("/create-preference", CreatePreferenceHandler)
	dashboardGrp.Get("/get-preference", GetPreferenceHandler)
	dashboardGrp.Put("/update-preference", UpdatePreferenceHandler)
	dashboardGrp.Delete("/delete-preference", DeletePreferenceHandler)
}

func recovery() {
	//extend this function to handle all type of errors
	if e := recover(); e != nil {
		err, ok := e.(utility.StrErr)
		if ok {
			//error define localy in package
			fmt.Println("Recovered in f", err)
		} else {
			//orginal error
			fmt.Println("Recovered in f", e)
		}
	}
}

func createProfileHandler(ctx *routing.Context) error {
	defer recovery()

	if string(ctx.Method()) != fasthttp.MethodPost {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.Write([]byte("Invalid request method"))
		return nil
	}

	firstname := string(ctx.FormValue("firstname"))
	secondname := string(ctx.FormValue("secondname"))
	email := string(ctx.FormValue("email"))
	phone := string(ctx.FormValue("phone"))
	gender := string(ctx.FormValue("gender"))
	dob := string(ctx.FormValue("dob"))
	looking := string(ctx.FormValue("looking_for"))
	religion := string(ctx.FormValue("religion"))
	country := string(ctx.FormValue("country"))
	language := string(ctx.FormValue("language"))
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	verified := 0

	err := CreateProfile(firstname, secondname, email, utility.Atoi(phone), gender, dob, looking, religion, country, language, createdAt, verified)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(fmt.Sprintf("failed to create profile: %v", err)))
		return nil
	}

	// Redirect to login page on success
	ctx.Response.Header.Set("Location", "/login")
	ctx.SetStatusCode(fasthttp.StatusSeeOther)
	return nil
}

func getProfileHandler(ctx *routing.Context) error {
	defer recovery()

	// Parse JSON payload from body
	var payload map[string]interface{}
	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(fmt.Sprintf("Invalid JSON input: %v", err)))
		return nil
	}

	idFloat, ok := payload["id"].(float64)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Missing or invalid 'id' in request payload"))
		return nil
	}
	id := int(idFloat)

	firstname, secondname, email, phone, gender, dob, looking, religion, country, language, createdAt, verified, err := GetProfile(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Failed to get profile: %v", err)))
		return nil
	}

	response := map[string]interface{}{
		"id":          id,
		"firstname":   firstname,
		"secondname":  secondname,
		"email":       email,
		"phone":       phone,
		"gender":      gender,
		"dob":         dob,
		"looking_for": looking,
		"religion":    religion,
		"country":     country,
		"language":    language,
		"created_at":  createdAt,
		"verified":    verified,
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(response)
	return nil
}

func updateProfileHandler(ctx *routing.Context) error {
	defer recovery()

	if string(ctx.Method()) != fasthttp.MethodPut && string(ctx.Method()) != fasthttp.MethodPatch {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.Write([]byte("Invalid request method"))
		return nil
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte(fmt.Sprintf("Invalid JSON input: %v", err)))
		return nil
	}

	idFloat, ok := payload["id"].(float64)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Missing or invalid 'id' in request payload"))
		return nil
	}
	id := int(idFloat)

	firstname, _ := payload["firstname"].(string)
	secondname, _ := payload["secondname"].(string)
	email, _ := payload["email"].(string)
	phone, _ := payload["phone"].(string)
	gender, _ := payload["gender"].(string)
	dob, _ := payload["dob"].(string)
	looking, _ := payload["looking_for"].(string)
	religion, _ := payload["religion"].(string)
	country, _ := payload["country"].(string)
	language, _ := payload["language"].(string)
	verifiedFloat, _ := payload["verified"].(float64)
	verified := int(verifiedFloat)

	err := UpdateProfile(id, firstname, secondname, email, phone, gender, dob, looking, religion, country, language, verified)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Failed to update profile: %v", err)))
		return nil
	}
	ctx.SetContentType("text/html")
	ctx.Write([]byte(fmt.Sprintf("<h1> Profile updated successfully for ID %d </h1>", id)))
	return nil
}

func deleteProfileHandler(ctx *routing.Context) error {

	defer recovery()

	idStr := string(ctx.QueryArgs().Peek("id"))
	if idStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Missing id parameter"))
		return nil
	}
	id := utility.Atoi(idStr)

	err := DeleteProfile(id)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Failed to delete profile: %v", err)))
		return nil
	}

	ctx.SetContentType("text/plain")
	ctx.Write([]byte("Profile deleted successfully!"))
	return nil
}

///CRUD for preference table

// CreatePreferenceHandler handles the creation of a user's preference
func CreatePreferenceHandler(ctx *routing.Context) error {
	if string(ctx.Method()) != fasthttp.MethodPost {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.Write([]byte("Invalid request method"))
		return nil
	}

	var payload struct {
		UserID        int    `json:"user_id"`
		Gender        string `json:"gender"`
		Religion      string `json:"religion"`
		Caste         string `json:"caste"`
		Language      string `json:"language"`
		State         string `json:"state"`
		Country       string `json:"country"`
		WorkingStatus string `json:"working_status"`
		SalaryMin     int    `json:"salary_min"`
		SalaryMax     int    `json:"salary_max"`
		MaritalStatus string `json:"marital_status"`
	}

	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Invalid JSON input"))
		return nil
	}

	err := CreatePreference(payload.UserID, payload.Gender, payload.Religion, payload.Caste, payload.Language, payload.State, payload.Country, payload.WorkingStatus, payload.SalaryMin, payload.SalaryMax, payload.MaritalStatus)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Failed to create preference: %v", err)))
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write([]byte("Preference created successfully!"))
	return nil
}

// GetPreferenceHandler handles retrieving a user's preference
func GetPreferenceHandler(ctx *routing.Context) error {
	if string(ctx.Method()) != fasthttp.MethodGet {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.Write([]byte("Invalid request method"))
		return nil
	}

	userIDStr := string(ctx.QueryArgs().Peek("user_id"))
	if userIDStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Missing user_id parameter"))
		return nil
	}

	preference, err := GetPreference(utility.Atoi(userIDStr))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Failed to get preference: %v", err)))
		return nil
	}

	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(preference)
	return nil
}

// UpdatePreferenceHandler handles updating a user's preference
func UpdatePreferenceHandler(ctx *routing.Context) error {
	if string(ctx.Method()) != fasthttp.MethodPut && string(ctx.Method()) != fasthttp.MethodPatch {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.Write([]byte("Invalid request method"))
		return nil
	}

	var payload struct {
		UserID  int                    `json:"user_id"`
		Updates map[string]interface{} `json:"updates"`
	}

	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Invalid JSON input"))
		return nil
	}

	err := UpdatePreference(payload.UserID, payload.Updates)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Failed to update preference: %v", err)))
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write([]byte("Preference updated successfully!"))
	return nil
}

// DeletePreferenceHandler handles deleting a user's preference
func DeletePreferenceHandler(ctx *routing.Context) error {
	if string(ctx.Method()) != fasthttp.MethodDelete {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		ctx.Write([]byte("Invalid request method"))
		return nil
	}

	userIDStr := string(ctx.QueryArgs().Peek("user_id"))
	if userIDStr == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Missing user_id parameter"))
		return nil
	}

	err := DeletePreference(utility.Atoi(userIDStr))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(fmt.Sprintf("Failed to delete preference: %v", err)))
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Write([]byte("Preference deleted successfully!"))
	return nil
}
