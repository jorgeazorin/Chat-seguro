//////
//MAIN
//////

package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var username = "sds"
var password = "sds"
var adress = ""
var database = "sds"

func main() {
	//Conexión BD
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)

	if err != nil {
		panic(err.Error())
		return false
	}
	defer db.Close()
}
