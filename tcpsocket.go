package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func tcpServ() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Println("Error listening:", err)
		os.Exit(1)
	}

	// 循环接收并打印客户端连接的数据
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		defer conn.Close()

		go func(_cn net.Conn) {
			// 接收并打印数据
			for {
				buf := make([]byte, 1024)
				n, err := _cn.Read(buf)
				if err != nil {
					log.Println("Error reading data:", err)
					return
				}

				_cn.Write(buf[:n])
				log.Println("Received from client:", string(buf[:n]))
				time.Sleep(time.Second * 1)
			}
		}(conn)
	}
}

func tcpCli() {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", *port))
	if err != nil {
		log.Println("Error connecting:", err)
		os.Exit(1)
	}
	i := 100
	buf := make([]byte, 1024)
	for {
		i++
		conn.Write([]byte(strconv.Itoa(i)))
		log.Println("send i:", i)

		lenv, ee := conn.Read(buf)
		if ee != nil || lenv <= 0 {
			panic("conn closed")
			break
		}
		log.Println("recv i:", buf[:lenv])
		time.Sleep(time.Second)
	}
}
