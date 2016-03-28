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
		defer db.Close()
		return false
	}

	for rows.Next() {
		err = rows.Scan(&idusuario)
		if err != nil {
			defer db.Close()
			return false
		}
	}

	if idusuario == 0 {
		return false
	}

	//Obtenemos el la clave del usuario con id obtenido
	rows, err = db.Query("SELECT clave FROM clavesusuario WHERE usuario = " + strconv.Itoa(idusuario))
	if err != nil {
		defer db.Close()
		return false
	}

	for rows.Next() {
		err = rows.Scan(&claveusuariobd)
		if err != nil {
			defer db.Close()
			return false
		}
	}

	//Vemos si claves coinciden
	if claveusuario != claveusuariobd {
		return false
	}

	defer db.Close()

	return true
}

func main() {
	//insertUsuarioBD("maria", "clave1")

	var test bool
	test = comprobarUsuarioBD("pepe", "clave1cifrada")
	fmt.Println("Mira comprobando usuario:", test)

}
