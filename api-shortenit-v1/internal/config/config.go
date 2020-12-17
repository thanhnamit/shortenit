package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	AppName string
	TracerName string
	TraceCollector string
	AliasCon string
	MongoCon string
	Port string
}

func NewAppConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		AppName: os.Getenv("APP_NAME"),
		TracerName: os.Getenv("TRACER_NAME"),
		TraceCollector: os.Getenv("TRACER_COLLECTOR"),
		AliasCon: os.Getenv("ALIAS_CON"),
		MongoCon: os.Getenv("MONGO_CON"),
		Port: os.Getenv("PORT"),
	}
}