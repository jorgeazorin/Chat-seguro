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
	uu.Nombre = "alex"
	uu.Clavepubrsa = []byte("clavepubrsa")
	uu.Claveprivrsa = []byte("claveprivrsa")
	uu.Clavelogin = []byte("clavecifrada")
	//bd.insertUsuarioBD(uu)

	//Prueba Modificar Usuario
	var u Usuario
	u.Id = 15
	u.Nombre = "Usuario15"
	u.Clavepubrsa = []byte("clave15pubrsa")
	u.Claveprivrsa = []byte("clave15privrsa")
	u.Salt = []byte("clave15cifrada")
	u.Clavelogin = []byte("clave15cifrada")
	test = bd.modificarUsuarioBD(u)
	fmt.Println("Mira modificar usuario:", test)

	//Probar obtener nombre según id
	nombreusuario, _ := bd.getNombreUsuario(1)
	fmt.Println("Mira el nombre del usuario:", nombreusuario)

	//Probar obtener pub según id
	clavepub, _ := bd.getClavePubUsuario(1)
	fmt.Println("Mira clave publica del usuario:", clavepub)

	//Probar obtener usuario según id
	usuario, _ := bd.getUsuarioById(1)
	fmt.Println("Mira el usuario:", usuario.Id, usuario.Nombre, usuario.Clavepubrsa, usuario.Claveprivrsa, usuario.Clavelogin)

	//Probar obtener usuario según nombre
	usuario, _ = bd.getUsuarioByNombreBD("Pepe")
	fmt.Println("Mira el usuario:", usuario.Id, usuario.Nombre, usuario.Clavepubrsa, usuario.Claveprivrsa, usuario.Clavelogin)

	//get usuarios de un chat
	usuarios, _ := bd.getUsuariosChatBD(1)
	fmt.Println("Mira:", usuarios)

	claves, _ := bd.getClavesMensajes(1)
	fmt.Println("Mira:", claves)

	//Prueba comprobar usuario
	//var miusuario Usuario
	//miusuario, test = bd.loginUsuarioBD("Maria", "clave15cifrada")
	//fmt.Println("Mira comprobando usuario:", test, " tiene:", miusuario.clavepubrsa, miusuario.id)
	//fmt.Println("-")

	//////////
	//MENSAJES
	//////////

	//Prueba guardar mensaje
	var m Mensaje
	m.Texto = "Hola que tal?? :)"
	m.Chat = 1
	m.Emisor = 1
	m.Clave = 1
	//test = bd.guardarMensajeBD(m)
	//fmt.Println("Mira guardar mensaje:", test)

	//Prueba crear nueva clavesmensajes
	//id, _ := bd.CrearNuevaClaveMensajesBD()
	//fmt.Println("Mira id clavesmensajes creado:", id)
	//fmt.Println("-")

	//Prueba insertar clave de un usuario para x mensajes
	var clavesusuario Clavesusuario
	clavesusuario.Idusuario = 1
	clavesusuario.Idclavesmensajes = 4
	clavesusuario.Clavemensajes = []byte("claveusuario1")
	//test = bd.GuardarClaveUsuarioMensajesBD(clavesusuario)
	//fmt.Println("Mira guardar clave usuario de x mensaje:", test)
	//fmt.Println("-")

	mismensajes, _ := bd.getMensajesChatBD(1, 1)
	fmt.Println("mira los mensajes:", mismensajes)

	clavemen, _ := bd.getClaveMensaje(11, 31)
	fmt.Println("mira clave mensaje 11 usuario 31:", clavemen)

	men, _ := bd.getMensajeBD(11)
	fmt.Println("mira  mensaje 11:", men)

	clavlast, _ := bd.getLastKeyMensaje(1, 1)
	fmt.Println("mira ultima clave chat 1 usu 1:", string(clavlast))

	test = bd.marcarLeidoPorUsuarioBD(23, 15)
	fmt.Println("marcado como leido mensaje 23 usu 15", test)
	test = bd.marcarLeidoPorUsuarioBD(23, 31)
	fmt.Println("marcado como leido mensaje 23 usu 31", test)

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
	chats = bd.getChatsUsuarioBD(15)
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
	fmt.Println("-")/
}
*/
