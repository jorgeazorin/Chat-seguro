/*
	Encarna Amorós Beneite
	MTIS Práctica 1: Publicador y subscriptor asíncrono (mayorista Go)
*/

package main

import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jjeffery/stomp"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

//Authentication
var username = "mtis"
var password = "mtis"

//Convert Bytes to String
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

//Search in the BD if there are stocks for the request of the store
func searchBD(elementoPedido string) (int, float32) {
	//Connection with the BD MTISMayorista
	db, err := sql.Open("mysql", "mtis:mtis@/MTISMayorista")
	if err != nil {
		//panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		return 0, 0
	}

	// we "scan" the result in here
	var cantidadReal int = 0
	var precio float32 = 0

	// Query the quantity (cantidad) of nombre = rueda
	rows, err := db.Query("SELECT cantidad, precio, nombre FROM almacen WHERE nombre = '" + elementoPedido + "'")
	if err != nil {
		//Closing connection
		defer db.Close()
		return 0, 0
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&cantidadReal, &precio, &name)
		if err != nil {
			//Closing connection
			defer db.Close()
			return 0, 0
		}
	}

	//Closing connection
	defer db.Close()

	return cantidadReal, precio
}

//Given a string with the data with an XML format obtain the element and the quantity
func readXML(solicitud string) (string, int) {
	var elementoPedido string
	var cantidadPedida int

	//If has a wrong format "" and 0
	if !(strings.Contains(solicitud, "<pieza>") && strings.Contains(solicitud, "</pieza>") &&
		strings.Contains(solicitud, "<cantidad>") && strings.Contains(solicitud, "</cantidad>")) {
		return "", 0
	}

	//Obtain indexs beginning and end to be able to obtain the substrings
	index_inicio_pieza := strings.Index(solicitud, "<pieza>") + len("<pieza>")
	index_fin_pieza := len(solicitud) - strings.Index(solicitud, "</pieza>")
	index_inicio_cantidad := strings.Index(solicitud, "<cantidad>") + len("<cantidad>")
	index_fin_cantidad := len(solicitud) - strings.Index(solicitud, "</cantidad>")

	//Obtain the substrings
	elementoPedido = solicitud[index_inicio_pieza : len(solicitud)-index_fin_pieza]
	cantidadPedidaString := solicitud[index_inicio_cantidad : len(solicitud)-index_fin_cantidad]

	cantidadPedida, err := strconv.Atoi(cantidadPedidaString)
	if err != nil {
		return "", 0
	}

	return elementoPedido, cantidadPedida
}

func responder(solicitud string) {
	var respuesta, disponibilidad string = "", "Si"
	var errorformato bool = false

	//Error XML format
	elementoPedido, cantidadPedida := readXML(solicitud)
	if elementoPedido == "" {
		respuesta = "<error>Go detecta Error formato XML</error>"
		errorformato = true
	}

	//Disponibility
	cantidadReal, precio := searchBD(elementoPedido)
	if cantidadReal < cantidadPedida || cantidadReal == 0 {
		disponibilidad = "No"
		precio = 0
	} else {
		precio = float32(cantidadPedida) * float32(precio)
	}

	if errorformato == false {
		respuesta = "<respuesta componente>"
		respuesta += "<pieza>" + elementoPedido + "</pieza>"
		respuesta += "<disponibilidad>" + disponibilidad + "</disponibilidad>"
		respuesta += "<precio>" + strconv.FormatFloat(float64(precio), 'f', 2, 64) + "</precio>"
		respuesta += "</respuesta componente>"
	}

	//Open connection
	conn, _ := stomp.Dial("tcp", "localhost:61613", stomp.Options{Login: username, Passcode: password})

	//Sending response
	conn.Send("/topic/Respuesta Disponibilidad", "", []byte(respuesta), nil)
	println("Enviamos respuesta a tienda: ", respuesta)

	conn.Disconnect()
}

func main() {

	//Open the connection
	conn, err := stomp.Dial("tcp", "127.0.0.1:61613", stomp.Options{Login: username, Passcode: password})
	if err != nil {
		println("Error con la conexión")
		println(err)
		return
	}

	//Asynchronous subscription
	for {
		//Subscription to the topic
		sub, suberror := conn.Subscribe("/topic/Consulta Disponibilidad", stomp.AckAuto)
		if suberror != nil {
			println("Error con la subscripción")
			println(err)
			conn.Disconnect()
			return
		}

		println("Esperando peticiones de tienda.")

		//When we have some request we respond
		msg := <-sub.C
		if msg.Err != nil {
			println("Error al recibir petición")
			return
		} else {
			stringPedido := BytesToString(msg.Body)
			println("Hemos recibido de tienda: ", stringPedido)
			responder(stringPedido)
		}
	}

	conn.Disconnect()
}
