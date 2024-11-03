package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)


var  upgrader = websocket.Upgrader {
	
	ReadBufferSize :1024,
	WriteBufferSize :1024,
}
type  seats struct {
	row int
	seat int
	booked bool
}

var seat [][] *seats

func bind (message string)(string,string,string){
	parts:=strings.Split(message, "/")
	mes := parts[0]
	seat:=parts[1]
	row:=parts[2]
	fmt.Println(mes,seat,row)
	return mes,seat,row
}

func reserveSeat(seatNumber int, rowNumber int) (bool, string) {
	if seat[rowNumber][seatNumber].booked {
		return false, "already booked"
	}

	seat[rowNumber][seatNumber].booked = true
	return true, "booked"
}

func initseats(){
	
	seat =make([][]*seats,10)
	for i:= range seat{
		seat[i]=make([]*seats,10)
		for j:=0 ;j< 10;j++{
	    
	    seat[i][j]=&seats{i,j,false}
		}
	
	}
}
func getSeats()([]byte){
	var jsonSeat []byte
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
		jsonSeat, err := json.Marshal(seat)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		fmt.Println(string(jsonSeat))
		}
		
	}
	return jsonSeat
}
func connectionHandler(w http.ResponseWriter, r *http.Request) {
	initseats()
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
func treatmessages(message string){
	msg,seat,row:=bind(message)
	sea,err:=strconv.Atoi(seat)
	print(err)
	rowi,err:=strconv.Atoi(row)
	print(err)

	if  msg==""{
		return
	}
	if msg=="make reservation"{
	 ok,err:=reserveSeat(sea,rowi)
	 if !ok{
	 	fmt.Println(err)
	 	return
	 }
	 if err==""{
	 	fmt.Println("booked successfully")
	 }

     

	}else if msg=="get seats"{
		getSeats()
	}


	}



func main() {
    http.HandleFunc("/ws", connectionHandler) 
    log.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
