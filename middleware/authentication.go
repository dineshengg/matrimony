package middleware

import (
	"fmt"
	"time"

	"sync"

	"github.com/golang-jwt/jwt/v5"
	routing "github.com/qiangxue/fasthttp-routing"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// 512 byte secret key for JWT signing
const secret_key string = "95e6f08f3aa2401e0e979d9b6f9b8d1fcc57814516d150c148b32e1de93cf9176ae7f92eebc0bcca00a1610696421f1526920d66c048919795684833b109181abc939cef8809846fb3c10af09d9decaa5426e595eb23cc7ef509b6fb285b14371223325dc69dd1937813f53d3e79862f3b558e0a0a55ca873f572bf7a472cc344a9496d1e5b2534cb011fdb6ee2223229141ed7df245c2e57a08ddc39e6846e2cbf7da026ed721992c20e752fd4cd2b0431f4b742c45dd178fbd74d86525eba2eede1726de86ac01621f6c00406cccda0e684b01163463b62deae166636b117bd536a17b25c3adeda4d3fa09b4f19067fb3321b9fd860bb34e36d9176419b58f752631852095340ab13a698e54131ed598bdcf715df5fa5e18a65cafa2b59544aa86a2467900335fafda864db55ac621b8b0f39667ff7f1f0005036781de6475a4bd4719d1dd7d9ae02d9a6d15e1ab82ed9b018eaee79a7b73339c408e5d4863b723894c6d77c969de1eef5bd74956d94a6bd5371ab262b0418ae878e77e82346be4967aadde83cc0d90d1932a6ee5cd4434a10fd8e1d0a86495a275661b38fd72adac2f1224008db6e3af9b4838247f87b844288e9e69fc95b30956e074745abd922ca367724e0259c78520b2534375ae6dadbca05c010a07deac0a04b5f611c7b98ad4b2831ed85162385eb31d1b166c58d338007c94d7207570f3771c6741"
const Tokentype string = "tokentype"
const Token string = "token"

const (
	//JWT token type
	EnrollTokenType  = 0
	ProfileTokenType = 1
)

type Authentication struct {
	auth_router *routing.RouteGroup
}

// global variable as it hold states that are common for all requests not for any specific request
var auth *Authentication
var once sync.Once

type ProfileInfo struct {
	FirstName            string `json:"firstname"`
	SecondName           string `json:"secondname"`
	EnrollToken          Enroll `json:"enroll"`
	jwt.RegisteredClaims `json:"registeredclaims"`
}

type Enroll struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Looking     string `json:"looking"`
	MatrimonyID string `json:"matrimonyid"`
}

func NewAuthentication(router *routing.RouteGroup) *Authentication {
	if auth == nil {
		once.Do(func() {
			auth = &Authentication{
				auth_router: router,
			}
		})
	}
	return auth
}

func (auth *Authentication) Authenticate(ctx *routing.Context) error {
	//raw request is no copyable need to use it as a pointer only
	jwt_token := ctx.Request.Header.Cookie("jwt_token")

	if len(jwt_token) == 0 {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return fmt.Errorf("JWT token is missing")
	}

	//unmarshal the JWT token to get the profile info
	claims := ProfileInfo{}
	_, err := jwt.ParseWithClaims(string(jwt_token), &claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret_key), nil // Use the same secret as used during token generation
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return fmt.Errorf("JWT token is invalid or tampered: %v", err)
	}
	// Extract profile info from claims
	// profileInfo := ProfileInfo{
	// 	FirstName:  claims["firstname"].(string),
	// 	SecondName: claims["secondname"].(string),
	// 	Email:      claims["email"].(string),
	// 	Role:       claims["role"].(string),
	// 	Exp:        int64(claims["exp"].(float64)), // Convert float64 to int64
	// }
	// Set profile info in context for further use
	// Note: ctx.Set is not a standard method in fasthttp, you might want to use a custom context or pass it differently
	// ctx.Set("profileinfo", profileInfo) // Uncomment if you have a custom context that supports Set method
	// For demonstration, we will just log the profile info
	log.Printf("Profile Info: %+v\n", claims)
	// If you want to use ctx.Set, you need to define a custom context or use a map

	if claims.FirstName == "" || claims.SecondName == "" || claims.FirstName == "newuser" || claims.SecondName == "newuser" {
		ctx.SetUserValue(Tokentype, EnrollTokenType)
	} else {
		ctx.SetUserValue(Tokentype, ProfileTokenType)
	}
	ctx.Set(Token, claims)
	return nil
}

func (auth *Authentication) CreateJWTToken(ctx *routing.Context) error {

	//1. Email id or phone number is not present => redirect to login page with GET method
	//2. Password is in invalid => Show error message in the same page
	//3. Any web server error while processing this request => Show a generic error page or same page with error message
	//4. Successful login => redirect to profile dashboard
	log.Debug("create jwt token for user")
	var profile ProfileInfo
	jwtRegisteredClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(ctx.Time().Add(time.Hour * 48)), // Set expiration to 48 hours
		IssuedAt:  jwt.NewNumericDate(ctx.Time()),
		Issuer:    "kandanmatrimony.com",
		Subject:   string(ctx.FormValue("email")), // Use email as subject
		NotBefore: jwt.NewNumericDate(ctx.Time()),
		ID:        fmt.Sprintf("%d", ctx.Time().UnixNano()), // Unique ID for the token
		Audience:  []string{"profileusers"},
	}
	tokentype := ctx.Value(Tokentype)
	switch tokentype {
	case EnrollTokenType:
		log.Debug("creating enroll token not supported yet")
		panic("Enroll token creation is not supported yet, please use profile token creation")
	case ProfileTokenType:
		// Handle Profile Token Creation
		// TODO - Get the current enroll token from request and form the whole claims after user creates all mandatory fields
		// otherwise they continue with enroll token which will be redirecting again to fill profile into page
		log.Debug("creating profile token")
		firstname := string(ctx.FormValue("firstname"))
		secondname := string(ctx.FormValue("secondname"))

		// jwt_token := ctx.Request.Header.Cookie("jwt_token")
		// if len(jwt_token) == 0 {
		// 	ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		// 	return fmt.Errorf("Previous enroll JWT token is missing")
		// }

		// //unmarshal the JWT token to get the profile info
		// enroll_temp := Enroll{}
		// _, err := jwt.ParseWithClaims(string(jwt_token), &enroll_temp, func(token *jwt.Token) (interface{}, error) {
		// 	// Validate the signing method
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// 	}
		// 	return []byte(secret_key), nil // Use the same secret as used during token generation
		// })
		// if err != nil {
		// 	ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		// 	return fmt.Errorf("JWT token is invalid or tampered: %v", err)
		// }
		profile = ProfileInfo{
			FirstName:  firstname,
			SecondName: secondname,
			EnrollToken: Enroll{
				Email:       string(ctx.FormValue("email")),
				MatrimonyID: string(ctx.FormValue("matrimonyid")),
				Phone:       string(ctx.FormValue("phone")),
				Looking:     string(ctx.FormValue("looking")),
			},
			RegisteredClaims: jwtRegisteredClaims,
		}
	default:
		log.Debug("unknown token type")
		return fmt.Errorf("Unknown token type")
	}

	secret := []byte(secret_key) // Use a secure secret in production
	var token *jwt.Token
	if tokentype == EnrollTokenType {
		//reduce signing method to 256 if there is a latency
		//token = jwt.NewWithClaims(jwt.SigningMethodHS512, EnrollToken)
	} else {
		//use 512 bit signing method for profile token
		token = jwt.NewWithClaims(jwt.SigningMethodHS512, profile)
	}

	tokenString, err := token.SignedString(secret)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return fmt.Errorf("Failed to generate token signing - %v", err)
	}

	//set this context to be used by other handlers like
	ctx.Set(Token, profile)
	log.Debugf("generated JWT token: %s,  %v", tokenString, profile)
	// Set JWT as a cookie (optional, you can also send in response body)
	cookie := fasthttp.AcquireCookie()
	cookie.SetKey("jwt_token")
	cookie.SetValue(tokenString)
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	cookie.SetSecure(true)
	cookie.SetMaxAge(24 * 3600)
	ctx.Response.Header.SetCookie(cookie)
	ctx.Response.Header.Set("Authorization", "Bearer "+tokenString)
	return nil
}
