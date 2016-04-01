package main

func (conexion *Conexion) ProcesarMensajeSocket(mensaje MensajeSocket) {

	if mensaje.Funcion == "login" {
		conexion.usuario.obtenerUsuario()
		mesj := MensajeSocket{}
		mesj.From = conexion.usuario.nombre
		mesj.MensajeSocket = "Usuario online"
		for i := 0; i < len(conexion.conexiones.conexiones); i++ {
			conexion.conexiones.conexiones[i].EnviarMensajeSocketSocket(mesj)
		}
	}

	if mensaje.Funcion == "enviar" {
		mesj := MensajeSocket{}
		mesj.From = conexion.usuario.nombre
		mesj.MensajeSocket = mensaje.MensajeSocket
		mesj.Funcion = "enviar"
		for i := 0; i < len(conexion.conexiones.conexiones); i++ {
			if conexion.conexiones.conexiones[i].conexion != conexion.conexion {
				conexion.conexiones.conexiones[i].EnviarMensajeSocketSocket(mesj)
			}
		}
	}

}
