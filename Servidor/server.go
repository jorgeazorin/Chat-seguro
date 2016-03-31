package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"sync"
)

//Struct de conexión tiene el id del usuario y la conexión
type Conexion struct {
	conexion net.Conn
	usuario  int
}

/*
	Esto es la "Clase" de lo que va a ser el objeto c del main
	Tiene un map de conexiones y lo del mutex es para la memoria
	compartida entre los diferentes procesos

*/
type Conexiones struct {
	conexiones []Conexion
	mux        sync.Mutex
}

//	Estructura del mensaje que vamos a recibir de los clientes
type Mensaje struct {
	From     string   `json:"From"`
	To       int      `json:"To"`
	Password string   `json:"Password"`
	Funcion  string   `json:"Funcion"`
	Datos    []string `json:"Datos"`
	Mensaje  string   `json:"Mensaje"`
}

/*
	Función que guarda un socket en el map de conexiones y que se queda
	en un bucle infinito por si envia el cliente un mensaje
*/
func (c *Conexiones) handleClientRead(conexion Conexion) {

	conn := conexion.conexion
	defer conn.Close()

	///////////////////////////////////
	//    Añadimos al map la conexión con el usuario
	//////////////////////////////////

	//bloqueamos la memoria compartida
	c.mux.Lock()
	//La añadimos
	c.conexiones = append(c.conexiones, conexion)
	//Y claro la debloqueamos
	c.mux.Unlock()

	///////////////////////////////////
	//    Bucle infinito que lee cosas que envia el usuario
	//////////////////////////////////
	var mensaje Mensaje
	for {
		buf := make([]byte, 256)
		//Lee el mensaje
		n, err := conn.Read(buf)
		if err != nil {
			break
			conn.Close()
		}
		json.Unmarshal(buf[:n], &mensaje)
		c.ProcesarMensaje(conexion, mensaje)
	}
}

func (c *Conexiones) ProcesarMensaje(conexion Conexion, mensaje Mensaje) {
	var ConexionDelVector *Conexion
	for i := 0; i < len(c.conexiones); i++ {
		if c.conexiones[i] == conexion {
			ConexionDelVector = &c.conexiones[i]
		}
	}
	if mensaje.Funcion == "login" {
		ServerLogin(mensaje.From, "pass", ConexionDelVector)
		mesj := Mensaje{}
		mesj.From = mensaje.From
		mesj.Mensaje = "Usuario online"
		for i := 0; i < len(c.conexiones); i++ {
			EnviarMensajeSocket(c.conexiones[i], mesj)
		}
	}
	if mensaje.Funcion == "enviar" {
		mesj := Mensaje{}
		mesj.From = mensaje.From
		mesj.Mensaje = mensaje.Mensaje
		mesj.Funcion = "enviar"
		for i := 0; i < len(c.conexiones); i++ {
			if c.conexiones[i].conexion != conexion.conexion {
				EnviarMensajeSocket(c.conexiones[i], mesj)
			}
		}
	}
}

//FUncion que envia un mensaje a un cliente mediante un id y un string
func EnviarMensajeSocket(conexion Conexion, s Mensaje) {
	b, _ := json.Marshal(s)
	_, err := conexion.conexion.Write(b)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}

}

func main() {

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

	//escuchar atodos
	service := "0.0.0.0:443"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}

	log.Print("server: listening")

	/*
		Esto es para paralelizar usando memoria compartida:
		se crea como un objeto c que tiene en su estructura un map de conexiones
		y ese objeto tiene las funciones de crear y enviar mensaje

		Nota para entender :
		Fíjate que solo va a haber un objeto c inicializado, este objeto tiene
		un map que es un vector de conexiones lo que hay son muchas ramas que
		ejecutan una función de c

	*/
	c := Conexiones{conexiones: make([]Conexion, 0)}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())

		//creamos una nueva conexión y se la pasamos al objeto c
		conexion := Conexion{}
		conexion.conexion = conn
		go c.handleClientRead(conexion)
	}
}

func ServerLogin(user string, pass string, conexion *Conexion) bool {
	conexion.usuario, _ = strconv.Atoi(user)
	return true
}

func EnviarMensajeAChat(texto string, idchat int, idemisor int, idclave int) {

}

func CrearChat(usuarios []string, nombrechat string) {

}

func CrearUsuario(nombre string, clavepubrsa string) {

}

func CrearNuevaClaveParaMensajesBD() {

}

func GuardarClaveUsuarioMensajesBD(idclavesmensajes int, claveusuario string, idusuario int) {

}

func ObtenerChat(usuario string) {

}

func ObtenerMensajesChat(usuario string) {

}

func ObtenerUsuarios() {

}

/*
	-ModificarChatBD(idchat, nombre) OK
	-AddUsuarioChatBD(idchat, idusuariosslice) OK
	-RemoveUsuarioChatBD(idchat, idusuariosslice) OK
	-obtenerchats???
	-EditarUsuario()
*/
