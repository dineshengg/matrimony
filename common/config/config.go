package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	log.Info("Loading configuration...")
	LoadConfig()
	BindFlags()
	log.Info("Configuration loaded successfully")
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig() {
	//database configurations
	flagset := flag.NewFlagSet("db.postgres", flag.ExitOnError)
	flagset.String("db.postgres.user", "postgres", "Database user")
	flagset.String("db.postgres.password", "postgres", "Database password")
	flagset.String("db.postgres.name", "postgres", "Database name")
	flagset.String("db.postgres.host", "localhost", "Database host")
	flagset.Int("db.postgres.port", 5432, "Database port")
	flag.CommandLine.AddFlagSet(flagset)

	//memcache configurations
	flagset = flag.NewFlagSet("memcache.servers", flag.ExitOnError)
	flagset.StringSlice("memcache.servers", []string{"127.0.0.1:11211"}, "Memcache servers (comma separated)")
	flagset.Duration("memcache.timeout", 1000, "Memcache operation timeout in milliseconds")
	flagset.Int("memcache.max_idle_conns", 10, "Maximum idle connections for memcache")
	flag.CommandLine.AddFlagSet(flagset)

	//server configurations
	flagset = flag.NewFlagSet("server", flag.ExitOnError)
	flagset.String("server.host", "localhost", "Server host")
	flagset.Int("server.port", 8080, "Server port")
	flagset.String("server.static_dir", "../resources", "Static files directory")
	flag.CommandLine.AddFlagSet(flagset)

	//grafana loki configurations
	flagset = flag.NewFlagSet("grafanaloki", flag.ExitOnError)
	flagset.String("loki.app", "userprofile", "Application name using this grafana loki service")
	flagset.String("loki.environment", "development", "Environment used by this application")
	flag.CommandLine.AddFlagSet(flagset)

	//resources file path
	exePath, err := os.Executable()
	if err != nil {
		//TODO - Hardcode server path
	}
	log.Info("Setting resources file path to: ", exePath)
	exePath = exePath[:len(exePath)-len("Userprofile/main")]
	log.Info("Setting resources file path to: ", exePath)

	flagset = flag.NewFlagSet("resources", flag.ExitOnError)
	flagset.String("resources.filepath", exePath, "Path to the resources directory")
	flag.CommandLine.AddFlagSet(flagset)

}

func BindFlags() {
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
	viper.AutomaticEnv()
}

func Args() []string {
	return flag.Args()
}
