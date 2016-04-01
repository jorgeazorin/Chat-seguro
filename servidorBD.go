/*
	Encarna Amorós Beneite, Jorge Azorín Martí
	Práctica SDS
*/

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"reflect"
	"strconv"
	//"strings"
	//"unsafe"
)

//Datos autentificación en BD
//var username = "root"
//var password = "ViadJid3"
//var adress = "51.255.44.18" //vps222360.ovh.net
//var database = "sds"
var username = "sds"
var password = "sds"
var adress = ""
var database = "sds"

/////////
//USUARIO
/////////

//Insertamos a un nuevo usuario en BD
func insertUsuarioBD(nombre string, clavepubrsa string, claveprivrsa string, claveusuariocifrada string) bool {

	//Conexión BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO usuario VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
		return false
	}

	//Insertamos
	_, err = stmtIns.Exec("DEFAULT", nombre, clavepubrsa, claveprivrsa, claveusuariocifrada)
	if err != nil {
		panic(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

//Obtenemos nombre de usuario según id usuario
func getNombreUsuario(id int) string {

	var nombreusuario string

	//Conexión BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return ""
	}
	defer db.Close()

	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT nombre FROM usuario WHERE id = " + strconv.Itoa(id))
	if err != nil {
		panic(err.Error())
		defer db.Close()
		return ""
	}

	for rows.Next() {
		err = rows.Scan(&nombreusuario)
		if err != nil {
			panic(err.Error())
			defer db.Close()
			return ""
		}
	}

	return nombreusuario
}

func modificarUsuarioBD(idusuario int, clavepubrsa string, claveusuariocifrada string) bool {

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("UPDATE usuario set clavepubrsa=?, claveusuario=? where id=?")
	if err != nil {
		panic(err.Error())
		return false
	}

	//Insertamos crear el chat
	_, err = stmtIns.Exec(clavepubrsa, claveusuariocifrada, idusuario)
	if err != nil {
		panic(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

// Comprobamos un usuario con su nombre y clave cifrada
func comprobarUsuarioBD(nombre string, claveusuario string) bool {

	var idusuario int
	var claveusuariobd string

	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Obtenemos el id del usuario
	rows, err := db.Query("SELECT id FROM usuario WHERE nombre = '" + nombre + "'")
	if err != nil {
		panic(err.Error())
		defer db.Close()
		return false
	}

	for rows.Next() {
		err = rows.Scan(&idusuario)
		if err != nil {
			panic(err.Error())
			defer db.Close()
			return false
		}
	}

	if idusuario == 0 {
		return false
	}

	//Obtenemos el la clave del usuario con id obtenido
	rows, err = db.Query("SELECT claveusuario FROM usuario WHERE id = " + strconv.Itoa(idusuario))
	if err != nil {
		panic(err.Error())
		defer db.Close()
		return false
	}

	for rows.Next() {
		err = rows.Scan(&claveusuariobd)
		if err != nil {
			panic(err.Error())
			defer db.Close()
			return false
		}
	}

	//Vemos si claves coinciden
	if claveusuario != claveusuariobd {
		return false
	}

	return true
}

//////
//CHAT
//////

//Creamos nuevo chat en BD
func crearChatBD(idusuarios []int, nombrechat string) bool {

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("INSERT INTO chat VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
		return false
	}

	//Insertamos crear el chat
	res, err := stmtIns.Exec("DEFAULT", nombrechat)
	if err != nil {
		panic(err.Error())
		return false
	}

	//Obtenemos id del chat creado
	idchat, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
		return false
	}
	println("Id del chat creado:", idchat)

	defer stmtIns.Close()

	//Insertamos usuarios a dicho chat
	for i := 0; i < len(idusuarios); i++ {
		//Preparamos insertar usuario al chat
		stmtIns, err := db.Prepare("INSERT INTO usuarioschat VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
			return false
		}

		//Insertamos usuario al chat
		_, err = stmtIns.Exec(idusuarios[i], idchat)
		if err != nil {
			panic(err.Error())
			return false
		}
	}

	defer stmtIns.Close()

	return true
}

func modificarChatBD(idchat int, nombre string) bool {

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("UPDATE chat set nombre=? where id=?")
	if err != nil {
		panic(err.Error())
		return false
	}

	//Insertamos crear el chat
	_, err = stmtIns.Exec(nombre, idchat)
	if err != nil {
		panic(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

//Añade una serie de usuarios a un chat
func addUsuariosChatBD(idchat int, nuevosusuarios []int) bool {

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Insertamos usuarios a dicho chat
	for i := 0; i < len(nuevosusuarios); i++ {
		//Preparamos insertar usuario al chat
		stmtIns, err := db.Prepare("INSERT INTO usuarioschat VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
			return false
		}

		//Insertamos usuario al chat
		_, err = stmtIns.Exec(nuevosusuarios[i], idchat)
		if err != nil {
			panic(err.Error())
			return false
		}

		defer stmtIns.Close()
	}

	return true
}

//Elimina una serie de usuarios a un chat
func removeUsuariosChatBD(idchat int, usuariosexpulsados []int) bool {

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Insertamos usuarios a dicho chat
	for i := 0; i < len(usuariosexpulsados); i++ {
		//Preparamos insertar usuario al chat
		stmtIns, err := db.Prepare("DELETE FROM usuarioschat where idusuario=? and idchat=?")
		if err != nil {
			panic(err.Error())
			return false
		}

		//Insertamos usuario al chat
		_, err = stmtIns.Exec(usuariosexpulsados[i], idchat)
		if err != nil {
			panic(err.Error())
			return false
		}

		defer stmtIns.Close()
	}

	return true
}

////////////////////////////
//MENSAJES Y CLAVES MENSAJES
////////////////////////////

//Guarda un mensaje para todos los receptores posibles del chat
func guardarMensajeBD(texto string, idchat int, idemisor int, idclave int) bool {

	var idreceptoraux = -1
	idreceptores := make([]int, 0, 1)

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Obtenemos los id de los usuarios que recibiráne el mensaje en este chat
	rows, err := db.Query("SELECT idusuario FROM usuarioschat WHERE idchat = " + strconv.Itoa(idchat) + " and idusuario !=" + strconv.Itoa(idemisor))
	if err != nil {
		panic(err.Error())
		defer db.Close()
		return false
	}

	//Guardamos id de los usuarios receptores en slice de ids
	for rows.Next() {
		err = rows.Scan(&idreceptoraux)

		if err != nil {
			panic(err.Error())
			defer db.Close()
			return false
		}
		idreceptores = append(idreceptores, idreceptoraux)
	}

	//Preparamos la creación del mensaje
	stmtIns, err := db.Prepare("INSERT INTO mensaje VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
		return false
	}

	//Insertamos el mensaje
	res, err := stmtIns.Exec("DEFAULT", texto, idemisor, idchat, idclave)
	if err != nil {
		panic(err.Error())
		return false
	}

	//Obtenemos id del mensaje creado
	idmensaje, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
		return false
	}
	println("Id del mensaje creado:", idmensaje)

	//Por cada receptor
	for i := 0; i < len(idreceptores); i++ {

		//Preparamos la insercion del receptor del mensaje
		stmtIns, err := db.Prepare("INSERT INTO receptoresmensaje VALUES(?, ?, ?)")
		if err != nil {
			panic(err.Error())
			return false
		}

		//Insertamos el receptor con el mensaje
		_, err = stmtIns.Exec(idmensaje, idreceptores[i], "DEFAULT")
		if err != nil {
			panic(err.Error())
			return false
		}
	}
	return true
}

//Guardamos la clave de un usuario para leer x mensajes
func GuardarClaveUsuarioMensajesBD(idclavesmensajes int, claveusuario string, idusuario int) bool {

	//Conexión BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO clavesusuario VALUES(?, ?, ?)")
	if err != nil {
		panic(err.Error())
		return false
	}

	//Insertamos
	_, err = stmtIns.Exec(idusuario, idclavesmensajes, claveusuario)
	if err != nil {
		panic(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

//Crear nuevo id para nuevo grupo de claves para los mensajes
func CrearNuevaClaveParaMensajesBD() int64 {

	//Conexión BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return 0
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO clavesmensajes VALUES(?)")
	if err != nil {
		panic(err.Error())
		return 0
	}

	//Insertamos
	res, err := stmtIns.Exec("DEFAULT")
	if err != nil {
		panic(err.Error())
		return 0
	}

	//Obtenemos id de lo creado
	idclavesmensajes, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
		return 0
	}

	defer stmtIns.Close()

	return idclavesmensajes
}

//Para guardar un mensaje con sus datos
type Mensaje struct {
	id           int
	texto        string
	idemisor     int
	nombreemisor string
	leido        bool
}

//Para guardar un chat con sus datos y mensajes que tenga
type Chat struct {
	id       int
	nombre   string
	mensajes []Mensaje
}

func obtenerChatsUsuarioBD(idusuario int) []Chat {

	chats := make([]Chat, 0, 1)       //Todos los chats del usuario
	mensajes := make([]Mensaje, 0, 1) //Los mensajes de un chat

	var chat Chat       // Para ir introduciendo chats al slice
	var mensaje Mensaje //Para ir introduciendo mensajes al slice

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return nil
	}
	defer db.Close()

	//Obtenemos los id de los chats en los que está el usuario
	rows, err := db.Query("SELECT idchat FROM usuarioschat WHERE idusuario = " + strconv.Itoa(idusuario))
	if err != nil {
		panic(err.Error())
		defer db.Close()
		return nil
	}

	//Guardamos id de cada chat en el slice de chats
	for rows.Next() {
		err = rows.Scan(&chat.id)

		if err != nil {
			panic(err.Error())
			defer db.Close()
			return nil
		}
		chats = append(chats, chat)
	}

	//Por cada id de chat obtenemos datos del chat y los mensajes del chat
	for i := 0; i < len(chats); i++ {

		//De cada chat obtenemos sus datos (nombre...)
		rows, err := db.Query("SELECT nombre FROM chat WHERE id = " + strconv.Itoa(chats[i].id))
		if err != nil {
			panic(err.Error())
			defer db.Close()
			return nil
		}

		for rows.Next() {
			var nombrechat string
			err = rows.Scan(&nombrechat)

			if err != nil {
				chats[i].nombre = ""
			} else {
				chats[i].nombre = nombrechat
			}

		}

		//De cada chat buscamos los datos de los mensajes de dicho chat
		rows, err = db.Query("SELECT id, texto, emisor FROM mensaje WHERE chat = " + strconv.Itoa(chats[i].id))
		if err != nil {
			panic(err.Error())
			defer db.Close()
			return nil
		}

		for rows.Next() {

			//Obtenemos los datos del mensaje
			err = rows.Scan(&mensaje.id, &mensaje.texto, &mensaje.idemisor)

			if err != nil {
				panic(err.Error())
				defer db.Close()
				return nil
			}

			mensaje.nombreemisor = getNombreUsuario(mensaje.idemisor)

			//Para ver si un mensaje aparece como leido o no
			rows2, err2 := db.Query("SELECT leido from receptoresmensaje where idmensaje = " + strconv.Itoa(mensaje.id))
			if err2 != nil {
				panic(err2.Error())
				defer db.Close()
				return nil
			}
			for rows2.Next() {
				err2 = rows2.Scan(&mensaje.leido)
				//Si no aparece, el mensaje es suyo propio, siempre lo habrá leido
				if err2 != nil {
					mensaje.leido = true
				}
			}

			//Guardamos el mensaje en el array de mensajes
			mensajes = append(mensajes, mensaje)
		}

		//Añadimos el array de mensajes a este chat
		chats[i].mensajes = mensajes

		//Vaciamos el array de mensajes, para rellenar el próximo chat
		mensajes = make([]Mensaje, 0, 1)
	}

	return chats
}

func main() {
	var test bool

	//Prueba insertar usuario
	//insertUsuarioBD("lolo", "clave4rsa", "clave4cifrada")

	//Prueba Modificar Usuario
	//test = modificarUsuarioBD(15, "clavepubrsa15", "clave15cifrada-")
	//fmt.Println("Mira modificar usuario:", test)
	//fmt.Println("-")

	//Prueba comprobar usuario
	test = comprobarUsuarioBD("pepe", "clave1cifrada")
	fmt.Println("Mira comprobando usuario:", test)
	fmt.Println("-")

	//Probar obtener nombre según id
	nombreusuario := getNombreUsuario(1)
	fmt.Println("Mira el nombre del usuario:", nombreusuario)

	//Prueba crear chat
	usuarios := make([]int, 0, 1)
	usuarios = append(usuarios, 1)
	usuarios = append(usuarios, 2)
	usuarios = append(usuarios, 3)
	//test = crearChatBD(usuarios, "")
	//fmt.Println("Mira crear chat:", test)
	//fmt.Println("\n")

	//Prueba modificar chat
	//test = modificarChatBD(5, "grupo molon")
	//fmt.Println("Mira modificar chat:", test)
	//fmt.Println("\n")

	//Prueba añadir usuarios a un char
	nuevosusuarios := make([]int, 0, 1)
	nuevosusuarios = append(nuevosusuarios, 4)
	nuevosusuarios = append(nuevosusuarios, 5)
	//test = addUsuariosChatBD(7, nuevosusuarios)
	//fmt.Println("Mira añadir nuevos usuarios a chat:", test)

	//Prueba eliminar usuarios de un char
	usuariosexpulsados := make([]int, 0, 1)
	usuariosexpulsados = append(usuariosexpulsados, 4)
	usuariosexpulsados = append(usuariosexpulsados, 5)
	//test = removeUsuariosChatBD(7, usuariosexpulsados)
	//fmt.Println("Mira eliminar usuarios a chat:", test)

	//Prueba guardar mensaje
	//test = guardarMensajeBD("Hola que tal?? :)", 5, 1, 1)
	//fmt.Println("Mira guardar mensaje:", test)
	//fmt.Println("-")

	//Prueba obtener mensajes
	//test = obtenerMensajeBD("Hola que tal?? :)", 5, 1, 1)
	//fmt.Println("Mira guardar mensaje:", test)
	//fmt.Println("-")

	//Prueba crear nueva clavesmensajes
	//id := CrearNuevaClaveParaMensajesBD()
	//fmt.Println("Mira id clavesmensajes creado:", id)
	//fmt.Println("-")

	//Prueba insertar clave de un usuario para x mensajes
	//test = GuardarClaveUsuarioMensajesBD(1, "claveusuario1", 1)
	//fmt.Println("Mira guardar clave usuario de x mensaje:", test)
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
