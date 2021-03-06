var myApp = angular.module('myApp',[]);

//Controlador Angular JS
myApp.controller('controlador', ['$scope', function($scope) {

  //Inicializamos datos, mostramos htmls
  $scope.mostarlogin = true;
  $scope.mostrarregistro = false;
  $scope.editarnombrechat = false;
  $scope.chatactual='chat';
  $scope.editarchat = false;
  $scope.verdatosusuario = false;
  $scope.vericonoperfil = false;
  $scope.datosusuarioeditar = false;
  $scope.modificarDatosUsuarioValue = false;
  $scope.verlistausuarios = false;
  var ws = new WebSocket("wss://localhost:"+location.port+"/echo");

  //Usuario se registra
  $scope.Registro = function() {
      $scope.greeting="Registro";
      
      usuario = {}
      usuario.Nombre = $scope.username
      usuario.Claveenclaro = $scope.password

      ws.send("registro@/@"+JSON.stringify(usuario));
  };

  //Usuario se logea
  $scope.Login = function() {
      usuario = {}
      usuario.Nombre = $scope.username
      usuario.Claveenclaro = $scope.password
      ws.send("login@/@"+JSON.stringify(usuario));
  };

  //Ver todos los mensajes del chat
  $scope.verChat = function(id) {

    for(i=0;i<chats.length;i++) {

      if(chats[i].Chat.Id == id) {
        $scope.mensajes = chats[i].Mensajes;
        $scope.chatactual=chats[i].Chat.Nombre;
        $scope.idchatactual=chats[i].Chat.Id;
        $scope.clavechatactual=chats[i].Chat.UltimaClave;

        //Llamamos a marcar como leidos
        mensaje = {}
        mensaje.Chat = $scope.idchatactual

        ws.send("leidos@/@"+JSON.stringify(mensaje));

        //Lo dejamos como leido
        $scope.chats[i].Chat.Leido = true
        $scope.$apply()
      }        
    }

    //Mostramos que se puede editar
    if(id != undefined)
      $scope.editarchat = true;

    $scope.$apply()
  }

  //Enviando mensaje por el chat
  $scope.enviarMensaje = function() {
    //Si es vacio no se hace nada
    if($scope.textoaenviar == "" || $scope.textoaenviar == undefined) {
      return
    }

    mensaje = {}
    mensaje.Chat = $scope.idchatactual
    mensaje.MensajeSocket = $scope.textoaenviar

    ws.send("enviarmensaje@/@"+JSON.stringify(mensaje))
    $scope.textoaenviar = ""
  }

  //Add usuarios al chat
  $scope.addUsuario = function() {
      mensaje = {}
      mensaje.MensajeSocket = $scope.usuarioadd
      mensaje.Chat = $scope.idchatactual

      $scope.usuarioadd = ""
      ws.send("addusuariochat@/@"+JSON.stringify(mensaje));
  }

  //Remove usuarios del chat
  $scope.removeUsuario = function() {        
      mensaje = {}
      mensaje.MensajeSocket = $scope.usuarioremove
      mensaje.Chat = $scope.idchatactual

      $scope.usuarioremove = ""
      ws.send("removeusuariochat@/@"+JSON.stringify(mensaje));
  }

  //Editar nombre del chat
  $scope.editarChat = function() {
    
    //Modo editar
    if($scope.editarnombrechat == false) {
      $scope.editarnombrechat = true;
    } 
    //Modo normal y guardar lo editado
    else {
      chat = {}
      chat.Nombre = $scope.nuevonombrechat
      chat.Id = $scope.idchatactual
      chat.UltimaClave = $scope.clavechatactual

      ws.send("editarchat@/@"+JSON.stringify(chat))
      $scope.editarnombrechat = false;
      $scope.$apply()
    }

  }

  //Editar nombre del chat
  $scope.crearChat = function() {    
    datos = "Nuevo chat"

    ws.send("crearchat@/@"+datos)
  }

  //Vemos el perfil de un usuario
  $scope.verDatosUsuario = function(nombre, estado) {

    //Es perfil de usuario
    if(nombre == $scope.username || nombre == undefined) {
      $scope.botonmodificardatosusuario = true
      $scope.usuariousername = $scope.username
      $scope.usuarioestadousuario = $scope.estadousuario
    } else {
      $scope.botonmodificardatosusuario = false
      $scope.usuariousername = nombre
      $scope.usuarioestadousuario = estado
    }

    if($scope.usuarioestadousuario==undefined || $scope.usuarioestadousuario=="")
      $scope.usuarioestadousuario = "Sin estado."

    //Modo editar
    if($scope.verdatosusuario == false || nombre != undefined) {
      $scope.verdatosusuario = true;
    } 
    //Modo normal y guardar lo editado
    else {
      $scope.verdatosusuario = false;
    }

    $scope.$apply()
  }

  //Modificamos el perfil del propio usuario
  $scope.modificarDatosUsuario = function() {
    //Modo editar
    if($scope.modificarDatosUsuarioValue == false) {

      datosusuariomodonoeditar.className = "oculto"
      datosusuariomodoeditar.className = ""
      $scope.modificarDatosUsuarioValue = true;
    } 
    //Modo normal y guardar lo editado
    else {
      usuario = {}
      usuario.Nombre = $scope.nuevonombreusuario
      usuario.Estado = $scope.nuevoestadousuario

      if($scope.nuevonombreusuario == "" || $scope.nuevonombreusuario == undefined)
        usuario.Nombre = $scope.username

      ws.send("editarusuario@/@"+JSON.stringify(usuario))

      datosusuariomodonoeditar.className = ""
      datosusuariomodoeditar.className = "oculto"
      $scope.modificarDatosUsuarioValue = false;
      $scope.$apply()
    }
  }


  //Vemos los usuarios si se busca algo si no, los chats del usuario
  $scope.verchatsousuarios = function() {

    if($scope.verlistausuarios == false) {
      divverchats.className = "oculto"
      divverusuarios.className = ""
      $scope.verlistausuarios = true
      $scope.placeholderbusqueda = "Buscador de usuarios"
    } else {
      divverchats.className = ""
      divverusuarios.className = "oculto"
      $scope.verlistausuarios = false
      $scope.placeholderbusqueda = "Buscador de chats"
    }

    $scope.$apply()
  }

  //Vemos si hay mensajes sin leer
  function versiestanleidos() {
    for(i=0;i<$scope.chats.length;i++) {
      $scope.chats[i].Chat.Leido = true
      $scope.chats[i].numsinleer = 0
      
      //Para un chat      
      for(j=0;j<$scope.chats[i].Mensajes.length;j++) {
        if($scope.chats[i].Mensajes[j].Leido == false && $scope.chats[i].Mensajes[j].Mensaje.Emisor != $scope.idusuario) {
          {
            $scope.chats[i].Chat.Leido = false
            $scope.chats[i].numsinleer ++;
          }
        }
      }
    } 
  }


  //Socket abierto, conexión establecida
  ws.onopen = function(){  
    console.log("Socket has been opened!");  
  };

  /////////////////////////////////////////////
  //Cuando cliente servidor http nos envía algo
  /////////////////////////////////////////////
  ws.onmessage = function (event) {

    respuesta = JSON.parse(event.data)

    //Datos usuario
    if(respuesta.MensajeSocket == "DatosUsuario") {

      if(respuesta.Datos.length != 0) {

        //Ontenemos datos usuario
        usuario = JSON.parse(respuesta.Datos)
        $scope.idusuario = usuario.Id
        $scope.username = usuario.Nombre
        $scope.estadousuario = usuario.Estado
        $scope.vericonoperfil = true;

        //Pedimos chats
        ws.send("chats@/@")

        //Pedimos usuarios
        ws.send("getusuarios@/@");
      }
    }

    //Obtenemos chats al principio
    if(respuesta.MensajeSocket == "chats") {
      $scope.mostarlogin = false;
      chats = JSON.parse(respuesta.Datos[0])
      $scope.chats = chats
      $scope.verChat($scope.idchatactual)
      versiestanleidos()
      var div = document.getElementById('centro');
      div.scrollTop = '999999';
      $scope.$apply()
    }

    if(respuesta.MensajeSocket=="Error"){
      alert(respuesta.Datos[0])
    }

    //Pedimos chats si algo cambia
    if (respuesta.MensajeSocket == "mensajeenviado:" || respuesta.MensajeSocket == "chatcambiadook" || respuesta.MensajeSocket == "chatcreadook" || respuesta.MensajeSocket == "mensajedeotrocliente") {
      ws.send("chats@/@")
    }
    if(respuesta.MensajeSocket=="mensajeadmindeotrocliente"){
       usuario = {}
      usuario.Nombre = $scope.username
      usuario.Claveenclaro = $scope.password
      ws.send("login@/@"+JSON.stringify(usuario));
    }

    //Cuando usuario es modificado
    if(respuesta.MensajeSocket == "usuariocambiaok") {
      $scope.username = respuesta.Datos[0]
      $scope.estadousuario = respuesta.Datos[1]
      $scope.usuariousername = $scope.username
      $scope.usuarioestadousuario = $scope.estadousuario
      $scope.$apply()
    }

    //Obtenemos los usuarios
    if(respuesta.MensajeSocket == "getusuariosok") {
      
      if(respuesta.Datos.length != 0) {
        datos = eval(respuesta.Datos)
        usuarios = eval(datos[0])
        $scope.usuarios = usuarios  
      }
    }

    /////////
    //Alertas
    /////////
    if (true) {
      console.log(respuesta.MensajeSocket + ", " + respuesta.Funcion)
    }

  }

}]);
