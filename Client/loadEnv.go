package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() (string, string) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	leaderIP := os.Getenv("LEADER_IP")
	leaderHttpPort := os.Getenv("LEADER_HTTP_PORT")
	return leaderIP, leaderHttpPort

}
