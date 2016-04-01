package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Mensaje struct {
	From     string   `json:"From"`
	To       int      `json:"To"`
	Password string   `json:"Password"`
	Funcion  string   `json:"Funcion"`
	Datos    []string `json:"Datos"`
	Mensaje  string   `json:"MensajeSocket"`
}

func main() {
	//var window ui.Window
	//LEER CERTIFICADOS DE LOS ARCHIVOS (ESTOS SON LOS CERTIFICADOS DEL CLIENTE)
	cert2_b, _ := ioutil.ReadFile("cert2.pem")
	priv2_b, _ := ioutil.ReadFile("cert2.key")
	priv2, _ := x509.ParsePKCS1PrivateKey(priv2_b)

	//CONFIGURAR TLS CON LOS CERTIFICADOS
	cert := tls.Certificate{
		Certificate: [][]byte{cert2_b},
		PrivateKey:  priv2,
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	///////////////////////////////////
	//    Conectar    /////////////////
	//////////////////////////////////
	conn, err := tls.Dial("tcp", "127.0.0.1:443", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	///////////////////////////////////
	//    Login      /////////////////
	//////////////////////////////////
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("User:")
	message, _ := reader.ReadString('\n')
	message = message[0 : len(message)-2]
	mensaje := Mensaje{}
	mensaje.From = message
	mensaje.Funcion = "login"
	mensaje.To = -1
	//Convertir a json
	b, _ := json.Marshal(mensaje)
	log.Printf(string(b))
	//Escribe json en el socket
	conn.Write(b)

	///////////////////////////////////
	//    Enviar  y recibir      /////
	//////////////////////////////////
	//Por si envia algo el servidor
	go handleServerRead(conn)

	//Enviar mensajes
	go handleClientWrite(conn, mensaje.From)

	//Para que no se cierre la consola
	for {
	}
}

//Si envia algo el servidor a este cliente lo muestra en pantalla
func handleServerRead(conn net.Conn) {
	var mensaje Mensaje
	//bucle infinito
	for {
		defer conn.Close()
		reply := make([]byte, 256)
		n, err := conn.Read(reply)
		if err != nil {
			break
			conn.Close()
		}
		json.Unmarshal(reply[:n], &mensaje)
		fmt.Println("" + mensaje.From + " -> " + mensaje.Mensaje)

	}
}

//SI escribe algo lo envia al servidor
func handleClientWrite(conn net.Conn, from string) {
	mensaje := Mensaje{}

	//bucle infinito
	for {
		defer conn.Close()
		//Cuando escribe algo y le da a enter
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		//Rellenar datos
		mensaje.From = from
		mensaje.Password = "1"
		mensaje.Funcion = "enviar"
		mensaje.Mensaje = message[0 : len(message)-2]
		mensaje.To = 2

		//Convertir a json
		b, _ := json.Marshal(mensaje)

		//Escribe json en el socket
		conn.Write(b)
	}

}
