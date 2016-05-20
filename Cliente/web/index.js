  var myApp = angular.module('myApp',[]);

  myApp.controller('controlador', ['$scope', function($scope) {

    //Inicializamos datos, mostramos htmls...
    $scope.greeting = 'Hola!';
  	$scope.mostarlogin=true;

    var ws = new WebSocket("wss://localhost:10443/echo");

    //Usuario se logea
    $scope.Login = function() {
        $scope.greeting="Login";
        
        usuario = {}
        usuario.Nombre = $scope.username
        usuario.Claveenclaro = $scope.password

        ws.send("login");
        ws.send(JSON.stringify(usuario));
        ws.send("prueba");
    };

    //Socket abierto, conexión establecida
    ws.onopen = function(){  
      console.log("Socket has been opened!");  
    };

    //Cliente servidor http nos envía algo
    ws.onmessage = function (event) {
    	console.log("Hemos recibido"+event.data);
    }

  }]);
