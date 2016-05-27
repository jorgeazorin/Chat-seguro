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

/////////////////////////////
//Envio de mensajes a cliente
/////////////////////////////
func echoHandler(ws *websocket.Conn) {
	wbSocket = ws

	for {
		datos := leerDatosWS(ws)
		fmt.Println(datos)

		///////
		//Login
		///////
		if datos == "login" {
			datos := leerDatosWS(ws)
			var usuario Usuario
			json.Unmarshal([]byte(datos), &usuario)
			loginweb(usuario.Nombre, usuario.Claveenclaro)

			obtenermensajesAdmin()

			b, _ := json.Marshal(ClientUsuario)
			mensaje := MensajeSocket{Mensaje: "DatosUsuario", Datos: []string{string(b)}}
			escribirSocketCliente(mensaje)
		}

		//////////
		//Registro
		//////////
		if datos == "registro" {
			datos := leerDatosWS(ws)
			var usuario Usuario
			json.Unmarshal([]byte(datos), &usuario)
			test := registrarUsuario(usuario)

			if test == false {
				websocket.Message.Send(ws, "registronook")
			}

			getUsuarios()

			b, _ := json.Marshal(ClientUsuario)
			mensaje := MensajeSocket{Mensaje: "DatosUsuario", Datos: []string{string(b)}}
			escribirSocketCliente(mensaje)
		}

		//////////////
		//Get Usuarios
		//////////////
		if datos == "getusuarios" {
			usuarios := getUsuarios()

			b, _ := json.Marshal(usuarios)
			mensaje := MensajeSocket{Mensaje: "getusuariosok", Datos: []string{string(b)}}
			escribirSocketCliente(mensaje)
		}

		////////////////
		//Obtener chats
		////////////////
		if datos == "chats" {
			chats := obtenerChats()
			for i := 0; i < len(chats); i++ {
				chats[i].MensajesDatos = obtenerMensajesChat(chats[i].Chat.Id)
			}

			b, _ := json.Marshal(chats)
			mensaje := MensajeSocket{Mensaje: "chats", Datos: []string{string(b)}}
			escribirSocketCliente(mensaje)
		}

		////////////////
		//Enviar mensaje
		////////////////
		if datos == "enviarmensaje" {
			datos := leerDatosWS(ws)
			var mensaje MensajeSocket
			json.Unmarshal([]byte(datos), &mensaje)
			mensaje.Mensajechat = []byte(mensaje.Mensaje)

			fmt.Println("Mira:", mensaje)

			test := enviarMensaje(mensaje)
			if test == false {
				mensaje := MensajeSocket{Mensaje: "Error al enviar el mensaje."}
				websocket.Message.Send(ws, mensaje)
			} else {
				mensaje := MensajeSocket{Mensaje: "mensajeenviado:"}
				escribirSocketCliente(mensaje)
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
			mensaje = MensajeSocket{Mensaje: "mensajesleidos"}
			escribirSocketCliente(mensaje)
		}

		////////////////
		//Modificar chat
		////////////////
		if datos == "editarchat" {
			datos := leerDatosWS(ws)
			var chat Chat
			json.Unmarshal([]byte(datos), &chat)

			editarChat(chat)
			mensaje := MensajeSocket{Mensaje: "chatcambiadook"}
			escribirSocketCliente(mensaje)
		}

		////////////
		//Crear chat
		////////////
		if datos == "crearchat" {
			datos := leerDatosWS(ws)

			crearChat(datos)
			mensaje := MensajeSocket{Mensaje: "chatcreadook"}
			escribirSocketCliente(mensaje)
		}

		///////////////////
		//Modificar usuario
		///////////////////
		if datos == "editarusuario" {
			datos := leerDatosWS(ws)
			var usuario Usuario
			json.Unmarshal([]byte(datos), &usuario)

			editarUsuario(usuario)
			mensaje := MensajeSocket{Mensaje: "usuariocambiaok", Datos: []string{usuario.Nombre, usuario.Estado}}
			escribirSocketCliente(mensaje)
		}

	}

}
