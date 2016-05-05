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

	usuariobd, _ = bd.getUsuarioByNombreBD(usuario.Nombre)

	return usuariobd, true
}

// Comprobamos un usuario con su nombre y clave cifrada
func (bd *BD) loginUsuarioBD(nombre string, clavehashlogin []byte) (Usuario, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return Usuario{}, true
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT * FROM usuario WHERE nombre = ?", nombre)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return Usuario{}, true
	}

	//Comparamos las claves generadas con la Salt de la BD y las guardadas en la BD
	clavelogin := generarClaveLoginConSalt(clavehashlogin, usuario.Salt)

	if string(clavelogin) != string(usuario.Clavelogin) {
		return usuario, false
	}

	return usuario, true
}

//Obtenemos una instancia de usuario según id usuario
func (bd *BD) getUsuarioById(id int) (Usuario, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return Usuario{}, true
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT * FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return Usuario{}, true
	}

	return usuario, false
}

//Obtenemos usuario según id usuario
func (bd *BD) getUsuarioByNombreBD(user string) (Usuario, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return Usuario{}, true
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT * FROM usuario WHERE nombre = ?", user)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return Usuario{}, true
	}

	return usuario, false
}

//Obtenemos nombre de usuario según id usuario
func (bd *BD) getNombreUsuario(id int) (string, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return "", true
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT nombre FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return "", true
	}

	return usuario.Nombre, false
}

//Obtenemos clave pub de usuario según id usuario
func (bd *BD) getClavePubUsuario(id int) ([]byte, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return []byte{}, true
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT clavepubrsa FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []byte{}, true
	}

	return usuario.Clavepubrsa, false
}

//Obtener los id de usuarios de un chat
func (bd *BD) getUsuariosChatBD(id int) ([]int, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return []int{}, true
	}

	//Select
	usuarios := []int{}
	_, err := dbmap.Select(&usuarios, "SELECT idusuario FROM usuarioschat WHERE idchat = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []int{}, true
	}

	return usuarios, false
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

func (bd *BD) getClavesMensajes(usuario int) ([]string, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == true {
		return []string{}, true
	}

	//Select
	claves := make([]string, 0, 1) //Los mensajes de un chat
	_, err := dbmap.Select(&claves, "SELECT idclavesmensajes FROM clavesusuario WHERE idusuario = ?", usuario)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []string{}, true
	}

	return claves, false
}
