# KVLite

A simple in-memory key-value store with a TCP server interface, implemented in Go.

## Features

- Thread-safe operations with read-write mutex
- Basic key-value operations: Set, Get, Delete
- Simple TCP server interface
- Concurrent client handling

## Usage

Start the server:
```bash
go run main.go
```

Connect using netcat or telnet:
```
$ nc localhost 8080
> SET foo bar
OK
> GET foo
bar
> DEL foo
OK
> GET foo
(not found)
```

## Commands

- `SET <key> <value>` - Store a key-value pair
- `GET <key>` - Retrieve a value by key
- `DEL <key>` - Delete a key-value pair
