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
var username = "root"
var password = "ViadJid3"

//var username = "root"
//var password = "chn6NQXM"
var adress = "51.255.44.18" //vps222360.ovh.net
var database = "sds"

//Insertamos a un nuevo usuario en BD
func insertUsuarioBD(nombre string, clavepubrsa string) {
	//db, err := sql.Open("mysql", username+":"+password+"@tcp("+adress+":8080)/"+database)
	db, err := sql.Open("mysql", username+":"+password+"@"+adress+"/"+database)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO usuario VALUES(" + nombre + "," + clavepubrsa + ")")
	if err != nil {
		panic(err.Error()) // Error handling
	}

	defer stmtIns.Close() // Close the statement
}

/*/Search in the BD if there are stocks for the request of the store
func searchBD(elementoPedido string) (int, float32) {
	//Connection with the BD MTISMayorista
	db, err := sql.Open("mysql", "mtis:mtis@/MTISMayorista")
	if err != nil {
		//panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return 0, 0
	}

	// we "scan" the result in here
	var cantidadReal int = 0
	var precio float32 = 0

	// Query the quantity (cantidad) of nombre = rueda
	rows, err := db.Query("SELECT cantidad, precio, nombre FROM almacen WHERE nombre = '" + elementoPedido + "'")
	if err != nil {
		//Closing connection
		defer db.Close()
		return 0, 0
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&cantidadReal, &precio, &name)
		if err != nil {
			//Closing connection
			defer db.Close()
			return 0, 0
		}
	}

	//Closing connection
	defer db.Close()

	return cantidadReal, precio
}*/

func main() {
	insertUsuarioBD("pepe", "clave1")
}
