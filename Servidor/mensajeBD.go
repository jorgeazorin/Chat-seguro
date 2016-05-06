/////////
//MENSAJE
/////////

package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

//Para guardar un mensaje con sus datos
type Mensaje struct {
	Id     int    `json:"Id"`
	Texto  string `json:"Texto"`
	Emisor int    `json:"Emisor"`
	Chat   int    `json:"Chat"`
	Clave  int    `json:"Clave"`
}

type Receptoresmensaje struct {
	Idmensaje  int  `json:"Idmensaje"`
	Idreceptor int  `json:"Idreceptor"`
	Leido      bool `json:"Leido"`
}

type Clavesmensajes struct {
	Id int `json:"Id"`
}

type Clavesusuario struct {
	Idusuario        int    `json:"Idusuario"`
	Idclavesmensajes int    `json:"Idclavesmensajes"`
	Clavemensajes    []byte `json:"Clavemensajes"`
}

type MensajeDatos struct {
	Mensaje Mensaje `json:"Mensaje"`
	Leido   bool    `json:"Leido"`
}

//Guarda un mensaje para todos los receptores posibles del chat
func (bd *BD) guardarMensajeBD(mensaje Mensaje) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Select
	idreceptores := make([]int, 0, 1)
	_, err := dbmap.Select(&idreceptores, "SELECT idusuario FROM usuarioschat WHERE idchat = ? and idusuario != ?", mensaje.Chat, mensaje.Emisor)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return false
	}

	//Insert
	err = dbmap.Insert(&mensaje)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return false
	}

	//Rellenamos un receptor e insertamos, asi con todos
	for i := 0; i < len(idreceptores); i++ {
		var receptor Receptoresmensaje
		receptor.Idmensaje = mensaje.Id
		receptor.Idreceptor = idreceptores[i]
		receptor.Leido = false

		err = dbmap.Insert(&receptor)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return false
		}
	}

	return true
}

//Crear nuevo id para nuevo grupo de claves para siguientes mensajes
func (bd *BD) CrearNuevaClaveMensajesBD() (int, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return 0, false
	}

	//Insert
	var clavesmensajes Clavesmensajes
	err := dbmap.Insert(&clavesmensajes)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return 0, false
	}

	return clavesmensajes.Id, true
}

//Guardamos la clave de un usuario para leer x mensajes
func (bd *BD) GuardarClaveUsuarioMensajesBD(clavesusuario Clavesusuario) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Insert
	err := dbmap.Insert(&clavesusuario)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return false
	}

	return true
}

//Obtenemos todos los mensajes de un chat
func (bd *BD) getMensajesChatBD(idchat int, idusuario int) ([]MensajeDatos, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return []MensajeDatos{}, false
	}

	//De el chat buscamos los datos de los mensajes de dicho chat
	mensajes := make([]Mensaje, 0, 1)
	_, err := dbmap.Select(&mensajes, "SELECT * FROM mensaje WHERE chat = ?", idchat)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []MensajeDatos{}, false
	}

	mismensajes := make([]MensajeDatos, 0, 1)

	for i := 0; i < len(mensajes); i++ {
		//Vemos más datos como si el mensaje está leído
		var recetoresmensajes Receptoresmensaje
		_, err := dbmap.Select(&recetoresmensajes, "SELECT * FROM receptoresmensaje WHERE idmensaje = ? and idreceptor = ?", mensajes[i].Id, idusuario)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return []MensajeDatos{}, false
		}
		//Rellenamos los datos
		var mimensaje MensajeDatos
		mimensaje.Mensaje = mensajes[i]
		mimensaje.Leido = recetoresmensajes.Leido
		mismensajes = append(mismensajes, mimensaje)
	}

	return mismensajes, true
}

//Obtiene la clave cifrada con la que se cifran los mensajes
func (bd *BD) getClaveMensaje(idmensaje int, idusuario int) ([]byte, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return []byte{}, false
	}

	//Buscamos idclave de ese mensaje
	var mensaje Mensaje
	err := dbmap.SelectOne(&mensaje, "SELECT * FROM mensaje WHERE id = ?", idmensaje)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []byte{}, false
	}

	//Buscamos la clave que tiene ese idclave
	var clavesusuario Clavesusuario
	err = dbmap.SelectOne(&clavesusuario, "SELECT * FROM clavesusuario WHERE idclavesmensajes = ? AND idusuario = ?", mensaje.Clave, idusuario)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []byte{}, false
	}

	return clavesusuario.Clavemensajes, true
}

//Obtiene los datos de un mensaje
func (bd *BD) getMensajeBD(idmensaje int) (Mensaje, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return Mensaje{}, false
	}

	//Buscamos los datos de mensaje en concreto
	var mensaje Mensaje
	err := dbmap.SelectOne(&mensaje, "SELECT * FROM mensaje WHERE id = ?", idmensaje)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return Mensaje{}, false
	}

	return mensaje, true
}

//Obtiene la última clave (con la que se están cifrando ahora los mensajes)
func (bd *BD) getLastKeyMensaje(idchat int, idusuario int) ([]byte, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return []byte{}, false
	}

	//Buscamos los datos de mensaje en concreto (el último mensaje de este chat)
	var mensaje Mensaje
	err := dbmap.SelectOne(&mensaje, "SELECT * FROM mensaje WHERE chat = ? AND id = (SELECT max(id) FROM mensaje WHERE chat = ?)", idchat, idchat)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []byte{}, false
	}

	//Buscamos la clave que tiene ese idclave del mensaje
	var clavesusuario Clavesusuario
	err = dbmap.SelectOne(&clavesusuario, "SELECT * FROM clavesusuario WHERE idclavesmensajes = ? AND idusuario = ?", mensaje.Clave, idusuario)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []byte{}, false
	}

	return clavesusuario.Clavemensajes, true
}

func (bd *BD) marcarLeidoPorUsuarioBD(idmensaje int, idreceptor int) bool {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Marcamos como leido
	_, err := dbmap.Exec("UPDATE receptoresmensaje SET leido = true WHERE idmensaje = ? and idreceptor = ?", idmensaje, idreceptor)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
