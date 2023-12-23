package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() (string, string, string, string) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	leaderIP := os.Getenv("LEADER_IP")
	leaderHttpPort := os.Getenv("LEADER_HTTP_PORT")
	leaderTcpPort := os.Getenv("LEADER_TCP_PORT")
	leaderNodesTcpPort := os.Getenv("LEADER_NODES_TCP_PORT")
	return leaderIP, leaderHttpPort, leaderTcpPort, leaderNodesTcpPort

}

func loadPolicy() Policy {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	policy := Policy{
		BLOCK_TIME: os.Getenv("BLOCK_TIME"),
		BLOCK_SIZE: os.Getenv("BLOCK_SIZE"),
	}

	return policy
}
