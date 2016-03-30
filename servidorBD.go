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

//Insertamos a un nuevo usuario en BD
func insertUsuarioBD(nombre string, clavepubrsa string) {
	//db, err := sql.Open("mysql", username+":"+password+"@tcp(:3306)/"+database)
	//db, err := sql.Open("mysql", username+":"+password+"@"+adress+"/"+database)
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO usuario VALUES(?, ?, ?)")
	if err != nil {
		panic(err.Error()) // Error handling
	}

	//Insertamos
	_, err = stmtIns.Exec("DEFAULT", nombre, clavepubrsa)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtIns.Close()
}

// Comprobamos un usuario con su nombre y clave cifrada
func comprobarUsuarioBD(nombre string, claveusuario string) bool {

	var idusuario int
	var claveusuariobd string

	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
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

	//	defer db.Close()

	return true
}

//Creamos nuevo chat en BD
func crearChatBD(usuarios []string, nombrechat string) bool {

	idusuarios := make([]int, 0, 1)

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Obtenemos los id de los usuarios según su nombre
	for i := 0; i < len(usuarios); i++ {

		//Consulta para saber id del usuario
		rows, err := db.Query("SELECT id FROM usuario WHERE nombre = '" + usuarios[i] + "'")
		if err != nil {
			panic(err.Error())
			defer db.Close()
			return false
		}

		//Guardamos id del usuario en slice de ids
		for rows.Next() {
			var idusuario int
			err = rows.Scan(&idusuario)
			if err != nil {
				panic(err.Error())
				defer db.Close()
				return false
			}

			idusuarios = append(idusuarios, idusuario)
		}
	}

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

/*type Mensaje struct {
	texto     string
	idusuario int
	usuario   string
}

func obtenerMensajesBD(idusuario int) []bool {

}*/

func main() {
	//insertUsuarioBD("maria", "clave1")

	//Prueba comprobar usuario
	var test bool
	test = comprobarUsuarioBD("pepe", "clave1cifrada")
	fmt.Println("Mira comprobando usuario:", test)
	fmt.Println("-")

	//Prueba crear chat
	usuarios := make([]string, 0, 1)
	usuarios = append(usuarios, "pepe")
	usuarios = append(usuarios, "lucia")
	usuarios = append(usuarios, "maria")

	//test = crearChatBD(usuarios)
	//fmt.Println("Mira crear chat:", test)
	//fmt.Println("\n")

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
}
