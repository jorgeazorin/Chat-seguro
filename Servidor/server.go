package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net"
)

type Conexion struct {
	conexion net.Conn
	usuario  int
}

func main() {
	conexiones := []Conexion{}
	ca_b, _ := ioutil.ReadFile("ca.pem")
	ca, _ := x509.ParseCertificate(ca_b)
	priv_b, _ := ioutil.ReadFile("ca.key")
	priv, _ := x509.ParsePKCS1PrivateKey(priv_b)

	pool := x509.NewCertPool()
	pool.AddCert(ca)

	cert := tls.Certificate{
		Certificate: [][]byte{ca_b},
		PrivateKey:  priv,
	}

	config := tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
	}
	config.Rand = rand.Reader
	service := "0.0.0.0:443"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		conexion := Conexion{conn, -1}
		conexiones := append(conexiones, conexion)
		go handleClientRead(conexion, conexiones)
	}
}

func handleClientRead(conexion Conexion, conexiones []Conexion) {
	conn := conexion.conexion
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read: %s", err)
			break
		}
		log.Printf("Read %q (%d bytes)", string(buf[:n]), n)
	}
}

func handleClientWrite(conn net.Conn, s string) {

	n, err := io.WriteString(conn, s)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("server: wrote %d bytes", n)
}
