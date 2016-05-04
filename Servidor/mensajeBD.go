/////////
//MENSAJE
/////////

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

//Para guardar un mensaje con sus datos
type Mensaje struct {
	Id           int    `json:"Id"`
	Texto        string `json:"Texto"`
	Idemisor     int    `json:"Idemisor"`
	Nombreemisor string `json:"Nombreemisor"`
	Leido        bool   `json:"Leido"`
	Idchat       int    `json:"Idchat"`
	Idclave      int    `json:"Idclave"`
}

//Guarda un mensaje para todos los receptores posibles del chat
func (bd *BD) guardarMensaje(mensaje Mensaje) bool {

	var idreceptoraux = -1
	idreceptores := make([]int, 0, 1)

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	//Obtenemos los id de los usuarios que recibiráne el mensaje en este chat
	rows, err := db.Query("SELECT idusuario FROM usuarioschat WHERE idchat = " + strconv.Itoa(mensaje.Idchat) + " and idusuario !=" + strconv.Itoa(mensaje.Idemisor))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return false
	}

	//Guardamos id de los usuarios receptores en slice de ids
	for rows.Next() {
		err = rows.Scan(&idreceptoraux)

		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return false
		}
		idreceptores = append(idreceptores, idreceptoraux)
	}

	//Preparamos la creación del mensaje
	stmtIns, err := db.Prepare("INSERT INTO mensaje VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos el mensaje
	res, err := stmtIns.Exec("DEFAULT", mensaje.Texto, mensaje.Idemisor, mensaje.Idchat, mensaje.Idclave)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Obtenemos id del mensaje creado
	idmensaje, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	println("Id del mensaje creado:", idmensaje)

	//Por cada receptor
	for i := 0; i < len(idreceptores); i++ {

		//Preparamos la insercion del receptor del mensaje
		stmtIns, err := db.Prepare("INSERT INTO receptoresmensaje VALUES(?, ?, ?)")
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		//Insertamos el receptor con el mensaje
		_, err = stmtIns.Exec(idmensaje, idreceptores[i], "DEFAULT")
		if err != nil {
			fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		return 0
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO clavesmensajes VALUES(?)")
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	//Insertamos
	res, err := stmtIns.Exec("DEFAULT")
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	//Obtenemos id de lo creado
	idclavesmensajes, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	defer stmtIns.Close()

	return idclavesmensajes
}

//Guardamos la clave de un usuario para leer x mensajes
func (bd *BD) GuardarClaveUsuarioMensajesBD(idusuario int, idclavesmensajes int, clavemensajes string) bool {

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO clavesusuario VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos
	_, err = stmtIns.Exec(idusuario, idclavesmensajes, clavemensajes)
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		return nil
	}
	defer db.Close()

	//De el chat buscamos los datos de los mensajes de dicho chat
	rows, err := db.Query("SELECT id, texto, emisor FROM mensaje WHERE chat = " + strconv.Itoa(idchat))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return nil
	}

	for rows.Next() {
		//Obtenemos los datos del mensaje
		err = rows.Scan(&mensaje.Id, &mensaje.Texto, &mensaje.Idemisor)

		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return nil
		}

		mensaje.Nombreemisor = bd.getNombreUsuario(mensaje.Idemisor)

		//Para ver si un mensaje aparece como leido o no
		rows2, err2 := db.Query("SELECT leido from receptoresmensaje where idmensaje = " + strconv.Itoa(mensaje.Id))
		if err2 != nil {
			fmt.Println(err2.Error())
			defer db.Close()
			return nil
		}
		for rows2.Next() {
			err2 = rows2.Scan(&mensaje.Leido)
			//Si no aparece, el mensaje es suyo propio, siempre lo habrá leido
			if err2 != nil {
				mensaje.Leido = true
			}
		}

		//Guardamos el mensaje en el array de mensajes
		mensajes = append(mensajes, mensaje)
	}

	return mensajes
}

//Obtiene la clave cifrada con la que se cifran los mensajes
func (bd *BD) getClaveMensaje(idmensaje int) (string, bool) {

	var clavemensaje string

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}
	defer db.Close()

	//De el chat buscamos los datos de los mensajes de dicho chat
	rows, err := db.Query("SELECT claveusuario FROM mensaje, clavesusuario WHERE id = " + strconv.Itoa(idmensaje))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return "", false
	}

	for rows.Next() {
		//Obtenemos los datos del mensaje
		err = rows.Scan(&clavemensaje)

		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return "", false
		}

	}

	return clavemensaje, true
}

//Obtiene los datos de un mensaje
func (bd *BD) getMensaje(idmensaje int) (Mensaje, bool) {

	var mensaje Mensaje
	mensaje.Id = idmensaje

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return mensaje, false
	}
	defer db.Close()

	//De el chat buscamos los datos de los mensajes de dicho chat
	rows, err := db.Query("SELECT texto, emisor, chat, clave FROM mensaje WHERE id = " + strconv.Itoa(idmensaje))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return mensaje, false
	}

	for rows.Next() {
		//Obtenemos los datos del mensaje
		err = rows.Scan(&mensaje.Texto, &mensaje.Idemisor, &mensaje.Idchat, &mensaje.Idclave)

		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return mensaje, false
		}

	}

	return mensaje, true
}

//Obtiene la última clave (con la que se están cifrando ahora los mensajes)
func (bd *BD) getLastKeyMensaje(idchat int, idusuario int) (string, bool) {

	var clave string

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}
	defer db.Close()

	//De el chat buscamos los datos de los mensajes de dicho chat
	rows, err := db.Query("SELECT claveusuario FROM clavesusuario, mensaje WHERE idusuario = " + strconv.Itoa(idusuario) + " AND chat = " + strconv.Itoa(idchat) + " AND clave = (select max(clave) from clavesusuario, mensaje WHERE idusuario = " + strconv.Itoa(idusuario) + " AND chat = " + strconv.Itoa(idchat) + ")")
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return "", false
	}

	for rows.Next() {
		//Obtenemos los datos del mensaje
		err = rows.Scan(&clave)

		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return "", false
		}

	}

	return clave, true
}

func (bd *BD) marcarLeido(idMensaje int) {
	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	//De el chat buscamos los datos de los mensajes de dicho chat
	_, err = db.Query("UPDATE `receptoresmensaje` SET `leido`=true WHERE `idmensaje`=" + strconv.Itoa(idMensaje))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
	}
}
