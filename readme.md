# KVLite

A minimal, high-performance key-value store with persistence and Pub/Sub capabilities, implemented in Go.

## Features

- **Thread-safe Operations**: Concurrent access handling with read-write mutex
- **Core Commands**: SET, GET, DEL for key-value operations
- **TCP Server**: Simple text-based protocol over TCP
- **Concurrency**: Handles multiple clients simultaneously using goroutines
- **Data Persistence**: Automatic saving to disk with atomic writes
- **Pub/Sub**: Basic publish-subscribe pattern implementation

## Usage

Build the server:
```bash
go build -o kvlite main.go persistence.go
```

Start the server:
```bash
./kvlite
```

Connect using netcat or telnet:
```
$ nc localhost 4000
> SET foo bar
OK
> GET foo
bar
> DEL foo
OK
> GET foo
(not found)
```

## Data Persistence

- Data is automatically saved to `data.json` on every write operation (SET/DEL)
- On server start, data is loaded from `data.json` if it exists
- Temporary files (`data.json.tmp`) are used during writes to prevent data loss on system failures
- The server creates `data.json` automatically on first write

## Commands

### Key-Value Operations
- `SET <key> <value>` - Store a key-value pair
- `GET <key>` - Retrieve a value by key
- `DEL <key>` - Delete a key-value pair

### Pub/Sub Operations
- `SUBSCRIBE <channel>` - Subscribe to a channel
- `PUBLISH <channel> <message>` - Publish a message to a channel
- `UNSUBSCRIBE <channel>` - Unsubscribe from a channel
