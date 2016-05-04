package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
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
	From        string   `json:"From"`
	To          int      `json:"To"`
	Password    string   `json:"Password"`
	Funcion     string   `json:"Funcion"`
	Datos       []string `json:"Datos"`
	DatosClaves [][]byte `json:"DatosClaves"`
	Chat        int      `json:"Chat"`
	Mensaje     string   `json:"MensajeSocket"`
}

var nombre_usuario_from string

//Para pasar los datos de un usuario
type Usuario struct {
	id           int
	nombre       string
	clavepubrsa  string
	claveprivrsa string
	clavelogin   []byte
	salt         []byte
	clavecifrado []byte
}

var ClientUsuario Usuario

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
	//obtenerMensajesChat(conn, 1)

	//Usuario 1 en el chat 7 al usuario 15
	//agregarUsuariosChat(conn, 7, []string{"15"})
	//Usuario 1 en el chat 7 al usuario 15
	//eliminarUsuariosChat(conn, 7, []string{"15"})

	//getClavePubUsuario(conn, 1)
	//getClaveMensaje(conn, 2)
	//getClaveCifrarMensajeChat(conn, 1)

	//CrearNuevaClaveMensajes(conn)
	//nuevaClaveUsuarioConIdConjuntoClaves(conn, 1, "minuevaclave1")
	var u Usuario
	u.nombre = "Prueba"
	u.clavepubrsa = "Prueba"
	u.claveprivrsa = "Prueba"
	//registrarUsuario(conn, u, "miclave1")

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
		reply := make([]byte, 1048576) //256
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

		//Si nos devuelven el usuario lo rellenamos
		if mensaje.Funcion == "DatosUsuario" {
			fmt.Println("cliente guardado:")
			idusuario, _ := strconv.Atoi(mensaje.Datos[0])
			ClientUsuario.id = idusuario
			ClientUsuario.nombre = mensaje.Datos[1]
			ClientUsuario.clavepubrsa = mensaje.Datos[2]
			ClientUsuario.claveprivrsa = mensaje.Datos[3]
			fmt.Println(ClientUsuario)
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

//De la ontraseña en claro se realiza hash y se divide en 2 (login y cifrado)
func generarHashClaves(clave string) {
	//Hash con SHA-2 (256) para la contraseña en general
	clavebytes := []byte(clave)
	clavebytesconsha2 := sha256.Sum256(clavebytes)

	//Dividimos dicho HASH
	clavehashlogin := clavebytesconsha2[0 : len(clavebytesconsha2)/2]
	clavehashcifrado := clavebytesconsha2[len(clavebytesconsha2)/2 : len(clavebytesconsha2)]

	ClientUsuario.clavecifrado = clavehashcifrado
	ClientUsuario.clavelogin = clavehashlogin
}

//Registrar a un usuario
func registrarUsuario(conn net.Conn, usuario Usuario, clave string) {

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.From = nombre_usuario_from
	mensaje.Funcion = "registrarusuario"

	//Generamos los hash de las claves
	generarHashClaves(clave)

	mensaje.Datos = []string{usuario.nombre, usuario.clavepubrsa, usuario.claveprivrsa}
	mensaje.DatosClaves = [][]byte{ClientUsuario.clavelogin}

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
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

	//Generamos los hash de las claves
	generarHashClaves(password)

	mensaje := Mensaje{}
	mensaje.From = nombreusuario
	mensaje.DatosClaves = [][]byte{ClientUsuario.clavelogin}
	mensaje.Funcion = "login"
	mensaje.To = -1

	//Rellenamos variable nombre usuario global
	nombre_usuario_from = nombreusuario
	ClientUsuario.nombre = nombreusuario
	//ClientUsuario.clavelogin = clavebcryptlogin
	//ClientUsuario.clavecifrado = clavehashcifrado
	//ClientUsuario.salt = salt

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
	mensaje.Funcion = "getclavesmensajes"
	cadena_idmensaje := strconv.Itoa(idmensaje)
	mensaje.Datos = []string{cadena_idmensaje}

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}

//Cliente pide clave cifrada para descifrar mensajes
func getClaveCifrarMensajeChat(conn net.Conn, idchat int) {

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.From = nombre_usuario_from
	mensaje.Password = "1"
	mensaje.Funcion = "getclavecifrarmensajechat"
	cadena_idchat := strconv.Itoa(idchat)
	mensaje.Datos = []string{cadena_idchat}

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}

//Cliente crea nuevo id clave para un nuevo conjunto de claves
func CrearNuevaClaveMensajes(conn net.Conn) {

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.From = nombre_usuario_from
	mensaje.Password = "1"
	mensaje.Funcion = "crearnuevoidparanuevaclavemensajes"

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}

func cifrarAES(datos string, clave []byte) (cipher.Block, bool) {

	fmt.Println("mira la clave:", clave)
	fmt.Println("mira la :", ClientUsuario.clavecifrado)
	fmt.Println("mira la :", ClientUsuario.nombre)

	var nulo cipher.Block

	clavecifrada, err := aes.NewCipher(clave)

	if err != nil {
		return nulo, true
	}

	fmt.Println("mira la clave:", clavecifrada)

	return clavecifrada, false
}

func descifrarAES(datos cipher.Block, clave []byte) {

}

//Asocia nueva clave de un usuario con el id que indica ese nuevo conjunto de claves
func nuevaClaveUsuarioConIdConjuntoClaves(conn net.Conn, idconjuntoclaves int, clavemensajes string) {

	//Cifrar la clave para los mensajes
	clavecifradamensajes, err := cifrarAES(clavemensajes, ClientUsuario.clavecifrado)

	if err == true {
		fmt.Println("Error al generar clave con cifrado AES")
		return
	}

	fmt.Println("Miraa con aes:", clavecifradamensajes)

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.From = nombre_usuario_from
	mensaje.Password = "1"
	mensaje.Funcion = "nuevaclaveusuarioconidconjuntoclaves"
	cadena_idconjuntoclaves := strconv.Itoa(idconjuntoclaves)
	mensaje.Datos = []string{cadena_idconjuntoclaves}
	//mensaje.DatosClaves = [][]byte{clavecifradamensajes}

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}
