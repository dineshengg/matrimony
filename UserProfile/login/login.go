package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/dineshengg/matrimony/common/utils"
	"github.com/dineshengg/matrimony/middleware"

	log "github.com/sirupsen/logrus"
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
	newProfile.Post("/create-account", createAccountHandler, auth.CreateJWTToken)

	//no authentication flow where jwt token check is not required
	noauth := router.Group("/noauth")
	noauth.Post("/forgot-password", forgotPasswordHandler)
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
	//2. Password is in invalid => Show error message in the same page with attempts count
	//3. Any web server error while processing this request => Show a generic error page or same page with error message
	//4. Successful login => redirect to profile dashboard
	//5. If user is unsuccessfull after 5 attempts => show error page with message user password is invalid try resetting the password with form

	// Verify username and password from PostgreSQL
	emailid := ctx.FormValue("emailid")
	password := ctx.FormValue("password")
	if len(emailid) == 0 || len(password) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Email and password are required"))
		log.Debug("Email and password are required")
		return nil
	}

	// Check if the user exists in the database
	var firstname, secondname, email, hs_Password string
	db := utils.GetDB()
	err := db.Exec("SELECT firstname, secondname, email, password FROM users WHERE email = ?", emailid).Row().Scan(&firstname, &secondname, &email, &hs_Password)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Write([]byte("User not found or invalid credentials"))
		log.Debugf("User not found or invalid credentials - %v", err)
		return nil
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hs_Password), []byte(password)); err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Write([]byte("Invalid username or password"))
		log.Debugf("Invalid username or password - %v", err)
		return nil
	}

	// Redirect to profile dashboard
	ctx.Response.Header.Set("Location", "/profile/dashboard")
	ctx.SetStatusCode(fasthttp.StatusSeeOther)
	return nil

}

func dashboardHandler(ctx *routing.Context) error {
	//Validation of JWT is done at middleware authentication hence only show the profile dashboard
	//profileinfo := ctx.Get("profileinfo").(middleware.ProfileInfo)
	//Add the profile home page in the write io writer
	ctx.Write([]byte("Welcome, new user here is your upcoming dashboard"))
	return nil
}

func createAccountHandler(ctx *routing.Context) error {
	//1. Email id or phone number is already present => Show its already present in the same page
	//2. Once account is created successfully => redirect to dashboard page
	//3. Any other error happenes => show error message some problem happened and its not your fault come after some time

	if string(ctx.Method()) != http.MethodPost {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		//ctx.Write([]byte("Method not allowed, please use POST method"))
		log.Debug("Method not allowed, please use POST method")
		return fmt.Errorf("Method not allowed")
	}

	// Parse request body
	var user struct {
		Email   string `form:"email"`
		Phone   string `form:"phone"`
		Looking string `form:"looking"`
	}

	user.Email = string(ctx.FormValue("email"))
	user.Phone = string(ctx.FormValue("phone"))
	user.Looking = string(ctx.FormValue("looking"))

	log.Println("received", string(ctx.PostBody()))
	// err := json.Unmarshal(ctx.PostBody(), &user)
	// if err != nil {
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	// 	ctx.Write([]byte("Invalid request body"))
	// 	log.Debugf("Invalid request body - %v", err)
	// 	return nil
	// }

	if len(user.Email) <= 0 || len(user.Phone) <= 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//ctx.Write([]byte("User input is empty"))
		log.Debug("User input is empty - email, phone or password")
		return fmt.Errorf("User input is empty - email, phone or password")
	}

	// Validate if email or phone already exists
	exists, err := checkIfUserExists(user.Email, user.Phone)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//Land to an error page
		//ctx.Write([]byte("checkIfUserExists failed"))
		log.Debugf("checkIfUserExists failed with error - %v", err)
		return fmt.Errorf("checkIfUserExists failed with error - %v", err)
	}

	if exists {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//show in the same form that email and phone is already present
		//ctx.Write([]byte(fmt.Sprintf("User already exists with this email or phone number please use a different one")))
		//ctx.Redirect("/api/profile/dashboard", fasthttp.StatusFound)
		log.Debugf("User already exists with this email or phone number - %s, %s", user.Email, user.Phone)
		return fmt.Errorf("User already exists with email or phone number please login with your credentials")
	}

	// Hash the password
	// var hashedPassword []byte
	// if len(user.Password) > 8 {
	// 	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	// 		ctx.Write([]byte("Failed to hash password"))
	// 		log.Debugf("Failed to hash password - %v", err)
	// 		return nil
	// 	}
	// }

	// Insert user into the database
	enrolledUser, err := createUser(user.Email, user.Phone, user.Looking)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		//ctx.Write([]byte("Failed to create user account"))
		log.Debugf("Failed to create user account - %v", err)
		return fmt.Errorf("Failed to create user account - %v", err)
	}
	log.Debugf("user account created in db - %v", *enrolledUser)
	//auto create matrimony id is populated in the form to be used in cookie and jwt creation
	mf, err1 := ctx.MultipartForm()
	if err1 != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		//ctx.Write([]byte("Failed to process multipart form data"))
		log.Debugf("Failed to process multipart form data - %v", err1)
		return fmt.Errorf("Failed to process multipart form data - %v", err1)
	}
	enrolledUser.Matrimonyid = fmt.Sprintf("KAN%06d", enrolledUser.Id)
	mf.Value["matrimonyid"] = []string{enrolledUser.Matrimonyid}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write([]byte("User account created successfully!"))
	//telemetry reporting new user created for analyzing user growth for this site
	var msg string = fmt.Sprintf("New user account created with email: %s", user.Email)
	utils.LogTelemetry("UserAccountCreated", msg)
	// Sent even to email server for welcome email
	utils.LogEmailEvent(1, "Welcome to our esteemed matrimony services provided by Kandan Matrimony")
	log.Infof("User account created successfully for email: %s", user.Email)
	//Passing token type to auth middleware to create enroll token
	ctx.SetUserValue(middleware.Tokentype, middleware.EnrollTokenType)
	utils.Redirect(ctx, "/api/profile/dashboard")
	return nil
}

func forgotPasswordHandler(ctx *routing.Context) error {
	//1. Email id is present => send email to reset password
	//2. Email id is not present => show error message that email id doesnt exists and provide a link to create account
	//3. Any other error happens => show error message its not your fault please come after some time

	var req struct {
		Email string `json:"email"`
	}
	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Invalid request body"))
		log.Debugf("Invalid request body - %v", err)
		return nil
	}

	// Validate if email exists
	exists, err := checkIfUserExists(req.Email, "")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte("Failed to check user existence"))
		log.Debugf("Failed to check user existence - %v", err)
		return nil
	}
	if exists {
		// TODO: Send email to this email id
		ctx.Write([]byte(fmt.Sprintf("Email id exists, email was sent to this id")))
		ctx.Response.Header.Set("Location", "/services/sent-email")
	} else {
		ctx.SetStatusCode(fasthttp.StatusSeeOther)
		ctx.Write([]byte(fmt.Sprintf("Email id doesnt exists, please register here <a href=\"/new-profile/create-account\">Create Account</a>")))
	}
	return nil
}
