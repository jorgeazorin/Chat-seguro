package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

//Mapa con las conexiones de los usuarios [ key=nombreusuario: value=conexion del usuario ]
var conexiones map[int]net.Conn

//Función que envia un mensaje a un cliente mediante un id y un string
func EnviarMensajeSocketSocket(conexion net.Conn, s MensajeSocket) {
	fmt.Println("Mensaje enviado: ", conexion.RemoteAddr(), "Funcion ", s.Funcion)
	fmt.Println()
	//Codifica el mensaje en json
	b, _ := json.Marshal(s)
	//log.Println(string(b))
	//Lo escribe en el socket
	_, err := conexion.Write(b)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

}

func main() {
	//Mapa de conexiones que estará en todo el programa
	conexiones = make(map[int]net.Conn)

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

	//Escuchar a todos
	service := "0.0.0.0:444"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("Servidor: escucha: %s", err)
	}

	log.Print("Servidor: escuchando")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Servidor: acepta: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("Servidor: aceptado de %s", conn.RemoteAddr())

		//Escuchamos la conxion de forma concurrente en el archivo conexion.go
		go func(conn net.Conn) {

			defer conn.Close()
			usuario := Usuario{}
			var mensaje MensajeSocket //Struct se guarda el mensaje a descodificar

			// Bucle infinito que lee cosas que envia el usuario
			for {
				buf := make([]byte, 1048576) //256

				//Lee el mensaje
				n, err := conn.Read(buf)

				if err != nil {
					_, ok := conexiones[usuario.Id]
					if ok {
						delete(conexiones, usuario.Id)
					}
					conn.Close()
					break
				}

				//Descodificar el mensaje recibido (estaba en json y se pasa a struct)
				json.Unmarshal(buf[:n], &mensaje)

				//Procesa el mensaje, esto lo hace en el archivo router.go
				ProcesarMensajeSocket(mensaje, conn, &usuario)
			}
		}(conn)
	}
}
