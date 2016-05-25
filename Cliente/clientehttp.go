package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type T struct {
	Msg   string
	Count int
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//Generar cadena aleatoria
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
	//var message = "hello"
	//websocket.Message.Send(ws, message)
}

func leerDatosWS(ws *websocket.Conn) string {
	receivedtext := make([]byte, 100)
	n, err := ws.Read(receivedtext)
	if err != nil {
		fmt.Printf("Error obteniendo datos:", err)
	}
	s := string(receivedtext[:n])

	return s
}

func echoHandler(ws *websocket.Conn) {
	wbSocket = ws

	//Esto para que es Jorge?
	var message = "hello"
	websocket.Message.Send(ws, message)

	for {
		datos := leerDatosWS(ws)

		///////
		//Login
		///////
		if datos == "login" {
			datos := leerDatosWS(ws)
			var usuario Usuario
			json.Unmarshal([]byte(datos), &usuario)
			loginweb(usuario.Nombre, usuario.Claveenclaro)
			obtenermensajesAdmin()
		}

		//////////
		//Registro
		//////////
		if datos == "registro" {
			datos := leerDatosWS(ws)

			var usuario Usuario
			json.Unmarshal([]byte(datos), &usuario)
			test := registrarUsuario(usuario)

			if test == true {
				websocket.Message.Send(ws, "registrook")
			} else {
				websocket.Message.Send(ws, "registronook")
			}
		}

		////////////////
		//Obtener chats
		////////////////
		if datos == "chats" {
			obtenerChats()
		}

		////////////////
		//Enviar mensaje
		////////////////
		if datos == "enviarmensaje" {
			datos := leerDatosWS(ws)
			var mensaje MensajeSocket
			json.Unmarshal([]byte(datos), &mensaje)
			mensaje.Mensajechat = []byte(mensaje.Mensaje)

			test := enviarMensaje(mensaje)
			if test == false {
				mensaje := MensajeSocket{Mensaje: "Error al enviar el mensaje."}
				websocket.Message.Send(ws, mensaje)
			}
		}

		////////////////////
		//Add usuario a chat
		////////////////////
		if datos == "addusuariochat" {
			datos := leerDatosWS(ws)
			var mensaje MensajeSocket
			json.Unmarshal([]byte(datos), &mensaje)

			agregarUsuariosChat(mensaje.Chat, []string{mensaje.Mensaje})
		}

		///////////////////
		//Marcar chat leido
		///////////////////
		if datos == "leidos" {
			datos := leerDatosWS(ws)
			var mensaje MensajeSocket
			json.Unmarshal([]byte(datos), &mensaje)

			MarcarChatComoLeido(mensaje.Chat)
		}

	}

}
