//////
//BD
//////

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

//Datos para la conexión con la BD
type BD struct {
	username string
	password string
	adress   string
	database string
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
		return dbmap, db, false
	}

	//Construye un mapa gorp DbMap
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	//Añade la tabla especificando el nombre, con true el id automático o no
	dbmap.AddTableWithName(Usuario{}, "usuario").SetKeys(true, "Id")
	dbmap.AddTableWithName(Mensaje{}, "mensaje").SetKeys(true, "Id")
	dbmap.AddTableWithName(Receptoresmensaje{}, "receptoresmensaje")
	dbmap.AddTableWithName(Clavesmensajes{}, "clavesmensajes").SetKeys(true, "Id")
	dbmap.AddTableWithName(Clavesusuario{}, "clavesusuario")
	dbmap.AddTableWithName(Chat{}, "chat").SetKeys(true, "Id")
	dbmap.AddTableWithName(UsuariosChat{}, "usuarioschat")

	return dbmap, db, true
}
