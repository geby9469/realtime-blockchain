package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var nodes map[string]string
var blockPool []Block
var mutex = new(sync.Mutex)

func main() {
	leaderIp, leaderHttpPort, leaderTcpPort, _ := loadEnv()
	fmt.Printf("Leader Node Information [IP:%s][HTTP PORT:%s][TCP PORT:%s]\n", leaderIp, leaderHttpPort, leaderTcpPort)

	runtime.GOMAXPROCS(runtime.NumCPU())

	go addNode(leaderTcpPort, leaderHttpPort)
	go receiveFromClient(leaderHttpPort)

	select {}
}

func receiveFromClient(leaderHttpPort string) {
	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var transactions []Transaction
		if err := decoder.Decode(&transactions); err != nil {
			log.Fatal(err)
		}

		go func(transactions []Transaction) {
			if len(nodes) == 0 {
				log.Fatal("No connected node")
			}
			n := make([]string, 0, len(nodes))
			for key := range nodes {
				n = append(n, key)
			}
			rand.New(rand.NewSource(time.Now().Unix()))
			rand.Shuffle(len(n), func(i, j int) { n[i], n[j] = n[j], n[i] })

			if !sendTxToNode(n[0], nodes[n[0]], transactions) {
				os.Exit(1)
			}
		}(transactions)
	})

	http.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var block Block
		if err := decoder.Decode(&block); err != nil {
			log.Fatal(err)
		}

		// need to synchronize the block pool with processes. -> grpc
	})
	http.ListenAndServe(fmt.Sprintf(":%s", leaderHttpPort), nil)
}

func sendTxToNode(nodeAddr, httpport string, transactions []Transaction) bool {
	url := fmt.Sprintf("http://%s/transactions", strings.Split(nodeAddr, ":")[0]+":"+httpport)

	txs, _ := json.Marshal(transactions)
	httptxs := bytes.NewBuffer([]byte(txs))
	response, err := http.Post(url, "application/json", httptxs)
	if err != nil {
		log.Println(err)
		return false
	}

	fmt.Printf("[TRANSMIT TRANSACTIONS TO MINER][Status : %d][Node: %s]\n", response.StatusCode, nodeAddr)
	return true
}

func addNode(leaderTcpPort, leaderHttpPort string) {
	nodes = make(map[string]string)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", leaderTcpPort))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()

			for {
				addr := strings.Split(c.LocalAddr().String(), ":")
				node := Nodes{
					IP:       addr[0],
					TcpPort:  addr[1],
					HttpPort: leaderHttpPort,
				}
				encoder := json.NewEncoder(c)
				if err := encoder.Encode(node); err != nil {
					delete(nodes, c.RemoteAddr().String())
					log.Println("Encode connection was forcibly closed by the remote host")
					return
				}

				decoder := json.NewDecoder(c)
				if err := decoder.Decode(&node); err != nil {
					delete(nodes, c.RemoteAddr().String())
					log.Println("Decode connection was forcibly closed by the remote host")
					return
				}
				nodei := fmt.Sprintf("%s:%s", node.IP, node.TcpPort)
				nodes[nodei] = node.HttpPort

				getNodes(nodes) // heartbeat with miner

				setPolicy(c) // set policy of network

				time.Sleep(time.Second * 5)
			}
		}(conn)
	}
}

func setPolicy(c net.Conn) {
	policy := loadPolicy()

	encoder := json.NewEncoder(c)
	if err := encoder.Encode(policy); err != nil {
		log.Println("Policy Encode connection was forcibly closed by the remote host")
		return
	}
}

func getNodes(nodes map[string]string) {
	fmt.Print("Connected Nodes : ")
	for addr, httpport := range nodes {
		fmt.Printf("[TCP: %s][HTTP PORT: %s]", addr, httpport)
	}
	fmt.Println()
}
