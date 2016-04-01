package main

//Struct de usuario
type Usuario struct {
	id                  int
	nombre              string
	clavePublicaRsa     string
	claveUsuarioCifrada string
	claveusuario        string
}

//Funcion para obtener los datos del usuario cuando se loguea
func (usuario *Usuario) login(nombre string) bool {
	usuario.nombre = nombre
	user := getUsuarioBD(nombre)
	usuario.id = user.id
	usuario.nombre = user.nombre
	usuario.clavePublicaRsa = user.clavePublicaRsa
	usuario.claveUsuarioCifrada = user.claveUsuarioCifrada
	usuario.claveusuario = user.claveusuario
	usuario.nombre = user.nombre
	return true
}
