package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

type saveOutput struct {
	reader io.Reader
	conn   net.Conn
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	fmt.Println("SERVER -> CLIENT:", string(p))
	if so.conn != nil {
		n, err = so.conn.Write(p)
	}
	return n, err
}

type proxyConn struct {
	conn   net.Conn
	cmd    *exec.Cmd
	writer io.Writer
}

func main() {
	reader, writer := io.Pipe()
	so := &saveOutput{reader: reader}
	se := &saveOutput{reader: reader}
	cmd := exec.Command("gopls")
	cmd.Dir = "/Users/pjlast/workspace/lsproxy"
	cmd.Stdin = reader
	cmd.Stdout = so
	cmd.Stderr = so
	go cmd.Run()

	l, err := net.Listen("tcp4", ":1337")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		pc := proxyConn{conn: c, writer: writer}
		so.conn = c
		se.conn = c
		go pc.handleConnection()
	}
}

func (pc *proxyConn) handleConnection() {
	fmt.Printf("Serving %s\n", pc.conn.RemoteAddr().String())
	bufferedReader := bufio.NewReader(pc.conn)
	for {
		header, err := bufferedReader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			println("END OF FILE")
			break
		}
		numToRead, _ := strconv.Atoi(strings.TrimRight(string(bytes.Split(header, []byte(" "))[1]), "\r\n"))
		msg := make([]byte, numToRead+2)
		io.ReadFull(bufferedReader, msg)

		msg = slices.Concat(header, msg)
		fmt.Println("CLIENT -> SERVER", string(msg))
		pc.writer.Write(msg)
	}
}
