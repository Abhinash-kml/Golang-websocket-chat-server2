package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"v/models"
	database "v/pkg"
)

func main() {
	port := "8000"       /* os.Getenv("PORT") */
	database.ConnectDB() // Connect to postgress
	serveraddres := fmt.Sprintf(":%s", port)
	chatserver := models.NewChatServer()

	fmt.Printf("----- Starting server on localhost:%s -----\n", port)

	http.HandleFunc("/ws", chatserver.HandleConnections)
	go http.ListenAndServe(serveraddres, nil)
	go chatserver.HandleMessages()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	<-sigs

	database.DB.Close()
	fmt.Println("----- Shutting down server -----")
}
