package main

import (
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

		if test == false {
			mesj := MensajeSocket{From: usuario.nombre, MensajeSocket: "Login incorrecto"}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Añadimos la conexion al map de conexiones bloqueando la memoria compartida
		conexiones[usuario.id] = conexion

		//Enviamos un mensaje a las demás conexiones mostrando que está diponible el usuario
		//Preparamos el mensaje que vamos a enviar
		mesj := MensajeSocket{From: usuario.nombre, MensajeSocket: "Logeado correctamente"}

		//Enviamos un mensaje
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

}
