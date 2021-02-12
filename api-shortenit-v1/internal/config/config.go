package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName        string
	TracerName     string
	TraceCollector string
	AliasCon       string
	MongoCon       string
	Port           string
	BrokerList     []string
	GetUrlTopic    string
}

func NewAppConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		AppName:        os.Getenv("APP_NAME"),
		TracerName:     os.Getenv("TRACER_NAME"),
		TraceCollector: os.Getenv("TRACER_COLLECTOR"),
		AliasCon:       os.Getenv("ALIAS_CON"),
		MongoCon:       os.Getenv("MONGO_CON"),
		Port:           os.Getenv("PORT"),
		BrokerList:     strings.Split(os.Getenv("KAFKA_PEERS"), ","),
		GetUrlTopic:    os.Getenv("GETURL_EVENT_TOPIC"),
	}
}
