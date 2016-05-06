//////
//CHAT
//////

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

//Para guardar un chat con sus datos y mensajes que tenga
type Chat struct {
	Id       int       `json:"Id"`
	Nombre   string    `json:"Nombre"`
	Mensajes []Mensaje `json:"Mensajes"`
}

//Creamos nuevo chat en BD
func (bd *BD) crearChatBD(idusuarios []int, nombrechat string) bool {

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("INSERT INTO chat VALUES(?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos crear el chat
	res, err := stmtIns.Exec("DEFAULT", nombrechat)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Obtenemos id del chat creado
	idchat, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	println("Id del chat creado:", idchat)

	defer stmtIns.Close()

	//Insertamos usuarios a dicho chat
	for i := 0; i < len(idusuarios); i++ {
		//Preparamos insertar usuario al chat
		stmtIns, err := db.Prepare("INSERT INTO usuarioschat VALUES(?, ?)")
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		//Insertamos usuario al chat
		_, err = stmtIns.Exec(idusuarios[i], idchat)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	defer stmtIns.Close()

	return true
}

//Modifica los datos del chat
func (bd *BD) modificarChatBD(chat Chat) bool {

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("UPDATE chat set nombre=? where id=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos crear el chat
	_, err = stmtIns.Exec(chat.Nombre, chat.Id)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

//Añade una serie de usuarios a un chat
func (bd *BD) addUsuariosChatBD(idchat int, nuevosusuarios []int) bool {

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	//Insertamos usuarios a dicho chat
	for i := 0; i < len(nuevosusuarios); i++ {
		//Preparamos insertar usuario al chat
		stmtIns, err := db.Prepare("INSERT INTO usuarioschat VALUES(?, ?)")
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		//Insertamos usuario al chat
		_, err = stmtIns.Exec(nuevosusuarios[i], idchat)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		defer stmtIns.Close()
	}

	return true
}

//Elimina una serie de usuarios a un chat
func (bd *BD) removeUsuariosChatBD(idchat int, usuariosexpulsados []int) bool {

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	//Insertamos usuarios a dicho chat
	for i := 0; i < len(usuariosexpulsados); i++ {
		//Preparamos insertar usuario al chat
		stmtIns, err := db.Prepare("DELETE FROM usuarioschat where idusuario=? and idchat=?")
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		//Insertamos usuario al chat
		_, err = stmtIns.Exec(usuariosexpulsados[i], idchat)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		defer stmtIns.Close()
	}

	return true
}

func (bd *BD) getChatsUsuarioBD(idusuario int) []Chat {

	chats := make([]Chat, 0, 1)            //Todos los chats del usuario
	mensajes := make([]MensajeDatos, 0, 1) //Los mensajes de un chat
	var chat Chat                          // Para ir introduciendo chats al slice

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer db.Close()

	//Obtenemos los id de los chats en los que está el usuario
	rows, err := db.Query("SELECT idchat FROM usuarioschat WHERE idusuario = " + strconv.Itoa(idusuario))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return nil
	}
	//Guardamos id de cada chat en el slice de chats
	for rows.Next() {
		err = rows.Scan(&chat.Id)

		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return nil
		}
		chats = append(chats, chat)
	}

	//Por cada id de chat obtenemos datos del chat y los mensajes del chat
	for i := 0; i < len(chats); i++ {

		//De cada chat obtenemos sus datos (nombre...)
		rows, err := db.Query("SELECT nombre FROM chat WHERE id = " + strconv.Itoa(chats[i].Id))
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return nil
		}
		for rows.Next() {
			var nombrechat string
			err = rows.Scan(&nombrechat)

			if err != nil {
				chats[i].Nombre = ""
			} else {
				chats[i].Nombre = nombrechat
			}
		}

		//De cada chat buscamos los datos de los mensajes de dicho chat
		mensajes, _ = bd.getMensajesChatBD(chats[i].Id, idusuario)

		//Añadimos el array de mensajes a este chat
		for j := 0; j < len(mensajes); j++ {
			chats[i].Mensajes[j] = mensajes[j].Mensaje
		}

		//Vaciamos el array de mensajes, para rellenar el próximo chat
		mensajes = make([]MensajeDatos, 0, 1)
	}

	return chats
}

//Ver si un usuario está en un chat
func (bd *BD) usuarioEnChat(idusuario int, idchat int) bool {

	var count int

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	//Vemos si el usuario esta en el chat
	rows, err := db.Query("SELECT count(idusuario) FROM usuarioschat WHERE idusuario = " + strconv.Itoa(idusuario) + " and idchat = " + strconv.Itoa(idchat))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return false
	}

	//Guardamos id del usuario
	for rows.Next() {
		err = rows.Scan(&count)

		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return false
		}
	}

	if count != 0 {
		return true
	}

	return false
}
