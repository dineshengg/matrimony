package config

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	LoadConfig()
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

	//redis configurations
	flagset = flag.NewFlagSet("db.redis", flag.ExitOnError)
	flagset.String("db.redis.host", "localhost", "Redis host")
	flagset.Int("db.redis.port", 6379, "Redis port")
	flagset.String("db.redis.password", "", "Redis password")
	flagset.String("db.redis.usename", "default", "Redis username")
	flag.CommandLine.AddFlagSet(flagset)

	//server configurations
	flagset = flag.NewFlagSet("server", flag.ExitOnError)
	flagset.String("server.host", "localhost", "Server host")
	flagset.Int("server.port", 8080, "Server port")
	flagset.String("server.static_dir", "../resources", "Static files directory")

}

func BindFlags() {
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)
	viper.AutomaticEnv()
}

func Args() []string {
	return flag.Args()
}
