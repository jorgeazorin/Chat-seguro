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
	"gopkg.in/gorp.v1"
	"io"
	"log"
	"strconv"
)

type Usuario struct {
	Id           int    `json:"Id"`
	Nombre       string `json:"Nombre"`
	Clavepubrsa  []byte `json:"Clavepubrsa"`
	Claveprivrsa []byte `json:"Claveprivrsa"`
	Clavelogin   []byte `json:"Clavelogin"`
	Salt         []byte `json:"Salt"`
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

//Conexión con BD y mapa para gorp
func (bd *BD) conectarBD() (*gorp.DbMap, *sql.DB, bool) {

	var dbmap *gorp.DbMap
	var db *sql.DB
	var err error

	//Conexión BD
	db, err = sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return dbmap, db, true
	}

	//Construye un mapa gorp DbMap
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	//Añade la tabla especificando el nombre, con true el id automático
	dbmap.AddTableWithName(Usuario{}, "usuario").SetKeys(true, "Id")

	return dbmap, db, false
}

//Insertamos a un nuevo usuario en BD
func (bd *BD) insertUsuarioBD(usuario Usuario) (Usuario, bool) {

	var usuariobd Usuario

	//genera la clavelogin con scrypt (y el salt con el que se crea)
	usuario.Clavelogin, usuario.Salt = generarClaveLoginClaves(usuario.Clavelogin)

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return usuariobd, false
	}

	//Insert
	err := dbmap.Insert(&usuario)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return usuariobd, false
	}

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

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return ""
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT nombre FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return ""
	}

	return usuario.Nombre
}

//Obtenemos clave pub de usuario según id usuario
func (bd *BD) getClavePubUsuario(id int) []byte {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return []byte{}
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT clavepubrsa FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []byte{}
	}

	return usuario.Clavepubrsa
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

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return false
	}

	//Update
	_, err := dbmap.Update(&usuario)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return false
	}

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
