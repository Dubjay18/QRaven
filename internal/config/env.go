package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Server       ServerConfiguration
	Database     Database
	TestDatabase Database
	App          App
	Cloudinary   Cloudinary
	// IPStack      IPStack
	// Mail         MAIL
	Redis Redis
}

type BaseConfig struct {
	SERVER_PORT                      string  `mapstructure:"SERVER_PORT"`
	SERVER_SECRET                    string  `mapstructure:"SERVER_SECRET"`
	SERVER_ACCESSTOKENEXPIREDURATION int     `mapstructure:"SERVER_ACCESSTOKENEXPIREDURATION"`
	REQUEST_PER_SECOND               float64 `mapstructure:"REQUEST_PER_SECOND"`
	TRUSTED_PROXIES                  string  `mapstructure:"TRUSTED_PROXIES"`
	EXEMPT_FROM_THROTTLE             string  `mapstructure:"EXEMPT_FROM_THROTTLE"`

	APP_NAME                string `mapstructure:"APP_NAME"`
	APP_MODE                string `mapstructure:"APP_MODE"`
	APP_URL                 string `mapstructure:"APP_URL"`
	MAGIC_LINK_DURATION     int    `mapstructure:"MAGIC_LINK_DURATION"`
	RESET_PASSWORD_DURATION int    `mapstructure:"RESET_PASSWORD_DURATION"`

	DB_HOST               string `mapstructure:"DB_HOST"`
	DB_PORT               string `mapstructure:"DB_PORT"`
	DB_CONNECTION         string `mapstructure:"DB_CONNECTION"`
	TIMEZONE              string `mapstructure:"TIMEZONE"`
	SSLMODE               string `mapstructure:"SSLMODE"`
	USERNAME              string `mapstructure:"USERNAME"`
	PASSWORD              string `mapstructure:"PASSWORD"`
	DB_NAME               string `mapstructure:"DB_NAME"`
	MIGRATE               bool   `mapstructure:"MIGRATE"`
	CLOUDINARY_CLOUD_NAME string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CLOUDINARY_API_KEY    string `mapstructure:"CLOUDINARY_API_KEY"`
	CLOUDINARY_API_SECRET string `mapstructure:"CLOUDINARY_API_SECRET"`

	TEST_DB_HOST       string `mapstructure:"TEST_DB_HOST"`
	TEST_DB_PORT       string `mapstructure:"TEST_DB_PORT"`
	TEST_DB_CONNECTION string `mapstructure:"TEST_DB_CONNECTION"`
	TEST_TIMEZONE      string `mapstructure:"TEST_TIMEZONE"`
	TEST_SSLMODE       string `mapstructure:"TEST_SSLMODE"`
	TEST_USERNAME      string `mapstructure:"TEST_USERNAME"`
	TEST_PASSWORD      string `mapstructure:"TEST_PASSWORD"`
	TEST_DB_NAME       string `mapstructure:"TEST_DB_NAME"`
	TEST_MIGRATE       bool   `mapstructure:"TEST_MIGRATE"`

	IPSTACK_KEY      string `mapstructure:"IPSTACK_KEY"`
	IPSTACK_BASE_URL string `mapstructure:"IPSTACK_BASE_URL"`

	MAIL_SERVER   string `mapstructure:"MAIL_SERVER"`
	MAIL_PASSWORD string `mapstructure:"MAIL_PASSWORD"`
	MAIL_USERNAME string `mapstructure:"MAIL_USERNAME"`
	MAIL_PORT     string `mapstructure:"MAIL_PORT"`

	REDIS_PORT string `mapstructure:"REDIS_PORT"`
	REDIS_HOST string `mapstructure:"REDIS_HOST"`
	REDIS_DB   string `mapstructure:"REDIS_DB"`
}

func (config *BaseConfig) SetupConfigurationn() *Configuration {
	trustedProxies := []string{}
	exemptFromThrottle := []string{}
	json.Unmarshal([]byte(config.TRUSTED_PROXIES), &trustedProxies)
	json.Unmarshal([]byte(config.EXEMPT_FROM_THROTTLE), &exemptFromThrottle)

	if config.SERVER_PORT == "" {
		config.SERVER_PORT = os.Getenv("PORT")
	}
	return &Configuration{
		Server: ServerConfiguration{
			Port:                      config.SERVER_PORT,
			Secret:                    config.SERVER_SECRET,
			AccessTokenExpireDuration: config.SERVER_ACCESSTOKENEXPIREDURATION,
			RequestPerSecond:          config.REQUEST_PER_SECOND,
			TrustedProxies:            trustedProxies,
			ExemptFromThrottle:        exemptFromThrottle,
		},
		App: App{
			Name:                  config.APP_NAME,
			Mode:                  config.APP_MODE,
			Url:                   config.APP_URL,
			MagicLinkDuration:     config.MAGIC_LINK_DURATION,
			ResetPasswordDuration: config.RESET_PASSWORD_DURATION,
		},
		Database: Database{
			DB_HOST:       config.DB_HOST,
			DB_PORT:       config.DB_PORT,
			DB_CONNECTION: config.DB_CONNECTION,
			USERNAME:      config.USERNAME,
			PASSWORD:      config.PASSWORD,
			TIMEZONE:      config.TIMEZONE,
			SSLMODE:       config.SSLMODE,
			DB_NAME:       config.DB_NAME,
			Migrate:       config.MIGRATE,
		},
		TestDatabase: Database{
			DB_HOST:       config.TEST_DB_HOST,
			DB_PORT:       config.TEST_DB_PORT,
			DB_CONNECTION: config.TEST_DB_CONNECTION,
			USERNAME:      config.TEST_USERNAME,
			PASSWORD:      config.TEST_PASSWORD,
			TIMEZONE:      config.TEST_TIMEZONE,
			SSLMODE:       config.TEST_SSLMODE,
			DB_NAME:       config.TEST_DB_NAME,
			Migrate:       config.TEST_MIGRATE,
		},
		Cloudinary: Cloudinary{
			CloudName: config.CLOUDINARY_CLOUD_NAME,
			ApiKey:    config.CLOUDINARY_API_KEY,
			ApiSecret: config.CLOUDINARY_API_SECRET,
		},
		// IPStack: IPStack{
		// 	Key:     config.IPSTACK_KEY,
		// 	BaseUrl: config.IPSTACK_BASE_URL,
		// },

		// Mail: MAIL{
		// 	Server:   config.MAIL_SERVER,
		// 	Password: config.MAIL_PASSWORD,
		// 	Port:     config.MAIL_PORT,
		// 	Username: config.MAIL_USERNAME,
		// },

		Redis: Redis{
			REDIS_PORT: config.REDIS_PORT,
			REDIS_HOST: config.REDIS_HOST,
			REDIS_DB:   config.REDIS_DB,
		},
	}
}
