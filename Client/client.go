package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Hash [32]byte
type Transaction struct {
	TransactionID Hash
	Priority      int
	Threshold     int
	Timestamp     time.Time
	Deadline      time.Time
	RemainingTime int
}

func main() {
	// submit transactions into the blockchain
	submitTransaction(createTransaction(20), 10)
}

func submitTransaction(transactions []Transaction, sendRate int) {
	for i := 0; i < len(transactions); i += sendRate {
		txs, _ := json.Marshal(transactions[i : i+sendRate])
		httptxs := bytes.NewBuffer([]byte(txs))

		var leaderIP, leaderHttpPort = loadEnv()
		url := fmt.Sprintf("http://%s:%s/transactions", leaderIP, leaderHttpPort)
		response, err := http.Post(url, "application/json", httptxs)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("[Status : %d][Transaction %d~%d]\n", response.StatusCode, i, i+sendRate-1)
		time.Sleep(time.Second)
	}
}

func createTransaction(transactionSize int) []Transaction {
	var transactions []Transaction
	for i := 0; i < transactionSize; i++ {
		transaction := new(Transaction)

		now := time.Now()
		deadline := now.Add(time.Second * 10)
		priority, threshold := rand.Intn(2)+1, rand.Intn(2)+1
		remainingtime := 0

		transaction.Priority = priority
		transaction.Threshold = threshold
		transaction.Timestamp = now
		transaction.Deadline = deadline
		transaction.RemainingTime = remainingtime

		transactionSum := bytes.Join([][]byte{
			[]byte(strconv.Itoa(i)),
			[]byte(stringWithCharset(16, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")),
			[]byte(now.String()),
			[]byte(deadline.String()),
			[]byte(strconv.Itoa(priority)),
			[]byte(strconv.Itoa(threshold)),
			[]byte(strconv.Itoa(remainingtime)),
		}, []byte{})
		transaction.TransactionID = sha256.Sum256(transactionSum)

		transactions = append(transactions, *transaction)
	}

	return transactions
}

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-32:]
	}

	copy(h[32-len(b):], b)
}

func stringWithCharset(length int, charset string) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
