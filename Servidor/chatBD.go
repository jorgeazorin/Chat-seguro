//////
//CHAT
//////

package main

import (
	//	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//	"strconv"
)

type Chat struct {
	Id     int    `json:"Id"`
	Nombre string `json:"Nombre"`
}

type UsuariosChat struct {
	Idusuario int `json:"Idusuario"`
	Idchat    int `json:"Idchat"`
}

//Para guardar un chat con sus datos y mensajes que tenga
type ChatDatos struct {
	Chat          Chat           `json:"Chat"`
	MensajesDatos []MensajeDatos `json:"Mensajes"`
}

//Creamos nuevo chat en BD
func (bd *BD) crearChatBD(idusuarios []int, nombrechat string) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Insert de un nuevo chat
	chat := Chat{Nombre: nombrechat}
	err := dbmap.Insert(&chat)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return false
	}

	//Formamos usuarioschat para el insert e insertamos
	for i := 0; i < len(idusuarios); i++ {
		usuariochat := UsuariosChat{Idchat: chat.Id, Idusuario: idusuarios[i]}
		err = dbmap.Insert(&usuariochat)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return false
		}
	}

	return true
}

//Modifica los datos del chat
func (bd *BD) modificarChatBD(chat Chat) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Insert de un nuevo chat
	_, err := dbmap.Update(&chat)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return false
	}

	return true
}

//Añade una serie de usuarios a un chat
func (bd *BD) addUsuariosChatBD(idchat int, nuevosusuarios []int) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Insertamos usuarios a dicho chat
	for i := 0; i < len(nuevosusuarios); i++ {
		_, err := dbmap.Exec("INSERT INTO usuarioschat VALUES(?, ?)", nuevosusuarios[i], idchat)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	return true
}

//Elimina una serie de usuarios a un chat
func (bd *BD) removeUsuariosChatBD(idchat int, usuariosexpulsados []int) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Insertamos usuarios a dicho chat
	for i := 0; i < len(usuariosexpulsados); i++ {
		_, err := dbmap.Exec("DELETE FROM usuarioschat where idusuario = ? and idchat = ?", usuariosexpulsados[i], idchat)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	return true
}

func (bd *BD) getChatsUsuarioBD(idusuario int) ([]ChatDatos, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return []ChatDatos{}, false
	}

	//Obtenemos ids de los chats del usuario obteniendo usuarioschat
	usuarioschat := make([]UsuariosChat, 0, 1)
	_, err := dbmap.Select(&usuarioschat, "SELECT * FROM usuarioschat WHERE idusuario = ?", idusuario)
	if err != nil {
		fmt.Println("Error1:", err.Error())
		return []ChatDatos{}, false
	}

	//Obtenemos cada chat del usuario con sus mensajes
	chatsdatos := make([]ChatDatos, 0, 1)
	for i := 0; i < len(usuarioschat); i++ {

		//Obtenemos info Chat al completo
		var chatdatos ChatDatos
		var chat Chat
		err := dbmap.SelectOne(&chat, "SELECT * FROM chat WHERE id = ?", usuarioschat[i].Idchat)
		if err != nil {
			fmt.Println("Error2:", usuarioschat[i].Idchat, err.Error())
			return []ChatDatos{}, false
		}
		chatdatos.Chat = chat

		//Obtenemos mensajes del chat
		mensajesdatos, test := bd.getMensajesChatBD(usuarioschat[i].Idchat, idusuario)
		if test == false {
			fmt.Println("Error3")
			return []ChatDatos{}, false
		}

		//Nombre al emisor
		for i := 0; i < len(mensajesdatos); i++ {
			mensajesdatos[i].Mensaje.NombreEmisor, _ = bd.getNombreUsuario(mensajesdatos[i].Mensaje.Emisor)
		}

		//Introducimos mensajes al chat y chat al array de chats
		chatdatos.MensajesDatos = mensajesdatos
		chatsdatos = append(chatsdatos, chatdatos)
	}

	return chatsdatos, true
}

//Ver si un usuario está en un chat
func (bd *BD) usuarioEnChat(idusuario int, idchat int) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Vemos si el usuario esta en el chat
	count, err := dbmap.SelectInt("SELECT count(idusuario) FROM usuarioschat WHERE idusuario = ? AND idchat = ?", idusuario, idchat)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if count != 0 {
		return true
	}

	return false
}
