# KVLite

A minimal, high-performance key-value store with persistence and Pub/Sub capabilities, implemented in Go.

## Features

- **Thread-safe Operations**: Concurrent access handling with read-write mutex
- **Core Commands**: SET, GET, DEL for key-value operations
- **Pub/Sub System**: Publish/Subscribe pattern with channel support
- **TCP Server**: Simple text-based protocol over TCP (default port: 4000)
- **Concurrency**: Handles multiple clients simultaneously using goroutines
- **Data Persistence**: Automatic saving to disk with atomic writes
- **Graceful Shutdown**: Handles client disconnections cleanly

## Usage

Build the server:
```bash
go build -o kvlite main.go persistence.go pubsub.go
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

## Available Commands

### Key-Value Operations
- `SET <key> <value>` - Store a key-value pair
  - Example: `SET username john_doe`
  - Returns: `OK` on success

- `GET <key>` - Retrieve a value by key
  - Example: `GET username`
  - Returns: The value if found, `(nil)` if key doesn't exist

- `DEL <key>` - Delete a key-value pair
  - Example: `DEL username`
  - Returns: `1` if key was deleted, `0` if key didn't exist

### Pub/Sub Operations
- `SUBSCRIBE <channel>` - Subscribe to a channel
  - Example: `SUBSCRIBE news`
  - Returns: `> Subscribed to <channel>`

- `UNSUBSCRIBE <channel>` - Unsubscribe from a channel
  - Example: `UNSUBSCRIBE news`
  - Returns: `> Unsubscribed from <channel>`

- `PUBLISH <channel> <message>` - Publish a message to a channel
  - Example: `PUBLISH news "Hello, subscribers!"`
  - Returns: `> Published to <channel>`

### Connection Management
- `QUIT` or `EXIT` - Close the connection
  - Returns: `> BYE`
