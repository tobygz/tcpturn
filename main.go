package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/gorilla/websocket"
)

var exetype = flag.String("exetype", "ts", "")
var port = flag.String("port", "8025", "")

// for turn server
var htarget = flag.String("htarget", "127.0.0.1:8020", "")
var wtarget = flag.String("wtarget", "127.0.0.1:8021", "")

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
func hsServ() {
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}

func main() {
	flag.Parse()

	if *exetype == "ws" {
		wsServ()
	} else if *exetype == "ts" {
		tcpServ()
	} else if *exetype == "tc" {
		tcpCli()
	} else if *exetype == "wc" {
		wscli()
	} else if *exetype == "hs" {
		hsServ()
	} else if *exetype == "ss" {
		StartTurnServ()
	}
	fmt.Println("unknown exetype:", *exetype)
}
