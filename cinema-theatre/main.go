package main

import (
	"encoding/json"
	
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

func bind(message string) (string, string, string) {
	parts := strings.Split(message, "/")
	log.Printf("Received message: %s", message)
	log.Printf("Split message into parts: %v", parts)
	if len(parts) != 3 {
		log.Println("Error: message is not in the correct format")
		return "", "", ""
	}

	mes := parts[0]
	seat := parts[1]
	row := parts[2]
	log.Printf("Extracted message: %s, row: %s, seat: %s", mes, row, seat)
	return mes, seat, row
}

func reserveSeat(seatNumber int, rowNumber int) (bool, string) {
	log.Printf("Attempting to reserve seat %d, row %d", seatNumber, rowNumber)
	if seat[rowNumber][seatNumber].booked {
		log.Println("Seat already booked")
		return false, "already booked"
	}

	seat[rowNumber][seatNumber].booked = true
	log.Printf("Seat reserved at row %d, seat %d", rowNumber, seatNumber)
	return true, "good"
}

func initseats() {
	log.Println("Initializing seats")
	seat = make([][]*seats, 10)
	for i := range seat {
		seat[i] = make([]*seats, 10)
		for j := 0; j < 10; j++ {
			seat[i][j] = &seats{i, j, false}
			log.Printf("Seat initialized at row %d, seat %d", i, j)
		}
	}
	log.Println("Seats initialization complete")
}
func getSeats()([]byte){
	var jsonSeat []byte
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			log.Printf("getSeats: row %d, seat %d", i, j)
			jsonSeat, err := json.Marshal(seat)
			if err != nil {
				log.Println(err)
				return nil
			}
			log.Printf("getSeats: %s", string(jsonSeat))
		}
		
	}
	return jsonSeat
}
func connectionHandler(w http.ResponseWriter, r *http.Request) {
	initseats()
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	log.Println("Connection established")
	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		p = treatmessages(string(p))
		log.Printf("Received message of type %d: %s", messageType, string(p))
		if err := c.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
		log.Println("Sent message")

	}
}
func treatmessages(message string) []byte {
	log.Printf("Treating message: %s", message)
	msg, seat, row := bind(message)
	if msg==""||seat==""||row==""{
		log.Println("Error: message is not in the correct format")
		return nil
	}
	sea, err := strconv.Atoi(seat)
	log.Printf("Converted seat to int: %d, err: %v", sea, err)
	rowi, err := strconv.Atoi(row)
	log.Printf("Converted row to int: %d, err: %v", rowi, err)

	if msg == "" {
		log.Println("Ignoring empty message")
		err := []byte("empty message")
		return err
	}
	if msg == "make reservation" {
		log.Println("Attempting to make reservation")
		ok, err := reserveSeat(sea, rowi)
		log.Printf("Reservation status: %t, err: %s", ok, err)
		if !ok {
			log.Println(err)
			errBytes := []byte(err)
			return errBytes
		}
		
		if err == "" {
			log.Println("Seat booked successfully")
			return nil
		}
		return []byte("u did it ")
	} else if msg == "get seats" {
		log.Println("Attempting to get seats")
		return getSeats()
	}

	return []byte("Unknown message")
}


func main() {
    log.Println("Initializing server...")

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        log.Println("Received new WebSocket connection request")
        connectionHandler(w, r)
    })

    log.Println("Server started on :8080, waiting for connections...")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatalf("ListenAndServe failed: %v", err)
    }
}
