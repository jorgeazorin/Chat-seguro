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
	Id          int    `json:"Id"`
	Nombre      string `json:"Nombre"`
	UltimaClave int    `json:"UltimaClave"`
}

type UsuariosChat struct {
	Idusuario int `json:"Idusuario"`
	Idchat    int `json:"Idchat"`
}

//Para guardar un chat con sus datos y mensajes que tenga
type ChatDatos struct {
	Chat          Chat           `json:"Chat"`
	MensajesDatos []MensajeDatos `json:"Mensajes"`
	Clave         []byte         `json:"Clave"`
	IdClave       int            `json:"IdClave"`
}

//Creamos nuevo chat en BD
func (bd *BD) crearChatBD(idusuarios []int, nombrechat string) (bool, int) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false, 0
	}

	//Insert de un nuevo chat
	chat := Chat{Nombre: nombrechat}
	err := dbmap.Insert(&chat)
	if err != nil {
		fmt.Println("Error1:", err.Error())
		return false, 0
	}
	//Formamos usuarioschat para el insert e insertamos
	for i := 0; i < len(idusuarios); i++ {
		usuariochat := UsuariosChat{Idchat: chat.Id, Idusuario: idusuarios[i]}
		err = dbmap.Insert(&usuariochat)
		if err != nil {
			fmt.Println("Error2:", err.Error())
			return false, 0
		}
	}

	return true, chat.Id
}

//Añade una serie de usuarios a un chat
func (bd *BD) AsociarNuevaClaveAChat(idchat int, idClave int) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Insertamos usuarios a dicho chat
	_, err := dbmap.Exec("UPDATE `chat` SET `ultimaClave`=? WHERE `id`=?", idClave, idchat)
	if err != nil {
		fmt.Println(err.Error())
		return false

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
			fmt.Println("Error:", usuarioschat[i].Idchat, err.Error())
			return []ChatDatos{}, false
		}
		chatdatos.Chat = chat

		//Obtenemos mensajes del chat
		mensajesdatos, test := bd.getMensajesChatBD(usuarioschat[i].Idchat, idusuario)
		if test == false {
			fmt.Println("Error")
			return []ChatDatos{}, false
		}

		//Introducimos mensajes al chat y chat al array de chats
		chatdatos.MensajesDatos = mensajesdatos
		clavechat, idclavechat, _ := bd.getLastKeyMensaje(usuarioschat[i].Idchat, idusuario)
		chatdatos.Clave = clavechat
		chatdatos.IdClave = idclavechat

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

//Marcar todos los mensajes del chat como leidos
func (bd *BD) marcarChatLeidoPorUsuarioBD(idchat int, idreceptor int) bool {

	mensajes, test := bd.getMensajesChatBD(idchat, idreceptor)
	if test == false {
		return false
	}

	for i := 0; i < len(mensajes); i++ {
		if !mensajes[i].Mensaje.Admin {
			test = bd.marcarLeidoPorUsuarioBD(mensajes[i].Mensaje.Id, idreceptor)
			if test == false {
				return false
			}
		}
	}

	return true
}

func (bd *BD) usuariosEnChat(idchat int) ([]int, bool) {
	usuarios := make([]int, 0, 1)
	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return usuarios, false
	}
	//usuariosChat := UsuariosChat[]{}
	usuarioschat := make([]UsuariosChat, 0, 1)

	//Vemos si el usuario esta en el chat
	_, err := dbmap.Select(&usuarioschat, "SELECT idusuario FROM usuarioschat WHERE idchat = ?", idchat)
	if err != nil {
		fmt.Println(err.Error())
		return usuarios, false
	}

	for i := 0; i < len(usuarioschat); i++ {
		usuarios = append(usuarios, usuarioschat[i].Idusuario)
	}

	return usuarios, true
}
