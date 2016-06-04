package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

var puerto = ""

func IniciarServidorWeb() {
	go http.Handle("/echo", websocket.Handler(echoHandler))

	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/index.js", js)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	puerto, _ := reader.ReadString('\n')

	var err = http.ListenAndServeTLS(":"+puerto, "cert.pem", "key.pem", nil)
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
		return "-1"
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
		//	fmt.Println(datos)
		if datos == "-1" {

			break
		}

		///////
		//Login
		///////
		var datos1 = strings.Split(datos, "@/@")
		fmt.Println(datos1)
		if datos1[0] == "login" {
			var usuario Usuario
			json.Unmarshal([]byte(datos1[1]), &usuario)
			loginweb(usuario.Nombre, usuario.Claveenclaro)
			b, _ := json.Marshal(ClientUsuario)
			mensaje := MensajeSocket{Mensaje: "DatosUsuario", Datos: []string{string(b)}}
			escribirSocketCliente(mensaje)
			datos1[0] = "chats"
		}

		//////////
		//Registro
		//////////
		if datos1[0] == "registro" {
			var usuario Usuario
			json.Unmarshal([]byte(datos1[1]), &usuario)
			test := registrarUsuario(usuario)

			if test == false {
				websocket.Message.Send(ws, "registronook")
			}
			loginweb(usuario.Nombre, usuario.Claveenclaro)
			b, _ := json.Marshal(ClientUsuario)
			mensaje := MensajeSocket{Mensaje: "DatosUsuario", Datos: []string{string(b)}}
			escribirSocketCliente(mensaje)
		}

		//////////////
		//Get Usuarios
		//////////////
		if datos1[0] == "getusuarios" {
			usuarios := getUsuarios()

			b, _ := json.Marshal(usuarios)
			mensaje := MensajeSocket{Mensaje: "getusuariosok", Datos: []string{string(b)}}
			escribirSocketCliente(mensaje)
		}

		////////////////
		//Obtener chats
		////////////////
		if datos1[0] == "chats" {
			//fmt.Println(ClientUsuario)
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
		if datos1[0] == "enviarmensaje" {
			//datos := leerDatosWS(ws)
			var mensaje MensajeSocket
			json.Unmarshal([]byte(datos1[1]), &mensaje)
			mensaje.Mensajechat = []byte(mensaje.Mensaje)
			mensaje.Idfrom = ClientUsuario.Id
			//fmt.Println("Mira:", mensaje)

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
		if datos1[0] == "addusuariochat" {
			//datos := leerDatosWS(ws)
			var mensaje MensajeSocket
			json.Unmarshal([]byte(datos1[1]), &mensaje)

			agregarUsuariosChat(mensaje.Chat, []string{mensaje.Mensaje})
		}

		////////////////////////
		//Remove usuario de chat
		////////////////////////
		if datos1[0] == "removeusuariochat" {
			//datos := leerDatosWS(ws)
			var mensaje MensajeSocket
			json.Unmarshal([]byte(datos1[1]), &mensaje)

			eliminarUsuariosChat(mensaje.Chat, []string{mensaje.Mensaje})
		}

		///////////////////
		//Marcar chat leido
		///////////////////
		if datos1[0] == "leidos" {
			//datos := leerDatosWS(ws)
			var mensaje MensajeSocket
			json.Unmarshal([]byte(datos1[1]), &mensaje)

			MarcarChatComoLeido(mensaje.Chat)
			mensaje = MensajeSocket{Mensaje: "mensajesleidos"}
			escribirSocketCliente(mensaje)
		}

		////////////////
		//Modificar chat
		////////////////
		if datos1[0] == "editarchat" {
			//datos := leerDatosWS(ws)
			var chat Chat
			json.Unmarshal([]byte(datos1[1]), &chat)

			editarChat(chat)
			mensaje := MensajeSocket{Mensaje: "chatcambiadook"}
			escribirSocketCliente(mensaje)
		}

		////////////
		//Crear chat
		////////////
		if datos1[0] == "crearchat" {
			//datos := leerDatosWS(ws)

			crearChat(datos1[1])
			mensaje := MensajeSocket{Mensaje: "chatcreadook"}
			escribirSocketCliente(mensaje)
		}

		///////////////////
		//Modificar usuario
		///////////////////
		if datos1[0] == "editarusuario" {
			//datos := leerDatosWS(ws)
			var usuario Usuario
			json.Unmarshal([]byte(datos1[1]), &usuario)

			editarUsuario(usuario)
			mensaje := MensajeSocket{Mensaje: "usuariocambiaok", Datos: []string{usuario.Nombre, usuario.Estado}}
			escribirSocketCliente(mensaje)
		}

	}

}
