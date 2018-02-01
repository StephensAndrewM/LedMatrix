var SCREEN_WIDTH = 32;
var SCREEN_HEIGHT = 32;

var screenPixels = null;
var socket;

document.addEventListener('DOMContentLoaded', function() {
   
	var screen = document.getElementById('screen');
	screenPixels = [];
	for (var j = 0; j < SCREEN_HEIGHT; j++) {
		screenPixels[j] = [];
		for (var i = 0; i < SCREEN_WIDTH; i++) {
			var pixel = document.createElement('div');
			pixel.className += ' pixel';
			screen.appendChild(pixel);
			screenPixels[j][i] = pixel;		
		}
		screen.appendChild(document.createElement('br'));
	}

	socket = new WebSocket("ws://localhost:8000/ws");
	socket.onopen = function() {
	    console.log("Socket is open");
	};
	socket.onmessage = function (e) {
	    var display = JSON.parse(e.data);
	    for (var j = 0; j < display.Grid.length; j++) {
	    	for (var i = 0; i < display.Grid[j].length; i++) {
	    		var pixel = display.Grid[j][i];
	    		console.log(j, i, pixel);
	    		console.log(screenPixels[j][i]);
	    		var rgb = "rgb(" + pixel.R + "," + pixel.G + "," + pixel.B + ")";
	    		console.log(rgb);
	    		screenPixels[j][i].style.background = rgb;
	    	}
	    }

	}
	socket.onclose = function () {
	    console.log("Socket closed");
	}

}, false);