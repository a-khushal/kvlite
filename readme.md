# KVLite

A simple in-memory key-value store with persistence, implemented in Go.

## Features

- Thread-safe operations with read-write mutex
- Basic key-value operations: Set, Get, Delete
- Simple TCP server interface
- Concurrent client handling
- Data persistence to disk (data.json)
- Atomic writes using temporary files to prevent data corruption

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

- `SET <key> <value>` - Store a key-value pair
- `GET <key>` - Retrieve a value by key
- `DEL <key>` - Delete a key-value pair
