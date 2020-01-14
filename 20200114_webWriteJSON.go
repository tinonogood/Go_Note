package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Id      int
	Name    string
	Age     int
}

var allUsers = []User {
	{Id: 1,Name: "john", Age: 20},
	{Id: 2,Name: "mary", Age: 18},
}



func jsonHandler (writer http.ResponseWriter, request *http.Request) {
	js, err := json.Marshal(allUsers[0])
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
	fmt.Println(allUsers[0].Id)
}

func main() {
	for _, user := range allUsers {
		path := "/"
		path += strconv.Itoa(user.Id)
		http.HandleFunc(path, func(w http.ResponseWriter, req *http.Request){
			js, err := json.Marshal(allUsers[user.Id - 1 ])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		})
		fmt.Println(path)
	}
	http.HandleFunc("/json", jsonHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
