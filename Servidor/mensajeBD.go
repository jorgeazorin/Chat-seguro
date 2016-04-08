/////////
//MENSAJE
/////////

package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

//Para guardar un mensaje con sus datos
type Mensaje struct {
	id           int    `json:"Id"`
	texto        string `json:"Texto"`
	idemisor     int    `json:"Idemisor"`
	nombreemisor string `json:"Nombreemisor"`
	leido        bool   `json:"Leido"`
	idchat       int    `json:"Idchat"`
	idclave      int    `json:"Idclave"`
}

//Guarda un mensaje para todos los receptores posibles del chat
func (bd *BD) guardarMensaje(mensaje Mensaje) bool {

	var idreceptoraux = -1
	idreceptores := make([]int, 0, 1)

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Obtenemos los id de los usuarios que recibiráne el mensaje en este chat
	rows, err := db.Query("SELECT idusuario FROM usuarioschat WHERE idchat = " + strconv.Itoa(mensaje.idchat) + " and idusuario !=" + strconv.Itoa(mensaje.idemisor))
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
	res, err := stmtIns.Exec("DEFAULT", mensaje.texto, mensaje.idemisor, mensaje.idchat, mensaje.idclave)
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

//Crear nuevo id para nuevo grupo de claves para siguientes mensajes
func (bd *BD) CrearNuevaClaveMensajesBD() int64 {

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

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

//Guardamos la clave de un usuario para leer x mensajes
func (bd *BD) GuardarClaveUsuarioMensajesBD(idclavesmensajes int, claveusuario string, idusuario int) bool {

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

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

//Obtenemos todos los mensajes de un chat
func (bd *BD) getMensajesChatBD(idchat int) []Mensaje {

	mensajes := make([]Mensaje, 0, 1) //Los mensajes de un chat
	var mensaje Mensaje               //Para ir introduciendo mensajes al slice

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		panic(err.Error())
		return nil
	}
	defer db.Close()

	//De cada chat buscamos los datos de los mensajes de dicho chat
	rows, err := db.Query("SELECT id, texto, emisor FROM mensaje WHERE chat = " + strconv.Itoa(idchat))
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

		mensaje.nombreemisor = bd.getNombreUsuario(mensaje.idemisor)

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

	return mensajes
}
