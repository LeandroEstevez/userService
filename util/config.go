package util

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	db "userMicroService/db/sqlc"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// Config stores the configuration of the app
// The values are read by viper from a congif file or enviroment variables
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig read configuration from file
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadKafkaConfig(configFile string) kafka.ConfigMap {
	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	m["group.id"] = "email-service"
	m["auto.offset.reset"] = "earliest"

	return m
}

var Store db.Store
var Conf Config

func SetUpConnAndStore() {
	var err error
	Conf, err = LoadConfig(".")
	fmt.Println("Config here")
	fmt.Println(Conf)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(Conf.DBDriver, Conf.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	Store = db.NewStore(conn)
	fmt.Println("This is Store", Store)
}
