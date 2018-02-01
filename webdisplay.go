package main

import(
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type WebDisplay struct {
	Conn *websocket.Conn
}

func NewWebDisplay() *WebDisplay {
	d := new(WebDisplay)
	fmt.Println("Initialized web display")
	return d
}

func (d *WebDisplay) Initialize() {
	upgrader := websocket.Upgrader{}

	http.Handle("/", http.FileServer(http.Dir("public_html/")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		fmt.Println("Saved socket!")
		if err != nil {
			fmt.Println("Socket Upgrade Error:", err)
			return
		}
		d.Conn = c
	})
	go http.ListenAndServe(":8000", nil)
	fmt.Println("Serving HTTP traffic...")
}

func (d *WebDisplay) Redraw(g *PixelGrid) {
	fmt.Println("Redraw")
	fmt.Println(g)
	if d.Conn != nil {
		json, err := json.Marshal(g)
		if err != nil {
			fmt.Sprintln("JSON Error: %s", err)
			return
		}
		fmt.Println(json)
		d.Conn.WriteMessage(websocket.TextMessage, json)
	}
}