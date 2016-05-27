package main

const (
	Constantes_registrarusuario     = 101
	Constantes_registrarusuario_ok  = 201
	Constantes_registrarusuario_err = 401

	Constantes_login     = 102
	Constantes_login_ok  = 202
	Constantes_login_err = 402

	Constantes_obtenermensajesAdmin     = 103
	Constantes_obtenermensajesAdmin_ok  = 203
	Constantes_obtenermensajesAdmin_err = 403

	Constantes_obtenerchats     = 104
	Constantes_obtenerchats_ok  = 204
	Constantes_obtenerchats_err = 404

	Constantes_enviar     = 105
	Constantes_enviar_ok  = 205
	Constantes_enviar_err = 405

	Constantes_obtenermensajeschat     = 106
	Constantes_obtenermensajeschat_ok  = 206
	Constantes_obtenermensajeschat_err = 406

	Constantes_agregarusuarioschat     = 107
	Constantes_agregarusuarioschat_ok  = 207
	Constantes_agregarusuarioschat_err = 407

	Constantes_eliminarusuarioschat     = 108
	Constantes_eliminarusuarioschat_ok  = 208
	Constantes_eliminarusuarioschat_err = 408

	Constantes_getclavepubusuario     = 109
	Constantes_getclavepubusuario_ok  = 209
	Constantes_getclavepubusuario_err = 409

	Constantes_getclavecifrarmensajechat     = 110
	Constantes_getclavecifrarmensajechat_ok  = 210
	Constantes_getclavecifrarmensajechat_err = 410

	Constantes_crearnuevoidparanuevaclavemensajes     = 111
	Constantes_crearnuevoidparanuevaclavemensajes_ok  = 211
	Constantes_crearnuevoidparanuevaclavemensajes_err = 411

	Constantes_nuevaclaveusuarioconidconjuntoclaves     = 112
	Constantes_nuevaclaveusuarioconidconjuntoclaves_ok  = 212
	Constantes_nuevaclaveusuarioconidconjuntoclaves_err = 412

	Constantes_marcarmensajeleido     = 113
	Constantes_marcarmensajeleido_ok  = 213
	Constantes_marcarmensajeleido_err = 413

	Constantes_marcarchatcomoleido     = 114
	Constantes_marcarchatcomoleido_ok  = 214
	Constantes_marcarchatcomoleido_err = 414

	Constantes_getclavesmensajes     = 115
	Constantes_getclavesmensajes_ok  = 215
	Constantes_getclavesmensajes_err = 415

	Constantes_getClavesDeUnUsuario     = 116
	Constantes_getClavesDeUnUsuario_ok  = 216
	Constantes_getClavesDeUnUsuario_err = 416

	Constantes_obtenerClavesDeMuchosUsuarios     = 117
	Constantes_obtenerClavesDeMuchosUsuarios_ok  = 217
	Constantes_obtenerClavesDeMuchosUsuarios_err = 417

	Constantes_modificarusuario     = 118
	Constantes_modificarusuario_ok  = 218
	Constantes_modificarusuario_err = 418

	Constantes_modificarchat     = 119
	Constantes_modificarchat_ok  = 219
	Constantes_modificarchat_err = 419

	Constantes_crearchat     = 120
	Constantes_crearchat_ok  = 220
	Constantes_crearchat_err = 420
)

type Clavesusuario struct {
	Idusuario        int    `json:"Idusuario"`
	Idclavesmensajes int    `json:"Idclavesmensajes"`
	Clavemensajes    []byte `json:"Clavemensajes"`
}

//Struct de los mensajes que se envian por el socket
type MensajeSocket struct {
	From        string   `json:"From"`
	Idfrom      int      `json:"Idfrom"`
	To          int      `json:"To"`
	Password    string   `json:"Password"`
	Funcion     int      `json:"Funcion"`
	Datos       []string `json:"Datos"`
	DatosClaves [][]byte `json:"DatosClaves"`
	Chat        int      `json:"Chat"`
	Mensaje     string   `json:"MensajeSocket"`
	Mensajechat []byte   `json:"Mensajechat"`
}

//Para pasar los datos de un usuario
type Usuario struct {
	Id               int    `json:"Id"`
	Nombre           string `json:"Nombre"`
	Clavepubrsa      []byte `json:"Clavepubrsa"`
	Claveprivrsa     []byte `json:"Claveprivrsa"`
	Claveenclaro     string `json:"Claveenclaro"`
	Clavehashcifrado []byte `json:"Clavehashcifrado"`
	Clavehashlogin   []byte `json:"Clavehashlogin"`
}

//Para obtener todos los datos de un mensaje
type MensajeTodo struct {
	Id           int    `json:"Id"`
	Texto        []byte `json:"Texto"`
	TextoClaro   string `json:"TextoClaro"`
	Emisor       int    `json:"Emisor"`
	Chat         int    `json:"Chat"`
	IdClave      int    `json:"IdClave"`
	NombreEmisor string `json:"NombreEmisor"`
	Clave        []byte `json:"Clave"`
}

//Para obtener los datos de un mensaje
type MensajeDatos struct {
	Mensaje MensajeTodo `json:"Mensaje"`
	Leido   bool        `json:"Leido"`
}

//Para obtener los datos del chat
type Chat struct {
	Id     int    `json:"Id"`
	Nombre string `json:"Nombre"`
}

//Para obtener todos los datos del chat
type ChatDatos struct {
	Chat          Chat           `json:"Chat"`
	MensajesDatos []MensajeDatos `json:"Mensajes"`
	Clave         []byte         `json:"Clave"`
	IdClave       int            `json:"IdClave"`
}

type MensajeAdmin struct {
	idclavesmensajes int
	Clave            []byte
}
