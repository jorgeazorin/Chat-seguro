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
			ClientUsuario.Estado = mensaje.Datos[2]
			ClientUsuario.Clavepubrsa = mensaje.DatosClaves[0]
			ClientUsuario.Claveprivrsa = mensaje.DatosClaves[1]
			ClientUsuario.Claveprivrsa, _ = descifrarAES(ClientUsuario.Claveprivrsa, ClientUsuario.Clavehashcifrado)
			ClientUsuario.Claveprivrsa = ClientUsuario.Claveprivrsa[0:1192]

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
	cliente.Claveprivrsa, cliente.Clavepubrsa = generarClavesRSA()
	//	fmt.Println("Privada: ", cliente.Claveprivrsa)

	//Clave privada del usuario cifrar
	cliente.Claveprivrsa, err = cifrarAES(cliente.Claveprivrsa, cliente.Clavehashcifrado)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return false
	}

	//_, _ := descifrarAES(cliente.Claveprivrsa, cliente.Clavehashcifrado)
	//fmt.Println("Privada: ", sdf)

	//Rellenar datos del mensaje
	mensaje := MensajeSocket{From: cliente.Nombre, Funcion: Constantes_registrarusuario, Datos: []string{cliente.Nombre}, DatosClaves: [][]byte{cliente.Clavehashlogin, cliente.Clavepubrsa, cliente.Claveprivrsa}}
	escribirSocket(mensaje)

	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket
		if mensaje.Funcion == Constantes_registrarusuario_ok {
			validar = 1
			return true
		}
		if mensaje.Funcion == Constantes_registrarusuario_err {
			validar = 1
			return false
		}
	}

	return true

}

//Enviar un mensaje
func enviarMensaje(mensaje MensajeSocket) bool {

	//Obtenemos la clave para cifrar mensaje
	var clavecifrarmensajes []byte
	var idclavecifrarmensajes int
	idUltimaClaveDelChat := -1
	chatsusuario = obtenerChats()

	for i := 0; i < len(chatsusuario); i++ {
		if chatsusuario[i].Chat.Id == mensaje.Chat {
			idUltimaClaveDelChat = chatsusuario[i].Chat.UltimaClave
			fmt.Println("Se ha encontrado una clave para este chat")
		}
	}

	//Si no existe la creamos y si si que existe pues no la creamos y la cogemos
	if idUltimaClaveDelChat <= 0 {
		idclavecifrarmensajes, clavecifrarmensajes, _ = CrearNuevaClaveMensajes(mensaje.Chat)
	} else {
		for i := 0; i < len(_clavesUsuarioDeMensajes); i++ {
			if _clavesUsuarioDeMensajes[i].Idclavesmensajes == idUltimaClaveDelChat {
				clavecifrarmensajes = _clavesUsuarioDeMensajes[i].Clavemensajes
			}
		}
		idclavecifrarmensajes = idUltimaClaveDelChat
	}
	clavecifrarmensajes = clavecifrarmensajes[0:32]
	mensajecifrado, err := cifrarAES(mensaje.Mensajechat, clavecifrarmensajes)
	//key := []byte
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AESs ", mensaje.Mensajechat, "  :   ", clavecifrarmensajes)
		return false
	}
	mensaje.Mensajechat = mensajecifrado

	//Rellenar datos
	mensaje = MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_enviar, Datos: []string{strconv.Itoa(idclavecifrarmensajes)}, Mensajechat: mensaje.Mensajechat, Chat: mensaje.Chat}
	escribirSocket(mensaje)
	return true
}

//Cliente pide mensajes de un chat
func obtenerMensajesChat(idchat int) []MensajeDatos {
	returnmensajes := []MensajeDatos{}
	mensaje := MensajeSocket{Chat: idchat, Idfrom: ClientUsuario.Id, From: ClientUsuario.Nombre, Funcion: Constantes_obtenermensajeschat}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket
		if mensaje.Funcion == Constantes_obtenermensajeschat_err {
			validar = 1
			fmt.Println("Error obteniendo chats")
		}
		if mensaje.Funcion == Constantes_obtenermensajeschat_ok {
			validar = 1
			chatsusuario = nil
			for i := 0; i < len(mensaje.Datos); i++ {
				var mensajeDatos MensajeDatos
				json.Unmarshal([]byte(mensaje.Datos[i]), &mensajeDatos)
				idClaveMensaje := mensajeDatos.Mensaje.IdClave
				Clave := []byte{}
				for j := 0; j < len(_clavesUsuarioDeMensajes); j++ {
					if _clavesUsuarioDeMensajes[j].Idclavesmensajes == idClaveMensaje {
						Clave = _clavesUsuarioDeMensajes[j].Clavemensajes

						Clave = Clave[0:32]
						mensajeDatos.Mensaje.Texto, _ = descifrarAES(mensajeDatos.Mensaje.Texto, Clave)

						mensajeDatos.Mensaje.TextoClaro = string(mensajeDatos.Mensaje.Texto)
						mensajeDatos.Mensaje.Texto = []byte{}
						returnmensajes = append(returnmensajes, mensajeDatos)
					}
				}
			}
		}

	}

	return returnmensajes
}

//Cliente pide añadir usuarios a un chat
func agregarUsuariosChat(idchat int, usuarios []string) {
	mensaje := MensajeSocket{Chat: idchat, Idfrom: ClientUsuario.Id, From: ClientUsuario.Nombre, Funcion: Constantes_agregarusuarioschat, Datos: usuarios}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket
		if mensaje.Funcion == Constantes_agregarusuarioschat_err {
			validar = -1
		}
		if mensaje.Funcion == Constantes_agregarusuarioschat_ok {
			CrearNuevaClaveMensajes(idchat)
			validar = 1
		}
	}
}

//Cliente pide eliminar usuarios en un chat
func eliminarUsuariosChat(idchat int, usuarios []string) {

	mensaje := MensajeSocket{Chat: idchat, From: ClientUsuario.Nombre, Funcion: Constantes_eliminarusuarioschat, Datos: usuarios}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket
		if mensaje.Funcion == Constantes_eliminarusuarioschat_err {
			validar = -1
		}
		if mensaje.Funcion == Constantes_eliminarusuarioschat_ok {
			CrearNuevaClaveMensajes(idchat)
			validar = 1
		}
	}
}

//Cliente lee mensajes del chay
func MarcarChatComoLeido(idchat int) {

	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Chat: idchat, Funcion: Constantes_marcarchatcomoleido}
	escribirSocket(mensaje)
}

func MarcarMensajeComoLeido(id int) {
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Datos: []string{strconv.Itoa(id)}, Funcion: Constantes_marcarmensajeleido}
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
				json.Unmarshal([]byte(mensaje.Datos[i]), &chatsusuario)
				chats = append(chats, chatsusuario)
			}

		}
	}
	return chats
}

//Modificar chat
func editarChat(chat Chat) {

	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Chat: chat.Id, Datos: []string{chat.Nombre}, Funcion: Constantes_modificarchat}
	escribirSocket(mensaje)
}

//Modificar usuario
func editarUsuario(usuario Usuario) {

	mensaje := MensajeSocket{From: usuario.Nombre, Idfrom: ClientUsuario.Id, Datos: []string{usuario.Estado}, Funcion: Constantes_modificarusuario}
	escribirSocket(mensaje)
}

//Obtener usuarios
func getUsuarios() []Usuario {

	var allusuarios = []Usuario{}
	mensaje := MensajeSocket{Idfrom: ClientUsuario.Id, From: ClientUsuario.Nombre, Funcion: Constantes_getUsuarios}
	escribirSocket(mensaje)

	//Leemos los mensajes que recibimos del servidor hasta que sea un error de getUsuarios o un ok
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		//Si el getUsuarios es incorrecto lo mostramos por pantalla
		if mensaje.Funcion == Constantes_getUsuarios_err {
			validar = 1
			fmt.Println("Error al obtener los usuarios")
		}

		//Si el getUsuarios es correcto rellenamos usuarios
		if mensaje.Funcion == Constantes_getUsuarios_ok {
			validar = 1
			for i := 0; i < len(mensaje.Datos); i++ {
				var usuario = Usuario{}
				json.Unmarshal([]byte(mensaje.Datos[i]), &usuario)
				allusuarios = append(allusuarios, usuario)
			}
		}
	}

	return allusuarios
}

////////////////////////////////////////////////////
////////////////////////////////////////////////////

//
//
//

////////////////////////////////////////////////////
////////// FUNCIOS AUXILIARES       /////////////
/////////////////////////////////////////////////

func obtenermensajesAdmin() {
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_obtenermensajesAdmin}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_obtenermensajesAdmin_ok {
			validar = 1
			for i := 0; i < len(mensaje.Datos); i++ {
				b, _ := descifrarRSA(mensaje.DatosClaves[i], ClientUsuario.Claveprivrsa)

				var mensaje1 = MensajeDatos{}
				json.Unmarshal([]byte(mensaje.Datos[i]), &mensaje1)
				idClave := mensaje1.Mensaje.IdClave
				GuardarClaveUsuarioConIdConjuntoClaves(idClave, b)
				MarcarMensajeComoLeido(mensaje1.Mensaje.Id)
			}
		}
		if mensaje.Funcion == Constantes_obtenermensajesAdmin_err {
			validar = 1
			fmt.Println("Error al recibir mensajes Admin.")
		}
	}
}

//Obtiene todas las claves del usuario para descifrar mensajes
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
				clavesUsuarioDeMensajes[i].Clavemensajes, _ = descifrarAES(clavesUsuarioDeMensajes[i].Clavemensajes, ClientUsuario.Clavehashcifrado)
			}

		} else if mensaje.Funcion == Constantes_getClavesDeUnUsuario_err {
			validar = 1
			fmt.Println("Error obteniendo las claves de un usuario")
		}
	}
	return clavesUsuarioDeMensajes
}

//Cliente crea nuevo id clave para un nuevo conjunto de claves
// y se la envia a los usuarios del chat que se le pasa por parametro
func CrearNuevaClaveMensajes(idChat int) (int, []byte, []byte) {
	id := -1
	clavenuevacifrada := []byte{}
	clavenuevassincifrar := []byte{}

	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Funcion: Constantes_crearnuevoidparanuevaclavemensajes}
	escribirSocket(mensaje)
	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_crearnuevoidparanuevaclavemensajes_ok {
			id, _ = strconv.Atoi(mensaje.Datos[0])
			clavenuevassincifrar, clavenuevacifrada = nuevaClaveUsuarioConIdConjuntoClaves(id)

			//Ahora se la vamos a enviar a los usuarios que pertenecen al chat
			//Obtenemos sus claves privadas
			UsuariosAEnviarLasClaves := ObtenerClavesPrivadasDeMuchosUsuarios(ObtenerUsuariosDeUnChat(idChat))
			for i := 0; i < len(UsuariosAEnviarLasClaves); i++ {
				//Ciframos la clave que le vamos a enviar con la clave publica del usuario
				claveCifradaAEnviar, _ := cifrarRSA(clavenuevassincifrar, UsuariosAEnviarLasClaves[i].Clavepubrsa)
				//Enviamos la clave en un mensaje de administracion al usuario
				mensaje = MensajeSocket{From: "-200", Idfrom: ClientUsuario.Id, Funcion: Constantes_enviar, Datos: []string{strconv.Itoa(id), strconv.Itoa(UsuariosAEnviarLasClaves[i].Id)}, Mensajechat: claveCifradaAEnviar, Chat: idChat}
				escribirSocket(mensaje)
			}
			//Asociar ultima clave creada al chat
			mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_AsociarNuevaClaveAChat, Datos: []string{strconv.Itoa(idChat), strconv.Itoa(id)}}
			escribirSocket(mensaje)
			validar = 1

		}
		if mensaje.Funcion == Constantes_crearnuevoidparanuevaclavemensajes_err {
			validar = 1
		}
	}
	return id, clavenuevassincifrar, clavenuevacifrada
}

//Crea nueva clave de un usuario con el id que indica ese nuevo conjunto de claves
func nuevaClaveUsuarioConIdConjuntoClaves(idconjuntoclaves int) ([]byte, []byte) {
	key := make([]byte, 32)
	_, errr := rand.Read(key)
	if errr != nil {
	}
	var keycifrada, err = cifrarAES(key, ClientUsuario.Clavehashcifrado)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return []byte{}, []byte{}
	}
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_nuevaclaveusuarioconidconjuntoclaves, Datos: []string{strconv.Itoa(idconjuntoclaves)}, DatosClaves: [][]byte{keycifrada}}
	escribirSocket(mensaje)

	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_nuevaclaveusuarioconidconjuntoclaves_ok {
			validar = 2
		}
		if mensaje.Funcion == Constantes_nuevaclaveusuarioconidconjuntoclaves_err {
			validar = 1
			key = []byte{}
			keycifrada = []byte{}
		}
	}

	return key, keycifrada
}

//Asocia nueva clave de un usuario con el id que indica ese nuevo conjunto de claves
func GuardarClaveUsuarioConIdConjuntoClaves(idconjuntoclaves int, clave []byte) bool {
	fmt.Println("Se va a ver un follon ", idconjuntoclaves)
	var keycifrada, err = cifrarAES(clave, ClientUsuario.Clavehashcifrado)
	if err == true {
		fmt.Println("Error al cifrar clave con cifrado AES")
		return false
	}
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Idfrom: ClientUsuario.Id, Funcion: Constantes_nuevaclaveusuarioconidconjuntoclaves, Datos: []string{strconv.Itoa(idconjuntoclaves)}, DatosClaves: [][]byte{keycifrada}}
	escribirSocket(mensaje)

	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_nuevaclaveusuarioconidconjuntoclaves_ok {
			validar = 2
			return true
		}
		if mensaje.Funcion == Constantes_nuevaclaveusuarioconidconjuntoclaves_err {
			validar = 1
			return false
		}
	}

	return true
}

//Obtiene las claves PUBLICAS PUBLICAS PUBLICAS de los usuarios que se le pasan
func ObtenerClavesPrivadasDeMuchosUsuarios(usuarios []int) []Usuario {
	returnusuarios := []Usuario{}
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
				user.Clavepubrsa = mensaje.DatosClaves[i]
				if user.Id != ClientUsuario.Id {
					returnusuarios = append(returnusuarios, user)
				}

			}
		}
		if mensaje.Funcion == Constantes_obtenerClavesDeMuchosUsuarios_err {
			validar = 1
		}
	}
	return returnusuarios
}

func ObtenerUsuariosDeUnChat(chat int) []int {
	returnusuarios := []int{}
	usuariosenv := make([]string, 0, 1)
	usuariosenv = append(usuariosenv, strconv.Itoa(chat))
	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Datos: usuariosenv, Idfrom: ClientUsuario.Id, Funcion: Constantes_getUsuariosDeUnChat}
	escribirSocket(mensaje)

	for validar := 0; validar == 0; {
		mensaje = <-_canalMensajeSocket

		if mensaje.Funcion == Constantes_getUsuariosDeUnChat_ok {
			validar = 1
			for i := 0; i < len(mensaje.Datos); i++ {
				user, _ := strconv.Atoi(mensaje.Datos[i])
				returnusuarios = append(returnusuarios, user)
			}
		}
		if mensaje.Funcion == Constantes_getUsuariosDeUnChat_err {
			validar = 1
		}
	}
	return returnusuarios
}

//
//
//
//
//
//
///////////////////
/// Estas ahora sirven para algo?¿?¿?
///////////////////

//Cliente pide clave pública de un usuario
func getClavePubUsuario(idusuario int) {

	mensaje := MensajeSocket{From: ClientUsuario.Nombre, Funcion: Constantes_getclavepubusuario, Datos: []string{strconv.Itoa(idusuario)}}
	escribirSocket(mensaje)
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

////////////////////////////////////////////////
////////////////////////////////////////////////
