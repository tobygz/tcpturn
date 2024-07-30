package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func StartTurnServ() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		fmt.Println("ListenRemoteClient failed:", err)
		panic("")
	}
	for {
		clientConn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept err:", err)
			return
		}
		go connectPipe(clientConn)
	}
}

func connectPipe(clientConn net.Conn) {
	iswebsocket := false
	buf := make([]byte, 4096)
	n := 0
	var err error
	for {
		//determine use http or websocket
		reader := bufio.NewReader(clientConn)
		n, err = reader.Read(buf)
		if err != nil {
			fmt.Println("conn failed:", err)
			return
		}
		if strings.Contains(string(buf[:n]), "\r\n") {
			log.Println("check header data:", string(buf[:n]))
			if strings.Contains(string(buf[:n]), "Sec-WebSocket-Key") ||
				strings.Contains(string(buf[:n]), "Sec-WebSocket-Accept") {
				//is websocket
				iswebsocket = true
			} else {
				//is http
				iswebsocket = false
			}
			break
		}
	}

	target := *htarget
	if iswebsocket {
		target = *wtarget
	}
	serverConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Println("net.Dial err:", err)
		_ = clientConn.Close()
		return
	}
	pipe(clientConn, serverConn, buf[:n])
}

func pipe(src net.Conn, dest net.Conn, buf []byte) {
	errChan := make(chan error, 1)
	onClose := func(err error) {
		_ = dest.Close()
		_ = src.Close()
	}
	go func() {
		_, err := io.Copy(src, dest)
		errChan <- err
		onClose(err)
	}()
	n, err := dest.Write(buf)
	if err != nil {
		log.Println("dest.Write failed:", err)
		return
	}
	log.Println(" send to dest size:", n, " buf len:", len(buf))
	go func() {
		_, err := io.Copy(dest, src)
		errChan <- err
		onClose(err)
	}()
	<-errChan
}
