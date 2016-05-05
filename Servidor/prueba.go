package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

type Usuario struct {
	Id           int    `json:"Id"`
	Nombre       string `json:"Nombre"`
	Clavepubrsa  []byte `json:"Clavepubrsa"`
	Claveprivrsa []byte `json:"Claveprivrsa"`
	Clavelogin   []byte `json:"Clavelogin"`
	Salt         []byte `json:"Salt"`
	Clavecifrado []byte `json:"Clavecifrado"`
}

func main() {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	username := "root"
	password := ""
	database := "sds"
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	var u Usuario
	u.Id = 6
	u.Nombre = "recorcholes"
	u.Clavepubrsa = []byte("recorcholes1")
	u.Claveprivrsa = []byte("recorchole2")
	u.Clavelogin = []byte("recorchole3")
	u.Clavelogin = []byte("recorchole4")

	//Crear la tabla
	//dbmap.AddTableWithName(u, "usuario").SetKeys(true, "Id")

	//insert
	err = dbmap.Insert(u)

	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}
