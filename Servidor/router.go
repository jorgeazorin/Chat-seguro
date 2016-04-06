package main

import (
	//"log"
	"fmt"
	"strconv"
)

func (conexion *Conexion) ProcesarMensajeSocket(mensaje MensajeSocket) {

	//Para las operaciones con la BD
	var bd BD
	bd.username = "sds"
	bd.password = "sds"
	bd.adress = ""
	bd.database = "sds"

	if mensaje.Funcion == "login" {

		//Rellenamos el usuario de la conexión con el login
		test := conexion.usuario.login(mensaje.From, mensaje.Password)

		if test == false {
			mesj := MensajeSocket{From: conexion.usuario.nombre, MensajeSocket: "Login incorrecto"}
			conexion.EnviarMensajeSocketSocket(mesj)
			return
		}

		//Añadimos la conexion al map de conexiones bloqueando la memoria compartida
		conexiones[conexion.usuario.id] = conexion

		//Enviamos un mensaje a las demás conexiones mostrando que está diponible el usuario
		//Preparamos el mensaje que vamos a enviar
		mesj := MensajeSocket{From: conexion.usuario.nombre, MensajeSocket: "Logeado correctamente"}

		//Enviamos un mensaje
		conexiones[conexion.usuario.id].EnviarMensajeSocketSocket(mesj)

	}

	if mensaje.Funcion == "enviar" {
		//Guardamos los mensajes en la BD
		var m Mensaje
		m.texto = mensaje.MensajeSocket
		m.idchat = 1
		m.idemisor = conexion.usuario.id
		m.idclave = 1
		//bd.guardarMensajeBD(m)

		//Obtenemos los usuarios que pertenecen en el chat
		idChat, _ := strconv.Atoi(mensaje.Datos[0])
		idusuarios := bd.getUsuariosChatBD(idChat)

		//Enviamos el mensaje a todos los usuarios de ese chat con el socket abierto (incluido el emisor)
		for i := 0; i < len(idusuarios); i++ {
			fmt.Println("Mira que usuario:", idusuarios[i])
			conexiones[idusuarios[i]].EnviarMensajeSocketSocket(mensaje)
		}

	}

}
