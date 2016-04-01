package main

import (
	"encoding/json"
	"log"
	"net"
)

//Struct de conexión tiene el id del usuario y la conexión
type Conexion struct {
	conexiones *Conexiones //Esto es el vector con todos los sockets que hay online
	conexion   net.Conn    //la propia conexión
	usuario    Usuario     //el usuario que habla por el socket se rellena al hacer el login()
}

//Struct de los mensajes que se envian por el socket
type MensajeSocket struct {
	From          string   `json:"From"`
	To            int      `json:"To"`
	Password      string   `json:"Password"`
	Funcion       string   `json:"Funcion"`
	Datos         []string `json:"Datos"`
	MensajeSocket string   `json:"MensajeSocket"`
}

//Función que se encarga de leer un socket infinitamente
func (conexion *Conexion) escuchar() {
	defer conexion.Close()
	var mensaje MensajeSocket //Struct donde se guarda el mensaje que se descodifia

	for { // Bucle infinito que lee cosas que envia el usuario
		buf := make([]byte, 256)
		n, err := conn.Read(buf) //Lee el mensaje
		if err != nil {
			break
			conn.Close()
		}
		json.Unmarshal(buf[:n], &mensaje)       //Descodificar el mensaje recibido (estaba en json y se pasa a struct)
		conexion.ProcesarMensajeSocket(mensaje) //Procesa el mensaje, esto lo hace en el archivo router.go
	}
}

//Función que envia un mensaje a un cliente mediante un id y un string
func (conexion *Conexion) EnviarMensajeSocketSocket(s MensajeSocket) {

	b, _ := json.Marshal(s)              //Codifica el mensaje en json
	_, err := conexion.conexion.Write(b) //lo escribe en el socket
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

}
