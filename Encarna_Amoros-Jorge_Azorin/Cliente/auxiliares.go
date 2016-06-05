package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"io/ioutil"
	"log"
)

var conn *tls.Conn

var _canalMensajeSocket = make(chan MensajeSocket)
var ClientUsuario Usuario
var chatsusuario []ChatDatos
var _clavesUsuarioDeMensajes []Clavesusuario

func main() {

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
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())
	go handleServerRead()
	IniciarServidorWeb()
	for {
	}
}

//Obtenemos las respuestas que lleguen del servidor
func handleServerRead() {
	var mensaje MensajeSocket
	for {
		defer conn.Close()
		reply := make([]byte, 1048576)
		n, err := conn.Read(reply)
		if err != nil {
			break
			conn.Close()
		}

		//Mostramos lo recibido
		json.Unmarshal(reply[:n], &mensaje)
		fmt.Println("Recibido:", mensaje.From, " -> ", mensaje.Funcion)

		//Diferenciamos entre mensajes de admin o no y lo enviamos a cliente web
		if mensaje.Funcion == Constantes_MensajeOtroClienteConectado {
			mensaje := MensajeSocket{Mensaje: "mensajedeotrocliente", Datos: []string{}}
			escribirWebSocket(mensaje)
		} else if mensaje.Funcion == Constantes_MensajeAdminOtroClienteConectado {
			mensaje := MensajeSocket{Mensaje: "mensajeadmindeotrocliente", Datos: []string{}}
			escribirWebSocket(mensaje)
		} else {
			_canalMensajeSocket <- mensaje

		}
	}

}

//Convertir a json y escribir en el socket del servidor
func escribirSocket(mensaje MensajeSocket) {
	fmt.Println("Enviado: ", mensaje.Idfrom, " -> ", mensaje.Funcion)
	mensaje.Idfrom = ClientUsuario.Id
	b, _ := json.Marshal(mensaje)
	conn.Write(b)
}

//Convertir a json y escribir en el socket con cliente web
func escribirWebSocket(mensaje MensajeSocket) {
	var s string
	b, _ := json.Marshal(mensaje)
	s = string(b)
	websocket.Message.Send(wbSocket, s)
}

//De la ontraseña en claro se realiza hash y se divide en 2 (clave login y clave cifrado)
func generarHashClaves(clave string) ([]byte, []byte) {

	//Hash con SHA-2 (256) para la contraseña en general
	clavebytes := []byte(clave)
	clavebytesconsha2 := sha256.Sum256(clavebytes)

	//Dividimos dicho HASH
	clavehashlogin := clavebytesconsha2[0 : len(clavebytesconsha2)/2]
	clavehashcifrado := clavebytesconsha2[len(clavebytesconsha2)/2 : len(clavebytesconsha2)]

	return clavehashlogin, clavehashcifrado
}

//Genera una clave pública y otra privada RSA
func generarClavesRSA() ([]byte, []byte) {
	claveprivada, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
	}

	clavepublica := &claveprivada.PublicKey
	pemblock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(claveprivada)}

	fd, _ := x509.MarshalPKIXPublicKey(clavepublica)
	pemblockPublica := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: fd}

	return pemblock.Bytes, pemblockPublica.Bytes
}

//Ciframos con RSA
func cifrarRSA(textocifrar []byte, clave []byte) ([]byte, bool) {
	r, _ := x509.ParsePKIXPublicKey(clave)
	rsaPub, _ := r.(*rsa.PublicKey)
	out, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, textocifrar, []byte{})
	return out, true
}

//Desciframos con RSA
func descifrarRSA(textocifrar []byte, clave []byte) ([]byte, bool) {
	privateKey, _ := x509.ParsePKCS1PrivateKey(clave)
	out, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, textocifrar, []byte{})
	return out, true
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

//Descifrar con AES en modo CTR
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

	for {
		if textodescifrado[len(textodescifrado)-1] != 0 {
			return textodescifrado, false
		} else {
			textodescifrado = textodescifrado[0 : len(textodescifrado)-2]
		}
	}
	return textodescifrado, false
}
