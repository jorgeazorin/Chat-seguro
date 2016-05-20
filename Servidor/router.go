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
	Funcion       string   `json:"Funcion"`
	Datos         []string `json:"Datos"`
	DatosClaves   [][]byte `json:"DatosClaves"`
	Chat          int      `json:"Chat"`
	MensajeSocket string   `json:"MensajeSocket"`
	Mensajechat   []byte   `json:"Mensajechat"`
}

type TodosLosDatos struct {
	Usuario Usuario  `json:"Usuario"`
	Claves  []string `json:"Claves"`
	Chats   []Chat   `json:"Chats"`
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
	if mensaje.Funcion == "registrarusuario" {

		var usuarionuevo Usuario

		usuarionuevo.Nombre = mensaje.Datos[0]
		usuarionuevo.Clavelogin = mensaje.DatosClaves[0]
		usuarionuevo.Clavepubrsa = mensaje.DatosClaves[1]
		usuarionuevo.Claveprivrsa = mensaje.DatosClaves[2]

		/*
			PERMISOS USUARIO?
			if usuario.Id != 1 {
				mesj := MensajeSocket{From: usuario.Nombre, MensajeSocket: "Error, no tienes permiso para registrar a un usuario."}
				EnviarMensajeSocketSocket(conexion, mesj)
				return
			}*/

		usuario, test := bd.insertUsuarioBD(usuarionuevo)

		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Registro incorrecto"}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: "DatosUsuario", Datos: []string{strconv.Itoa(usuario.Id), usuario.Nombre}, DatosClaves: [][]byte{usuario.Clavepubrsa, usuario.Claveprivrsa}, MensajeSocket: "Registrado correctamente"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////
	//INICIAR SESIÓN
	////////////////
	if mensaje.Funcion == "login" {

		//Rellenamos el usuario de la conexión con el login
		usuario, test := bd.loginUsuarioBD(mensaje.From, mensaje.DatosClaves[0])

		//Si login incorrecto se lo decimos al cliente
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Login incorrecto"}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Si es correcto, aAñadimos la conexion al map de conexiones
		conexiones[usuario.Id] = conexion

		//Enviamos un mensaje de todo OK al usuario logeado
		mesj := MensajeSocket{From: mensaje.From, Funcion: "DatosUsuario", Datos: []string{strconv.Itoa(usuario.Id), usuario.Nombre}, DatosClaves: [][]byte{usuario.Clavepubrsa, usuario.Claveprivrsa}, MensajeSocket: "Logeado correctamente"}
		EnviarMensajeSocketSocket(conexion, mesj)

	}

	////////////////////////
	//CREAR MENSAJE Y ENVIAR
	////////////////////////
	if mensaje.Funcion == "enviar" {

		//Guardamos los mensajes en la BD
		//m := Mensaje{Texto: mensaje.MensajeSocket, Chat: 1, Emisor: mensaje.Idfrom, Clave: 1}
		//bd.guardarMensajeBD(m)

		//Obtenemos los usuarios que pertenecen en el chat
		idusuarios, test := bd.getUsuariosChatBD(mensaje.Chat)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al enviar mensaje"}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos el mensaje a todos los usuarios de ese chat (incluido el emisor)
		for i := 0; i < len(idusuarios); i++ {
			conexion, ok := conexiones[idusuarios[i]]
			if ok {
				EnviarMensajeSocketSocket(conexion, mensaje)
			}
		}

	}

	/////////////////////////////
	//OBTENER MENSAJES DE UN CHAT
	/////////////////////////////
	if mensaje.Funcion == "obtenermensajeschat" {

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(mensaje.Idfrom, mensaje.Chat)
		if permitido == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "No perteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Obtenemos los mensajes
		mensajes, test := bd.getMensajesChatBD(mensaje.Chat, mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al obtener los mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Los convertimos con marshall
		datos := make([]string, 0, 1)
		for i := 0; i < len(mensajes); i++ {
			men := Mensaje{Id: mensajes[i].Mensaje.Id, Texto: mensajes[i].Mensaje.Texto}
			b, _ := json.Marshal(men)
			datos = append(datos, string(b))
		}

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Datos: datos, MensajeSocket: "Mensajes recibidos:"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////////////////
	//AGREGAMOS USUARIOS AL CHAT
	////////////////////////////
	if mensaje.Funcion == "agregarusuarioschat" {

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(usuario.Id, mensaje.Chat)
		if permitido == false {
			//Enviamos mensaje error
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "No tienes permiso para realizar esta acción, noperteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Todos los usuarios del mensaje
		idusuarios := make([]int, 0, 1)
		for i := 0; i < len(mensaje.Datos); i++ {
			idusuario, _ := strconv.Atoi(mensaje.Datos[i])
			idusuarios = append(idusuarios, idusuario)
		}

		//Los agregamos llamando a la BD
		test := bd.addUsuariosChatBD(mensaje.Chat, idusuarios)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Hubo un error al añadir los usuarios al chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Operación realizada correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	/////////////////////////////
	//ELIMINAMOS USUARIOS AL CHAT
	/////////////////////////////
	if mensaje.Funcion == "eliminarusuarioschat" {

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(usuario.Id, mensaje.Chat)
		if permitido == false {
			//Enviamos mensaje error
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "No tienes permiso para realizar esta acción, noperteneces al chat de estos mensajes."}
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
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Hubo un error al eliminar los usuarios del chat."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Operación realizada correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	//////////////////////////////////
	//OBTENER CLAVE PÚBLICA DE USUARIO
	//////////////////////////////////
	if mensaje.Funcion == "getclavepubusuario" {

		//Obtenemos id del usuario
		idusuario, _ := strconv.Atoi(mensaje.Datos[0])

		//Llamamos a la bd para obtener la clave pública de un usuario
		clavepub, test := bd.getClavePubUsuario(idusuario)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al obtener clave del usuario."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, DatosClaves: [][]byte{clavepub}, MensajeSocket: "Clave pub usuario obtenida correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////////////////////////////////
	//OBTENER CLAVE PARA CIFRAR MENSAJES DE CHAT
	////////////////////////////////////////////
	if mensaje.Funcion == "getclavecifrarmensajechat" {

		//Obtenemos id del mensaje
		idchat, _ := strconv.Atoi(mensaje.Datos[0])

		prueba := bd.usuarioEnChat(mensaje.Idfrom, idchat)
		if prueba == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "No tienes permiso de acceso a los datos de este mensaje."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Lamamos a la BD para obtener la última clave con la que cifrar los mensajes
		clavemensaje, test := bd.getLastKeyMensaje(idchat, mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al obtener clave para cifrar mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Funcion: "DatosClaveCifrarMensajeChat", DatosClaves: [][]byte{clavemensaje}, MensajeSocket: "Clave para mensajes obtenida correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	//////////////////////////////////
	//NUEVA CLAVE PARA CIFRAR MENSAJES
	//////////////////////////////////
	if mensaje.Funcion == "crearnuevoidparanuevaclavemensajes" {

		//Llamamos a la BD para crear nuevi id clave para nuevo conjunto de clave para los mensajes
		idclavemensajes, test := bd.CrearNuevaClaveMensajesBD()
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al crear id para nuevo conjunto de claves."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, Datos: []string{strconv.Itoa(idclavemensajes)}, MensajeSocket: "Id nuevo conjunto claves creado correctamente."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////////////////////////////////
	//NUEVA CLAVE USUARIO CON ID CONJUNTO CLAVES
	////////////////////////////////////////////
	if mensaje.Funcion == "nuevaclaveusuarioconidconjuntoclaves" {

		idconjuntoclaves, _ := strconv.Atoi(mensaje.Datos[0])
		clavemensajes := mensaje.DatosClaves[0]
		clavesusuario := Clavesusuario{Idusuario: mensaje.Idfrom, Idclavesmensajes: idconjuntoclaves, Clavemensajes: clavemensajes}

		//Llamamos a la BD para guardar la nueva clave del usuario para ese conjunto de claves
		test := bd.GuardarClaveUsuarioMensajesBD(clavesusuario)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al asociar la clave del usuario con el id del conjunto de claves."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Enviamos mensaje contestación
		mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Clave usuario asociada a id del conjunto de claves."}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	///////////////////
	//OBTENER LOS CHATS
	///////////////////
	if mensaje.Funcion == "obtenerchats" {

		//Llamada BD obtener chats del usuario
		chats, test := bd.getChatsUsuarioBD(mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al obtener los chats."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Codificamos con Marshal
		datos := make([]string, 0, 1)
		for i := 0; i < len(chats); i++ {
			b, _ := json.Marshal(chats[i])
			datos = append(datos, string(b))
		}

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: mensaje.From, Datos: datos, MensajeSocket: "Chats:"}
		EnviarMensajeSocketSocket(conexion, mesj)

	}

	///////////////////////////
	//MARCAR MENSAJE COMO LEIDO
	///////////////////////////
	if mensaje.Funcion == "marcarmensajeleido" {
		i, _ := strconv.Atoi(mensaje.Datos[0])

		//Llamamos a la BD para marcar mensaje del usuario como leido
		test := bd.marcarLeidoPorUsuarioBD(i, mensaje.Idfrom)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al marcar mensaje como leído."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}
	}

	/////////////////////////
	//OBTENER CLAVES MENSAJES
	/////////////////////////
	if mensaje.Funcion == "getclavesmensajes" {

		//Llamada a BD obtener claves de los mensajes
		claves, test := bd.getClavesMensajes(usuario.Id)
		if test == false {
			mesj := MensajeSocket{From: mensaje.From, MensajeSocket: "Error al obtener las claves de los mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Codigicamos con marshal
		datos := make([]string, 0, 1)
		for i := 0; i < len(claves); i++ {
			fmt.Println("::::", claves[i])
			datos = append(datos, claves[i])
		}

		//Enviamos los mensajes al usuario que los pidió
		mesj := MensajeSocket{From: usuario.Nombre, Datos: datos, MensajeSocket: "Claves:"}
		EnviarMensajeSocketSocket(conexion, mesj)
	}

	////////////////
	//EDITAR USUARIO
	////////////////
	if mensaje.Funcion == "modificarusuario" {

		//Llamada a la BD para modificar usuario
		usuarioAux := Usuario{Id: mensaje.Idfrom, Nombre: mensaje.From, Claveprivrsa: mensaje.DatosClaves[1], Clavepubrsa: mensaje.DatosClaves[2], Clavelogin: mensaje.DatosClaves[0]}
		boolean := bd.modificarUsuarioBD(usuarioAux)

		if boolean {
			mesj := MensajeSocket{From: usuario.Nombre, MensajeSocket: "Usuario cambiado correctamente"}
			EnviarMensajeSocketSocket(conexion, mesj)
		} else {
			mesj := MensajeSocket{From: usuario.Nombre, MensajeSocket: "Error al cambiar usuario"}
			EnviarMensajeSocketSocket(conexion, mesj)
		}

	}

	/////////////
	//EDITAR CHAT
	/////////////
	if mensaje.Funcion == "modificarchat" {

		i, err := strconv.Atoi(mensaje.Datos[0])
		if err != nil {
			mesj := MensajeSocket{From: usuario.Nombre, MensajeSocket: "Error con los parámetros recibidos."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Creamos el char con sus datos
		chat := Chat{Id: i, Nombre: mensaje.Datos[1]}

		//Comprobamos si ese usuario está en ese chat
		permitido := bd.usuarioEnChat(usuario.Id, chat.Id)
		if permitido == false {
			//Enviamos mensaje error
			mesj := MensajeSocket{From: usuario.Nombre, MensajeSocket: "No tienes permiso para realizar esta acción, noperteneces al chat de estos mensajes."}
			EnviarMensajeSocketSocket(conexion, mesj)
			return
		}

		//Llamamos a la BD para modificar al chat
		boolean := bd.modificarChatBD(chat)
		if boolean {
			mesj := MensajeSocket{From: usuario.Nombre, MensajeSocket: "Chat cambiado correctamente"}
			EnviarMensajeSocketSocket(conexion, mesj)
		} else {
			mesj := MensajeSocket{From: usuario.Nombre, MensajeSocket: "Error al cambiar el chat"}
			EnviarMensajeSocketSocket(conexion, mesj)
		}
	}

}
