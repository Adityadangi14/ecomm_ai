package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	RabbitMQ RabbitMQ
	Postgres PostgresConfig
	Redis    RedisConfig
	Weaviate WeaviateConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

type RabbitMQ struct {
	Host           string
	Port           string
	User           string
	Password       string
	Exchange       string
	Queue          string
	RoutingKey     string
	ConsumerTag    string
	WorkerPoolSize int
}

type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

type WeaviateConfig struct {
	Host      string
	Scheme    string
	TimeoutMs int
}

type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

type KafkaConfig struct {
	Brokers          []string
	Topic            string
	GroupID          string
	ClientID         string
	TimeoutMs        int
	MinBytes         int
	MaxBytes         int
	CommitIntervalMs int
	RequiredAcks     int
	Compression      string
	SASL             KafkaSASL
	TLS              KafkaTLS
	MaxRetry         int
}

type KafkaSASL struct {
	Enabled   bool
	Mechanism string
	Username  string
	Password  string
}

type KafkaTLS struct {
	Enabled            bool
	InsecureSkipVerify bool
}

func LoadConfig(filepath string) (*viper.Viper, error) {
	v := viper.New()

	// Ensure extension exists
	if !strings.HasSuffix(filepath, ".yml") && !strings.HasSuffix(filepath, ".yaml") {
		filepath = filepath + ".yml"
	}

	// Configure Viper
	v.SetConfigFile(filepath)
	v.SetConfigType("yaml")

	// Add search paths
	v.AddConfigPath(".")           // current dir
	v.AddConfigPath("./config")    // config folder
	v.AddConfigPath("../config")   // parent (useful for running from subfolders)
	v.AddConfigPath("/app/config") // Docker path

	// Read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found: " + filepath)
		}
		return nil, err
	}

	return v, nil
}

// Parse into struct
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

// high-level helper
func GetConfig(configName string) (*Config, error) {
	path := GetConfigPath(configName)

	v, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(v)
}

// helper to pick file
func GetConfigPath(env string) string {
	fmt.Println("Got from env", env)
	if env == "docker" {
		fmt.Println("Returned value", "config/docker-config.yml")
		return "config/docker-config.yml"
	}
	return "config/config-local.yml"
}
