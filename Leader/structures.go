package main

import (
	"time"
)

type Policy struct {
	BLOCK_TIME string `json:"blocktime"`
	BLOCK_SIZE string `json:"blocksize"`
}

type Nodes struct {
	IP       string `json:"nodeip"`
	TcpPort  string `json:"nodetcpport"`
	HttpPort string `json:"nodehttpport"`
}

type Hash [32]byte
type Transaction struct {
	TransactionID Hash      `json:"transactionid"`
	Priority      int       `json:"priority"`
	Threshold     int       `json:"threshold"`
	Timestamp     time.Time `json:"timestamp"`
	Deadline      time.Time `json:"deadline"`
	RemainingTime int       `json:"remainingtime"`
}

type Block struct {
	Index            int           `json:"index"`
	Timestamp        time.Time     `json:"timestamp"`
	PreviousHash     string        `json:"previousHash"`
	Hash             string        `json:"hash"`
	Nonce            int           `json:"nonce"`
	Difficulty       int           `json:"difficulty"`
	DifficultyTarget string        `json:"difficultyTarget"`
	TransactionList  []Transaction `json:"transactionList"`
}
