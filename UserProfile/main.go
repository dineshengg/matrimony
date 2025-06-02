package main

import (
	"fmt"
	"log"

	"github.com/dineshengg/matrimony/common/config"
	"github.com/dineshengg/matrimony/common/utils"
	"github.com/dineshengg/matrimony/userprofile/login"
	"github.com/valyala/fasthttp"

	_ "github.com/lib/pq"
	routing "github.com/qiangxue/fasthttp-routing"
)

var db *utils.MyDatabase
var redisClient *utils.RedisClient

func init() {

}

func clean() {
	defer db.CloseDB()
	defer redisClient.Close()
}

func main() {
	fmt.Println("starting web server")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		panic("Failed to load json configuration file")
	}

	defer clean()

	//Initialize database and redis client in a separate go routine
	go func() {
		//// Initialize PostgreSQL connection
		dbConnStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

		redisAddressStr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
		db = utils.NewDatabaseClient(dbConnStr)
		redisClient = utils.NewRedisClient(redisAddressStr)
	}()

	//all routing starts with /api/v1
	router := routing.New()
	router.Get("/", Index)
	//No other route path found this handler will call the home page
	router.NotFound(Index)

	api := router.Group("/api")

	login.LoginRoutingFunctions(api)
	//fasthttp.ListenAndServe(":8080", router.HandleRequest)

	//profile.ProfileRoutingFunctions(mux)
	//extprofile.ExtProfileRoutingFunctions(mux)
	//dashboard.DashboardRoutingFunctions(mux)
	//membership.MembershipRoutingFunctions(mux)

	//wrappedMux := datasource.NewMiddleWare(mux)

	fmt.Println("Starting web server on :8080")
	log.Fatal(fasthttp.ListenAndServe(":8080", router.HandleRequest))

}

func Index(ctx *routing.Context) error {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.Write([]byte("Welcome to Kandan Matrimony application your life saviour for finding a perfect partner!"))

	//check if cookie is present and user is logged in if so redirect to dashboard
	return nil
}
