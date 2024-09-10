package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		AuthServer  AuthConfig
		EmailSender EmailSenderConfig
		Server      ServerConfig
		Database    DatabaseConfig
		Redis       RedisConfig
		JWT         JWTConfig
		RabbitMQ    RabbitMQConfig
		Auth        string
		Booking     string
	}
	EmailSenderConfig struct {
		SMTPServer  string
		SMTPPort    string
		Password    string
		SenderEmail string
	}
	AuthConfig struct {
		Host string
		Port string
	}

	JWTConfig struct {
		SecretKey string
	}

	ServerConfig struct {
		Host string
		Port string
	}
	DatabaseConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	RedisConfig struct {
		Host string
		Port string
	}
	RabbitMQConfig struct {
		RabbitMQ string
	}
)

func (c *Config) Load() error {

	if err := godotenv.Load(); err != nil {
		return err
	}
	c.AuthServer.Host = os.Getenv("AUTH_HOST")
	c.AuthServer.Port = os.Getenv("AUTH_PORT")

	c.Server.Host = os.Getenv("SERVER_HOST")
	c.Server.Port = os.Getenv("SERVER_PORT")

	c.Database.Host = os.Getenv("DB_HOST")
	c.Database.Port = os.Getenv("DB_PORT")
	c.Database.User = os.Getenv("DB_USER")
	c.Database.Password = os.Getenv("DB_PASS")
	c.Database.DBName = os.Getenv("DB_NAME")

	c.Redis.Host = os.Getenv("REDIS_HOST")
	c.Redis.Port = os.Getenv("REDIS_PORT")

	c.EmailSender.SMTPServer = os.Getenv("SMTP_SERVER")
	c.EmailSender.SMTPPort = os.Getenv("SMTP_PORT")
	c.EmailSender.Password = os.Getenv("EMAIL_PASS")
	c.EmailSender.SenderEmail = os.Getenv("SENDER_EMAIL")

	c.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")

	c.RabbitMQ.RabbitMQ = os.Getenv("RABBITMQ_URI")

	// pp.Println(c)

	return nil
}

func NewConfig() *Config {
	return &Config{}
}
