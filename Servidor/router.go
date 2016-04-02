package main

import (
	"strconv"
)

func (conexion *Conexion) ProcesarMensajeSocket(mensaje MensajeSocket) {

	var bd BD //Para las operaciones con la BD

	if mensaje.Funcion == "login" {

		//Rellenamos el usuario de la conexión con el login
		conexion.usuario.login(mensaje.From)

		//Enviamos un mensaje a las demás conexiones mostrando que está diponible el usuario
		//Preparamos el mensaje que vamos a enviar
		mesj := MensajeSocket{From: conexion.usuario.nombre, MensajeSocket: "Usuario online"}

		//recorremos el vector de conexiones
		for i := 0; i < len(conexion.conexiones.conexiones); i++ {
			//si la conexión es distinta de nuestro socket guardamos los datos del usuario
			if conexion.conexiones.conexiones[i].conexion != conexion.conexion {
				//enviamos un mensaje al resto de usuarios conectados
				conexion.conexiones.conexiones[i].EnviarMensajeSocketSocket(mesj)
			}

		}
	}

	if mensaje.Funcion == "enviar" {
		//Guardamos los mensajes en la BD
		var m Mensaje
		m.texto = mensaje.MensajeSocket
		m.idchat = 1
		m.idemisor = conexion.usuario.id
		m.idclave = 1
		bd.guardarMensajeBD(m)

		//Obtenemos los usuarios que pertenecen en el chat
		idChat, _ := strconv.Atoi(mensaje.Datos[0])
		usuarios := bd.getUsuariosChatBD(idChat)

		//Enviamos el mensaje a todos los que tienen el socket abierto que esten el chat
		for i := 0; i < len(conexion.conexiones.conexiones); i++ {
			for j := 0; j < len(usuarios); j++ {
				if conexion.conexiones.conexiones[i].usuario.id == usuarios[j] {
					conexion.conexiones.conexiones[i].EnviarMensajeSocketSocket(mensaje)
				}
			}

		}
	}

}
