var SCREEN_WIDTH = 128;
var SCREEN_HEIGHT = 32;
var LED_SIZE = 8;
var LED_SPACE = 2;
var MIN_BRIGHTNESS = 40;
var SOCKET_ADDR = "ws://localhost:8000/ws";

var canvas;
var screenPixels = null;
var socket = null;
var refreshInterval = null;

document.addEventListener('DOMContentLoaded', function() {

	initScreen();
	initSocket();

}, false);

var initScreen = function() {
	canvas = document.getElementById('screen');
    canvas.width = SCREEN_WIDTH * (LED_SIZE + (2 * LED_SPACE));
    canvas.height = SCREEN_HEIGHT * (LED_SIZE + (2 * LED_SPACE));
}

var initSocket = function() {
	socket = new WebSocket(SOCKET_ADDR);
	socket.onopen = function() {
	    console.log("Socket is open");
	    if (refreshInterval != null) {
	    	console.log("Socket reconnected.")
	    	clearTimeout(refreshInterval);
	    	refreshInterval = null;
	    }
	};
	socket.onmessage = function (e) {
		var context = canvas.getContext('2d');
		context.clearRect(0, 0, canvas.width, canvas.height);

	    var display = JSON.parse(e.data);
	    for (var j = 0; j < display.length; j++) {
	    	for (var i = 0; i < display[j].length; i++) {
	    		var pixel = display[j][i];
	    		var r = Math.max(pixel[0], MIN_BRIGHTNESS);
	    		var g = Math.max(pixel[1], MIN_BRIGHTNESS);
	    		var b = Math.max(pixel[2], MIN_BRIGHTNESS);
	    		var rgb = "rgb(" + r + "," + g + "," + b + ")";
	    		
	    		var centerX = getCenter(i);
	    		var centerY = getCenter(j);
	    		context.beginPath();
	    		context.arc(centerX, centerY, (LED_SIZE / 2), 0, 2 * Math.PI, false);
	    		context.fillStyle = rgb;
	    		if (r+g+b > 100) {
		    		context.shadowColor = rgb;
		    		context.shadowBlur = 15;
		    		context.shadowOffsetX = 0;
		    		context.shadowOffsetY = 0;
		    	} else {
		    		context.shadowBlur = 0;
		    	}
	    		context.fill();
	    	}
	    }

	}
	socket.onclose = function () {
		console.log("Socket closed");
	    if (refreshInterval != null) {
	    	return;
	    }
	    refreshInterval = window.setInterval(function() {
	    	initSocket();
	    }, 500)
	}
}

var getCenter = function(i) {
	return (i * ((LED_SPACE * 2) + LED_SIZE)) + LED_SPACE + (LED_SIZE / 2);
}