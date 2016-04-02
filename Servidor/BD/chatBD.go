//////
//CHAT
//////

package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

//Para guardar un chat con sus datos y mensajes que tenga
type Chat struct {
	id       int
	nombre   string
	mensajes []Mensaje
}

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

func obtenerChatsUsuarioBD(idusuario int) []Chat {

	chats := make([]Chat, 0, 1)       //Todos los chats del usuario
	mensajes := make([]Mensaje, 0, 1) //Los mensajes de un chat
	var chat Chat                     // Para ir introduciendo chats al slice

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
		mensajes = obtenerMensajesChatBD(chats[i].id)

		//Añadimos el array de mensajes a este chat
		chats[i].mensajes = mensajes

		//Vaciamos el array de mensajes, para rellenar el próximo chat
		mensajes = make([]Mensaje, 0, 1)
	}

	return chats
}
