package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

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
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

	//Imprime por pantalla lo que envia al servidor
	log.Printf("client: wrote %q (%d bytes)", message, n)

	//Por si envia algo el servidor
	go handleServerRead(conn)

	//Enviar mensajes
	go handleClientWrite(conn)

	//Para que no se cierre la consola
	for {
	}
}

//Si envia algo el servidor a este cliente lo muestra en pantalla
func handleServerRead(conn net.Conn) {
	//bucle infinito
	for {
		defer conn.Close()
		reply := make([]byte, 256)
		n, _ := conn.Read(reply)
		log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	}
}

//SI escribe algo lo envia al servidor
func handleClientWrite(conn net.Conn) {
	//bucle infinito
	for {
		//Cuando escribe algo y le da a enter
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		//Escribe esto en el socket
		n, err := io.WriteString(conn, message)
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("client: wrote %d bytes", n)
	}

}
