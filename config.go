package main

import (
	"log"
	"os"
)

type PostgresConfig struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

type Config struct {
	Database         PostgresConfig
	TelegramApiKey   string
	TelegramBotDebug bool
}

var Cfg Config

func init() {
	dbConfig := PostgresConfig{
		User:       getEnvStr("DB_USER", ""),
		Password:   getEnvStr("DB_PASS", ""),
		Host:       getEnvStr("DB_HOST", ""),
		Name:       getEnvStr("DB_NAME", ""),
		DisableTLS: len(getEnvStr("DB_DISABLE_TLS", "")) > 0,
	}

	Cfg = Config{
		Database:         dbConfig,
		TelegramApiKey:   getEnvStr("TG_API_KEY", ""),
		TelegramBotDebug: getEnvBool("TG_DEBUG", false),
	}

	if Cfg.Database.Host == "" {
		log.Fatal("Database env var not set")
	}

	if Cfg.TelegramApiKey == "" {
		log.Fatal("Telegram API key not set")
	}
}

func getEnvStr(key string, defaultVal string) string {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	} else {
		return defaultVal
	}
}

func getEnvBool(key string, defaultVal bool) bool {
	val, exists := os.LookupEnv(key)
	if exists && len(val) > 0 {
		return true
	} else {
		return defaultVal
	}
}

// func getEnvInt(key string, defaultVal int) int {
// 	val, exists := os.LookupEnv(key)
// 	if exists {
// 		intVal, err := strconv.Atoi(val)
// 		if err != nil {
// 			panic(fmt.Sprintf("Invalid value for key %s = %s", key, val))
// 		}
// 		return intVal
// 	} else {
// 		return defaultVal
// 	}
// }
