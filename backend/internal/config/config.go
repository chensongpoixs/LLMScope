package config

import (
	"log"
	"os"
)

type Config struct {
	Server     ServerConfig
	LlamaCpp   LlamaCppConfig
	CORS       CORSConfig
}

type ServerConfig struct {
	Address string
}

type LlamaCppConfig struct {
	Host    string
	Port    string
	BaseURL string
}

type CORSConfig struct {
	AllowOrigins []string
}

func Load() *Config {
	llamaCppHost := getEnv("LLAMA_CPP_HOST", "221.10.121.80")
	llamaCppPort := getEnv("LLAMA_CPP_PORT", "20002")
	serverPort := getEnv("SERVER_PORT", "8080")

	cfg := &Config{
		Server: ServerConfig{
			Address: ":" + serverPort,
		},
		LlamaCpp: LlamaCppConfig{
			Host:    llamaCppHost,
			Port:    llamaCppPort,
			BaseURL: "http://" + llamaCppHost + ":" + llamaCppPort,
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
	}

	log.Printf("Config loaded: llama.cpp at %s", cfg.LlamaCpp.BaseURL)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
