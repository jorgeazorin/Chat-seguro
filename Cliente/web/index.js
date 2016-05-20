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

      //AHORA RECIBES UN OBJETO Object BLOB
      //Ni idea de que es y como obtener sus cosillas por dentro
      //Cuando lo sepas ya cambias esto xD

    	console.log("Hemos recibido"+event.data);
      console.log(event.data);

      if(event.data == "loginok") {
        $scope.mostarlogin = false;
        alert('¡Usuario logeado correctamente!')        
      }

      if(event.data == "loginnook") {
        alert('Error al iniciar sesión, pruebe con otras credenciales.')
      }

      if(event.data == "registrook") {
        alert('¡Usuario registrado correctamente!')
        $scope.mostarlogin = false;
      }

      if(event.data == "registronook") {
        alert('¡Error al registrar usuario!')
      }

    }

  }]);
