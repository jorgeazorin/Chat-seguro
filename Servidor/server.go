package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"sync"
)

//Struct con las conexiones (esto hay que hacerlo con lo del mux porque, con la memoria compartida,
// como hay muchos hilos se puede estropear si dos acceden al vector a la vez, si no se podria hacer
//el vector solo en el main e ir pasandolo)
type Conexiones struct {
	//Vector con todas las conexiones de sockets online
	conexiones []Conexion
	//El mux es para la memoria compartida
	mux sync.Mutex
}

func main() {
	//Vector de conexiones que estará en todo el programa con memoria compartida
	conexiones := Conexiones{conexiones: make([]Conexion, 0)}

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
	//escuchar atodos
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

		conexion := Conexion{conexion: conn, conexiones: &conexiones} //Creamos una nueva conexión
		//Añadimos la conxion al vector conexiones bloqueando la memoria compartida
		conexiones.mux.Lock()
		conexiones.conexiones = append(conexiones.conexiones, conexion)
		conexion.conexiones.mux.Unlock()
		//Escuchamos la conxion paralelamente o como se diga en el archivo conexion.go
		go conexion.escuchar()

	}
}
