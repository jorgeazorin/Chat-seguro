//////
//MAIN
//////

package main

import (
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type BD struct {
	username string
	password string
	adress   string
	database string
}

/*
func main() {
	var test bool
	var bd BD

	//BD
	bd.username = "sds"
	bd.password = "sds"
	bd.adress = ""
	bd.database = "sds"

	/////////
	//USUARIO
	/////////

	//Prueba insertar usuario
	var uu Usuario
	uu.nombre = "alex"
	uu.clavepubrsa = "clavepubrsa"
	uu.claveprivrsa = "claveprivrsa"
	uu.claveusuario = "clavecifrada"
	//bd.insertUsuarioBD(uu)

	//Prueba Modificar Usuario
	var u Usuario
	u.id = 15
	u.clavepubrsa = "clave15pubrsa"
	u.claveprivrsa = "clave15privrsa"
	u.claveusuario = "clave15cifrada"
	test = bd.modificarUsuarioBD(u)
	fmt.Println("Mira modificar usuario:", test)

	//Probar obtener nombre según id
	nombreusuario := bd.getNombreUsuario(1)
	fmt.Println("Mira el nombre del usuario:", nombreusuario)

	//Probar obtener usuario según id
	usuario := bd.getUsuario(1)
	fmt.Println("Mira el usuario:", usuario.id, usuario.nombre, usuario.clavepubrsa, usuario.claveprivrsa, usuario.claveusuario)

	//Prueba comprobar usuario
	var miusuario Usuario
	miusuario, test = bd.comprobarUsuarioBD("Maria", "clave15cifrada")
	fmt.Println("Mira comprobando usuario:", test, " tiene:", miusuario.clavepubrsa, miusuario.id)
	fmt.Println("-")

	//////////
	//MENSAJES
	//////////

	//Prueba guardar mensaje
	var m Mensaje
	m.texto = "Hola que tal?? :)"
	m.idchat = 1
	m.idemisor = 1
	m.idclave = 1
	//	test = bd.guardarMensajeBD(m)
	fmt.Println("Mira guardar mensaje:", test)

	//Prueba crear nueva clavesmensajes
	//id := bd.CrearNuevaClaveParaMensajesBD()
	//fmt.Println("Mira id clavesmensajes creado:", id)
	//fmt.Println("-")

	//Prueba insertar clave de un usuario para x mensajes
	//test = bd.GuardarClaveUsuarioMensajesBD(1, "claveusuario1", 1)
	//fmt.Println("Mira guardar clave usuario de x mensaje:", test)
	fmt.Println("-")

	//////
	//CHAT
	//////

	//Prueba crear chat
	usuarios := make([]int, 0, 1)
	usuarios = append(usuarios, 1)
	usuarios = append(usuarios, 2)
	usuarios = append(usuarios, 3)
	//test = bd.crearChatBD(usuarios, "")
	//fmt.Println("Mira crear chat:", test)

	//Prueba modificar chat
	var c Chat
	c.id = 5
	c.nombre = "grupo molon :)"
	test = bd.modificarChatBD(c)
	fmt.Println("Mira modificar chat:", test)

	//Prueba añadir usuarios a un char
	nuevosusuarios := make([]int, 0, 1)
	nuevosusuarios = append(nuevosusuarios, 4)
	nuevosusuarios = append(nuevosusuarios, 5)
	//test = bd.addUsuariosChatBD(7, nuevosusuarios)
	//fmt.Println("Mira añadir nuevos usuarios a chat:", test)

	//Prueba eliminar usuarios de un char
	usuariosexpulsados := make([]int, 0, 1)
	usuariosexpulsados = append(usuariosexpulsados, 4)
	usuariosexpulsados = append(usuariosexpulsados, 5)
	//test = bd.removeUsuariosChatBD(7, usuariosexpulsados)
	//fmt.Println("Mira eliminar usuarios a chat:", test)

	//Obtener mensajes de un usuario
	chats := make([]Chat, 0, 1)
	chats = bd.obtenerChatsUsuarioBD(15)
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
	fmt.Println("-")
}
*/
