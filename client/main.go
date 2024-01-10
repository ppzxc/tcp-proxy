package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	var host string
	var port int

	flag.StringVar(&host, "host", "localhost", "--host localhost")
	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.Parse()

	ctx, cancelFunc := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
				if err != nil {
					panic(err)
				}
				writeLoop(ctx, conn)
			}
		}
	}()
	sig := <-signals
	log.Println(sig)
	cancelFunc()
}

func writeLoop(ctx context.Context, conn net.Conn) {
	go func() {
		reader := bufio.NewReader(conn)
		for {
			read, err := reader.ReadString('\n')
			if err != nil {
				log.Println("close", err)
				return
			}
			log.Printf("INBOUND [%s]\n", strings.ReplaceAll(read, "\n", ""))
		}
	}()

	closeTicker := time.NewTicker(11 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			writeString := "Hello World!\n"
			log.Printf("OUTBOUND [%s]\n", strings.ReplaceAll(writeString, "\n", ""))
			_, err := conn.Write([]byte(writeString))
			if err != nil {
				log.Println(err)
				return
			}
		case <-closeTicker.C:
			log.Println("close tick")
			err := conn.Close()
			if err != nil {
				log.Println(err)
			}
			break
		case <-ctx.Done():
			log.Println("interrupted")
			break
		}
	}
}

