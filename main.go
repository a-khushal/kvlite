package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type KVStore struct {
	mu   sync.RWMutex      // read write mutex
	data map[string]string // actual in-memory db
}

// Constructor function
func NewKVStore() *KVStore {
	return &KVStore{data: make(map[string]string)}
}

func (kv *KVStore) Set(key, value string) {
	kv.mu.Lock() // exclusive lock for writing
	kv.data[key] = value
	kv.mu.Unlock()
}

// RWMutex to let concurrent reads to happen
func (kv *KVStore) Get(key string) (string, bool) {
	kv.mu.RLock() // read only lock, shared lock for reading
	v, ok := kv.data[key]
	kv.mu.RUnlock()
	return v, ok
}

func (kv *KVStore) Del(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock() // schedule Unlock no matter where it returns

	_, exists := kv.data[key]
	if exists {
		delete(kv.data, key)
		return true
	}

	return false
}

func handleConnection(conn net.Conn, store *KVStore) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("client connected: %s\n", addr)

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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v\n", err)
			continue
		}
		go handleConnection(conn, store)
	}
}
