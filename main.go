package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// `chan` is short for channel: a built-in Go type used for communication between goroutines.
// A channel lets one goroutine send data to another safely, without extra locks.
// Publisher’s goroutine is sending messages to the channel, and the subscriber’s writer goroutine(the go func() { for msg := range client.send { … } }()), [anonymous function written below] is reading from the channel.
type Client struct {
	conn net.Conn
	send chan string
}

type KVStore struct {
	mu          sync.RWMutex
	data        map[string]string
	subscribers map[string]map[*Client]bool
	subsMu      sync.RWMutex
}

func NewKVStore() *KVStore {
	return &KVStore{
		data:        make(map[string]string),
		subscribers: make(map[string]map[*Client]bool),
	}
}

func (kv *KVStore) Set(key, value string) {
	kv.mu.Lock()
	kv.data[key] = value
	kv.mu.Unlock()
	_ = kv.SaveToFile("data.json")
}

func (kv *KVStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	v, ok := kv.data[key]
	kv.mu.RUnlock()
	return v, ok
}

func (kv *KVStore) Del(key string) bool {
	kv.mu.Lock()
	_, exists := kv.data[key]
	if exists {
		delete(kv.data, key)
	}
	kv.mu.Unlock()
	if exists {
		_ = kv.SaveToFile("data.json")
	}
	return exists
}

func handleConnection(conn net.Conn, store *KVStore) {
	client := &Client{
		conn: conn,
		send: make(chan string, 10), // buffered channel for outgoing messages
	}

	// a goroutine whose sole job is to send messages from the client's buffer (client.send) out to the TCP connection. an anonymous function
	go func() {
		for msg := range client.send {
			fmt.Fprintln(client.conn, msg)
		}
	}()

	addr := conn.RemoteAddr().String()
	log.Printf("client connected: %s\n", addr)
	defer func() {
		conn.Close()
		log.Printf("client disconnected: %s\n", addr)
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 3)
		cmd := strings.ToUpper(parts[0])

		switch cmd {
		case "SET":
			if len(parts) < 3 {
				fmt.Fprintln(conn, "> Err usage: SET key value")
				continue
			}
			store.Set(parts[1], parts[2])
			fmt.Fprintln(conn, "> OK")

		case "GET":
			if len(parts) < 2 {
				fmt.Fprintln(conn, "> Err usage: GET key")
				continue
			}
			value, exists := store.Get(parts[1])
			if exists {
				fmt.Fprintln(conn, "> ", value)
			} else {
				fmt.Fprintln(conn, "> (nil)")
			}

		case "DEL":
			if len(parts) < 2 {
				fmt.Fprintln(conn, "> Err usage: DEL key")
				continue
			}
			if store.Del(parts[1]) {
				fmt.Fprintln(conn, "> 1")
			} else {
				fmt.Fprintln(conn, "> 0")
			}

		case "SUBSCRIBE":
			if len(parts) < 2 {
				fmt.Fprintln(conn, "> Err usage: SUBSCRIBE channel")
				continue
			}
			store.Subscribe(parts[1], client)
			fmt.Fprintln(conn, "> Subscribed to", parts[1])

		case "UNSUBSCRIBE":
			if len(parts) < 2 {
				fmt.Fprintln(conn, "> Err usage: UNSUBSCRIBE channel")
				continue
			}
			store.Unsubscribe(parts[1], client)
			fmt.Fprintln(conn, "> Unsubscribed from", parts[1])

		case "PUBLISH":
			if len(parts) < 3 {
				fmt.Fprintln(conn, "> Err usage: PUBLISH channel message")
				continue
			}
			store.Publish(parts[1], parts[2])
			fmt.Fprintln(conn, "> Published to", parts[1])

		case "QUIT", "EXIT":
			fmt.Fprintln(conn, "> BYE")
			log.Printf("client disconnected (requested): %s\n", addr)
			return

		default:
			fmt.Fprintln(conn, "> Err unknown command")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read error from %s: %v\n", addr, err)
	}
	log.Printf("client disconnected: %s\n", addr)
}

func main() {
	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	defer listener.Close()
	log.Println("KVLite listening on :4000")

	store := NewKVStore()
	if err := store.LoadFromFile("data.json"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to load data file: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		go handleConnection(conn, store)
	}
}
