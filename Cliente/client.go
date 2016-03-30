package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
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

	//CONECTAR
	conn, err := tls.Dial("tcp", "127.0.0.1:443", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	message := "Hello\n"
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

	log.Printf("client: wrote %q (%d bytes)", message, n)
	for {
		defer conn.Close()
		reply := make([]byte, 256)
		n, err = conn.Read(reply)
		log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	}
}
