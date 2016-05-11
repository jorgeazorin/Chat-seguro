package main

import (
	//"crypto/sha256"

	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"net/http"
)

type T struct {
	Msg   string
	Count int
}

// receive JSON type T
var data T
var wbSocket *websocket.Conn

func HelloServer(w http.ResponseWriter, req *http.Request) {
	webpage, err := ioutil.ReadFile("web/index.html")
	if err != nil {
		panic(err)

	}
	io.WriteString(w, string(webpage))
}

func js(w http.ResponseWriter, req *http.Request) {
	webpage, err := ioutil.ReadFile("web/index.js")
	if err != nil {
		panic(err)

	}
	io.WriteString(w, string(webpage))
}

func IniciarServidorWeb() {
	go http.Handle("/echo", websocket.Handler(echoHandler))

	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/index.js", js)

	var err = http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	if err != nil {
		panic(err)

	}

}

func escribitWebSocket(ws *websocket.Conn) {

	var message = "hello"
	websocket.Message.Send(ws, message)
}

func echoHandler(ws *websocket.Conn) {
	wbSocket = ws
	receivedtext := make([]byte, 100)

	n, err := ws.Read(receivedtext)

	if err != nil {
		fmt.Printf("Received: %d bytes\n", n)
	}

	s := string(receivedtext[:n])
	fmt.Printf("Received: %d bytes: %s\n", n, s)
	var message = "hello"
	websocket.Message.Send(ws, message)

	for {
		receivedtext := make([]byte, 100)
		n, err := ws.Read(receivedtext)
		s := string(receivedtext[:n])
		fmt.Printf("Received: %d bytes: %s\n", n, s)
		if err != nil {
			break
		}
		//	websocket.Message.Send(ws, message)
	}

}
