  var myApp = angular.module('myApp',[]);

  myApp.controller('controlador', ['$scope', function($scope) {

    //Inicializamos datos, mostramos htmls...
    $scope.greeting = 'Hola!';
  	$scope.mostarlogin = true;
    $scope.mostrarregistro = false;
    $scope.chatactual='chat';
    var ws = new WebSocket("wss://localhost:10443/echo");

    //Usuario se registra
        $scope.Registro = function() {
        $scope.greeting="Registro";
        
        usuario = {}
        usuario.Nombre = $scope.username
        usuario.Claveenclaro = $scope.password

        ws.send("registro");
        ws.send(JSON.stringify(usuario));
    };

    //Usuario se logea
    $scope.Login = function() {
        $scope.greeting="Login";
        
        usuario = {}
        usuario.Nombre = $scope.username
        usuario.Claveenclaro = $scope.password

        ws.send("login");
        ws.send(JSON.stringify(usuario));
    };

    //Ver todos los mensajes del chat
    $scope.verChat = function(id) {
      //Recorremos chats buscando el seleccionado
      for(i=0;i<chats.length;i++) {

        if(chats[i].Chat.Id == id) {
          $scope.mensajes = chats[i].Mensajes;
          $scope.chatactual=chats[i].Chat.Nombre;
          $scope.idchatactual=chats[i].Chat.Id;
          console.log(chats[i].Mensajes)
        }        
      }
      $scope.$apply()
    }

    //Enviando mensaje por el chat
    $scope.enviarMensaje = function() {
      mensaje = {}
      mensaje.Chat = $scope.idchatactual
      mensaje.MensajeSocket = $scope.textoaenviar
      ws.send("enviarmensaje")
      ws.send(JSON.stringify(mensaje))
      $scope.textoaenviar = ""
      $scope.$apply()
    }

    //Socket abierto, conexión establecida
    ws.onopen = function(){  
      console.log("Socket has been opened!");  
    };

    //Cliente servidor http nos envía algo
    ws.onmessage = function (event) {

      respuesta = JSON.parse(event.data)

      if(respuesta.MensajeSocket == "Chats:") {
        console.log(respuesta);

        if(respuesta.Datos.length != 0) {
          
          $scope.mostarlogin = false;
          chats = eval(respuesta.Datos)

          for(i=0;i<chats.length;i++) {
            chats[i] = JSON.parse(chats[i])
          }

          $scope.chats = chats          
        }
        $scope.verChat($scope.idchatactual)
        $scope.$apply()
      }

      //Cuando el usuario se rellene se piden los chats
      else if(respuesta.Funcion == "DatosUsuario") {
        ws.send("chats")
      }

      else if(respuesta.MensajeSocket == "MensajeEnviado:") {
        ws.send("chats")
      }

      else {
        alert(respuesta.MensajeSocket)
      }

    }

  }]);
