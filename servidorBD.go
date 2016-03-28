/*
	Encarna Amorós Beneite, Jorge Azorín Martí
	Práctica SDS
*/

package main

import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"reflect"
	//"strconv"
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

func main() {
	insertUsuarioBD("maria", "clave1")
}
