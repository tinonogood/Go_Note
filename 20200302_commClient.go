package main

import (
	"encoding/gob"
	"log"
	"net"
)


const (
	serverHost = "localhost"
	serverPort = ":8080"
)


// ===================
// Serialize
// ===================

type User struct {
	Id      int
	Name    string
	Age     int
}

func updateUser() {
	serializeUserEncode User{Id: 1, Name:"john", Age: 20}
	log.Println("start client");

	conn, err := net.Dial("tcp", serverHost + serverPort)
	if err != nil {
		log.Fatal("Connection error", err)
	}

	encoder := gob.NewEncoder(conn)
	encoder.Encode(studentEncode)

	conn.Close()
	log.Println("Serialize User update done")
}
