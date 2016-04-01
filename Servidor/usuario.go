package main

//Struct de usuario
type Usuario struct {
	nombre              string
	clavePublicaRsa     string
	claveUsuarioCifrada string
}

//Funcion para obtener los datos del usuario cuando se loguea
func (usuario *Usuario) obtenerUsuario() {
	usuario.nombre = "jogre"
	usuario.clavePublicaRsa = "clavePublica"
	usuario.claveUsuarioCifrada = "claveCifrada"
}
