package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

//Struct de los mensajes que se envian por el socket
type Mensaje struct {
	From        string   `json:"From"`
	Idfrom      int      `json:"Idfrom"`
	To          int      `json:"To"`
	Password    string   `json:"Password"`
	Funcion     string   `json:"Funcion"`
	Datos       []string `json:"Datos"`
	DatosClaves [][]byte `json:"DatosClaves"`
	Chat        int      `json:"Chat"`
	Mensaje     string   `json:"MensajeSocket"`
	Mensajechat []byte   `json:"Mensajechat"`
}

var nombre_usuario_from string

//Para pasar los datos de un usuario
type Usuario struct {
	id               int
	nombre           string
	clavepubrsa      []byte
	claveprivrsa     []byte
	claveenclaro     string
	clavehashcifrado []byte
	clavehashlogin   []byte
}

var ClientUsuario Usuario

//Variable global, de momento para guardar clave cifrar mensajes chat1
var clavecifrarmensajes []byte

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

	//probando a cifrar descifrar con AES
	/*mensajecifrado, _ := cifrarAES([]byte("hola amigos"), ClientUsuario.clavehashcifrado)
	mensajedescifrado, _ := descifrarAES(mensajecifrado, ClientUsuario.clavehashcifrado)
	fmt.Println("Tu mensaje era:", string(mensajedescifrado))*/

	//obtenerMensajesChat(conn, 1)

	//Usuario 1 en el chat 7 al usuario 15
	//agregarUsuariosChat(conn, 7, []string{"15"})
	//Usuario 1 en el chat 7 al usuario 15
	//eliminarUsuariosChat(conn, 7, []string{"15"})

	//getClavePubUsuario(conn, 1)
	//getClaveMensaje(conn, 2)
	//getClaveCifrarMensajeChat(conn, 1)

	//CrearNuevaClaveMensajes(conn)
	//nuevaClaveUsuarioConIdConjuntoClaves(conn, 1, "nuevaclave1")

	//Registrar un usuario
	//ClientUsuario.nombre = "Prueba3"
	//ClientUsuario.claveenclaro = "miclave3"
	//ClientUsuario.clavepubrsa, ClientUsuario.claveprivrsa = generarClavesRSA()
	//registrarUsuario(conn)

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

		//Descifrar el mensaje
		if len(mensaje.Mensajechat) != 0 {
			mensajedescifrado, err2 := descifrarAES([]byte(mensaje.Mensajechat), clavecifrarmensajes)
			if err2 == true {
				fmt.Println("Error al descifrar clave con cifrado AES")
				continue
			}
			mensaje.Mensajechat = mensajedescifrado
		}

		//Mostramos mensaje por pantalla
		fmt.Println("" + mensaje.From + " -> " + mensaje.Mensaje + " -> " + string(mensaje.Mensajechat) + " Datos: ->")
		for i := 0; i < len(mensaje.Datos); i++ {
			fmt.Println("dato:", i, "->", mensaje.Datos[i])
		}
		/*De momento no mostramos que es mucho dato grande
		for i := 0; i < len(mensaje.DatosClaves); i++ {
			fmt.Println("dato clave:", i, "->", mensaje.DatosClaves[i])
		}*/
		fmt.Println()

		//Si nos devuelven el usuario lo rellenamos. (menos claveenclaro, clavehashcifrado clavehashlogin ya estan rellenos)
		if mensaje.Funcion == "DatosUsuario" {
			idusuario, _ := strconv.Atoi(mensaje.Datos[0])
			ClientUsuario.id = idusuario
			ClientUsuario.nombre = mensaje.Datos[1]
			ClientUsuario.clavepubrsa = mensaje.DatosClaves[0]
			ClientUsuario.claveprivrsa = mensaje.DatosClaves[1]

			////////////////////////
			//PRUEBAS AFTER DE LOGIN
			////////////////////////
			//nuevaClaveUsuarioConIdConjuntoClaves(conn, 1, "nuevaclave1")
			getClaveCifrarMensajeChat(conn, 1)
		}

		if mensaje.Funcion == "DatosClaveCifrarMensajeChat" {
			laclave := mensaje.DatosClaves[0]
			clavecifrarmensajes, _ = descifrarAES(laclave, ClientUsuario.clavehashcifrado)
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
		mensaje.From = ClientUsuario.nombre
		mensaje.Password = "1"
		mensaje.Funcion = "enviar"
		mensaje.Mensajechat = []byte(message[0 : len(message)-2])

		//getClaveCifrarMensajeChat(conn, 1)

		//Ciframos el mensaje
		//De momento clave guardada globalmente pero sería hacer llamadas al servidor, channels, etc
		mensajecifrado, err := cifrarAES([]byte(mensaje.Mensajechat), clavecifrarmensajes)
		if err == true {
			fmt.Println("Error al cifrar clave con cifrado AES")
			continue
		}
		mensaje.Mensajechat = mensajecifrado

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
func generarHashClaves(clave string) ([]byte, []byte) {

	//Hash con SHA-2 (256) para la contraseña en general
	clavebytes := []byte(clave)
	clavebytesconsha2 := sha256.Sum256(clavebytes)

	//Dividimos dicho HASH
	clavehashlogin := clavebytesconsha2[0 : len(clavebytesconsha2)/2]
	clavehashcifrado := clavebytesconsha2[len(clavebytesconsha2)/2 : len(clavebytesconsha2)]

	return clavehashlogin, clavehashcifrado
}

//Genera una clave pública y otra privada
func generarClavesRSA() ([]byte, []byte) {
	claveprivada, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
	}

	clavepublica := &claveprivada.PublicKey

	priv := x509.MarshalPKCS1PrivateKey(claveprivada)
	pub, _ := x509.MarshalPKIXPublicKey(clavepublica)

	return pub, priv
}

//Cifrar con AES en modo CTR
func cifrarAES(textocifrar []byte, clave []byte) ([]byte, bool) {

	//Calculamos block con clave
	block, err := aes.NewCipher(clave)
	if err != nil {
		fmt.Println(err)
		return []byte{}, true
	}

	// IV necesita ser único aunque no seguro, se incluye al principio del textocifrado
	ciphertext := make([]byte, aes.BlockSize+len(textocifrar))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, true
	}

	//Ciframos
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], textocifrar)

	return ciphertext, false
}

//Con CTR se descifra como se cifra, con NewCTR
func descifrarAES(ciphertext []byte, clave []byte) ([]byte, bool) {

	//Calculamos block con clave
	block, err := aes.NewCipher(clave)
	if err != nil {
		fmt.Println(err)
		return []byte{}, true
	}

	//Volvemos a calcular iv (ahora sin rand, iv está al principio del textocifrado)
	iv := ciphertext[:aes.BlockSize]

	//Desciframos
	textodescifrado := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(textodescifrado, ciphertext[aes.BlockSize:])

	return textodescifrado, false
}

//Registrar a un usuario
func registrarUsuario(conn net.Conn) {

	var err bool

	//Rellenar datos del mensaje
	mensaje := Mensaje{}
	mensaje.From = ClientUsuario.nombre
	mensaje.Funcion = "registrarusuario"

	//Generamos los hash de las claves
	ClientUsuario.clavehashlogin, ClientUsuario.clavehashcifrado = generarHashClaves(ClientUsuario.claveenclaro)

	//Clave privada del usuario cifrar
	ClientUsuario.claveprivrsa, err = cifrarAES(ClientUsuario.claveprivrsa, ClientUsuario.clavehashcifrado)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return
	}

	mensaje.Datos = []string{ClientUsuario.nombre}
	mensaje.DatosClaves = [][]byte{ClientUsuario.clavehashlogin, ClientUsuario.clavepubrsa, ClientUsuario.claveprivrsa}

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
	ClientUsuario.nombre, _ = reader.ReadString('\n')
	ClientUsuario.nombre = ClientUsuario.nombre[0 : len(ClientUsuario.nombre)-2]

	fmt.Print("Password:")
	ClientUsuario.claveenclaro, _ = reader.ReadString('\n')
	ClientUsuario.claveenclaro = ClientUsuario.claveenclaro[0 : len(ClientUsuario.claveenclaro)-2]

	//Generamos los hash de las claves
	ClientUsuario.clavehashlogin, ClientUsuario.clavehashcifrado = generarHashClaves(ClientUsuario.claveenclaro)

	mensaje := Mensaje{}
	mensaje.From = ClientUsuario.nombre
	mensaje.DatosClaves = [][]byte{ClientUsuario.clavehashlogin}
	mensaje.Funcion = "login"
	mensaje.To = -1

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
	mensaje.From = ClientUsuario.nombre
	mensaje.Password = "1"
	mensaje.Funcion = "obtenermensajeschat"

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
	mensaje.From = ClientUsuario.nombre
	mensaje.Password = "1"
	mensaje.Funcion = "agregarusuarioschat"
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
	mensaje.From = ClientUsuario.nombre
	mensaje.Password = "1"
	mensaje.Funcion = "eliminarusuarioschat"
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
	mensaje.From = ClientUsuario.nombre
	mensaje.Password = "1"
	mensaje.Funcion = "getclavepubusuario"

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
	mensaje.From = ClientUsuario.nombre
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
	mensaje.From = ClientUsuario.nombre
	mensaje.Idfrom = ClientUsuario.id
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
	mensaje.From = ClientUsuario.nombre
	mensaje.Password = "1"
	mensaje.Funcion = "crearnuevoidparanuevaclavemensajes"

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}

//Asocia nueva clave de un usuario con el id que indica ese nuevo conjunto de claves
func nuevaClaveUsuarioConIdConjuntoClaves(conn net.Conn, idconjuntoclaves int, clavemensajes string) {

	clavesmensajeshash, _ := generarHashClaves(clavemensajes)

	//Cifrar la clave para los mensajes
	clavecifradamensajes, err := cifrarAES(clavesmensajeshash, ClientUsuario.clavehashcifrado)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return
	}

	mensaje := Mensaje{}

	//Rellenar datos
	mensaje.From = ClientUsuario.nombre
	mensaje.Idfrom = ClientUsuario.id
	mensaje.Chat = 1
	mensaje.Funcion = "nuevaclaveusuarioconidconjuntoclaves"
	cadena_idconjuntoclaves := strconv.Itoa(idconjuntoclaves)
	mensaje.Datos = []string{cadena_idconjuntoclaves}
	mensaje.DatosClaves = [][]byte{clavecifradamensajes}

	//Convertir a json
	b, _ := json.Marshal(mensaje)

	log.Printf(string(b))

	//Escribe json en el socket
	conn.Write(b)
}
