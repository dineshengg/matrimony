package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"gorm.io/datatypes"

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
	//TODO - Need to add rate limiting to prevent DDOS attacks
	profileGroup := router.Group("/profile")
	//_ = middleware.NewMiddleWare(profileGroup, true, true, true)
	auth := middleware.NewAuthentication(profileGroup)
	profileGroup.Get("/dashboard", auth.Authenticate, dashboardHandler)
	profileGroup.Get("/newprofiledetails", auth.Authenticate, newProfileDetails)

	//Doesnt need authentication middleware and hence creating a separate group end points
	newProfile := router.Group("/new-profile")
	//_ = middleware.NewMiddleWare(newProfile, true, true, false)
	//auth1 := middleware.NewAuthentication(newProfile)
	newProfile.Post("/login", loginHandler, auth.CreateJWTToken)
	//newProfile.Post("/login", Timeouthandler{handler: loginHandler}.TimeOutHandler, auth.CreateJWTToken)
	// handler to add new profile from home page with email and phone number
	newProfile.Post("/create-account", createAccountHandler, auth.CreateJWTToken)
	// handler to add new profile from nav bar with full profile details
	newProfile.Post("/create-new-account", createNewProfileAccountHandler, auth.CreateJWTToken)
	newProfile.Post("/create-full-account", createFullProfileAccountHandler, auth.CreateJWTToken)
	newProfile.Post("/validate", checkIfUserExistsHandler)
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
	emailid := ctx.FormValue("email")
	password := ctx.FormValue("password")
	if len(emailid) == 0 || len(password) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("Email and password are required"))
		log.Debug("Email and password are required")
		return fmt.Errorf("Email and password are required")
	}

	// Check if the user exists in the database
	var firstname, secondname, email, phone, hs_Password, matrimonyid, looking string
	db := utils.GetDB()
	err := db.Raw("SELECT matrimonyid, firstname, secondname, email, phone, password, looking FROM profiles WHERE email = ?", string(emailid)).
		Row().Scan(&matrimonyid, &firstname, &secondname, &email, &phone, &hs_Password, &looking)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Write([]byte("User not found or invalid credentials"))
		log.Debugf("User not found or invalid credentials - %v", err)
		return fmt.Errorf("User not found or invalid credentials - %v", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hs_Password), []byte(password)); err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Write([]byte("Invalid username or password"))
		log.Debugf("Invalid username or password - %v", err)
		return fmt.Errorf("Invalid username or password - %v", err)
	}

	// Update user status to active and set last_active timestamp
	if err := updateUserActivityStatus(matrimonyid, "active"); err != nil {
		log.Warnf("Failed to update user activity status: %v", err)
		// Don't fail login if status update fails
	}

	// Get or create multipart form to add values
	mf, err1 := ctx.MultipartForm()
	if err1 != nil {
		// If multipart form doesn't exist, create form values manually
		ctx.Request.PostArgs().Set("matrimonyid", matrimonyid)
		ctx.Request.PostArgs().Set("firstname", firstname)
		ctx.Request.PostArgs().Set("secondname", secondname)
		ctx.Request.PostArgs().Set("phone", phone)
		ctx.Request.PostArgs().Set("looking", looking)
		ctx.Request.PostArgs().Set("email", email)
	} else {
		// Add values to existing multipart form
		mf.Value["matrimonyid"] = []string{matrimonyid}
		mf.Value["firstname"] = []string{firstname}
		mf.Value["secondname"] = []string{secondname}
		mf.Value["phone"] = []string{phone}
		mf.Value["looking"] = []string{looking}
		mf.Value["email"] = []string{email}
	}

	// Redirect to profile dashboard
	ctx.SetUserValue(middleware.Tokentype, middleware.ProfileTokenType)
	fmt.Println("token type", ctx.UserValue(middleware.Tokentype))
	utils.Redirect(ctx, "/api/profile/dashboard")
	ctx.SetStatusCode(fasthttp.StatusSeeOther)
	return nil

}

func dashboardHandler(ctx *routing.Context) error {
	//Validation of JWT is done at middleware authentication hence only show the profile dashboard
	//profileinfo := ctx.Get("profileinfo").(middleware.ProfileInfo)
	//Add the profile home page in the write io writer
	log.Debug("Welcome to user dashboard")
	user := ctx.Get(middleware.Token).(middleware.ProfileInfo)
	ctx.Write([]byte(fmt.Sprintf("Welcome %s %s to your profile dashboard", user.FirstName, user.SecondName)))
	ctx.SetStatusCode(fasthttp.StatusOK)
	return nil
}

// This a form filler to send back the user with enrolled token user to fill further details.
func newProfileDetails(ctx *routing.Context) error {
	//1. If user is not logged in => redirect to login page
	//2. If user is logged in => show the profile details page with form to fill the details
	//3. Any other error happens => show error message its not your fault please come after some time

	// Check if the user is authenticated
	var profileInfo middleware.ProfileInfo
	claims := ctx.Get(middleware.Token)
	if claims == nil {
		log.Debug("User is not having any valid enroll token hence will show email id and phone number fields")
	} else {
		profileInfo = claims.(middleware.ProfileInfo)
	}

	// fill the form with user details and send it to user
	if page, err := utils.RenderTemplatePage(ctx, "resources/login/newuser.html", profileInfo); err == nil {
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.Write(page)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return nil
	} else {
		ctx.SetStatusCode(fasthttp.StatusSeeOther)
		utils.Redirect(ctx, "/static/login/blankuser.html")
	}
	return nil
}

// This is home page new user handler to create an account in database and to check if user already exists.
func createAccountHandler(ctx *routing.Context) error {
	//1. Email id or phone number is already present => Show its already present in the same page
	//2. Once account is created successfully => redirect to new profile details page
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
	matid, err := checkIfUserExists(user.Email, user.Phone)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//Land to an error page
		//ctx.Write([]byte("checkIfUserExists failed"))
		log.Debugf("checkIfUserExists failed with error - %v", err)
		return fmt.Errorf("checkIfUserExists failed with error - %v", err)
	}

	if matid != "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//show in the same form that email and phone is already present
		//ctx.Write([]byte(fmt.Sprintf("User already exists with this email or phone number please use a different one")))
		//ctx.Redirect("/api/profile/dashboard", fasthttp.StatusFound)
		log.Debugf("User already exists with this email or phone number - %s, %s", user.Email, user.Phone)
		return fmt.Errorf("User already exists with email or phone number please login with your credentials")
	}

	// Insert user into the database
	enrolledUser, err := createEnrolledUser(user.Email, user.Phone, user.Looking)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		//ctx.Write([]byte("Failed to create user account"))
		log.Debugf("Failed to create user account - %v", err)
		return fmt.Errorf("Failed to create user account - %v", err)
	}
	log.Debugf("user account created in db - %v", *enrolledUser)

	var profiles *Profiles
	profiles, err = createFullProfile(&Profiles{Email: user.Email, Phone: user.Phone, Looking: user.Looking, SubscriptionType: "enrolled_only"})
	//auto create matrimony id is populated in the form to be used in cookie and jwt creation
	mf, err1 := ctx.MultipartForm()
	if err1 != nil {
		// If multipart form doesn't exist, create form values manually
		ctx.Request.PostArgs().Set("matrimonyid", fmt.Sprintf("KAN%020d", profiles.Id))
		ctx.Request.PostArgs().Set("firstname", "newuser")
		ctx.Request.PostArgs().Set("secondname", "newuser")
	} else {
		// Add values to existing multipart form
		mf.Value["matrimonyid"] = []string{fmt.Sprintf("KAN%020d", profiles.Id)}
		mf.Value["firstname"] = []string{"newuser"}
		mf.Value["secondname"] = []string{"newuser"}
	}
	if err1 != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		//ctx.Write([]byte("Failed to process multipart form data"))
		log.Debugf("Failed to process multipart form data - %v", err1)
		return fmt.Errorf("Failed to process multipart form data - %v", err1)
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.Write([]byte("User account created successfully!"))
	//telemetry reporting new user created for analyzing user growth for this site
	var msg string = fmt.Sprintf("New user account created with email: %s", user.Email)
	utils.LogTelemetry("UserAccountCreated", msg)
	// Sent even to email server for welcome email
	utils.LogEmailEvent(1)
	log.Infof("User account created successfully for email: %s %s", user.Email, fmt.Sprintf("KAN%020d", profiles.Id))
	//Passing token type to auth middleware to create enroll token
	ctx.SetUserValue(middleware.Tokentype, middleware.ProfileTokenType)
	fmt.Println("token type", ctx.UserValue(middleware.Tokentype))
	utils.Redirect(ctx, "/api/profile/newprofiledetails")
	return nil
}

// This is home page new user handler to create an account in database and to check if user already exists.
func createNewProfileAccountHandler(ctx *routing.Context) error {
	//1. Email id or phone number is already present => Show its already present in the same page
	//2. Once account is created successfully => redirect to new profile details page
	//3. Any other error happenes => show error message some problem happened and its not your fault come after some time

	if string(ctx.Method()) != http.MethodPost {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		//ctx.Write([]byte("Method not allowed, please use POST method"))
		log.Debug("Method not allowed, please use POST method")
		return fmt.Errorf("Method not allowed")
	}

	// Parse request body
	var err error
	var profile Profiles
	profile.Matrimonyid = string(ctx.FormValue("matrimonyid"))
	profile.FirstName = string(ctx.FormValue("firstname"))
	profile.SecondName = string(ctx.FormValue("secondname"))
	profile.Email = string(ctx.FormValue("email"))
	profile.Phone = string(ctx.FormValue("phone"))
	profile.Looking = string(ctx.FormValue("looking"))
	t, err := time.Parse("2006-01-02", string(ctx.FormValue("dob")))
	if err != nil {
		log.Debug("dob parsing failed - ", err)
	}
	dt := datatypes.Date(t)
	log.Debugf("dob %s, %v", string(ctx.FormValue("dob")), dt)
	profile.DOB = dt
	profile.Gender = string(ctx.FormValue("gender"))
	profile.Country = string(ctx.FormValue("country"))
	profile.Religion = string(ctx.FormValue("religion"))
	profile.Language = string(ctx.FormValue("language"))
	profile.PrefReligion = string(ctx.FormValue("prefreligion"))
	profile.PrefCountry = string(ctx.FormValue("prefcountry"))
	profile.PrefLanguage = string(ctx.FormValue("preflanguage"))
	profile.PrefAgeMin, _ = strconv.Atoi(string(ctx.FormValue("prefagemin")))
	profile.PrefAgeMax, _ = strconv.Atoi(string(ctx.FormValue("prefagemax")))
	profile.Password = string(ctx.FormValue("password"))
	ConfirmPassword := string(ctx.FormValue("confirmpassword"))
	profile.Status = "active"
	profile.Verified = false
	profile.SubscriptionType = "trial"
	log.Println("received", ctx.PostBody())

	log.Println("profiles filled from form", profile)
	log.Println("password", string(ctx.FormValue("password")))

	// Validate password before any database operations
	if len(profile.Password) == 0 || len(ConfirmPassword) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("Password or confirm password is empty")
		return fmt.Errorf("Password or confirm password is empty")
	}

	if profile.Password != ConfirmPassword {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("Password and confirm password do not match")
		return fmt.Errorf("Password and confirm password do not match")
	}

	if len(profile.Password) < 8 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("Password must be at least 8 characters long")
		return fmt.Errorf("Password must be at least 8 characters long")
	}

	// Validate password complexity (A-Z, a-z, 0-9, special chars)
	// passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@#$%^&*]).{8,}$`)
	// if !passwordRegex.MatchString(profile.Password) {
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	// 	log.Debug("Password must include A-Z, a-z, 0-9, and special characters (@#$%^&*)")
	// 	return fmt.Errorf("Password must include A-Z, a-z, 0-9, and special characters")
	// }

	if len(profile.Email) <= 0 || len(profile.Phone) <= 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("User input is empty - email or phone")
		return fmt.Errorf("User input is empty - email or phone")
	}

	// Validate if email or phone already exists
	matid, err := checkIfUserExists(profile.Email, profile.Phone)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//Land to an error page
		//ctx.Write([]byte("checkIfUserExists failed"))
		log.Debugf("checkIfUserExists failed with error - %v", err)
		return fmt.Errorf("checkIfUserExists failed with error - %v", err)
	}

	if len(matid) <= 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//show in the same form that email and phone is already present
		//ctx.Write([]byte(fmt.Sprintf("User already exists with this email or phone number please use a different one")))
		//ctx.Redirect("/api/profile/dashboard", fasthttp.StatusFound)
		log.Debugf("Matrimony id not generated for enrollment - %s, %s", profile.Email, profile.Phone)
		return fmt.Errorf("Matrimony id not generated for enrollment - %s, %s", profile.Email, profile.Phone)
	}

	// Hash the password
	var hashedPassword []byte
	if len(profile.Password) >= 8 {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			//ctx.Write([]byte("Failed to hash password"))
			log.Debugf("Failed to hash password - %v", err)
			return fmt.Errorf("Failed to hash password - %v", err)
		}
	}
	profile.Password = string(hashedPassword)
	log.Println("hashed password length", len(profile.Password))

	profile.Age = utils.CalculateAge(profile.DOB)

	fmt.Println("age calculated is ", profile.Age)

	// Insert user into the database
	enrolledUser, err := createNewProfile(&profile)
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
	enrolledUser.Matrimonyid = fmt.Sprintf("KAN%020d", enrolledUser.Id)
	mf.Value["matrimonyid"] = []string{enrolledUser.Matrimonyid}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	//ctx.Write([]byte("User account created successfully!"))
	//telemetry reporting new user created for analyzing user growth for this site
	var msg string = fmt.Sprintf("New user account created with email, phone, matrimony id: %s, %s, %s", enrolledUser.Email, enrolledUser.Phone, enrolledUser.Matrimonyid)
	utils.LogTelemetry("UserAccountCreated", msg)
	// Sent even to email server for welcome email
	err = utils.LogEmailEvent(utils.WelcomeEmail, enrolledUser.Email, enrolledUser.Matrimonyid)
	if err != nil {
		log.Errorf("failed to send welcome email to user %s:%v", enrolledUser.Email, err)
	}
	log.Infof("User account created successfully for email: %s", enrolledUser.Email)
	//Passing token type to auth middleware to create enroll token
	ctx.SetUserValue(middleware.Tokentype, middleware.ProfileTokenType)
	utils.Redirect(ctx, "/api/profile/dashboard")
	return nil
}

func createFullProfileAccountHandler(ctx *routing.Context) error {
	//1. Email id or phone number is already present => Show its already present in the same page
	//2. Once account is created successfully => redirect to new profile details page
	//3. Any other error happenes => show error message some problem happened and its not your fault come after some time

	if string(ctx.Method()) != http.MethodPost {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
		//ctx.Write([]byte("Method not allowed, please use POST method"))
		log.Debug("Method not allowed, please use POST method")
		return fmt.Errorf("Method not allowed")
	}

	// Parse request body
	var err error
	var profile Profiles
	profile.FirstName = string(ctx.FormValue("firstname"))
	profile.SecondName = string(ctx.FormValue("secondname"))
	profile.Email = string(ctx.FormValue("email"))
	profile.Phone = string(ctx.FormValue("phone"))
	profile.Looking = string(ctx.FormValue("looking"))
	t, err := time.Parse("2006-01-02", string(ctx.FormValue("dob")))
	if err != nil {
		log.Debug("dob parsing failed - ", err)
	}
	dt := datatypes.Date(t)
	log.Debugf("dob %s, %v", string(ctx.FormValue("dob")), dt)
	profile.DOB = dt
	profile.Gender = string(ctx.FormValue("gender"))
	profile.Country = string(ctx.FormValue("country"))
	profile.Religion = string(ctx.FormValue("religion"))
	profile.Language = string(ctx.FormValue("language"))
	profile.PrefReligion = string(ctx.FormValue("prefreligion"))
	profile.PrefCountry = string(ctx.FormValue("prefcountry"))
	profile.PrefLanguage = string(ctx.FormValue("preflanguage"))
	profile.PrefAgeMin, _ = strconv.Atoi(string(ctx.FormValue("prefagemin")))
	profile.PrefAgeMax, _ = strconv.Atoi(string(ctx.FormValue("prefagemax")))
	profile.Password = string(ctx.FormValue("password"))
	ConfirmPassword := string(ctx.FormValue("confirmpassword"))
	profile.Status = "active"
	profile.Verified = false
	profile.SubscriptionType = "trial"
	log.Println("received", ctx.PostBody())

	log.Println("profiles filled from form", profile)
	log.Println("password", string(ctx.FormValue("password")))

	// Validate password before any database operations
	if len(profile.Password) == 0 || len(ConfirmPassword) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("Password or confirm password is empty")
		return fmt.Errorf("Password or confirm password is empty")
	}

	if profile.Password != ConfirmPassword {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("Password and confirm password do not match")
		return fmt.Errorf("Password and confirm password do not match")
	}

	if len(profile.Password) < 8 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("Password must be at least 8 characters long")
		return fmt.Errorf("Password must be at least 8 characters long")
	}

	// Validate password complexity (A-Z, a-z, 0-9, special chars)
	// passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@#$%^&*]).{8,}$`)
	// if !passwordRegex.MatchString(profile.Password) {
	// 	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	// 	log.Debug("Password must include A-Z, a-z, 0-9, and special characters (@#$%^&*)")
	// 	return fmt.Errorf("Password must include A-Z, a-z, 0-9, and special characters")
	// }

	if len(profile.Email) <= 0 || len(profile.Phone) <= 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debug("User input is empty - email or phone")
		return fmt.Errorf("User input is empty - email or phone")
	}

	// Validate if email or phone already exists
	matid, err := checkIfUserExists(profile.Email, profile.Phone)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//Land to an error page
		//ctx.Write([]byte("checkIfUserExists failed"))
		log.Debugf("checkIfUserExists failed with error - %v", err)
		return fmt.Errorf("checkIfUserExists failed with error - %v", err)
	}

	if len(matid) > 10 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//show in the same form that email and phone is already present
		//ctx.Write([]byte(fmt.Sprintf("User already exists with this email or phone number please use a different one")))
		//ctx.Redirect("/api/profile/dashboard", fasthttp.StatusFound)
		log.Debugf("User already exists for this email and phone no - %s, %s", profile.Email, profile.Phone)
		return fmt.Errorf("User already exists for this email and phone no - %s, %s", profile.Email, profile.Phone)
	}

	// Hash the password
	var hashedPassword []byte
	if len(profile.Password) >= 8 {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			//ctx.Write([]byte("Failed to hash password"))
			log.Debugf("Failed to hash password - %v", err)
			return fmt.Errorf("Failed to hash password - %v", err)
		}
	}
	profile.Password = string(hashedPassword)
	log.Println("hashed password length", len(profile.Password))

	profile.Age = utils.CalculateAge(profile.DOB)

	fmt.Println("age calculated is ", profile.Age)

	// Insert user into the database
	enrolledUser, err := createFullProfile(&profile)
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
	enrolledUser.Matrimonyid = fmt.Sprintf("KAN%020d", enrolledUser.Id)
	mf.Value["matrimonyid"] = []string{enrolledUser.Matrimonyid}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	//ctx.Write([]byte("User account created successfully!"))
	//telemetry reporting new user created for analyzing user growth for this site
	var msg string = fmt.Sprintf("New user account created with email, phone, matrimony id: %s, %s, %s", enrolledUser.Email, enrolledUser.Phone, enrolledUser.Matrimonyid)
	utils.LogTelemetry("UserAccountCreated", msg)
	// Sent even to email server for welcome email
	err = utils.LogEmailEvent(utils.WelcomeEmail, enrolledUser.Email, enrolledUser.Matrimonyid)
	if err != nil {
		log.Errorf("failed to send welcome email to user %s:%v", enrolledUser.Email, err)
	}
	log.Infof("User account created successfully for email: %s", enrolledUser.Email)
	//Passing token type to auth middleware to create enroll token
	ctx.SetUserValue(middleware.Tokentype, middleware.ProfileTokenType)
	utils.Redirect(ctx, "/api/profile/dashboard")
	return nil
}

func checkIfUserExistsHandler(ctx *routing.Context) error {
	Email := string(ctx.FormValue("email"))
	Phone := string(ctx.FormValue("phone"))

	if len(Email) <= 0 || len(Phone) <= 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		//ctx.Write([]byte("User input is empty"))
		log.Debug("User input is empty - email, phone or password")
		return fmt.Errorf("User input is empty - email, phone or password")
	}

	type Info struct {
		Email bool `json:"email"`
		Phone bool `json:"phone"`
	}

	info := Info{}

	// Validate if email or phone already exists
	exists, err := checkIfEmailExists(Email)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debugf("checkIfEmailExists failed with error - %v", err)
		return fmt.Errorf("check email failed with error - %v", err)
	}

	info.Email = exists
	exists, err = checkIfPhoneExists(Phone)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Debugf("checkIfEmailExists failed with error - %v", err)
		return fmt.Errorf("check email failed with error - %v", err)
	}
	info.Phone = exists

	jsonData, err1 := json.Marshal(info)
	if err1 != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Debugf("Failed to marshal response - %v", err1)
		return fmt.Errorf("Failed to marshal response - %v", err1)
	}

	ctx.Write(jsonData)
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	log.Debugf("User existence check completed for email: %s, phone: %s", Email, Phone)
	return nil

}

// Add this function at the end of the file
func updateUserActivityStatus(matrimonyid string, status string) error {
	db := utils.GetDB()
	query := `UPDATE profiles SET status = ?, last_active = NOW() WHERE matrimonyid = ?`
	err := db.Exec(query, status, matrimonyid).Error
	if err != nil {
		log.Debugf("Failed to update user activity status - %v", err)
		return fmt.Errorf("Failed to update user activity status - %v", err)
	}
	return nil
}
