//////
//MAIN
//////

package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var username = "sds"
var password = "sds"
var adress = ""
var database = "sds"

func main() {
	var test bool

	/////////
	//USUARIO
	/////////

	//Prueba insertar usuario
	var uu Usuario
	uu.nombre = "alex"
	uu.clavepubrsa = "clavepubrsa"
	uu.claveprivrsa = "claveprivrsa"
	uu.claveusuario = "clavecifrada"
	//insertUsuarioBD(uu)

	//Prueba Modificar Usuario
	var u Usuario
	u.id = 15
	u.clavepubrsa = "clave15pubrsa"
	u.claveprivrsa = "clave15privrsa"
	u.claveusuario = "clave15cifrada"
	test = modificarUsuarioBD(u)
	fmt.Println("Mira modificar usuario:", test)

	//Probar obtener nombre según id
	nombreusuario := getNombreUsuario(1)
	fmt.Println("Mira el nombre del usuario:", nombreusuario)

	//Probar obtener usuario según id
	usuario := getUsuario(1)
	fmt.Println("Mira el usuario:", usuario.id, usuario.nombre, usuario.clavepubrsa, usuario.claveprivrsa, usuario.claveusuario)

	//Prueba comprobar usuario
	test = comprobarUsuarioBD("pepe", "clave1cifrada")
	fmt.Println("Mira comprobando usuario:", test)
	fmt.Println("-")

	//////////
	//MENSAJES
	//////////

	//Prueba guardar mensaje
	//test = guardarMensajeBD("Hola que tal?? :)", 5, 1, 1)
	//fmt.Println("Mira guardar mensaje:", test)
	//fmt.Println("-")

	//Prueba obtener mensajes
	//test = obtenerMensajeBD("Hola que tal?? :)", 5, 1, 1)
	//fmt.Println("Mira guardar mensaje:", test)
	//fmt.Println("-")

	//Obtener mensajes de un usuario
	chats := make([]Chat, 0, 1)
	chats = obtenerChatsUsuarioBD(15)
	fmt.Println("-")
	fmt.Println("Mira mensajes usuario 15 Maria")

	for i := 0; i < len(chats); i++ {
		fmt.Println("Mira mi chat id", chats[i].id, "es", chats[i].nombre)

		for j := 0; j < len(chats[i].mensajes); j++ {
			if chats[i].mensajes[j].idemisor != 15 {
				fmt.Println("De", chats[i].mensajes[j].nombreemisor, "-> ", chats[i].mensajes[j].id, ": '", chats[i].mensajes[j].texto, "' / leido:", chats[i].mensajes[j].leido)
			} else {
				fmt.Println("De", chats[i].mensajes[j].nombreemisor, "-> ", chats[i].mensajes[j].id, ": '", chats[i].mensajes[j].texto, "' / leido: es un mensaje mio")
			}
		}
	}

}
