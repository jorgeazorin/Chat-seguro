  var myApp = angular.module('myApp',[]);

  myApp.controller('controlador', ['$scope', function($scope) {
    $scope.greeting = 'Hola!';
  	$scope.mostarlogin=true;

    var ws = new WebSocket("wss://localhost:10443/echo");
    $scope.Login = function() {
        $scope.greeting="Login";
        console.log("Mierda");  
        ws.send("hola");
      //  $scope.mostarlogin=false;
    };

	ws.onopen = function(){  
        console.log("Socket has been opened!");  
    };

    ws.onmessage = function (event) {
  		console.log("recibido"+event.data);
	}

  }]);
