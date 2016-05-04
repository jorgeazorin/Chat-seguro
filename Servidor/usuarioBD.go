/////////
//USUARIO
/////////

package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
	"strconv"
)

type Usuario struct {
	Id           int    `json:"Id"`
	Nombre       string `json:"Nombre"`
	Clavepubrsa  string `json:"Clavepubrsa"`
	Claveprivrsa string `json:"Claveprivrsa"`
	Clavelogin   []byte `json:"Clavelogin"`
	Salt         []byte `json:"Salt"`
	Clavecifrado []byte `json:"Clavecifrado"`
}

//Generar clave con scrypt dada una salt
func generarClaveLoginConSalt(clavehashlogin []byte, salt []byte) []byte {

	//La parte del login hacemos BCRYPT con SALT que nos dan
	clavebcryptlogin, _ := scrypt.Key(clavehashlogin, salt, 16384, 8, 1, 64)

	return clavebcryptlogin
}

//Generar clave con scrypt generando una salt
func generarClaveLoginClaves(clavehashlogin []byte) ([]byte, []byte) {

	//La parte del login hacemos BCRYPT con SALT
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}

	clavebcryptlogin, _ := scrypt.Key(clavehashlogin, salt, 16384, 8, 1, 64)

	return clavebcryptlogin, salt
}

//Insertamos a un nuevo usuario en BD
func (bd *BD) insertUsuarioBD(usuario Usuario) (Usuario, bool) {

	var usuariobd Usuario

	//genera la clavelogin con scrypt (y el salt con el que se crea)
	clavelogin, salt := generarClaveLoginClaves(usuario.Clavelogin)

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return usuariobd, false
	}
	defer db.Close()

	//Preparamos consulta
	stmtIns, err := db.Prepare("INSERT INTO usuario VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return usuariobd, false
	}

	//Insertamos
	_, err = stmtIns.Exec("DEFAULT", usuario.Nombre, usuario.Clavepubrsa, usuario.Claveprivrsa, clavelogin, salt)
	if err != nil {
		fmt.Println(err.Error())
		return usuariobd, false
	}

	defer stmtIns.Close()

	usuariobd = bd.getUsuarioByNombreBD(usuario.Nombre)

	return usuariobd, true
}

// Comprobamos un usuario con su nombre y clave cifrada
func (bd *BD) loginUsuarioBD(nombre string, clavehashlogin []byte) (Usuario, bool) {

	var usuario Usuario
	usuario.Nombre = nombre

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return usuario, false
	}
	defer db.Close()

	//Obtenemos los datos del usuario según su nombre
	rows, err := db.Query("SELECT id, clavepubrsa, claveprivrsa, clavelogin, salt FROM usuario WHERE nombre = '" + nombre + "'")
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return usuario, false
	}

	for rows.Next() {
		err = rows.Scan(&usuario.Id, &usuario.Clavepubrsa, &usuario.Claveprivrsa, &usuario.Clavelogin, &usuario.Salt)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return usuario, false
		}
	}

	//Comparamos las claves generadas con la Salt de la BD y las guardadas en la BD
	clavelogin := generarClaveLoginConSalt(clavehashlogin, usuario.Salt)

	if string(clavelogin) != string(usuario.Clavelogin) {
		return usuario, false
	}

	return usuario, true
}

//Obtenemos usuario según id usuario
func (bd *BD) getUsuarioByNombreBD(user string) Usuario {
	usuario := Usuario{}
	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT id, nombre, clavepubrsa, claveprivrsa, clavelogin, salt FROM usuario WHERE nombre = '" + user + "'")
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
	}
	for rows.Next() {
		err = rows.Scan(&usuario.Id, &usuario.Nombre, &usuario.Clavepubrsa, &usuario.Claveprivrsa, &usuario.Clavelogin, &usuario.Salt)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
		}
	}
	return usuario
}

//Obtenemos una instancia de usuario según id usuario
func (bd *BD) getUsuarioById(id int) Usuario {

	var usuario Usuario

	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return usuario
	}
	defer db.Close()

	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT id, nombre, clavepubrsa, claveprivrsa, clavelogin, salt FROM usuario WHERE id = " + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return usuario
	}

	for rows.Next() {
		err = rows.Scan(&usuario.Id, &usuario.Nombre, &usuario.Clavepubrsa, &usuario.Claveprivrsa, &usuario.Clavelogin, &usuario.Salt)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return usuario
		}
	}

	return usuario
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

//Obtenemos clave pub de usuario según id usuario
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

//Modificamos los datos de un usuario
func (bd *BD) modificarUsuarioBD(usuario Usuario) bool {

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()

	nombreu := bd.getNombreUsuario(usuario.Id)
	if nombreu == "" {
		return false
	}

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("UPDATE usuario set clavepubrsa=?, claveprivrsa=?, clavelogin=?, salt=? where id=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos crear el chat
	_, err = stmtIns.Exec(usuario.Clavepubrsa, usuario.Claveprivrsa, usuario.Clavelogin, usuario.Salt, usuario.Id)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

func (bd *BD) getClavesMensajes(usuario int) []string {
	claves := make([]string, 0, 1) //Los mensajes de un chat

	//Conexion BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer db.Close()

	//De el chat buscamos los datos de los mensajes de dicho chat
	rows, err := db.Query("SELECT `idusuario`,`idclavesmensajes` FROM `clavesusuario` WHERE `idusuario`= " + strconv.Itoa(usuario))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return nil
	}

	var clave string
	var id string

	for rows.Next() {
		//Obtenemos los datos del mensaje
		err = rows.Scan(&clave, &id)

		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return nil
		}

		//Guardamos el mensaje en el array de mensajes
		claves = append(claves, id+":"+clave)
	}

	return claves
}
