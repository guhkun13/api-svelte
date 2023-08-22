package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var msgChan chan string

func sendTime(w http.ResponseWriter) {
	fmt.Println("send time here!")
	for {
		time.Sleep(time.Second * 1)
		msg := time.Now().Format("15:04:05")
		msgChan <- msg
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	go sendTime(w)
	msgChan = make(chan string)

	defer func() {
		close(msgChan)
		msgChan = nil
		fmt.Println("Client close connection")
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("fail init http.Flusher")
	}

	for {
		select {
		case msg := <-msgChan:
			fmt.Println(msg)
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("client close connection")
			return
		}
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/event", sseHandler)

	fmt.Println("running server on port 3123!")
	log.Fatal(http.ListenAndServe("localhost:3123", router))
}
