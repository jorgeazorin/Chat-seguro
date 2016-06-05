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

	fmt.Println("Mensaje recibido: ", usuario.Id, "Funcion ", mensaje.Funcion)

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
		var test = false

		//Rellenamos el usuario de la conexión con el login
		usuario1, test := bd.loginUsuarioBD(mensaje.From, mensaje.DatosClaves[0])

		//Si login incorrecto se lo decimos al cliente
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_login_err, MensajeSocket: "Nombre de usuario o contraseña incorrectos."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
		usuario.Id = usuario1.Id
		//Si es correcto, aAñadimos la conexion al map de conexiones
		conexiones[usuario1.Id] = conexion

		//Enviamos un mensaje de todo OK al usuario logeado
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_login_ok, Datos: []string{strconv.Itoa(usuario1.Id), usuario1.Nombre, usuario1.Estado}, DatosClaves: [][]byte{usuario1.Clavepubrsa, usuario1.Claveprivrsa}, MensajeSocket: "Logeado correctamente"}
		EnviarMensajeSocketSocket(conexion, mesj)

	}

	////////////////////////
	//OBTENER MENSAJES ADMIN
	////////////////////////
	if mensaje.Funcion == Constantes_obtenermensajesAdmin {
		mensajes, test := bd.getMensajesAdmin(usuario.Id)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenermensajesAdmin_err, MensajeSocket: "Error obtenermensajesAdmin."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
		claves := make([][]byte, 0, 1)
		datos := make([]string, 0, 1)
		for i := 0; i < len(mensajes); i++ {
			men := mensajes[i]
			men.Mensaje.Texto = []byte{}
			b, _ := json.Marshal(men)
			datos = append(datos, string(b))
			claves = append(claves, mensajes[i].Mensaje.Texto)
		}
		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenermensajesAdmin_ok, Datos: datos, DatosClaves: claves, MensajeSocket: "obtenermensajesAdmin:"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	///////////////////
	//OBTENER LOS CHATS
	///////////////////
	if mensaje.Funcion == Constantes_obtenerchats {

		//Llamada BD obtener chats del usuario
		chats, test := bd.getChatsUsuarioBD(usuario.Id)
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
					if idusuarios[i] != usuario.Id {
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
		idTO := -1
		m := Mensaje{Texto: mensaje.Mensajechat, Chat: mensaje.Chat, Emisor: usuario.Id, Clave: idclavemensaje}
		if mensaje.From == "-200" {
			m.Admin = true
			idTO, _ = strconv.Atoi(mensaje.Datos[1])
		}

		bd.guardarMensajeBD(m, idTO)

		//Obtenemos los usuarios que pertenecen en el chat
		idusuarios, test := bd.getUsuariosChatBD(mensaje.Chat)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_enviar_err, MensajeSocket: "Error al enviar mensaje."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_enviar_ok, MensajeSocket: "Mensaje enviado."}
		EnviarMensajeSocketSocket(conexion, mesj)

		fmt.Println("Mensajes a clientes conectados")
		//Enviamos el mensaje a todos los usuarios de ese chat (incluido el emisor)

		for i := 0; i < len(idusuarios); i++ {
			conexion, ok := conexiones[idusuarios[i]]
			if ok {
				mensaje.To = idusuarios[i]
				mensaje.MensajeSocket = "Mensaje de otro usuario al chat:"
				fmt.Println("enviandomensaje a usuario ", idusuarios[i])
				if m.Admin {
					mensaje.Funcion = Constantes_MensajeAdminOtroClienteConectado

				} else {
					mensaje.Funcion = Constantes_MensajeOtroClienteConectado
				}
				EnviarMensajeSocketSocket(conexion, mensaje)
			}

		}

	}

	/////////////////////////////
	//OBTENER MENSAJES DE UN CHAT
	/////////////////////////////
	if mensaje.Funcion == Constantes_obtenermensajeschat {
		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(usuario.Id, mensaje.Chat)
		if permitido == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_obtenermensajeschat_err, MensajeSocket: "Error de permiso. No perteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Obtenemos los mensajes
		mensajes, test := bd.getMensajesChatBD(mensaje.Chat, usuario.Id)
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
		permitido := bd.usuarioEnChat(usuario.Id, mensaje.Chat)
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

		//Los añadimos al chat
		test := bd.addUsuariosChatBD(mensaje.Chat, nombresusuarios)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_agregarusuarioschat_err, MensajeSocket: "Hubo un error al añadir usuarios al chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
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

		//Todos los usuarios del mensaje
		nombresusuarios := make([]string, 0, 1)
		for i := 0; i < len(mensaje.Datos); i++ {
			nombresusuarios = append(nombresusuarios, mensaje.Datos[i])
		}

		//Los eliminamos llamando a la BD
		test := bd.removeUsuariosChatBD(mensaje.Chat, nombresusuarios)
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
		clavesusuario := Clavesusuario{Idusuario: usuario.Id, Idclavesmensajes: idconjuntoclaves, Clavemensajes: clavemensajes}

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
		test := bd.marcarLeidoPorUsuarioBD(i, usuario.Id)
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
		test := bd.marcarChatLeidoPorUsuarioBD(mensaje.Chat, usuario.Id)

		if test == false {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_marcarchatcomoleido_err, MensajeSocket: "Error al marcar chat como leído."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
		mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_marcarchatcomoleido_ok, MensajeSocket: "OK al marcar chat como leído."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////////////
	//OBTENER CLAVES USUARIO
	////////////////////////
	if mensaje.Funcion == Constantes_getClavesDeUnUsuario {

		claves, test := bd.getClavesMensajesdeUnUsuario(usuario.Id)
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

	////////////////////////////////
	//OBTENER CLAVES MUCHOS USUARIOS
	////////////////////////////////
	if mensaje.Funcion == Constantes_obtenerClavesDeMuchosUsuarios {

		usuarios := make([]Usuario, 0, 1)

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
		usuarioAux := Usuario{Id: usuario.Id, Nombre: mensaje.From, Estado: mensaje.Datos[0]}
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
		permitido := bd.usuarioEnChat(usuario.Id, chat.Id)
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

	////////////
	//CREAR CHAT
	////////////
	if mensaje.Funcion == Constantes_crearchat {

		var idusuarios = make([]int, 0, 1)
		idusuarios = append(idusuarios, usuario.Id)
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

	/////////////////////////
	//GET USUARIOS DE UN CHAT
	/////////////////////////
	if mensaje.Funcion == Constantes_getUsuariosDeUnChat {
		idChat, _ := strconv.Atoi(mensaje.Datos[0])
		usuarios, test := bd.usuariosEnChat(idChat)
		if test {

			//Los convertimos con marshall
			datos := make([]string, 0, 1)
			for i := 0; i < len(usuarios); i++ {
				datos = append(datos, strconv.Itoa(usuarios[i]))
			}

			mesj := MensajeSocket{From: mensaje.From, Datos: datos, Funcion: Constantes_getUsuariosDeUnChat_ok, MensajeSocket: "getusuariosDeUnChatok"}
			EnviarMensajeSocketSocket(conexion, mesj)
		} else {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_getUsuariosDeUnChat_err, MensajeSocket: "Error al obtener usuarios."}
			EnviarMensajeSocketSocket(conexion, mesj)
		}
	}

	////////////////////////////
	//ASOCIAR NUEVA CLAVE A CHAT
	////////////////////////////
	if mensaje.Funcion == Constantes_AsociarNuevaClaveAChat {
		idChat, _ := strconv.Atoi(mensaje.Datos[0])
		idClave, _ := strconv.Atoi(mensaje.Datos[1])
		test := bd.AsociarNuevaClaveAChat(idChat, idClave)
		if test {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_AsociarNuevaClaveAChat_ok, MensajeSocket: "Ok aosciar nueva clave chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
		} else {
			mesj := MensajeSocket{From: mensaje.From, Funcion: Constantes_AsociarNuevaClaveAChat_err, MensajeSocket: "Error asociar nueva clave chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
		}

	}

}
