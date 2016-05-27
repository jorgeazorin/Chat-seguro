package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

//Struct de los mensajes que se envian por el socket
type MensajeSocket struct {
	From          string   `json:"From"`
	Idfrom        int      `json:"Idfrom"`
	To            int      `json:"To"`
	Password      string   `json:"Password"`
	Funcion       int      `json:"Funcion"`
	Datos         []string `json:"Datos"`
	DatosClaves   [][]byte `json:"DatosClaves"`
	Chat          int      `json:"Chat"`
	MensajeSocket string   `json:"MensajeSocket"`
	Mensajechat   []byte   `json:"Mensajechat"`
}

func ProcesarMensajeSocket(mensaje MensajeSocket, conexion net.Conn, usuario *Usuario) {

	//Para las operaciones con la BD
	var bd BD
	bd.username = "root"
	bd.password = ""
	bd.adress = ""
	bd.database = "sds"

	///////////////////
	//REGISTRAR USUARIO
	///////////////////
	if mensaje.Funcion == Constantes_registrarusuario {

		var usuarionuevo Usuario

		usuarionuevo.Nombre = mensaje.Datos[0]
		usuarionuevo.Clavelogin = mensaje.DatosClaves[0]
		usuarionuevo.Clavepubrsa = mensaje.DatosClaves[1]
		usuarionuevo.Claveprivrsa = mensaje.DatosClaves[2]

		usuario, test := bd.insertUsuarioBD(usuarionuevo)

		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_registrarusuario_err, MensajeSocket: "Error en el registro."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_registrarusuario_ok, Datos: []string{strconv.Itoa(usuario.Id), usuario.Nombre, usuario.Estado}, DatosClaves: [][]byte{usuario.Clavepubrsa, usuario.Claveprivrsa}, MensajeSocket: "Registrado correctamente"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////
	//INICIAR SESIÓN
	////////////////
	if mensaje.Funcion == Constantes_login {

		//Rellenamos el usuario de la conexión con el login
		usuario, test := bd.loginUsuarioBD(mensaje.From, mensaje.DatosClaves[0])

		//Si login incorrecto se lo decimos al cliente
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_login_err, MensajeSocket: "Nombre de usuario o contraseña incorrectos."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Si es correcto, aAñadimos la conexion al map de conexiones
		conexiones[usuario.Id] = conexion

		//Enviamos un mensaje de todo OK al usuario logeado
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_login_ok, Datos: []string{strconv.Itoa(usuario.Id), usuario.Nombre, usuario.Estado}, DatosClaves: [][]byte{usuario.Clavepubrsa, usuario.Claveprivrsa}, MensajeSocket: "Logeado correctamente"}
		EnviarMensajeSocketSocket(conexion, mesj)

	}

	if mensaje.Funcion == Constantes_obtenermensajesAdmin {
		mensajes, test := bd.getMensajesAdmin(mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenermensajesAdmin_err, MensajeSocket: "Error obtenermensajesAdmin."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
		datos := make([]string, 0, 1)
		for i := 0; i < len(mensajes); i++ {
			men := Mensaje{Id: mensajes[i].Mensaje.Id, Texto: mensajes[i].Mensaje.Texto}
			b, _ := json.Marshal(men)
			datos = append(datos, string(b))
		}
		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenermensajesAdmin_ok, Datos: datos, MensajeSocket: "obtenermensajesAdmin:"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	///////////////////
	//OBTENER LOS CHATS
	///////////////////
	if mensaje.Funcion == Constantes_obtenerchats {

		//Llamada BD obtener chats del usuario
		chats, test := bd.getChatsUsuarioBD(mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenerchats_err, MensajeSocket: "Error al obtener los chats."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Ponemos nombres del chat bien (los no grupales nombre del receptor)
		for i := 0; i < len(chats); i++ {
			if chats[i].Chat.Nombre == "" {
				idusuarios, _ := bd.getUsuariosChatBD(chats[i].Chat.Id)
				for j := 0; j < len(idusuarios); j++ {
					if idusuarios[i] != mensaje.Idfrom {
						chats[i].Chat.Nombre, _ = bd.getNombreUsuario(idusuarios[i])
					}
				}
			}
		}

		//Codificamos con Marshal
		datos := make([]string, 0, 1)
		for i := 0; i < len(chats); i++ {
			b, _ := json.Marshal(chats[i])
			datos = append(datos, string(b))
		}

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenerchats_ok, Datos: datos, MensajeSocket: "Chats:"}
		EnviarMensajeSocketSocket(conexion, mesj)

	}

	////////////////////////
	//CREAR MENSAJE Y ENVIAR
	////////////////////////
	if mensaje.Funcion == Constantes_enviar {

		//Guardamos los mensajes en la BD
		idclavemensaje, _ := strconv.Atoi(mensaje.Datos[0])
		m := Mensaje{Texto: mensaje.Mensajechat, Chat: mensaje.Chat, Emisor: mensaje.Idfrom, Clave: idclavemensaje}
		bd.guardarMensajeBD(m)

		//Obtenemos los usuarios que pertenecen en el chat
		idusuarios, test := bd.getUsuariosChatBD(mensaje.Chat)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_enviar_err, MensajeSocket: "Error al enviar mensaje."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos el mensaje a todos los usuarios de ese chat (incluido el emisor)
		for i := 0; i < len(idusuarios); i++ {
			conexion, ok := conexiones[idusuarios[i]]
			if ok {
				mensaje.MensajeSocket = "MensajeEnviado:"
				mensaje.Funcion = Constantes_enviar_ok
				EnviarMensajeSocketSocket(conexion, mensaje)
			}
		}

	}

	/////////////////////////////
	//OBTENER MENSAJES DE UN CHAT
	/////////////////////////////
	if mensaje.Funcion == Constantes_obtenermensajeschat {

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(mensaje.Idfrom, mensaje.Chat)
		if permitido == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenermensajeschat_err, MensajeSocket: "Error de permiso. No perteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Obtenemos los mensajes
		mensajes, test := bd.getMensajesChatBD(mensaje.Chat, mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenermensajeschat_err, MensajeSocket: "Error al obtener los mensajes del chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Los convertimos con marshall
		datos := make([]string, 0, 1)
		for i := 0; i < len(mensajes); i++ {
			men := mensajes[i]
			b, _ := json.Marshal(men)
			datos = append(datos, string(b))
		}

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Datos: datos, Funcion: Constantes_obtenermensajeschat_ok, MensajeSocket: "Mensajes recibidos:"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////////////////
	//AGREGAMOS USUARIOS AL CHAT
	////////////////////////////
	if mensaje.Funcion == Constantes_agregarusuarioschat {

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(mensaje.Idfrom, mensaje.Chat)
		if permitido == false {
			//Enviamos mensaje error
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_agregarusuarioschat_err, MensajeSocket: "Error de permiso. No perteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Todos los usuarios del mensaje
		nombresusuarios := make([]string, 0, 1)
		for i := 0; i < len(mensaje.Datos); i++ {
			nombresusuarios = append(nombresusuarios, mensaje.Datos[i])
		}
		fmt.Println(nombresusuarios)
		/*/Los agregamos llamando a la BD
		test := bd.addUsuariosChatBD(mensaje.Chat, nombresusuarios)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Hubo un error al añadir usuarios al chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Creamos nueva clave para el nuevo conjunto de mensajes
		idclave, test := bd.CrearNuevaClaveMensajesBD()
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Hubo un error al añadir usuarios al chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}*/

		//Asociamos a todos los usuarios con la nueva clave

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_agregarusuarioschat_ok, MensajeSocket: "Usuarios añadidos correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	/////////////////////////////
	//ELIMINAMOS USUARIOS AL CHAT
	/////////////////////////////
	if mensaje.Funcion == Constantes_eliminarusuarioschat {

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(usuario.Id, mensaje.Chat)
		if permitido == false {
			//Enviamos mensaje error
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_eliminarusuarioschat_err, MensajeSocket: "Error de permiso. No perteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Obtenemos los usuarios del mensaje
		idusuarios := make([]int, 0, 1)
		for i := 0; i < len(mensaje.Datos); i++ {
			idusuario, _ := strconv.Atoi(mensaje.Datos[i])
			idusuarios = append(idusuarios, idusuario)
		}

		//Los eliminamos llamando a la BD
		test := bd.removeUsuariosChatBD(mensaje.Chat, idusuarios)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_eliminarusuarioschat_err, MensajeSocket: "Hubo un error al eliminar usuarios del chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_eliminarusuarioschat_ok, MensajeSocket: "Usuarios eliminados correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	//////////////////////////////////
	//OBTENER CLAVE PÚBLICA DE USUARIO
	//////////////////////////////////
	if mensaje.Funcion == Constantes_getclavepubusuario {

		//Obtenemos id del usuario
		idusuario, _ := strconv.Atoi(mensaje.Datos[0])

		//Llamamos a la bd para obtener la clave pública de un usuario
		clavepub, test := bd.getClavePubUsuario(idusuario)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getclavepubusuario_err, MensajeSocket: "Error al obtener la clave del usuario."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getclavepubusuario_ok, DatosClaves: [][]byte{clavepub}, MensajeSocket: "Clave publica del usuario obtenida correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////////////////////////////////
	//OBTENER CLAVE PARA CIFRAR MENSAJES DE CHAT
	////////////////////////////////////////////
	if mensaje.Funcion == Constantes_getclavecifrarmensajechat {

		//Obtenemos id del mensaje
		idchat, _ := strconv.Atoi(mensaje.Datos[0])

		prueba := bd.usuarioEnChat(mensaje.Idfrom, idchat)
		if prueba == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getclavecifrarmensajechat_err, MensajeSocket: "No tienes permiso para acceder a los datos de este mensaje."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Lamamos a la BD para obtener la última clave con la que cifrar los mensajes
		clavemensaje, idclavemensaje, test := bd.getLastKeyMensaje(idchat, mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getclavecifrarmensajechat_err, MensajeSocket: "Error al obtener clave para cifrar mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getclavecifrarmensajechat_ok, Datos: []string{strconv.Itoa(idclavemensaje)}, DatosClaves: [][]byte{clavemensaje}, MensajeSocket: "Clave para mensajes obtenida correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	//////////////////////////////////
	//NUEVA CLAVE PARA CIFRAR MENSAJES
	//////////////////////////////////
	if mensaje.Funcion == Constantes_crearnuevoidparanuevaclavemensajes {

		//Llamamos a la BD para crear nuevi id clave para nuevo conjunto de clave para los mensajes
		idclavemensajes, test := bd.CrearNuevaClaveMensajesBD()
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_crearnuevoidparanuevaclavemensajes_err, MensajeSocket: "Error al crear id para nuevo conjunto de claves."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_crearnuevoidparanuevaclavemensajes_ok, Datos: []string{strconv.Itoa(idclavemensajes)}, MensajeSocket: "Id nuevo conjunto claves creado correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////////////////////////////////
	//NUEVA CLAVE USUARIO CON ID CONJUNTO CLAVES
	////////////////////////////////////////////
	if mensaje.Funcion == Constantes_nuevaclaveusuarioconidconjuntoclaves {

		idconjuntoclaves, _ := strconv.Atoi(mensaje.Datos[0])
		clavemensajes := mensaje.DatosClaves[0]
		clavesusuario := Clavesusuario{Idusuario: mensaje.Idfrom, Idclavesmensajes: idconjuntoclaves, Clavemensajes: clavemensajes}

		//Llamamos a la BD para guardar la nueva clave del usuario para ese conjunto de claves
		test := bd.GuardarClaveUsuarioMensajesBD(clavesusuario)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_nuevaclaveusuarioconidconjuntoclaves_err, MensajeSocket: "Error al asociar la clave del usuario con el id del conjunto de claves."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_nuevaclaveusuarioconidconjuntoclaves_ok, MensajeSocket: "Clave usuario asociada a id del conjunto de claves."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	///////////////////////////
	//MARCAR MENSAJE COMO LEIDO
	///////////////////////////
	if mensaje.Funcion == Constantes_marcarmensajeleido {
		i, _ := strconv.Atoi(mensaje.Datos[0])

		//Llamamos a la BD para marcar mensaje del usuario como leido
		test := bd.marcarLeidoPorUsuarioBD(i, mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_marcarmensajeleido_err, MensajeSocket: "Error al marcar mensaje como leído."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_marcarmensajeleido_ok, MensajeSocket: "ok marcar mensaje como leído."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////////////
	//MARCAR CHAT COMO LEIDO
	////////////////////////
	if mensaje.Funcion == Constantes_marcarchatcomoleido {

		//Llamamos a la BD para marcar chat como leido
		test := bd.marcarChatLeidoPorUsuarioBD(mensaje.Chat, mensaje.Idfrom)
		fmt.Println("mm<-->", mensaje.Chat, mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_marcarchatcomoleido_err, MensajeSocket: "Error al marcar chat como leído."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_marcarchatcomoleido_ok, MensajeSocket: "OK al marcar chat como leído."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	/////////////////////////
	//OBTENER CLAVES MENSAJES
	/////////////////////////
	if mensaje.Funcion == Constantes_getclavesmensajes {
		//Llamada a BD obtener claves de los mensajes
		claves, test := bd.getClavesMensajes(mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getclavesmensajes_err, MensajeSocket: "Error al obtener las claves de los mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Codigicamos con marshal
		datos := make([]string, 0, 1)
		for i := 0; i < len(claves); i++ {
			datos = append(datos, claves[i])
		}

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getclavesmensajes_ok, Datos: datos, MensajeSocket: "getclavesmensajes"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	if mensaje.Funcion == Constantes_getClavesDeUnUsuario {

		fmt.Println("getclaves de ", mensaje.Idfrom)
		claves, test := bd.getClavesMensajesdeUnUsuario(mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getClavesDeUnUsuario_err, MensajeSocket: "Error al obtener las claves de los mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Codigicamos con marshal
		datos := make([]string, 0, 1)
		datosClaves := make([][]byte, 0, 1)

		for i := 0; i < len(claves); i++ {
			b, _ := json.Marshal(claves[i])
			datos = append(datos, string(b))
			datosClaves = append(datosClaves, claves[i].Clavemensajes)
		}

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getClavesDeUnUsuario_ok, Datos: datos, DatosClaves: datosClaves, MensajeSocket: "getClavesDeUnUsuario"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	if mensaje.Funcion == Constantes_obtenerClavesDeMuchosUsuarios {
		usuarios := make([]Usuario, 0, 1)
		//fmt.Println("getclaves de ", mensaje.Idfrom)
		for i := 0; i < len(mensaje.Datos); i++ {
			usuario, _ := strconv.Atoi(mensaje.Datos[i])
			user, test := bd.obtenerClavePrivadaUsuario(usuario)
			if test == false {
				mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenerClavesDeMuchosUsuarios_err, MensajeSocket: "Error al obtener las claves de los mensajes."}
				EnviarMensajeSocketSocket(conexion, mesj)
				return
			}
			usuarios = append(usuarios, user)
		}

		//Codigicamos con marshal
		datos := make([]string, 0, 1)
		datosClaves := make([][]byte, 0, 1)

		for i := 0; i < len(usuarios); i++ {
			usuarios[i].Clavelogin = []byte{}
			usuarios[i].Claveprivrsa = []byte{}
			usuarios[i].Salt = []byte{}
			b, _ := json.Marshal(usuarios[i])
			datos = append(datos, string(b))
			datosClaves = append(datosClaves, usuarios[i].Clavepubrsa)
		}

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Datos: datos, Funcion: Constantes_obtenerClavesDeMuchosUsuarios_ok, DatosClaves: datosClaves, MensajeSocket: "getClavesDeMuchosUsuario"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////
	//EDITAR USUARIO
	////////////////
	if mensaje.Funcion == Constantes_modificarusuario {

		//Llamada a la BD para modificar usuario
		usuarioAux := Usuario{Id: mensaje.Idfrom, Nombre: mensaje.From, Estado: mensaje.Datos[0]}
		boolean := bd.modificarUsuarioNombreEstadoBD(usuarioAux)

		if boolean {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_modificarusuario_ok, Datos: []string{usuarioAux.Nombre, usuarioAux.Estado}, MensajeSocket: "usuariocambiaok"}
			EnviarMensajeSocketSocket(conexion, mesj)
		} else {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_modificarusuario_err, MensajeSocket: "Error al cambiar usuario"}
			EnviarMensajeSocketSocket(conexion, mesj)
		}

	}

	/////////////
	//EDITAR CHAT
	/////////////
	if mensaje.Funcion == Constantes_modificarchat {

		//Creamos el char con sus datos
		chat := Chat{Id: mensaje.Chat, Nombre: mensaje.Datos[0]}

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(mensaje.Idfrom, chat.Id)
		if permitido == false {
			//Enviamos mensaje error
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_modificarchat_err, MensajeSocket: "No tienes permiso para realizar esta acción, noperteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Llamamos a la BD para modificar al chat
		boolean := bd.modificarChatBD(chat)
		if boolean {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_modificarchat_ok, MensajeSocket: "chatcambiadook"}
			EnviarMensajeSocketSocket(conexion, mesj)
		} else {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_modificarchat_err, MensajeSocket: "Error al cambiar el chat"}
			EnviarMensajeSocketSocket(conexion, mesj)
		}
	}

	if mensaje.Funcion == Constantes_crearchat {

		var idusuarios = make([]int, 0, 1)
		idusuarios = append(idusuarios, mensaje.Idfrom)
		_, id := bd.crearChatBD(idusuarios, mensaje.Datos[0])
		datos := make([]string, 0, 1)
		datos = append(datos, strconv.Itoa(id))
		mesj := MensajeSocket{From: "usuario", Datos: datos, Funcion: Constantes_crearchat_ok, MensajeSocket: "Chatcreado:"}
		EnviarMensajeSocketSocket(conexion, mesj)

	}

	//////////////
	//GET USUARIOS
	//////////////
	if mensaje.Funcion == Constantes_getUsuarios {

		usuarios, test := bd.getUsuarios()
		if test {

			//Los convertimos con marshall
			datos := make([]string, 0, 1)
			for i := 0; i < len(usuarios); i++ {
				usu := Usuario{Id: usuarios[i].Id, Nombre: usuarios[i].Nombre, Estado: usuarios[i].Estado}
				b, _ := json.Marshal(usu)
				datos = append(datos, string(b))
			}

			mesj := MensajeSocket{From: mensaje.From, Datos: datos, Funcion: Constantes_getUsuarios_ok, MensajeSocket: "getusuariosok"}
			EnviarMensajeSocketSocket(conexion, mesj)
		} else {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getUsuarios_err, MensajeSocket: "Error al obtener usuarios."}
			EnviarMensajeSocketSocket(conexion, mesj)
		}
	}
}
