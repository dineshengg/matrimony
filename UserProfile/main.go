package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dineshengg/matrimony/common/config"
	"github.com/dineshengg/matrimony/common/utils"

	"github.com/dineshengg/matrimony/userprofile/login"
	"github.com/valyala/fasthttp"

	"os/signal"

	_ "github.com/lib/pq"
	routing "github.com/qiangxue/fasthttp-routing"
	log "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
)

var db *utils.MyDatabase
var redisClient *utils.RedisClient

// css static files serving
// fs will be initialized in main
var fs fasthttp.RequestHandler

func init() {

}

func clean() {
	defer db.CloseDB()
	defer redisClient.Close()
}

//TODO - add panic recovery for different services

// TODO - add both http and https support with certtificate then add a reverse proxy like nginx to offload ssl traffic
func main() {
	fmt.Println("starting web server")

	//Add signal handling for CTRL+C and graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Initialize static file handler
	fs := fasthttp.FSHandler("../resources", 1)

	// Load configuration happens at init of config package
	//err := config.LoadConfig()

	//Bind the environment variables and command line argurments
	config.BindFlags()
	log.SetLevel(log.DebugLevel)

	defer clean()

	//Initialize database and redis client in a separate go routine
	go func() {
		//// Initialize PostgreSQL connection
		dbConnStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", viper.GetString("db.postgres.host"),
			viper.GetString("db.postgres.user"), viper.GetString("db.postgres.password"), viper.GetString("db.postgres.name"),
			viper.GetInt("db.postgres.port"))

		redisAddressStr := fmt.Sprintf("%s:%s", viper.GetString("db.redis.host"), viper.GetString("db.redis.port"))
		db = utils.NewDatabaseClient(dbConnStr)
		redisClient = utils.NewRedisClient(redisAddressStr)

		if db == nil || redisClient == nil {
			log.Errorf("Failed to initialize database  - %v, or redis client - %v", db, redisClient)
			panic("Failed to initialize database or redis client")

		}

		if !utils.IsTelemetryInit() {
			log.Error("Telemtery is not initialized")
			panic("Telemtery opensearch is not initialized")
		}

		if !utils.IsTracingInitialized() {
			utils.LogTelemetry("init", "Tracing is not initialized")
		}
		if !utils.IsEmailClientInitialized() {

			utils.LogTelemetry("init", "email client is not initialized")
			log.Error("Email client is not initialized")
		}
	}()

	//all routing starts with /api/v1
	router := routing.New()
	static := router.Group("/static")
	static.Any("/*", func(ctx *routing.Context) error {
		log.Println("serving static files -", string(ctx.RequestCtx.Path()))
		fs(ctx.RequestCtx)
		return nil
	},
	)

	index := router.Group("/")
	index.Get("/", Index)
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

	srvErr := make(chan error, 1)
	go func() {

		// Start the HTTP server
		log.Println("Starting web server on :8080")
		err := fasthttp.ListenAndServe(":8080", router.HandleRequest)
		if err != nil {
			log.Println("Error starting HTTP server:", err)
			srvErr <- err
			return
		}
	}()

	//handle ctr+c gracefully
	select {
	case err := <-srvErr:
		// Error when starting HTTP server
		log.Println("Server listen error for http port 8080", err)
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

}

func Index(ctx *routing.Context) error {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/plain; charset=utf-8")
	//ctx.Write([]byte("Welcome to Kandan Matrimony application your life saviour for finding a perfect partner!"))
	filepath := "../resources/home/home.html"
	index, err := os.ReadFile(filepath)
	if err != nil {
		log.Println("Error while reading index.html file - ", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString("Internal server error")
		return fmt.Errorf("Error while reading index.html file - %v", err)
	}
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.Write(index)
	ctx.SetStatusCode(fasthttp.StatusOK)
	//check if cookie is present and user is logged in if so redirect to dashboard
	return nil
}
