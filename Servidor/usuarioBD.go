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

type Usuario struct {
	id           int
	nombre       string
	clavepubrsa  string
	claveprivrsa string
	claveusuario string
}

//Funcion para obtener los datos del usuario cuando se loguea
func (usuario *Usuario) login(nombre string, password string) bool {

	//Para las operaciones con la BD
	var bd BD
	bd.username = "root"
	bd.password = ""
	bd.adress = ""
	bd.database = "sds"

	usuario.nombre = nombre
	user, test := bd.comprobarUsuarioBD(nombre, password)

	if test == false {
		return false
	}

	usuario.id = user.id
	usuario.nombre = user.nombre
	usuario.clavepubrsa = user.clavepubrsa
	usuario.claveprivrsa = user.claveprivrsa
	usuario.claveusuario = user.claveusuario
	usuario.nombre = user.nombre
	return true
}

// Comprobamos un usuario con su nombre y clave cifrada
func (bd *BD) comprobarUsuarioBD(nombre string, claveusuario string) (Usuario, bool) {

	var usuario Usuario

	usuario.nombre = nombre
	usuario.claveusuario = claveusuario

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return usuario, false
	}
	defer db.Close()

	//Obtenemos el id del usuario
	rows, err := db.Query("SELECT id, clavepubrsa, claveprivrsa FROM usuario WHERE nombre = '" + nombre + "' and claveusuario= '" + claveusuario + "'")
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return usuario, false
	}

	for rows.Next() {
		err = rows.Scan(&usuario.id, &usuario.clavepubrsa, &usuario.claveprivrsa)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return usuario, false
		}
	}

	if usuario.id == 0 {
		return usuario, false
	}

	return usuario, true
}

//Obtenemos nombre de usuario según id usuario
func (bd *BD) getUsuarioBD(user string) Usuario {
	usuario := Usuario{}
	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT id, nombre, clavepubrsa, claveprivrsa, claveusuario FROM usuario WHERE nombre = '" + user + "'")
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
	}
	for rows.Next() {
		err = rows.Scan(&usuario.id, &usuario.nombre, &usuario.clavepubrsa, &usuario.claveprivrsa, &usuario.claveusuario)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
		}
	}
	return usuario
}

//Obtener los id de usuarios de un chat
func (bd *BD) getUsuariosChatBD(id int) []int {
	usuarios := []int{}

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT idusuario FROM usuarioschat WHERE idchat = " + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
	}
	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
		}
		usuarios = append(usuarios, i)

	}
	return usuarios
}

//Obtenemos nombre de usuario según id usuario
func (bd *BD) getNombreUsuario(id int) string {

	var nombreusuario string

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer db.Close()

	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT nombre FROM usuario WHERE id = " + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return ""
	}

	for rows.Next() {
		err = rows.Scan(&nombreusuario)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return ""
		}
	}

	return nombreusuario
}

//Obtenemos nombre de usuario según id usuario
func (bd *BD) getClavePubUsuario(id int) string {

	var clavepub string

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer db.Close()

	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT clavepubrsa FROM usuario WHERE id = " + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return ""
	}

	for rows.Next() {
		err = rows.Scan(&clavepub)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return ""
		}
	}

	return clavepub
}

//Obtenemos una instancia de usuario según id usuario
func (bd *BD) getUsuario(id int) Usuario {

	var usuario Usuario

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return usuario
	}
	defer db.Close()

	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT id, nombre, clavepubrsa, claveprivrsa, claveusuario FROM usuario WHERE id = " + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return usuario
	}

	for rows.Next() {
		err = rows.Scan(&usuario.id, &usuario.nombre, &usuario.clavepubrsa, &usuario.claveprivrsa, &usuario.claveusuario)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return usuario
		}
	}

	return usuario
}

func (bd *BD) modificarUsuarioBD(usuario Usuario) bool {

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	nombreu := bd.getNombreUsuario(usuario.id)
	if nombreu == "" {
		return false
	}

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("UPDATE usuario set clavepubrsa=?, claveprivrsa=?, claveusuario=? where id=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos crear el chat
	_, err = stmtIns.Exec(usuario.clavepubrsa, usuario.claveprivrsa, usuario.claveusuario, usuario.id)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

//Insertamos a un nuevo usuario en BD
func (bd *BD) insertUsuarioBD(usuario Usuario) bool {

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO usuario VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos
	_, err = stmtIns.Exec("DEFAULT", usuario.nombre, usuario.clavepubrsa, usuario.claveprivrsa, usuario.claveusuario)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}
