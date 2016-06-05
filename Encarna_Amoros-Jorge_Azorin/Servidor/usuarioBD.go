/////////
//USUARIO
/////////

package main

import (
	"crypto/rand"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
)

//Datos tabla usuario
type Usuario struct {
	Id           int    `json:"Id"`
	Nombre       string `json:"Nombre"`
	Estado       string `json:"Estado"`
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

//Insertamos a un nuevo usuario en BD
func (bd *BD) insertUsuarioBD(usuario Usuario) (Usuario, bool) {

	var usuariobd Usuario

	//genera la clavelogin con bcrypt (y el salt con el que se crea)
	usuario.Clavelogin, usuario.Salt = generarClaveLoginClaves(usuario.Clavelogin)

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
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
	if test == false {
		return Usuario{}, false
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT * FROM usuario WHERE nombre = ?", nombre)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return Usuario{}, false
	}

	//Comparamos las claves generadas con la Salt de la BD y las guardadas en la BD
	clavelogin := generarClaveLoginConSalt(clavehashlogin, usuario.Salt)

	if string(clavelogin) != string(usuario.Clavelogin) {
		return usuario, false
	}

	return usuario, true
}

//Devuelve la clave privada de un usuario
func (bd *BD) obtenerClavePrivadaUsuario(id int) (Usuario, bool) {
	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return Usuario{}, false
	}
	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT * FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return Usuario{}, false
	}

	return usuario, true
}

//Obtenemos una instancia de usuario según id usuario
func (bd *BD) getUsuarioById(id int) (Usuario, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return Usuario{}, false
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT * FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return Usuario{}, false
	}

	return usuario, true
}

//Obtenemos usuario según id usuario
func (bd *BD) getUsuarioByNombreBD(user string) (Usuario, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return Usuario{}, false
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT * FROM usuario WHERE nombre = ?", user)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return Usuario{}, false
	}

	return usuario, true
}

//Obtenemos nombre de usuario según id usuario
func (bd *BD) getNombreUsuario(id int) (string, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return "", false
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT nombre FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return "", false
	}

	return usuario.Nombre, true
}

//Obtenemos clave pub de usuario según id usuario
func (bd *BD) getClavePubUsuario(id int) ([]byte, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return []byte{}, false
	}

	//Select
	var usuario Usuario
	err := dbmap.SelectOne(&usuario, "SELECT clavepubrsa FROM usuario WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []byte{}, false
	}

	return usuario.Clavepubrsa, true
}

//Obtener los id de usuarios de un chat
func (bd *BD) getUsuariosChatBD(id int) ([]int, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return []int{}, false
	}

	//Select
	usuarios := []int{}
	_, err := dbmap.Select(&usuarios, "SELECT idusuario FROM usuarioschat WHERE idchat = ?", id)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []int{}, false
	}

	return usuarios, true
}

//Modificamos los datos de un usuario
func (bd *BD) modificarUsuarioNombreEstadoBD(usuario Usuario) bool {

	usuarionuevo, test := bd.getUsuarioById(usuario.Id)
	if test == false {
		return false
	}
	usuarionuevo.Nombre = usuario.Nombre
	usuarionuevo.Estado = usuario.Estado

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return false
	}

	//Update
	_, err := dbmap.Update(&usuarionuevo)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return false
	}

	return true
}

//Obtiene las claves de los mensajes
func (bd *BD) getClavesMensajes(usuario int) ([]string, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return []string{}, false
	}

	//Select
	claves := make([]string, 0, 1) //Los mensajes de un chat
	_, err := dbmap.Select(&claves, "SELECT idclavesmensajes FROM clavesusuario WHERE idusuario = ?", usuario)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []string{}, false
	}

	return claves, true
}

//Obtenemos una instancia de usuario según id usuario
func (bd *BD) getUsuarios() ([]Usuario, bool) {

	//Conexion y dbmapa
	dbmap, db, test := bd.conectarBD()
	defer db.Close()
	if test == false {
		return []Usuario{}, false
	}

	//Select
	var usuarios []Usuario
	_, err := dbmap.Select(&usuarios, "SELECT * FROM usuario")
	if err != nil {
		fmt.Println("Error:", err.Error())
		return []Usuario{}, false
	}

	return usuarios, true
}
