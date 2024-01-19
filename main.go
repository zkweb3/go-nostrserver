package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	balance    = 100
)

var (
	addr     = flag.String("addr", "localhost:443", "https service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	reply = []string{
		"000003dcad81050de3a507b37c929a68cad06711b3993d67ccc43da766e1593c",
		"00000147028ba43a4d73d2b86e972248739a50bd7983ac0f9498e5a1b30c2bdf",
		"000001ed5df72d80e42c6be870c33fe1e2a8274a41a91467f5e4f25a401aaaef",
		"000002098eda9a22879715ff9dd9c1fd8584fe93a88e35ea0090e2334ce417df",
		"00000142bb9b3a25ba093afe5c5cad46dfc1965f34735f1b841f461f0d612595",
	}
)

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

func WebsocketHandle(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer c.Close()

	done := make(chan struct{})
	go ping(c, done)
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		for _, r := range reply {
			s := fmt.Sprintf(`{"eventId":"%s"}`, r)
			c.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.WriteMessage(websocket.TextMessage, []byte(s))
			if err != nil {
				log.Println(err)
				close(done)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func PostEventHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-type", "application/json")
		fmt.Fprintf(w, `{"balance":"%v"}`, balance)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	log.Printf("listen at %s", *addr)
	http.HandleFunc("/", WebsocketHandle)
	http.HandleFunc("/postEvent", PostEventHandle)
	http.HandleFunc("/getBalance", GetBalance)
	log.Fatalln(http.ListenAndServeTLS(*addr, "ca/cacert.pem", "ca/cakey.pem", nil))
}
