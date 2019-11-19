package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func sender(conn net.Conn, timeout time.Duration) {
	var inStr string
	var err error
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		if inStr, err = stdinReader.ReadString('\n'); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("sending %#v ...\n", inStr)
		if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			log.Fatal(err)
		}
		if _, err = io.WriteString(conn, inStr); err != nil {
			log.Fatal(err)
		}
		//fmt.Println("... sent")
	}
}

func receiver(conn net.Conn) {
	var recvStr string
	var err error
	var r = bufio.NewReader(conn)
	for {
		//fmt.Println("receiving...")
		if recvStr, err = r.ReadString('\n'); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Print(recvStr)
	}
}

func tcpTalk(host, port string, timeout time.Duration) error {
	addr := net.JoinHostPort(host, port)
	fmt.Println("connecting to", addr)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	go receiver(conn)
	sender(conn, timeout)
	return nil
}
