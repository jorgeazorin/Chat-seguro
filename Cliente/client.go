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
	"strconv"
)

//Struct de los mensajes que se envian por el socket
type Mensaje struct {
	From     string   `json:"From"`
	To       int      `json:"To"`
	Password string   `json:"Password"`
	Funcion  string   `json:"Funcion"`
	Datos    []string `json:"Datos"`
	Chat     int      `json:"Chat"`
	Mensaje  string   `json:"MensajeSocket"`
}

var nombre_usuario_from string

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
	conn, err := tls.Dial("tcp", "127.0.0.1:444", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	//Por si envia algo el servidor
	go handleServerRead(conn)

	///////////////////////////////////
	//    PRUEBAS    /////////////////
	//////////////////////////////////
	login(conn)
	obtenerMensajesChat(conn, 1)
	//Usuario 1 en el chat 7 al usuario 15
	//agregarUsuariosChat(conn, 7, []string{"15"})
	//Usuario 1 en el chat 7 al usuario 15
	//eliminarUsuariosChat(conn, 7, []string{"15"})
	getClavePubUsuario(conn, 1)
	getClaveMensaje(conn, 2)

	///////////////////////////////////
	//    Enviar  y recibir      /////
	//////////////////////////////////

	//Enviar mensajes
	go handleClientWrite(conn) //	go handleClientWrite(conn, mensaje.From)

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
		reply := make([]byte, 524288) //256
		n, err := conn.Read(reply)
		if err != nil {
			break
			conn.Close()
		}
		json.Unmarshal(reply[:n], &mensaje)

		fmt.Println("" + mensaje.From + " -> " + mensaje.Mensaje + " Datos: ->")

		for i := 0; i < len(mensaje.Datos); i++ {
			fmt.Println("dato:", i, "->", mensaje.Datos[i])
		}

	}
}

//SI escribe algo lo envia al servidor
func handleClientWrite(conn net.Conn) {
	mensaje := Mensaje{}

	//bucle infinito
	for {
		defer conn.Close()

		//Cuando escribe algo y le da a enter
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		//Rellenar datos
		mensaje.From = nombre_usuario_from
		mensaje.Password = "1"
		mensaje.Funcion = "enviar"
		mensaje.Mensaje = message[0 : len(message)-2]
		mensaje.To = 2
		datos := []string{""}
		mensaje.Datos = datos
		mensaje.Chat = 1

		//Convertir a json
		b, _ := json.Marshal(mensaje)

		//Escribe json en el socket
		conn.Write(b)
	}

}

//Cliente realiza login
func login(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	//Pedimos los datos
	fmt.Print("Usuario:")
	nombreusuario, _ := reader.ReadString('\n')
	nombreusuario = nombreusuario[0 : len(nombreusuario)-2]

	fmt.Print("Password:")
	password, _ := reader.ReadString('\n')
	password = password[0 : len(password)-2]

	mensaje := Mensaje{}
	mensaje.From = nombreusuario
	mensaje.Password = password
	mensaje.Funcion = "login"
	mensaje.To = -1

	//Rellenamos variable nombre usuario global
	nombre_usuario_from = nombreusuario

	//Convertir a json
	b, _ := json.Marshal(mensaje)
	log.Printf(string(b))

	//Escribe peticion json en el socket
	conn.Write(b)
}

//Cliente pide mensajes de un chat
func obtenerMensajesChat(conn net.Conn, idchat int) {

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.Chat = idchat
	mensaje.From = nombre_usuario_from
	mensaje.Password = "1"
	mensaje.Funcion = "obtenermensajeschat"
	mensaje.Mensaje = ""

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)

}

//Cliente pide añadir usuarios a un chat
func agregarUsuariosChat(conn net.Conn, idchat int, usuarios []string) {

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.Chat = idchat
	mensaje.From = nombre_usuario_from
	mensaje.Password = "1"
	mensaje.Funcion = "agregarusuarioschat"
	mensaje.Mensaje = ""
	mensaje.Datos = usuarios

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}

//Cliente pide eliminar usuarios en un chat
func eliminarUsuariosChat(conn net.Conn, idchat int, usuarios []string) {

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.Chat = idchat
	mensaje.From = nombre_usuario_from
	mensaje.Password = "1"
	mensaje.Funcion = "eliminarusuarioschat"
	mensaje.Mensaje = ""
	mensaje.Datos = usuarios

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}

//Cliente pide clave pública de un usuario
func getClavePubUsuario(conn net.Conn, idusuario int) {

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.From = nombre_usuario_from
	mensaje.Password = "1"
	mensaje.Funcion = "getclavepubusuario"
	mensaje.Mensaje = ""
	cadena_idusuario := strconv.Itoa(idusuario)
	mensaje.Datos = []string{cadena_idusuario}

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}

//Cliente pide clave cifrada para descifrar mensajes
func getClaveMensaje(conn net.Conn, idmensaje int) {

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.From = nombre_usuario_from
	mensaje.Password = "1"
	mensaje.Funcion = "getclavemensaje"
	cadena_idmensaje := strconv.Itoa(idmensaje)
	mensaje.Datos = []string{cadena_idmensaje}

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}
