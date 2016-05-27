package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////
////////// FUNCIONES DIRECTAS DE WEB /////////////
/////////////////////////////////////////////////

//Cliente realiza login
func loginweb(usuario string, password string) {
	//reader := bufio.NewReader(os.Stdin)
	ClientUsuario.Nombre = usuario
	ClientUsuario.Claveenclaro = password
	//Generamos los hash de las claves
	ClientUsuario.Clavehashlogin, ClientUsuario.Clavehashcifrado = generarHashClaves(ClientUsuario.Claveenclaro)
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, DatosClaves: [][]byte{ClientUsuario.Clavehashlogin}, Funcion: Constantes_login, To: -1}
	escribirSocket(mensaje)

	//Leemos los mensajes que recibimos del servidor hasta que sea un error de login o un login correcto
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		//Si el login es incorrecto lo mostramos por pantalla
		if mensaje.Funcion == Constantes_getClavesDeUnUsuario_err {
			validar = 1
			fmt.Println("Error validando usuario")
		}

		//Si el login es correcto rellenamos el usuario
		if mensaje.Funcion == Constantes_login_ok {
			validar = 2
			idusuario, _ := strconv.Atoi(mensaje.Datos[0])
			ClientUsuario.Id = idusuario
			ClientUsuario.Nombre = mensaje.Datos[1]
			ClientUsuario.Clavepubrsa = mensaje.DatosClaves[0]
			ClientUsuario.Claveprivrsa = mensaje.DatosClaves[1]
			ClientUsuario.Claveprivrsa, _ = descifrarAES(ClientUsuario.Claveprivrsa, ClientUsuario.Clavehashcifrado)
		}
		//Si el login es correcto y se ha rellenado correctamente obtemenos sus mensajes de administración
		if validar == 2 {
			obtenermensajesAdmin()
			chatsusuario = obtenerChats()
			_clavesUsuarioDeMensajes = getTodasLasClavesDeUnUsuario()
		}

	}

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
	mensaje := MensajeSocket{From: cliente.Nombre, Funcion: Constantes_registrarusuario, Datos: []string{cliente.Nombre}, DatosClaves: [][]byte{cliente.Clavehashlogin, cliente.Clavepubrsa, cliente.Claveprivrsa}}
	escribirSocket(mensaje)
	return true
}

//Enviar un mensaje
func enviarMensaje(mensaje MensajeSocket) bool {

	//Obtenemos la clave para cifrar mensaje
	var clavecifrarmensajes []byte
	var idclavecifrarmensajes int
	for i := 0; i < len(chatsusuario); i++ {
		if chatsusuario[i].Chat.Id == mensaje.Chat {
			clavecifrarmensajes = chatsusuario[i].Clave
			idclavecifrarmensajes = chatsusuario[i].IdClave
		}
	}

	mensajecifrado, err := cifrarAES(mensaje.Mensajechat, clavecifrarmensajes)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return false
	}
	mensaje.Mensajechat = mensajecifrado

	//Rellenar datos
	mensaje = MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_enviar, Datos: []string{strconv.Itoa(idclavecifrarmensajes)}, Mensajechat: mensaje.Mensajechat, Chat: mensaje.Chat}
	escribirSocket(mensaje)
	return true
}

//Cliente pide mensajes de un chat
func obtenerMensajesChat(idchat int) {
	mensaje := MensajeSocket{Chat: idchat, From: ClientUsuario.Nombre, Funcion: Constantes_obtenermensajeschat}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_obtenermensajeschat_ok {
			validar = 1
			chatsusuario = nil
			var chatdatos ChatDatos
			var err bool

			for i := 0; i < len(mensaje.Datos); i++ {
				json.Unmarshal([]byte(mensaje.Datos[i]), &chatdatos)

				//Descifrando cada mensaje
				for j := 0; j < len(chatdatos.MensajesDatos); j++ {

					chatdatos.MensajesDatos[j].Mensaje.Texto, err = descifrarAES(chatdatos.MensajesDatos[j].Mensaje.Texto, chatdatos.MensajesDatos[j].Mensaje.Clave)
					if err == true {
						fmt.Println("Error al descifrar mensajes.")
						mensaje.Mensaje = "Error al obtener los mensajes."
						return
					}
					chatdatos.MensajesDatos[j].Mensaje.TextoClaro = string(chatdatos.MensajesDatos[j].Mensaje.Texto)
					chatdatos.MensajesDatos[j].Mensaje.Texto = []byte{}
				}

				chatsusuario = append(chatsusuario, chatdatos)
				chatdatosdescifrados, _ := json.Marshal(chatdatos)
				mensaje.Datos[i] = string(chatdatosdescifrados)
			}
		}
		if mensaje.Funcion == Constantes_obtenermensajeschat_err {
			validar = 1
			fmt.Println("Error obteniendo chats")
		}
	}
}

//Cliente pide añadir usuarios a un chat
func agregarUsuariosChat(idchat int, usuarios []string) {

	fmt.Println(randSeq(20))
	fmt.Println(idchat, usuarios)

	mensaje := MensajeSocket{Chat: idchat, Idfrom: ClientUsuario.Id, From: ClientUsuario.Nombre, Funcion: Constantes_agregarusuarioschat, Datos: usuarios}
	escribirSocket(mensaje)
}

//Cliente pide eliminar usuarios en un chat
func eliminarUsuariosChat(idchat int, usuarios []string) {

	mensaje := MensajeSocket{Chat: idchat, From: ClientUsuario.Nombre, Funcion: Constantes_eliminarusuarioschat, Datos: usuarios}
	escribirSocket(mensaje)
}

//Cliente lee mensajes del chay
func MarcarChatComoLeido(idchat int) {

	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Chat: idchat, Funcion: Constantes_marcarchatcomoleido}
	escribirSocket(mensaje)
}

func crearChat(nombrechat string) int {
	r := -1
	datos := make([]string, 0, 1)
	datos = append(datos, nombrechat)
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
	}

	mensaje := MensajeSocket{Idfrom: ClientUsuario.Id, Datos: datos, Funcion: Constantes_crearchat}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		//Si el login es incorrecto lo mostramos por pantalla
		if mensaje.Funcion == Constantes_crearchat_err {
			validar = 1

		}
		if mensaje.Funcion == Constantes_crearchat_ok {

			validar = 2
			r, _ = strconv.Atoi(mensaje.Datos[0])
		}
	}
	return r
}

////////////////////////////////////////////////////
////////////////////////////////////////////////////

//
//
//

////////////////////////////////////////////////////
////////// FUNCIONES AUXILIARES       /////////////
/////////////////////////////////////////////////

func obtenermensajesAdmin() {
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Funcion: Constantes_obtenermensajesAdmin}
	escribirSocket(mensaje)
}

//Cliente pide todos los chats
func obtenerChats() []ChatDatos {
	chats := []ChatDatos{}
	mensaje := MensajeSocket{Idfrom: ClientUsuario.Id, From: ClientUsuario.Nombre, Funcion: Constantes_obtenerchats}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		//Si el login es incorrecto lo mostramos por pantalla
		if mensaje.Funcion == Constantes_obtenerchats_err {
			validar = 1
			fmt.Println("Error obteniendo Chats")
		}

		//Si el login es correcto rellenamos el usuario
		if mensaje.Funcion == Constantes_obtenerchats_ok {
			validar = 1
			for i := 0; i < len(mensaje.Datos); i++ {
				var chatsusuario = ChatDatos{}
				json.Unmarshal([]byte(mensaje.Datos[0]), &chatsusuario)
				chats = append(chats, chatsusuario)
			}

		}
	}
	return chats
}

//Cliente pide clave pública de un usuario
func getClavePubUsuario(idusuario int) {

	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Funcion: Constantes_getclavepubusuario, Datos: []string{strconv.Itoa(idusuario)}}
	escribirSocket(mensaje)
}

func getTodasLasClavesDeUnUsuario() []Clavesusuario {
	clavesUsuarioDeMensajes := []Clavesusuario{}
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_getClavesDeUnUsuario, Datos: []string{}}
	escribirSocket(mensaje)
	//for validar := 0; validar == 0; {

	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_getClavesDeUnUsuario_ok {
			validar = 1
			for i := 0; i < len(mensaje.Datos); i++ {
				var clave = Clavesusuario{}
				json.Unmarshal([]byte(mensaje.Datos[i]), &clave)
				clave.Clavemensajes = mensaje.DatosClaves[i]
				clavesUsuarioDeMensajes = append(clavesUsuarioDeMensajes, clave)
			}
			for i := 0; i < len(clavesUsuarioDeMensajes); i++ {
				clavesUsuarioDeMensajes[i].Clavemensajes, _ = descifrarAES(ClientUsuario.Clavehashcifrado, clavesUsuarioDeMensajes[i].Clavemensajes)
				fmt.Println("clave recibida del usuario ", clavesUsuarioDeMensajes[i].Idusuario, " y con calveid ", clavesUsuarioDeMensajes[i].Idclavesmensajes, ": ", clavesUsuarioDeMensajes[i].Clavemensajes)
			}
		} else if mensaje.Funcion == Constantes_getClavesDeUnUsuario_err {
			validar = 1
			fmt.Println("error obteniendo las claves de un usuario")
		}
	}
	return clavesUsuarioDeMensajes
}

//Cliente pide clave cifrada para descifrar mensajes
func getClaveMensaje(idmensaje int) {

	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Funcion: Constantes_getclavesmensajes, Datos: []string{strconv.Itoa(idmensaje)}}
	escribirSocket(mensaje)
}

//Cliente pide clave cifrada para descifrar mensajes
func getClaveCifrarMensajeChat(idchat int) {
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_getclavecifrarmensajechat, Datos: []string{strconv.Itoa(idchat)}}
	escribirSocket(mensaje)
}

//Cliente crea nuevo id clave para un nuevo conjunto de claves
func CrearNuevaClaveMensajes() {

	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Funcion: Constantes_crearnuevoidparanuevaclavemensajes}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_crearnuevoidparanuevaclavemensajes_ok {
			id, _ := strconv.Atoi(mensaje.Datos[0])
			nuevaClaveUsuarioConIdConjuntoClaves(id)
			validar = 1
		}
		if mensaje.Funcion == Constantes_crearnuevoidparanuevaclavemensajes_err {
			validar = 1
		}
	}
}

//Asocia nueva clave de un usuario con el id que indica ese nuevo conjunto de claves
func nuevaClaveUsuarioConIdConjuntoClaves(idconjuntoclaves int) {
	key := make([]byte, 64)
	_, errr := rand.Read(key)
	if errr != nil {
	}
	var keycifrada, err = cifrarAES(key, ClientUsuario.Clavehashcifrado)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return
	}
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_nuevaclaveusuarioconidconjuntoclaves, Datos: []string{strconv.Itoa(idconjuntoclaves)}, DatosClaves: [][]byte{keycifrada}}
	escribirSocket(mensaje)
	fmt.Println("!!!Se va a crear la clave idmensaje ", idconjuntoclaves, " clave descifrada ", key, " clavecifrada: ", keycifrada)
}

func ObtenerClavesPrivadasDeMuchosUsuarios(usuarios []int) {
	usuariosenv := make([]string, 0, 1)
	for i := 0; i < len(usuarios); i++ {
		usuariosenv = append(usuariosenv, strconv.Itoa(usuarios[i]))
	}
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Datos: usuariosenv, Idfrom: ClientUsuario.Id, Funcion: Constantes_obtenerClavesDeMuchosUsuarios}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_obtenerClavesDeMuchosUsuarios_ok {
			validar = 1
			for i := 0; i < len(mensaje.Datos); i++ {
				var user = Usuario{}

				json.Unmarshal([]byte(mensaje.Datos[i]), &user)
				user.Claveprivrsa = mensaje.DatosClaves[i]
				fmt.Println("clave recibida de usuario ", user.Nombre, "  clave privada ", user.Claveprivrsa)
			}
		}
		if mensaje.Funcion == Constantes_obtenerClavesDeMuchosUsuarios_err {
			validar = 1
		}
	}
}

////////////////////////////////////////////////
////////////////////////////////////////////////
