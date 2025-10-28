package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host, Port, User, Password, Name, SSLMode string
}

type RedisConfig struct {
	Addr, Password, DB, RedisHost string
}

type AppConfig struct {
	Port string
}
type Twillow struct {
	Sid,
	Token,
	WatappNum string
}
type Auth struct {
	key string
}

type Config struct {
	DB    DBConfig
	Redis RedisConfig
	App   AppConfig
	Twill *Twillow
	Key   *Auth
}

var (
	cfg  *Config
	once sync.Once
)

func LoadEnv() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using system environment variables")
		}

		cfg = &Config{
			DB: DBConfig{
				Host:     os.Getenv("DB_HOST"),
				Port:     os.Getenv("DB_PORT"),
				User:     os.Getenv("DB_USER"),
				Password: os.Getenv("DB_PASSWORD"),
				Name:     os.Getenv("DB_NAME"),
				SSLMode:  os.Getenv("DB_SSLMODE"),
			},
			Redis: RedisConfig{
				Addr:      os.Getenv("REDIS_ADDR"),
				Password:  os.Getenv("REDIS_PASSWORD"),
				DB:        os.Getenv("REDIS_DB"),
				RedisHost: os.Getenv("REDIS_HOST"),
			},
			App: AppConfig{
				Port: os.Getenv("APP_PORT"),
			},
			Twill: &Twillow{
				Sid:       os.Getenv("ACCOUNT_SID"),
				Token:     os.Getenv("Twillo_AUTH_TOKEN"),
				WatappNum: os.Getenv("Twillo_WATAPPS_NUM"),
			},
			Key: &Auth{
				key: os.Getenv("SECRATE_KEY"),
			},
		}

		log.Println("Environment variables loaded successfully")
	})

	return cfg
}

func (c *Config) GetDBURL() string {
	db := c.DB
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

func (c *Config) GetRedis() *RedisConfig {
	return &c.Redis
}

func (c *Config) GetAppPort() string {
	return fmt.Sprintf(":%s", c.App.Port)
}
func (c *Config) GetTillowInfo() *Twillow {
	return c.Twill
}
func (c *Config) GetSecrateKey() string {
	return c.Key.key
}
