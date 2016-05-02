/////////
//USUARIO
/////////

package main

import (
	//"crypto/sha256"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"golang.org/x/crypto/scrypt"
	"strconv"
	//"strings"
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

//Funcion para obtener los datos del usuario cuando se loguea
func (usuario *Usuario) login(nombre string, password string) bool {

	//Para las operaciones con la BD
	var bd BD
	bd.username = "root"
	bd.password = ""
	bd.adress = ""
	bd.database = "sds"

	usuario.Nombre = nombre
	user, test := bd.comprobarUsuarioBD(nombre, password)

	if test == false {
		return false
	}

	usuario.Id = user.Id
	usuario.Nombre = user.Nombre
	usuario.Clavepubrsa = user.Clavepubrsa
	usuario.Claveprivrsa = user.Claveprivrsa
	usuario.Clavelogin = user.Clavelogin
	usuario.Nombre = user.Nombre
	return true
}

// Comprobamos un usuario con su nombre y clave cifrada
func (bd *BD) comprobarUsuarioBD(nombre string, clave string) (Usuario, bool) {

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
	rows, err := db.Query("SELECT id, clavepubrsa, claveprivrsa, clavelogin, salt, clavecifrado FROM usuario WHERE nombre = '" + nombre + "'")
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return usuario, false
	}

	var clavelogin []byte
	var salt []byte
	var clavecifrado []byte

	for rows.Next() {
		err = rows.Scan(&usuario.Id, &usuario.Clavepubrsa, &usuario.Claveprivrsa, &clavelogin, &salt, &clavecifrado)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return usuario, false
		}
	}

	fmt.Println("MIRA:", clavecifrado)

	//Obtenemos los datos del usuario según su nombre
	rows, err = db.Query("SELECT id, clavepubrsa, claveprivrsa, clavelogin, salt, clavecifrado FROM usuario WHERE nombre = '" + nombre + "'")
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return usuario, false
	}

	for rows.Next() {
		err = rows.Scan(&usuario.Id, &usuario.Clavepubrsa, &usuario.Claveprivrsa, &usuario.Clavelogin, &usuario.Salt, &usuario.Clavecifrado)
		if err != nil {
			fmt.Println(err.Error())
			defer db.Close()
			return usuario, false
		}
	}

	/*/Comprobamos las claves

	//Hash con SHA-2 (256) para la contraseña en general
	clavebytes := []byte(clave)
	clavebytesconsha2 := sha256.Sum256(clavebytes)

	//fmt.Println("Mira la clave:", clavebytesconsha2)
	claveconsha2 := string(clavebytesconsha2[:])

	//Dividimos dicho HASH
	clavehashlogin := string(claveconsha2[0 : len(claveconsha2)/2])
	clavehashcifrado := string(claveconsha2[len(claveconsha2)/2 : len(claveconsha2)])

	fmt.Println("Mira la clave:", claveconsha2)
	fmt.Println("Mira la clave1:", clavehashlogin)
	fmt.Println("Mira la clave2:", clavehashcifrado)

	//Si la clave cifrado no coincide ya es falso
	if strings.Compare(clavehashcifrado, usuario.Clavecifrado) != 0 {
		return usuario, false
	}

	//La parte del login hacemos BCRYPT con SALT de la BD
	clavebcryptlogin, _ := scrypt.Key([]byte(clavehashlogin), []byte(usuario.Salt), 16384, 8, 1, 64)

	if strings.Compare(string(clavebcryptlogin), usuario.Clavelogin) != 0 {
		fmt.Println("Mal2", string(clavebcryptlogin), usuario.Clavelogin)
		return usuario, false
	}*/

	return usuario, true
}

//Obtenemos usuario según id usuario
func (bd *BD) getUsuarioBD(user string) Usuario {
	usuario := Usuario{}
	//Conexión BD
	db, err := sql.Open("mysql", bd.username+":"+bd.password+"@/"+bd.database)

	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	//Obtenemos el nombre del usuario
	rows, err := db.Query("SELECT id, nombre, clavepubrsa, claveprivrsa, clavelogin FROM usuario WHERE nombre = '" + user + "'")
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
	}
	for rows.Next() {
		err = rows.Scan(&usuario.Id, &usuario.Nombre, &usuario.Clavepubrsa, &usuario.Claveprivrsa, &usuario.Clavelogin)
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
	rows, err := db.Query("SELECT id, nombre, clavepubrsa, claveprivrsa, clavelogin FROM usuario WHERE id = " + strconv.Itoa(id))
	if err != nil {
		fmt.Println(err.Error())
		defer db.Close()
		return usuario
	}

	for rows.Next() {
		err = rows.Scan(&usuario.Id, &usuario.Nombre, &usuario.Clavepubrsa, &usuario.Claveprivrsa, &usuario.Clavelogin)
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

	nombreu := bd.getNombreUsuario(usuario.Id)
	if nombreu == "" {
		return false
	}

	//Preparamos crear el chat
	stmtIns, err := db.Prepare("UPDATE usuario set clavepubrsa=?, claveprivrsa=?, clavelogin=? where id=?")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos crear el chat
	_, err = stmtIns.Exec(usuario.Clavepubrsa, usuario.Claveprivrsa, usuario.Clavelogin, usuario.Id)
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
	stmtIns, err := db.Prepare("INSERT INTO usuario VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	//Insertamos
	_, err = stmtIns.Exec("DEFAULT", usuario.Nombre, usuario.Clavepubrsa, usuario.Claveprivrsa, usuario.Clavelogin, usuario.Salt, usuario.Clavecifrado)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmtIns.Close()

	return true
}

func (bd *BD) getClaves(usuario int) []string {
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
