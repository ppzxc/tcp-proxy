# tcp forward proxy using envoy

## envoy

- listen 10000 port 
- proxy to echo.example.com:9000

## client

- send "hello world" per 1s
- reconnect per 10s

## echo

- logging and echo

## usage

- require docker

```
make
```
