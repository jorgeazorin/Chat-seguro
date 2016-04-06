package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	//"sync"
)

var conexiones map[int]*Conexion

func main() {
	//Mapa de conexiones que estará en todo el programa
	conexiones = make(map[int]*Conexion)

	///////////////////////////////////
	//              TLS           ////
	//////////////////////////////////

	//Leer los ficheros de los certificados
	ca_b, _ := ioutil.ReadFile("ca.pem")
	ca, _ := x509.ParseCertificate(ca_b)
	priv_b, _ := ioutil.ReadFile("ca.key")
	priv, _ := x509.ParsePKCS1PrivateKey(priv_b)

	//Configurar los certificados en tls
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

	///////////////////////////////////
	//    ESCUCHAR LAS PETICIONES  ////
	//////////////////////////////////
	//escuchar a todos
	service := "0.0.0.0:444"
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

		conexion := Conexion{conexion: conn, usuario: &Usuario{}} //Creamos una nueva conexión

		//Escuchamos la conxion paralelamente o como se diga en el archivo conexion.go
		go conexion.escuchar()

	}
}
