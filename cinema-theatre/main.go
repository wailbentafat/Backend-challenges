package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


var  upgrader = websocket.Upgrader {
	
	ReadBufferSize :1024,
	WriteBufferSize :1024,
}
func connectionHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()
	for {
		message,p,err:=c.ReadMessage()
		if err!=nil{
			fmt.Println(err)
			return
		}
		fmt.Println(string(message))
		if err:=c.WriteMessage(websocket.TextMessage,p);err!=nil{
			fmt.Println(err)
			return
		}

	}
}
func treatmessages(){
	
}


func main() {
    http.HandleFunc("/ws", connectionHandler) 
    log.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
