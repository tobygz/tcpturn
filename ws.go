package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("ws recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func wsServ() {
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HomeTemplate.Execute(w, "ws://"+r.Host+"/echo")
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

func wscli() {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1:%s", *port), Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	i := 0
	for {
		i++
		str := fmt.Sprintf("%d", i)
		err := c.WriteMessage(websocket.TextMessage, []byte(str))
		if err != nil {
			log.Println("write:", err)
			return
		}
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
		time.Sleep(time.Second)
	}
}
