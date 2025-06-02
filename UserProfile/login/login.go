package login

import (
	"encoding/json"
	"fmt"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/dineshengg/matrimony/common/utils"
	"github.com/dineshengg/matrimony/middleware"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ResponseCode struct {
	err error
}
type Timeouthandler struct {
	handler  func(ctx *routing.Context) error
	login_ch chan ResponseCode
}

func LoginRoutingFunctions(router *routing.RouteGroup) {
	//all below end points needs middleware both authentication, centralized logging and request logging
	profileGroup := router.Group("/profile")
	_ = middleware.NewMiddleWare(profileGroup, true, true, true)
	auth := middleware.NewAuthentication(profileGroup)
	profileGroup.Post("/login", Timeouthandler{handler: loginHandler}.TimeOutHandler, auth.CreateJWTToken)
	profileGroup.Get("/dashboard", dashboardHandler)
	//Doesnt need authentication middleware and hence creating a separate group end points
	newProfile := router.Group("/new-profile")
	_ = middleware.NewMiddleWare(newProfile, true, true, false)
	newProfile.Post("/create-account", createAccountHandler)
	newProfile.Post("/forgot-password", forgotPasswordHandler)
}

func (h Timeouthandler) TimeOutHandler(ctx *routing.Context) error {

	h.login_ch = make(chan ResponseCode, 1)
	go func() {
		err1 := h.handler(ctx)
		h.login_ch <- ResponseCode{err: err1}
	}()

	select {
	case h.login_ch <- ResponseCode{}:
		// call completed successfully
		return nil

	case <-time.After(5 * time.Second):
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		ctx.Write([]byte("Request timed out"))
		return fmt.Errorf("request timed out")
	case <-ctx.Done():
		//server is shutting down or request is cancelled
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.Write([]byte("Request aborted"))
		return fmt.Errorf("request aborted")
	}
}

func loginHandler(ctx *routing.Context) error {

	//1. Email id or phone number is not present => redirect to login page with GET method
	//2. Password is in invalid => Show error message in the same page
	//3. Any web server error while processing this request => Show a generic error page or same page with error message
	//4. Successful login => redirect to profile dashboard

	// Verify username and password from PostgreSQL
	emailid := ctx.FormValue("emailid")
	password := ctx.FormValue("password")
	if len(emailid) == 0 || len(password) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Email and password are required"))
		return fmt.Errorf("email and password are required")
	}

	// Check if the user exists in the database
	var firstname, secondname, email, hs_Password string
	db := utils.GetDB()
	err := db.Exec("SELECT firstname, secondname, email, password FROM users WHERE email = ?", emailid).Row().Scan(&firstname, &secondname, &email, &hs_Password)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Write([]byte("Invalid username or password"))
		return fmt.Errorf("failed to get user name and password from database: %v", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hs_Password), []byte(password)); err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Write([]byte("Invalid username or password"))
		return fmt.Errorf("Invalid password: %v", err)
	}

	// JWT claims
	claims := middleware.ProfileInfo{
		FirstName:  firstname,
		SecondName: secondname,
		Email:      email,
		Role:       "user",
		Exp:        time.Now().Add(24 * time.Hour).Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "MatrimonyApp",
			Subject:   email,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()), // Unique ID for the token
			Audience:  []string{"MatrimonyAppUsers"},
		},
	}

	ctx.Set("profileinfo", claims)

	// Redirect to profile dashboard
	ctx.Response.Header.Set("Location", "/profile/dashboard")
	ctx.SetStatusCode(fasthttp.StatusSeeOther)
	return nil

}

func dashboardHandler(ctx *routing.Context) error {
	//Validation of JWT is done at middleware authentication hence only show the profile dashboard
	profileinfo := ctx.Get("profileinfo").(middleware.ProfileInfo)

	//Add the profile home page in the write io writer
	ctx.Write([]byte(fmt.Sprintf("Welcome, %s!", profileinfo.FirstName+" "+profileinfo.SecondName)))
	return nil
}

func createAccountHandler(ctx *routing.Context) error {
	// Parse request body
	var user struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(ctx.PostBody(), &user)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Invalid request body"))
		return fmt.Errorf("failed to parse request body: %v", err)
	}

	// Validate if email or phone already exists
	exists, err := checkIfUserExists(user.Email, user.Phone)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//Land to an error page
		ctx.Write([]byte("Invalid request body"))
		return nil
	}
	if exists {
		ctx.SetStatusCode(fasthttp.StatusProxyAuthRequired)
		//show in the same form that email and phone is already present
		ctx.Write([]byte("Invalid request body"))
		return nil
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte("Failed to hash password"))
		return nil
	}

	// Insert user into the database
	err = createUser(user.Email, user.Phone, string(hashedPassword))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte("Failed to create user account"))
		return nil
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write([]byte("User account created successfully!"))
	return nil
}

func forgotPasswordHandler(ctx *routing.Context) error {
	var req struct {
		Email string `json:"email"`
	}
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Invalid request body"))
		return nil
	}

	// Validate if email exists
	exists, err := checkIfUserExists(req.Email, "")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte("Failed to check user existence"))
		return nil
	}
	if exists {
		// TODO: Send email to this email id
		ctx.Write([]byte(fmt.Sprintf("Email id exists, email was sent to this id")))
	} else {
		ctx.SetStatusCode(fasthttp.StatusSeeOther)
		ctx.Response.Header.Set("Location", "/services/sent-email")
	}
	return nil
}
