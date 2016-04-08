package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

//Struct de los mensajes que se envian por el socket
type MensajeSocket struct {
	From          string   `json:"From"`
	To            int      `json:"To"`
	Password      string   `json:"Password"`
	Funcion       string   `json:"Funcion"`
	Datos         []string `json:"Datos"`
	MensajeSocket string   `json:"MensajeSocket"`
}

func ProcesarMensajeSocket(mensaje MensajeSocket, conexion net.Conn, usuario *Usuario) {

	//Para las operaciones con la BD
	var bd BD
	bd.username = "root"
	bd.password = ""
	bd.adress = ""
	bd.database = "sds"

	if mensaje.Funcion == "login" {

		//Rellenamos el usuario de la conexión con el login
		test := usuario.login(mensaje.From, mensaje.Password)

		//Si login incorrecto se lo decimos al cliente
		if test == false {
			mesj := MensajeSocket{From: usuario.nombre, MensajeSocket: "Login incorrecto"}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Si es correcto, aAñadimos la conexion al map de conexiones
		conexiones[usuario.id] = conexion

		//Enviamos un mensaje de todo OK al usuario logeado
		mesj := MensajeSocket{From: usuario.nombre, MensajeSocket: "Logeado correctamente"}
		EnviarMensajeSocketSocket(conexion, mesj)

	}

	if mensaje.Funcion == "enviar" {

		//Guardamos los mensajes en la BD
		var m Mensaje
		m.texto = mensaje.MensajeSocket
		m.idchat = 1
		m.idemisor = usuario.id
		m.idclave = 1
		//bd.guardarMensajeBD(m)

		//Obtenemos los usuarios que pertenecen en el chat
		idChat, _ := strconv.Atoi(mensaje.Datos[0])
		idusuarios := bd.getUsuariosChatBD(idChat)

		//Enviamos el mensaje a todos los usuarios de ese chat (incluido el emisor)
		for i := 0; i < len(idusuarios); i++ {
			conexion, ok := conexiones[idusuarios[i]]
			if ok {
				EnviarMensajeSocketSocket(conexion, mensaje)
			}
		}

	}

	if mensaje.Funcion == "obtenermensajeschat" {

		//Obtenemos los mensajes de ese chat
		idChat, _ := strconv.Atoi(mensaje.Datos[0])

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(usuario.id, idChat)

		fmt.Println("Lo tiene:", permitido)

		if permitido == false {
			//Enviamos mensaje error
			mesj := MensajeSocket{From: usuario.nombre, MensajeSocket: "No perteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Obtenemos los mensajes
		mensajes := bd.getMensajesChatBD(idChat)

		for i := 0; i < len(mensajes); i++ {
			fmt.Println("::::", mensajes[i].id, mensajes[i].texto)
		}

		//Codificamos los mensajes en json
		b, _ := json.Marshal(mensajes)

		fmt.Println(b)

		fmt.Println(string(b))

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: usuario.nombre, MensajeSocket: string(b)}
		EnviarMensajeSocketSocket(conexion, mesj)

	}

}
