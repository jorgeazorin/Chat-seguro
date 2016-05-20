  var myApp = angular.module('myApp',[]);

  myApp.controller('controlador', ['$scope', function($scope) {

    //Inicializamos datos, mostramos htmls...
    $scope.greeting = 'Hola!';
  	$scope.mostarlogin = true;
    $scope.mostrarregistro = false;

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
        
        //Probamos a obtener chats
        ws.send("chats");
    };

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
          console.log("BIENNNNNNNNNNN")
          $scope.mostarlogin = false;

        }
      } else {
        alert(respuesta.MensajeSocket)
      }

    }

  }]);
