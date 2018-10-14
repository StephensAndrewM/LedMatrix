package main

import(
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket"
    "image"
)

type WebDisplay struct {
	Conn *websocket.Conn
}

func NewWebDisplay() *WebDisplay {
	d := new(WebDisplay)
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

func (d *WebDisplay) Redraw(img *image.RGBA) {
    width := img.Bounds().Dx()
    height:= img.Bounds().Dy()

    // Convert to 2D array of 3-tuples
    data := make([][][]int, height)
    for j := 0; j < height; j++ {
        data[j] = make([][]int, width)
        for i := 0; i < width; i++ {
            data[j][i] = make([]int, 3)
            rgba := img.RGBAAt(i,j)
            data[j][i][0] = int(rgba.R)
            data[j][i][1] = int(rgba.G)
            data[j][i][2] = int(rgba.B)
        }
    }

	if d.Conn != nil {
		json, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("JSON Error: %s\n", err)
			return
		}
		d.Conn.WriteMessage(websocket.TextMessage, json)
	}
}