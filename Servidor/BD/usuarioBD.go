/////////
//USUARIO
/////////

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

var username = "sds"
var password = "sds"
var adress = ""
var database = "sds"

type Usuario struct {
	id           int
	nombre       string
	clavepubrsa  string
	claveprivrsa string
	claveusuario string
}

//Insertamos a un nuevo usuario en BD
func insertUsuarioBD(usuario Usuario) bool {

	//Conexión BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO usuario VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
		return false
	}

	//Insertamos
	_, err = stmtIns.Exec("DEFAULT", usuario.nombre, usuario.clavepubrsa, usuario.claveprivrsa, usuario.claveusuario)
	if err != nil {
		panic(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

//Obtenemos nombre de usuario según id usuario
func getNombreUsuario(id int) string {

	var nombreusuario string

	//Conexión BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return ""
	}
	defer db.Close()

	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT nombre FROM usuario WHERE id = " + strconv.Itoa(id))
	if err != nil {
		panic(err.Error())
		defer db.Close()
		return ""
	}

	for rows.Next() {
		err = rows.Scan(&nombreusuario)
		if err != nil {
			panic(err.Error())
			defer db.Close()
			return ""
		}
	}

	return nombreusuario
}

//Obtenemos una instancia de usuario según id usuario
func getUsuario(id int) Usuario {

	var usuario Usuario

	//Conexión BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return usuario
	}
	defer db.Close()

	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT id, nombre, clavepubrsa, claveprivrsa, claveusuario FROM usuario WHERE id = " + strconv.Itoa(id))
	if err != nil {
		panic(err.Error())
		defer db.Close()
		return usuario
	}

	for rows.Next() {
		err = rows.Scan(&usuario.id, &usuario.nombre, &usuario.clavepubrsa, &usuario.claveprivrsa, &usuario.claveusuario)
		if err != nil {
			panic(err.Error())
			defer db.Close()
			return usuario
		}
	}

	return usuario
}

func modificarUsuarioBD(usuario Usuario) bool {

	//Conexion BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()

	nombreu := getNombreUsuario(usuario.id)
	if nombreu == "" {
		return false
	}

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("UPDATE usuario set clavepubrsa=?, claveprivrsa=?, claveusuario=? where id=?")
	if err != nil {
		panic(err.Error())
		return false
	}

	//Insertamos crear el chat
	_, err = stmtIns.Exec(usuario.clavepubrsa, usuario.claveprivrsa, usuario.claveusuario, usuario.id)
	if err != nil {
		panic(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

// Comprobamos un usuario con su nombre y clave cifrada
func comprobarUsuarioBD(nombre string, claveusuario string) bool {

	var idusuario int
	var claveusuariobd string

	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
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

	return true
}

func main() {
	var test bool

	//Prueba insertar usuario
	var uu Usuario
	uu.nombre = "alex"
	uu.clavepubrsa = "clavepubrsa"
	uu.claveprivrsa = "claveprivrsa"
	uu.claveusuario = "clavecifrada"
	//insertUsuarioBD(uu)

	//Prueba Modificar Usuario
	var u Usuario
	u.id = 15
	u.clavepubrsa = "clave15pubrsa"
	u.claveprivrsa = "clave15privrsa"
	u.claveusuario = "clave15cifrada"
	test = modificarUsuarioBD(u)
	fmt.Println("Mira modificar usuario:", test)

	//Probar obtener nombre según id
	nombreusuario := getNombreUsuario(1)
	fmt.Println("Mira el nombre del usuario:", nombreusuario)

	//Probar obtener usuario según id
	usuario := getUsuario(1)
	fmt.Println("Mira el usuario:", usuario.id, usuario.nombre, usuario.clavepubrsa, usuario.claveprivrsa, usuario.claveusuario)

	//Prueba comprobar usuario
	test = comprobarUsuarioBD("pepe", "clave1cifrada")
	fmt.Println("Mira comprobando usuario:", test)
	fmt.Println("-")

}
