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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	Id               int    `json:"Id"`
	Nombre           string `json:"Nombre"`
	Clavepubrsa      []byte `json:"Clavepubrsa"`
	Claveprivrsa     []byte `json:"Claveprivrsa"`
	Claveenclaro     string `json:"Claveenclaro"`
	Clavehashcifrado []byte `json:"Clavehashcifrado"`
	Clavehashlogin   []byte `json:"Clavehashlogin"`
}

var ClientUsuario Usuario

//Variable global, de momento para guardar clave cifrar mensajes chat1
var clavecifrarmensajes []byte

//var conn Connection
var conn *tls.Conn

// función para codificar de []bytes a string (Base64)
func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

// función para decodificar de string a []bytes (Base64)
func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	fmt.Println(err)                             // comprobamos el error
	return b                                     // devolvemos los datos originales
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
	conn, _ = tls.Dial("tcp", "127.0.0.1:444", &config)
	/*if err != nil {
		log.Fatalf("client: dial: %s", err)
	}*/
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	//Por si envia algo el servidor
	go handleServerRead()

	///////////////////////////////////
	//    PRUEBAS    /////////////////
	//////////////////////////////////
	IniciarServidorWeb()
	//login()

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
	//ClientUsuario.Nombre = "Prueba3"
	//ClientUsuario.Claveenclaro = "miclave3"
	//registrarUsuario(conn)

	///////////////////////////////////
	//    Enviar  y recibir      /////
	//////////////////////////////////

	//Enviar mensajes
	go handleClientWrite()

	//Para que no se cierre la consola
	for {
	}
}

//Convertir a json y escribir en el socket
func escribirSocket(mensaje Mensaje) {
	b, _ := json.Marshal(mensaje)
	conn.Write(b)
}

//Si envia algo el servidor a este cliente lo muestra en pantalla
func handleServerRead() {
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
			ClientUsuario.Id = idusuario
			ClientUsuario.Nombre = mensaje.Datos[1]
			ClientUsuario.Clavepubrsa = mensaje.DatosClaves[0]
			ClientUsuario.Claveprivrsa = mensaje.DatosClaves[1]

			////////////////////////
			//PRUEBAS AFTER DE LOGIN
			////////////////////////
			//nuevaClaveUsuarioConIdConjuntoClaves(conn, 1, "nuevaclave1")
			//Obtenemos la clave para cifrar mensajes del chat1
			getClaveCifrarMensajeChat(1)
		}

		if mensaje.Funcion == "DatosClaveCifrarMensajeChat" {
			laclave := mensaje.DatosClaves[0]
			clavecifrarmensajes, _ = descifrarAES(laclave, ClientUsuario.Clavehashcifrado)
		}
	}
}

//SI escribe algo lo envia al servidor
func handleClientWrite() {

	//bucle infinito
	for {
		defer conn.Close()

		//Cuando escribe algo y le da a enter
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		//Rellenar datos
		mensaje := Mensaje{From: ClientUsuario.Nombre, Funcion: "enviar", Mensajechat: []byte(message[0 : len(message)-2]), Chat: 1}

		//getClaveCifrarMensajeChat(conn, 1)

		//Ciframos el mensaje
		//De momento clave guardada globalmente pero sería hacer llamadas al servidor, channels, etc
		mensajecifrado, err := cifrarAES([]byte(mensaje.Mensajechat), clavecifrarmensajes)
		if err == true {
			fmt.Println("Error al cifrar clave con cifrado AES")
			continue
		}
		mensaje.Mensajechat = mensajecifrado

		escribirSocket(mensaje)
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
func registrarUsuario(cliente Usuario) bool {

	var err bool

	//Generamos los hash de las claves
	cliente.Clavehashlogin, cliente.Clavehashcifrado = generarHashClaves(cliente.Claveenclaro)

	//Generamos clave pública y privada RSA
	cliente.Clavepubrsa, cliente.Claveprivrsa = generarClavesRSA()

	//Clave privada del usuario cifrar
	cliente.Claveprivrsa, err = cifrarAES(cliente.Claveprivrsa, cliente.Clavehashcifrado)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return false
	}

	//Rellenar datos del mensaje
	mensaje := Mensaje{From: cliente.Nombre, Funcion: "registrarusuario", Datos: []string{cliente.Nombre}, DatosClaves: [][]byte{cliente.Clavehashlogin, cliente.Clavepubrsa, cliente.Claveprivrsa}}
	escribirSocket(mensaje)
	return true
}

//Cliente realiza login
func loginweb(usuario string, password string) {

	//reader := bufio.NewReader(os.Stdin)
	ClientUsuario.Nombre = usuario
	ClientUsuario.Claveenclaro = password

	//Generamos los hash de las claves
	ClientUsuario.Clavehashlogin, ClientUsuario.Clavehashcifrado = generarHashClaves(ClientUsuario.Claveenclaro)

	mensaje := Mensaje{From: ClientUsuario.Nombre, DatosClaves: [][]byte{ClientUsuario.Clavehashlogin}, Funcion: "login", To: -1}
	escribirSocket(mensaje)
}

//Cliente pide mensajes de un chat
func obtenerMensajesChat(idchat int) {

	mensaje := Mensaje{Chat: idchat, From: ClientUsuario.Nombre, Funcion: "obtenermensajeschat"}
	escribirSocket(mensaje)
}

//Cliente pide añadir usuarios a un chat
func agregarUsuariosChat(idchat int, usuarios []string) {

	mensaje := Mensaje{Chat: idchat, From: ClientUsuario.Nombre, Funcion: "agregarusuarioschat", Datos: usuarios}
	escribirSocket(mensaje)
}

//Cliente pide eliminar usuarios en un chat
func eliminarUsuariosChat(idchat int, usuarios []string) {

	mensaje := Mensaje{Chat: idchat, From: ClientUsuario.Nombre, Funcion: "eliminarusuarioschat", Datos: usuarios}
	escribirSocket(mensaje)
}

//Cliente pide clave pública de un usuario
func getClavePubUsuario(idusuario int) {

	mensaje := Mensaje{From: ClientUsuario.Nombre, Funcion: "getclavepubusuario", Datos: []string{strconv.Itoa(idusuario)}}
	escribirSocket(mensaje)
}

//Cliente pide clave cifrada para descifrar mensajes
func getClaveMensaje(idmensaje int) {

	mensaje := Mensaje{From: ClientUsuario.Nombre, Funcion: "getclavesmensajes", Datos: []string{strconv.Itoa(idmensaje)}}
	escribirSocket(mensaje)
}

//Cliente pide clave cifrada para descifrar mensajes
func getClaveCifrarMensajeChat(idchat int) {

	mensaje := Mensaje{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: "getclavecifrarmensajechat", Datos: []string{strconv.Itoa(idchat)}}
	escribirSocket(mensaje)
}

//Cliente crea nuevo id clave para un nuevo conjunto de claves
func CrearNuevaClaveMensajes() {

	mensaje := Mensaje{From: ClientUsuario.Nombre, Funcion: "crearnuevoidparanuevaclavemensajes"}
	escribirSocket(mensaje)
}

//Asocia nueva clave de un usuario con el id que indica ese nuevo conjunto de claves
func nuevaClaveUsuarioConIdConjuntoClaves(idconjuntoclaves int, clavemensajes string, idchat int) {

	clavesmensajeshash, _ := generarHashClaves(clavemensajes)

	//Cifrar la clave para los mensajes
	clavecifradamensajes, err := cifrarAES(clavesmensajeshash, ClientUsuario.Clavehashcifrado)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return
	}

	mensaje := Mensaje{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Chat: idchat, Funcion: "nuevaclaveusuarioconidconjuntoclaves", Datos: []string{strconv.Itoa(idconjuntoclaves)}, DatosClaves: [][]byte{clavecifradamensajes}}
	escribirSocket(mensaje)
}
